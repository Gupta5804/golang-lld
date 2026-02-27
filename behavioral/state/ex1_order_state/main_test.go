package main

import "testing"

func TestOrder_StateMachine(t *testing.T) {
	// Scenario A (Happy Path)
	order := NewOrder()
	if state := order.GetStateName(); state != "Pending" {
		t.Errorf("Expected state to be 'Pending',got %s", state)
	}
	if err := order.Pay(); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if err := order.Ship(); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Scenario B(Invalid Action)
	order = NewOrder()
	if err := order.Ship(); err == nil {
		t.Error("Expected error , got nil")
	} else if err.Error() != "cannot ship a pending order" {
		t.Errorf("Expected 'cannot ship a pending order' error, got %v", err)
	}
	if state := order.GetStateName(); state != "Pending" {
		t.Errorf("Expected state to be 'Pending',got %s", state)
	}

	// Scenario C(Cancellation and terminal state)
	order = NewOrder()
	if err := order.Cancel(); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if state := order.GetStateName(); state != "Cancelled" {
		t.Errorf("Expected state 'Cancelled',got %s", state)
	}
	if err := order.Pay();err == nil{
		t.Error("Expected error, got nil")
	} else if err.Error() != "cannot pay for a cancelled order"{
		t.Errorf("Expected 'cannot pay for a cancelled order', got %v", err)
	}

}
