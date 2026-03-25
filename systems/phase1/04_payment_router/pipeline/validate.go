package pipeline

import (
	"context"

	"github.com/Gupta5804/golang-lld/systems/phase1/04_payment_router/domain"
)

type Validatable interface {
	GetAmount() int64
	GetRegion() string
}

type ValidateHandler struct {
	BaseHandler
}

func (h *ValidateHandler) Execute(ctx context.Context, r *domain.PaymentRequest) error {
	if err := h.validate(ctx, r); err != nil {
		return err
	}
	return h.executeNext(ctx, r)
}
func (h *ValidateHandler) validate(ctx context.Context, v Validatable) error {
	amount := v.GetAmount()
	region := v.GetRegion()
	if amount <= 0 || region == "" {
		return domain.ErrValidationFailed
	}
	return nil
}
