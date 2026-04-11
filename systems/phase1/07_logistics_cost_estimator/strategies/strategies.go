package strategies

type StandardStrategy struct {
}

func (s *StandardStrategy) Calculate(weightKG float64, distanceKM int) int64 {
	cost := int64(weightKG * float64(distanceKM) * 5)
	return cost
}

type HolidayStrategy struct {
}

func (s *HolidayStrategy) Calculate(weightKG float64, distanceKM int) int64 {
	cost := int64(weightKG * float64(distanceKM) * 8)
	return cost
}
