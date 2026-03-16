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

func newHealPotion(id string, weight float64, heal int, col color.RGBA) *Item {
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
	ItemSmallHealthPotion  = newHealPotion("small_health_potion", 0.3, 5, color.RGBA{210, 120, 120, 255})
	ItemMediumHealthPotion = newHealPotion("medium_health_potion", 0.5, 10, color.RGBA{210, 80, 80, 255})
	ItemLargeHealthPotion  = newHealPotion("large_health_potion", 0.8, 20, color.RGBA{180, 30, 30, 255})

	// HealthPotions is the pool used for random potion spawning.
	HealthPotions = []*Item{ItemSmallHealthPotion, ItemMediumHealthPotion, ItemLargeHealthPotion}
)
