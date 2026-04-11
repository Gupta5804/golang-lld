package estimator

import (
	"context"

	"github.com/Gupta5804/golang-lld/systems/phase1/07_logistics_cost_estimator/domain"
)

type PricingStrategy interface { // strategy interface
	Calculate(weightKG float64, distanceKM int) int64
}
type RateFetcher interface { //adapter interface
	Fetch(ctx context.Context, weightKG float64, distanceKM int) int64
}
type EstimatorService struct {
	fetcher RateFetcher
}

func (s *EstimatorService) CalculateTotal(ctx context.Context, req domain.Shipment, strategy PricingStrategy) (*domain.Quote, error) {

	quote := &domain.Quote{
		LineHaulCost: strategy.Calculate(req.WeightKG, req.DistanceKM),
		LastMileCost: s.fetcher.Fetch(ctx, req.WeightKG, req.DistanceKM),
	}
	return quote, nil
}
func NewEstimatorService(fetcher RateFetcher) *EstimatorService {
	return &EstimatorService{fetcher: fetcher}
}
