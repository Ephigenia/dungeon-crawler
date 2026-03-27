package game

import "fmt"

// updateEnemies ticks AI for every living enemy once per call.
func (g *Game) updateEnemies() {
	for _, e := range g.enemies {
		if !e.IsAlive() {
			continue
		}
		e.moveTick--
		if e.moveTick > 0 {
			continue
		}
		e.moveTick = enemyMoveInterval

		dist := iabs(g.player.X-e.X) + iabs(g.player.Y-e.Y)
		if dist <= enemyChaseRange {
			e.state = enemyStateChase
		} else {
			e.state = enemyStateIdle
		}

		switch e.state {
		case enemyStateChase:
			g.enemyChaseMove(e)
		case enemyStateIdle:
			g.enemyWander(e)
		}
	}
}

// enemyChaseMove steps the enemy one tile toward the player.
// Tries the dominant axis first, then the other as a fallback.
// If the target tile is the player, attacks instead of moving.
func (g *Game) enemyChaseMove(e *Enemy) {
	dx := g.player.X - e.X
	dy := g.player.Y - e.Y

	var moves [2][2]int
	n := 0
	if iabs(dx) >= iabs(dy) {
		if dx != 0 {
			moves[n] = [2]int{isign(dx), 0}
			n++
		}
		if dy != 0 {
			moves[n] = [2]int{0, isign(dy)}
			n++
		}
	} else {
		if dy != 0 {
			moves[n] = [2]int{0, isign(dy)}
			n++
		}
		if dx != 0 {
			moves[n] = [2]int{isign(dx), 0}
			n++
		}
	}

	for i := 0; i < n; i++ {
		nx, ny := e.X+moves[i][0], e.Y+moves[i][1]
		if nx == g.player.X && ny == g.player.Y {
			g.resolveEnemyAttack(e)
			return
		}
		if g.dungeon.IsWalkable(nx, ny) && g.enemyAt(nx, ny) == nil {
			if o := g.objectAt(nx, ny); o == nil || o.Type.PassableByEnemy {
				e.X, e.Y = nx, ny
				return
			}
		}
	}
}

// enemyWander moves the enemy to a random adjacent walkable tile.
func (g *Game) enemyWander(e *Enemy) {
	dirs := [4][2]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	g.rng.Shuffle(len(dirs), func(i, j int) { dirs[i], dirs[j] = dirs[j], dirs[i] })
	for _, d := range dirs {
		nx, ny := e.X+d[0], e.Y+d[1]
		if g.dungeon.IsWalkable(nx, ny) && g.enemyAt(nx, ny) == nil &&
			!(nx == g.player.X && ny == g.player.Y) {
			if o := g.objectAt(nx, ny); o == nil || o.Type.PassableByEnemy {
				e.X, e.Y = nx, ny
				return
			}
		}
	}
}

// counterChance returns the probability (0–100) that the player can counter
// an incoming enemy attack. Base 30% + 2% per agility point, capped at 80%.
func (g *Game) counterChance() int {
	chance := 30 + g.player.EffectiveAgility()*2
	if chance > 80 {
		chance = 80
	}
	return chance
}

// resolveEnemyAttack applies damage from an enemy to the player, then gives
// the player a chance to counter-attack based on agility.
func (g *Game) resolveEnemyAttack(e *Enemy) {
	// Enemy attacks player.
	hpBefore := g.player.HP
	g.player.TakeDamage(e.Type.Attack, g.rng)
	dmg := hpBefore - g.player.HP

	chance := g.counterChance()
	countered := g.rng.Intn(100) < chance

	if countered {
		// Player counter-attacks.
		counterDmg := calcPlayerDamage(
			g.player.EffectiveAttack(), g.player.WeaponPower(), g.player.WeaponSpeed(),
			g.player.EffectiveAgility(), g.player.Level, e.Type.Defense, g.rng,
		)
		e.HP -= counterDmg
		if e.HP < 0 {
			e.HP = 0
		}

		if e.IsAlive() {
			g.combatLines = []string{
				fmt.Sprintf("%s attacks you for %d damage!", e.Type.Name, dmg),
				fmt.Sprintf("You counter for %d damage! (%d%% chance)", counterDmg, chance),
				fmt.Sprintf("%s: %d / %d HP", e.Type.Name, e.HP, e.Type.MaxHP),
				fmt.Sprintf("You: %d / %d HP", g.player.HP, g.player.EffectiveMaxHP()),
			}
		} else {
			g.player.AddEXP(20)
			g.combatLines = []string{
				fmt.Sprintf("%s attacks you for %d damage!", e.Type.Name, dmg),
				fmt.Sprintf("You counter for %d — %s defeated! (%d%% chance)", counterDmg, e.Type.Name, chance),
				fmt.Sprintf("You: %d / %d HP", g.player.HP, g.player.EffectiveMaxHP()),
			}
			dropChance := 10 + (g.player.Level - 1)
			if g.rng.Intn(100) < dropChance {
				g.potions = append(g.potions, newPotion(e.X, e.Y, g.rng))
			}
		}
	} else {
		g.combatLines = []string{
			fmt.Sprintf("%s attacks you for %d damage! (no counter, %d%% chance)", e.Type.Name, dmg, chance),
			fmt.Sprintf("You: %d / %d HP", g.player.HP, g.player.EffectiveMaxHP()),
		}
	}
	g.combatFrames = 120
}

func isign(x int) int {
	if x > 0 {
		return 1
	}
	return -1
}

func iabs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
