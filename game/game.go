package game

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
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
	rng     *rand.Rand

	inventoryOpen   bool
	inventoryFocus  bool // true = item grid, false = equipment slots
	inventoryCursor int
	equipmentCursor int

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
	g := &Game{dungeon: d, rng: rand.New(rand.NewSource(rand.Int63()))}

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
	if f, err := assets.Open("assets/player/player.png"); err == nil {
		if img, _, err := image.Decode(f); err == nil {
			g.playerImg = ebiten.NewImageFromImage(img)
		}
		f.Close()
	}
	if g.playerImg == nil {
		g.playerImg = ebiten.NewImage(PlayerSize, PlayerSize)
		g.playerImg.Fill(colorPlayer)
	}
	g.enemyImg = ebiten.NewImage(PlayerSize, PlayerSize)
	g.enemyImg.Fill(colorEnemy)

	// Load item sprites
	for _, entry := range []struct {
		path string
		item *Item
	}{
		{"assets/items/health_potion_large.png", ItemLargeHealthPotion},
		{"assets/items/health_potion_medium.png", ItemMediumHealthPotion},
		{"assets/items/health_potion_small.png", ItemSmallHealthPotion},
		// food
		{"assets/items/food/apple.png", ItemApple},
		{"assets/items/food/bread_roll.png", ItemBreadRoll},
		{"assets/items/food/egg_fried.png", ItemFriedEgg},
		{"assets/items/food/grapes.png", ItemGrapes},
		{"assets/items/food/meat.png", ItemMeat},
		{"assets/items/food/mushroom.png", ItemMushroom},
		{"assets/items/food/pizza_slice.png", ItemPizzaSlice},
		// items // gear // backpack
		{"assets/items/gear/backpack/basic.png", ItemSmallBackpack},
		{"assets/items/gear/backpack/medium.png", ItemMediumBackpack},
		{"assets/items/gear/backpack/large.png", ItemLargeBackpack},
		// items // gear // legs
		{"assets/items/gear/legs/pants.png", ItemPants},
		// items // accessories
		{"assets/items/accessories/necklace_diamond.png", ItemNecklaceDiamond},
		{"assets/items/accessories/necklace_skull.png", ItemNecklaceSkull},
		{"assets/items/accessories/necklace_star.png", ItemNecklaceStar},
		{"assets/items/accessories/necklace_tooth.png", ItemNecklaceTooth},
		{"assets/items/accessories/ring_diamond.png", ItemDiamondRing},
		{"assets/items/accessories/ring_diamond_2.png", ItemDiamondRing2},
		{"assets/items/accessories/ring_gold.png", ItemGoldRing},
		{"assets/items/accessories/ring_silver.png", ItemSilverRing},
		// items // gear // gloves
		{"assets/items/gear/gloves/gloves_finger.png", ItemGlovesFinger},
		{"assets/items/gear/gloves/gloves_leather_metal.png", ItemGlovesLeatherMetal},
		{"assets/items/gear/gloves/gloves_leather.png", ItemGlovesLeather},
		{"assets/items/gear/gloves/gloves_metal.png", ItemGlovesMetal},
		// items // gear // head
		{"assets/items/gear/head/basic_helmet.png", ItemBasicHelmet},
		{"assets/items/gear/head/coif.png", ItemCoif},
		{"assets/items/gear/head/full_helmet.png", ItemFullHelmet},
		{"assets/items/gear/head/gold_helmet.png", ItemGoldHelmet},
		{"assets/items/gear/head/horn_helmet.png", ItemHornHelmet},
		// items // gear // shoes
		{"assets/items/gear/shoes/shoes_gold.png", ItemGoldShoes},
		{"assets/items/gear/shoes/shoes_leather.png", ItemLeatherShoes},
		{"assets/items/gear/shoes/shoes_metal.png", ItemMetalShoes},
		{"assets/items/gear/shoes/shoes_simple.png", ItemSimpleShoes},
		// items // gear // shield
		{"assets/items/shield/shield_metal.png", ItemMetalShield},
		{"assets/items/shield/shield_wood.png", ItemWoodenShield},
		// items // gear // armor
		{"assets/items/armor/basic.png", ItemBasicArmor},
		{"assets/items/armor/bronze.png", ItemBronzeArmor},
		{"assets/items/armor/complex.png", ItemComplexArmor},
		{"assets/items/armor/gold.png", ItemGoldArmor},
	} {
		if f, err := assets.Open(entry.path); err == nil {
			if img, _, err := image.Decode(f); err == nil {
				entry.item.Image = ebiten.NewImageFromImage(img)
			}
			f.Close()
		}
	}

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

	g.enemies = g.enemies[:0]
	// Spawn one enemy per room, skipping the starting room
	for i := 1; i < len(d.Rooms); i++ {
		ex, ey := d.Rooms[i].Center()
		g.enemies = append(g.enemies, spawnEnemy(ex, ey, g.rng))
	}

	// Spawn 0-2 potions per room at random offsets from the center
	g.potions = g.potions[:0]
	for _, room := range d.Rooms {
		count := g.rng.Intn(4) + 1 // 1, 2, 3, or 4
		for n := 0; n < count; n++ {
			cx, cy := room.Center()
			// Random offset within the room, avoiding the exact center (enemy/player tile)
			ox := g.rng.Intn(room.W) - room.W/2
			oy := g.rng.Intn(room.H) - room.H/2
			if ox == 0 && oy == 0 {
				ox = 1
			}
			px, py := cx+ox, cy+oy
			if d.IsWalkable(px, py) {
				g.potions = append(g.potions, newPotion(px, py, g.rng))
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
				dropChance := 10 + (g.player.Level - 1)
				if g.rng.Intn(100) < dropChance {
					drop := newPotion(e.X, e.Y, g.rng)
					g.potions = append(g.potions, drop)
				}
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
	// Tab switches focus between item grid and equipment slots.
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

	usePressed := inpututil.IsKeyJustPressed(ebiten.KeyU) || inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	if usePressed && g.inventoryCursor < len(inv.Items) {
		item := inv.Items[g.inventoryCursor]
		switch item.Category {
		case CategoryConsumable:
			if item.OnUse != nil && item.OnUse(g.player) {
				inv.Remove(g.inventoryCursor)
				if g.inventoryCursor >= len(inv.Items) && g.inventoryCursor > 0 {
					g.inventoryCursor--
				}
			}
		case CategoryEquipment, CategoryBackpack:
			if g.player.IsEquipped(item) {
				// Toggle off: find the slot and unequip
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
		item := inv.Items[g.inventoryCursor]
		if !g.player.IsEquipped(item) {
			inv.Remove(g.inventoryCursor)
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
	const potionSize = 16
	for _, p := range g.potions {
		if p.Taken {
			continue
		}
		px := float32(float64(p.X*TileSize) + offsetX + float64(TileSize-potionSize)/2)
		py := float32(float64(p.Y*TileSize) + offsetY + float64(TileSize-potionSize)/2)
		if p.Item.Image != nil {
			iop := &ebiten.DrawImageOptions{}
			iw, ih := p.Item.Image.Bounds().Dx(), p.Item.Image.Bounds().Dy()
			iop.GeoM.Scale(float64(potionSize)/float64(iw), float64(potionSize)/float64(ih))
			iop.GeoM.Translate(float64(px), float64(py))
			screen.DrawImage(p.Item.Image, iop)
		} else {
			vector.DrawFilledRect(screen, px, py, potionSize, potionSize, p.Item.Color, false)
		}
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
	hudDim := color.RGBA{160, 160, 160, 255}
	hudY := 14

	text.Draw(screen, "HP", g.hudFont, 4, hudY, hudDim)
	text.Draw(screen, fmt.Sprintf("%d / %d", g.player.HP, g.player.MaxHP), g.hudFont, 32, hudY, color.White)
	vector.DrawFilledRect(screen, 94, float32(hudY-9), 80, 4, color.RGBA{50, 20, 20, 220}, false)
	vector.DrawFilledRect(screen, 94, float32(hudY-9), 80*float32(g.player.HP)/float32(g.player.MaxHP), 4, color.RGBA{200, 60, 60, 255}, false)
	hudY += 12

	text.Draw(screen, "ATK", g.hudFont, 4, hudY, hudDim)
	text.Draw(screen, fmt.Sprintf("%d", g.player.Attack), g.hudFont, 32, hudY, color.RGBA{224, 180, 100, 255})
	hudY += 12

	text.Draw(screen, "DEF", g.hudFont, 4, hudY, hudDim)
	text.Draw(screen, fmt.Sprintf("%d", g.player.Defense), g.hudFont, 32, hudY, color.RGBA{100, 160, 220, 255})
	hudY += 12

	text.Draw(screen, "LVL", g.hudFont, 4, hudY, hudDim)
	text.Draw(screen, fmt.Sprintf("%d", g.player.Level), g.hudFont, 32, hudY, color.White)
	hudY += 12

	text.Draw(screen, "EXP", g.hudFont, 4, hudY, hudDim)
	text.Draw(screen, fmt.Sprintf("%d / %d", g.player.EXP, g.player.NextLevelEXP), g.hudFont, 32, hudY, color.White)
	vector.DrawFilledRect(screen, 94, float32(hudY-9), 80, 4, color.RGBA{20, 30, 50, 220}, false)
	vector.DrawFilledRect(screen, 94, float32(hudY-9), 80*float32(g.player.EXP)/float32(g.player.NextLevelEXP), 4, color.RGBA{100, 160, 240, 255}, false)
	hudY += 12

	text.Draw(screen, fmt.Sprintf("INV: %d / %d items  %.2f / %.2f kg", len(inv.Items), inv.MaxItems, inv.CurrentWeight(), inv.MaxWeight), g.hudFont, 4, hudY, hudDim)

	// Draw player
	playerPx := float64(g.player.X*TileSize) + offsetX + float64(TileSize-PlayerSize)/2
	playerPy := float64(g.player.Y*TileSize) + offsetY + float64(TileSize-PlayerSize)/2
	op.GeoM.Reset()
	iw, ih := g.playerImg.Bounds().Dx(), g.playerImg.Bounds().Dy()
	op.GeoM.Scale(float64(PlayerSize)/float64(iw), float64(PlayerSize)/float64(ih))
	op.GeoM.Translate(playerPx, playerPy)
	screen.DrawImage(g.playerImg, &op)

	if g.inventoryOpen {
		g.drawInventory(screen)
	}
}

// drawInventory renders the inventory overlay with item grid and equipment slots.
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

	// Overlay + panel
	vector.DrawFilledRect(screen, 0, 0, ScreenW, ScreenH, color.RGBA{0, 0, 0, 180}, false)
	vector.DrawFilledRect(screen, panelX, panelY, panelW, panelH, bg, false)
	vector.DrawFilledRect(screen, panelX, panelY, panelW, 1, border, false)
	vector.DrawFilledRect(screen, panelX, panelY+panelH, panelW, 1, border, false)
	vector.DrawFilledRect(screen, panelX, panelY, 1, panelH, border, false)
	vector.DrawFilledRect(screen, panelX+panelW, panelY, 1, panelH, border, false)

	// Vertical divider
	vector.DrawFilledRect(screen, dividerX, panelY+40, 1, panelH-40, border, false)
	// Horizontal divider above detail area
	vector.DrawFilledRect(screen, panelX, detailY-8, panelW, 1, border, false)

	// --- Section headers ---
	white := color.RGBA{255, 255, 255, 255}
	itemsFocusColor := dim
	equipFocusColor := dim
	if g.inventoryFocus {
		itemsFocusColor = white
	} else {
		equipFocusColor = white
	}
	inv := g.player.Inventory
	text.Draw(screen, fmt.Sprintf("ITEMS  %d/%d  %.1f/%.1fkg", len(inv.Items), inv.MaxItems, inv.CurrentWeight(), inv.MaxWeight),
		g.hudFont, gridX, panelY+18, itemsFocusColor)
	text.Draw(screen, "EQUIPMENT", g.hudFont, equipX, panelY+18, equipFocusColor)

	// --- Item grid ---
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
			if idx < len(inv.Items) {
				item := inv.Items[idx]
				if item.Image != nil {
					iop := &ebiten.DrawImageOptions{}
					iw, ih := item.Image.Bounds().Dx(), item.Image.Bounds().Dy()
					iop.GeoM.Scale(float64(slotSize)/float64(iw), float64(slotSize)/float64(ih))
					iop.GeoM.Translate(float64(sx), float64(sy))
					screen.DrawImage(item.Image, iop)
				} else {
					vector.DrawFilledRect(screen, sx+2, sy+2, slotSize-4, slotSize-4, item.Color, false)
				}
			}
			borderCol := color.RGBA{80, 88, 108, 255}
			if idx < len(inv.Items) && g.player.IsEquipped(inv.Items[idx]) {
				borderCol = color.RGBA{220, 200, 60, 255}
			}
			if selected {
				borderCol = color.RGBA{180, 200, 255, 255}
			}
			vector.DrawFilledRect(screen, sx, sy, slotSize, 1, borderCol, false)
			vector.DrawFilledRect(screen, sx, sy+slotSize-1, slotSize, 1, borderCol, false)
			vector.DrawFilledRect(screen, sx, sy, 1, slotSize, borderCol, false)
			vector.DrawFilledRect(screen, sx+slotSize-1, sy, 1, slotSize, borderCol, false)
		}
	}

	// --- Character stats (below item grid) ---
	statsY := gridY + rows*slotStride + 16
	vector.DrawFilledRect(screen, float32(gridX), float32(statsY-7), float32(dividerX-gridX-12), 1, border, false)
	text.Draw(screen, "CHARACTER", g.hudFont, gridX, statsY+2, dim)
	statsY += 16

	text.Draw(screen, "HP", g.hudFont, gridX, statsY, dim)
	text.Draw(screen, fmt.Sprintf("%d / %d", g.player.HP, g.player.MaxHP), g.hudFont, gridX+28, statsY, white)
	vector.DrawFilledRect(screen, float32(gridX+90), float32(statsY-9), 80, 4, color.RGBA{50, 20, 20, 220}, false)
	vector.DrawFilledRect(screen, float32(gridX+90), float32(statsY-9), 80*float32(g.player.HP)/float32(g.player.MaxHP), 4, color.RGBA{200, 60, 60, 255}, false)
	statsY += 12

	text.Draw(screen, "ATK", g.hudFont, gridX, statsY, dim)
	text.Draw(screen, fmt.Sprintf("%d", g.player.Attack), g.hudFont, gridX+28, statsY, color.RGBA{224, 180, 100, 255})
	statsY += 12

	text.Draw(screen, "DEF", g.hudFont, gridX, statsY, dim)
	text.Draw(screen, fmt.Sprintf("%d", g.player.Defense), g.hudFont, gridX+28, statsY, color.RGBA{100, 160, 220, 255})
	statsY += 12

	text.Draw(screen, "LVL", g.hudFont, gridX, statsY, dim)
	text.Draw(screen, fmt.Sprintf("%d", g.player.Level), g.hudFont, gridX+28, statsY, white)
	statsY += 12

	text.Draw(screen, "EXP", g.hudFont, gridX, statsY, dim)
	text.Draw(screen, fmt.Sprintf("%d / %d", g.player.EXP, g.player.NextLevelEXP), g.hudFont, gridX+28, statsY, white)
	vector.DrawFilledRect(screen, float32(gridX+90), float32(statsY-9), 80, 4, color.RGBA{20, 30, 50, 220}, false)
	vector.DrawFilledRect(screen, float32(gridX+90), float32(statsY-9), 80*float32(g.player.EXP)/float32(g.player.NextLevelEXP), 4, color.RGBA{100, 160, 240, 255}, false)

	// --- Equipment slot list ---
	// columns: label | swatch | name | weight | effect
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
			if equipped.Image != nil {
				iop := &ebiten.DrawImageOptions{}
				iw, ih := equipped.Image.Bounds().Dx(), equipped.Image.Bounds().Dy()
				iop.GeoM.Scale(6.0/float64(iw), 6.0/float64(ih))
				iop.GeoM.Translate(float64(colSwatch), float64(ey-8))
				screen.DrawImage(equipped.Image, iop)
			} else {
				vector.DrawFilledRect(screen, float32(colSwatch), float32(ey-8), 6, 6, equipped.Color, false)
			}
			text.Draw(screen, equipped.ID, g.hudFont, colName, ey, equipped.Color)
			text.Draw(screen, fmt.Sprintf("%.1fkg", equipped.Weight), g.hudFont, colWeight, ey, dim)
			if equipped.Effect != "" {
				text.Draw(screen, equipped.Effect, g.hudFont, colEffect, ey, green)
			}
		} else {
			text.Draw(screen, "(empty)", g.hudFont, colName, ey, color.RGBA{70, 70, 70, 255})
		}
	}

	// --- Detail area (shared) ---
	var selectedItem *Item
	var fromEquipment bool
	if g.inventoryFocus {
		if g.inventoryCursor < len(inv.Items) {
			selectedItem = inv.Items[g.inventoryCursor]
		}
	} else {
		slot := EquipmentSlotOrder[g.equipmentCursor]
		selectedItem = g.player.Equipment.Slots[slot]
		fromEquipment = true
	}

	dy := detailY + 10
	if selectedItem != nil {
		text.Draw(screen, selectedItem.ID, g.hudFont, gridX, dy, white)
		dy += 16
		text.Draw(screen, fmt.Sprintf("Type: %s   Weight: %.1f kg", selectedItem.Category, selectedItem.Weight), g.hudFont, gridX, dy, dim)
		dy += 14
		if selectedItem.Effect != "" {
			text.Draw(screen, fmt.Sprintf("Effect: %s", selectedItem.Effect), g.hudFont, gridX, dy, green)
			dy += 14
		}
		if selectedItem.StatMods != (StatModifiers{}) {
			mods := selectedItem.StatMods
			modsStr := ""
			if mods.Attack != 0 {
				modsStr += fmt.Sprintf("ATK %+d  ", mods.Attack)
			}
			if mods.Defense != 0 {
				modsStr += fmt.Sprintf("DEF %+d  ", mods.Defense)
			}
			if mods.HP != 0 {
				modsStr += fmt.Sprintf("HP %+d", mods.HP)
			}
			text.Draw(screen, modsStr, g.hudFont, gridX, dy, green)
			dy += 14
		}
		dy += 4
		if fromEquipment {
			text.Draw(screen, "[U/Enter] Unequip", g.hudFont, gridX, dy, yellow)
		} else {
			switch selectedItem.Category {
			case CategoryConsumable:
				text.Draw(screen, "[U/Enter] Use    [X] Destroy", g.hudFont, gridX, dy, yellow)
			case CategoryEquipment, CategoryBackpack:
				if g.player.IsEquipped(selectedItem) {
					text.Draw(screen, "[U/Enter] Unequip", g.hudFont, gridX, dy, yellow)
				} else {
					text.Draw(screen, "[U/Enter] Equip  [X] Destroy", g.hudFont, gridX, dy, yellow)
				}
			default:
				text.Draw(screen, "[X] Destroy", g.hudFont, gridX, dy, red)
			}
		}
	} else {
		if fromEquipment {
			text.Draw(screen, "(slot empty)", g.hudFont, gridX, dy, color.RGBA{70, 70, 70, 255})
		} else {
			text.Draw(screen, "(empty slot)", g.hudFont, gridX, dy, color.RGBA{70, 70, 70, 255})
		}
	}

	// Controls hint
	text.Draw(screen, "[Tab] Switch   [Arrows/WASD] Navigate   [U/Enter] Action   [X] Destroy   [I] Close",
		g.hudFont, panelX+6, panelY+panelH-10, color.RGBA{100, 100, 100, 255})
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
