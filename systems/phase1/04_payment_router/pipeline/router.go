package pipeline

import (
	"context"

	"github.com/Gupta5804/golang-lld/systems/phase1/04_payment_router/domain"
	"github.com/Gupta5804/golang-lld/systems/phase1/04_payment_router/gateways"
)

type RouterHandler struct {
	BaseHandler
	registry map[string]gateways.PaymentGateway
}

func NewRouterHandler(registry map[string]gateways.PaymentGateway) *RouterHandler {
	return &RouterHandler{
		registry: registry,
	}
}

func (h *RouterHandler) Execute(ctx context.Context, r *domain.PaymentRequest) error{
	if gateway,ok := h.registry[r.Region]; !ok{
		return domain.ErrRegionNotFound
	} else {
		return gateway.Pay(ctx,r.Amount)
	}
}
