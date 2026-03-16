package game

// Player holds the player's state and stats.
type Player struct {
	X, Y         int
	HP           int
	MaxHP        int
	Attack       int
	Defense      int
	Level        int
	EXP          int
	NextLevelEXP int
	Inventory    *Inventory
	Equipment    *Equipment
}

func newPlayer(x, y int) *Player {
	return &Player{
		X:            x,
		Y:            y,
		HP:           30,
		MaxHP:        30,
		Attack:       5,
		Defense:      2,
		Level:        1,
		EXP:          0,
		NextLevelEXP: 100,
		Inventory:    newInventory(),
		Equipment:    newEquipment(),
	}
}

// applyStatMods adds or subtracts stat modifiers (sign = +1 or -1).
func (p *Player) applyStatMods(mods StatModifiers, sign int) {
	p.Attack += sign * mods.Attack
	p.Defense += sign * mods.Defense
	p.MaxHP += sign * mods.HP
	if p.HP > p.MaxHP {
		p.HP = p.MaxHP
	}
	p.Inventory.MaxItems += sign * mods.InvSlots
	p.Inventory.MaxWeight += float64(sign) * mods.InvWeight
}

// Equip moves the item at invIdx from inventory into its designated equipment slot.
// If the slot is already occupied the old item is swapped back to inventory.
func (p *Player) Equip(invIdx int) bool {
	inv := p.Inventory
	if invIdx >= len(inv.Items) {
		return false
	}
	item := inv.Items[invIdx]
	if item.Slot == "" {
		return false
	}
	old := p.Equipment.Slots[item.Slot]
	if old != nil {
		// After removing the new item, check if old item fits weight-wise.
		weightAfter := inv.CurrentWeight() - item.Weight + old.Weight
		if weightAfter > inv.MaxWeight {
			return false
		}
	}
	inv.Remove(invIdx)
	if old != nil {
		inv.Items = append(inv.Items, old)
		p.applyStatMods(old.StatMods, -1)
	}
	p.Equipment.Slots[item.Slot] = item
	p.applyStatMods(item.StatMods, 1)
	return true
}

// Unequip moves the item in the given slot back to inventory.
func (p *Player) Unequip(slot EquipmentSlot) bool {
	item := p.Equipment.Slots[slot]
	if item == nil {
		return false
	}
	inv := p.Inventory
	// After removing capacity bonuses, verify the inventory stays valid.
	newMaxItems := inv.MaxItems - item.StatMods.InvSlots
	newMaxWeight := inv.MaxWeight - item.StatMods.InvWeight
	if len(inv.Items)+1 > newMaxItems || inv.CurrentWeight()+item.Weight > newMaxWeight {
		return false
	}
	inv.Add(item)
	p.applyStatMods(item.StatMods, -1)
	p.Equipment.Slots[slot] = nil
	return true
}

// levelUp increases the player's level and improves stats.
func (p *Player) levelUp() {
	p.Level++
	p.MaxHP = p.MaxHP * 110 / 100
	p.HP = p.MaxHP
	p.NextLevelEXP = p.NextLevelEXP * 125 / 100
	p.Inventory.levelUp()
}

// AddEXP adds exp points and calls levelUp each time the threshold is reached.
func (p *Player) AddEXP(amount int) {
	p.EXP += amount
	for p.EXP >= p.NextLevelEXP {
		p.EXP -= p.NextLevelEXP
		p.levelUp()
	}
}

// IsAlive returns true if the player has HP remaining.
func (p *Player) IsAlive() bool {
	return p.HP > 0
}

// TakeDamage reduces HP by the incoming attack minus defense, minimum 1.
func (p *Player) TakeDamage(attack int) {
	dmg := attack - p.Defense
	if dmg < 1 {
		dmg = 1
	}
	p.HP -= dmg
	if p.HP < 0 {
		p.HP = 0
	}
}
