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

// allObjectTypes returns AllObjectTypes plus any types not in the spawnable pool.
func allObjectTypes() []*ObjectType {
	return append(AllObjectTypes, ObjectTypeShelf)
}

// loadObjectImages loads standalone sprites for all object types from the embedded FS.
// Types that use the shared spritesheet (UsesSpritesheet == true) are skipped.
func loadObjectImages(assets fs.FS) {
	for _, ot := range allObjectTypes() {
		if ot.UsesSpritesheet {
			continue
		}
		if ot.SpritesheetPath != "" {
			sheet := loadImageFile(assets, ot.SpritesheetPath)
			if sheet != nil {
				cols := sheet.Bounds().Dx() / TileSize
				spriteRect := func(idx int) image.Rectangle {
					col := idx % cols
					row := idx / cols
					x := col * TileSize
					y := row * TileSize
					return image.Rect(x, y, x+TileSize, y+TileSize)
				}
				ot.Image = sheet.SubImage(spriteRect(ot.SpritesheetIndex)).(*ebiten.Image)
				if ot.HasDestroyedSprite {
					ot.DestroyedImage = sheet.SubImage(spriteRect(ot.DestroyedSpritesheetIndex)).(*ebiten.Image)
				}
			}
			continue
		}
		ot.Image = loadImageFile(assets, ot.ImagePath)
	}
}
