package main

import (
	"testing"
)


func TestPizza_WithToppings(t *testing.T) {
	pizza := &VeggieMania{}
	cheesePizza := &CheeseTopping{
		pizza: pizza,
	}
	fullPizza := &TomatoTopping{
		pizza:cheesePizza,
	}
	if price:=fullPizza.GetPrice() ; price != 30 {
		t.Errorf("Expected price 30, got %d",price)
	}
}
