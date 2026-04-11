package domain

type Shipment struct {
	WeightKG   float64
	DistanceKM int
}

type Quote struct {
	LineHaulCost int64
	LastMileCost int64
}

func (q *Quote) TotalCost() int64 {
	return q.LineHaulCost + q.LastMileCost
}


