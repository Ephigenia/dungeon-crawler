package game

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
	Effect   string
	OnUse    func(p *Player) bool // returns true if the item is consumed on use
}

// Predefined item definitions.
var (
	ItemHealthPotion = &Item{
		ID:       "health_potion",
		Weight:   0.5,
		Category: CategoryConsumable,
		Effect:   "Restores 10 HP",
		OnUse: func(p *Player) bool {
			if p.HP >= p.MaxHP {
				return false
			}
			p.HP += 10
			if p.HP > p.MaxHP {
				p.HP = p.MaxHP
			}
			return true
		},
	}
)
