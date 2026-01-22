package main

import (
	"testing"
)


func TestNewPaymentGateway_Stripe(t *testing.T) {
	gateway, err := NewPaymentGateway("stripe")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, ok := gateway.(*StripeGateway)
	if !ok {
		t.Errorf("Expected *StripeGateway, got %T", gateway)
	}
}

func TestNewPaymentGateway_Bitcoin(t *testing.T) {
	gateway, err := NewPaymentGateway("bitcoin")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, ok := gateway.(*BitcoinGateway)
	if !ok {
		t.Errorf("Expected *BitcoinGateway, got %T", gateway)
	}
}

func TestNewPaymentGateway_Invalid(t *testing.T) {
	_, err := NewPaymentGateway("invalid_gateway")

	if err == nil {
		t.Error("Expected error here, but got nil")
	}
}



type MockGateway struct{
	CapturedAmount float64
	PayCalled bool
}

func (m *MockGateway) Pay(amount float64) error {
	m.CapturedAmount = amount
	m.PayCalled = true

	return nil
}

func TestPay_ValidAmount(t *testing.T){
	mock := &MockGateway{}
	service := NewCheckoutService(mock)
	
	err := service.ProcessPayment(100)
	
	if err != nil {
		t.Errorf("Expected success, got error:%v", err)
	}

	if !mock.PayCalled {
		t.Error("Expected to Pay function to be called but didnt")
	}
	if mock.CapturedAmount != 100 {
		t.Errorf("Expected amount to be '100' but got : %f", mock.CapturedAmount)
	} 
}
func TestPay_InvalidAmount(t *testing.T){
	mock := &MockGateway{}
	service := NewCheckoutService(mock)
	
	err := service.ProcessPayment(-100)
	

	if err == nil {
		t.Error("Expected error for negative amount, got nil")
	}
	if mock.PayCalled {
		t.Error("Expected to Pay function to not be called but did")
	}
}