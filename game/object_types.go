package game

import (
	"math/rand"
)

// Predefined object type definitions.
var (
	ObjectTypeWoodenChest = &ObjectType{
		Name:            "Wooden Chest",
		Openable:        true,
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
	ObjectTypeBarrel = &ObjectType{
		Name:                  "Barrel",
		MaxHP:                 20,
		DestroyedImagePath:    "assets/map/debris1.png",
		ImagePath:             "assets/map/barrel.png",
		WalkableWhenDestroyed: true,
	}
	ObjectTypeBookshelf = &ObjectType{
		Name:                  "Bookshelf",
		MaxHP:                 120,
		ImagePath:             "assets/map/bookshelf1.png",
		DestroyedImagePath:    "assets/map/debris1.png",
		Openable:              true,
		Destructable:          true,
		WalkableWhenDestroyed: true,
		Loot: func(_ *rand.Rand) []*Item {
			return []*Item{ItemSmallHealthPotion}
		},
	}
	ObjectTypeBookshelf2 = &ObjectType{
		Name:                  "Bookshelf",
		MaxHP:                 120,
		ImagePath:             "assets/map/bookshelf2.png",
		DestroyedImagePath:    "assets/map/debris1.png",
		Openable:              true,
		Destructable:          true,
		WalkableWhenDestroyed: true,
		Loot: func(_ *rand.Rand) []*Item {
			return []*Item{ItemSmallHealthPotion}
		},
	}
	ObjectTypeCloset = &ObjectType{
		Name:      "Closet",
		ImagePath: "assets/map/closet.png",
	}
	ObjectTypeCrate = &ObjectType{
		Name:                  "Crate",
		MaxHP:                 23,
		ImagePath:             "assets/map/crate.png",
		DestroyedImagePath:    "assets/map/debris1.png",
		WalkableWhenDestroyed: true,
		Loot: func(rng *rand.Rand) []*Item {
			if rng.Intn(10) != 0 {
				return nil
			}
			return []*Item{SpawnableItems[rng.Intn(len(SpawnableItems))]}
		},
	}
	ObjectTypeCrate2 = &ObjectType{
		Name:                  "Crate",
		MaxHP:                 10,
		ImagePath:             "assets/map/crate2.png",
		DestroyedImagePath:    "assets/map/debris1.png",
		WalkableWhenDestroyed: true,
		Loot: func(rng *rand.Rand) []*Item {
			if rng.Intn(10) != 0 {
				return nil
			}
			return []*Item{SpawnableItems[rng.Intn(len(SpawnableItems))]}
		},
	}
	ObjectTypeTable = &ObjectType{
		Name:                  "Table",
		MaxHP:                 90,
		ImagePath:             "assets/map/table.png",
		DestroyedImagePath:    "assets/map/debris1.png",
		WalkableWhenDestroyed: true,
	}
	ObjectTypeTable2 = &ObjectType{
		Name:                  "Table",
		MaxHP:                 90,
		ImagePath:             "assets/map/table2.png",
		DestroyedImagePath:    "assets/map/debris1.png",
		WalkableWhenDestroyed: true,
	}
	ObjectTypeBones = &ObjectType{
		Name:             "Bones",
		PassableByPlayer: true,
		PassableByEnemy:  true,
		ImagePath:        "assets/map/bones.png",
	}
	ObjectTypeBones2 = &ObjectType{
		Name:             "Bones",
		PassableByPlayer: true,
		PassableByEnemy:  true,
		ImagePath:        "assets/map/bones2.png",
	}
)

// AllObjectTypes is the pool used when spawning objects in a new dungeon.
var AllObjectTypes = []*ObjectType{
	ObjectTypeBarrel,
	ObjectTypeBones,
	ObjectTypeBones2,
	ObjectTypeBookshelf,
	ObjectTypeBookshelf2,
	ObjectTypeCloset,
	ObjectTypeCrate,
	ObjectTypeCrate,
	ObjectTypeCrate2,
	ObjectTypeIronChest,
	ObjectTypeTable,
	ObjectTypeTable2,
	ObjectTypeWoodenChest,
}
