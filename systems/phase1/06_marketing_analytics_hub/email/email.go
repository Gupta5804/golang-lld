package email

import (
	"context"
	"time"
)

type EmailService struct{}

func (s *EmailService) GetOpenRates(ctx context.Context) (float64, error) {
	time.Sleep(20 * time.Millisecond)
	return 568.63,nil
}
