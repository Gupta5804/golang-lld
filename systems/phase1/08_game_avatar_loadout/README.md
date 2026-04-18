# LLD Study Notes: Game Avatar Loadout System

**Goal (Initial Prompt):** We are launching a new RPG feature in our app. Players get a basic avatar, but during a match, they can pick up all sorts of power-ups like 'speed boosts', 'fire auras', or 'heavy armor'. They can stack these buffs however they want, and the system needs to calculate their total stats on the fly. Also, the monetization team wants to sell pre-configured 'starter bundles', so we need a really clean, readable way in the backend to spawn these new avatars with different combinations of starting health, base damage, and special VIP flags without making a mess of the creation logic.

(*Design an RPG avatar configuration and runtime loadout system. The monetization team requires a clean, readable API to spawn "Starter Bundles" with heavily customized base stats and VIP flags without creating bloated constructors. During a match, players can equip stacking power-ups (buffs/debuffs) that modify their base stats on the fly. The system must rapidly calculate the final theoretical max stats based on the active loadout combination.*)

**Phase:** 1 (Foundation)  
**Core Patterns:** Functional Options (Creational), Decorator (Structural)

---

## 1. The 4 Pillars (Requirement Gathering)

Before writing any code, we defined the strict boundaries of the system to prevent domain bleed and establish clear mathematical constraints.

* **Scope & Boundaries:**
  * **Capacity vs. State:** This module is strictly responsible for calculating the *maximum theoretical stats* based on the active loadout. It does not track real-time combat mutations (e.g., taking damage).
  * **Permanent Equip:** For the MVP, equipped buffs are permanent for the duration of the match. There is no requirement to "un-equip" a decorator mid-chain.
* **Business Logic & Edge Cases:**
  * **Multi-Stat Mutations:** A single item (e.g., Heavy Armor) can increase one stat while simultaneously penalizing another (handling negative values).
  * **Order Dependence:** Because modifiers can be multiplicative (e.g., +20% Damage) or additive (e.g., +5 Damage), the order in which items are picked up mathematically alters the final state.
* **Entity & Domain Modeling:**
  * **Decoupled Attributes:** We separated static identity attributes (`ID`, `VIPFlags`) from dynamic combat attributes (`Stats`).
  * **Float Precision:** All calculable attributes (`Health`, `Damage`, `Speed`) use `float64` to cleanly handle percentage-based math without invisible integer truncation bugs.
* **System Constraints:**
  * **High-Frequency Reads:** The combat engine reads these final calculations 60 times a second. The loadout chain must be 100% thread-safe.
  * **High-Throughput Spawning:** Matchmaking spawns hundreds of avatars simultaneously; creation must be clean, declarative, and allocation-efficient.

---

## 2. Architectural Design (The Security Scanner Analogy)

We utilized the **Functional Options** pattern for the initial instantiation and the **Decorator** pattern (Structural variant) for the runtime power-ups.

**The Analogy:**
Imagine a security scanner at an airport that only asks, "What is your total weight?"

* **The Interface (`StatsProvider`):** The scanner's rule. Anything passing through must have a `CalculateStats()` method.
* **The Base (`Avatar`):** A person walks through. They answer by providing their base body weight.
* **The Decorator (`HeavyArmor`):** Put a metal box over the person. The scanner asks the box for its weight. The box turns around, asks the person inside for their weight, adds 50 lbs to it, and gives the final total to the scanner. The scanner is completely blind to the wrapping.

**Directory Blueprint:**

```text
08_game_avatar_loadout/
├── domain/
│   └── models.go       <-- Pure entities (Avatar, Stats) and the StatsProvider interface.
├── avatar/
│   ├── avatar.go       <-- NewAvatar constructor logic
│   └── options.go      <-- Creational Functional Options (WithBaseHealth, WithVIPFlags)
├── buffs/
│   ├── buffs.go        <-- Structural Decorators (SpeedBoots, HeavyArmor, BerserkerPotion)
│   └── buffs_test.go   <-- Black-box TDD, Order-Dependence Matrix, Chaos Concurrency Proof
└── main.go             <-- Composition Root & Simulation
```

---

## 3. Data Structure & Contract Decisions (Staff-Level Choices)

* **Structural over Functional Decorators:** Instead of passing functions that return functions (which creates a black-box closure impossible to inspect), we defined explicit Structs containing the `domain.StatsProvider` interface. This keeps the nodes tangible in memory, allowing for future inspection or serialization.
* **Return by Value Immutability:** The `StatsProvider` interface dictates that `CalculateStats()` returns the `domain.Stats` struct by *value*, not by pointer. This creates an unmovable mathematical barrier protecting the base avatar from downstream corruption.
* **Declarative Configurations:** Kept the `NewAvatar` signature strictly limited to `(id string, opts ...Option)`. This guarantees the core constructor never has to change, even if the monetization team invents 50 new configuration parameters next year.

---

## 4. Go-Specific Traps Avoided (The "Gotcha Ledger")

### Trap 1: The Pointer Bleed (Catastrophic State Mutation)

* **The Trap:** A junior engineer writing a decorator might call `stats := s.Provider.CalculateStats()`, assuming it's safe. If `CalculateStats()` returned a pointer (`*domain.Stats`), executing `stats.Speed += 5` inside the `SpeedBoots` buff would permanently mutate the *base avatar's* health in memory.
* **The Fix:** The interface contract mandates returning `domain.Stats` by value. Decorators fetch a copy, mutate the copy locally, and return the newly modified copy. The core entity remains pristine.

### Trap 2: State Bleed in Table-Driven Tests

* **The Trap:** When building the TDD matrix, defining the `StatsProvider` directly as a field in the test struct allows the same memory allocation to be reused across different test loops, causing tests to mysteriously fail depending on their execution order.
* **The Fix:** Used an anonymous closure `setup: func() domain.StatsProvider` in the test struct. Executing this closure inside `t.Run()` guarantees every single test case spins up a completely fresh, isolated memory allocation.

### Trap 3: The Premature WaitGroup Evaluation

* **The Trap:** In the `-race` concurrency test, placing `wg.Add(1)` *inside* the goroutine execution loop. The test function can hit `wg.Wait()` and exit cleanly before the Go scheduler even has time to spin up the goroutines, resulting in a false positive test pass.
* **The Fix:** Placed `wg.Add(workers)` securely outside the loop, establishing an unmovable synchronization barrier before any concurrent execution begins.

---

## 5. Table-Driven TDD (The Red Phase)

We mathematically proved the system's behavior using Table-Driven Tests prior to implementation.

* **The Order-Dependence Proof:** We explicitly created two mirrored test cases in the matrix (Case 4: Base + Armor + Potion vs. Case 5: Base + Potion + Armor). By asserting against two entirely different expected totals (135 HP vs 140 HP), we mathematically proved that our Decorator chain correctly respects the strict application order of multiplicative vs. additive modifiers.
* **The Chaos Concurrency Test:** Deployed a `sync.WaitGroup` to fire 1,000 concurrent goroutines against a deeply nested Avatar (wrapped in all three buffs). Running `go test -race` proved that iterating through the deeply nested interface chain for massive read-throughput is completely thread-safe and free of data races.

---

## 6. Quick Start

To verify the test suite (and the race detector) or run the end-to-end Composition Root simulation:

```bash
# Run the black-box test suite to prove the math, order-dependence, and concurrency safety
go test -v -race ./buffs/...

# Execute the Composition Root to view the lifecycle simulation and Immutability Check
go run main.go
```
