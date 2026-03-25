package pipeline

import (
	"context"

	"github.com/Gupta5804/golang-lld/systems/phase1/04_payment_router/domain"
)

type Handler interface {
	SetNext(handler Handler) Handler
	Execute(ctx context.Context, req *domain.PaymentRequest) error
}

type BaseHandler struct {
	next Handler
}

func (h *BaseHandler) SetNext(next Handler) Handler {
	h.next = next
	return next
}
func (h *BaseHandler) executeNext(ctx context.Context, req *domain.PaymentRequest) error {
	if h.next != nil{
		return h.next.Execute(ctx,req)
	}
	return nil
}
