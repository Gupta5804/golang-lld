# LLD Study Notes: Admin Task Execution Engine

**Goal(Interview Prompt):** We need a new toolkit for our customer support admins to manage user accounts. They want the ability to trigger single maintenance actions, like resetting a password, or trigger a whole bundle of actions at once, like a complete account wipe and notification blast. We also need to be able to line up these actions, execute them sequentially, and potentially reverse them if an admin makes a mistake.

(*Design an in-memory task execution toolkit for customer support admins. The system must provide a unified interface to trigger atomic administrative actions (e.g., suspending a user) or bundled composite actions (e.g., resetting a password and sending a notification simultaneously). The engine must sequentially execute these tasks, maintain a history of completed operations, and guarantee a flawless, reverse-order rollback if a bundled operation fails midway or an admin requests an undo.*)

---

## 1. The 4 Pillars (Requirement Gathering)

Before writing the contracts, we established strict system constraints to ensure thread safety and architectural decoupling.

* **Scope & Boundaries:**
  * **Unified Execution:** The API exposes a single `Execute(task)` endpoint. The caller (e.g., an HTTP handler) is completely blind to whether it is executing a single command or a massive tree of bundled actions.
  * **Storage Agnostic:** While the MVP operates in-memory, the domain and task layers are completely decoupled from how entities are persisted.
* **Business Logic & Edge Cases:**
  * **Partial Failure Rollbacks:** If a bundled action containing 5 tasks fails on task #3, the system must immediately halt and trigger `Undo()` on tasks #2 and #1 in strict reverse order (LIFO).
  * **Supported Actions:** MVP supports Account Suspension, Password Resets, and System Notifications.
* **Entity & Domain Modeling:**
  * **Thread-Safe Entities:** The `domain.User` entity governs its own state and concurrency protection, exposing behavior-rich methods rather than bare fields.
* **System Constraints:**
  * **High Concurrency:** Multiple admins may attempt to modify the same user or interact with the engine simultaneously. Thread-safety (via tight mutex locking) is mandatory at both the engine and domain levels.

---

## 2. Architectural Design (The Restaurant Kitchen Analogy)

We utilized the **Command Pattern** to encapsulate execution logic and the **Composite Pattern** to treat individual and grouped actions uniformly.

**The Analogy:**

* **The Chef (Engine):** Takes an order, executes it, and spikes the finished ticket.
* **The Order Ticket (Command):** Isolated instructions (e.g., "Reset Password").
* **The Stapled Tickets (Composite Bundle):** A batch of tickets. The Chef treats the staple exactly like a single ticket, unaware of the internal orchestration until execution begins.

**Directory Blueprint:**

```text
05_admin_task_engine/
├── domain/
│   └── models.go       <-- Pure entities with internal RWMutexes
├── task/
│   ├── command.go      <-- Task interface, Atomic Commands, Composite Bundle, and Repositories
│   └── command_test.go <-- Black-box Spy Mock TDD suite
├── engine/
│   └── engine.go       <-- The stateful Invoker (Execution Engine)
└── main.go             <-- Composition Root & Integration Proof
```

---

## 3. Data Structure & Contract Decisions (Staff-Level Choices)

* **The Single Interface Contract:** Instead of maintaining separate interfaces for commands and composite nodes, everything converges on a single `Task` interface requiring only `Execute() error` and `Undo() error`. This is the purest manifestation of the Composite pattern.
* **Persistent-Ready Repositories:** Rather than passing raw memory pointers (`*domain.User`) into commands, tasks are instantiated with a `userID` string and a `UserRepository` interface. This abstracts the data layer, allowing the system to easily migrate from an in-memory map to a PostgreSQL database without altering a single line of business logic.
* **Strict LIFO Iteration:** Rollbacks are mathematically structured using reverse iteration loops (`for j := i - 1; j >= 0; j--`) to guarantee dependent operations are unraveled properly.

---

## 4. Go-Specific Traps Avoided (The "Gotcha Ledger")

### Trap 1: The "Double Interface" Anti-Pattern

* **The Trap:** Creating an unnecessary wrapper struct (e.g., `SingleCommand`) to force an atomic command to satisfy a composite component interface.
* **The Fix:** Simplified the design so both the Leaf (`SuspendUserTask`) and the Composite (`BundleTask`) directly satisfy the exact same `Task` interface.

### Trap 2: Engine Deadlocks on Undo

* **The Trap:** Holding the Engine's `sync.RWMutex` lock for the entirety of the `UndoLast()` function. If the underlying task's `Undo()` method takes 5 seconds to query a database, the entire engine is paralyzed.
* **The Fix:** Minimized lock scoping. Locked the engine, popped the task off the history slice, immediately unlocked the engine, and *then* performed the heavy `lastTask.Undo()` operation outside the critical section.

### Trap 3: The Off-By-One Rollback Corruption

* **The Trap:** When catching a failure at index `i` during a bundle execution, starting the rollback loop at index `i` instead of `i - 1`. This inadvertently attempts to undo a task that never successfully completed, potentially corrupting database state.
* **The Fix:** Implemented precise index targeting, guaranteeing only fully successful tasks trigger their compensation logic.

### Trap 4: CI/CD Nightmare Assertions

* **The Trap:** Writing test assertions like `t.Errorf("Expected execute to be called")` inside loops, making it impossible to identify which mock failed in headless CI environments.
* **The Fix:** Enforced hyper-specific, indexed error logging (`t.Errorf("Task %d: Expected Execute...", i)`) and utilized `errors.Is()` for future-proof error value comparison.

---

## 5. Table-Driven TDD (The Red Phase)

We mathematically proved the Composite pattern's orchestration logic in complete isolation using **Spy Mocks**.

* **The Spy Mock:** Created a fake `MockTask` that recorded boolean flags for `ExecuteCalled` and `UndoCalled`, and allowed programmatic failure injection via a `ShouldFail` flag.
* **Black-Box Testing:** Placed tests in `package task_test` to enforce the use of constructor functions (`NewBundleTask`), proving the API works flawlessly from the perspective of an external caller.
* **Matrix Scenarios:**
  * **Happy Path:** Verified all tasks in a bundle execute and none trigger undo.
  * **Partial Failure:** Injected a failure in the middle of a bundle to mathematically prove execution halted and previously successful tasks were rolled back in reverse order.
  * **Explicit Undo:** Proved that a successfully executed bundle iterates backwards through its entire slice when manually asked to reverse via the Engine.

---

## 6. Quick Start

To verify the test suite and run the end-to-end orchestration simulation:

```bash
# Run the black-box test suite to prove the Composite rollback logic
go test -v ./task/...

# Execute the Composition Root to see the Engine orchestrate and rollback tasks
go run main.go
```
