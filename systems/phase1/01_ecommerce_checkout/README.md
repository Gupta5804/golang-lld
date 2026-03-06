# E-Commerce Checkout Pricing Engine

**Phase:** 1 (Foundation)  
**Core Patterns:** Strategy, Decorator  

## 1. Initial Problem Statement

Design an in-memory pricing engine for a new checkout flow. The system needs to calculate the final price of a user's cart by applying various dynamic discounts, promotional codes, and regional taxes. Pricing rules can change frequently and might need to be stacked or combined in different ways depending on the ongoing sale.

## 2. Requirement Gathering

### *Scope and boundaries*

* **Input Contract:** The engine will accept a `Cart` Object. The cart contains a list of `Item` objects (*which have their own base price*) and a method to get the subtotal. You need the full Cart because in future we might have item-specific discounts.
* **Rule Configuration:** You are building the algorithms (the generic strategies), not the specific campaigns. For e.g., you should build a generic `PercentageDiscount` rule that can be configured at `10%` or `20%` at runtime. Do not build the database that fetches the active rules;*assume that the active rules for a given checkout session are passed into your engine's constructor or execution method.*
  
### *Business Logic & Edge Cases*

* **Rule Types:** You must support `PercentDiscount`, `FixedAmountDiscount`, and `Tax` (*which adds to the total instead of subtracting*).
* **Zero-Bound Behavior:** The cart total can never go below zero. If a $50 fixed discount is applied to a $40 cart, the total simply floors at $0. Do not return an error;just make it free (*Taxes, however can obviously push it above zero again*).
* **Order of Operations:** The engine should be dumb regarding order. It should simply apply the rules in the exact sequence they are provided in the input slice/chain.
* **Validation:** Assume all rules are pre-validated.You will just receive the hydrated mathematical strategy objects.
* **Rounding:** We will represent all money as `int` (cents) to avoid floating-point issues.Standard rounding applies.
  
### *Entity & Doma•in Modeling*

* **Item Schema:** To keep it extensible for future phases, an `Item` needs a `SKU`(string), a `Price`(int, in cents), and a `Quality`(int).
* **Cart Encapsulation:** The `Cart` will have a slice of `Item`s. It should indeed have a `BaseTotal() int` method attached to it via a pointer receiver. It does not need a region code; assume the specific tax strategy is chosen by a higher-level service before the engine is even called.
* **Immutability:** The engine should treat the `Cart` as immutable. Do not add fields to the `Cart` to track discounts. The engine's `Calculate` method should just return the final calculated `int`(and an `error`).
  
### *System Constraints*

* **Traffic & SLA:** This is a tier-1 system.During Sales, it will hit 10,000 TPS. The `Calculate` method must be blisteringly fast--essentially zero-allocation if possible. GC pauses are our enemy.
* **Concurrency(The Engine):** The `PricingEngine` will be instantiated once at system startup (like a singleton or injected dependency) and shared across thousands of concurrent HTTP handler routines.
* **Concurrency(The Cart):** The `Cart` is strictly request-scoped. Only *one* goroutine will ever touch a specific `Cart` at a time. No mutexes are needed on the Cart itself.
* **Memory Bounds:** Carts generally have fewer than 50 items.You will rarely evaluate more than 10 rules per checkout.

## 3. API & Model Design

* **Entities (`domain` package):** `Cart` encapsulates `Item`s and provides a `BaseTotal()` method via a pointer receiver to prevent memory copying.
* **Engine (`pricing` package):** Uses the Strategy Pattern. `PricingRule` is an interface requiring an `Apply` method.
* **Thread Safety:** The `PricingEngine.rules` slice is unexported. It is injected via `NewEngine()` and cannot be mutated by external concurrent HTTP handlers.

## 4. Red Phase (TDD)

Implemented a Table-Driven Test suite covering:

* Atomic rule applications (Percentage, Fixed, Tax).
* Zero-bound flooring (preventing negative cart totals).
* Sequential stacking proofs (Order of Operations).
* Complex integration (Fixed -> Percent -> Tax).

## 5. Green Phase

* Leveraged Go's implicit interfaces.
* Ensured zero-allocation during the `Calculate` loop.
* Implemented strict integer-based percentage math to protect system accounting.
