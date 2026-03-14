package game

import (
	"image/color"

	"github.com/ephigenia/ebit-engine-game-1/dungeon"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	TileSize   = 16
	ScreenW    = 640
	ScreenH    = 480
	PlayerSize = 14

	// Key repeat: move on first frame, then after this many frames move every frame while held
	repeatDelayFrames = 4  // ~67ms before repeat starts
	repeatIntervalFrames = 1 // move every frame when holding
)

var (
	colorWall   = color.RGBA{40, 44, 52, 255}
	colorFloor  = color.RGBA{60, 64, 72, 255}
	colorPlayer = color.RGBA{152, 195, 121, 255}
)

// Game implements ebiten.Game.
type Game struct {
	dungeon *dungeon.Dungeon
	playerX int
	playerY int
	cameraX float64
	cameraY float64
	holdFramesUp, holdFramesDown, holdFramesLeft, holdFramesRight int

	// Cached images (created once, reused every frame to avoid allocation and GPU upload)
	wallTileImg  *ebiten.Image
	floorTileImg *ebiten.Image
	playerImg    *ebiten.Image
}

// New creates a new game with a generated dungeon.
func New() *Game {
	cfg := dungeon.DefaultConfig()
	d := dungeon.Generate(cfg)
	g := &Game{dungeon: d}

	// Create tile and player images once (reused every frame)
	g.wallTileImg = ebiten.NewImage(TileSize-1, TileSize-1)
	g.wallTileImg.Fill(colorWall)
	g.floorTileImg = ebiten.NewImage(TileSize-1, TileSize-1)
	g.floorTileImg.Fill(colorFloor)
	g.playerImg = ebiten.NewImage(PlayerSize, PlayerSize)
	g.playerImg.Fill(colorPlayer)

	if len(d.Rooms) > 0 {
		g.playerX, g.playerY = d.Rooms[0].Center()
	} else {
		g.playerX, g.playerY = d.Width/2, d.Height/2
	}
	g.cameraX = float64(g.playerX * TileSize)
	g.cameraY = float64(g.playerY * TileSize)
	return g
}

// Update handles input and moves the player.
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
		nx, ny := g.playerX+dx, g.playerY+dy
		if g.dungeon.IsWalkable(nx, ny) {
			g.playerX, g.playerY = nx, ny
			g.cameraX = float64(g.playerX * TileSize)
			g.cameraY = float64(g.playerY * TileSize)
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
		return true // immediate response on first frame
	}
	if holdFrames < repeatDelayFrames {
		return false
	}
	return (holdFrames-repeatDelayFrames)%repeatIntervalFrames == 0
}

// Draw renders the dungeon and player. Uses cached tile/player images to avoid per-frame allocations.
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

	playerPx := float64(g.playerX*TileSize) + offsetX + float64(TileSize-PlayerSize)/2
	playerPy := float64(g.playerY*TileSize) + offsetY + float64(TileSize-PlayerSize)/2
	op.GeoM.Reset()
	op.GeoM.Translate(playerPx, playerPy)
	screen.DrawImage(g.playerImg, &op)
}

// Layout returns the logical screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenW, ScreenH
}

// Regenerate creates a new dungeon and resets the player (for future "new level" or R key).
func (g *Game) Regenerate() {
	cfg := dungeon.DefaultConfig()
	g.dungeon = dungeon.Generate(cfg)
	if len(g.dungeon.Rooms) > 0 {
		g.playerX, g.playerY = g.dungeon.Rooms[0].Center()
	} else {
		g.playerX, g.playerY = g.dungeon.Width/2, g.dungeon.Height/2
	}
	g.cameraX = float64(g.playerX * TileSize)
	g.cameraY = float64(g.playerY * TileSize)
}
