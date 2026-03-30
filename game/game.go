package game

import (
	"image/color"
	"math/rand"

	"github.com/ephigenia/ebit-engine-game-1/dungeon"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

const (
	TileSize      = 16
	EnemyTileSize = 32
	ScreenW       = 640
	ScreenH       = 480
	PlayerSize    = 14

	repeatDelayFrames    = 8 // ~133ms before repeat starts
	repeatIntervalFrames = 3 // move every 3 frames when holding
	inventoryCols        = 5
)

var (
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
	objects []*Object
	rng     *rand.Rand

	combatLines  []string
	combatFrames int
	particles    ParticleSystem

	inventoryOpen   bool
	inventoryFocus  bool // true = item grid, false = equipment slots
	inventoryCursor int
	equipmentCursor int

	hudFont font.Face

	tilemap    *Spritesheet      // source for wall tiles
	floormap   *Spritesheet      // source for floor tiles (2×2 sheet)
	wallTiles  [16]*ebiten.Image // autotile sprites indexed by neighbor mask (1=up,2=left,4=right,8=down)
	floorTiles [9]*ebiten.Image  // lazily extracted on first draw
	playerImg  *ebiten.Image
	enemyImg   *ebiten.Image
	objectImg  *ebiten.Image // shared spritesheet for animated chest types
}

// Layout returns the logical screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenW, ScreenH
}
