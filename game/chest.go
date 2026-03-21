package game

import "math/rand"

// ChestKind identifies the chest variant and maps to the spritesheet row.
type ChestKind int

const (
	ChestWooden ChestKind = iota // row 0 in spritesheet
	ChestIron                    // row 1 in spritesheet
)

// ChestState drives which animation frame is shown.
type ChestState int

const (
	ChestStateClosed  ChestState = iota // frame 0
	ChestStateOpening                   // frame 1 (transient)
	ChestStateOpened                    // frame 2
)

// chestOpeningFrames is how many game ticks the opening animation plays.
const chestOpeningFrames = 20

// Chest is an openable container placed in the dungeon.
type Chest struct {
	X, Y        int
	Kind        ChestKind
	State       ChestState
	openingTick int // countdown while State == ChestStateOpening
}

// newChest creates a chest of random kind at (x, y).
func newChest(x, y int, rng *rand.Rand) *Chest {
	kind := ChestKind(rng.Intn(2))
	return &Chest{X: x, Y: y, Kind: kind, State: ChestStateClosed}
}

// spritesheetCol returns the column index (0–2) for the current state.
func (c *Chest) spritesheetCol() int {
	switch c.State {
	case ChestStateOpening:
		return 1
	case ChestStateOpened:
		return 2
	default:
		return 0
	}
}

// isAdjacentTo reports whether the chest is directly next to (x, y).
func (c *Chest) isAdjacentTo(x, y int) bool {
	dx := c.X - x
	dy := c.Y - y
	return (dx == 0 && (dy == 1 || dy == -1)) || (dy == 0 && (dx == 1 || dx == -1))
}
