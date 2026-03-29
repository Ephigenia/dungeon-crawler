package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Update handles input, combat, and movement.
func (g *Game) Update() error {
	if g.combatFrames > 0 {
		g.combatFrames--
	}
	g.particles.Update()
	g.player.tickStaminaRegen()
	g.updateEnemies()
	for _, o := range g.objects {
		if o.State == ObjectStateOpening {
			o.openingTick--
			if o.openingTick <= 0 {
				o.State = ObjectStateOpened
				items := o.Type.Loot(g.rng)
				for _, item := range items {
					g.potions = append(g.potions, &Potion{X: g.player.X, Y: g.player.Y, Item: item})
				}
			}
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		g.inventoryOpen = !g.inventoryOpen
		g.inventoryCursor = 0
		g.equipmentCursor = 0
		g.inventoryFocus = true
		return nil
	}
	if g.inventoryOpen {
		g.updateInventory()
		return nil
	}

	suppressLeft := false
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		if items := g.potionsAt(g.player.X, g.player.Y); len(items) > 0 {
			suppressLeft = true
			for _, p := range items {
				if g.player.Inventory.Add(p.Item) {
					p.Taken = true
				}
			}
		}
	}

	up := ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW)
	down := ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS)
	left := ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || (!suppressLeft && ebiten.IsKeyPressed(ebiten.KeyA))
	right := ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD)

	if up {
		g.holdFramesUp++
	} else {
		g.holdFramesUp = 0
	}
	if down {
		g.holdFramesDown++
	} else {
		g.holdFramesDown = 0
	}
	if left {
		g.holdFramesLeft++
	} else {
		g.holdFramesLeft = 0
	}
	if right {
		g.holdFramesRight++
	} else {
		g.holdFramesRight = 0
	}

	dx, dy := 0, 0
	if g.shouldMove(g.holdFramesUp) {
		dy = -1
	}
	if g.shouldMove(g.holdFramesDown) {
		dy = 1
	}
	if g.shouldMove(g.holdFramesLeft) {
		dx = -1
	}
	if g.shouldMove(g.holdFramesRight) {
		dx = 1
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.Regenerate()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyO) {
		if o := g.closedObjectAdjacentTo(g.player.X, g.player.Y); o != nil {
			o.State = ObjectStateOpening
			if o.Type.SkipOpeningAnimation {
				o.openingTick = 1
			} else {
				o.openingTick = objectOpeningFrames
			}
			g.player.SpendStamina(staminaCostAction)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		if p := g.potionAt(g.player.X, g.player.Y); p != nil {
			if g.player.Inventory.Add(p.Item) {
				p.Taken = true
			}
		}
	}

	if (dx != 0 || dy != 0) && g.player.Stamina > 0 {
		nx, ny := g.player.X+dx, g.player.Y+dy
		if e := g.enemyAt(nx, ny); e != nil {
			g.resolveCombat(e)
			g.player.SpendStamina(staminaCostAction)
		} else if g.dungeon.IsWalkable(nx, ny) {
			if o := g.objectAt(nx, ny); o == nil || o.Type.PassableByPlayer {
				g.player.X, g.player.Y = nx, ny
				g.cameraX = float64(g.player.X * TileSize)
				g.cameraY = float64(g.player.Y * TileSize)
				g.player.SpendStamina(staminaCostMove)
			} else if o.Type.Destructable {
				dmg := calcPlayerDamage(
					g.player.EffectiveAttack(), g.player.WeaponPower(), g.player.WeaponSpeed(),
					g.player.EffectiveAgility(), g.player.Level, 0, g.rng,
				)
				if dmg < 1 {
					dmg = 1
				}
				o.HP -= dmg
				if o.HP <= 0 {
					o.Destroyed = true
					if o.Type.Loot != nil {
						for _, item := range o.Type.Loot(g.rng) {
							g.potions = append(g.potions, &Potion{X: o.X, Y: o.Y, Item: item})
						}
					}
				}
				g.player.SpendStamina(staminaCostAction)
			}
		}
	}
	return nil
}

