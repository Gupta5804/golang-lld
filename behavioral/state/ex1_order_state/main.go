package main

import "fmt"

type Order struct {
	currentState State
}

func NewOrder() *Order {
	order := &Order{}
	order.SetState(&PendingState{order: order})
	return order
}
func (o *Order) SetState(s State) {
	o.currentState = s
}
func (o *Order) GetStateName() string {
	return o.currentState.Name()
}
func (o *Order) Pay() error {
	return o.currentState.Pay()
}
func (o *Order) Ship() error {
	return o.currentState.Ship()
}
func (o *Order) Cancel() error {
	return o.currentState.Cancel()
}

type State interface {
	Name() string
	Pay() error
	Ship() error
	Cancel() error
}

// Pending state
type PendingState struct {
	order *Order
}

func (s *PendingState) Name() string {
	return "Pending"
}
func (s *PendingState) Pay() error {
	s.order.SetState(&PaidState{order: s.order})
	return nil
}
func (s *PendingState) Ship() error {
	return fmt.Errorf("cannot ship a pending order")
}
func (s *PendingState) Cancel() error {
	s.order.SetState(&CancelledState{order: s.order})
	return nil
}

// Paid State
type PaidState struct {
	order *Order
}

func (s *PaidState) Name() string {
	return "Paid"
}
func (s *PaidState) Pay() error {
	
	return fmt.Errorf("already paid")
}
func (s *PaidState) Ship() error {
	s.order.SetState(&ShippedState{order:s.order})
	return nil
}
func (s *PaidState) Cancel() error {
	return fmt.Errorf("cannot cancel a paid order")
}

// Shipped State
type ShippedState struct {
	order *Order
}

func (s *ShippedState) Name() string {
	return "Shipped"
}
func (s *ShippedState) Pay() error {
	return fmt.Errorf("already paid and shipped")
}
func (s *ShippedState) Ship() error {
	return fmt.Errorf("already shipped")
}
func (s *ShippedState) Cancel() error {
	return fmt.Errorf("cannot cancel a shipped order")
}

// Cancelled State
type CancelledState struct {
	order *Order
}

func (s *CancelledState) Name() string {
	return "Cancelled"
}
func (s *CancelledState) Pay() error {
	return fmt.Errorf("cannot pay for a cancelled order")
}
func (s *CancelledState) Ship() error {
	return fmt.Errorf("cannot ship a cancelled order")
}
func (s *CancelledState) Cancel() error {
	return fmt.Errorf("already cancelled")
}
