# Dungeon Crawler

A 2D dungeon crawler written in Go using [Ebitengine](https://ebitengine.org/). Levels are **procedurally generated** with a rooms-and-corridors algorithm: random rectangular rooms are placed and connected by L-shaped corridors.

## Run

```bash
go run .
# or
go build -o dungeon-crawler . && ./dungeon-crawler
```

## Controls

- **WASD** or **Arrow keys** – move (one tile per keypress)
- **R** – generate a new dungeon

## Structure

- `main.go` – entry point, window setup
- `game/` – game loop, player, camera, rendering
- `dungeon/` – procedural dungeon generation (rooms + corridors)

The camera follows the player; the dungeon is larger than the screen so you explore by moving.


## Ideas

- items require a certain level or stats to be able to equipped
- list of slots an item can be equipped on should be an array so that items can be equipped on left or right hand 
- quick belt
  - additional item
- pickup items with key?
- drop items with key?
- diablo like
  - attributes influence weapons
  - different attack modes
- add image reader with offsets
- images with animations
- load assets in main instead of in game.go