package main

import (
	"testing"
	"time"
)

func TestEventBus_Broadcast(t *testing.T) {
	bus := NewEventBus()
	sub1 := bus.Subscribe("order.placed")
	sub2 := bus.Subscribe("order.placed")

	bus.Publish("order.placed", "ORD-1234")

	// sequential verification
	// verify subscriber 1
	select {
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout waiting for the event")
	case msg1 := <-sub1:
		if msg1 != "ORD-1234" {
			t.Errorf("Expected msg1 'ORD-1234',got %s", msg1)
		}

	}

	// verify subscriber 2
	select {
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout waiting for the event")
	case msg2 := <-sub2:
		if msg2 != "ORD-1234" {
			t.Errorf("Expected msg2 'ORD-1234',got %s", msg2)
		}
	}
}
