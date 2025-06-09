package primitives

import (
	"image/color"

	"github.com/IvanMolodtsov/GoEngine/sdl"
	// "github.com/jupiterrider/purego-sdl3/sdl"
)

func ToHex(color color.Color, l float64) uint32 {
	R, G, B, A := color.RGBA()
	A = uint32(255 * l)
	return sdl.MapRGBA(sdl.GetPixelFormatDetails(sdl.PixelFormatRGBA32), nil, uint8(R), uint8(G), uint8(B), uint8(A))
}
