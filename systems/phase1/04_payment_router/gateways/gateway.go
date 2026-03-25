package gateways

import "context"

type PaymentGateway interface {
	Pay(ctx context.Context, amount int64) error
}
