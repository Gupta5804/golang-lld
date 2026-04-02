# LLD Study Notes: Marketing Analytics Report Hub

**Goal(Interview Prompt):** We need a unified reporting feature for our marketing team because right now they are wasting hours manually cross-referencing data from our user CRM, the legacy payment gateway, and the email delivery platform. I just want them to be able to ask the system for a "Weekly Engagement", "Monthly Revenue", or "Quarterly Churn" report and get it immediately without caring about the underlying data sources. Build a system that figures out which report they need and handles all the complex data gathering behind the scenes.

(*Design a reporting engine that abstracts the complexity of multiple underlying data systems. The engine must accept a simple request for a specific report type, dynamically generate the correct execution strategy, query only the necessary subsystems, and return a unified domain object. The system must be robust against external failures and maintain strict architectural boundaries.*)

---

## 1. The 4 Pillars (Requirement Gathering)

Before writing the contracts, we established strict system constraints to ensure loose coupling, fail-fast fault tolerance, and stateless concurrency.

* **Scope & Boundaries:**
  * **Core Domain Library:** The system is an internal library package, not an HTTP server or file formatter. It outputs a structured `domain.Report` object; presentation rendering (PDF/JSON) is explicitly out of scope.
  * **Interface-Driven Subsystems:** The actual CRM, Payment, and Email systems are external to the reporting domain and are mocked via interfaces.
* **Business Logic & Edge Cases:**
  * **Fixed-Format Aggregations:** Reports are bound strictly to predefined temporal enums (Weekly, Monthly, Quarterly).
  * **Precise Routing:** A Weekly report must query the CRM and Email platforms, but strictly ignore the Billing system.
  * **Fail-Fast Fault Tolerance:** If an underlying subsystem (e.g., the CRM) errors out during a report generation, the entire operation halts and returns an error immediately. No partial data is returned.
* **Entity & Domain Modeling:**
  * **Pure Domain Objects:** `domain.ReportRequest` (holding the type enum) and `domain.Report` (holding a string Title and a `map[string]string` Payload) remain entirely isolated from implementation logic.
* **System Constraints:**
  * **Stateless Concurrency:** The central Hub acts as a singleton-like orchestrator. It must be completely stateless to handle concurrent report generation safely across thousands of goroutines.

---

## 2. Architectural Design (The Concierge Analogy)

We utilized the **Facade Pattern** to hide the chaotic complexity of the underlying APIs, and the **Simple Factory & Strategy Patterns** to dynamically route report generation.

**The Analogy:**

* **The Concierge (Facade/Hub):** The client walks up to the desk and asks for a "Weekly Report." The concierge doesn't build the report, but they hold the phone numbers (Interfaces) for all the departments.
* **The Dispatcher (Simple Factory):** The concierge asks the dispatcher *who* handles Weekly Reports.
* **The Specialist (Strategy/Generator):** The dispatcher assigns a specialist (`weeklyGenerator`). The concierge hands the specialist the phones, the specialist calls the exact departments needed, writes the report, and hands it back to the concierge.

**Directory Blueprint:**

```text
06_marketing_analytics_hub/
├── domain/
│   └── model.go        <-- Pure structs and Enums (ReportType, Report)
├── crm/
│   └── crm.go          <-- Concrete external subsystem (Dummy adapter)
├── billing/
│   └── billing.go      <-- Concrete external subsystem (Dummy adapter)
├── email/
│   └── email.go        <-- Concrete external subsystem (Dummy adapter)
├── report/
│   ├── hub.go          <-- The Facade and Consumer-Defined Interfaces
│   ├── factory.go      <-- The Simple Factory logic
│   ├── strategies.go   <-- Private concrete generators (Business Logic)
│   └── hub_test.go     <-- Black-box Spy Mock TDD suite
└── main.go             <-- Composition Root & Integration Proof
```

---

## 3. Data Structure & Contract Decisions (Staff-Level Choices)

* **Consumer Rules (Dependency Inversion):** The `report` package defines the `CRM`, `Email`, and `Billing` interfaces *it* requires. It does not import the concrete `crm` or `billing` packages. This inversion guarantees that the `report` package is entirely immune to breaking changes made by the external subsystem teams.
* **File Segmentation (Avoiding the God File):** Instead of stuffing the entire architecture into `report.go`, we strictly separated concerns. `hub.go` manages external connections (Orchestration), `factory.go` manages struct initialization (Creational), and `strategies.go` manages the actual data extraction (Behavioral).
* **Enum-Driven Safety:** Utilizing Go's `iota` for the `domain.ReportType` prevents raw string typos and enforces type safety at compile time.

---

## 4. Go-Specific Traps Avoided (The "Gotcha Ledger")

### Trap 1: The Tight Coupling Anti-Pattern

* **The Trap:** Defining the `ReportHub` struct using pointers to concrete external services (e.g., `crm *crm.CRMService`). This makes mocking impossible and tightly couples the domain to the infrastructure.
* **The Fix:** Swapped concrete pointers for interface types (`crm CRM`) defined within the consumer package itself.

### Trap 2: The Stateful Facade (Concurrency Panic)

* **The Trap:** Storing the `generator ReportGenerator` as a field on the `ReportHub` struct. If multiple requests hit the Hub simultaneously, they would overwrite each other's generator state, leading to catastrophic race conditions.
* **The Fix:** Removed the generator from the struct. The Hub requests a new, ephemeral generator from the Factory *inside* the `GenerateReport` method scope for every single call.

### Trap 3: The Nil Map Panic

* **The Trap:** Attempting to assign keys to the `domain.Report` payload (`payload["active_users"] = "100"`) without initializing the map. Go maps are nil by default, and writing to a nil map triggers a runtime panic.
* **The Fix:** Explicitly utilized `payload := make(map[string]string)` inside the generator strategies before formatting and injecting the extracted data.

### Trap 4: Silent Data Loss

* **The Trap:** Using the blank identifier (`_`) to ignore the actual integer/float values returned from the external services just to satisfy the compiler and the "Red Phase" error assertions.
* **The Fix:** Captured the variables, utilized `fmt.Sprintf` to convert integers and floats to strings, and hydrated the `domain.Report` payload completely.

---

## 5. Table-Driven TDD (The Red Phase)

We mathematically proved the Facade's routing logic and fault tolerance using a heavily instrumented **Spy Mock** matrix.

* **The Spy Mocks:** Created `MockCRM`, `MockEmail`, and `MockBilling` that explicitly tracked call counts (e.g., `CountGetActiveUsersCalled`).
* **Routing Verification:** Our table-driven tests proved that requesting a `WeeklyEngagement` report resulted in exactly 1 CRM call, 1 Email call, and strictly **0** Billing calls.
* **Subsystem Failure Testing:** We added a `ShouldFail` boolean to `MockCRM`. We proved that if the CRM returns an error, the `ReportHub` respects the Fail-Fast constraint, immediately bubbling the error up to the client without calling subsequent systems.

---

## 6. Quick Start

To verify the test suite and run the end-to-end orchestration simulation:

```bash
# Run the black-box test suite to prove routing and error bubbling
go test -v ./report/...

# Execute the Composition Root to see the Facade orchestrate the subsystems
go run main.go
```
