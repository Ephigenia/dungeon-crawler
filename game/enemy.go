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

// TakeDamage reduces HP by the incoming attack minus defense, minimum 1.
func (e *Enemy) TakeDamage(attack int) {
	dmg := attack - e.Defense
	if dmg < 1 {
		dmg = 1
	}
	e.HP -= dmg
	if e.HP < 0 {
		e.HP = 0
	}
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
