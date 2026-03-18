package game

import "image/color"

// Predefined item definitions.
var (
	// Health potions
	ItemSmallHealthPotion  = newConsumable("small_health_potion", 0.3, 5, color.RGBA{210, 120, 120, 255}, "assets/items/health_potion_small.png")
	ItemMediumHealthPotion = newConsumable("medium_health_potion", 0.5, 10, color.RGBA{210, 80, 80, 255}, "assets/items/health_potion_medium.png")
	ItemLargeHealthPotion  = newConsumable("large_health_potion", 0.8, 20, color.RGBA{180, 30, 30, 255}, "assets/items/health_potion_large.png")

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
		ID: "iron_sword", Weight: 2.5, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotRightWeapon, SlotLeftWeapon},
		Color: color.RGBA{180, 185, 198, 255}, Effect: "+5 ATK",
		StatMods: StatModifiers{Attack: 5},
	}

	// Legs
	ItemPants = &Item{
		ID: "pants", Weight: 0.5, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLegs},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/legs/pants.png",
	}

	// Armor
	ItemBasicArmor = &Item{
		ID: "basic_armor", Weight: 2.0, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotBody},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/armor/basic.png",
		Effect: "+1 DEF", StatMods: StatModifiers{Defense: 1},
	}
	ItemComplexArmor = &Item{
		ID: "complex_armor", Weight: 2.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotBody},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/armor/complex.png",
		Effect: "+2 DEF", StatMods: StatModifiers{Defense: 2},
	}
	ItemBronzeArmor = &Item{
		ID: "bronze_armor", Weight: 2.5, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotBody},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/armor/bronze.png",
		Effect: "+3 DEF", StatMods: StatModifiers{Defense: 3},
	}
	ItemGoldArmor = &Item{
		ID: "gold_armor", Weight: 3, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotBody},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/armor/gold.png",
		Effect: "+5 DEF", StatMods: StatModifiers{Defense: 5},
	}

	// Shields (equipped in weapon slots)
	ItemWoodenShield = &Item{
		ID: "wooden_shield", Weight: 1.0, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftWeapon, SlotRightWeapon},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/shield/shield_wood.png",
		Effect: "+2 DEF", StatMods: StatModifiers{Defense: 2},
	}
	ItemMetalShield = &Item{
		ID: "metal_shield", Weight: 2.0, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftWeapon, SlotRightWeapon},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/shield/shield_metal.png",
		Effect: "+4 DEF", StatMods: StatModifiers{Defense: 4},
	}
	ItemGoldShield = &Item{
		ID: "gold_shield", Weight: 2.5, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftWeapon, SlotRightWeapon},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/shield/shield_gold.png",
		Effect: "+4 DEF", StatMods: StatModifiers{Defense: 4},
	}
	ItemBronzeShield = &Item{
		ID: "bronze_shield", Weight: 1.5, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftWeapon, SlotRightWeapon},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/shield/shield_bronze.png",
		Effect: "+2 DEF", StatMods: StatModifiers{Defense: 2},
	}

	// Gloves (fit both hand slots)
	ItemGlovesFinger = &Item{
		ID: "gloves_finger", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftHand, SlotRightHand},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/gloves/gloves_finger.png",
		Effect: "+1 DEF",
	}
	ItemGlovesLeather = &Item{
		ID: "gloves_leather", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftHand, SlotRightHand},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/gloves/gloves_leather.png",
		Effect: "+2 DEF", StatMods: StatModifiers{Defense: 1},
	}
	ItemGlovesLeatherMetal = &Item{
		ID: "gloves_leather_metal", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftHand, SlotRightHand},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/gloves/gloves_leather_metal.png",
		Effect: "+2 DEF", StatMods: StatModifiers{Defense: 2},
	}
	ItemGlovesMetal = &Item{
		ID: "gloves_metal", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftHand, SlotRightHand},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/gloves/gloves_metal.png",
		Effect: "+3 DEF", StatMods: StatModifiers{Defense: 3},
	}

	// Helmets
	ItemCoif = &Item{
		ID: "coif", Weight: 0.1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotHead},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/head/coif.png",
	}
	ItemBasicHelmet = &Item{
		ID: "basic_helmet", Weight: 1.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotHead},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/head/basic_helmet.png",
		Effect: "+1 DEF", StatMods: StatModifiers{Defense: 1},
	}
	ItemFullHelmet = &Item{
		ID: "full_helmet", Weight: 1.4, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotHead},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/head/full_helmet.png",
		Effect: "+2 DEF", StatMods: StatModifiers{Defense: 2},
	}
	ItemHornHelmet = &Item{
		ID: "horn_helmet", Weight: 1.3, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotHead},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/head/horn_helmet.png",
		Effect: "+2 DEF", StatMods: StatModifiers{Defense: 2},
	}
	ItemGoldHelmet = &Item{
		ID: "gold_helmet", Weight: 2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotHead},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/head/gold_helmet.png",
		Effect: "+3 DEF", StatMods: StatModifiers{Defense: 3},
	}

	// Shoes
	ItemSimpleShoes = &Item{
		ID: "simple_shoes", Weight: 1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotFeet},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/shoes/shoes_simple.png",
	}
	ItemLeatherShoes = &Item{
		ID: "leather_shoes", Weight: 1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotFeet},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/shoes/shoes_leather.png",
	}
	ItemMetalShoes = &Item{
		ID: "metal_shoes", Weight: 1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotFeet},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/shoes/shoes_metal.png",
		Effect: "+1 DEF", StatMods: StatModifiers{Defense: 1},
	}
	ItemGoldShoes = &Item{
		ID: "gold_shoes", Weight: 1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotFeet},
		Color: color.RGBA{160, 120, 75, 255}, ImagePath: "assets/items/gear/shoes/shoes_gold.png",
		Effect: "+2 DEF", StatMods: StatModifiers{Defense: 2},
	}

	// Necklaces
	ItemNecklaceSkull = &Item{
		ID: "skull_necklace", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotNecklace},
		Color: color.RGBA{218, 188, 48, 255}, ImagePath: "assets/items/accessories/necklace_skull.png",
		Effect: "+20 HP", StatMods: StatModifiers{HP: 20},
	}
	ItemNecklaceDiamond = &Item{
		ID: "diamond_necklace", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotNecklace},
		Color: color.RGBA{218, 188, 48, 255}, ImagePath: "assets/items/accessories/necklace_diamond.png",
		Effect: "+5 HP", StatMods: StatModifiers{HP: 5},
	}
	ItemNecklaceStar = &Item{
		ID: "star_necklace", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotNecklace},
		Color: color.RGBA{218, 188, 48, 255}, ImagePath: "assets/items/accessories/necklace_star.png",
		Effect: "+5 HP", StatMods: StatModifiers{HP: 5},
	}
	ItemNecklaceTooth = &Item{
		ID: "tooth_necklace", Weight: 0.2, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotNecklace},
		Color: color.RGBA{218, 188, 48, 255}, ImagePath: "assets/items/accessories/necklace_tooth.png",
		Effect: "+5 HP", StatMods: StatModifiers{HP: 5},
	}

	// Rings (fit either ring slot)
	ItemGoldRing = &Item{
		ID: "gold_ring", Weight: 0.1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftRing, SlotRightRing},
		Color: color.RGBA{220, 195, 55, 255}, ImagePath: "assets/items/accessories/ring_gold.png",
		Effect: "+2 ATK", StatMods: StatModifiers{Attack: 2},
	}
	ItemSilverRing = &Item{
		ID: "silver_ring", Weight: 0.1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftRing, SlotRightRing},
		Color: color.RGBA{195, 198, 210, 255}, ImagePath: "assets/items/accessories/ring_silver.png",
		Effect: "+1 DEF", StatMods: StatModifiers{Defense: 1},
	}
	ItemDiamondRing = &Item{
		ID: "diamond_ring", Weight: 0.1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftRing, SlotRightRing},
		Color: color.RGBA{220, 195, 55, 255}, ImagePath: "assets/items/accessories/ring_diamond.png",
		Effect: "+3 ATK", StatMods: StatModifiers{Attack: 3},
	}
	ItemDiamondRing2 = &Item{
		ID: "diamond_ring2", Weight: 0.1, Category: CategoryEquipment, Slots: []EquipmentSlot{SlotLeftRing, SlotRightRing},
		Color: color.RGBA{195, 198, 210, 255}, ImagePath: "assets/items/accessories/ring_diamond_2.png",
		Effect: "+2 DEF", StatMods: StatModifiers{Defense: 2},
	}

	// Backpacks
	ItemSmallBackpack = &Item{
		ID: "small_backpack", Weight: 1.0, Category: CategoryBackpack, Slots: []EquipmentSlot{SlotBackpack},
		Color: color.RGBA{180, 138, 88, 255}, ImagePath: "assets/items/gear/backpack/basic.png",
		Effect: "+10 slots, +5 kg", StatMods: StatModifiers{InvSlots: 10, InvWeight: 5.0},
	}
	ItemMediumBackpack = &Item{
		ID: "medium_backpack", Weight: 1.5, Category: CategoryBackpack, Slots: []EquipmentSlot{SlotBackpack},
		Color: color.RGBA{158, 112, 68, 255}, ImagePath: "assets/items/gear/backpack/medium.png",
		Effect: "+15 slots, +7 kg", StatMods: StatModifiers{InvSlots: 15, InvWeight: 7.0},
	}
	ItemLargeBackpack = &Item{
		ID: "large_backpack", Weight: 2.0, Category: CategoryBackpack, Slots: []EquipmentSlot{SlotBackpack},
		Color: color.RGBA{138, 92, 50, 255}, ImagePath: "assets/items/gear/backpack/large.png",
		Effect: "+15 slots, +20 kg", StatMods: StatModifiers{InvSlots: 15, InvWeight: 20.0},
	}

	// AllItems lists every item so loadItemImages can find them all.
	AllItems = []*Item{
		ItemSmallHealthPotion, ItemMediumHealthPotion, ItemLargeHealthPotion,
		ItemBreadRoll, ItemGrapes, ItemFriedEgg, ItemMeat, ItemApple, ItemMushroom, ItemPizzaSlice,
		ItemIronSword,
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
