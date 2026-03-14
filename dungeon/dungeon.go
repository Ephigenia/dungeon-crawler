package dungeon

import (
	"math/rand"
)

// Tile represents a single cell in the dungeon.
type Tile int

const (
	Wall Tile = iota
	Floor
)

// Room represents a rectangular room.
type Room struct {
	X, Y, W, H int
}

// Center returns the center cell of the room.
func (r Room) Center() (cx, cy int) {
	return r.X + r.W/2, r.Y + r.H/2
}

// Dungeon holds the map grid and metadata.
type Dungeon struct {
	Width  int
	Height int
	Tiles  [][]Tile
	Rooms  []Room
}

// At returns the tile at (x, y). Returns Wall for out-of-bounds.
func (d *Dungeon) At(x, y int) Tile {
	if x < 0 || x >= d.Width || y < 0 || y >= d.Height {
		return Wall
	}
	return d.Tiles[y][x]
}

// IsWalkable returns true if the tile can be walked on.
func (d *Dungeon) IsWalkable(x, y int) bool {
	return d.At(x, y) == Floor
}

// Config for procedural generation.
type Config struct {
	Width       int
	Height      int
	RoomCount   int
	MinRoomW    int
	MinRoomH    int
	MaxRoomW    int
	MaxRoomH    int
	CorridorW   int // 1 = single-tile corridors
	RandomSeed  int64
}

// DefaultConfig returns a sensible default.
func DefaultConfig() Config {
	return Config{
		Width:      80,
		Height:     60,
		RoomCount:  12,
		MinRoomW:   5,
		MinRoomH:   4,
		MaxRoomW:   12,
		MaxRoomH:   8,
		CorridorW:  1,
		RandomSeed: 0,
	}
}

// Generate creates a new dungeon with rooms and corridors.
func Generate(cfg Config) *Dungeon {
	rng := rand.New(rand.NewSource(cfg.RandomSeed))
	if cfg.RandomSeed == 0 {
		rng = rand.New(rand.NewSource(rand.Int63()))
	}
	d := &Dungeon{
		Width:  cfg.Width,
		Height: cfg.Height,
		Tiles:  make([][]Tile, cfg.Height),
	}
	for y := 0; y < cfg.Height; y++ {
		d.Tiles[y] = make([]Tile, cfg.Width)
		for x := 0; x < cfg.Width; x++ {
			d.Tiles[y][x] = Wall
		}
	}

	var rooms []Room
	margin := 2
	for i := 0; i < cfg.RoomCount*3; i++ { // try more than needed
		if len(rooms) >= cfg.RoomCount {
			break
		}
		w := cfg.MinRoomW + rng.Intn(cfg.MaxRoomW-cfg.MinRoomW+1)
		h := cfg.MinRoomH + rng.Intn(cfg.MaxRoomH-cfg.MinRoomH+1)
		x := margin + rng.Intn(cfg.Width-w-2*margin)
		y := margin + rng.Intn(cfg.Height-h-2*margin)
		if x < 0 || y < 0 {
			continue
		}
		r := Room{X: x, Y: y, W: w, H: h}
		overlap := false
		for _, other := range rooms {
			if roomsOverlap(r, other, 1) {
				overlap = true
				break
			}
		}
		if overlap {
			continue
		}
		carveRoom(d, r)
		rooms = append(rooms, r)
	}
	d.Rooms = rooms

	// Connect each room to the next with L-shaped corridors
	for i := 0; i < len(rooms)-1; i++ {
		cx1, cy1 := rooms[i].Center()
		cx2, cy2 := rooms[i+1].Center()
		carveCorridor(d, cx1, cy1, cx2, cy2)
	}

	return d
}

func roomsOverlap(a, b Room, padding int) bool {
	return a.X-padding < b.X+b.W &&
		a.X+a.W+padding > b.X &&
		a.Y-padding < b.Y+b.H &&
		a.Y+a.H+padding > b.Y
}

func carveRoom(d *Dungeon, r Room) {
	for y := r.Y; y < r.Y+r.H && y < d.Height; y++ {
		if y < 0 {
			continue
		}
		for x := r.X; x < r.X+r.W && x < d.Width; x++ {
			if x < 0 {
				continue
			}
			d.Tiles[y][x] = Floor
		}
	}
}

func carveCorridor(d *Dungeon, x1, y1, x2, y2 int) {
	x, y := x1, y1
	for x != x2 {
		if x < x2 {
			x++
		} else {
			x--
		}
		d.Tiles[y][x] = Floor
	}
	for y != y2 {
		if y < y2 {
			y++
		} else {
			y--
		}
		d.Tiles[y][x] = Floor
	}
}
