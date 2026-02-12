package main

import (
	"testing"
	"sync"
)

func TestGetInstance_IsSingleton(t *testing.T){
	instance1 := GetInstance()
	instance2 := GetInstance()

	if instance1 != instance2 {
		t.Errorf("Expected same instance , but got different instances: %p vs %p",instance1,instance2)
	}
}

func TestAuditLogger_ConcurrentLogging(t *testing.T){
	instance = GetInstance()
	var wg sync.WaitGroup
	
	wg.Add(100)
	for i:=0;i<100;i++{
		go func(){
			defer wg.Done()
			instance.LogEvent("Alice","Login")
		}()
	}
	wg.Wait()

	count := instance.GetEventCount("Alice")
	if count != 100 {
		t.Errorf("Expected 100 count, instead got %d",count)
	}
	
}