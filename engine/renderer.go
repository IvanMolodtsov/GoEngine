package engine

import (
	"math"
	"unsafe"

	"github.com/jupiterrider/purego-sdl3/sdl"
)

const (
	fNear float64 = 0.1
	fFar  float64 = 1000.0
	FOV   float64 = 90.0
)

type Renderer struct {
	renderer         *sdl.Renderer
	texture          *sdl.Texture
	bufferSurface    *sdl.Surface
	buffer           []uint32
	screenWidth      int64
	screenHeight     int64
	ProjectionMatrix Matrix4x4
}

func InitRenderer(window *sdl.Window, width int64, height int64) (*Renderer, error) {
	var renderer Renderer
	r := sdl.CreateRenderer(window, "")
	if r == nil {
		return nil, GetError()
	}
	renderer.renderer = r

	renderer.screenWidth = width
	renderer.screenHeight = height

	t := sdl.CreateTexture(renderer.renderer, sdl.PixelFormatRGBA32, sdl.TextureAccessStreaming, int32(renderer.screenWidth), int32(renderer.screenHeight))
	if t == nil {
		return nil, GetError()
	}
	renderer.texture = t
	renderer.ProjectionMatrix = ProjectionMatrix(float64(renderer.screenHeight)/float64(renderer.screenWidth), fNear, fFar, FOV)

	renderer.ProjectionMatrix.Print()
	return &renderer, nil
}

func (renderer *Renderer) destroy() {
	sdl.DestroyRenderer(renderer.renderer)
	sdl.DestroyTexture(renderer.texture)
}

func (renderer *Renderer) clearScreenBuffer() {
	sdl.FillSurfaceRect(renderer.bufferSurface, nil, 0)
}

func (renderer *Renderer) Render(renderF func()) {
	sdl.RenderClear(renderer.renderer)

	sdl.LockTextureToSurface(renderer.texture, nil, &renderer.bufferSurface)

	renderer.clearScreenBuffer()

	renderer.buffer = unsafe.Slice((*uint32)(renderer.bufferSurface.Pixels), renderer.bufferSurface.Pitch*int32(renderer.screenHeight))

	renderF()

	sdl.UnlockTexture(renderer.texture)
	sdl.RenderTexture(renderer.renderer, renderer.texture, nil, nil)
	sdl.RenderPresent(renderer.renderer)
}

func (renderer *Renderer) DrawPixel(x, y float64, color uint32) {
	dx := int64(x)
	dy := int64(y)

	isOutOfBounds := dx < 0 || dx > renderer.screenWidth || dy < 0 || dy > renderer.screenHeight
	isOutsideMemBuff := (renderer.screenWidth*dy + dx) >= renderer.screenWidth*renderer.screenHeight
	if isOutOfBounds || isOutsideMemBuff {
		return
	}

	renderer.buffer[dy*renderer.screenWidth+dx] = color
}

func (renderer *Renderer) DrawLine(x1, y1, x2, y2 float64, color uint32) {
	var x, y, dx, dy, dx1, dy1, px, py, xe, ye float64
	dx = x2 - x1
	dy = y2 - y1
	dx1 = float64(math.Abs(float64(dx)))
	dy1 = float64(math.Abs(float64(dy)))
	px = 2*dy1 - dx1
	py = 2*dx1 - dy1
	if dy1 <= dx1 {
		if dx >= 0 {
			x = x1
			y = y1
			xe = x2
		} else {
			x = x2
			y = y2
			xe = x1
		}
		renderer.DrawPixel(x, y, color)
		for i := 0; x < xe; i++ {
			x = x + 1
			if px < 0 {
				px = px + 2*dy1

			} else {
				if (dx < 0 && dy < 0) || (dx > 0 && dy > 0) {
					y = y + 1
				} else {
					y = y - 1
				}
				px = px + 2*(dy1-dx1)
			}
			renderer.DrawPixel(x, y, color)
		}
	} else {
		if dy >= 0 {
			x = x1
			y = y1
			ye = y2
		} else {
			x = x2
			y = y2
			ye = y1
		}
		renderer.DrawPixel(x, y, color)
		for i := 0; y < ye; i++ {
			y = y + 1
			if py <= 0 {
				py = py + 2*dx1
			} else {
				if (dx < 0 && dy < 0) || (dx > 0 && dy > 0) {
					x = x + 1

				} else {
					x = x - 1
				}
				py = py + 2*(dx1-dy1)
			}
			renderer.DrawPixel(x, y, color)
		}
	}
}
