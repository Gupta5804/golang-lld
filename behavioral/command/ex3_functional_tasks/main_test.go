package main

import (
	"errors"
	"testing"
)

func TestTaskRunner_RunAll(t *testing.T) {
	tr := &TaskRunner{}
	executionCount := 0
	task1 := func() error {
		executionCount++
		return nil
	}
	task2 := func() error {
		return errors.New("cache clear failed")
	}
	task3 := func() error{
		executionCount++
		return nil
	}

	tr.AddTask(task1)
	tr.AddTask(task2)
	tr.AddTask(task3)
	errorList := tr.RunAll()
	if executionCount != 2{
		t.Errorf("Expected exectionCount 2, got %d",executionCount)
	}
	if len(errorList) != 1 {
		t.Errorf("Expected 1 error to be returned,got %d",len(errorList))
	}
	if errorList[0].Error() != "cache clear failed"{
		t.Errorf("Expected error 'cache clear failed', got %v",errorList[0].Error())
	}
}
