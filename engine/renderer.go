package engine

import (
	"image/color"
	"math"
	"sync"
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
	screenClipPlane  Plane
	planes           []Plane
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

	renderer.screenClipPlane = NewPlane(NewVector3d(0, 0, 0.1), NewVector3d(0, 0, 1.0))
	top := NewPlane(NewVector3d(0, 0, 0), NewVector3d(0, 1, 0))
	bottom := NewPlane(NewVector3d(0, float64(renderer.screenHeight-1), 0), NewVector3d(0, -1, 0))
	left := NewPlane(NewVector3d(0, 0, 0), NewVector3d(1, 0, 0))
	right := NewPlane(NewVector3d(float64(renderer.screenWidth-1), 0, 0), NewVector3d(-1, 0, 0))

	renderer.planes = []Plane{
		top, bottom, left, right,
	}
	return &renderer, nil
}

func (renderer *Renderer) destroy() {
	sdl.DestroyRenderer(renderer.renderer)
	sdl.DestroyTexture(renderer.texture)
}

func (renderer *Renderer) clearScreenBuffer() {
	sdl.FillSurfaceRect(renderer.bufferSurface, nil, 0)
}

func (renderer *Renderer) Render(tris []Triangle) {
	sdl.RenderClear(renderer.renderer)

	sdl.LockTextureToSurface(renderer.texture, nil, &renderer.bufferSurface)

	renderer.clearScreenBuffer()

	renderer.buffer = unsafe.Slice((*uint32)(renderer.bufferSurface.Pixels), renderer.bufferSurface.Pitch*int32(renderer.screenHeight))

	trisToRaster := NewQueue[Triangle]()

	for _, t := range tris {
		trisToRaster.Start()
		go func() {
			listTriangles := make([]Triangle, 0)
			listTriangles = append(listTriangles, t)
			var newT = 1
			for _, p := range renderer.planes {
				for newT > 0 {
					test := listTriangles[0]
					listTriangles = listTriangles[1:]
					newT -= 1

					clipped := p.Clip(test)
					listTriangles = append(listTriangles, clipped...)

				}
				newT = len(listTriangles)
			}

			for _, res := range listTriangles {
				trisToRaster.Push(res)
			}
			trisToRaster.End()
		}()

	}

	for t := range trisToRaster.Values {
		renderer.FillTriangle(t)
		t.Render(renderer)
	}

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

func (renderer *Renderer) Project(o Object, camera *Camera, trisToRender *Queue[Triangle]) {
	var wg sync.WaitGroup
	worldMatrix := o.GetWorld()
	view := camera.GetView()
	screenClipPlane := renderer.screenClipPlane
	for _, t := range o.Mesh.Tris {
		wg.Add(1)
		go func() {
			var transformed, viewed, projected Triangle

			// World Matrix Transform
			transformed.Points[0] = worldMatrix.MulV(t.Points[0])
			transformed.Points[1] = worldMatrix.MulV(t.Points[1])
			transformed.Points[2] = worldMatrix.MulV(t.Points[2])

			// Get surface normal
			normal := transformed.Normal()
			cameraRay := transformed.Points[0].Sub(camera.Position)

			// check if triangle if visible
			if normal.DotProduct(cameraRay) < 0.0 {
				lightDir := NewVector3d(0, 0, -1.0)
				lightDir = lightDir.Normalize()

				luminance := normal.DotProduct(lightDir)

				viewed.Points[0] = view.MulV(transformed.Points[0])
				viewed.Points[1] = view.MulV(transformed.Points[1])
				viewed.Points[2] = view.MulV(transformed.Points[2])

				//Clip Triangle
				clipped := screenClipPlane.Clip(viewed)

				for _, c := range clipped {
					// Project triangles from 3D --> 2D
					projected.Points[0] = renderer.ProjectionMatrix.MulV(c.Points[0])
					projected.Points[1] = renderer.ProjectionMatrix.MulV(c.Points[1])
					projected.Points[2] = renderer.ProjectionMatrix.MulV(c.Points[2])
					projected.Points[0] = projected.Points[0].Div(projected.Points[0].W)
					projected.Points[1] = projected.Points[1].Div(projected.Points[1].W)
					projected.Points[2] = projected.Points[2].Div(projected.Points[2].W)

					offset := NewVector3d(1.0, 1.0, 0)
					projected.Points[0] = projected.Points[0].Add(offset)
					projected.Points[1] = projected.Points[1].Add(offset)
					projected.Points[2] = projected.Points[2].Add(offset)

					projected.Points[0].X *= 0.5 * float64(renderer.screenWidth)
					projected.Points[0].Y *= 0.5 * float64(renderer.screenHeight)
					projected.Points[1].X *= 0.5 * float64(renderer.screenWidth)
					projected.Points[1].Y *= 0.5 * float64(renderer.screenHeight)
					projected.Points[2].X *= 0.5 * float64(renderer.screenWidth)
					projected.Points[2].Y *= 0.5 * float64(renderer.screenHeight)

					projected.Color = color.RGBA{R: 255, G: 255, B: 255, A: uint8(luminance * 255)}

					trisToRender.Push(projected)

				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	trisToRender.End()
}

func (r *Renderer) FillTriangle(t Triangle) {
	x1 := int64(t.Points[0].X)
	y1 := int64(t.Points[0].Y)
	x2 := int64(t.Points[1].X)
	y2 := int64(t.Points[1].Y)
	x3 := int64(t.Points[2].X)
	y3 := int64(t.Points[2].Y)

	SWAP := func(x, y *int64) {
		t := *x
		*x = *y
		*y = t
	}

	DRAWLINE := func(sx, ex, ny int64) {
		for i := sx; i <= ex; i++ {
			r.DrawPixel(float64(i), float64(ny), toHex(t.Color))
		}
	}

	changed1, changed2 := false, false

	// Sort vertices
	if y1 > y2 {
		SWAP(&y1, &y2)
		SWAP(&x1, &x2)
	}
	if y1 > y3 {
		SWAP(&y1, &y3)
		SWAP(&x1, &x3)
	}
	if y2 > y3 {
		SWAP(&y2, &y3)
		SWAP(&x2, &x3)
	}

	// Starting points
	t1x, t2x := x1, x1
	y := y1

	var signx1 int64
	var dx1 int64 = x2 - x1
	if dx1 < 0 {
		dx1 = -dx1
		signx1 = -1
	} else {
		signx1 = 1
	}

	dy1 := y2 - y1

	var signx2 int64
	dx2 := x3 - x1
	if dx2 < 0 {
		dx2 = -dx2
		signx2 = -1
	} else {
		signx2 = 1
	}

	dy2 := y3 - y1

	if dy1 > dx1 { // swap values
		SWAP(&dx1, &dy1)
		changed1 = true
	}
	if dy2 > dx2 { // swap values
		SWAP(&dy2, &dx2)
		changed2 = true
	}
	var minx, maxx, t1xp, t2xp, e1, i int64

	e2 := dx2 >> 1
	// Flat top, just process the second half
	if y1 == y2 {
		goto NEXT
	}
	e1 = dx1 >> 1

	for i = 0; i < dx1; {
		t1xp = 0
		t2xp = 0
		if t1x < t2x {
			minx = t1x
			maxx = t2x
		} else {
			minx = t2x
			maxx = t1x
		}
		// process first line until y value is about to change
		for i < dx1 {
			i++
			e1 += dy1
			for e1 >= dx1 {
				e1 -= dx1
				if changed1 { //t1x += signx1;
					t1xp = signx1
				} else {
					goto next1
				}
			}
			if changed1 {
				break
			} else {
				t1x += signx1
			}
		}
		// Move line
	next1:
		// process second line until y value is about to change
		for {
			e2 += dy2
			for e2 >= dx2 {
				e2 -= dx2
				if changed2 {
					//t2x += signx2;
					t2xp = signx2
				} else {
					goto next2
				}
			}
			if changed2 {
				break
			} else {
				t2x += signx2
			}
		}
	next2:
		if minx > t1x {
			minx = t1x
			if minx > t2x {
				minx = t2x
			}
		}
		if maxx < t1x {
			maxx = t1x
			if maxx < t2x {
				maxx = t2x
			}
		}
		DRAWLINE(minx, maxx, y) // Draw line from min to max points found on the y
		// Now increase y
		if !changed1 {
			t1x += signx1
		}
		t1x += t1xp
		if !changed2 {
			t2x += signx2
		}
		t2x += t2xp
		y += 1
		if y == y2 {
			break
		}

	}
NEXT:
	// Second half
	dx1 = x3 - x2
	if dx1 < 0 {
		dx1 = -dx1
		signx1 = -1
	} else {
		signx1 = 1
	}
	dy1 = (y3 - y2)
	t1x = x2

	if dy1 > dx1 { // swap values
		SWAP(&dy1, &dx1)
		changed1 = true
	} else {
		changed1 = false
	}

	e1 = dx1 >> 1

	for i = 0; i <= dx1; i++ {
		t1xp = 0
		t2xp = 0
		if t1x < t2x {
			minx = t1x
			maxx = t2x
		} else {
			minx = t2x
			maxx = t1x
		}
		// process first line until y value is about to change
		for i < dx1 {
			e1 += dy1
			for e1 >= dx1 {
				e1 -= dx1
				if changed1 {
					t1xp = signx1
					break
					//t1x += signx1;
				} else {
					goto next3
				}
			}
			if changed1 {
				break
			} else {
				t1x += signx1
			}
			if i < dx1 {
				i++
			}
		}
	next3:
		// process second line until y value is about to change
		for t2x != x3 {
			e2 += dy2
			for e2 >= dx2 {
				e2 -= dx2
				if changed2 {
					t2xp = signx2
				} else {
					goto next4
				}
			}
			if changed2 {
				break
			} else {
				t2x += signx2
			}
		}
	next4:

		if minx > t1x {
			minx = t1x
			if minx > t2x {
				minx = t2x
			}
		}
		if maxx < t1x {
			maxx = t1x
			if maxx < t2x {
				maxx = t2x
			}
		}
		DRAWLINE(minx, maxx, y)
		if !changed1 {
			t1x += signx1
		}
		t1x += t1xp
		if !changed2 {
			t2x += signx2
		}
		t2x += t2xp
		y += 1
		if y > y3 {
			return
		}
	}
}
