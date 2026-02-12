package main

import (
	"testing"
	"sync"
	"fmt"
)

func TestGetInstance_IsSingleton(t *testing.T){
	instance1 := GetInstance()
	instance2 := GetInstance()

	if instance1 != instance2 {
		t.Errorf("Expected same instance, but got different instances : %p vs %p", instance1,instance2)
	}
}

func TestConfigManager_ConcurrentAccess(t *testing.T){
	instance := GetInstance()
	var wg sync.WaitGroup
	wg.Add(10)
	for i:=0;i<10;i++ {
		go func(index int) {
			defer wg.Done()

			// Create unique key-value pairs
			key := fmt.Sprintf("config_%d",index)
			value := fmt.Sprintf("value_%d",index)

			// 1. Write(this method doesn't exist yet -> RED)
			instance.Set(key,value)
			// 2. Read
			got := instance.Get(key)

			// 3. verify
			if got != value {
				t.Errorf("Expected %s, got %s",value,got)
			}
		}(i)
	}
	wg.Wait()
}