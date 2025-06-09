package sdl

// #cgo CFLAGS: -I${SRCDIR}/include/SDL
// #cgo LDFLAGS: -L${SRCDIR}/libs -lSDL3
// #include "SDL.h"
import "C"

type InitFlags = C.SDL_InitFlags

const (
	InitAudio    InitFlags = 0x00000010
	InitVideo    InitFlags = 0x00000020
	InitJoystick InitFlags = 0x00000200
	InitHaptic   InitFlags = 0x00001000
	InitGamepad  InitFlags = 0x00002000
	InitEvents   InitFlags = 0x00004000
	InitSensor   InitFlags = 0x00008000
	InitCamera   InitFlags = 0x00010000
)

func Init(flags InitFlags) error {
	if !C.SDL_Init(flags) {
		return GetError()
	}
	return nil
}

func Quit() {
	C.SDL_Quit()
}
