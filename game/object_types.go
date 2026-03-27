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
	}
	ObjectTypeIronChest = &ObjectType{
		Name:            "Iron Chest",
		Openable:        true,
		FallbackColor:   color.RGBA{160, 160, 170, 255},
		UsesSpritesheet: true,
		SpritesheetRow:  1,
	}
	ObjectTypeVase = &ObjectType{
		Name:             "Vase",
		Openable:         false,
		FallbackColor:    color.RGBA{180, 160, 100, 255},
		SpritesheetPath:  "assets/map/map_objects.png",
		SpritesheetIndex: 0,
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
}
