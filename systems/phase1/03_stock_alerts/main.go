package main

import (
	"fmt"
	"time"

	"github.com/Gupta5804/golang-lld/systems/phase1/03_stock_alerts/domain"
	"github.com/Gupta5804/golang-lld/systems/phase1/03_stock_alerts/feed"
	"github.com/Gupta5804/golang-lld/systems/phase1/03_stock_alerts/notifier"
)

func main() {
	fmt.Println("---- Booting Stock Price Alerting System ----")

	commandChan := make(chan domain.Command, 100)

	dispatcher := notifier.NewDispatcherHub(commandChan, 3)

	dispatcher.Start()

	defer dispatcher.Stop()

	engine := feed.NewEngine(commandChan)

	// ---- Simulation ----
	fmt.Println("[SYS] User 1 subscribing to AAPL > $150.00")
	engine.Subscribe(domain.AlertCondition{
		AlertID:      "alert_1",
		UserID:       "user_1",
		TickerSymbol: "AAPL",
		TargetPrice:  15000,
		Direction:    domain.Above,
	})

	// Simulate market ticks pushing data to the engine
	fmt.Println("[MARKET] AAPL ticks at $149.00...")
	engine.UpdatePrice("AAPL", 14900) // Should NOT trigger

	fmt.Println("[MARKET] AAPL ticks at $151.00...")
	engine.UpdatePrice("AAPL", 15100) // SHOULD trigger!

	fmt.Println("[MARKET] AAPL ticks at $152.00...")
	engine.UpdatePrice("AAPL", 15200) // Should NOT trigger (One-Shot rule deleted it)

	// Give the async workers 100 milliseconds to print to the console before main() exits
	time.Sleep(100 * time.Millisecond)
	fmt.Println("--- Shutting Down System ---")
	close(commandChan)
}
