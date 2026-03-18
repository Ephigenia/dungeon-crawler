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
