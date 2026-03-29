package game

import "image/color"

func init() {
	ItemSmallHealthPotion.MaxStack = 5
	ItemMediumHealthPotion.MaxStack = 5
	ItemLargeHealthPotion.MaxStack = 5
}

// Predefined item definitions.
var (
	// Health potions
	ItemSmallHealthPotion  = newConsumable("small_health_potion", 0.3, 5, color.RGBA{210, 120, 120, 255}, "assets/items/health_potion_small.png")
	ItemMediumHealthPotion = newConsumable("medium_health_potion", 0.5, 10, color.RGBA{210, 80, 80, 255}, "assets/items/health_potion_medium.png")
	ItemLargeHealthPotion  = newConsumable("large_health_potion", 0.8, 20, color.RGBA{180, 30, 30, 255}, "assets/items/health_potion_large.png")

	// Buff potions
	ItemStrengthPotion = newTimedBuff("strength_potion", 0.3, 30*60, 50, color.RGBA{230, 210, 80, 255}, "assets/items/potion_small_yellow.png")

	// Food
	ItemBreadRoll  = newConsumable("bread_roll", 0.1, 2, color.RGBA{210, 175, 125, 255}, "assets/items/food/bread_roll.png")
	ItemGrapes     = newConsumable("grapes", 0.1, 3, color.RGBA{158, 100, 58, 255}, "assets/items/food/grapes.png")
	ItemFriedEgg   = newConsumable("fried_egg", 0.1, 1, color.RGBA{158, 100, 58, 255}, "assets/items/food/egg_fried.png")
	ItemMeat       = newConsumable("meat", 0.7, 5, color.RGBA{158, 100, 58, 255}, "assets/items/food/meat.png")
	ItemApple      = newConsumable("apple", 0.2, 3, color.RGBA{168, 88, 65, 255}, "assets/items/food/apple.png")
	ItemMushroom   = newConsumable("mushroom", 0.1, 1, color.RGBA{172, 138, 98, 255}, "assets/items/food/mushroom.png")
	ItemPizzaSlice = newConsumable("pizza_slice", 0.1, 1, color.RGBA{200, 100, 50, 255}, "assets/items/food/pizza_slice.png")

	// Weapons
	ItemIronSword = &Item{
		ID: "iron_sword", Weight: 2.5, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotRightWeapon},
		Color: color.RGBA{180, 185, 198, 255}, ImagePath: "assets/items/weapons/iron_sword.png",
		Power: 3, Speed: 5, CritChance: 2,
		MaxDurability: 80, DurabilityLossRate: 1.0,
	}
	Broadsword = &Item{
		ID: "broadsword", Weight: 3.0, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotRightWeapon},
		Color: color.RGBA{180, 185, 198, 255}, ImagePath: "assets/items/weapons/broadsword.png",
		Power: 5, Speed: 3, CritChance: 1,
		MaxDurability: 100, DurabilityLossRate: 0.8,
	}
	GoldenSword = &Item{
		ID: "golden_sword", Weight: 2.5, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotRightWeapon},
		Color: color.RGBA{180, 185, 198, 255}, ImagePath: "assets/items/weapons/golden_sword.png",
		Power: 7, Speed: 5, CritChance: 3,
		MaxDurability: 90, DurabilityLossRate: 0.9,
	}
	SwordJeweled = &Item{
		ID: "sword_jeweled", Weight: 2.5, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotRightWeapon},
		Color: color.RGBA{180, 185, 198, 255}, ImagePath: "assets/items/weapons/sword_jeweled.png",
		Effect: "DEF +1", Power: 8, Speed: 5, CritChance: 4, StatMods: StatModifiers{Defense: 1},
		MaxDurability: 100, DurabilityLossRate: 0.7,
	}
	MegaSword = &Item{
		ID: "mega_sword", Weight: 4.5, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotRightWeapon},
		Color: color.RGBA{180, 185, 198, 255}, ImagePath: "assets/items/weapons/mega_sword.png",
		Effect: "DEF +2", Power: 14, Speed: 2, CritChance: 0.5, StatMods: StatModifiers{Defense: 2},
		MaxDurability: 120, DurabilityLossRate: 0.6,
	}
	Saber = &Item{
		ID: "saber", Weight: 1.8, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotRightWeapon},
		Color: color.RGBA{180, 185, 198, 255}, ImagePath: "assets/items/weapons/saber.png",
		Power: 4, Speed: 7, CritChance: 5,
		MaxDurability: 70, DurabilityLossRate: 1.2,
	}
	Rapier1 = &Item{
		ID: "rapier", Weight: 1.5, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotRightWeapon},
		Color: color.RGBA{180, 185, 198, 255}, ImagePath: "assets/items/weapons/rapier1.png",
		Power: 3, Speed: 9, CritChance: 6,
		MaxDurability: 60, DurabilityLossRate: 1.5,
	}
	Rapier2 = &Item{
		ID: "rapier", Weight: 1.5, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotRightWeapon},
		Color: color.RGBA{180, 185, 198, 255}, ImagePath: "assets/items/weapons/rapier2.png",
		Power: 6, Speed: 9, CritChance: 8,
		MaxDurability: 60, DurabilityLossRate: 1.5,
	}
	Axe = &Item{
		ID: "axe", Weight: 1.0, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotRightWeapon},
		Color: color.RGBA{180, 185, 198, 255}, ImagePath: "assets/items/weapons/axe.png",
		Power: 5, Speed: 4, CritChance: 1.5,
		MaxDurability: 80, DurabilityLossRate: 1.0,
	}
	Hatchet = &Item{
		ID: "hatchet", Weight: 1.1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotRightWeapon},
		Color: color.RGBA{180, 185, 198, 255}, ImagePath: "assets/items/weapons/hatchet.png",
		Power: 4, Speed: 7, CritChance: 4,
		MaxDurability: 80, DurabilityLossRate: 1.0,
	}
	KnightsAxe = &Item{
		ID: "knights_axe", Weight: 1.6, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotRightWeapon},
		Color: color.RGBA{180, 185, 198, 255}, ImagePath: "assets/items/weapons/knights_axe.png",
		Power: 9, Speed: 3, CritChance: 2,
		MaxDurability: 100, DurabilityLossRate: 0.8,
	}
	ExecutionersAxe = &Item{
		ID: "executioners_axe", Weight: 2.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotRightWeapon},
		Color: color.RGBA{180, 185, 198, 255}, ImagePath: "assets/items/weapons/executioners_axe.png",
		Power: 12, Speed: 2, CritChance: 0.5,
		MaxDurability: 120, DurabilityLossRate: 0.5,
	}

	// Legs
	ItemPants = &Item{
		ID: "pants", Weight: 0.5, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLegs},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/legs/pants.png",
		MaxDurability: 60, DurabilityLossRate: 0.3,
	}

	// Armor
	ItemBasicArmor = &Item{
		ID: "basic_armor", Weight: 2.0, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotBody},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/armor/basic.png",
		Effect: "+1 DEF", StatMods: StatModifiers{Defense: 1},
		MaxDurability: 100, DurabilityLossRate: 0.5,
	}
	ItemComplexArmor = &Item{
		ID: "complex_armor", Weight: 2.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotBody},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/armor/complex.png",
		Effect: "+2 DEF", StatMods: StatModifiers{Defense: 2},
		MaxDurability: 120, DurabilityLossRate: 0.4,
	}
	ItemBronzeArmor = &Item{
		ID: "bronze_armor", Weight: 2.5, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotBody},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/armor/bronze.png",
		Effect: "+3 DEF", StatMods: StatModifiers{Defense: 3},
		MaxDurability: 150, DurabilityLossRate: 0.3,
	}
	ItemGoldArmor = &Item{
		ID: "gold_armor", Weight: 3, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotBody},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/armor/gold.png",
		Effect: "+5 DEF", StatMods: StatModifiers{Defense: 5},
		MaxDurability: 200, DurabilityLossRate: 0.2,
	}

	// Shields (equipped in weapon slots)
	ItemWoodenShield = &Item{
		ID: "wooden_shield", Weight: 1.0, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftWeapon, SlotRightWeapon},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/shield/shield_wood.png",
		Effect: "+2 DEF", StatMods: StatModifiers{Defense: 2},
		MaxDurability: 80, DurabilityLossRate: 1.0,
	}
	ItemMetalShield = &Item{
		ID: "metal_shield", Weight: 2.0, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftWeapon, SlotRightWeapon},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/shield/shield_metal.png",
		Effect: "+4 DEF", StatMods: StatModifiers{Defense: 4},
		MaxDurability: 120, DurabilityLossRate: 0.6,
	}
	ItemGoldShield = &Item{
		ID: "gold_shield", Weight: 2.5, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftWeapon, SlotRightWeapon},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/shield/shield_gold.png",
		Effect: "+4 DEF", StatMods: StatModifiers{Defense: 4},
		MaxDurability: 150, DurabilityLossRate: 0.4,
	}
	ItemBronzeShield = &Item{
		ID: "bronze_shield", Weight: 1.5, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftWeapon, SlotRightWeapon},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/shield/shield_bronze.png",
		Effect: "+2 DEF", StatMods: StatModifiers{Defense: 2},
		MaxDurability: 100, DurabilityLossRate: 0.8,
	}

	// Gloves (fit both hand slots)
	ItemGlovesFinger = &Item{
		ID: "gloves_finger", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftHand, SlotRightHand},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/gloves/gloves_finger.png",
		Effect: "+1 DEF",
		MaxDurability: 50, DurabilityLossRate: 0.3,
	}
	ItemGlovesLeather = &Item{
		ID: "gloves_leather", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftHand, SlotRightHand},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/gloves/gloves_leather.png",
		Effect: "+2 DEF", StatMods: StatModifiers{Defense: 1},
		MaxDurability: 60, DurabilityLossRate: 0.3,
	}
	ItemGlovesLeatherMetal = &Item{
		ID: "gloves_leather_metal", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftHand, SlotRightHand},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/gloves/gloves_leather_metal.png",
		Effect: "+2 DEF", StatMods: StatModifiers{Defense: 2},
		MaxDurability: 80, DurabilityLossRate: 0.25,
	}
	ItemGlovesMetal = &Item{
		ID: "gloves_metal", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftHand, SlotRightHand},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/gloves/gloves_metal.png",
		Effect: "+3 DEF", StatMods: StatModifiers{Defense: 3},
		MaxDurability: 100, DurabilityLossRate: 0.2,
	}

	// Helmets
	ItemCoif = &Item{
		ID: "coif", Weight: 0.1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotHead},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/head/coif.png",
		MaxDurability: 60, DurabilityLossRate: 0.3,
	}
	ItemBasicHelmet = &Item{
		ID: "basic_helmet", Weight: 1.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotHead},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/head/basic_helmet.png",
		Effect: "+1 DEF", StatMods: StatModifiers{Defense: 1},
		MaxDurability: 80, DurabilityLossRate: 0.4,
	}
	ItemFullHelmet = &Item{
		ID: "full_helmet", Weight: 1.4, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotHead},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/head/full_helmet.png",
		Effect: "+2 DEF", StatMods: StatModifiers{Defense: 2},
		MaxDurability: 100, DurabilityLossRate: 0.3,
	}
	ItemHornHelmet = &Item{
		ID: "horn_helmet", Weight: 1.3, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotHead},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/head/horn_helmet.png",
		Effect: "+2 DEF", StatMods: StatModifiers{Defense: 2},
		MaxDurability: 100, DurabilityLossRate: 0.3,
	}
	ItemGoldHelmet = &Item{
		ID: "gold_helmet", Weight: 2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotHead},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/head/gold_helmet.png",
		Effect: "+3 DEF", StatMods: StatModifiers{Defense: 3},
		MaxDurability: 150, DurabilityLossRate: 0.2,
	}

	// Shoes
	ItemSimpleShoes = &Item{
		ID: "simple_shoes", Weight: 1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotFeet},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/shoes/shoes_simple.png",
		MaxDurability: 60, DurabilityLossRate: 0.2,
	}
	ItemLeatherShoes = &Item{
		ID: "leather_shoes", Weight: 1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotFeet},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/shoes/shoes_leather.png",
		MaxDurability: 80, DurabilityLossRate: 0.2,
	}
	ItemMetalShoes = &Item{
		ID: "metal_shoes", Weight: 1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotFeet},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/shoes/shoes_metal.png",
		Effect: "+1 DEF", StatMods: StatModifiers{Defense: 1},
		MaxDurability: 100, DurabilityLossRate: 0.15,
	}
	ItemGoldShoes = &Item{
		ID: "gold_shoes", Weight: 1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotFeet},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/shoes/shoes_gold.png",
		Effect: "+2 DEF", StatMods: StatModifiers{Defense: 2},
		MaxDurability: 120, DurabilityLossRate: 0.1,
	}

	// Necklaces
	ItemNecklaceSkull = &Item{
		ID: "skull_necklace", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotNecklace},
		Color: color.RGBA{218, 188, 48, 255}, ImagePath: "assets/items/accessories/necklace_skull.png",
		Effect: "+20 HP", StatMods: StatModifiers{HP: 20},
		MaxDurability: 200, DurabilityLossRate: 0.1,
	}
	ItemNecklaceDiamond = &Item{
		ID: "diamond_necklace", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotNecklace},
		Color: color.RGBA{218, 188, 48, 255}, ImagePath: "assets/items/accessories/necklace_diamond.png",
		Effect: "+5 HP", StatMods: StatModifiers{HP: 5},
		MaxDurability: 200, DurabilityLossRate: 0.1,
	}
	ItemNecklaceStar = &Item{
		ID: "star_necklace", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotNecklace},
		Color: color.RGBA{218, 188, 48, 255}, ImagePath: "assets/items/accessories/necklace_star.png",
		Effect: "+5 HP", StatMods: StatModifiers{HP: 5},
		MaxDurability: 200, DurabilityLossRate: 0.1,
	}
	ItemNecklaceTooth = &Item{
		ID: "tooth_necklace", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotNecklace},
		Color: color.RGBA{218, 188, 48, 255}, ImagePath: "assets/items/accessories/necklace_tooth.png",
		Effect: "+5 HP, +3% CRIT", StatMods: StatModifiers{HP: 5, CritChance: 3},
		MaxDurability: 200, DurabilityLossRate: 0.1,
	}

	// Rings (fit either ring slot)
	ItemGoldRing = &Item{
		ID: "gold_ring", Weight: 0.1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftRing, SlotRightRing},
		Color: color.RGBA{220, 195, 55, 255}, ImagePath: "assets/items/accessories/ring_gold.png",
		Effect: "+2 ATK, +2% CRIT", StatMods: StatModifiers{Attack: 2, CritChance: 2},
		MaxDurability: 200, DurabilityLossRate: 0.05,
	}
	ItemSilverRing = &Item{
		ID: "silver_ring", Weight: 0.1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftRing, SlotRightRing},
		Color: color.RGBA{195, 198, 210, 255}, ImagePath: "assets/items/accessories/ring_silver.png",
		Effect: "+1 DEF", StatMods: StatModifiers{Defense: 1},
		MaxDurability: 200, DurabilityLossRate: 0.05,
	}
	ItemDiamondRing = &Item{
		ID: "diamond_ring", Weight: 0.1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftRing, SlotRightRing},
		Color: color.RGBA{220, 195, 55, 255}, ImagePath: "assets/items/accessories/ring_diamond.png",
		Effect: "+3 ATK, +3% CRIT", StatMods: StatModifiers{Attack: 3, CritChance: 3},
		MaxDurability: 200, DurabilityLossRate: 0.05,
	}
	ItemDiamondRing2 = &Item{
		ID: "diamond_ring2", Weight: 0.1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftRing, SlotRightRing},
		Color: color.RGBA{195, 198, 210, 255}, ImagePath: "assets/items/accessories/ring_diamond_2.png",
		Effect: "+2 DEF", StatMods: StatModifiers{Defense: 2},
		MaxDurability: 200, DurabilityLossRate: 0.05,
	}

	// Backpacks
	ItemSmallBackpack = &Item{
		ID: "small_backpack", Weight: 1.0, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotBackpack},
		Color: color.RGBA{180, 138, 88, 255}, ImagePath: "assets/items/gear/backpack/small.png",
		Effect: "+10 slots, +5 kg", StatMods: StatModifiers{InvSlots: 10, InvWeight: 5.0},
		MaxDurability: 150, DurabilityLossRate: 0.1,
	}
	ItemMediumBackpack = &Item{
		ID: "medium_backpack", Weight: 1.5, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotBackpack},
		Color: color.RGBA{158, 112, 68, 255}, ImagePath: "assets/items/gear/backpack/medium.png",
		Effect: "+15 slots, +7 kg", StatMods: StatModifiers{InvSlots: 15, InvWeight: 7.0},
		MaxDurability: 200, DurabilityLossRate: 0.08,
	}
	ItemLargeBackpack = &Item{
		ID: "large_backpack", Weight: 2.0, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotBackpack},
		Color: color.RGBA{138, 92, 50, 255}, ImagePath: "assets/items/gear/backpack/large.png",
		Effect: "+15 slots, +20 kg", StatMods: StatModifiers{InvSlots: 15, InvWeight: 20.0},
		MaxDurability: 250, DurabilityLossRate: 0.06,
	}

	// AllItems lists every item so loadItemImages can find them all.
	AllItems = []*Item{
		ItemSmallHealthPotion, ItemMediumHealthPotion, ItemLargeHealthPotion,
		ItemStrengthPotion,
		ItemBreadRoll, ItemGrapes, ItemFriedEgg, ItemMeat, ItemApple, ItemMushroom, ItemPizzaSlice,
		ItemIronSword, Broadsword, GoldenSword, SwordJeweled, MegaSword,
		Saber, Rapier1, Rapier2, Axe, Hatchet, KnightsAxe, ExecutionersAxe,
		ItemPants,
		ItemBasicArmor, ItemComplexArmor, ItemBronzeArmor, ItemGoldArmor,
		ItemWoodenShield, ItemMetalShield, ItemGoldShield, ItemBronzeShield,
		ItemGlovesFinger, ItemGlovesLeather, ItemGlovesLeatherMetal, ItemGlovesMetal,
		ItemCoif, ItemBasicHelmet, ItemFullHelmet, ItemHornHelmet, ItemGoldHelmet,
		ItemSimpleShoes, ItemLeatherShoes, ItemMetalShoes, ItemGoldShoes,
		ItemNecklaceSkull, ItemNecklaceDiamond, ItemNecklaceStar, ItemNecklaceTooth,
		ItemGoldRing, ItemSilverRing, ItemDiamondRing, ItemDiamondRing2,
		ItemSmallBackpack, ItemMediumBackpack, ItemLargeBackpack,
	}

	// SpawnableItems is the pool used for random map pickup spawning.
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
		ItemIronSword, Broadsword, GoldenSword, SwordJeweled, MegaSword,
		Saber, Rapier1, Rapier2, Axe, Hatchet, KnightsAxe, ExecutionersAxe,
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
		ItemStrengthPotion,
		ItemWoodenShield,
	}
)
