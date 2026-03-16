package game

import "math/rand"

// Potion is a health potion placed on the map.
type Potion struct {
	X, Y  int
	Item  *Item
	Taken bool
}

// newPotion creates a pickup at (x, y) with a randomly chosen item from the spawn pool.
func newPotion(x, y int, rng *rand.Rand) *Potion {
	item := SpawnableItems[rng.Intn(len(SpawnableItems))]
	return &Potion{X: x, Y: y, Item: item}
}
