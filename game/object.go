package game

import "math/rand"

// ObjectKind identifies the type of map object.
type ObjectKind int

const (
	WoodenChest ObjectKind = iota // row 0 in chest spritesheet
	IronChest                     // row 1 in chest spritesheet
	Vase                          // standalone PNG, not openable
)

// ObjectState drives which animation frame is shown.
type ObjectState int

const (
	ObjectStateClosed  ObjectState = iota // frame 0
	ObjectStateOpening                    // frame 1 (transient)
	ObjectStateOpened                     // frame 2
)

// objectOpeningFrames is how many game ticks the opening animation plays.
const objectOpeningFrames = 20

// Object is an item placed in the dungeon map.
type Object struct {
	X, Y             int
	Kind             ObjectKind
	State            ObjectState
	Openable         bool // whether the player can open this object with O
	openingTick      int  // countdown while State == ObjectStateOpening
	PassableByPlayer bool // player can walk onto this tile
	PassableByEnemy  bool // enemies can walk onto this tile
}

// newObject creates an object of random kind at (x, y).
func newObject(x, y int, rng *rand.Rand) *Object {
	kind := ObjectKind(rng.Intn(3))
	openable := kind == WoodenChest || kind == IronChest
	return &Object{X: x, Y: y, Kind: kind, State: ObjectStateClosed, Openable: openable}
}

// spritesheetCol returns the column index (0–2) for the current state.
func (o *Object) spritesheetCol() int {
	switch o.State {
	case ObjectStateOpening:
		return 1
	case ObjectStateOpened:
		return 2
	default:
		return 0
	}
}

// isAdjacentTo reports whether the object is directly next to (x, y).
func (o *Object) isAdjacentTo(x, y int) bool {
	dx := o.X - x
	dy := o.Y - y
	return (dx == 0 && (dy == 1 || dy == -1)) || (dy == 0 && (dx == 1 || dx == -1))
}
