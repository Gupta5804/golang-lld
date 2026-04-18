package avatar

import "github.com/Gupta5804/golang-lld/systems/phase1/08_game_avatar_loadout/domain"

func WithBaseHealth(health float64) Option {
	return func(avatar *domain.Avatar) {
		avatar.Stats.Health = health
	}
}
func WithVIPFlags(flags ...string) Option {
	return func(avatar *domain.Avatar){
		avatar.VIPFlags = flags 
	}
}
func WithBaseDamage(dmg float64) Option {
	return func(avatar *domain.Avatar){
		avatar.Stats.Damage = dmg
	}
}
func WithBaseSpeed(spd float64) Option {
	return func(avatar *domain.Avatar){
		avatar.Stats.Speed = spd
	}
}
