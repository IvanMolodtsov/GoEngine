package sdl

// #cgo LDFLAGS: -L. -lSDL3
// #include <SDL3/SDL.h>
import "C"

type Renderer = C.SDL_Renderer

func CreateRenderer(window *Window, name string) (*Renderer, error) {
	cname := C.CString(name)
	defer Free(cname)
	var r *Renderer = C.SDL_CreateRenderer(window, cname)
	if r == nil {
		return nil, GetError()
	}
	return r, nil
}

func DestroyRenderer(renderer *Renderer) {
	C.SDL_DestroyRenderer(renderer)
}

func RenderClear(renderer *Renderer) error {
	if !C.SDL_RenderClear(renderer) {
		return GetError()
	}
	return nil
}

func RenderPresent(renderer *Renderer) error {
	if !C.SDL_RenderPresent(renderer) {
		return GetError()
	}
	return nil
}

func RenderTexture(renderer *Renderer, texture *Texture, r1, r2 *FRect) {
	C.SDL_RenderTexture(renderer, texture.cPtr(), r1.cPtr(), r2.cPtr())
}
