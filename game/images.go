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

// loadEnemyImages loads spritesheets for all enemy types.
// Frame size is auto-detected from each image's height (square frames).
func loadEnemyImages(assets fs.FS) {
	for _, et := range AllEnemyTypes {
		et.IdleSheet = LoadSpritesheetAutoFrame(assets, et.IdleImagePath)
		et.MoveSheet = LoadSpritesheetAutoFrame(assets, et.MoveImagePath)
		et.AttackSheet = LoadSpritesheetAutoFrame(assets, et.AttackImagePath)
		et.DeathSheet = LoadSpritesheetAutoFrame(assets, et.DeathImagePath)
	}
}

// loadObjectImages loads standalone sprites for all object types from the embedded FS.
// Types that use the shared spritesheet (UsesSpritesheet == true) are skipped.
func loadObjectImages(assets fs.FS) {
	for _, ot := range AllObjectTypes {
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
			}
			continue
		}
		ot.Image = loadImageFile(assets, ot.ImagePath)
		ot.DestroyedImage = loadImageFile(assets, ot.DestroyedImagePath)
	}
}
