package game

import "math/rand"

// Enemy holds an enemy's state and stats.
type Enemy struct {
	X, Y    int
	HP      int
	MaxHP   int
	Attack  int
	Defense int
	Name    string
}

// IsAlive returns true if the enemy has HP remaining.
func (e *Enemy) IsAlive() bool {
	return e.HP > 0
}

// TakeDamage reduces HP using the shared damage formula.
func (e *Enemy) TakeDamage(attack int) {
	dmg := calcDamage(attack, e.Defense)
	e.HP -= dmg
	if e.HP < 0 {
		e.HP = 0
	}
}

// calcPlayerDamage computes damage for a player attack, incorporating weapon
// power, weapon speed, agility, and level.
//
//   weaponContrib  = weaponPower × (1 + weaponSpeed × agility / 100)
//   effectiveAttack = (baseAttack + weaponContrib) × (1 + (level−1) × 0.05)
//   damage          = max(0, int(effectiveAttack) − defense)
func calcPlayerDamage(baseAttack, weaponPower, weaponSpeed, agility, level, defense int) int {
	weaponContrib := float64(weaponPower) * (1.0 + float64(weaponSpeed)*float64(agility)/100.0)
	effective := (float64(baseAttack) + weaponContrib) * (1.0 + float64(level-1)*0.05)
	dmg := int(effective) - defense
	if dmg < 0 {
		dmg = 0
	}
	return dmg
}

// calcDamage returns the damage dealt given an attack and defense value.
// The attack/defense ratio is used as a factor, so higher attack relative
// to defense amplifies damage. Result is clamped to 0 (no minimum of 1).
func calcDamage(attack, defense int) int {
	if defense <= 0 {
		return attack
	}
	factor := float64(attack) / float64(defense)
	dmg := int(float64(attack-defense) * factor)
	if dmg < 0 {
		dmg = 0
	}
	return dmg
}

var enemyTypes = []struct {
	name    string
	hp      int
	attack  int
	defense int
}{
	{"Goblin", 8, 3, 0},
	{"Orc", 15, 5, 2},
	{"Skeleton", 10, 4, 1},
	{"Troll", 22, 7, 3},
}

func spawnEnemy(x, y int, rng *rand.Rand) *Enemy {
	t := enemyTypes[rng.Intn(len(enemyTypes))]
	return &Enemy{
		X:       x,
		Y:       y,
		HP:      t.hp,
		MaxHP:   t.hp,
		Attack:  t.attack,
		Defense: t.defense,
		Name:    t.name,
	}
}
