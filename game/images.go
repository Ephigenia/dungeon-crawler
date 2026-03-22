package game

import (
	"image"
	_ "image/png"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
)

// loadImageFile opens, decodes, and returns an ebiten.Image from the embedded FS.
// Returns nil if the path is empty, the file is missing, or decoding fails.
func loadImageFile(assets fs.FS, path string) *ebiten.Image {
	if path == "" {
		return nil
	}
	f, err := assets.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil
	}
	return ebiten.NewImageFromImage(img)
}

// loadItemImages loads sprites for all items from the embedded FS.
func loadItemImages(assets fs.FS) {
	for _, item := range AllItems {
		item.Image = loadImageFile(assets, item.ImagePath)
	}
}

// loadEnemyImages loads sprites for all enemy types from the embedded FS.
func loadEnemyImages(assets fs.FS) {
	for _, et := range AllEnemyTypes {
		et.Image = loadImageFile(assets, et.ImagePath)
	}
}

// loadObjectImages loads standalone sprites for all object types from the embedded FS.
// Types that use the shared spritesheet (UsesSpritesheet == true) are skipped.
func loadObjectImages(assets fs.FS) {
	for _, ot := range AllObjectTypes {
		if !ot.UsesSpritesheet {
			ot.Image = loadImageFile(assets, ot.ImagePath)
		}
	}
}
