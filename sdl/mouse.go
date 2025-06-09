package sdl

// #cgo CFLAGS: -I${SRCDIR}/include/SDL
// #cgo LDFLAGS: -L${SRCDIR}/libs -lSDL3
// #include "SDL.h"
import "C"

type MouseButtonFlags = C.SDL_MouseButtonFlags

type MouseID = C.SDL_MouseID

func CaptureMouse(window *Window) bool {
	return bool(C.SDL_SetWindowMouseGrab(window, C.bool(true)))
}

func SetWindowRelativeMouseMode(window *Window) bool {
	return bool(C.SDL_SetWindowRelativeMouseMode(window, C.bool(true)))
}
