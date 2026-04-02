package report

import (
	"context"
	"errors"

	"github.com/Gupta5804/golang-lld/systems/phase1/06_marketing_analytics_hub/domain"
)

// strategy interface
type ReportGenerator interface {
	Generate(ctx context.Context, email Email, crm CRM, billing Billing) (*domain.Report, error)
}

// simple factory
func NewReportGenerator(reportType domain.ReportType) (ReportGenerator, error) {
	switch reportType {
	case domain.WeeklyEngagement:
		return &weeklyGenerator{}, nil
	case domain.MonthlyRevenue:
		return &monthlyGenerator{}, nil
	case domain.QuarterlyChurn:
		return &quarterlyGenerator{}, nil
	default:
		return nil, errors.New("unknown ReportType")
	}
}
