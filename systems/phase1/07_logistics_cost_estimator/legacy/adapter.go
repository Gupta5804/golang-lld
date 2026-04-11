package legacy

import (
	"context"
	"time"
)

type LegacyAdapter struct {
	externalAPI *RegionalCourierAPI
}

func (a *LegacyAdapter) Fetch(ctx context.Context, weightKG float64, distanceKM int) int64 {
	pounds := float32(weightKG * 2.204)
	miles := int(float64(distanceKM) * 0.621)

	responseChan := make(chan float32, 1)
	go func() {
		val := a.externalAPI.CalculateLastMile(pounds, miles)
		responseChan <- val
	}()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	select {
	case val := <-responseChan:
		return int64(val * 100)
	case <-ctx.Done():
		return 0
	}
}
func NewLegacyAdapter(api *RegionalCourierAPI) *LegacyAdapter {
	return &LegacyAdapter{externalAPI: api}
}
type RegionalCourierAPI struct{}

func (api *RegionalCourierAPI) CalculateLastMile(pounds float32, miles int) float32 {
	// This will just be a dummy implementation returning a flat rate later
	return 50
}
