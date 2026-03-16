package game

import (
	"fmt"
	"image/color"
)

// ItemCategory classifies an item's general purpose.
type ItemCategory string

const (
	CategoryConsumable ItemCategory = "consumable"
	CategoryEquipment  ItemCategory = "equipment"
	CategoryKeyItem    ItemCategory = "key-item"
	CategoryOther      ItemCategory = "other"
)

// Item describes a single type of item.
type Item struct {
	ID       string
	Weight   float64
	Category ItemCategory
	Color    color.RGBA
	Effect   string
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
	ItemBreadRoll    = newConsumable("bread_roll", 0.1, 2, color.RGBA{210, 175, 125, 255})
	ItemBreadLoaf    = newConsumable("bread_loaf", 0.4, 4, color.RGBA{190, 145, 90, 255})
	ItemFlatbread    = newConsumable("flatbread", 0.2, 3, color.RGBA{215, 190, 135, 255})
	ItemCrackers     = newConsumable("crackers", 0.1, 1, color.RGBA{200, 178, 130, 255})
	ItemSmokedSausa  = newConsumable("smoked_sausage", 0.3, 8, color.RGBA{148, 78, 48, 255})
	ItemDriedMeat    = newConsumable("dried_meat", 0.2, 4, color.RGBA{158, 100, 58, 255})
	ItemMeatPie      = newConsumable("meat_pie", 0.4, 8, color.RGBA{138, 88, 52, 255})
	ItemApple        = newConsumable("apple", 0.2, 3, color.RGBA{168, 88, 65, 255})
	ItemCarrot       = newConsumable("carrot", 0.1, 1, color.RGBA{198, 118, 55, 255})
	ItemMushroom     = newConsumable("mushroom", 0.1, 1, color.RGBA{172, 138, 98, 255})
	ItemHoneycomb    = newConsumable("honeycomb", 0.2, 3, color.RGBA{208, 158, 48, 255})

	// SpawnableItems is the combined pool used for random map pickup spawning.
	SpawnableItems = []*Item{
		ItemSmallHealthPotion, ItemMediumHealthPotion, ItemLargeHealthPotion,
		ItemBreadRoll, ItemBreadLoaf, ItemFlatbread, ItemCrackers,
		ItemSmokedSausa, ItemDriedMeat, ItemMeatPie,
		ItemApple, ItemCarrot, ItemMushroom, ItemHoneycomb,
	}
)
