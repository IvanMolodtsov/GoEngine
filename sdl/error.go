package sdl

// #cgo LDFLAGS: -L. -lSDL3
// #include <SDL3/SDL.h>
import "C"
import "errors"

func GetError() error {
	return errors.New(C.GoString(C.SDL_GetError()))
}
