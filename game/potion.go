package game

import "math/rand"

// Potion is a health potion placed on the map.
type Potion struct {
	X, Y  int
	Item  *Item
	Taken bool
}

// newPotion creates a potion at (x, y) with a randomly chosen potion type.
func newPotion(x, y int, rng *rand.Rand) *Potion {
	item := HealthPotions[rng.Intn(len(HealthPotions))]
	return &Potion{X: x, Y: y, Item: item}
}
