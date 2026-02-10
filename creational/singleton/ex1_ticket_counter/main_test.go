package main

import (
	"testing"
	"sync"
)


func TestGetInstance_ReturnsSamePointer(t *testing.T){
	ins1 := GetInstance()
	ins2 := GetInstance()
	
	if ins1 != ins2 {
		t.Errorf("Expected same instance , got distinct pointers: %p vs %p",ins1,ins2)
	}
}


func TestTicketInventory_BuyTicket_Concurrent(t *testing.T){
	var wg sync.WaitGroup
	inv := GetInstance()
	initialCount := inv.count
	ticketsToBuy := 50
	wg.Add(ticketsToBuy)

	for i:=0;i<ticketsToBuy;i++{
		go func(){
			defer wg.Done()
			inv.BuyTicket("test_user")
		}()
	}
	wg.Wait()
	expectedCount := initialCount-ticketsToBuy

	if inv.count != expectedCount {
		t.Errorf("Expected remaining tickets : %d, got %d",expectedCount, inv.count)
	}
}