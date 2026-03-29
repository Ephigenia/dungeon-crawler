package game

import (
	"image/color"
	"math/rand"
)

// Predefined object type definitions.
var (
	// Ideas:
	// sarcophagus, barrel, crate, urn, cabinet, locker, chest of drawers, box, trunk, coffer
	ObjectTypeWoodenChest = &ObjectType{
		Name:            "Wooden Chest",
		Openable:        true,
		FallbackColor:   color.RGBA{180, 120, 60, 255},
		UsesSpritesheet: true,
		SpritesheetRow:  0,
		Loot: func(rng *rand.Rand) []*Item {
			count := rng.Intn(5) + 1
			items := make([]*Item, count)
			for i := range items {
				items[i] = SpawnableItems[rng.Intn(len(SpawnableItems))]
			}
			return items
		},
	}
	ObjectTypeIronChest = &ObjectType{
		Name:            "Iron Chest",
		Openable:        true,
		FallbackColor:   color.RGBA{160, 160, 170, 255},
		UsesSpritesheet: true,
		SpritesheetRow:  1,
		Loot: func(rng *rand.Rand) []*Item {
			count := rng.Intn(5) + 1
			items := make([]*Item, count)
			for i := range items {
				items[i] = SpawnableItems[rng.Intn(len(SpawnableItems))]
			}
			return items
		},
	}
	ObjectTypeVase = &ObjectType{
		Name:             "Vase",
		Openable:         false,
		FallbackColor:    color.RGBA{180, 160, 100, 255},
		SpritesheetPath:  "assets/map/map_objects.png",
		SpritesheetIndex: 0,
	}
	ObjectTypeCrate = &ObjectType{
		Name:                      "Crate",
		Openable:                  false,
		FallbackColor:             color.RGBA{140, 100, 60, 255},
		SpritesheetPath:           "assets/map/map_objects.png",
		SpritesheetIndex:          2,
		Destructable:              true,
		MaxHP:                     3,
		HasDestroyedSprite:        true,
		DestroyedSpritesheetIndex: 5,
		WalkableWhenDestroyed:     true,
		Loot: func(rng *rand.Rand) []*Item {
			if rng.Intn(10) != 0 {
				return nil
			}
			return []*Item{SpawnableItems[rng.Intn(len(SpawnableItems))]}
		},
	}
	// ObjectTypeShelf is placed against walls only; spawned separately from AllObjectTypes.
	ObjectTypeShelf = &ObjectType{
		Name:                 "Shelf",
		Openable:             true,
		FallbackColor:        color.RGBA{160, 130, 90, 255},
		SpritesheetPath:      "assets/map/map_objects.png",
		SpritesheetIndex:     1,
		SkipOpeningAnimation: true,
		Loot: func(_ *rand.Rand) []*Item {
			return []*Item{ItemSmallHealthPotion}
		},
	}
)

// AllObjectTypes is the pool used when spawning objects in a new dungeon.
var AllObjectTypes = []*ObjectType{
	ObjectTypeWoodenChest,
	ObjectTypeIronChest,
	ObjectTypeVase,
	ObjectTypeCrate,
}
