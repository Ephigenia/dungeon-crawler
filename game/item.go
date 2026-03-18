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
	Slots    []EquipmentSlot // compatible slots for CategoryEquipment items
	Color    color.RGBA
	Image    *ebiten.Image // optional sprite; falls back to Color if nil
	Effect   string
	StatMods StatModifiers
	OnUse    func(p *Player) bool // returns true if the item is consumed on use
}

// FitsSlot reports whether the item can be placed in the given slot.
func (item *Item) FitsSlot(s EquipmentSlot) bool {
	for _, slot := range item.Slots {
		if slot == s {
			return true
		}
	}
	return false
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
	ItemBreadRoll  = newConsumable("bread_roll", 0.1, 2, color.RGBA{210, 175, 125, 255})
	ItemGrapes     = newConsumable("grapes", 0.1, 3, color.RGBA{158, 100, 58, 255})
	ItemFriedEgg   = newConsumable("fried_egg", 0.1, 1, color.RGBA{158, 100, 58, 255})
	ItemMeat       = newConsumable("meat", 0.7, 5, color.RGBA{158, 100, 58, 255})
	ItemApple      = newConsumable("apple", 0.2, 3, color.RGBA{168, 88, 65, 255})
	ItemMushroom   = newConsumable("mushroom", 0.1, 1, color.RGBA{172, 138, 98, 255})
	ItemPizzaSlice = newConsumable("pizza_slice", 0.1, 1, color.RGBA{200, 100, 50, 255})

	// weapons
	ItemIronSword = &Item{
		ID: "iron_sword", Weight: 2.5, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotRightWeapon, SlotLeftWeapon},
		Color: color.RGBA{180, 185, 198, 255}, Effect: "+5 ATK",
		StatMods: StatModifiers{Attack: 5},
	}

	// items // gear // legs
	ItemPants = &Item{
		ID: "pants", Weight: 0.5, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLegs},
		Color: color.RGBA{160, 120, 75, 255}, Effect: "",
		StatMods: StatModifiers{},
	}

	// gear // armor
	ItemBasicArmor = &Item{
		ID: "basic_armor", Weight: 2.0, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotBody},
		Color: color.RGBA{160, 120, 75, 255}, Effect: "+1 DEF",
		StatMods: StatModifiers{Defense: 1},
	}
	ItemComplexArmor = &Item{
		ID: "complex_armor", Weight: 2.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotBody},
		Color: color.RGBA{160, 120, 75, 255}, Effect: "+2 DEF",
		StatMods: StatModifiers{Defense: 2},
	}
	ItemBronzeArmor = &Item{
		ID: "bronze_armor", Weight: 2.5, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotBody},
		Color: color.RGBA{160, 120, 75, 255}, Effect: "+3 DEF",
		StatMods: StatModifiers{Defense: 3},
	}
	ItemGoldArmor = &Item{
		ID: "gold_armor", Weight: 3, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotBody},
		Color: color.RGBA{160, 120, 75, 255}, Effect: "+5 DEF",
		StatMods: StatModifiers{Defense: 5},
	}

	// gear // shields (equipped in weapon slots)
	ItemWoodenShield = &Item{
		ID: "wooden_shield", Weight: 1.0, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftWeapon, SlotRightWeapon},
		Color: color.RGBA{160, 120, 75, 255}, Effect: "+2 DEF",
		StatMods: StatModifiers{Defense: 2},
	}
	ItemMetalShield = &Item{
		ID: "metal_shield", Weight: 2.0, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftWeapon, SlotRightWeapon},
		Color: color.RGBA{160, 120, 75, 255}, Effect: "+4 DEF",
		StatMods: StatModifiers{Defense: 4},
	}
	ItemGoldShield = &Item{
		ID: "gold_shield", Weight: 2.5, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftWeapon, SlotRightWeapon},
		Color: color.RGBA{160, 120, 75, 255}, Effect: "+4 DEF",
		StatMods: StatModifiers{Defense: 4},
	}
	ItemBronzeShield = &Item{
		ID: "bronze_shield", Weight: 1.5, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftWeapon, SlotRightWeapon},
		Color: color.RGBA{160, 120, 75, 255}, Effect: "+2 DEF",
		StatMods: StatModifiers{Defense: 2},
	}

	// gear // gloves (fit both hand slots as a pair)
	ItemGlovesFinger = &Item{
		ID: "gloves_finger", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftHand, SlotRightHand},
		Color: color.RGBA{160, 120, 75, 255}, Effect: "+1 DEF",
		StatMods: StatModifiers{Defense: 0},
	}
	ItemGlovesLeather = &Item{
		ID: "gloves_leather", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftHand, SlotRightHand},
		Color: color.RGBA{160, 120, 75, 255}, Effect: "+2 DEF",
		StatMods: StatModifiers{Defense: 1},
	}
	ItemGlovesLeatherMetal = &Item{
		ID: "gloves_leather_metal", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftHand, SlotRightHand},
		Color: color.RGBA{160, 120, 75, 255}, Effect: "+2 DEF",
		StatMods: StatModifiers{Defense: 2},
	}
	ItemGlovesMetal = &Item{
		ID: "gloves_metal", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftHand, SlotRightHand},
		Color: color.RGBA{160, 120, 75, 255}, Effect: "+2 DEF",
		StatMods: StatModifiers{Defense: 3},
	}

	// items // gear // head
	ItemCoif = &Item{
		ID: "coif", Weight: 0.1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotHead},
		Color: color.RGBA{160, 120, 75, 255}, Effect: "",
		StatMods: StatModifiers{},
	}
	ItemBasicHelmet = &Item{
		ID: "basic_helmet", Weight: 1.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotHead},
		Color: color.RGBA{160, 120, 75, 255}, Effect: "+1 DEF",
		StatMods: StatModifiers{Defense: 1},
	}
	ItemFullHelmet = &Item{
		ID: "full_helmet", Weight: 1.4, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotHead},
		Color: color.RGBA{160, 120, 75, 255}, Effect: "+2 DEF",
		StatMods: StatModifiers{Defense: 2},
	}
	ItemHornHelmet = &Item{
		ID: "horn_helmet", Weight: 1.3, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotHead},
		Color: color.RGBA{160, 120, 75, 255}, Effect: "+2 DEF",
		StatMods: StatModifiers{Defense: 2},
	}
	ItemGoldHelmet = &Item{
		ID: "gold_helmet", Weight: 2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotHead},
		Color: color.RGBA{160, 120, 75, 255}, Effect: "+3 DEF",
		StatMods: StatModifiers{Defense: 3},
	}

	// items // gear // shoes
	ItemSimpleShoes = &Item{
		ID: "simple_shoes", Weight: 1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotFeet},
		Color: color.RGBA{160, 120, 75, 255}, Effect: "",
		StatMods: StatModifiers{},
	}
	ItemLeatherShoes = &Item{
		ID: "leather_shoes", Weight: 1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotFeet},
		Color: color.RGBA{160, 120, 75, 255}, Effect: "",
		StatMods: StatModifiers{},
	}
	ItemMetalShoes = &Item{
		ID: "metal_shoes", Weight: 1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotFeet},
		Color: color.RGBA{160, 120, 75, 255}, Effect: "+1 DEF",
		StatMods: StatModifiers{Defense: 1},
	}
	ItemGoldShoes = &Item{
		ID: "gold_shoes", Weight: 1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotFeet},
		Color: color.RGBA{160, 120, 75, 255}, Effect: "+2 DEF",
		StatMods: StatModifiers{Defense: 2},
	}

	// items // gear // necklaces and rings
	ItemNecklaceSkull = &Item{
		ID: "skull_necklace", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotNecklace},
		Color: color.RGBA{218, 188, 48, 255}, Effect: "+20 HP",
		StatMods: StatModifiers{HP: 20},
	}
	ItemNecklaceDiamond = &Item{
		ID: "diamond_necklace", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotNecklace},
		Color: color.RGBA{218, 188, 48, 255}, Effect: "+5 HP",
		StatMods: StatModifiers{HP: 5},
	}
	ItemNecklaceStar = &Item{
		ID: "star_necklace", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotNecklace},
		Color: color.RGBA{218, 188, 48, 255}, Effect: "+5 HP",
		StatMods: StatModifiers{HP: 5},
	}
	ItemNecklaceTooth = &Item{
		ID: "tooth_necklace", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotNecklace},
		Color: color.RGBA{218, 188, 48, 255}, Effect: "+5 HP",
		StatMods: StatModifiers{HP: 5},
	}

	// rings fit either ring slot
	ItemGoldRing = &Item{
		ID: "gold_ring", Weight: 0.1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftRing, SlotRightRing},
		Color: color.RGBA{220, 195, 55, 255}, Effect: "+2 ATK",
		StatMods: StatModifiers{Attack: 2},
	}
	ItemSilverRing = &Item{
		ID: "silver_ring", Weight: 0.1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftRing, SlotRightRing},
		Color: color.RGBA{195, 198, 210, 255}, Effect: "+1 DEF",
		StatMods: StatModifiers{Defense: 1},
	}
	ItemDiamondRing = &Item{
		ID: "diamond_ring", Weight: 0.1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftRing, SlotRightRing},
		Color: color.RGBA{220, 195, 55, 255}, Effect: "+2 ATK",
		StatMods: StatModifiers{Attack: 3},
	}
	ItemDiamondRing2 = &Item{
		ID: "diamond_ring2", Weight: 0.1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftRing, SlotRightRing},
		Color: color.RGBA{195, 198, 210, 255}, Effect: "+1 DEF",
		StatMods: StatModifiers{Defense: 2},
	}

	// Backpacks
	ItemSmallBackpack = &Item{
		ID: "small_backpack", Weight: 1.0, Category: CategoryBackpack, Slots: []EquipmentSlot{SlotBackpack},
		Color: color.RGBA{180, 138, 88, 255}, Effect: "+10 slots, +5 kg",
		StatMods: StatModifiers{InvSlots: 10, InvWeight: 5.0},
	}
	ItemMediumBackpack = &Item{
		ID: "medium_backpack", Weight: 1.5, Category: CategoryBackpack, Slots: []EquipmentSlot{SlotBackpack},
		Color: color.RGBA{158, 112, 68, 255}, Effect: "+15 slots, +7 kg",
		StatMods: StatModifiers{InvSlots: 15, InvWeight: 7.0},
	}
	ItemLargeBackpack = &Item{
		ID: "large_backpack", Weight: 2.0, Category: CategoryBackpack, Slots: []EquipmentSlot{SlotBackpack},
		Color: color.RGBA{138, 92, 50, 255}, Effect: "+15 slots, +20 kg",
		StatMods: StatModifiers{InvSlots: 15, InvWeight: 20.0},
	}

	// SpawnableItems is the combined pool used for random map pickup spawning.
	SpawnableItems = []*Item{
		ItemApple,
		ItemBasicArmor,
		ItemBasicHelmet,
		ItemBronzeArmor,
		ItemBronzeShield,
		ItemCoif,
		ItemComplexArmor,
		ItemDiamondRing,
		ItemDiamondRing2,
		ItemFullHelmet,
		ItemGlovesFinger,
		ItemGlovesLeather,
		ItemGlovesLeatherMetal,
		ItemGlovesMetal,
		ItemGoldArmor,
		ItemGoldHelmet,
		ItemGoldRing,
		ItemGoldShield,
		ItemGoldShoes,
		ItemGrapes,
		ItemHornHelmet,
		ItemIronSword,
		ItemLargeBackpack,
		ItemLargeHealthPotion,
		ItemLeatherShoes,
		ItemMeat,
		ItemMediumBackpack,
		ItemMediumHealthPotion,
		ItemMetalShield,
		ItemMetalShoes,
		ItemMushroom,
		ItemNecklaceDiamond,
		ItemNecklaceSkull,
		ItemNecklaceStar,
		ItemNecklaceTooth,
		ItemPants,
		ItemSilverRing,
		ItemSimpleShoes,
		ItemSmallBackpack,
		ItemSmallHealthPotion,
		ItemWoodenShield,
	}
)
