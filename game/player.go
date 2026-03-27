package game

import (
	"math"
	"math/rand"
)

const (
	xpBase     = 100
	xpExponent = 1.5

	staminaRegenInterval = 30 // frames between passive regen ticks (~2/sec at 60 fps)
	staminaCostMove      = 1
	staminaCostAction    = 2
)

// Player holds the player's state and stats.
// Base stats (BaseAttack, BaseDefense, BaseAgility, BaseMaxHP) grow on level-up
// and are never mutated by equipment. Effective stats are derived at call time
// via EffectiveXxx() methods, which add the sum of all equipped-item bonuses.
type Player struct {
	X, Y            int
	HP              int
	BaseMaxHP       int
	BaseAttack      int
	BaseDefense     int
	BaseAgility     int
	BaseMaxStamina  int
	Stamina         int
	Level           int
	EXP             int
	NextLevelEXP    int
	Inventory       *Inventory
	Equipment       *Equipment
	staminaRegenTick int
}

func newPlayer(x, y int) *Player {
	return &Player{
		X:              x,
		Y:              y,
		HP:             30,
		BaseMaxHP:      30,
		BaseAttack:     5,
		BaseDefense:    2,
		BaseAgility:    5,
		BaseMaxStamina: 20,
		Stamina:        20,
		Level:          1,
		EXP:            0,
		NextLevelEXP:   100,
		Inventory:      newInventory(),
		Equipment:      newEquipment(),
	}
}

// equipmentStatMods returns the cumulative stat modifiers of all equipped items.
func (p *Player) equipmentStatMods() StatModifiers {
	var total StatModifiers
	for _, item := range p.Equipment.Slots {
		if item != nil {
			total.Attack += item.StatMods.Attack
			total.Defense += item.StatMods.Defense
			total.Agility += item.StatMods.Agility
			total.HP += item.StatMods.HP
			total.AttackPct += item.StatMods.AttackPct
			total.DefensePct += item.StatMods.DefensePct
			total.AgilityPct += item.StatMods.AgilityPct
			total.HPPct += item.StatMods.HPPct
		}
	}
	return total
}

// applyPct applies a percentage modifier: result = int(base × (1 + pct/100)).
func applyPct(base int, pct float64) int {
	return int(float64(base) * (1.0 + pct/100.0))
}

// EffectiveMaxHP returns the HP cap including all equipment bonuses.
func (p *Player) EffectiveMaxHP() int {
	mods := p.equipmentStatMods()
	return applyPct(p.BaseMaxHP+mods.HP, mods.HPPct)
}

// EffectiveAttack returns the attack stat including all equipment bonuses.
func (p *Player) EffectiveAttack() int {
	mods := p.equipmentStatMods()
	return applyPct(p.BaseAttack+mods.Attack, mods.AttackPct)
}

// EffectiveDefense returns the defense stat including all equipment bonuses.
func (p *Player) EffectiveDefense() int {
	mods := p.equipmentStatMods()
	return applyPct(p.BaseDefense+mods.Defense, mods.DefensePct)
}

// EffectiveAgility returns the agility stat including all equipment bonuses.
func (p *Player) EffectiveAgility() int {
	mods := p.equipmentStatMods()
	return applyPct(p.BaseAgility+mods.Agility, mods.AgilityPct)
}

// EffectiveMaxStamina returns the stamina cap (no equipment modifiers yet).
func (p *Player) EffectiveMaxStamina() int {
	return p.BaseMaxStamina
}

// SpendStamina deducts amount from current stamina, clamped to 0.
func (p *Player) SpendStamina(amount int) {
	p.Stamina -= amount
	if p.Stamina < 0 {
		p.Stamina = 0
	}
}

// RestoreStamina adds amount to current stamina, clamped to the effective max.
func (p *Player) RestoreStamina(amount int) {
	p.Stamina += amount
	if p.Stamina > p.EffectiveMaxStamina() {
		p.Stamina = p.EffectiveMaxStamina()
	}
}

// tickStaminaRegen advances the passive regen timer and restores 1 stamina when it fires.
func (p *Player) tickStaminaRegen() {
	p.staminaRegenTick++
	if p.staminaRegenTick >= staminaRegenInterval {
		p.staminaRegenTick = 0
		p.RestoreStamina(1)
	}
}

// applyStatMods adjusts inventory capacity when equipping/unequipping items.
// Combat stats are derived dynamically; HP is clamped to the current effective cap.
func (p *Player) applyStatMods(mods StatModifiers, sign int) {
	p.Inventory.MaxItems += sign * mods.InvSlots
	p.Inventory.MaxWeight += float64(sign) * mods.InvWeight
	if p.HP > p.EffectiveMaxHP() {
		p.HP = p.EffectiveMaxHP()
	}
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
	// Re-clamp HP now that the effective cap may have dropped.
	if p.HP > p.EffectiveMaxHP() {
		p.HP = p.EffectiveMaxHP()
	}
	return true
}

// levelUp increases the player's level and improves all base stats.
//
// Stat growth per level:
//   - BaseMaxHP: +10%
//   - BaseAttack: +1 every level
//   - BaseDefense: +1 every 2 levels
//   - BaseAgility: +1 every 3 levels
//
// XP threshold follows an exponential curve: xpBase × level^xpExponent.
func (p *Player) levelUp() {
	p.Level++
	p.BaseMaxHP = p.BaseMaxHP * 110 / 100
	p.HP = p.EffectiveMaxHP()
	p.BaseAttack++
	if p.Level%2 == 0 {
		p.BaseDefense++
	}
	if p.Level%3 == 0 {
		p.BaseAgility++
	}
	p.BaseMaxStamina += 2
	p.RestoreStamina(p.EffectiveMaxStamina())
	p.NextLevelEXP = int(float64(xpBase) * math.Pow(float64(p.Level), xpExponent))
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
func (p *Player) TakeDamage(attack int, rng *rand.Rand) {
	dmg := calcDamage(attack, p.EffectiveDefense(), rng)
	p.HP -= dmg
	if p.HP < 0 {
		p.HP = 0
	}
}
