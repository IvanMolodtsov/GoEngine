package engine

import (
	"image/color"

	"github.com/jupiterrider/purego-sdl3/sdl"
)

func toHex(c color.RGBA) uint32 {
	alpha := float32(c.A) / float32(255)
	return sdl.MapRGB(sdl.GetPixelFormatDetails(sdl.PixelFormatRGBA32), nil, uint8(float32(c.R)*alpha), uint8(float32(c.G)*alpha), uint8(float32(c.B)*alpha))
}
