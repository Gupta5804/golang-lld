package gateways

import "context"

type PayPal struct{}

func (p *PayPal) Pay(ctx context.Context, amount int64) error {
	return nil
}
