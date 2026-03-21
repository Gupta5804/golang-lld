# LLD Study Notes: Stock Price Alerting System

**Goal:** Design an in-memory, highly concurrent stock alerting system where users can subscribe to specific stock tickers and receive asynchronous notifications when the price crosses a customized threshold.

---

## 1. The 4 Pillars (Requirement Gathering)

Before writing any code, we defined the strict boundaries of the system to prevent systemic rewrites later.

* **Scope & Boundaries:**
  * **Dynamic Lifecycles:** Users can subscribe and unsubscribe at runtime.
  * **Directional Triggers:** We need to know if a user wants to be alerted when a stock goes *Above* or *Below* a target.
  * **Extensibility:** The engine must not know *how* notifications are sent (SMS, Email, etc.).
* **Business Logic:**
  * **One-Shot Alerts:** To prevent "flapping" (notification spam when a stock hovers around a threshold), alerts are strictly one-shot. Once it fires, it is deleted.
  * **Inclusive Thresholds:** In financial systems, thresholds are inclusive (`>=` or `<=`), not strictly greater/less than.
* **Entity & Domain Modeling:**
  * **Precision (No Floats):** All monetary values (prices, target thresholds) must be represented as `int64` (cents) to prevent floating-point rounding errors during threshold evaluation.
  * **Identifiers:** The engine relies on deterministic routing, so Users and Tickers are identified by unique IDs (UUID/int for users, strings like "AAPL" for tickers) rather than mutable names.
  * **State Snapshotting:** The dispatched `Command` must be an immutable snapshot. It must carry the `TriggeredPrice` and `Timestamp` from the exact millisecond the threshold was crossed, ensuring the async worker has the correct context even if the stock price fluctuates a second later.
* **System Constraints:**
  * **Backpressure (Load Shedding):** We cannot block the core market feed if the SMS API is slow. If the async queue fills up, we drop the alert (load shed) rather than freezing the engine.
  * **Lock Contention:** We cannot use a single global mutex for the whole system, and we cannot iterate over thousands of users $O(N)$ on every price tick.

---

## 2. Architectural Design (The Factory Analogy)

**Directory Blueprint:**

```text
03_stock_alerts/
├── domain/
│   └── model.go          <-- Interfaces & pure structs (Zero dependencies)
├── feed/
│   ├── engine.go         <-- The Observer (Writes to channel)
│   └── engine_test.go
├── notifier/
│   ├── worker.go         <-- The Command Dispatcher (Reads from channel)
│   └── worker_test.go
└── main.go               <-- Composition Root (Wires them together)
```

We avoided a "Big Ball of Mud" by isolating the system into three distinct packages:

### A. `domain` (The Vocabulary)

* Contains pure, immutable structs (`User`, `Ticker`, `AlertCondition`).
* **Zero Dependencies:** This package imports nothing. This prevents import cycles and ensures data safely crosses boundaries.
* **The Snapshot:** When an alert fires, we capture an `AlertSnapshot` (the exact price and timestamp at that millisecond) so the background worker has accurate context, even if the stock price changes a second later.

### B. `feed` (The Intake Desk / Observer Pattern)

* **The Engine:** Synchronously ingests price ticks and evaluates them against active subscriptions.
* It instantiates a `Command` object and drops it onto a Go channel (the conveyor belt) without waiting for the execution.

### C. `notifier` (The Delivery Fleet / Command Pattern)

* **DispatcherHub:** Manages a bounded worker pool of goroutines.
* These workers listen to the channel, pick up `Command` objects, and execute them. They have zero knowledge of what a "Ticker" or a "Feed" is.

### D. The Composition Root (`main.go`)

* **Dependency Injection:** The `feed` and `notifier` packages have absolutely no knowledge of each other. They do not import each other. In `main.go`, we create the shared Go channel (`make(chan domain.Command)`) and strictly inject it into both constructors (`NewEngine` and `NewDispatcherHub`). This allows us to easily swap out the in-memory channel for a real message broker (like Kafka or RabbitMQ) in the future without changing a single line of core business logic.


---


## 3. Data Structure Decisions (Staff-Level Choices)

