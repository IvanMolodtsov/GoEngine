package sdl

// #cgo LDFLAGS: -L. -lSDL3
// #include <SDL3/SDL.h>
import "C"
import "unsafe"

type Texture struct {
	Format   PixelFormat
	W        int32
	H        int32
	Refcount int32
}

func (t *Texture) cPtr() *C.SDL_Texture {
	return (*C.SDL_Texture)(unsafe.Pointer(t))
}

type TextureAccess = C.SDL_TextureAccess

const (
	TextureAccessStatic TextureAccess = iota
	TextureAccessStreaming
	TextureAccessTarget
)

func CreateTexture(renderer *Renderer, format PixelFormat, access TextureAccess, w int32, h int32) (*Texture, error) {
	t := C.SDL_CreateTexture(renderer, format, access, C.int(w), C.int(h))
	if t == nil {
		return nil, GetError()
	}
	return (*Texture)(unsafe.Pointer(t)), nil
}

func DestroyTexture(texture *Texture) {
	C.SDL_DestroyTexture(texture.cPtr())
}

func UnlockTexture(texture *Texture) {
	C.SDL_UnlockTexture(texture.cPtr())
}

func LockTexture(texture *Texture, data *[]uint32, pitch *int) error {
	dataPtr := unsafe.Pointer(data)
	var temp C.int
	if !C.SDL_LockTexture(texture.cPtr(), nil, &dataPtr, &temp) {
		return GetError()
	}
	*pitch = int(temp)
	return nil
}

func LockTextureToSurface(texture *Texture, rect *Rect, surface *Surface) (*Surface, error) {

	var temp = surface.cPtr()

	if !C.SDL_LockTextureToSurface(texture.cPtr(), rect.cPtr(), &temp) {
		return nil, GetError()
	}
	return CastPointer[Surface](temp), nil
}
