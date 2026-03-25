package gateways

import (
	"context"
	"fmt"
)

type Stripe struct {
}

func (s *Stripe) Pay(ctx context.Context, amount int64) error {
	fmt.Printf("Processing %d cents with stripe",amount)
	return nil
}
