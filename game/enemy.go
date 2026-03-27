package game

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// EnemyType is the static definition of an enemy kind — its name and base stats.
// Instances on the map are represented by Enemy, which holds a pointer to its type
// plus mutable runtime state (current HP, position).
type EnemyType struct {
	Name         string
	MaxHP        int
	Attack       int
	Defense      int
	MoveInterval int           // frames between moves; lower = faster
	VisionRange  int           // Manhattan distance at which the enemy starts chasing
	ImagePath    string        // asset path for the sprite; empty = use color fallback
	Image        *ebiten.Image // loaded at startup from ImagePath; nil until then
}

type enemyState int

const (
	enemyStateIdle  enemyState = iota
	enemyStateChase
)


// Enemy is a live enemy on the map.
type Enemy struct {
	X, Y     int
	HP       int
	Type     *EnemyType
	state    enemyState
	moveTick int
}

// IsAlive returns true if the enemy has HP remaining.
func (e *Enemy) IsAlive() bool {
	return e.HP > 0
}

// TakeDamage reduces HP using the shared damage formula.
func (e *Enemy) TakeDamage(attack int, rng *rand.Rand) {
	dmg := calcDamage(attack, e.Type.Defense, rng)
	e.HP -= dmg
	if e.HP < 0 {
		e.HP = 0
	}
}

// spawnEnemy creates a live Enemy from a randomly chosen EnemyType.
func spawnEnemy(x, y int, rng *rand.Rand) *Enemy {
	t := AllEnemyTypes[rng.Intn(len(AllEnemyTypes))]
	return &Enemy{X: x, Y: y, HP: t.MaxHP, Type: t}
}

// calcPlayerDamage computes damage for a player attack, incorporating weapon
// power, weapon speed, agility, level, and a level-scaled random bonus.
//
//	weaponContrib   = weaponPower × (1 + weaponSpeed × agility / 100)
//	effectiveAttack = (baseAttack + weaponContrib) × (1 + (level−1) × 0.05)
//	randomBonus     = rng(0 … level×2)
//	damage          = max(0, int(effectiveAttack) − defense + randomBonus)
func calcPlayerDamage(baseAttack, weaponPower, weaponSpeed, agility, level, defense int, rng *rand.Rand) int {
	weaponContrib := float64(weaponPower) * (1.0 + float64(weaponSpeed)*float64(agility)/100.0)
	effective := (float64(baseAttack) + weaponContrib) * (1.0 + float64(level-1)*0.05)
	randomBonus := rng.Intn(level*2 + 1)
	dmg := int(effective) - defense + randomBonus
	if dmg < 0 {
		dmg = 0
	}
	return dmg
}

// calcDamage returns the damage dealt given an attack and defense value.
// The attack/defense ratio is used as a factor, so higher attack relative
// to defense amplifies damage. A random bonus of 0…attack/2 is added.
// Result is clamped to 0 (no minimum of 1).
func calcDamage(attack, defense int, rng *rand.Rand) int {
	randomBonus := rng.Intn(attack/2 + 1)
	if defense <= 0 {
		return attack + randomBonus
	}
	factor := float64(attack) / float64(defense)
	dmg := int(float64(attack-defense)*factor) + randomBonus
	if dmg < 0 {
		dmg = 0
	}
	return dmg
}
