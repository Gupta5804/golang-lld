package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Gupta5804/golang-lld/systems/phase1/07_logistics_cost_estimator/domain"
	"github.com/Gupta5804/golang-lld/systems/phase1/07_logistics_cost_estimator/estimator"
	"github.com/Gupta5804/golang-lld/systems/phase1/07_logistics_cost_estimator/legacy"
	"github.com/Gupta5804/golang-lld/systems/phase1/07_logistics_cost_estimator/strategies"
)

func main() {
	fmt.Println("=== Logistics Shipping Cost Estimator ===")

	// 1. Initialize the external dependencies (The Black Box)
	regionalAPI := &legacy.RegionalCourierAPI{}

	// 2. Initialize our Anti-Corruption Layer (The Adapter)
	adapter := legacy.NewLegacyAdapter(regionalAPI)

	// 3. Initialize the Core Engine (Estimator Service)
	// We inject the adapter so the engine never touches the black box directly.
	service := estimator.NewEstimatorService(adapter)

	// 4. Initialize our Pricing Strategies
	standard := &strategies.StandardStrategy{}
	holiday := &strategies.HolidayStrategy{}

	// 5. Create a dummy shipment: 50 KG, 100 KM
	shipment := domain.Shipment{
		WeightKG:   50.0,
		DistanceKM: 100,
	}

	ctx := context.Background()

	// --- Run Scenario 1: Standard Pricing ---
	fmt.Println("\n--- Scenario 1: Standard Shipping ---")
	quote1, err := service.CalculateTotal(ctx, shipment, standard)
	if err != nil {
		log.Fatalf("Error calculating quote: %v", err)
	}
	printQuote(quote1)

	// --- Run Scenario 2: Holiday Pricing ---
	fmt.Println("\n--- Scenario 2: Holiday Surge Shipping ---")
	quote2, err := service.CalculateTotal(ctx, shipment, holiday)
	if err != nil {
		log.Fatalf("Error calculating quote: %v", err)
	}
	printQuote(quote2)
}

// printQuote converts our internal int64 cents into human-readable dollar strings
func printQuote(q *domain.Quote) {
	fmt.Printf("Line-Haul Cost: $%.2f\n", float64(q.LineHaulCost)/100.0)
	fmt.Printf("Last-Mile Cost: $%.2f\n", float64(q.LastMileCost)/100.0)
	fmt.Printf("Total Cost:     $%.2f\n", float64(q.TotalCost())/100.0)
}