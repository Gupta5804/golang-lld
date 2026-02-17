package main

import "testing"


func TestPizza_Functional(t *testing.T){
	myPizza := VeggieMania()
	myPizza = AddCheese(myPizza)
	myPizza = AddTomato(myPizza)

	if price := myPizza() ; price != 30 {
		t.Errorf("Expected price 30, instead got %d", price)
	}
}
