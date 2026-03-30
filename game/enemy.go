package game

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// EnemyType is the static definition of an enemy kind — its name and base stats.
// Instances on the map are represented by Enemy, which holds a pointer to its type
// plus mutable runtime state (current HP, position).
//
// Each state can have its own spritesheet; frame 0 is displayed. Unset paths are
// ignored and the renderer falls back to the first available image.
type EnemyType struct {
	Name         string
	MaxHP        int
	Attack       int
	Defense      int
	MoveInterval int // frames between moves; lower = faster
	VisionRange  int // Manhattan distance at which the enemy starts chasing
	AnimSpeed    int // ticks per animation frame; 0 = no animation (always show frame 0)

	IdleImagePath   string
	MoveImagePath   string
	AttackImagePath string
	DeathImagePath  string

	IdleSheet   *Spritesheet // spritesheet for idle state
	MoveSheet   *Spritesheet // spritesheet for move/chase state
	AttackSheet *Spritesheet // spritesheet for attack state
	DeathSheet  *Spritesheet // spritesheet for death state
}

// SheetForState returns the best available spritesheet for the given state,
// falling back through idle → move → attack → death until one is found.
func (et *EnemyType) SheetForState(state enemyState) *Spritesheet {
	switch state {
	case enemyStateDead:
		if et.DeathSheet != nil {
			return et.DeathSheet
		}
	case enemyStateChase:
		if et.MoveSheet != nil {
			return et.MoveSheet
		}
	}
	if et.IdleSheet != nil {
		return et.IdleSheet
	}
	if et.MoveSheet != nil {
		return et.MoveSheet
	}
	if et.AttackSheet != nil {
		return et.AttackSheet
	}
	return et.DeathSheet
}

// FrameForState returns the sprite at the given frame index from the best available
// spritesheet for the given state. The frame index is wrapped to the sheet length.
// Returns nil if no spritesheet is loaded for this type.
func (et *EnemyType) FrameForState(state enemyState, frame int) *ebiten.Image {
	sheet := et.SheetForState(state)
	if sheet.Len() == 0 {
		return nil
	}
	return sheet.Sprite(frame % sheet.Len())
}

type enemyState int

const (
	enemyStateIdle  enemyState = iota
	enemyStateChase
	enemyStateDead
)


// Enemy is a live enemy on the map.
type Enemy struct {
	X, Y      int
	HP        int
	Type      *EnemyType
	state     enemyState
	moveTick  int
	animFrame int // current animation frame index
	animTick  int // ticks elapsed since last frame advance
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

// calcEnemyDangerLevel computes a 1–4 danger rating for an enemy relative to the player.
// It estimates average hits-to-kill in both directions and uses the ratio as a danger score.
func calcEnemyDangerLevel(e *Enemy, p *Player) int {
	defense := p.EffectiveDefense()
	var avgEnemyDmg float64
	if defense <= 0 {
		avgEnemyDmg = float64(e.Type.Attack)
	} else {
		avgEnemyDmg = float64(e.Type.Attack) * float64(e.Type.Attack) / float64(e.Type.Attack+defense)
	}
	playerDmg := p.EffectiveAttack() + p.WeaponPower() - e.Type.Defense
	if playerDmg < 1 {
		playerDmg = 1
	}
	hitsToKillPlayer := float64(p.HP) / math.Max(1, avgEnemyDmg)
	hitsToKillEnemy := float64(e.Type.MaxHP) / float64(playerDmg)
	dangerScore := hitsToKillEnemy / hitsToKillPlayer
	switch {
	case dangerScore < 0.15:
		return 1
	case dangerScore < 0.4:
		return 2
	case dangerScore < 0.8:
		return 3
	default:
		return 4
	}
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
