package main

import (
	"context"
	"fmt"
	"time"
)

// -- target interface (synchronous) --
type PaymentGateway interface {
	Pay(amount float64) error
}

// Adaptee (3rd Party) (asynchronous)
type CryptoProcessor struct{}

func (c *CryptoProcessor) PerformTransaction(ctx context.Context, amount float64) <-chan bool {
	resultChan := make(chan bool)

	go func() {
		defer close(resultChan)

		select {
		case <-time.After(500 * time.Millisecond):
			resultChan <- true
		case <-ctx.Done():
			return
		}
	}()

	return resultChan
}

// Adapter

type CryptoAdapter struct {
	processor *CryptoProcessor
	Timeout   time.Duration
}

func (c *CryptoAdapter) Pay(amount float64) error {
	// Create context with a Deadline
	// We use context.Background() as the root because Pay() didn't give us one
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()

	// call the async API (adaptee)
	// it returns immediately , giving us a channel to listen to
	resultChan := c.processor.PerformTransaction(ctx, amount)

	// the "bridge" (waiting block)
	// we use 'select' to block until EITHER the result arrives or time runs out

	select {
	case success := <-resultChan:
		if !success {
			return fmt.Errorf("crypto transaction rejected")
		}
		return nil // success!
	case <- ctx.Done():
		return fmt.Errorf("payment timed out")
	}
}
