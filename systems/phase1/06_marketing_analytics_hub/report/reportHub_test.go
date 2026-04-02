package report_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Gupta5804/golang-lld/systems/phase1/06_marketing_analytics_hub/domain"
	"github.com/Gupta5804/golang-lld/systems/phase1/06_marketing_analytics_hub/report"
)

type MockCRM struct {
	CountGetActiveUsersCalled  int
	CountGetChurnedUsersCalled int
	ShouldFail                 bool
}

func (m *MockCRM) GetActiveUsers(ctx context.Context) (int, error) {
	m.CountGetActiveUsersCalled++
	if m.ShouldFail {
		return 0, errors.New("503 Service Unavailable: CRM is down")
	}
	return 100, nil
}
func (m *MockCRM) GetChurnedUsers(ctx context.Context) (int, error) {
	m.CountGetChurnedUsersCalled++
	if m.ShouldFail {
		return 0, errors.New("503 Service Unavailable: CRM is down")
	}
	return 10, nil
}

type MockEmail struct {
	CountCalled int
}

func (m *MockEmail) GetOpenRates(ctx context.Context) (float64, error) {
	m.CountCalled++
	return 0, nil
}

type MockBilling struct {
	CountCalled int
}

func (m *MockBilling) GetTotalRevenue(ctx context.Context) (int64, error) {
	m.CountCalled++
	return 0, nil
}

func TestReportHub_GenerateReport(t *testing.T) {
	type SpyMocks struct {
		email   *MockEmail
		billing *MockBilling
		crm     *MockCRM
	}
	tests := []struct {
		name                 string
		setup                func() (*report.ReportHub, *SpyMocks)
		request              domain.ReportRequest
		expectedEmailCalls   int
		expectedBillingCalls int
		expectedCRMActive    int
		expectedCRMChurn     int
		expectError          bool
	}{
		{
			name: "Weekly Engagement Report(CRM Active + Email)",
			setup: func() (*report.ReportHub, *SpyMocks) {
				spies := &SpyMocks{
					email:   &MockEmail{},
					billing: &MockBilling{},
					crm:     &MockCRM{},
				}
				hub := report.NewReportHub(spies.crm, spies.email, spies.billing)
				return hub, spies
			},
			request: domain.ReportRequest{
				ReportType: domain.WeeklyEngagement,
			},
			expectedEmailCalls:   1,
			expectedBillingCalls: 0,
			expectedCRMActive:    1,
			expectedCRMChurn:     0,
			expectError:          false,
		},
		{
			name: "Monthly Revenue Report(CRM Active + Billing)",
			setup: func() (*report.ReportHub, *SpyMocks) {
				spies := &SpyMocks{
					email:   &MockEmail{},
					billing: &MockBilling{},
					crm:     &MockCRM{},
				}
				hub := report.NewReportHub(spies.crm, spies.email, spies.billing)
				return hub, spies
			},
			request: domain.ReportRequest{
				ReportType: domain.MonthlyRevenue,
			},
			expectedEmailCalls:   0,
			expectedBillingCalls: 1,
			expectedCRMActive:    1,
			expectedCRMChurn:     0,
			expectError:          false,
		},
		{
			name: "Quarterly Churn Report(CRM Churn)",
			setup: func() (*report.ReportHub, *SpyMocks) {
				spies := &SpyMocks{
					email:   &MockEmail{},
					billing: &MockBilling{},
					crm:     &MockCRM{},
				}
				hub := report.NewReportHub(spies.crm, spies.email, spies.billing)
				return hub, spies
			},
			request: domain.ReportRequest{
				ReportType: domain.QuarterlyChurn,
			},
			expectedEmailCalls:   0,
			expectedBillingCalls: 0,
			expectedCRMActive:    0,
			expectedCRMChurn:     1,
			expectError:          false,
		},
		{
			name: "Unsupported Report Type(Error)",
			setup: func() (*report.ReportHub, *SpyMocks) {
				spies := &SpyMocks{
					email:   &MockEmail{},
					billing: &MockBilling{},
					crm:     &MockCRM{},
				}
				hub := report.NewReportHub(spies.crm, spies.email, spies.billing)
				return hub, spies
			},
			request: domain.ReportRequest{
				ReportType: 4,
			},
			expectedEmailCalls:   0,
			expectedBillingCalls: 0,
			expectedCRMActive:    0,
			expectedCRMChurn:     0,
			expectError:          true,
		},
		{
			name: "Subsystem Failure(CRM down during weekly report)",
			setup: func() (*report.ReportHub, *SpyMocks) {
				spies := &SpyMocks{
					email:   &MockEmail{},
					billing: &MockBilling{},
					crm:     &MockCRM{ShouldFail: true},
				}
				hub := report.NewReportHub(spies.crm, spies.email, spies.billing)
				return hub, spies
			},
			request: domain.ReportRequest{
				ReportType: domain.WeeklyEngagement,
			},
			expectedEmailCalls:   0,
			expectedBillingCalls: 0,
			expectedCRMActive:    1,
			expectedCRMChurn:     0,
			expectError:          true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hub, spies := tc.setup()
			_, err := hub.GenerateReport(context.Background(), tc.request)
			if err != nil && !tc.expectError {
				t.Fatalf("Expected no error, got %v", err)
			}
			if err == nil && tc.expectError {
				t.Fatal("Expected error, got nil")
			}
			actualEmailCalls := spies.email.CountCalled
			actualBillingCalls := spies.billing.CountCalled
			actualCRMActive := spies.crm.CountGetActiveUsersCalled
			actualCRMChurn := spies.crm.CountGetChurnedUsersCalled
			if actualEmailCalls != tc.expectedEmailCalls {
				t.Errorf("Expected %d email calls, got %d", tc.expectedEmailCalls, actualEmailCalls)
			}
			if actualBillingCalls != tc.expectedBillingCalls {
				t.Errorf("Expected %d billing calls, got %d", tc.expectedBillingCalls, actualBillingCalls)
			}
			if actualCRMActive != tc.expectedCRMActive {
				t.Errorf("Expected %d CRM active calls, got %d", tc.expectedCRMActive, actualCRMActive)
			}
			if actualCRMChurn != tc.expectedCRMChurn {
				t.Errorf("Expected %d CRM churn calls, got %d", tc.expectedCRMChurn, actualCRMChurn)
			}

		})
	}
}
