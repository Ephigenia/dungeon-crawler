package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// drawInventory renders the full-screen inventory overlay.
func (g *Game) drawInventory(screen *ebiten.Image) {
	const (
		panelX     = 20
		panelY     = 20
		panelW     = 600
		panelH     = 440
		slotSize   = 16
		slotGap    = 4
		slotStride = slotSize + slotGap
		gridX      = panelX + 14
		gridY      = panelY + 54
		dividerX   = 310
		equipX     = dividerX + 12
		equipY     = panelY + 54
		equipRowH  = 18
		detailY    = panelY + 310
	)

	border := color.RGBA{100, 110, 140, 255}
	bg := color.RGBA{28, 32, 42, 255}
	dim := color.RGBA{160, 160, 160, 255}
	green := color.RGBA{152, 210, 152, 255}
	yellow := color.RGBA{220, 210, 100, 255}
	red := color.RGBA{200, 80, 80, 255}
	white := color.RGBA{255, 255, 255, 255}

	// Panel background + borders
	vector.DrawFilledRect(screen, 0, 0, ScreenW, ScreenH, color.RGBA{0, 0, 0, 180}, false)
	vector.DrawFilledRect(screen, panelX, panelY, panelW, panelH, bg, false)
	drawBorder(screen, panelX, panelY, panelW, panelH, border)
	vector.DrawFilledRect(screen, dividerX, panelY+40, 1, panelH-40, border, false)
	vector.DrawFilledRect(screen, panelX, detailY-8, panelW, 1, border, false)

	// Section headers
	inv := g.player.Inventory
	itemsFocusColor, equipFocusColor := dim, dim
	if g.inventoryFocus {
		itemsFocusColor = white
	} else {
		equipFocusColor = white
	}
	text.Draw(screen, fmt.Sprintf("ITEMS  %d/%d  %.1f/%.1fkg",
		len(inv.Items), inv.MaxItems, inv.CurrentWeight(), inv.MaxWeight),
		g.hudFont, gridX, panelY+18, itemsFocusColor)
	text.Draw(screen, "EQUIPMENT", g.hudFont, equipX, panelY+18, equipFocusColor)

	// Item grid
	maxSlots := inv.MaxItems
	rows := (maxSlots + inventoryCols - 1) / inventoryCols
	for row := 0; row < rows; row++ {
		for col := 0; col < inventoryCols; col++ {
			idx := row*inventoryCols + col
			if idx >= maxSlots {
				break
			}
			sx := float32(gridX + col*slotStride)
			sy := float32(gridY + row*slotStride)

			selected := g.inventoryFocus && idx == g.inventoryCursor
			bgCol := color.RGBA{45, 50, 62, 255}
			if selected {
				bgCol = color.RGBA{70, 90, 115, 255}
			}
			vector.DrawFilledRect(screen, sx, sy, slotSize, slotSize, bgCol, false)

			borderCol := color.RGBA{80, 88, 108, 255}
			if idx < len(inv.Items) {
				slot := inv.Items[idx]
				inst := slot.Instance
				drawItemSprite(screen, inst.Type, sx, sy, slotSize, 2)
				if slot.Count > 1 {
					countStr := fmt.Sprintf("%d", slot.Count)
					text.Draw(screen, countStr, g.hudFont, int(sx)+slotSize-len(countStr)*4, int(sy)+slotSize-1, color.RGBA{255, 255, 200, 255})
				}
				if g.player.IsEquipped(inst) {
					borderCol = color.RGBA{220, 200, 60, 255}
				}
				if inst.Type.MaxDurability > 0 {
					if inst.DurabilityPct() <= 0 {
						borderCol = color.RGBA{200, 60, 60, 255}
					}
					barH := float32(slotSize) * float32(inst.DurabilityPct())
					vector.DrawFilledRect(screen, sx+float32(slotSize)-2, sy+float32(slotSize)-barH, 2, barH, color.RGBA{80, 180, 220, 255}, false)
				}
			}
			if selected {
				borderCol = color.RGBA{180, 200, 255, 255}
			}
			drawBorder(screen, sx, sy, slotSize, slotSize, borderCol)
		}
	}

	// Character stats section
	statsY := gridY + rows*slotStride + 16
	vector.DrawFilledRect(screen, float32(gridX), float32(statsY-7), float32(dividerX-gridX-12), 1, border, false)
	text.Draw(screen, "CHARACTER", g.hudFont, gridX, statsY+2, dim)
	statsY += 16

	text.Draw(screen, "HP", g.hudFont, gridX, statsY, dim)
	text.Draw(screen, fmt.Sprintf("%d / %d", g.player.HP, g.player.EffectiveMaxHP()), g.hudFont, gridX+28, statsY, white)
	drawStatBar(screen, float32(gridX+90), float32(statsY-9), 80, g.player.HP, g.player.EffectiveMaxHP(),
		color.RGBA{50, 20, 20, 220}, color.RGBA{200, 60, 60, 255})
	statsY += 12

	text.Draw(screen, "ATK", g.hudFont, gridX, statsY, dim)
	text.Draw(screen, fmt.Sprintf("%d / %d", g.player.BaseAttack, g.player.EffectiveAttack()+g.player.WeaponPower()), g.hudFont, gridX+28, statsY, color.RGBA{224, 180, 100, 255})
	statsY += 12

	text.Draw(screen, "DEF", g.hudFont, gridX, statsY, dim)
	text.Draw(screen, fmt.Sprintf("%d / %d", g.player.BaseDefense, g.player.EffectiveDefense()), g.hudFont, gridX+28, statsY, color.RGBA{100, 160, 220, 255})
	statsY += 12

	text.Draw(screen, "AGI", g.hudFont, gridX, statsY, dim)
	text.Draw(screen, fmt.Sprintf("%d / %d", g.player.BaseAgility, g.player.EffectiveAgility()), g.hudFont, gridX+28, statsY, color.RGBA{152, 210, 152, 255})
	statsY += 12

	text.Draw(screen, "LVL", g.hudFont, gridX, statsY, dim)
	text.Draw(screen, fmt.Sprintf("%d", g.player.Level), g.hudFont, gridX+28, statsY, white)
	statsY += 12

	text.Draw(screen, "EXP", g.hudFont, gridX, statsY, dim)
	text.Draw(screen, fmt.Sprintf("%d / %d", g.player.EXP, g.player.NextLevelEXP), g.hudFont, gridX+28, statsY, white)
	drawStatBar(screen, float32(gridX+90), float32(statsY-9), 80, g.player.EXP, g.player.NextLevelEXP,
		color.RGBA{20, 30, 50, 220}, color.RGBA{100, 160, 240, 255})

	// Equipment slot list
	const (
		colSwatch = equipX + 58
		colName   = equipX + 66
		colWeight = equipX + 178
		colEffect = equipX + 222
	)
	for i, slot := range EquipmentSlotOrder {
		ey := equipY + i*equipRowH
		selected := !g.inventoryFocus && i == g.equipmentCursor
		if selected {
			vector.DrawFilledRect(screen, equipX-2, float32(ey-11), panelW-float32(equipX-panelX)-4, equipRowH, color.RGBA{60, 70, 90, 255}, false)
		}
		labelCol := dim
		if selected {
			labelCol = white
		}
		text.Draw(screen, slotLabel(slot), g.hudFont, equipX, ey, labelCol)

		equipped := g.player.Equipment.Slots[slot]
		if equipped != nil {
			drawItemSprite(screen, equipped.Type, float32(colSwatch), float32(ey-8), 6, 0)
			nameCol := equipped.Type.Color
			if equipped.Type.MaxDurability > 0 && equipped.DurabilityPct() <= 0 {
				nameCol = color.RGBA{200, 60, 60, 255}
			}
			text.Draw(screen, equipped.Type.ID, g.hudFont, colName, ey, nameCol)
			text.Draw(screen, fmt.Sprintf("%.1fkg", equipped.Type.Weight), g.hudFont, colWeight, ey, dim)
			if equipped.Type.Effect != "" {
				text.Draw(screen, equipped.Type.Effect, g.hudFont, colEffect, ey, green)
			}
		} else {
			text.Draw(screen, "(empty)", g.hudFont, colName, ey, color.RGBA{70, 70, 70, 255})
		}
	}

	// Detail panel
	g.drawInventoryDetail(screen, inv, gridX, detailY, white, dim, green, yellow, red)

	// Controls hint
	text.Draw(screen, "[Tab] Switch   [Arrows/WASD] Navigate   [U/Enter] Action   [D] Drop   [X] Destroy   [I] Close",
		g.hudFont, panelX+6, panelY+panelH-10, color.RGBA{100, 100, 100, 255})
}

