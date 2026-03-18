package game

import (
	"image"
	_ "image/png"
	"io/fs"
	"log"
	"math/rand"

	"github.com/ephigenia/ebit-engine-game-1/dungeon"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// New creates a new game with a generated dungeon.
func New(assets fs.FS) *Game {
	cfg := dungeon.DefaultConfig()
	d := dungeon.Generate(cfg)
	g := &Game{dungeon: d, rng: rand.New(rand.NewSource(rand.Int63()))}

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

	loadItemImages(assets)
	g.resetEntities(d)
	return g
}

// resetEntities places the player and spawns enemies/items for the given dungeon.
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
	for i := 1; i < len(d.Rooms); i++ {
		ex, ey := d.Rooms[i].Center()
		g.enemies = append(g.enemies, spawnEnemy(ex, ey, g.rng))
	}

	g.potions = g.potions[:0]
	for _, room := range d.Rooms {
		count := g.rng.Intn(4) + 1
		for n := 0; n < count; n++ {
			cx, cy := room.Center()
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

// potionAt returns the untaken pickup at (x, y), or nil.
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

// Regenerate creates a new dungeon and resets all entities.
func (g *Game) Regenerate() {
	cfg := dungeon.DefaultConfig()
	g.dungeon = dungeon.Generate(cfg)
	g.resetEntities(g.dungeon)
}
