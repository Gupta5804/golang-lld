package main

import (
	"fmt"
	"testing"
)

func TestDispatcher_WorkerPool(t *testing.T) {
	results := make(chan Result, 10)
	d := NewDispatcher(3, 10) // 3 workers, buffer of 10
	d.Start(results)
	for i := 1; i <= 5; i++ {
		job := Job{ID: i ,FilePath: fmt.Sprintf("C:/%d",i)}
		d.Publish(job)
	}
	d.Stop()
	close(results)
	receivedIDs := make(map[int]bool)
	counter := 0
	for result:= range results{
		receivedIDs[result.JobID] = true
		counter++
	}

	if counter != 5 {
		t.Errorf("Expected 5 IDs to be received , got %d",counter)
	}
	for i:=1;i<=5;i++{
		if !receivedIDs[i]{
			t.Errorf("Job with ID: %d not received",i)
		}
	}
}
