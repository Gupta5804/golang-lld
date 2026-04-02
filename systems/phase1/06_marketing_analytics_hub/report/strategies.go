package report

import (
	"context"
	"fmt"

	"github.com/Gupta5804/golang-lld/systems/phase1/06_marketing_analytics_hub/domain"
)

type weeklyGenerator struct{}

func (g *weeklyGenerator) Generate(ctx context.Context, email Email, crm CRM, billing Billing) (*domain.Report, error) {
	activeUsers, err := crm.GetActiveUsers(ctx)
	if err != nil {
		return nil, err
	}
	openRates, err := email.GetOpenRates(ctx)
	if err != nil {
		return nil, err
	}
	payload := make(map[string]string)
	payload["active_users"] = fmt.Sprintf("%d", activeUsers)
	payload["email_open_rate"] = fmt.Sprintf("%f", openRates)
	return &domain.Report{
		Title:   "Weekly Engagement Report",
		Payload: payload,
	}, nil
}

type monthlyGenerator struct{}

func (g *monthlyGenerator) Generate(ctx context.Context, email Email, crm CRM, billing Billing) (*domain.Report, error) {
	activeUsers, err := crm.GetActiveUsers(ctx)
	if err != nil {
		return nil, err
	}
	totalRevenue, err := billing.GetTotalRevenue(ctx)
	if err != nil {
		return nil, err
	}
	payload := make(map[string]string)
	payload["active_users"] = fmt.Sprintf("%d", activeUsers)
	payload["total_revenue"] = fmt.Sprintf("%d", totalRevenue)
	return &domain.Report{
		Title:   "Monthly Revenue Report",
		Payload: payload,
	}, nil
}

type quarterlyGenerator struct{}

func (g *quarterlyGenerator) Generate(ctx context.Context, email Email, crm CRM, billing Billing) (*domain.Report, error) {
	churnedUsers, err := crm.GetChurnedUsers(ctx)
	if err != nil {
		return nil, err
	}
	payload := make(map[string]string)
	payload["churned_users"] = fmt.Sprintf("%d", churnedUsers)
	return &domain.Report{
		Title:   "Quarterly Churn report",
		Payload: payload,
	}, nil
}
