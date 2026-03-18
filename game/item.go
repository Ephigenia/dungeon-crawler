package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// ItemCategory classifies an item's general purpose.
type ItemCategory string

const (
	CategoryConsumable ItemCategory = "consumable"
	CategoryEquipment  ItemCategory = "equipment"
	CategoryKeyItem    ItemCategory = "key-item"
	CategoryOther      ItemCategory = "other"
	CategoryBackpack   ItemCategory = "backpack"
)

// Item describes a single type of item.
type Item struct {
	ID       string
	Weight   float64
	Category ItemCategory
	Slot     EquipmentSlot // only set for CategoryEquipment items
	Color    color.RGBA
	Image    *ebiten.Image // optional sprite; falls back to Color if nil
	Effect   string
	StatMods StatModifiers
	OnUse    func(p *Player) bool // returns true if the item is consumed on use
}

func newConsumable(id string, weight float64, heal int, col color.RGBA) *Item {
	effect := fmt.Sprintf("Restores %d HP", heal)
	return &Item{
		ID:       id,
		Weight:   weight,
		Category: CategoryConsumable,
		Color:    col,
		Effect:   effect,
		OnUse: func(p *Player) bool {
			if p.HP >= p.MaxHP {
				return false
			}
			p.HP += heal
			if p.HP > p.MaxHP {
				p.HP = p.MaxHP
			}
			return true
		},
	}
}

// Predefined item definitions.
var (
	// Health potions
	ItemSmallHealthPotion  = newConsumable("small_health_potion", 0.3, 5, color.RGBA{210, 120, 120, 255})
	ItemMediumHealthPotion = newConsumable("medium_health_potion", 0.5, 10, color.RGBA{210, 80, 80, 255})
	ItemLargeHealthPotion  = newConsumable("large_health_potion", 0.8, 20, color.RGBA{180, 30, 30, 255})

	// Food items
	ItemBreadRoll = newConsumable("bread_roll", 0.1, 2, color.RGBA{210, 175, 125, 255})
	ItemGrapes    = newConsumable("grapes", 0.1, 3, color.RGBA{158, 100, 58, 255})
	ItemMeat      = newConsumable("meat", 0.2, 4, color.RGBA{158, 100, 58, 255})
	ItemApple     = newConsumable("apple", 0.2, 3, color.RGBA{168, 88, 65, 255})
	ItemMushroom  = newConsumable("mushroom", 0.1, 1, color.RGBA{172, 138, 98, 255})

	// Equipment items
	ItemLeatherHelmet = &Item{
		ID: "leather_helmet", Weight: 1.0, Category: CategoryEquipment, Slot: SlotHead,
		Color: color.RGBA{140, 120, 100, 255}, Effect: "+2 DEF",
		StatMods: StatModifiers{Defense: 2},
	}
	ItemLeatherArmor = &Item{
		ID: "leather_armor", Weight: 3.0, Category: CategoryEquipment, Slot: SlotBody,
		Color: color.RGBA{120, 100, 80, 255}, Effect: "+3 DEF",
		StatMods: StatModifiers{Defense: 3},
	}
	ItemLeatherBoots = &Item{
		ID: "leather_boots", Weight: 0.8, Category: CategoryEquipment, Slot: SlotFeet,
		Color: color.RGBA{130, 110, 90, 255}, Effect: "+1 DEF",
		StatMods: StatModifiers{Defense: 1},
	}
	ItemLeatherLegs = &Item{
		ID: "leather_legs", Weight: 1.2, Category: CategoryEquipment, Slot: SlotLegs,
		Color: color.RGBA{125, 105, 85, 255}, Effect: "+2 DEF",
		StatMods: StatModifiers{Defense: 2},
	}
	ItemIronSword = &Item{
		ID: "iron_sword", Weight: 2.5, Category: CategoryEquipment, Slot: SlotRightWeapon,
		Color: color.RGBA{180, 185, 198, 255}, Effect: "+5 ATK",
		StatMods: StatModifiers{Attack: 5},
	}
	ItemWoodenShield = &Item{
		ID: "wooden_shield", Weight: 2.0, Category: CategoryEquipment, Slot: SlotLeftHand,
		Color: color.RGBA{160, 120, 75, 255}, Effect: "+3 DEF",
		StatMods: StatModifiers{Defense: 3},
	}
	ItemGoldNecklace = &Item{
		ID: "gold_necklace", Weight: 0.2, Category: CategoryEquipment, Slot: SlotNecklace,
		Color: color.RGBA{218, 188, 48, 255}, Effect: "+10 HP",
		StatMods: StatModifiers{HP: 10},
	}
	ItemGoldRing = &Item{
		ID: "gold_ring", Weight: 0.1, Category: CategoryEquipment, Slot: SlotLeftRing,
		Color: color.RGBA{220, 195, 55, 255}, Effect: "+2 ATK",
		StatMods: StatModifiers{Attack: 2},
	}
	ItemSilverRing = &Item{
		ID: "silver_ring", Weight: 0.1, Category: CategoryEquipment, Slot: SlotRightRing,
		Color: color.RGBA{195, 198, 210, 255}, Effect: "+1 DEF",
		StatMods: StatModifiers{Defense: 1},
	}

	// Backpacks
	ItemSmallBackpack = &Item{
		ID: "small_backpack", Weight: 1.0, Category: CategoryBackpack, Slot: SlotBackpack,
		Color: color.RGBA{180, 138, 88, 255}, Effect: "+10 slots, +5 kg",
		StatMods: StatModifiers{InvSlots: 10, InvWeight: 5.0},
	}
	ItemMediumBackpack = &Item{
		ID: "medium_backpack", Weight: 1.5, Category: CategoryBackpack, Slot: SlotBackpack,
		Color: color.RGBA{158, 112, 68, 255}, Effect: "+15 slots, +7 kg",
		StatMods: StatModifiers{InvSlots: 15, InvWeight: 7.0},
	}
	ItemLargeBackpack = &Item{
		ID: "large_backpack", Weight: 2.0, Category: CategoryBackpack, Slot: SlotBackpack,
		Color: color.RGBA{138, 92, 50, 255}, Effect: "+15 slots, +20 kg",
		StatMods: StatModifiers{InvSlots: 15, InvWeight: 20.0},
	}

	// SpawnableItems is the combined pool used for random map pickup spawning.
	SpawnableItems = []*Item{
		ItemSmallHealthPotion, ItemMediumHealthPotion, ItemLargeHealthPotion,
		ItemApple, ItemMushroom, ItemGrapes, ItemMeat,
		ItemLeatherHelmet, ItemLeatherArmor, ItemLeatherBoots, ItemLeatherLegs,
		ItemIronSword, ItemWoodenShield, ItemGoldNecklace, ItemGoldRing, ItemSilverRing,
		ItemSmallBackpack, ItemMediumBackpack, ItemLargeBackpack,
	}
)
