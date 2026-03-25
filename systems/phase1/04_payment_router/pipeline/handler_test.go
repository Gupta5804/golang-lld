package pipeline_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Gupta5804/golang-lld/systems/phase1/04_payment_router/domain"
	"github.com/Gupta5804/golang-lld/systems/phase1/04_payment_router/pipeline"
)

type SpyHandler struct {
	next       pipeline.Handler
	Name       string
	Called     bool
	ShouldFail bool
}

func (h *SpyHandler) SetNext(next pipeline.Handler) pipeline.Handler {
	h.next = next
	return next
}
func (h *SpyHandler) Execute(ctx context.Context, req *domain.PaymentRequest) error {
	h.Called = true
	if h.ShouldFail {
		return fmt.Errorf("%s failed", h.Name)
	}
	if h.next != nil {
		return h.next.Execute(ctx, req)
	}
	return nil
}

func TestPipeline_ShortCircuiting(t *testing.T) {
	type Spies struct {
		validator *SpyHandler
		fraud     *SpyHandler
		router    *SpyHandler
	}
	tests := []struct {
		name          string
		setup         func() (pipeline.Handler, Spies)
		expectedError bool
		assertSpies   func(t *testing.T, spies Spies)
	}{
		{
			name: "Happy Path(All called)",
			setup: func() (pipeline.Handler, Spies) {
				spies := Spies{
					validator: &SpyHandler{Name: "Validator"},
					fraud:     &SpyHandler{Name: "Fraud"},
					router:    &SpyHandler{Name: "Router"},
				}
				spies.validator.SetNext(spies.fraud).SetNext(spies.router)

				return spies.validator, spies
			},
			expectedError: false,
			assertSpies: func(t *testing.T, spies Spies) {
				if !spies.validator.Called {
					t.Errorf("Expected Validator to be called")
				}
				if !spies.fraud.Called {
					t.Errorf("Expected fraud to be called")
				}
				if !spies.router.Called {
					t.Errorf("Expected router to be called")
				}
			},
		},
		{
			name: "Validation Failure",
			setup: func() (pipeline.Handler, Spies) {
				spies := Spies{
					validator: &SpyHandler{Name: "Validator",ShouldFail:true},
					fraud:     &SpyHandler{Name: "Fraud"},
					router:    &SpyHandler{Name: "Router"},
				}
				spies.validator.SetNext(spies.fraud).SetNext(spies.router)

				return spies.validator, spies
			},
			expectedError: true,
			assertSpies: func(t *testing.T, spies Spies) {
				if !spies.validator.Called {
					t.Errorf("Expected Validator to be called")
				}
				if spies.fraud.Called {
					t.Errorf("Expected fraud not to be called")
				}
				if spies.router.Called {
					t.Errorf("Expected router not to be called")
				}
			},
		},
		{
			name: "Fraud Failure",
			setup: func() (pipeline.Handler, Spies) {
				spies := Spies{
					validator: &SpyHandler{Name: "Validator"},
					fraud:     &SpyHandler{Name: "Fraud",ShouldFail:true},
					router:    &SpyHandler{Name: "Router"},
				}
				spies.validator.SetNext(spies.fraud).SetNext(spies.router)

				return spies.validator, spies
			},
			expectedError: true,
			assertSpies: func(t *testing.T, spies Spies) {
				if !spies.validator.Called {
					t.Errorf("Expected Validator to be called")
				}
				if !spies.fraud.Called {
					t.Errorf("Expected fraud to be called")
				}
				if spies.router.Called {
					t.Errorf("Expected router not to be called")
				}
			},
		},
	}

	for _,tc := range tests {
		t.Run(tc.name,func(t *testing.T){
			head, spies := tc.setup()
			err := head.Execute(context.Background(),&domain.PaymentRequest{})
		
			if tc.expectedError && err == nil{
				t.Fatalf("Expected error but got nil")
			}
			if !tc.expectedError && err != nil{
				t.Fatalf("Expected no error,got %v",err)
			}
			tc.assertSpies(t,spies)
		})
	}
}
