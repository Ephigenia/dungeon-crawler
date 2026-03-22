package game

// Predefined enemy type definitions.
var (
	EnemyGoblin = &EnemyType{
		Name:      "Goblin",
		MaxHP:     8,
		Attack:    3,
		Defense:   0,
		ImagePath: "assets/enemies/goblin.png",
	}
	EnemySkeleton = &EnemyType{
		Name:      "Skeleton",
		MaxHP:     10,
		Attack:    4,
		Defense:   1,
		ImagePath: "assets/enemies/skeleton.png",
	}
	EnemyOrc = &EnemyType{
		Name:      "Orc",
		MaxHP:     15,
		Attack:    5,
		Defense:   2,
		ImagePath: "assets/enemies/orc.png",
	}
	EnemyTroll = &EnemyType{
		Name:      "Troll",
		MaxHP:     22,
		Attack:    7,
		Defense:   3,
		ImagePath: "assets/enemies/troll.png",
	}
)

// AllEnemyTypes is the pool used when spawning enemies in a new dungeon.
var AllEnemyTypes = []*EnemyType{
	EnemyGoblin,
	EnemySkeleton,
	EnemyOrc,
	EnemyTroll,
}