// drawInventoryDetail renders the selected item's details in the bottom panel.
func (g *Game) drawInventoryDetail(screen *ebiten.Image, inv *Inventory, x, panelY int,
	white, dim, green, yellow, red color.RGBA) {

	var selectedInst *ItemInstance
	var selectedSlot *InventorySlot
	var fromEquipment bool
	if g.inventoryFocus {
		if g.inventoryCursor < len(inv.Items) {
			selectedSlot = inv.Items[g.inventoryCursor]
			selectedInst = selectedSlot.Instance
		}
	} else {
		slot := EquipmentSlotOrder[g.equipmentCursor]
		selectedInst = g.player.Equipment.Slots[slot]
		fromEquipment = true
	}
	var selectedItem *Item
	if selectedInst != nil {
		selectedItem = selectedInst.Type
	}

	dy := panelY + 10
	if selectedItem == nil {
		label := "(empty slot)"
		if fromEquipment {
			label = "(slot empty)"
		}
		text.Draw(screen, label, g.hudFont, x, dy, color.RGBA{70, 70, 70, 255})
		return
	}

	itemName := selectedItem.ID
	if selectedSlot != nil && selectedSlot.Count > 1 {
		itemName = fmt.Sprintf("%s  x%d", selectedItem.ID, selectedSlot.Count)
	}
	text.Draw(screen, itemName, g.hudFont, x, dy, white)
	dy += 16
	text.Draw(screen, fmt.Sprintf("Type: %s   Weight: %.1f kg", selectedItem.Category, selectedItem.Weight),
		g.hudFont, x, dy, dim)
	dy += 14

	if selectedItem.Power > 0 || selectedItem.Speed > 0 || selectedItem.CritChance > 0 {
		weaponStr := ""
		if selectedItem.Power > 0 {
			weaponStr += fmt.Sprintf("Power: +%d", selectedItem.Power)
		}
		if selectedItem.Speed > 0 {
			if weaponStr != "" {
				weaponStr += "   "
			}
			weaponStr += fmt.Sprintf("Speed: %d", selectedItem.Speed)
		}
		if selectedItem.CritChance > 0 {
			if weaponStr != "" {
				weaponStr += "   "
			}
			weaponStr += fmt.Sprintf("Crit: +%.1f%%", selectedItem.CritChance)
		}
		text.Draw(screen, weaponStr, g.hudFont, x, dy, green)
		dy += 14
	}
	if selectedItem.Effect != "" {
		text.Draw(screen, fmt.Sprintf("Effect: %s", selectedItem.Effect), g.hudFont, x, dy, green)
		dy += 14
	}

	if mods := selectedItem.StatMods; mods != (StatModifiers{}) {
		modsStr := ""
		if mods.Attack != 0 {
			modsStr += fmt.Sprintf("ATK %+d  ", mods.Attack)
		}
		if mods.AttackPct != 0 {
			modsStr += fmt.Sprintf("ATK %+.0f%%  ", mods.AttackPct)
		}
		if mods.Defense != 0 {
			modsStr += fmt.Sprintf("DEF %+d  ", mods.Defense)
		}
		if mods.DefensePct != 0 {
			modsStr += fmt.Sprintf("DEF %+.0f%%  ", mods.DefensePct)
		}
		if mods.Agility != 0 {
			modsStr += fmt.Sprintf("AGI %+d  ", mods.Agility)
		}
		if mods.AgilityPct != 0 {
			modsStr += fmt.Sprintf("AGI %+.0f%%  ", mods.AgilityPct)
		}
		if mods.HP != 0 {
			modsStr += fmt.Sprintf("HP %+d  ", mods.HP)
		}
		if mods.HPPct != 0 {
			modsStr += fmt.Sprintf("HP %+.0f%%", mods.HPPct)
		}
		text.Draw(screen, modsStr, g.hudFont, x, dy, green)
		dy += 14
	}

	if selectedItem.MaxDurability > 0 {
		durCol := color.RGBA{80, 180, 220, 255}
		if selectedInst.DurabilityPct() <= 0 {
			durCol = color.RGBA{200, 60, 60, 255}
		}
		text.Draw(screen, fmt.Sprintf("Durability: %d / %d", int(selectedInst.Durability), selectedItem.MaxDurability),
			g.hudFont, x, dy, durCol)
		drawStatBar(screen, float32(x+110), float32(dy-9), 80, int(selectedInst.Durability), selectedItem.MaxDurability,
			color.RGBA{20, 40, 50, 220}, color.RGBA{80, 180, 220, 255})
		dy += 14
	}

	dy += 4
	if fromEquipment {
		text.Draw(screen, "[U/Enter] Unequip", g.hudFont, x, dy, yellow)
		return
	}
	switch selectedItem.Category {
	case CategoryConsumable:
		text.Draw(screen, "[U/Enter] Use    [D] Drop    [X] Destroy", g.hudFont, x, dy, yellow)
	case CategoryEquipment:
		if g.player.IsEquipped(selectedInst) {
			text.Draw(screen, "[U/Enter] Unequip", g.hudFont, x, dy, yellow)
		} else {
			text.Draw(screen, "[U/Enter] Equip  [D] Drop   [X] Destroy", g.hudFont, x, dy, yellow)
		}
	default:
		text.Draw(screen, "[D] Drop    [X] Destroy", g.hudFont, x, dy, red)
	}
}
