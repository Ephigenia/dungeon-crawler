package game

import "github.com/hajimehoshi/ebiten/v2"

// Predefined enemy type definitions.
var (
	EnemyGoblin = &EnemyType{
		Name:         "Goblin",
		MaxHP:        8,
		Attack:       3,
		Defense:      0,
		MoveInterval: ebiten.DefaultTPS / 6, // nimble
		VisionRange:  10,                    // alert and skittish
		ImagePath:    "assets/enemies/goblin.png",
	}
	EnemySkeleton = &EnemyType{
		Name:         "Skeleton",
		MaxHP:        10,
		Attack:       4,
		Defense:      1,
		MoveInterval: ebiten.DefaultTPS / 4, // steady
		VisionRange:  8,                     // average awareness
		ImagePath:    "assets/enemies/skeleton.png",
	}
	EnemyOrc = &EnemyType{
		Name:         "Orc",
		MaxHP:        15,
		Attack:       5,
		Defense:      2,
		MoveInterval: ebiten.DefaultTPS / 3, // heavy
		VisionRange:  6,                     // focused but not perceptive
		ImagePath:    "assets/enemies/orc.png",
	}
	EnemyTroll = &EnemyType{
		Name:         "Troll",
		MaxHP:        22,
		Attack:       7,
		Defense:      3,
		MoveInterval: ebiten.DefaultTPS / 2, // lumbering
		VisionRange:  4,                     // dim-witted, poor awareness
		ImagePath:    "assets/enemies/troll.png",
	}
)

// AllEnemyTypes is the pool used when spawning enemies in a new dungeon.
var AllEnemyTypes = []*EnemyType{
	EnemyGoblin,
	EnemySkeleton,
	EnemyOrc,
	EnemyTroll,
}
