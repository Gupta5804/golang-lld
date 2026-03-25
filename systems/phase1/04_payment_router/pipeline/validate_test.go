package pipeline_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Gupta5804/golang-lld/systems/phase1/04_payment_router/domain"
	"github.com/Gupta5804/golang-lld/systems/phase1/04_payment_router/pipeline"
)

func TestValidateHandler_Execute(t *testing.T) {
	tests := []struct {
		name        string
		request     *domain.PaymentRequest
		expectedErr error
	}{
		{
			name:        "Zero Amount",
			request:     &domain.PaymentRequest{Amount: 0, Region: "US"},
			expectedErr: domain.ErrValidationFailed,
		},
		{
			name:        "Empty Region",
			request:     &domain.PaymentRequest{Amount: 100},
			expectedErr: domain.ErrValidationFailed,
		},
		{
			name:        "Valid",
			request:     &domain.PaymentRequest{Amount: 100, Region: "US"},
			expectedErr: nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			vh := &pipeline.ValidateHandler{}
			err := vh.Execute(context.Background(), tc.request)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected error %v,got %v",tc.expectedErr,err)
			}
		})
	}
}
