package main

import "time"

type BackoffStrategy interface {
	NextDelay(attempt int) time.Duration
}

type ConstantBackoff struct {
	delay time.Duration
}

func (c *ConstantBackoff) NextDelay(attempt int) time.Duration {
	return c.delay
}

type ExponentialBackoff struct {
	base time.Duration
}

func (e *ExponentialBackoff) NextDelay(attempt int) time.Duration {
	return e.base * time.Duration(1<<attempt)
}

type RetrySimulator struct {
	strategy BackoffStrategy
}

func NewRetrySimulator(strategy BackoffStrategy) *RetrySimulator {
	return &RetrySimulator{
		strategy: strategy,
	}
}

func (r *RetrySimulator) SimulateRetries(maxAttempts int) []time.Duration {
	result := []time.Duration{}
	for i:=0;i<maxAttempts;i++{
		result = append(result,r.strategy.NextDelay(i))
	}
	return result
}
