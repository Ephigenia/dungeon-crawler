package game

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	xpBase     = 100
	xpExponent = 1.5

	staminaRegenInterval = ebiten.DefaultTPS / 2 // frames between passive regen ticks (~2/sec)
	staminaCostMove      = 1
	staminaCostAction    = 2
)

// Player holds the player's state and stats.
// Base stats (BaseAttack, BaseDefense, BaseAgility, BaseMaxHP) grow on level-up
// and are never mutated by equipment. Effective stats are derived at call time
// via EffectiveXxx() methods, which add the sum of all equipped-item bonuses.
type Player struct {
	X, Y             int
	HP               int
	BaseMaxHP        int
	BaseAttack       int
	BaseDefense      int
	BaseAgility      int
	BaseMaxStamina   int
	Stamina          int
	Level            int
	EXP              int
	NextLevelEXP     int
	Inventory        *Inventory
	Equipment        *Equipment
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
		BaseMaxStamina: 40,
		Stamina:        40,
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
	for _, inst := range p.Equipment.Slots {
		if inst == nil {
			continue
		}
		pow := inst.EffectivePower()
		total.Attack += int(float64(inst.Type.StatMods.Attack) * pow)
		total.Defense += int(float64(inst.Type.StatMods.Defense) * pow)
		total.Agility += inst.Type.StatMods.Agility
		total.HP += inst.Type.StatMods.HP
		total.AttackPct += inst.Type.StatMods.AttackPct
		total.DefensePct += inst.Type.StatMods.DefensePct
		total.AgilityPct += inst.Type.StatMods.AgilityPct
		total.HPPct += inst.Type.StatMods.HPPct
		total.CritChance += inst.Type.StatMods.CritChance * pow
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

// EffectiveCritChance returns the player's critical hit probability in percent [0, 50].
// Base: agility×0.5 + attack×0.1. Weapon CritChance and StatMod bonuses are added on top.
func (p *Player) EffectiveCritChance() float64 {
	base := float64(p.EffectiveAgility())*0.5 + float64(p.EffectiveAttack())*0.1
	mods := p.equipmentStatMods()
	base += mods.CritChance
	for _, slot := range []EquipmentSlot{SlotRightWeapon, SlotLeftWeapon} {
		if inst := p.Equipment.Slots[slot]; inst != nil {
			base += inst.Type.CritChance * inst.EffectivePower()
		}
	}
	base *= 0.9
	if base > 50 {
		base = 50
	}
	return base
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

// staminaRegenRate returns the current frames-between-ticks, reduced by 1 per level (min 5).
func (p *Player) staminaRegenRate() int {
	interval := staminaRegenInterval - (p.Level - 1)
	if interval < 5 {
		interval = 5
	}
	return interval
}

// tickStaminaRegen advances the passive regen timer and restores 1 stamina when it fires.
func (p *Player) tickStaminaRegen() {
	p.staminaRegenTick++
	if p.staminaRegenTick >= p.staminaRegenRate() {
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

// WeaponPower returns the combined effective Power of all items in weapon slots,
// scaled by each weapon's durability.
func (p *Player) WeaponPower() int {
	power := 0
	for _, slot := range []EquipmentSlot{SlotRightWeapon, SlotLeftWeapon} {
		if inst := p.Equipment.Slots[slot]; inst != nil {
			power += int(float64(inst.Type.Power) * inst.EffectivePower())
		}
	}
	return power
}

// WeaponSpeed returns the highest Speed among equipped weapons.
func (p *Player) WeaponSpeed() int {
	speed := 0
	for _, slot := range []EquipmentSlot{SlotRightWeapon, SlotLeftWeapon} {
		if inst := p.Equipment.Slots[slot]; inst != nil && inst.Type.Speed > speed {
			speed = inst.Type.Speed
		}
	}
	return speed
}

// IsEquipped reports whether the given instance is currently equipped in any slot.
func (p *Player) IsEquipped(inst *ItemInstance) bool {
	for _, equipped := range p.Equipment.Slots {
		if equipped == inst {
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
	inst := inv.Items[invIdx].Instance
	if len(inst.Type.Slots) == 0 {
		return false
	}
	// Prefer an empty slot; fall back to the first candidate.
	target := inst.Type.Slots[0]
	for _, s := range inst.Type.Slots {
		if p.Equipment.Slots[s] == nil {
			target = s
			break
		}
	}
	if old := p.Equipment.Slots[target]; old != nil {
		p.applyStatMods(old.Type.StatMods, -1)
	}
	p.Equipment.Slots[target] = inst
	p.applyStatMods(inst.Type.StatMods, 1)
	return true
}

// Unequip removes the item from the given slot. The item stays in inventory.
func (p *Player) Unequip(slot EquipmentSlot) bool {
	inst := p.Equipment.Slots[slot]
	if inst == nil {
		return false
	}
	p.applyStatMods(inst.Type.StatMods, -1)
	p.Equipment.Slots[slot] = nil
	// Re-clamp HP now that the effective cap may have dropped.
	if p.HP > p.EffectiveMaxHP() {
		p.HP = p.EffectiveMaxHP()
	}
	return true
}

// durabilityLossFactor returns a [0.5, 1.0] multiplier for durability loss,
// reduced by 2% per level so higher-level players wear items down more slowly.
func (p *Player) durabilityLossFactor() float64 {
	factor := 1.0 - float64(p.Level-1)*0.02
	if factor < 0.5 {
		factor = 0.5
	}
	return factor
}

// reduceItemDurability subtracts scaled durability loss from inst, clamped to 0.
func (p *Player) reduceItemDurability(inst *ItemInstance) {
	if inst == nil || inst.Type.MaxDurability == 0 {
		return
	}
	loss := inst.Type.DurabilityLossRate * p.durabilityLossFactor() * 0.375
	inst.Durability -= loss
	if inst.Durability < 0 {
		inst.Durability = 0
	}
}

// WearWeapons reduces durability of equipped weapons (items with Power > 0 in weapon slots).
// Called after each player attack.
func (p *Player) WearWeapons() {
	for _, slot := range []EquipmentSlot{SlotRightWeapon, SlotLeftWeapon} {
		if inst := p.Equipment.Slots[slot]; inst != nil && inst.Type.Power > 0 {
			p.reduceItemDurability(inst)
		}
	}
}

// WearArmor reduces durability of equipped armor and shields.
// Called after the player takes damage.
func (p *Player) WearArmor() {
	for slot, inst := range p.Equipment.Slots {
		if inst == nil {
			continue
		}
		isWeaponSlot := slot == SlotLeftWeapon || slot == SlotRightWeapon
		isShield := isWeaponSlot && inst.Type.Power == 0
		if !isWeaponSlot || isShield {
			p.reduceItemDurability(inst)
		}
	}
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

// TakeDamage reduces HP using the shared damage formula and wears down armor.
// When the player is out of stamina they cannot defend and take full damage.
func (p *Player) TakeDamage(attack int, rng *rand.Rand) {
	defense := p.EffectiveDefense()
	if p.Stamina <= 0 {
		defense = 0
	}
	dmg := calcDamage(attack, defense, rng)
	p.HP -= dmg
	if p.HP < 0 {
		p.HP = 0
	}
	p.WearArmor()
}
