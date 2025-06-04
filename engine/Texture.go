package engine

import (
	"unsafe"

	"github.com/IvanMolodtsov/GoEngine/sdl"
)

type Texture struct {
	texture     *sdl.Texture
	width       int32
	height      int32
	surface     *sdl.Surface
	Buffer      []uint32
	DepthBuffer [][]float64
	// mutex       sync.Mutex
}

func InitTexture(renderer *sdl.Renderer, width, height int32) *Texture {
	var texture Texture
	t, err := sdl.CreateTexture(renderer, sdl.PixelFormatRGBA32, sdl.TextureAccessStreaming, width, height)
	if err != nil {
		return nil
	}
	texture.texture = t
	texture.DepthBuffer = make([][]float64, width)
	for i := range texture.DepthBuffer {
		texture.DepthBuffer[i] = make([]float64, height)
	}
	texture.width = width
	texture.height = height
	return &texture
}

func (texture *Texture) Lock() {
	var err error
	texture.surface, err = sdl.LockTextureToSurface(texture.texture, nil, texture.surface)
	if err != nil {
		panic(err)
	}
	texture.Buffer = unsafe.Slice((*uint32)(texture.surface.Pixels), texture.surface.Pitch*int32(texture.texture.H))
}

func (texture *Texture) Unlock(renderer *sdl.Renderer) {
	sdl.UnlockTexture(texture.texture)
	sdl.RenderTexture(renderer, texture.texture, nil, nil)
}

func (texture *Texture) ClearBuffer() {
	// texture.mutex.Lock()
	texture.DepthBuffer = make([][]float64, texture.width)
	for i := range texture.DepthBuffer {
		texture.DepthBuffer[i] = make([]float64, texture.height)
	}
	// texture.mutex.Unlock()
}

func (texture *Texture) Clear() {
	sdl.FillSurfaceRect(texture.surface, nil, 0)

}

func (texture *Texture) Destroy() {
	sdl.DestroyTexture(texture.texture)
}

func (texture *Texture) DrawPixel(x, y int64, w float64, color uint32) {
	if texture.DepthBuffer[x][y] < w {
		texture.DepthBuffer[x][y] = w
		texture.Buffer[y*int64(texture.width)+x] = color
	}
}
