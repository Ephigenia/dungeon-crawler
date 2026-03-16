package game

import (
	"fmt"
	"image/color"
	"io/fs"
	"log"
	"math/rand"

	"github.com/ephigenia/ebit-engine-game-1/dungeon"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	TileSize   = 16
	ScreenW    = 640
	ScreenH    = 480
	PlayerSize = 14

	// Key repeat: move on first frame, then after this many frames move every frame while held
	repeatDelayFrames    = 4 // ~67ms before repeat starts
	repeatIntervalFrames = 1 // move every frame when holding
)

var (
	colorWall   = color.RGBA{40, 44, 52, 255}
	colorFloor  = color.RGBA{60, 64, 72, 255}
	colorPlayer = color.RGBA{152, 195, 121, 255}
	colorEnemy  = color.RGBA{224, 108, 117, 255}
)

// Game implements ebiten.Game.
type Game struct {
	dungeon                                                       *dungeon.Dungeon
	player                                                        *Player
	enemies                                                       []*Enemy
	cameraX                                                       float64
	cameraY                                                       float64
	holdFramesUp, holdFramesDown, holdFramesLeft, holdFramesRight int

	potions []*Potion

	inventoryOpen   bool
	inventoryCursor int

	hudFont font.Face

	// Cached images (created once, reused every frame to avoid allocation and GPU upload)
	wallTileImg  *ebiten.Image
	floorTileImg *ebiten.Image
	playerImg    *ebiten.Image
	enemyImg     *ebiten.Image
}

// New creates a new game with a generated dungeon.
func New(assets fs.FS) *Game {
	cfg := dungeon.DefaultConfig()
	d := dungeon.Generate(cfg)
	g := &Game{dungeon: d}

	// Load HUD font
	fontData, err := fs.ReadFile(assets, "assets/Gorgeous-Pixel/GorgeousPixel.ttf")
	if err != nil {
		log.Fatalf("read font: %v", err)
	}
	tt, err := opentype.Parse(fontData)
	if err != nil {
		log.Fatalf("parse font: %v", err)
	}
	g.hudFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    12,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatalf("new font face: %v", err)
	}

	// Create tile and entity images once (reused every frame)
	g.wallTileImg = ebiten.NewImage(TileSize-1, TileSize-1)
	g.wallTileImg.Fill(colorWall)
	g.floorTileImg = ebiten.NewImage(TileSize-1, TileSize-1)
	g.floorTileImg.Fill(colorFloor)
	g.playerImg = ebiten.NewImage(PlayerSize, PlayerSize)
	g.playerImg.Fill(colorPlayer)
	g.enemyImg = ebiten.NewImage(PlayerSize, PlayerSize)
	g.enemyImg.Fill(colorEnemy)

	g.resetEntities(d)
	return g
}

// resetEntities places the player and spawns enemies for the given dungeon.
func (g *Game) resetEntities(d *dungeon.Dungeon) {
	var startX, startY int
	if len(d.Rooms) > 0 {
		startX, startY = d.Rooms[0].Center()
	} else {
		startX, startY = d.Width/2, d.Height/2
	}
	if g.player == nil {
		g.player = newPlayer(startX, startY)
	} else {
		g.player.X, g.player.Y = startX, startY
	}
	g.cameraX = float64(startX * TileSize)
	g.cameraY = float64(startY * TileSize)

	rng := rand.New(rand.NewSource(rand.Int63()))

	g.enemies = g.enemies[:0]
	// Spawn one enemy per room, skipping the starting room
	for i := 1; i < len(d.Rooms); i++ {
		ex, ey := d.Rooms[i].Center()
		g.enemies = append(g.enemies, spawnEnemy(ex, ey, rng))
	}

	// Spawn 0-2 potions per room at random offsets from the center
	g.potions = g.potions[:0]
	for _, room := range d.Rooms {
		count := rng.Intn(3) // 0, 1, or 2
		for n := 0; n < count; n++ {
			cx, cy := room.Center()
			// Random offset within the room, avoiding the exact center (enemy/player tile)
			ox := rng.Intn(room.W) - room.W/2
			oy := rng.Intn(room.H) - room.H/2
			if ox == 0 && oy == 0 {
				ox = 1
			}
			px, py := cx+ox, cy+oy
			if d.IsWalkable(px, py) {
				g.potions = append(g.potions, newPotion(px, py, rng))
			}
		}
	}
}

// potionAt returns the untaken potion at (x, y), or nil.
func (g *Game) potionAt(x, y int) *Potion {
	for _, p := range g.potions {
		if !p.Taken && p.X == x && p.Y == y {
			return p
		}
	}
	return nil
}

