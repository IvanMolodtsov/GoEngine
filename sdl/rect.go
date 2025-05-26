package sdl

// #cgo LDFLAGS: -lSDL3
// #include <SDL3/SDL.h>
import "C"
import "unsafe"

type Rect struct {
	X, Y, W, H int32
}

func (s *Rect) cPtr() *C.SDL_Rect {
	return (*C.SDL_Rect)(unsafe.Pointer(s))
}

// FRect is a rectangle, with the origin at the upper left (using floating point values).
type FRect struct {
	X, Y, W, H float32
}

func (s *FRect) cPtr() *C.SDL_FRect {
	return (*C.SDL_FRect)(unsafe.Pointer(s))
}

type FPoint struct {
	X, Y float32
}

type Point struct {
	X, Y int32
}
