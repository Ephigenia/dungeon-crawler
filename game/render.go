package game

import (
	"fmt"
	"image"
	"image/color"

	"github.com/ephigenia/ebit-engine-game-1/dungeon"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Draw renders the dungeon, entities, HUD, and any open UI overlays.
func (g *Game) Draw(screen *ebiten.Image) {
	g.drawWorld(screen)
	g.drawHUD(screen)
	g.drawCombatNotification(screen)

	if p := g.potionAt(g.player.X, g.player.Y); p != nil {
		text.Draw(screen, fmt.Sprintf("%s  [P] Pick up", p.Item.ID), g.hudFont, 4, ScreenH-6, color.RGBA{220, 210, 100, 255})
	}
	if o := g.closedObjectAdjacentTo(g.player.X, g.player.Y); o != nil {
		text.Draw(screen, fmt.Sprintf("%s  [O] Open", o.Type.Name), g.hudFont, 4, ScreenH-18, color.RGBA{220, 210, 100, 255})
	}
	if g.inventoryOpen {
		g.drawInventory(screen)
	}
}

// drawWorld renders dungeon tiles, pickups, enemies, and the player.
func (g *Game) drawWorld(screen *ebiten.Image) {
	offsetX := float64(ScreenW/2) - g.cameraX
	offsetY := float64(ScreenH/2) - g.cameraY

	startTileX := int((-offsetX) / TileSize)
	startTileY := int((-offsetY) / TileSize)
	endTileX := startTileX + ScreenW/TileSize + 2
	endTileY := startTileY + ScreenH/TileSize + 2

	var op ebiten.DrawImageOptions
	for ty := startTileY; ty <= endTileY; ty++ {
		for tx := startTileX; tx <= endTileX; tx++ {
			t := g.dungeon.At(tx, ty)
			px := float64(tx*TileSize) + offsetX + 0.5
			py := float64(ty*TileSize) + offsetY + 0.5
			op.GeoM.Reset()
			op.GeoM.Translate(px, py)
			if t == dungeon.Wall {
				screen.DrawImage(g.wallTileImg, &op)
			} else {
				screen.DrawImage(g.floorTileImg, &op)
			}
		}
	}

	const pickupSize = 16
	for _, p := range g.potions {
		if p.Taken {
			continue
		}
		px := float32(float64(p.X*TileSize) + offsetX + float64(TileSize-pickupSize)/2)
		py := float32(float64(p.Y*TileSize) + offsetY + float64(TileSize-pickupSize)/2)
		drawItemSprite(screen, p.Item, px, py, pickupSize, 0)
	}

	for _, o := range g.objects {
		ox := float64(o.X*TileSize) + offsetX
		oy := float64(o.Y*TileSize) + offsetY
		switch {
		case o.Type.UsesSpritesheet && g.objectImg != nil:
			col := o.spritesheetCol()
			row := o.Type.SpritesheetRow
			src := g.objectImg.SubImage(image.Rect(col*16, row*16, col*16+16, row*16+16)).(*ebiten.Image)
			op.GeoM.Reset()
			op.GeoM.Translate(ox, oy)
			screen.DrawImage(src, &op)
		case o.Type.Image != nil:
			iw, ih := o.Type.Image.Bounds().Dx(), o.Type.Image.Bounds().Dy()
			op.GeoM.Reset()
			op.GeoM.Scale(float64(TileSize)/float64(iw), float64(TileSize)/float64(ih))
			op.GeoM.Translate(ox, oy)
			screen.DrawImage(o.Type.Image, &op)
		default:
			vector.DrawFilledRect(screen, float32(ox)+1, float32(oy)+1, TileSize-2, TileSize-2, o.Type.FallbackColor, false)
		}
	}

	for _, e := range g.enemies {
		if !e.IsAlive() {
			continue
		}
		ex := float64(e.X*TileSize) + offsetX + float64(TileSize-PlayerSize)/2
		ey := float64(e.Y*TileSize) + offsetY + float64(TileSize-PlayerSize)/2
		eImg := g.enemyImg
		if e.Type.Image != nil {
			eImg = e.Type.Image
		}
		iw, ih := eImg.Bounds().Dx(), eImg.Bounds().Dy()
		op.GeoM.Reset()
		op.GeoM.Scale(float64(PlayerSize)/float64(iw), float64(PlayerSize)/float64(ih))
		op.GeoM.Translate(ex, ey)
		screen.DrawImage(eImg, &op)
		drawStatBar(screen, float32(ex), float32(ey)+PlayerSize+1, PlayerSize, e.HP, e.Type.MaxHP,
			color.RGBA{30, 30, 30, 200}, color.RGBA{224, 108, 117, 255})
	}

	playerPx := float64(g.player.X*TileSize) + offsetX + float64(TileSize-PlayerSize)/2
	playerPy := float64(g.player.Y*TileSize) + offsetY + float64(TileSize-PlayerSize)/2
	op.GeoM.Reset()
	iw, ih := g.playerImg.Bounds().Dx(), g.playerImg.Bounds().Dy()
	op.GeoM.Scale(float64(PlayerSize)/float64(iw), float64(PlayerSize)/float64(ih))
	op.GeoM.Translate(playerPx, playerPy)
	screen.DrawImage(g.playerImg, &op)
}

// drawHUD renders the player stats overlay in the top-left corner.
func (g *Game) drawHUD(screen *ebiten.Image) {
	inv := g.player.Inventory
	dim := color.RGBA{160, 160, 160, 255}
	hudY := 14

	text.Draw(screen, "HP", g.hudFont, 4, hudY, dim)
	text.Draw(screen, fmt.Sprintf("%d / %d", g.player.HP, g.player.EffectiveMaxHP()), g.hudFont, 32, hudY, color.White)
	drawStatBar(screen, 94, float32(hudY-9), 80, g.player.HP, g.player.EffectiveMaxHP(),
		color.RGBA{50, 20, 20, 220}, color.RGBA{200, 60, 60, 255})
	hudY += 12

	text.Draw(screen, "ATK", g.hudFont, 4, hudY, dim)
	text.Draw(screen, fmt.Sprintf("%d / %d", g.player.BaseAttack, g.player.EffectiveAttack()+g.player.WeaponPower()), g.hudFont, 32, hudY, color.RGBA{224, 180, 100, 255})
	hudY += 12

	text.Draw(screen, "DEF", g.hudFont, 4, hudY, dim)
	text.Draw(screen, fmt.Sprintf("%d / %d", g.player.BaseDefense, g.player.EffectiveDefense()), g.hudFont, 32, hudY, color.RGBA{100, 160, 220, 255})
	hudY += 12

	text.Draw(screen, "AGI", g.hudFont, 4, hudY, dim)
	text.Draw(screen, fmt.Sprintf("%d / %d", g.player.BaseAgility, g.player.EffectiveAgility()), g.hudFont, 32, hudY, color.RGBA{152, 210, 152, 255})
	hudY += 12

	text.Draw(screen, "LVL", g.hudFont, 4, hudY, dim)
	text.Draw(screen, fmt.Sprintf("%d", g.player.Level), g.hudFont, 32, hudY, color.White)
	hudY += 12

	text.Draw(screen, "EXP", g.hudFont, 4, hudY, dim)
	text.Draw(screen, fmt.Sprintf("%d / %d", g.player.EXP, g.player.NextLevelEXP), g.hudFont, 32, hudY, color.White)
	drawStatBar(screen, 94, float32(hudY-9), 80, g.player.EXP, g.player.NextLevelEXP,
		color.RGBA{20, 30, 50, 220}, color.RGBA{100, 160, 240, 255})
	hudY += 12

	text.Draw(screen, fmt.Sprintf("INV: %d / %d items  %.2f / %.2f kg",
		len(inv.Items), inv.MaxItems, inv.CurrentWeight(), inv.MaxWeight),
		g.hudFont, 4, hudY, dim)
}

// drawCombatNotification renders the centered combat result box for combatFrames.
func (g *Game) drawCombatNotification(screen *ebiten.Image) {
	if g.combatFrames <= 0 {
		return
	}
	lineH := 14
	boxW := float32(220)
	boxH := float32(len(g.combatLines)*lineH + 16)
	boxX := float32(ScreenW)/2 - boxW/2
	boxY := float32(ScreenH)/2 - boxH/2
	border := color.RGBA{100, 110, 140, 255}

	vector.DrawFilledRect(screen, boxX, boxY, boxW, boxH, color.RGBA{20, 22, 30, 220}, false)
	drawBorder(screen, boxX, boxY, boxW, boxH, border)

	lineColors := []color.RGBA{
		{220, 200, 60, 255},  // yellow: hit line
		{200, 120, 120, 255}, // red: enemy HP or defeat
		{224, 108, 117, 255}, // pink: retaliation
	}
	for i, line := range g.combatLines {
		col := lineColors[i%len(lineColors)]
		lx := int(boxX) + int(boxW)/2 - len(line)*3
		ly := int(boxY) + 12 + i*lineH
		text.Draw(screen, line, g.hudFont, lx, ly, col)
	}
}

// --- Shared drawing helpers ---

// drawItemSprite draws item.Image scaled to size×size at (x,y).
// Falls back to item.Color with an optional inset (use 2 for grid cells, 0 otherwise).
func drawItemSprite(screen *ebiten.Image, item *Item, x, y, size, inset float32) {
	if item.Image != nil {
		iop := &ebiten.DrawImageOptions{}
		iw, ih := item.Image.Bounds().Dx(), item.Image.Bounds().Dy()
		iop.GeoM.Scale(float64(size)/float64(iw), float64(size)/float64(ih))
		iop.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(item.Image, iop)
	} else {
		vector.DrawFilledRect(screen, x+inset, y+inset, size-inset*2, size-inset*2, item.Color, false)
	}
}

// drawBorder draws a 1-pixel border around a rectangle.
func drawBorder(screen *ebiten.Image, x, y, w, h float32, col color.RGBA) {
	vector.DrawFilledRect(screen, x, y, w, 1, col, false)
	vector.DrawFilledRect(screen, x, y+h-1, w, 1, col, false)
	vector.DrawFilledRect(screen, x, y, 1, h, col, false)
	vector.DrawFilledRect(screen, x+w-1, y, 1, h, col, false)
}

// drawStatBar draws a horizontal progress bar at (x, y) with height 4.
func drawStatBar(screen *ebiten.Image, x, y, width float32, current, max int, bg, fg color.RGBA) {
	vector.DrawFilledRect(screen, x, y, width, 4, bg, false)
	if max > 0 {
		vector.DrawFilledRect(screen, x, y, width*float32(current)/float32(max), 4, fg, false)
	}
}
