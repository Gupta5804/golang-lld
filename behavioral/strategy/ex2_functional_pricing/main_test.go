package main

import "testing"

func TestCheckout_PricingStrategies(t *testing.T) {
	checkout := &Checkout{Total: 100.0}
	if result := checkout.CalculateFinal(StandardPricing); result != 100 {
		t.Errorf("Expected 100, got %f",result)
	}
	if result := checkout.CalculateFinal(BlackFridayPricing); result != 80 {
		t.Errorf("Expected 80, got %f", result)
	}
}