// enemyAt returns the living enemy at (x, y), or nil.
func (g *Game) enemyAt(x, y int) *Enemy {
	for _, e := range g.enemies {
		if e.IsAlive() && e.X == x && e.Y == y {
			return e
		}
	}
	return nil
}

// Update handles input, combat, and movement.
func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyI) {
		g.inventoryOpen = !g.inventoryOpen
		g.inventoryCursor = 0
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
	if dx != 0 || dy != 0 {
		nx, ny := g.player.X+dx, g.player.Y+dy
		if e := g.enemyAt(nx, ny); e != nil {
			// Bump attack: player hits enemy, enemy retaliates
			e.TakeDamage(g.player.Attack)
			g.player.AddEXP(5)
			if e.IsAlive() {
				g.player.TakeDamage(e.Attack)
			} else {
				g.player.AddEXP(20)
			}
		} else if g.dungeon.IsWalkable(nx, ny) {
			g.player.X, g.player.Y = nx, ny
			g.cameraX = float64(g.player.X * TileSize)
			g.cameraY = float64(g.player.Y * TileSize)
			if p := g.potionAt(nx, ny); p != nil {
				if g.player.Inventory.Add(p.Item) {
					p.Taken = true
				}
			}
		}
	}
	return nil
}

const inventoryCols = 5

// updateInventory handles input while the inventory screen is open.
func (g *Game) updateInventory() {
	inv := g.player.Inventory
	maxSlots := inv.MaxItems

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
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

	if inpututil.IsKeyJustPressed(ebiten.KeyU) && g.inventoryCursor < len(inv.Items) {
		item := inv.Items[g.inventoryCursor]
		if item.Category == CategoryConsumable && item.OnUse != nil {
			if item.OnUse(g.player) {
				inv.Remove(g.inventoryCursor)
				if g.inventoryCursor >= len(inv.Items) && g.inventoryCursor > 0 {
					g.inventoryCursor--
				}
			}
		}
	}
}

// shouldMove returns true on the first frame the key is pressed, then after repeatDelayFrames, every repeatIntervalFrames.
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

// Draw renders the dungeon, enemies, and player.
func (g *Game) Draw(screen *ebiten.Image) {
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

	// Draw untaken potions
	const potionSize = 8
	for _, p := range g.potions {
		if p.Taken {
			continue
		}
		px := float32(float64(p.X*TileSize) + offsetX + float64(TileSize-potionSize)/2)
		py := float32(float64(p.Y*TileSize) + offsetY + float64(TileSize-potionSize)/2)
		vector.DrawFilledRect(screen, px, py, potionSize, potionSize, p.Item.Color, false)
	}

	// Draw living enemies with HP bar
	for _, e := range g.enemies {
		if !e.IsAlive() {
			continue
		}
		ex := float64(e.X*TileSize) + offsetX + float64(TileSize-PlayerSize)/2
		ey := float64(e.Y*TileSize) + offsetY + float64(TileSize-PlayerSize)/2
		op.GeoM.Reset()
		op.GeoM.Translate(ex, ey)
		screen.DrawImage(g.enemyImg, &op)

		// 1px HP bar just below the sprite
		barY := float32(ey) + PlayerSize + 1
		barX := float32(ex)
		vector.DrawFilledRect(screen, barX, barY, PlayerSize, 1, color.RGBA{30, 30, 30, 200}, false)
		hpWidth := float32(PlayerSize) * float32(e.HP) / float32(e.MaxHP)
		vector.DrawFilledRect(screen, barX, barY, hpWidth, 1, color.RGBA{224, 108, 117, 255}, false)
	}

	// HUD: player stats
	inv := g.player.Inventory
	text.Draw(screen, fmt.Sprintf("HP: %d / %d", g.player.HP, g.player.MaxHP), g.hudFont, 4, 14, color.White)
	text.Draw(screen, fmt.Sprintf("LVL: %d  EXP: %d / %d", g.player.Level, g.player.EXP, g.player.NextLevelEXP), g.hudFont, 4, 28, color.White)
	text.Draw(screen, fmt.Sprintf("INV: %d / %d items  %.2f / %.2f kg", len(inv.Items), inv.MaxItems, inv.CurrentWeight(), inv.MaxWeight), g.hudFont, 4, 42, color.White)

	// Draw player
	playerPx := float64(g.player.X*TileSize) + offsetX + float64(TileSize-PlayerSize)/2
	playerPy := float64(g.player.Y*TileSize) + offsetY + float64(TileSize-PlayerSize)/2
	op.GeoM.Reset()
	op.GeoM.Translate(playerPx, playerPy)
	screen.DrawImage(g.playerImg, &op)

	if g.inventoryOpen {
		g.drawInventory(screen)
	}
}

