package game

// Predefined enemy type definitions.
var (
	EnemyGoblin = &EnemyType{
		Name:         "Goblin",
		MaxHP:        8,
		Attack:       3,
		Defense:      0,
		MoveInterval: 10, // nimble, ~6 moves/sec
		VisionRange:  10, // alert and skittish
		ImagePath:    "assets/enemies/goblin.png",
	}
	EnemySkeleton = &EnemyType{
		Name:         "Skeleton",
		MaxHP:        10,
		Attack:       4,
		Defense:      1,
		MoveInterval: 15, // steady, ~4 moves/sec
		VisionRange:  8,  // average awareness
		ImagePath:    "assets/enemies/skeleton.png",
	}
	EnemyOrc = &EnemyType{
		Name:         "Orc",
		MaxHP:        15,
		Attack:       5,
		Defense:      2,
		MoveInterval: 20, // heavy, ~3 moves/sec
		VisionRange:  6,  // focused but not perceptive
		ImagePath:    "assets/enemies/orc.png",
	}
	EnemyTroll = &EnemyType{
		Name:         "Troll",
		MaxHP:        22,
		Attack:       7,
		Defense:      3,
		MoveInterval: 30, // lumbering, ~2 moves/sec
		VisionRange:  4,  // dim-witted, poor awareness
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
