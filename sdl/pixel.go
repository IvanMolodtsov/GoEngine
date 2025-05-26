package sdl

// #cgo LDFLAGS: -lSDL3
// #include <SDL3/SDL.h>
import "C"

type PixelFormat = C.SDL_PixelFormat

const (
	PixelFormatUnknown      PixelFormat = 0
	PixelFormatIndex1LSB    PixelFormat = 0x11100100
	PixelFormatIndex1MSB    PixelFormat = 0x11200100
	PixelFormatIndex2LSB    PixelFormat = 0x1C100200
	PixelFormatIndex2MSB    PixelFormat = 0x1C200200
	PixelFormatIndex4LSB    PixelFormat = 0x12100400
	PixelFormatIndex4MSB    PixelFormat = 0x12200400
	PixelFormatIndex8       PixelFormat = 0x13000801
	PixelFormatRGB332       PixelFormat = 0x14110801
	PixelFormatXRGB4444     PixelFormat = 0x15120C02
	PixelFormatXBGR4444     PixelFormat = 0x15520C02
	PixelFormatXRGB1555     PixelFormat = 0x15130F02
	PixelFormatXBGR1555     PixelFormat = 0x15530F02
	PixelFormatARGB4444     PixelFormat = 0x15321002
	PixelFormatRGBA4444     PixelFormat = 0x15421002
	PixelFormatABGR4444     PixelFormat = 0x15721002
	PixelFormatBGRA4444     PixelFormat = 0x15821002
	PixelFormatARGB1555     PixelFormat = 0x15331002
	PixelFormatRGBA5551     PixelFormat = 0x15441002
	PixelFormatABGR1555     PixelFormat = 0x15731002
	PixelFormatBGRA5551     PixelFormat = 0x15841002
	PixelFormatRGB565       PixelFormat = 0x15151002
	PixelFormatBGR565       PixelFormat = 0x15551002
	PixelFormatRGB24        PixelFormat = 0x17101803
	PixelFormatBGR24        PixelFormat = 0x17401803
	PixelFormatXRGB8888     PixelFormat = 0x16161804
	PixelFormatRGBX8888     PixelFormat = 0x16261804
	PixelFormatXBGR8888     PixelFormat = 0x16561804
	PixelFormatBGRX8888     PixelFormat = 0x16661804
	PixelFormatARGB8888     PixelFormat = 0x16362004
	PixelFormatRGBA8888     PixelFormat = 0x16462004
	PixelFormatABGR8888     PixelFormat = 0x16762004
	PixelFormatBGRA8888     PixelFormat = 0x16862004
	PixelFormatXRGB2101010  PixelFormat = 0x16172004
	PixelFormatXBGR2101010  PixelFormat = 0x16572004
	PixelFormatARGB2101010  PixelFormat = 0x16372004
	PixelFormatABGR2101010  PixelFormat = 0x16772004
	PixelFormatRGB48        PixelFormat = 0x18103006
	PixelFormatBGR48        PixelFormat = 0x18403006
	PixelFormatRGBA64       PixelFormat = 0x18204008
	PixelFormatARGB64       PixelFormat = 0x18304008
	PixelFormatBGRA64       PixelFormat = 0x18504008
	PixelFormatABGR64       PixelFormat = 0x18604008
	PixelFormatRGB48Float   PixelFormat = 0x1A103006
	PixelFormatBGR48Float   PixelFormat = 0x1A403006
	PixelFormatRGBA64Float  PixelFormat = 0x1A204008
	PixelFormatARGB64Float  PixelFormat = 0x1A304008
	PixelFormatBGRA64Float  PixelFormat = 0x1A504008
	PixelFormatABGR64Float  PixelFormat = 0x1A604008
	PixelFormatRGB96Float   PixelFormat = 0x1B10600C
	PixelFormatBGR96Float   PixelFormat = 0x1B40600C
	PixelFormatRGBA128Float PixelFormat = 0x1B208010
	PixelFormatARGB128Float PixelFormat = 0x1B308010
	PixelFormatBGRA128Float PixelFormat = 0x1B508010
	PixelFormatABGR128Float PixelFormat = 0x1B608010
	PixelFormatYV12         PixelFormat = 0x32315659
	PixelFormatIYUV         PixelFormat = 0x56555949
	PixelFormatYUY2         PixelFormat = 0x32595559
	PixelFormatUYVY         PixelFormat = 0x59565955
	PixelFormatYVYU         PixelFormat = 0x55595659
	PixelFormatNV12         PixelFormat = 0x3231564E
	PixelFormatNV21         PixelFormat = 0x3132564E
	PixelFormatP010         PixelFormat = 0x30313050
	PixelFormatExternalOES  PixelFormat = 0x2053454F
	PixelFormatRGBA32       PixelFormat = PixelFormatABGR8888
	PixelFormatARGB32       PixelFormat = PixelFormatBGRA8888
	PixelFormatBGRA32       PixelFormat = PixelFormatARGB8888
	PixelFormatABGR32       PixelFormat = PixelFormatRGBA8888
	PixelFormatRGBX32       PixelFormat = PixelFormatXBGR8888
	PixelFormatXRGB32       PixelFormat = PixelFormatBGRX8888
	PixelFormatBGRX32       PixelFormat = PixelFormatXRGB8888
	PixelFormatXBGR32       PixelFormat = PixelFormatRGBX8888
)

type PixelFormatDetails struct {
	Format                         PixelFormat
	BitsPerPixel                   uint8
	BytesPerPixel                  uint8
	Padding                        [2]uint8
	Rmask, Gmask, Bmask, Amask     uint32
	Rbits, Gbits, Bbits, Abits     uint8
	Rshift, Gshift, Bshift, Ashift uint8
}

type Color struct {
	R, G, B, A uint8
}

type FColor struct {
	R, G, B, A float32
}

type Palette struct {
	ncolors int32
	colors  *Color
	// Version is for internal use only, do not touch.
	Version uint32
	// Refcount is for internal use only, do not touch.
	Refcount int32
}

func GetPixelFormatDetails(format PixelFormat) *PixelFormatDetails {
	var details *C.SDL_PixelFormatDetails = C.SDL_GetPixelFormatDetails(format)
	return CastPointer[PixelFormatDetails](details)
}

func MapRGB(details *PixelFormatDetails, pallette *Palette, r, g, b uint8) uint32 {
	color := C.SDL_MapRGB(CastPointer[C.SDL_PixelFormatDetails](details), CastPointer[C.SDL_Palette](pallette), C.Uint8(r), C.Uint8(g), C.Uint8(b))
	return uint32(color)
}

func MapRGBA(details *PixelFormatDetails, pallette *Palette, r, g, b, a uint8) uint32 {
	color := C.SDL_MapRGBA(CastPointer[C.SDL_PixelFormatDetails](details), CastPointer[C.SDL_Palette](pallette), C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a))
	return uint32(color)
}
