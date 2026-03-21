package game

// Player holds the player's state and stats.
type Player struct {
	X, Y         int
	HP           int
	MaxHP        int
	Attack       int
	Defense      int
	Agility      int
	BaseAttack   int // base value unaffected by equipment
	BaseDefense  int
	BaseAgility  int
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
		Agility:      5,
		BaseAttack:   5,
		BaseDefense:  2,
		BaseAgility:  5,
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
	p.Agility += sign * mods.Agility
	p.MaxHP += sign * mods.HP
	if p.HP > p.MaxHP {
		p.HP = p.MaxHP
	}
	p.Inventory.MaxItems += sign * mods.InvSlots
	p.Inventory.MaxWeight += float64(sign) * mods.InvWeight
}

// WeaponPower returns the combined Power of all items in weapon slots.
func (p *Player) WeaponPower() int {
	power := 0
	for _, slot := range []EquipmentSlot{SlotRightWeapon, SlotLeftWeapon} {
		if item := p.Equipment.Slots[slot]; item != nil {
			power += item.Power
		}
	}
	return power
}

// WeaponSpeed returns the highest Speed among equipped weapons.
func (p *Player) WeaponSpeed() int {
	speed := 0
	for _, slot := range []EquipmentSlot{SlotRightWeapon, SlotLeftWeapon} {
		if item := p.Equipment.Slots[slot]; item != nil && item.Speed > speed {
			speed = item.Speed
		}
	}
	return speed
}

// IsEquipped reports whether item is currently equipped in any slot.
func (p *Player) IsEquipped(item *Item) bool {
	for _, equipped := range p.Equipment.Slots {
		if equipped == item {
			return true
		}
	}
	return false
}

// Equip marks the item at invIdx as equipped. The item stays in inventory.
// Among compatible slots an empty one is preferred; otherwise the first slot
// is used and its previous occupant is simply unmarked.
func (p *Player) Equip(invIdx int) bool {
	inv := p.Inventory
	if invIdx >= len(inv.Items) {
		return false
	}
	item := inv.Items[invIdx].Item
	if len(item.Slots) == 0 {
		return false
	}
	// Prefer an empty slot; fall back to the first candidate.
	target := item.Slots[0]
	for _, s := range item.Slots {
		if p.Equipment.Slots[s] == nil {
			target = s
			break
		}
	}
	if old := p.Equipment.Slots[target]; old != nil {
		p.applyStatMods(old.StatMods, -1)
	}
	p.Equipment.Slots[target] = item
	p.applyStatMods(item.StatMods, 1)
	return true
}

// Unequip removes the item from the given slot. The item stays in inventory.
func (p *Player) Unequip(slot EquipmentSlot) bool {
	item := p.Equipment.Slots[slot]
	if item == nil {
		return false
	}
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

// TakeDamage reduces HP using the shared damage formula.
func (p *Player) TakeDamage(attack int) {
	dmg := calcDamage(attack, p.Defense)
	p.HP -= dmg
	if p.HP < 0 {
		p.HP = 0
	}
}
