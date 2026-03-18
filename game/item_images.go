package game

import (
	"image"
	_ "image/png"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
)

// itemImageEntry maps an asset path to the item whose Image field should be set.
type itemImageEntry struct {
	path string
	item *Item
}

// itemImageTable lists every item sprite and its asset path.
var itemImageTable = []itemImageEntry{
	// potions
	{"assets/items/health_potion_large.png", ItemLargeHealthPotion},
	{"assets/items/health_potion_medium.png", ItemMediumHealthPotion},
	{"assets/items/health_potion_small.png", ItemSmallHealthPotion},
	// food
	{"assets/items/food/apple.png", ItemApple},
	{"assets/items/food/bread_roll.png", ItemBreadRoll},
	{"assets/items/food/egg_fried.png", ItemFriedEgg},
	{"assets/items/food/grapes.png", ItemGrapes},
	{"assets/items/food/meat.png", ItemMeat},
	{"assets/items/food/mushroom.png", ItemMushroom},
	{"assets/items/food/pizza_slice.png", ItemPizzaSlice},
	// backpacks
	{"assets/items/gear/backpack/basic.png", ItemSmallBackpack},
	{"assets/items/gear/backpack/medium.png", ItemMediumBackpack},
	{"assets/items/gear/backpack/large.png", ItemLargeBackpack},
	// legs
	{"assets/items/gear/legs/pants.png", ItemPants},
	// accessories
	{"assets/items/accessories/necklace_diamond.png", ItemNecklaceDiamond},
	{"assets/items/accessories/necklace_skull.png", ItemNecklaceSkull},
	{"assets/items/accessories/necklace_star.png", ItemNecklaceStar},
	{"assets/items/accessories/necklace_tooth.png", ItemNecklaceTooth},
	{"assets/items/accessories/ring_diamond.png", ItemDiamondRing},
	{"assets/items/accessories/ring_diamond_2.png", ItemDiamondRing2},
	{"assets/items/accessories/ring_gold.png", ItemGoldRing},
	{"assets/items/accessories/ring_silver.png", ItemSilverRing},
	// gloves
	{"assets/items/gear/gloves/gloves_finger.png", ItemGlovesFinger},
	{"assets/items/gear/gloves/gloves_leather_metal.png", ItemGlovesLeatherMetal},
	{"assets/items/gear/gloves/gloves_leather.png", ItemGlovesLeather},
	{"assets/items/gear/gloves/gloves_metal.png", ItemGlovesMetal},
	// helmets
	{"assets/items/gear/head/basic_helmet.png", ItemBasicHelmet},
	{"assets/items/gear/head/coif.png", ItemCoif},
	{"assets/items/gear/head/full_helmet.png", ItemFullHelmet},
	{"assets/items/gear/head/gold_helmet.png", ItemGoldHelmet},
	{"assets/items/gear/head/horn_helmet.png", ItemHornHelmet},
	// shoes
	{"assets/items/gear/shoes/shoes_gold.png", ItemGoldShoes},
	{"assets/items/gear/shoes/shoes_leather.png", ItemLeatherShoes},
	{"assets/items/gear/shoes/shoes_metal.png", ItemMetalShoes},
	{"assets/items/gear/shoes/shoes_simple.png", ItemSimpleShoes},
	// shields
	{"assets/items/shield/shield_bronze.png", ItemBronzeShield},
	{"assets/items/shield/shield_gold.png", ItemGoldShield},
	{"assets/items/shield/shield_metal.png", ItemMetalShield},
	{"assets/items/shield/shield_wood.png", ItemWoodenShield},
	// armor
	{"assets/items/armor/basic.png", ItemBasicArmor},
	{"assets/items/armor/bronze.png", ItemBronzeArmor},
	{"assets/items/armor/complex.png", ItemComplexArmor},
	{"assets/items/armor/gold.png", ItemGoldArmor},
}

// loadItemImages loads all item sprites from the embedded FS.
// Items whose file is missing or fails to decode keep their Color fallback.
func loadItemImages(assets fs.FS) {
	for _, e := range itemImageTable {
		f, err := assets.Open(e.path)
		if err != nil {
			continue
		}
		img, _, err := image.Decode(f)
		f.Close()
		if err != nil {
			continue
		}
		e.item.Image = ebiten.NewImageFromImage(img)
	}
}