* **Avoiding Floating Point Drift:** All monetary values (prices, thresholds) are stored as `int64` (cents). `float64` causes rounding errors during exact comparisons.
* **Idiomatic Enums:** Used custom types and `iota` for `Direction` (`Above`, `Below`) instead of raw strings or magic numbers.
* **$O(1)$ Lookups:** Instead of a `map[string][]AlertCondition` (which requires $O(N)$ iteration to unsubscribe), we used a nested map: `map[TickerSymbol]map[AlertID]AlertCondition`. This gives $O(1)$ ticker lookups AND $O(1)$ alert deletions.

---

## 4. Go-Specific Traps Avoided

This is the most critical section for revision. These are the concurrency traps we actively neutralized:

### Trap 1: The Nil Map Panic

* **The Trap:** Go maps are `nil` by default. Writing to them causes a panic.
* **The Fix:** In `NewEngine`, initialize the outer map: `make(map[string]map[string]domain.AlertCondition)`. In `Subscribe`, check if the inner map exists for a ticker; if not, `make` it before inserting.

### Trap 2: Concurrent Map Writes

* **The Trap:** If a user subscribes while the engine is evaluating a price tick, the program crashes instantly. Maps are not thread-safe.
* **The Fix:** Wrapped map mutations and reads inside `Subscribe`, `Unsubscribe`, and `UpdatePrice` with a `sync.RWMutex`.

### Trap 3: Safe Map Deletion during Iteration

* **The Trap:** In Java/C#, deleting from a collection while iterating over it causes a `ConcurrentModificationException`.
* **The Fix (Go Magic):** Go explicitly allows calling `delete(map, key)` during a `for...range` loop over that exact same map. We used this to easily implement the One-Shot constraint.

### Trap 4: Deadlocking the Channel (Load Shedding)

* **The Trap:** If the worker pool is slow, `commandChan <- cmd` will block the `UpdatePrice` goroutine forever, freezing the entire market feed.
* **The Fix:** Used a `select` block with a `default` case.
  
  ```go
  select {
  case e.commandChan <- cmd: // Success
  default: // Channel buffer is full. Drop the alert to save the system.
  }

### Trap 5: The Worker Pool Deadlock (The close Fix)

* **The Trap:** A `for cmd := range channel` loop never exits unless the channel is closed. If we call `wg.Wait()` on shutdown, but the workers are still looping on an open (but empty) channel, the system deadlocks.

* **The Fix:** In main.go, we must call close(commandChan) to signal to the workers that no more jobs are coming. This breaks their range loop, allowing them to call wg.Done(), which finally unblocks wg.Wait().

## 5. Table-Driven TDD (The Red Phase)

We proved the logic worked before we wrote it.

* **Unit Testing Asynchronous Channels:**
  To test the `Engine` without the real `DispatcherHub` stealing our messages, the test pretended to be the worker. We used a select block with a time.After(10 * time.Millisecond) timeout to safely assert whether a command was dropped into the channel or if the channel remained empty (preventing test deadlocks).

* **Testing the Worker Pool (The Spy Mock):**
To test the `DispatcherHub` without the real `Engine`, we injected `MockCommand` objects. The mock contained an `ExecuteFunc` closure that incremented a thread-safe `atomic.Int32` counter, mathematically proving the goroutines were actually picking up and executing the jobs.

## 6. The Graceful Shutdown Mechanics

A production system cannot just `exit(0)`. We used `sync.WaitGroup` to build a synchronization barrier:

  1. `main` boots the `DispatcherHub` and immediately calls `defer dispatcher.Stop()`.
  2. Inside `Start()`, we call `wg.Add(1)` outside the goroutine, and `defer wg.Done()` inside it.
  3. When `main` finishes its simulation, we `close(commandChan)`.
  4. `main` exits, triggering the deferred `dispatcher.Stop()`.
  5. `Stop()` calls `wg.Wait()`. It blocks until every worker finishes its current `Execute()` and exits its `range` loop. Zero dropped messages.

Quick Start
To see the system run end-to-end, proving the integration of the Engine and the Dispatcher:

```go
go run main.go
