package main

import (
	"embed"
	"log"

	"github.com/ephigenia/ebit-engine-game-1/game"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/Gorgeous-Pixel/GorgeousPixel.ttf
var assets embed.FS

func main() {
	scale := 2
	ebiten.SetWindowSize(game.ScreenW*scale, game.ScreenH*scale)
	ebiten.SetWindowTitle("Dungeon Crawler")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	g := game.New(assets)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
