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
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	colorPotion = color.RGBA{229, 192, 123, 255}
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

	// Per-tile floor variant index, dimensions match dungeon grid
	floorVariants [][]int

	// Cached images (created once, reused every frame to avoid allocation and GPU upload)
	wallTileImg   *ebiten.Image
	floorTileImgs []*ebiten.Image
	playerImg     *ebiten.Image
	enemyImg      *ebiten.Image
	potionImg     *ebiten.Image
}

// loadImage decodes a PNG from the given FS path into an *ebiten.Image.
func loadImage(fsys fs.FS, path string) *ebiten.Image {
	f, err := fsys.Open(path)
	if err != nil {
		log.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatalf("decode %s: %v", path, err)
	}
	return ebiten.NewImageFromImage(img)
}

// New creates a new game with a generated dungeon.
func New(assets fs.FS) *Game {
	cfg := dungeon.DefaultConfig()
	d := dungeon.Generate(cfg)
	g := &Game{dungeon: d}

	// Create tile and entity images once (reused every frame)
	g.wallTileImg = ebiten.NewImage(TileSize-1, TileSize-1)
	g.wallTileImg.Fill(colorWall)
	for i := 1; i <= 8; i++ {
		g.floorTileImgs = append(g.floorTileImgs, loadImage(assets, fmt.Sprintf("assets/grass_tile_%d.png", i)))
	}
	g.playerImg = ebiten.NewImage(PlayerSize, PlayerSize)
	g.playerImg.Fill(colorPlayer)
	g.enemyImg = ebiten.NewImage(PlayerSize, PlayerSize)
	g.enemyImg.Fill(colorEnemy)
	g.potionImg = ebiten.NewImage(8, 8)
	g.potionImg.Fill(colorPotion)

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
	g.player = newPlayer(startX, startY)
	g.cameraX = float64(startX * TileSize)
	g.cameraY = float64(startY * TileSize)

	rng := rand.New(rand.NewSource(rand.Int63()))

	// Assign a random floor tile variant to every cell in the grid
	g.floorVariants = make([][]int, d.Height)
	for y := 0; y < d.Height; y++ {
		g.floorVariants[y] = make([]int, d.Width)
		for x := 0; x < d.Width; x++ {
			g.floorVariants[y][x] = rng.Intn(len(g.floorTileImgs))
		}
	}

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
			if e.IsAlive() {
				g.player.TakeDamage(e.Attack)
			}
		} else if g.dungeon.IsWalkable(nx, ny) {
			g.player.X, g.player.Y = nx, ny
			g.cameraX = float64(g.player.X * TileSize)
			g.cameraY = float64(g.player.Y * TileSize)
			if p := g.potionAt(nx, ny); p != nil {
				p.Taken = true
				g.player.HP += p.Heal
				if g.player.HP > g.player.MaxHP {
					g.player.HP = g.player.MaxHP
				}
			}
		}
	}
	return nil
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
				v := g.floorVariants[ty][tx]
				screen.DrawImage(g.floorTileImgs[v], &op)
			}
		}
	}

	// Draw untaken potions
	const potionSize = 8
	for _, p := range g.potions {
		if p.Taken {
			continue
		}
		px := float64(p.X*TileSize) + offsetX + float64(TileSize-potionSize)/2
		py := float64(p.Y*TileSize) + offsetY + float64(TileSize-potionSize)/2
		op.GeoM.Reset()
		op.GeoM.Translate(px, py)
		screen.DrawImage(g.potionImg, &op)
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

	// HUD: player health
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP: %d / %d", g.player.HP, g.player.MaxHP), 4, 4)

	// Draw player
	playerPx := float64(g.player.X*TileSize) + offsetX + float64(TileSize-PlayerSize)/2
	playerPy := float64(g.player.Y*TileSize) + offsetY + float64(TileSize-PlayerSize)/2
	op.GeoM.Reset()
	op.GeoM.Translate(playerPx, playerPy)
	screen.DrawImage(g.playerImg, &op)
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
