package game

import (
	"fmt"
	"image/color"
	"math"

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

	// Durability: MaxDurability 0 means indestructible (consumables, etc.)
	// DurabilityLossRate is subtracted from Durability each time the item is used.
	MaxDurability      int
	DurabilityLossRate float64
}

// durabilityExponent is the power curve exponent so that 50% durability → 75% power.
// Derived from: 0.5^a = 0.75 → a = log(0.75)/log(0.5).
var durabilityExponent = math.Log(0.75) / math.Log(0.5)

// ItemInstance is a live equipped or carried item with its own durability state.
// Multiple slots may not share the same instance; the inventory slot and equipment
// slot both point to the same *ItemInstance so durability changes are reflected everywhere.
type ItemInstance struct {
	Type      *Item
	Durability float64 // current durability; starts at Type.MaxDurability
}

// newItemInstance creates a fresh instance at full durability.
func newItemInstance(item *Item) *ItemInstance {
	return &ItemInstance{Type: item, Durability: float64(item.MaxDurability)}
}

// DurabilityPct returns durability as a fraction in [0, 1].
// Returns 1 for indestructible items (MaxDurability == 0).
func (inst *ItemInstance) DurabilityPct() float64 {
	if inst.Type.MaxDurability == 0 {
		return 1.0
	}
	pct := inst.Durability / float64(inst.Type.MaxDurability)
	if pct < 0 {
		return 0
	}
	if pct > 1 {
		return 1
	}
	return pct
}

// EffectivePower returns a [0,1] multiplier for attack/defense contributions.
// The curve is exponential: 100% → 1.0, 50% → 0.75, 0% → 0.
func (inst *ItemInstance) EffectivePower() float64 {
	pct := inst.DurabilityPct()
	if pct <= 0 {
		return 0
	}
	return math.Pow(pct, durabilityExponent)
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
			hpFull := p.HP >= p.EffectiveMaxHP()
			staminaFull := p.Stamina >= p.EffectiveMaxStamina()
			if hpFull && staminaFull {
				return false
			}
			p.HP += heal
			if p.HP > p.EffectiveMaxHP() {
				p.HP = p.EffectiveMaxHP()
			}
			p.RestoreStamina(heal)
			return true
		},
	}
}
