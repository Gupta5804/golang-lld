package buffs_test

import (
	"sync"
	"testing"

	"github.com/Gupta5804/golang-lld/systems/phase1/08_game_avatar_loadout/avatar"
	"github.com/Gupta5804/golang-lld/systems/phase1/08_game_avatar_loadout/buffs"
	"github.com/Gupta5804/golang-lld/systems/phase1/08_game_avatar_loadout/domain"
)

func TestCalculateStats_Scenarios(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() domain.StatsProvider
		expectedStats domain.Stats
	}{
		{
			name: "Base Avatar + Speed-boots",
			setup: func() domain.StatsProvider {
				avatar := avatar.NewAvatar("avatar_01")
				modifiedAvatar := &buffs.SpeedBoots{Provider: avatar}
				return modifiedAvatar
			},
			expectedStats: domain.Stats{
				Health: 100.0,
				Damage: 10.0,
				Speed:  15.0,
			},
		},
		{
			name: "Base Avatar + Heavy Armor",
			setup: func() domain.StatsProvider {
				avatar := avatar.NewAvatar("avatar_01")
				modifiedAvatar := &buffs.HeavyArmor{Provider: avatar}
				return modifiedAvatar
			},
			expectedStats: domain.Stats{
				Health: 150.0,
				Damage: 10.0,
				Speed:  5.0,
			},
		},
		{
			name: "Base Avatar + Berserker Potion",
			setup: func() domain.StatsProvider {
				avatar := avatar.NewAvatar("avatar_01")
				modifiedAvatar := &buffs.BerserkerPotion{Provider: avatar}
				return modifiedAvatar
			},
			expectedStats: domain.Stats{
				Health: 90.0,
				Damage: 12.0,
				Speed:  10.0,
			},
		},
		{
			name: "Base Avatar + Heavy Armor + Berserker Potion",
			setup: func() domain.StatsProvider {
				avatar := avatar.NewAvatar("avatar_01")
				heavyArmorAvatar := &buffs.HeavyArmor{Provider: avatar}
				modifiedAvatar := &buffs.BerserkerPotion{Provider: heavyArmorAvatar}
				return modifiedAvatar
			},
			expectedStats: domain.Stats{
				Health: 135.0,
				Damage: 12.0,
				Speed:  5.0,
			},
		},
		{
			name: "Base Avatar + Berserker Potion + Heavy Armor",
			setup: func() domain.StatsProvider {
				avatar := avatar.NewAvatar("avatar_01")
				berserkedAvatar := &buffs.BerserkerPotion{Provider: avatar}
				modifiedAvatar := &buffs.HeavyArmor{Provider: berserkedAvatar}
				return modifiedAvatar
			},
			expectedStats: domain.Stats{
				Health: 140.0,
				Damage: 12.0,
				Speed:  5.0,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			statsProvider := tc.setup()
			actualStats := statsProvider.CalculateStats()

			if actualStats.Health != tc.expectedStats.Health {
				t.Errorf("Expected %f health, got %f", tc.expectedStats.Health, actualStats.Health)
			}
			if actualStats.Damage != tc.expectedStats.Damage {
				t.Errorf("Expected %f Damage, got %f", tc.expectedStats.Damage, actualStats.Damage)
			}
			if actualStats.Speed != tc.expectedStats.Speed {
				t.Errorf("Expected %f Speed, got %f", tc.expectedStats.Speed, actualStats.Speed)
			}

		})
	}
}

func TestCalculateStats_Concurrency(t *testing.T) {
	avatar := avatar.NewAvatar("avatar_01")
	armorAvatar := &buffs.HeavyArmor{Provider: avatar}
	speedAvatar := &buffs.SpeedBoots{Provider: armorAvatar}
	modifiedAvatar := &buffs.BerserkerPotion{Provider: speedAvatar}

	workers := 1000
	var wg sync.WaitGroup
	wg.Add(workers)
	for i:=0;i<workers;i++{
		go func(){
			defer wg.Done()
			modifiedAvatar.CalculateStats()
		}()
	}

	wg.Wait()

}