// drawInventory renders the inventory overlay.
func (g *Game) drawInventory(screen *ebiten.Image) {
	const (
		panelX      = 30
		panelY      = 30
		panelW      = 580
		panelH      = 380
		slotSize    = 16
		slotGap     = 4
		slotStride  = slotSize + slotGap
		gridX       = panelX + 20
		gridY       = panelY + 50
		detailX     = gridX + inventoryCols*slotStride + 30
	)

	// Dark overlay + panel background
	vector.DrawFilledRect(screen, 0, 0, ScreenW, ScreenH, color.RGBA{0, 0, 0, 180}, false)
	vector.DrawFilledRect(screen, panelX, panelY, panelW, panelH, color.RGBA{28, 32, 42, 255}, false)
	vector.DrawFilledRect(screen, panelX, panelY, panelW, 1, color.RGBA{100, 110, 140, 255}, false)
	vector.DrawFilledRect(screen, panelX, panelY+panelH, panelW, 1, color.RGBA{100, 110, 140, 255}, false)
	vector.DrawFilledRect(screen, panelX, panelY, 1, panelH, color.RGBA{100, 110, 140, 255}, false)
	vector.DrawFilledRect(screen, panelX+panelW, panelY, 1, panelH, color.RGBA{100, 110, 140, 255}, false)

	text.Draw(screen, "INVENTORY", g.hudFont, panelX+10, panelY+18, color.White)
	inv := g.player.Inventory
	infoStr := fmt.Sprintf("%d / %d items   %.2f / %.2f kg", len(inv.Items), inv.MaxItems, inv.CurrentWeight(), inv.MaxWeight)
	text.Draw(screen, infoStr, g.hudFont, panelX+10, panelY+34, color.RGBA{160, 160, 160, 255})

	// Slot grid
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

			selected := idx == g.inventoryCursor
			bgCol := color.RGBA{45, 50, 62, 255}
			if selected {
				bgCol = color.RGBA{70, 90, 115, 255}
			}
			vector.DrawFilledRect(screen, sx, sy, slotSize, slotSize, bgCol, false)

			// Item fill
			if idx < len(inv.Items) {
				vector.DrawFilledRect(screen, sx+2, sy+2, slotSize-4, slotSize-4, inv.Items[idx].Color, false)
			}

			// Border
			borderCol := color.RGBA{80, 88, 108, 255}
			if selected {
				borderCol = color.RGBA{180, 200, 255, 255}
			}
			vector.DrawFilledRect(screen, sx, sy, slotSize, 1, borderCol, false)
			vector.DrawFilledRect(screen, sx, sy+slotSize-1, slotSize, 1, borderCol, false)
			vector.DrawFilledRect(screen, sx, sy, 1, slotSize, borderCol, false)
			vector.DrawFilledRect(screen, sx+slotSize-1, sy, 1, slotSize, borderCol, false)
		}
	}

	// Item detail panel
	dy := gridY + 12
	if g.inventoryCursor < len(inv.Items) {
		item := inv.Items[g.inventoryCursor]
		text.Draw(screen, item.ID, g.hudFont, detailX, dy, color.White)
		dy += 20
		text.Draw(screen, fmt.Sprintf("Type:   %s", item.Category), g.hudFont, detailX, dy, color.RGBA{160, 160, 160, 255})
		dy += 16
		text.Draw(screen, fmt.Sprintf("Weight: %.1f kg", item.Weight), g.hudFont, detailX, dy, color.RGBA{160, 160, 160, 255})
		dy += 16
		if item.Effect != "" {
			text.Draw(screen, fmt.Sprintf("Effect: %s", item.Effect), g.hudFont, detailX, dy, color.RGBA{152, 210, 152, 255})
			dy += 16
		}
		if item.Category == CategoryConsumable {
			dy += 4
			text.Draw(screen, "[U] Use", g.hudFont, detailX, dy, color.RGBA{220, 210, 100, 255})
		}
	} else {
		text.Draw(screen, "Empty slot", g.hudFont, detailX, dy, color.RGBA{80, 80, 80, 255})
	}

	// Controls hint
	text.Draw(screen, "[Arrows/WASD] Navigate   [U] Use   [I] Close", g.hudFont, panelX+10, panelY+panelH-14, color.RGBA{120, 120, 120, 255})
}


// Layout returns the logical screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenW, ScreenH
}

// Regenerate creates a new dungeon and resets all entities.
func (g *Game) Regenerate() {
	cfg := dungeon.DefaultConfig()
	g.dungeon = dungeon.Generate(cfg)
	g.resetEntities(g.dungeon)
}
