package game

import "math/rand"

const (
	potionMinHeal = 5
	potionMaxHeal = 15
)

// Potion is a health potion placed on the map.
type Potion struct {
	X, Y  int
	Heal  int
	Taken bool
}

// newPotion creates a potion at (x, y) with a random heal amount.
func newPotion(x, y int, rng *rand.Rand) *Potion {
	heal := potionMinHeal + rng.Intn(potionMaxHeal-potionMinHeal+1)
	return &Potion{X: x, Y: y, Heal: heal}
}
