package buffs

import "github.com/Gupta5804/golang-lld/systems/phase1/08_game_avatar_loadout/domain"

type SpeedBoots struct {
	Provider domain.StatsProvider
}

func (s *SpeedBoots) CalculateStats() domain.Stats {
	avatarStats := s.Provider.CalculateStats()
	avatarStats.Speed = avatarStats.Speed + 5
	return avatarStats
}

type HeavyArmor struct {
	Provider domain.StatsProvider
}

func (h *HeavyArmor) CalculateStats() domain.Stats {
	avatarStats := h.Provider.CalculateStats()
	avatarStats.Health = avatarStats.Health + 50 
	avatarStats.Speed = avatarStats.Speed - 5
	return avatarStats
}

type BerserkerPotion struct {
	Provider domain.StatsProvider
}

func (b *BerserkerPotion) CalculateStats() domain.Stats {
	avatarStats := b.Provider.CalculateStats()
	avatarStats.Damage = avatarStats.Damage * 1.2
	avatarStats.Health = avatarStats.Health * 0.9
	return avatarStats
}
