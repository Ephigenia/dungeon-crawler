package game

import (
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
)

// Spritesheet holds a loaded image and sprite dimensions.
// Sprites are indexed left-to-right, top-to-bottom starting at 0.
// A 32×32 image with 16×16 sprites contains sprites 0–3:
//
//	0 1
//	2 3
type Spritesheet struct {
	img     *ebiten.Image
	spriteW int
	spriteH int
	cols    int
}

// LoadSpritesheet loads an image from assets and returns a Spritesheet
// that slices it into sprites of spriteW×spriteH pixels.
// Returns nil if the image cannot be loaded.
func LoadSpritesheet(assets fs.FS, path string, spriteW, spriteH int) *Spritesheet {
	img := loadImageFile(assets, path)
	if img == nil {
		return nil
	}
	return &Spritesheet{
		img:     img,
		spriteW: spriteW,
		spriteH: spriteH,
		cols:    img.Bounds().Dx() / spriteW,
	}
}

// Sprite returns a new image containing the sprite at the given zero-based index.
// Returns nil if the index is out of range or the sheet is nil.
func (s *Spritesheet) Sprite(index int) *ebiten.Image {
	if s == nil {
		return nil
	}
	total := s.cols * (s.img.Bounds().Dy() / s.spriteH)
	if index < 0 || index >= total {
		return nil
	}
	col := index % s.cols
	row := index / s.cols
	x := col * s.spriteW
	y := row * s.spriteH

	dst := ebiten.NewImage(s.spriteW, s.spriteH)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(-x), float64(-y))
	dst.DrawImage(s.img, op)
	return dst
}
