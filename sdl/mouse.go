package sdl

// #cgo LDFLAGS: -lSDL3
// #include <SDL3/SDL.h>
import "C"

type MouseButtonFlags = C.SDL_MouseButtonFlags

type MouseID = C.SDL_MouseID

func CaptureMouse(window *Window) bool {
	return bool(C.SDL_SetWindowMouseGrab(window, C.bool(true)))
}

func SetWindowRelativeMouseMode(window *Window) bool {
	return bool(C.SDL_SetWindowRelativeMouseMode(window, C.bool(true)))
}
