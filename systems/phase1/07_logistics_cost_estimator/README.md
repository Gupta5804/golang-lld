# LLD Study Notes: Logistics Shipping Cost Estimator

**Goal (Initial Prompt):** We need a new shipping cost estimator for our logistics network. The operations team wants the ability to dynamically switch out how we calculate freight prices based on the time of year (e.g., holiday surge vs. standard fuel-saver). Additionally, we recently acquired a regional courier and leadership insists we pass a portion of our estimates (Last-Mile) through their legacy, proprietary calculation engine. Their system is a complete black box, notoriously slow, and requires entirely different data inputs (Imperial units) than our internal metric systems.

**Phase:** 1 (Foundation)  
**Core Patterns:** Strategy (Behavioral), Adapter (Structural)

---

## 1. The 4 Pillars (Requirement Gathering)

Before writing any code, we defined the strict boundaries of the system to prevent race conditions and protect our domain from external pollution.

* **Scope & Boundaries:**
  * **Caller-Driven Routing:** The Estimator is "dumb" regarding seasons. The caller dynamically injects the chosen `PricingStrategy` at runtime.
  * **Mandatory Delegation:** The Legacy system is strictly utilized for the "Last-Mile" cost for all shipments.
* **Business Logic & Edge Cases:**
  * **Formulas:** * Standard: `(Weight * Distance * 5) / 100`
    * Holiday: `(Weight * Distance * 8) / 100`
  * **The Total:** Final Cost = `LineHaulCost` (Internal Strategy) + `LastMileCost` (Legacy System).
* **Entity & Domain Modeling:**
  * **Money Precision:** All costs are returned as a structured `Quote` object representing amounts in **cents** (`int64`) to prevent floating-point rounding drift.
  * **The Input Constraint:** Our internal `Shipment` entity uses Metric (KG, KM). The Legacy API expects Imperial (Lbs, Miles) and returns float dollars.
* **System Constraints:**
  * **Thread Safety:** The `EstimatorService` will be instantiated as a singleton and shared across thousands of concurrent HTTP requests. It must be completely stateless.
  * **Fail-Safe Timeouts:** The legacy system is slow. We must strictly enforce a 2-second SLA via `context.Context` to prevent our own threads from hanging.

---

## 2. Architectural Design (The Calculator Analogy)

We utilized the **Strategy Pattern** for the dynamic line-haul calculations and the **Adapter Pattern** to build an Anti-Corruption Layer against the legacy black box.

**The Analogy:**

* **The Calculator (Estimator):** Sits on the desk. It doesn't remember previous math problems (stateless).
* **The Formula Sheets (Strategies):** Swappable instructions handed to the calculator for different equations.
* **The Translator (Adapter):** Takes our Metric numbers, converts them to Imperial, calls the foreign branch, and translates the foreign currency back to our cents.

**Directory Blueprint:**

```text
07_logistics_cost_estimator/
├── domain/
│   └── model.go          <-- Pure entities (Shipment, Quote). Zero dependencies.
├── estimator/
│   ├── estimator.go      <-- Core engine and Ports (PricingStrategy, RateFetcher interfaces)
│   └── estimator_test.go <-- Black-box TDD & Concurrency/Chaos proofs
├── legacy/
│   └── adapter.go        <-- Anti-Corruption Layer (Adapter Pattern & External timeout wrapper)
├── strategies/
│   └── strategies.go     <-- Concrete implementations of PricingStrategy
└── main.go               <-- Composition Root & Wiring
```

---

## 3. Data Structure & Contract Decisions (Staff-Level Choices)

* **Stateless Singleton:** Rather than storing the `Shipment` and the active `PricingStrategy` as fields on the `EstimatorService` struct, we passed them strictly as method arguments. This allowed 10,000 goroutines to call the engine simultaneously without triggering a data race.
* **Idiomatic Go Interfaces:** Refrained from using the "Java-style" `I` prefix (e.g., `IPricingStrategy`). We relied on implicit interface satisfaction and named our ports based on behavior (e.g., `RateFetcher`).
* **Dependency Inversion:** The core `EstimatorService` depends entirely on abstract interfaces (`PricingStrategy`, `RateFetcher`). It has absolutely no knowledge of the concrete `HolidayStrategy` or the `LegacyAdapter`.

---

## 4. Go-Specific Traps Avoided (The "Gotcha Ledger")

This system presented severe concurrency and memory-leak traps. Here is how we neutralized them:

### Trap 1: The Stateful Web Service Race Condition

* **The Trap:** Storing per-request data (like `strategy`) inside the `EstimatorService` struct. Under high load, Goroutine A's "Standard" strategy would be overwritten by Goroutine B's "Holiday" strategy microseconds before the calculation, charging the wrong price silently.
* **The Fix:** Structs should only hold long-lived dependencies (like the `RateFetcher` adapter). Ephemeral data must live in the method signature (`CalculateTotal(ctx, shipment, strategy)`) so it stays on the isolated goroutine stack.

### Trap 2: Early Float Truncation (Revenue Loss)

* **The Trap:** `cost := int64(weightKG) * int64(distanceKM) * 5`. Casting a float like `10.9 kg` to `int64` directly truncates it to `10`, losing a full kilogram of revenue.
* **The Fix:** Perform all multiplications in `float64` first, and cast to `int64` at the very last moment: `int64(weight * float64(distance) * 5.0)`.

### Trap 3: The Ignored Context

* **The Trap:** Using `case <-time.After(2 * time.Seconds):` inside the Adapter. This hardcodes a timeout but completely ignores the caller's `ctx`, failing to respect upstream cancellations.
* **The Fix:** We wrapped the parent context with a deadline (`context.WithTimeout(ctx, 2*time.Second)`) and listened to `<-ctx.Done()` in the select block.

### Trap 4: The Zombie Goroutine Memory Leak (The Deadliest Trap)

* **The Trap:** Wrapping a synchronous black-box API inside a goroutine and sending the result to an **unbuffered** channel (`make(chan float32)`). If the `select` block times out and exits, the goroutine remains permanently blocked trying to send its value to a reader that no longer exists. 10,000 timeouts = 10,000 leaked goroutines.
* **The Fix:** Initialized the channel with a buffer of 1: `make(chan float32, 1)`. If the main thread times out and leaves, the worker goroutine drops its payload into the buffer and peacefully dies, allowing the GC to clean it up.

---

## 5. Table-Driven TDD (The Red Phase)

We mathematically proved the system's resilience before implementing the math.

* **The Stateful Spy Mock:** Created a `SpyRateFetcher` equipped with its own `sync.Mutex` to safely record interactions (`CallCount`, `LastWeight`) even under heavy concurrent load.
* **Asserting Anatomy:** The table tests did not just assert the final Total; they explicitly asserted the `LineHaulCost` and `LastMileCost` independently to ensure the internal routing math was perfectly executed.
* **The Chaos Concurrency Test:** Deployed a `sync.WaitGroup` to fire 1,000 concurrent goroutines at a single `EstimatorService` instance, alternating strategies randomly. Running `go test -race` mathematically proved our stateless architecture prevented data races.

---

## 6. Quick Start

To verify the test suite (and the race detector) or run the end-to-end Composition Root simulation:

```bash
# Run the black-box test suite to prove the math and concurrency safety
go test -v -race ./...

# Execute the Composition Root to see the patterns at work
go run main.go
```
