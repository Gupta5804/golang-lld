package main

import (
	"testing"
	"time"
)

func TestMetricsMiddleware(t *testing.T) {
	mockMetrics := make(map[string]string)
	finalHandler := func(r Request) error {
		time.Sleep(50 * time.Millisecond)
		return nil
	}
	chain := BuildChain(finalHandler, MetricsMiddleware(mockMetrics))

	req := map[string]string{
		"dummy": "dummy",
	}
	if err := chain(req); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if val,ok := mockMetrics["status"]; !ok{
		t.Fatal("no value found in mockMetrics['status']")
	} else if val != "success"{
		t.Errorf("Expected 'success', got %s",val)
	}
}
