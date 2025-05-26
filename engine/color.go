package engine

import (
	"image/color"

	"github.com/IvanMolodtsov/GoEngine/sdl"
	// "github.com/jupiterrider/purego-sdl3/sdl"
)

func toHex(c color.RGBA) uint32 {
	// alpha := float32(c.A) / float32(255)
	return sdl.MapRGBA(sdl.GetPixelFormatDetails(sdl.PixelFormatRGBA32), nil, c.R, c.G, c.B, c.A)
}
