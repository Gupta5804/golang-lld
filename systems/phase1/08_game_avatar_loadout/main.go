package main

import (
	"fmt"

	"github.com/Gupta5804/golang-lld/systems/phase1/08_game_avatar_loadout/avatar"
	"github.com/Gupta5804/golang-lld/systems/phase1/08_game_avatar_loadout/buffs"
	"github.com/Gupta5804/golang-lld/systems/phase1/08_game_avatar_loadout/domain"
)

func main() {
	fmt.Println("=== 🛡️ Game Avatar Loadout System ===")

	// ---------------------------------------------------------
	// Phase 1: The Creational Pattern (Functional Options)
	// ---------------------------------------------------------
	fmt.Println("\n[Phase 1] Spawning 'Juggernaut VIP' Starter Bundle...")
	
	// Notice how clean this API is. We don't have a constructor with 10 nil parameters.
	baseAvatar := avatar.NewAvatar("player_99_juggernaut",
		avatar.WithBaseHealth(200.0), // Non-standard starting health
		avatar.WithBaseDamage(20.0),  // Non-standard starting damage
		avatar.WithVIPFlags("Premium_Tier_1", "Founder"),
	)

	fmt.Printf("-> Base Identity : ID: %s | VIP: %v\n", baseAvatar.ID, baseAvatar.VIPFlags)
	fmt.Printf("-> Base Stats    : %+v\n", baseAvatar.CalculateStats())


	// ---------------------------------------------------------
	// Phase 2: The Structural Pattern (Decorator)
	// ---------------------------------------------------------
	fmt.Println("\n[Phase 2] Match Started: Equipping Items...")

	// We cast the base struct to our Interface to begin the wrapping process.
	var currentLoadout domain.StatsProvider = baseAvatar

	fmt.Println("-> Player picks up [Heavy Armor] (+50 HP, -5 Speed)")
	currentLoadout = &buffs.HeavyArmor{Provider: currentLoadout}
	fmt.Printf("   Current Stats : %+v\n", currentLoadout.CalculateStats())

	fmt.Println("-> Player picks up [Speed Boots] (+5 Speed)")
	currentLoadout = &buffs.SpeedBoots{Provider: currentLoadout}
	fmt.Printf("   Current Stats : %+v\n", currentLoadout.CalculateStats())

	fmt.Println("-> Player drinks [Berserker Potion] (+20% Damage, -10% HP)")
	currentLoadout = &buffs.BerserkerPotion{Provider: currentLoadout}
	fmt.Printf("   Final Loadout : %+v\n", currentLoadout.CalculateStats())


	// ---------------------------------------------------------
	// Phase 3: The Immutability Proof
	// ---------------------------------------------------------
	fmt.Println("\n[Phase 3] Architectural Immutability Check...")
	fmt.Println("Checking if the base avatar in memory was accidentally mutated by the buffs...")
	
	baseStats := baseAvatar.CalculateStats()
	if baseStats.Health == 200.0 && baseStats.Damage == 20.0 {
		fmt.Printf("-> SUCCESS: Base Stats are completely untouched: %+v\n", baseStats)
	} else {
		fmt.Printf("-> FATAL: Base Stats were mutated! %+v\n", baseStats)
	}
	
	fmt.Println("\n=== Session Complete ===")
}