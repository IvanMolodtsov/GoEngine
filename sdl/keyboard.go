package sdl

// #cgo LDFLAGS: -L. -lSDL3
// #include <SDL3/SDL.h>
import "C"

type KeyboardID = C.SDL_KeyboardID
