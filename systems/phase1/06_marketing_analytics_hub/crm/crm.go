package crm

import (
	"context"
	"time"
)

type CRMService struct{}

func (s *CRMService) GetActiveUsers(ctx context.Context) (int, error) {
	time.Sleep(20 * time.Millisecond)
	return 56, nil
}
func (s *CRMService) GetChurnedUsers(ctx context.Context) (int, error) {
	time.Sleep(20 * time.Millisecond)
	return 200,nil
}
