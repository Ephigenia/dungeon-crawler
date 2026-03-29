package game

// EquipmentSlot identifies a specific equipment position on the player.
type EquipmentSlot string

const (
	SlotHead        EquipmentSlot = "head"
	SlotBody        EquipmentSlot = "body"
	SlotLegs        EquipmentSlot = "legs"
	SlotFeet        EquipmentSlot = "feet"
	SlotNecklace    EquipmentSlot = "necklace"
	SlotLeftHand    EquipmentSlot = "left_hand"
	SlotRightHand   EquipmentSlot = "right_hand"
	SlotLeftRing    EquipmentSlot = "left_ring"
	SlotRightRing   EquipmentSlot = "right_ring"
	SlotLeftWeapon  EquipmentSlot = "left_weapon"
	SlotRightWeapon EquipmentSlot = "right_weapon"
	SlotBackpack    EquipmentSlot = "backpack"
)

// EquipmentSlotOrder defines the display order for the equipment panel.
var EquipmentSlotOrder = []EquipmentSlot{
	SlotHead, SlotBody, SlotLegs, SlotFeet, SlotNecklace,
	SlotLeftHand, SlotRightHand, SlotLeftRing, SlotRightRing,
	SlotLeftWeapon, SlotRightWeapon, SlotBackpack,
}

func slotLabel(s EquipmentSlot) string {
	switch s {
	case SlotHead:
		return "HEAD"
	case SlotBody:
		return "BODY"
	case SlotLegs:
		return "LEGS"
	case SlotFeet:
		return "FEET"
	case SlotNecklace:
		return "NECK"
	case SlotLeftHand:
		return "L.HAND"
	case SlotRightHand:
		return "R.HAND"
	case SlotLeftRing:
		return "L.RING"
	case SlotRightRing:
		return "R.RING"
	case SlotLeftWeapon:
		return "L.WPN"
	case SlotRightWeapon:
		return "R.WPN"
	case SlotBackpack:
		return "BACKPCK"
	}
	return string(s)
}

// StatModifiers holds stat changes applied when an item is equipped.
// Flat bonuses (HP, Attack, Defense, Agility) are added to the base stat.
// Percentage bonuses (*Pct) multiply the result: final = (base + flat) × (1 + pct/100).
// Multiple percentage bonuses from different items are summed before applying.
type StatModifiers struct {
	Agility    int
	AgilityPct float64 // % bonus to agility
	Attack     int
	AttackPct  float64 // % bonus to attack
	Defense    int
	DefensePct float64 // % bonus to defense
	HP         int
	HPPct      float64 // % bonus to max HP
	InvSlots   int     // extra inventory slots
	InvWeight  float64 // extra carry weight in kg
}

// Equipment tracks the item instance equipped in each slot.
type Equipment struct {
	Slots map[EquipmentSlot]*ItemInstance
}

func newEquipment() *Equipment {
	return &Equipment{Slots: make(map[EquipmentSlot]*ItemInstance)}
}
