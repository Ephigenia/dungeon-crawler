package game

import (
	"image"
	_ "image/png"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
)

// loadItemImages loads all item sprites from the embedded FS.
// Items whose ImagePath is empty or whose file fails to load keep their Color fallback.
func loadItemImages(assets fs.FS) {
	for _, item := range AllItems {
		if item.ImagePath == "" {
			continue
		}
		f, err := assets.Open(item.ImagePath)
		if err != nil {
			continue
		}
		img, _, err := image.Decode(f)
		f.Close()
		if err != nil {
			continue
		}
		item.Image = ebiten.NewImageFromImage(img)
	}
}

// loadEnemyImages loads all enemy sprites from the embedded FS.
// Enemy types whose ImagePath is empty or whose file fails to load keep the color fallback.
func loadEnemyImages(assets fs.FS) {
	for _, et := range AllEnemyTypes {
		if et.ImagePath == "" {
			continue
		}
		f, err := assets.Open(et.ImagePath)
		if err != nil {
			continue
		}
		img, _, err := image.Decode(f)
		f.Close()
		if err != nil {
			continue
		}
		et.Image = ebiten.NewImageFromImage(img)
	}
}
