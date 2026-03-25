package pipeline_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Gupta5804/golang-lld/systems/phase1/04_payment_router/domain"
	"github.com/Gupta5804/golang-lld/systems/phase1/04_payment_router/gateways"
	"github.com/Gupta5804/golang-lld/systems/phase1/04_payment_router/pipeline"
)

type MockGateway struct {
	Called bool
}

func (m *MockGateway) Pay(ctx context.Context, amount int64) error {
	m.Called = true
	return nil
}

func TestRouterHandler_Execute(t *testing.T) {
	tests := []struct {
		name        string
		request     *domain.PaymentRequest
		setup       func() (*pipeline.RouterHandler, *MockGateway)
		expectedErr error
	}{
		{
			name:    "Supported Region",
			request: &domain.PaymentRequest{Region: "US", Amount: 100},
			setup: func() (*pipeline.RouterHandler, *MockGateway) {
				mock := &MockGateway{}
				registry := map[string]gateways.PaymentGateway{
					"US": mock,
				}
				router := pipeline.NewRouterHandler(registry)
				return router, mock
			},
			expectedErr: nil,
		},
		{
			name:    "Unsupported Region",
			request: &domain.PaymentRequest{Region: "EU", Amount: 100},
			setup: func() (*pipeline.RouterHandler, *MockGateway) {
				mock := &MockGateway{}
				registry := map[string]gateways.PaymentGateway{
					"US":mock,
				}
				router := pipeline.NewRouterHandler(registry)
				return router, mock
			},
			expectedErr: domain.ErrRegionNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router, mock := tc.setup()
			err := router.Execute(context.Background(), tc.request)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("Expected %v, got %v",tc.expectedErr,err)
			}

			if tc.expectedErr == nil && !mock.Called{
				t.Errorf("Expected the gateway to be called but it was not")
			}
		})
	}
}
