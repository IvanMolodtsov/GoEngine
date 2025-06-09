package sdl

// #cgo CFLAGS: -I${SRCDIR}/include/SDL
// #cgo LDFLAGS: -L${SRCDIR}/../ -lSDL3
// #include "SDL.h"
import "C"
import "errors"

func GetError() error {
	return errors.New(C.GoString(C.SDL_GetError()))
}
