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
	for _, c := range g.chests {
		if c.State == ChestStateOpening {
			c.openingTick--
			if c.openingTick <= 0 {
				c.State = ChestStateOpened
				count := g.rng.Intn(5) + 1
				for i := 0; i < count; i++ {
					item := SpawnableItems[g.rng.Intn(len(SpawnableItems))]
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

	up := ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW)
	down := ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS)
	left := ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA)
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
		if c := g.closedChestAdjacentTo(g.player.X, g.player.Y); c != nil {
			c.State = ChestStateOpening
			c.openingTick = chestOpeningFrames
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		if p := g.potionAt(g.player.X, g.player.Y); p != nil {
			if g.player.Inventory.Add(p.Item) {
				p.Taken = true
			}
		}
	}

	if dx != 0 || dy != 0 {
		nx, ny := g.player.X+dx, g.player.Y+dy
		if e := g.enemyAt(nx, ny); e != nil {
			g.resolveCombat(e)
		} else if g.dungeon.IsWalkable(nx, ny) {
			g.player.X, g.player.Y = nx, ny
			g.cameraX = float64(g.player.X * TileSize)
			g.cameraY = float64(g.player.Y * TileSize)
		}
	}
	return nil
}

// resolveCombat handles a bump attack between the player and an enemy.
func (g *Game) resolveCombat(e *Enemy) {
	hpBefore := e.HP
	dmg := calcPlayerDamage(
		g.player.Attack, g.player.WeaponPower(), g.player.WeaponSpeed(),
		g.player.Agility, g.player.Level, e.Defense,
	)
	e.HP -= dmg
	if e.HP < 0 {
		e.HP = 0
	}
	dmg = hpBefore - e.HP
	g.player.AddEXP(5)

	if e.IsAlive() {
		playerHpBefore := g.player.HP
		g.player.TakeDamage(e.Attack)
		playerDmg := playerHpBefore - g.player.HP
		g.combatLines = []string{
			fmt.Sprintf("Hit %s for %d damage", e.Name, dmg),
			fmt.Sprintf("%s  %d / %d HP", e.Name, e.HP, e.MaxHP),
			fmt.Sprintf("%s hits you for %d damage", e.Name, playerDmg),
		}
	} else {
		g.player.AddEXP(20)
		g.combatLines = []string{
			fmt.Sprintf("Hit %s for %d damage", e.Name, dmg),
			fmt.Sprintf("%s defeated!", e.Name),
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
		item := inv.Items[g.inventoryCursor].Item
		switch item.Category {
		case CategoryConsumable:
			if item.OnUse != nil && item.OnUse(g.player) {
				inv.Consume(g.inventoryCursor)
				if g.inventoryCursor >= len(inv.Items) && g.inventoryCursor > 0 {
					g.inventoryCursor--
				}
			}
		case CategoryEquipment:
			if g.player.IsEquipped(item) {
				for _, slot := range EquipmentSlotOrder {
					if g.player.Equipment.Slots[slot] == item {
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
		item := inv.Items[g.inventoryCursor].Item
		if !g.player.IsEquipped(item) {
			inv.Consume(g.inventoryCursor)
			if g.inventoryCursor >= len(inv.Items) && g.inventoryCursor > 0 {
				g.inventoryCursor--
			}
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) && g.inventoryCursor < len(inv.Items) {
		item := inv.Items[g.inventoryCursor].Item
		if !g.player.IsEquipped(item) {
			g.potions = append(g.potions, &Potion{X: g.player.X, Y: g.player.Y, Item: item})
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

