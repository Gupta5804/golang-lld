package estimator_test

import (
	"context"
	"sync"
	"testing"

	"github.com/Gupta5804/golang-lld/systems/phase1/07_logistics_cost_estimator/domain"
	"github.com/Gupta5804/golang-lld/systems/phase1/07_logistics_cost_estimator/estimator"
	"github.com/Gupta5804/golang-lld/systems/phase1/07_logistics_cost_estimator/strategies"
)

type SpyRateFetcher struct {
	CallCount    int
	LastWeight   float64
	LastDistance int
	MockReturn   int64
	mu           sync.Mutex
}

func (s *SpyRateFetcher) Fetch(ctx context.Context, weightKG float64, distanceKM int) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.CallCount++
	s.LastWeight = weightKG
	s.LastDistance = distanceKM
	return s.MockReturn
}
func TestEstimator_CalculateTotal(t *testing.T) {
	tests := []struct {
		name             string
		setup            func() (*estimator.EstimatorService, estimator.PricingStrategy, *SpyRateFetcher)
		shipment         domain.Shipment
		expectedLineHaul int64
		expectedLastMile int64
		expectedTotal    int64
	}{
		{
			name: "Standard Strategy",
			setup: func() (*estimator.EstimatorService, estimator.PricingStrategy, *SpyRateFetcher) {
				spy := &SpyRateFetcher{MockReturn: 100}
				service := estimator.NewEstimatorService(spy)
				strategy := &strategies.StandardStrategy{}
				return service, strategy, spy
			},
			shipment: domain.Shipment{
				WeightKG:   10.0,
				DistanceKM: 10,
			},
			expectedLineHaul: 500,
			expectedLastMile: 100,
			expectedTotal:    600,
		},
		{
			name: "Holiday Strategy",
			setup: func() (*estimator.EstimatorService, estimator.PricingStrategy, *SpyRateFetcher) {
				spy := &SpyRateFetcher{MockReturn: 150}
				service := estimator.NewEstimatorService(spy)
				strategy := &strategies.HolidayStrategy{}
				return service, strategy, spy
			},
			shipment: domain.Shipment{
				WeightKG:   10.0,
				DistanceKM: 10,
			},
			expectedLineHaul: 800,
			expectedLastMile: 150,
			expectedTotal:    950,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service, strategy, spy := tc.setup()
			ctx := context.Background()
			quote, err := service.CalculateTotal(ctx, tc.shipment, strategy)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if quote == nil {
				t.Fatalf("Expected quote to not be nil")
			}
			// Assert Spy interactions (Did the estimator pass the right data to the port?)
			if spy.CallCount != 1 {
				t.Errorf("Expected RateFetcher to be called exactly 1 time, got %d", spy.CallCount)
			}
			if spy.LastWeight != tc.shipment.WeightKG || spy.LastDistance != tc.shipment.DistanceKM {
				t.Errorf("Expected spy to receive %f KG and %d KM, got %f KG and %d KM",
					tc.shipment.WeightKG, tc.shipment.DistanceKM, spy.LastWeight, spy.LastDistance)
			}

			// Assert Anatomy
			if quote.LineHaulCost != tc.expectedLineHaul {
				t.Errorf("Expected LineHaul %d, got %d", tc.expectedLineHaul, quote.LineHaulCost)
			}
			if quote.LastMileCost != tc.expectedLastMile {
				t.Errorf("Expected LastMile %d, got %d", tc.expectedLastMile, quote.LastMileCost)
			}

			if actualTotal := quote.TotalCost(); actualTotal != tc.expectedTotal {
				t.Errorf("Expected total to be %d, got %d", tc.expectedTotal, actualTotal)
			}
		})
	}
}

// concurrency test
func TestEstimator_Concurrency(t *testing.T) {
	spy := &SpyRateFetcher{MockReturn: 200}
	service := estimator.NewEstimatorService(spy)

	standardStrategy := &strategies.StandardStrategy{}
	holidayStrategy := &strategies.HolidayStrategy{}
	shipment := domain.Shipment{WeightKG: 50, DistanceKM: 100}

	workers := 1000

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(workerID int) {
			defer wg.Done()
			ctx := context.Background()

			var strat estimator.PricingStrategy
			if workerID%2 == 0 {
				strat = standardStrategy
			} else {
				strat = holidayStrategy
			}

			_, err := service.CalculateTotal(ctx, shipment, strat)
			if err != nil {
				t.Errorf("Worker %d failed: %v", workerID, err)
			}
		}(i)
	}
	wg.Wait()

	if spy.CallCount != workers {
		t.Errorf("Expected %d total calls to legacy adapter, got %d", workers, spy.CallCount)
	}
}
