package sdl

// #cgo CFLAGS: -I${SRCDIR}/include/SDL
// #cgo LDFLAGS: -L${SRCDIR}/libs -lSDL3
// #include "SDL.h"
import "C"
import "unsafe"

type SurfaceFlags = C.SDL_SurfaceFlags

const (
	SurfacePreallocated SurfaceFlags = 1 << iota
	SurfaceLockNeeded
	SurfaceLocked
	SurfaceSIMDAligned
)

type Surface struct {
	Flags    SurfaceFlags
	Format   PixelFormat
	W        int32
	H        int32
	Pitch    int32
	Pixels   unsafe.Pointer
	Refcount int32
	Reserved unsafe.Pointer
}

type SDLSurface = C.SDL_Surface

func (s *Surface) cPtr() *SDLSurface {
	return (*SDLSurface)(unsafe.Pointer(s))
}

func FillSurfaceRect(surface *Surface, rect *Rect, color uint32) error {
	if !C.SDL_FillSurfaceRect(surface.cPtr(), rect.cPtr(), C.uint(color)) {
		return GetError()
	}
	return nil
}
