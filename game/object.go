package game

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// ObjectType is the static definition of a map object — its name, behaviour, and visuals.
// Instances on the map are represented by Object, which holds a pointer to its type
// plus mutable runtime state (position, open/close state).
type ObjectType struct {
	Name             string
	Openable         bool // player can open with O key
	PassableByPlayer bool // player can walk onto this tile
	PassableByEnemy  bool // enemies can walk onto this tile
	FallbackColor    color.RGBA // drawn when no image is available

	// Standalone image (loaded from ImagePath at startup).
	ImagePath string
	Image     *ebiten.Image

	// Spritesheet-based animation (shared animated_chests.png).
	// When UsesSpritesheet is true the object is rendered by slicing a 16×16
	// tile from the shared spritesheet at SpritesheetRow.
	UsesSpritesheet bool
	SpritesheetRow  int
}

// ObjectState drives which animation frame / spritesheet column is shown.
type ObjectState int

const (
	ObjectStateClosed  ObjectState = iota // column 0
	ObjectStateOpening                    // column 1 (transient)
	ObjectStateOpened                     // column 2
)

// objectOpeningFrames is how many game ticks the opening animation plays.
const objectOpeningFrames = 20

// Object is a live map object — a placed instance of an ObjectType.
type Object struct {
	X, Y        int
	Type        *ObjectType
	State       ObjectState
	openingTick int // countdown while State == ObjectStateOpening
}

// newObject places a randomly chosen ObjectType at (x, y).
func newObject(x, y int, rng *rand.Rand) *Object {
	t := AllObjectTypes[rng.Intn(len(AllObjectTypes))]
	return &Object{X: x, Y: y, Type: t, State: ObjectStateClosed}
}

// spritesheetCol returns the column index (0–2) in the chest spritesheet for the current state.
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
