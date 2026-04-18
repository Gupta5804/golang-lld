package avatar

import "github.com/Gupta5804/golang-lld/systems/phase1/08_game_avatar_loadout/domain"

type Option func(*domain.Avatar)

func NewAvatar(id string, opts ...Option) *domain.Avatar {
	stats := &domain.Stats{
		Health: 100.0,
		Damage: 10.0,
		Speed:  10.0,
	}
	var flags []string
	a := &domain.Avatar{
		ID:       id,
		VIPFlags: flags,
		Stats:    stats,
	}

	for _, opt := range opts {
		opt(a)
	}
	return a
}
