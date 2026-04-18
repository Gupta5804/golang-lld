package domain

type Avatar struct {
	ID       string
	VIPFlags []string
	Stats    *Stats
}

func (a *Avatar) CalculateStats() Stats {
	return *a.Stats
}

type Stats struct {
	Damage float64
	Health float64
	Speed  float64
}

type StatsProvider interface {
	CalculateStats() Stats
}
