package main

import (
	"testing"
	"time"
)

func TestMetricsBroker_Unsubscribe(t *testing.T) {
	broker := NewMetricsBroker()

	sub := broker.Subscribe("throughput")
	broker.Unsubscribe("throughput",sub)

	broker.Publish("throughput", 99)
	select {
	case <-time.After(100*time.Millisecond):
		t.Fatal("timeout")
	case val,ok :=<- sub:
		if ok {
			t.Errorf("expected channel to be closed and channel is not closed, got %d",val)
		}
	}
}

func TestMetricsBroker_DropSlowSubscriber(t *testing.T) {
	broker := NewMetricsBroker()
	sub := broker.Subscribe("throughput")

	broker.Publish("throughput",42)
	broker.Publish("throughput",100)
	select {
	case <-time.After(100 * time.Millisecond):
    	t.Fatal("Timeout waiting for first metric (42)")
	case val := <-sub:
    	if val != 42 { t.Errorf("Expected 42, got %d", val) }
	}
	select {
	case <- time.After(100*time.Millisecond):
	case val:=<-sub :
		t.Errorf("Expected message to be dropped, but received %d", val) 
	}
}
