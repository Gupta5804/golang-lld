package billing

import (
	"context"
	"time"
)

type BillingService struct{}

func (s *BillingService) GetTotalRevenue(ctx context.Context) (int64, error) {
	time.Sleep(20 * time.Millisecond)
	return 2000000, nil
}
