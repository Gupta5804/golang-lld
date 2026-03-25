package pipeline_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Gupta5804/golang-lld/systems/phase1/04_payment_router/domain"
	"github.com/Gupta5804/golang-lld/systems/phase1/04_payment_router/pipeline"
)

func TestFraudHandler_Execute(t *testing.T) {
	tests := []struct {
		name        string
		request     *domain.PaymentRequest
		expectedErr error
	}{
		{
			name: "Bad IP",
			request: &domain.PaymentRequest{
				ClientIP:          "192.168.1.100",
				DeviceFingerprint: "safedevice",
			},
			expectedErr: domain.ErrFraudDetected,
		},
		{
			name: "Bad Device",
			request: &domain.PaymentRequest{
				ClientIP:          "192.168.1.120",
				DeviceFingerprint: "blacklisted_device_id",
			},
			expectedErr: domain.ErrFraudDetected,
		},
		{
			name: "Safe",
			request: &domain.PaymentRequest{
				ClientIP:          "192.168.1.120",
				DeviceFingerprint: "safedevice",
			},
			expectedErr: nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fh := &pipeline.FraudHandler{}
			err := fh.Execute(context.Background(), tc.request)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
