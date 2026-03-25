package pipeline

import (
	"context"

	"github.com/Gupta5804/golang-lld/systems/phase1/04_payment_router/domain"
)

type FraudCheckable interface {
	GetIP() string
	GetFingerprint() string
}

type FraudHandler struct {
	BaseHandler
}

func (h *FraudHandler) Execute(ctx context.Context, r *domain.PaymentRequest) error {
	// return h.performCheck(ctx, r)
	if err := h.performCheck(ctx, r); err != nil {
		return err
	}
	return h.executeNext(ctx, r)
}
func (h *FraudHandler) performCheck(ctx context.Context, f FraudCheckable) error {
	ip := f.GetIP()
	fingerprint := f.GetFingerprint()
	if ip == "192.168.1.100" || fingerprint == "blacklisted_device_id" {
		return domain.ErrFraudDetected
	}
	return nil
}
