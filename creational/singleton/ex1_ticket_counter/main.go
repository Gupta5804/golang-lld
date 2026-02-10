package main

import (
	"fmt"
	"sync"
	"time"
)

// --- The Shared Resource ---
type TicketInventory struct {
	count int
	mu sync.Mutex
}

// method to buy a ticket
func (t * TicketInventory) BuyTicket(user string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	time.Sleep(100 * time.Millisecond) // to prove concurrency works
	if t.count > 0 {
		t.count--
		fmt.Printf("User %s bought a ticket. Remaining: %d\n", user, t.count)
	} else {
		fmt.Println("Sold Out!")
	}
}
var (
	instance *TicketInventory
	once sync.Once
)
// --- The Factory ---
// This creates a NEW inventory every time it is called
// func NewTicketInventory() *TicketInventory {
// 	return &TicketInventory{
// 		count: 100, // intial stock
// 	}
// }

func GetInstance() *TicketInventory {
	once.Do(func(){
		fmt.Println("--Running Initialization Logic (runs only once)")
		instance = &TicketInventory{
			count:100,
		}
	})
	return instance
}

func main() {
	var wg sync.WaitGroup
	users := []string{"Alice","Bob","Charlie","Denim","Eve","Fischer"}

	for _,user := range users{
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			inv := GetInstance()
			inv.BuyTicket(u)
		}(user)
	}

	wg.Wait()
	fmt.Println("All sales Done!")
}
