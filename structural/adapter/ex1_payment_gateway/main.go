package main

import (
	"fmt"
)

type PaymentGateway interface{
	Pay(amount float64) error
}

// 1st Adaptee
type PayPal struct{}
func (p *PayPal) MakePayment(amount float64, currency string) bool {
	return true
}


// 2nd Adaptee
type Stripe struct{}
func (s *Stripe) PayRequest(amountCents int) (string, error) {
	// simulate processing
	return "txn_12345", nil
}

// Adapter 1

type PayPalAdapter struct{
	paypal *PayPal
}

func (p *PayPalAdapter) Pay(amount float64) error {
	if success := p.paypal.MakePayment(amount,"USD"); !success{
		return fmt.Errorf("payment failed")
	}
	return nil
}

// Adapter 2

type StripeAdapter struct{
	stripe *Stripe
}
func (s *StripeAdapter) Pay(amount float64) error {
	if _,err := s.stripe.PayRequest(int(amount*100)); err != nil{
		return fmt.Errorf("stripe gateway error: %w",err)
	}
	return nil
}