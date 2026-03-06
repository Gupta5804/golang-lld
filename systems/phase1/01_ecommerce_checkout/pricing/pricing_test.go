package pricing_test

import (
	"testing"

	"github.com/Gupta5804/golang-lld/systems/phase1/01_ecommerce_checkout/domain"
	"github.com/Gupta5804/golang-lld/systems/phase1/01_ecommerce_checkout/pricing"
)

func TestPricingEngine_Calculate(t *testing.T) {
	tests := []struct {
		name          string
		cart          domain.Cart
		rules         []pricing.PricingRule
		expectedTotal int64
	}{
		// --- Atomic/Happy paths ---
		{
			name: "empty cart should return 0",
			cart: domain.Cart{
				Items: []domain.Item{},
			},
			rules:         nil,
			expectedTotal: 0,
		},
		{
			name: "apply single percentage discount",
			cart: domain.Cart{
				Items: []domain.Item{
					{Price: 10000, Quantity: 1},
				},
			},
			rules: []pricing.PricingRule{
				&pricing.PercentageDiscount{Percent: 10},
			},
			expectedTotal: 9000, // assuming $100 cart
		},
		{
			name: "apply single fixed discount",
			cart: domain.Cart{
				Items: []domain.Item{
					{Price: 10000, Quantity: 1},
				},
			},
			rules: []pricing.PricingRule{
				&pricing.FixedDiscount{FixedAmount: 1000},
			},
			expectedTotal: 9000,
		},
		{
			name: "apply single tax",
			cart: domain.Cart{
				Items: []domain.Item{
					{Price: 10000, Quantity: 1},
				},
			},
			rules: []pricing.PricingRule{
				&pricing.TaxRule{Percent: 5},
			},
			expectedTotal: 10500,
		},

		// --- Zero-bound edge cases ---
		{
			name: "fixed discount exceeding cart total floors at zero",
			cart: domain.Cart{
				Items: []domain.Item{
					{Price: 4000, Quantity: 1},
				},
			},
			rules: []pricing.PricingRule{
				&pricing.FixedDiscount{FixedAmount: 5000},
			},
			expectedTotal: 0, // does not go negative or gives error
		},
		// --- Stacking & Ordering Proofs ---
		{
			name: "prove order matters: percent then fixed",
			cart: domain.Cart{
				Items: []domain.Item{
					{Price: 10000, Quantity: 1},
				},
			},
			rules: []pricing.PricingRule{
				&pricing.PercentageDiscount{Percent: 10},
				&pricing.FixedDiscount{FixedAmount: 1000},
			},
			expectedTotal: 8000,
		},
		{
			name: "prove order matters: fixed then percent",
			cart: domain.Cart{
				Items: []domain.Item{
					{Price: 10000, Quantity: 1},
				},
			},
			rules: []pricing.PricingRule{
				&pricing.FixedDiscount{FixedAmount: 1000},
				&pricing.PercentageDiscount{Percent: 10},
			},
			expectedTotal: 8100,
		},

		// --- Complex Integration ---
		{
			name: "realistic checkout: standard order with tax",
			cart: domain.Cart{
				Items: []domain.Item{
					{Price: 10000, Quantity: 1},
				},
			},
			rules:         []pricing.PricingRule{ 
				&pricing.FixedDiscount{FixedAmount: 1000},
				&pricing.PercentageDiscount{Percent: 10},
				&pricing.TaxRule{Percent:5},
			},
			expectedTotal: 8505,
		},
	}
	// The runner
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			engine := pricing.NewEngine(tc.rules)
			total, err := engine.Calculate(&tc.cart)
			if err != nil { t.Fatalf("unexpected error: %v", err)}
			if total != tc.expectedTotal{
				t.Errorf("Expected total %d, got %d",tc.expectedTotal,total)
			}
		})
	}
}
