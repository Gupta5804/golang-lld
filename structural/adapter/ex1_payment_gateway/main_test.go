package main

import (
	"testing"
)



func TestAdapter_Pay_Success(t *testing.T) {
	// setup the Adaptee
	paypal := &PayPal{}
	
	// Initialize the Adapter
	// We Explicitly define the type as PaymentGateway
	// This ensures our adapter actually implements the interface we need
	var gateway PaymentGateway = &PayPalAdapter{
		paypal : paypal,
	}

	amount := 100.0
	err := gateway.Pay(amount)
	if err != nil {
		t.Errorf("Expected success (nil error), but got :%v", err)
	}
}

func TestPayment_Polymorphism(t *testing.T){
	paypal := &PayPal{}
	stripe := &Stripe{}

	gateways := []PaymentGateway{
		&PayPalAdapter{paypal:paypal}, 
		&StripeAdapter{stripe:stripe},
	}

	for _,gateway := range gateways{
		if err := gateway.Pay(100.0); err != nil{
			t.Errorf("Expected success , but got this error:%v", err)
		}
	}
}