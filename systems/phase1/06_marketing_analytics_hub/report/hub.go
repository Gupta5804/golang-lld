package report

import (
	"context"

	"github.com/Gupta5804/golang-lld/systems/phase1/06_marketing_analytics_hub/domain"
)

type CRM interface {
	GetActiveUsers(ctx context.Context) (int, error)
	GetChurnedUsers(ctx context.Context) (int, error)
}
type Email interface {
	GetOpenRates(ctx context.Context) (float64, error)
}
type Billing interface {
	GetTotalRevenue(ctx context.Context) (int64, error)
}
type ReportHub struct {
	crm     CRM
	email   Email
	billing Billing
}

func (h *ReportHub) GenerateReport(ctx context.Context, req domain.ReportRequest) (*domain.Report, error) {
	generator,err := NewReportGenerator(req.ReportType)
	if err != nil{
		return nil,err
	}
	return generator.Generate(ctx,h.email,h.crm,h.billing)
}
func NewReportHub(crm CRM, email Email, billing Billing) *ReportHub {
	return &ReportHub{
		crm:     crm,
		email:   email,
		billing: billing,
	}
}

