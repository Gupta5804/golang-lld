# LLD Study Notes: Multi-Vendor Payment Gateway Router

**Goal(Initial prompt):** Design an in-memory payment routing engine for a global marketplace. When a user submits a checkout request, the system must pass it through a series of validation and fraud-check steps before routing it to the appropriate external third-party payment processor based on the region. We need a seamless way to integrate multiple external, incompatible payment APIs while keeping our core checkout pipeline clean and easily extensible for future steps.

---

## 1. The 4 Pillars (Requirement Gathering)

Before writing any code, we defined the strict boundaries of the system to prevent over-engineering and establish clear domain constraints.

* **Scope & Boundaries:**
  * **Core Domain Only:** The engine is purely a synchronous domain service exposing `ProcessPayment(ctx, req)`. It is not responsible for the HTTP/gRPC transport layer, JSON parsing, or idempotency keys (which are handled upstream).
  * **Static Pipeline:** The validation and fraud pipeline is statically assembled at application startup, not dynamically reconfigured per request.
* **Business Logic:**
  * **Strict Short-Circuiting:** If any step in the pipeline (e.g., Validation or Fraud) fails, the chain immediately halts, returning a typed sentinel error. There are no retries for domain-level business rule violations.
  * **Explicit Routing:** If a region is unmapped or missing, the system strictly fails rather than guessing a default gateway.
* **Entity & Domain Modeling:**
  * **Pure Entities:** The `PaymentRequest` struct contains only business-relevant fields (TransactionID, UserID, Amount, Region, ClientIP, DeviceFingerprint). It is isolated from vendor-specific infrastructure pollution (like Stripe JSON tags).
  * **Precision (No Floats):** All monetary values are represented as `int64` (cents) to completely eliminate floating-point arithmetic drift.
* **System Constraints:**
  * **Stateless Handlers:** The engine is built once and shared across thousands of concurrent goroutines. Pipeline steps must be 100% stateless to guarantee thread safety.
  * **Strict SLAs:** External vendor network calls are unpredictable. The entire execution is bounded by a strict `context.Context` timeout (e.g., 5 seconds) propagated through the entire chain.

---

## 2. Architectural Design (The Factory Analogy)

We utilized the **Chain of Responsibility (CoR)** for the workflow and the **Adapter** pattern for external integrations, organized strictly by Domain and Behavior.

**Directory Blueprint:**

```text
04_payment_router/
├── domain/
│   ├── models.go       <-- Pure structs (PaymentRequest)
│   └── errors.go       <-- Custom sentinel errors (e.g., ErrFraudDetected)
├── pipeline/
│   ├── handler.go      <-- The CoR interface and BaseHandler plumbing
│   ├── validate.go     <-- Step 1: Basic validation
│   ├── fraud.go        <-- Step 2: Fraud detection
│   └── router.go       <-- Step 3: Terminal region router
├── gateways/
│   ├── gateway.go      <-- Unified internal PaymentGateway interface
│   └── stripe.go       <-- External Adapter (Anti-Corruption Layer)
└── main.go             <-- Composition Root
```

### A. `domain` (The Core Entities)

* Contains pure structs and sentinel errors. Imports absolutely nothing. This guarantees the business logic is never polluted by external frameworks.

### B. `pipeline` (The Conveyor Belt)

* Manages the Chain of Responsibility. Handlers act as workers on the assembly line inspecting the `PaymentRequest`. If a check fails, the item is removed from the belt.

### C. `gateways` (The Shipping Fleet)

* The Anti-Corruption Layer (Adapter pattern). Wraps chaotic, vendor-specific APIs into a single, predictable `PaymentGateway` interface.

### D. The Composition Root (`main.go`)

* The only file that knows about all packages. It constructs the pipeline, injects the adapter registry into the router, and wires the chain together via `SetNext()`.

---

## 3. Data Structure & Contract Decisions (Staff-Level Choices)

* **Dependency Injection Registry:** Instead of hardcoding vendor logic inside the router, `RouterHandler` accepts a `map[string]gateways.PaymentGateway` at startup. This enables $O(1)$ routing lookups and makes testing trivial by injecting mock gateways.
* **Sentinel Errors:** Extracted all error definitions (e.g., `var ErrFraudDetected = errors.New("fraud detected")`) into the `domain` package. This creates a strongly-typed universal vocabulary, allowing clients to use `errors.Is()` instead of brittle string matching.
* **Context Propagation:** Added `context.Context` to every interface method signature (`Execute`, `Pay`) from day one to enforce network boundaries and cancellation propagation.

---

## 4. Go-Specific Traps Avoided

### Trap 1: The Interface Segregation Paradox

* **The Trap:** A CoR pipeline's `Execute` method must accept the master `*domain.PaymentRequest` object to move it down the chain. However, a step like Fraud Detection only needs the IP and Fingerprint, violating the "Need to Know" principle if it accesses the raw struct.
* **The Fix:** We separated the CoR *plumbing* from the *business logic*. `Execute` receives the full struct, but delegates the actual work to a private `performCheck(f FraudCheckable)` method. The domain entity inherently satisfies this narrow interface, decoupling the algorithm from the master struct.

### Trap 2: White-Box Testing Brittleness

* **The Trap:** Writing tests in `package pipeline` allows access to unexported methods (like `executeNext`), leading to tests that verify internal implementation details rather than external behavior.
* **The Fix:** Strictly used `package pipeline_test` (Black-Box testing). We built a self-contained `SpyHandler` that independently satisfied the public `Handler` interface, proving the API was usable from an external caller's perspective.

### Trap 3: Breaking the Conveyor Belt

* **The Trap:** When a CoR handler successfully completes its specific check, simply returning `nil` terminates the request prematurely, preventing downstream handlers from executing.
* **The Fix:** Handlers must explicitly return `h.executeNext(ctx, r)` upon success to intentionally pass control to the next link in the chain.

### Trap 4: The "Else" Trap (Line of Sight)

* **The Trap:** Using `else` blocks for error handling creates nested, hard-to-read code.
* **The Fix:** Enforced the Go idiom of early returns. Handled errors immediately and returned, keeping the "happy path" aligned strictly to the left margin.

---

## 5. Table-Driven TDD (The Red Phase)

We mathematically proved the system's behavior using Table-Driven Tests prior to implementation.

* **Short-Circuit Verification:** Tested the CoR mechanics by assembling chains of `SpyHandlers` equipped with `ShouldFail` flags. We verified that an upstream validation failure correctly bypassed the fraud and routing steps.
* **Component Rules:** Tested isolated handlers (`validate.go`, `fraud.go`) against our domain rules, asserting strictly against our `domain.Err...` sentinel variables using `errors.Is()`.
* **Dependency Mocking:** For the router, we injected a fake `map[string]gateways.PaymentGateway` loaded with a `MockGateway` to verify that supported regions triggered the `Pay` method and unsupported regions correctly halted execution.

---

## 6. Quick Start

To verify the test suite and run the end-to-end composition root simulation:

```bash
# Run the test suite to prove the business logic and short-circuiting
go test -v ./...

# Execute the Composition Root to see the end-to-end pipeline
go run main.go
```