// resolveCombat handles a bump attack between the player and an enemy.
func (g *Game) resolveCombat(e *Enemy) {
	hpBefore := e.HP
	dmg := calcPlayerDamage(
		g.player.EffectiveAttack(), g.player.WeaponPower(), g.player.WeaponSpeed(),
		g.player.EffectiveAgility(), g.player.Level, e.Type.Defense, g.rng,
	)
	isCrit := g.rng.Float64()*100 < g.player.EffectiveCritChance()
	if isCrit {
		dmg = int(float64(dmg) * 2.0)
	}
	e.HP -= dmg
	if e.HP < 0 {
		e.HP = 0
	}
	dmg = hpBefore - e.HP
	g.player.WearWeapons()
	g.player.AddEXP(5)

	hitLine := fmt.Sprintf("Hit %s for %d damage", e.Type.Name, dmg)
	if isCrit {
		hitLine = fmt.Sprintf("CRITICAL! Hit %s for %d damage", e.Type.Name, dmg)
	}

	// Splat on the enemy.
	offsetX := float64(ScreenW/2) - g.cameraX
	offsetY := float64(ScreenH/2) - g.cameraY
	ex := float32(float64(e.X*TileSize)+offsetX) + float32(TileSize)/2
	ey := float32(float64(e.Y*TileSize)+offsetY) + float32(TileSize)/2
	count := 8
	if isCrit {
		count = 18
	}
	g.particles.SpawnParticles(ex, ey, count, 60, 200, 100, g.rng)

	if e.IsAlive() {
		playerHpBefore := g.player.HP
		g.player.TakeDamage(e.Type.Attack, g.rng)
		playerDmg := playerHpBefore - g.player.HP

		// Blood splat on the player when hurt.
		if playerDmg > 0 {
			px := float32(float64(g.player.X*TileSize)+offsetX) + float32(TileSize)/2
			py := float32(float64(g.player.Y*TileSize)+offsetY) + float32(TileSize)/2
			g.particles.SpawnBlood(px, py, 8, g.rng)
		}
		g.combatLines = []string{
			hitLine,
			fmt.Sprintf("%s  %d / %d HP", e.Type.Name, e.HP, e.Type.MaxHP),
			fmt.Sprintf("%s hits you for %d damage", e.Type.Name, playerDmg),
		}
	} else {
		g.player.AddEXP(20)
		g.combatLines = []string{
			hitLine,
			fmt.Sprintf("%s defeated!", e.Type.Name),
		}
		dropChance := 10 + (g.player.Level - 1)
		if g.rng.Intn(100) < dropChance {
			g.potions = append(g.potions, newPotion(e.X, e.Y, g.rng))
		}
	}
	g.combatFrames = 120
}

// updateInventory routes input to the focused inventory panel.
func (g *Game) updateInventory() {
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		g.inventoryFocus = !g.inventoryFocus
		return
	}
	if g.inventoryFocus {
		g.updateInventoryItems()
	} else {
		g.updateInventoryEquipment()
	}
}

func (g *Game) updateInventoryItems() {
	inv := g.player.Inventory
	maxSlots := inv.MaxItems

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		if g.inventoryCursor < maxSlots-1 {
			g.inventoryCursor++
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		if g.inventoryCursor > 0 {
			g.inventoryCursor--
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if g.inventoryCursor+inventoryCols < maxSlots {
			g.inventoryCursor += inventoryCols
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		if g.inventoryCursor >= inventoryCols {
			g.inventoryCursor -= inventoryCols
		}
	}

	usePressed := inpututil.IsKeyJustPressed(ebiten.KeyU) || inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	if usePressed && g.inventoryCursor < len(inv.Items) {
		inst := inv.Items[g.inventoryCursor].Instance
		switch inst.Type.Category {
		case CategoryConsumable:
			if inst.Type.OnUse != nil && inst.Type.OnUse(g.player) {
				inv.Consume(g.inventoryCursor)
				if g.inventoryCursor >= len(inv.Items) && g.inventoryCursor > 0 {
					g.inventoryCursor--
				}
			}
		case CategoryEquipment:
			if g.player.IsEquipped(inst) {
				for _, slot := range EquipmentSlotOrder {
					if g.player.Equipment.Slots[slot] == inst {
						g.player.Unequip(slot)
						break
					}
				}
			} else {
				g.player.Equip(g.inventoryCursor)
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyX) && g.inventoryCursor < len(inv.Items) {
		inst := inv.Items[g.inventoryCursor].Instance
		if !g.player.IsEquipped(inst) {
			inv.Consume(g.inventoryCursor)
			if g.inventoryCursor >= len(inv.Items) && g.inventoryCursor > 0 {
				g.inventoryCursor--
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) && g.inventoryCursor < len(inv.Items) {
		inst := inv.Items[g.inventoryCursor].Instance
		if !g.player.IsEquipped(inst) {
			g.potions = append(g.potions, &Potion{X: g.player.X, Y: g.player.Y, Item: inst.Type})
			inv.Consume(g.inventoryCursor)
			if g.inventoryCursor >= len(inv.Items) && g.inventoryCursor > 0 {
				g.inventoryCursor--
			}
		}
	}
}

func (g *Game) updateInventoryEquipment() {
	numSlots := len(EquipmentSlotOrder)

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		if g.equipmentCursor < numSlots-1 {
			g.equipmentCursor++
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		if g.equipmentCursor > 0 {
			g.equipmentCursor--
		}
	}

	usePressed := inpututil.IsKeyJustPressed(ebiten.KeyU) || inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	if usePressed {
		slot := EquipmentSlotOrder[g.equipmentCursor]
		g.player.Unequip(slot)
	}
}

// shouldMove returns true on the first frame a key is held, then after
// repeatDelayFrames, every repeatIntervalFrames.
func (g *Game) shouldMove(holdFrames int) bool {
	if holdFrames <= 0 {
		return false
	}
	if holdFrames == 1 {
		return true
	}
	if holdFrames < repeatDelayFrames {
		return false
	}
	return (holdFrames-repeatDelayFrames)%repeatIntervalFrames == 0
}
