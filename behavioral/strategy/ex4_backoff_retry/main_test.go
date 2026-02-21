package main

import (
	"slices"
	"testing"
	"time"
)

func TestSimulateRetries_ConstantBackoff(t *testing.T) {
	constant := &ConstantBackoff{delay: 2 * time.Second}
	simulator := NewRetrySimulator(constant)

	attempts := simulator.SimulateRetries(3)
	expected := []time.Duration{2 * time.Second, 2 * time.Second, 2 * time.Second}
	if !slices.Equal(attempts, expected) {
		t.Errorf("Expected [2,2,2], got %v", attempts)
	}

}

func TestSimulateRetries_ExponentialBackoffs(t *testing.T) {
	exponential := &ExponentialBackoff{base : 100 * time.Millisecond}
	simulator := NewRetrySimulator(exponential)
	attempts := simulator.SimulateRetries(4)
	expected := []time.Duration{100 * time.Millisecond, 200* time.Millisecond, 400 * time.Millisecond, 800 * time.Millisecond}

	if !slices.Equal(attempts,expected){
		t.Errorf("Expected [100,200,400,800], got %v", attempts)
	}
}
