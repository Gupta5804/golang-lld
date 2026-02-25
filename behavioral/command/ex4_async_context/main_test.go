package main

import (
	"context"
	"testing"
	"time"
)

func TestCommandRunner_Success(t *testing.T) {
	cr := &CommandRunner{}
	cmd := &SleepCommand{
		Duration: 10 * time.Millisecond,
	}
	cr.Add(cmd)
	if err := cr.RunWithTimeout(50 * time.Millisecond); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
func TestCommandRunner_Timeout(t *testing.T) {
	cr := &CommandRunner{}
	cmd := &SleepCommand{
		Duration: 100 * time.Millisecond,
	}
	cr.Add(cmd)
	err := cr.RunWithTimeout(10 * time.Millisecond)
	if err == nil {
		t.Fatal("Expected error, instead got nil")
	}
	if err != nil && err != context.DeadlineExceeded {
		t.Errorf("Expected dealine exceeded error, got %v",err)
	}
}
