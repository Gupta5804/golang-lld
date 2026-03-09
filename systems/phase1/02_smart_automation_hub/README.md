# Smart Home Automation Hub

**Phase:** 1 (Foundation)  
**Core Patterns:** Composite (Structural), Observer / Mediator (Behavioral)

## 1. Initial Problem Statement

Design the core domain model for a smart home hub that manages a heterogeneous collection of connected devices like lights, thermostats, and security cameras. Users need the ability to organize these devices into hierarchical zones or rooms, and perform actions or read states either on individual devices or across an entire zone simultaneously. Additionally, the system should allow certain devices, like motion sensors, to seamlessly trigger automated reactions in other devices when their state changes.

## 2. Requirement Gathering

* **Constraints:**
  * The system adheres to a strictly top-down Directed Acyclic Graph (DAG). Leaf nodes do not maintain pointers to their parent composites to simplify memory management and avoid cyclic dependencies.
  * The execution model expects high read-throughput (users polling home state) versus write-throughput (sensors tripping).
* **Edge Cases:**
  * A single component failure in a Composite node (e.g., one dead bulb in a Room) must not halt the execution of the broadcast command for the rest of the zone.
* **Scale/Performance considerations:**
  * The environment is highly concurrent. External API requests and internal sensor streams will attempt to read and mutate node states simultaneously, requiring strict `sync.RWMutex` encapsulation.

## 3. API & Model Design

* **Core Interfaces:** `HubNode` (acting as the uniform Component for both Devices and Zones).
* **Entities:** Separated into two distinct packages:
  * `domain`: Contains the physical representations (`Device`, `Zone`). Ignorant of any automation logic.
  * `hub`: Contains the `AutomationHub` registry. Acts as the execution engine mapping publishers to subscribers.
* **Design Decisions:** * **Separation of Concerns:** By keeping the Observer registry out of the `domain` package, we prevented our "dumb" physical devices from becoming tightly coupled to the behavioral routing rules.
  * **Memory Safety:** Placed a `sync.RWMutex` directly on the `Zone`'s `AddNode` method to prevent underlying slice reallocation panics when appending to the topology concurrently.

## 4. Red Phase (TDD)

* [x] **Scenario 1:** 1-to-1 Mapping. Verify a single sensor triggers a single specific device.
* [x] **Scenario 2:** Composite Cascade. Verify a sensor triggering a `Zone` correctly recursively toggles all nested devices.
* [x] **Scenario 3:** Many-to-Many Mapping. Verify a single sensor can independently trigger multiple distinct devices.
* [x] **Scenario 4:** Isolation Principle. Verify triggering Device A has zero side effects on Device B.
* [x] **Scenario 5:** Chaos/Concurrency Test. Spawned 100 reader goroutines and 100 writer goroutines using a `sync.WaitGroup` to ensure the Hub's routing map is safe from concurrent read/write panics under the `-race` detector.

## 5. Green Phase (Implementation Notes & Lessons Learned)

* **Testing Architecture:** Used setup closures inside the table-driven test matrix to isolate physical topology creation. Returned an *Assertion Dictionary* (`map[string]domain.HubNode`) from the setup to decouple the system's internal pub-sub state from the test loop's verification mechanism.
* **Zero-Allocation Aggregation:** Leveraged Go's `strings.Builder` inside the Composite `Zone.GetState()` method to prevent expensive string concatenations and allocations during recursive tree traversal.
* **Lock Ordering & Anti-Patterns:** Explored the "Senior Go" paradigm of copying a slice of subscribers and releasing the Hub's `RLock` *before* iterating over interface methods. While our strict DAG made holding the lock safe here, dropping locks before executing interface callbacks is a critical defensive pattern to prevent deadlocks in complex cyclical systems.
* **Go Type Conversions:** Learned that `string(integer)` converts an int to its Unicode character representation, not its numeric string equivalent (e.g., use `fmt.Sprintf` or `strconv.Itoa` instead).
  