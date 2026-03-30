package game

import "github.com/hajimehoshi/ebiten/v2"

// Predefined enemy type definitions.
var (
	EnemyGhost = &EnemyType{
		Name:            "Ghost",
		MaxHP:           14,
		Attack:          5,
		Defense:         3, // ethereal — hard to damage
		MoveInterval:    ebiten.DefaultTPS / 7,
		VisionRange:     12, // senses presence from afar
		AnimSpeed:       15, // drifting — 4 fps
		IdleImagePath:   "assets/enemies/ghost_idle_strip.png",
		MoveImagePath:   "assets/enemies/ghost_move_strip.png",
		AttackImagePath: "assets/enemies/ghost_attack_strip.png",
		DeathImagePath:  "assets/enemies/ghost_death_strip.png",
	}
	EnemyRat = &EnemyType{
		Name:            "SkeleRatton",
		MaxHP:           5,
		Attack:          3,
		Defense:         0,
		MoveInterval:    ebiten.DefaultTPS / 8, // very fast scurry
		VisionRange:     7,
		AnimSpeed:       6, // quick — 10 fps
		IdleImagePath:   "assets/enemies/rat_idle_strip.png",
		MoveImagePath:   "assets/enemies/rat_move_strip.png",
		AttackImagePath: "assets/enemies/rat_attack_strip.png",
		DeathImagePath:  "assets/enemies/rat_death_strip.png",
	}
	EnemySkeleton = &EnemyType{
		Name:            "Skeleton",
		MaxHP:           15,
		Attack:          5,
		Defense:         2, // bony frame deflects blows
		MoveInterval:    ebiten.DefaultTPS / 4,
		VisionRange:     8,
		AnimSpeed:       8, // steady march — ~7 fps
		IdleImagePath:   "assets/enemies/skeleton_idle_up_strip.png",
		MoveImagePath:   "assets/enemies/skeleton_run_up_strip.png",
		AttackImagePath: "assets/enemies/skeleton_up_strip.png",
		DeathImagePath:  "assets/enemies/skeleton_death_strip.png",
	}
	EnemySpider = &EnemyType{
		Name:            "Spider",
		MaxHP:           18,
		Attack:          7,
		Defense:         1,
		MoveInterval:    ebiten.DefaultTPS / 6, // fast skitter
		VisionRange:     9,                     // eight eyes
		AnimSpeed:       6,                     // 10 fps
		IdleImagePath:   "assets/enemies/spider_idle_strip.png",
		MoveImagePath:   "assets/enemies/spider_move_strip.png",
		AttackImagePath: "assets/enemies/spider_attack_strip.png",
		DeathImagePath:  "assets/enemies/spider_death_strip.png",
	}
	EnemySlime = &EnemyType{
		Name:            "Slime",
		MaxHP:           20,
		Attack:          2,
		Defense:         1, // absorbs hits
		MoveInterval:    ebiten.DefaultTPS * 3, // very slow
		VisionRange:     4,                     // dim senses
		AnimSpeed:       12,                    // slow bounce — 5 fps
		IdleImagePath:   "assets/enemies/slime_idle_strip.png",
		MoveImagePath:   "assets/enemies/slime_move_strip.png",
		AttackImagePath: "assets/enemies/slime_attack_strip.png",
		DeathImagePath:  "assets/enemies/slime_death_strip.png",
	}
	EnemyBlueSlime = &EnemyType{
		Name:            "Blue Slime",
		MaxHP:           12,
		Attack:          2,
		Defense:         0,
		MoveInterval:    ebiten.DefaultTPS * 2,
		VisionRange:     5,
		AnimSpeed:       10, // 6 fps
		IdleImagePath:   "assets/enemies/slime_blue_idle_strip.png",
		MoveImagePath:   "assets/enemies/slime_blue_move_strip.png",
		AttackImagePath: "assets/enemies/slime_blue_attack_strip.png",
		DeathImagePath:  "assets/enemies/slime_blue_death_strip.png",
	}
)

// AllEnemyTypes is the pool used when spawning enemies in a new dungeon.
var AllEnemyTypes = []*EnemyType{
	EnemyBlueSlime,
	EnemyGhost,
	EnemyRat,
	EnemySkeleton,
	EnemySlime,
	EnemySpider,
}
