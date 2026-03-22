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
)

// Item describes a single type of item.
type Item struct {
	ID        string
	Weight    float64
	Category  ItemCategory
	Slots     []EquipmentSlot // compatible slots for CategoryEquipment items
	Color     color.RGBA
	ImagePath string        // asset path for the sprite; empty = use Color fallback
	Image     *ebiten.Image // loaded at startup from ImagePath; nil until then
	Effect    string
	StatMods  StatModifiers
	MaxStack  int                  // max items per inventory slot (0 or 1 = not stackable)
	Power     int                  // weapon attack power, added to player.Attack only on hit
	Speed     int                  // weapon attack speed (higher = faster)
	OnUse     func(p *Player) bool // returns true if the item is consumed on use
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

// newConsumable builds a healing consumable item.
func newConsumable(id string, weight float64, heal int, col color.RGBA, imagePath string) *Item {
	effect := fmt.Sprintf("Restores %d HP", heal)
	return &Item{
		ID:        id,
		Weight:    weight,
		Category:  CategoryConsumable,
		Color:     col,
		ImagePath: imagePath,
		Effect:    effect,
		OnUse: func(p *Player) bool {
			if p.HP >= p.EffectiveMaxHP() {
				return false
			}
			p.HP += heal
			if p.HP > p.EffectiveMaxHP() {
				p.HP = p.EffectiveMaxHP()
			}
			return true
		},
	}
}
