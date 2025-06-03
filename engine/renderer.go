package engine

import (
	"math"
	"sync"

	"github.com/IvanMolodtsov/GoEngine/primitives"
	"github.com/IvanMolodtsov/GoEngine/sdl"
)

const (
	fNear float64 = 0.1
	fFar  float64 = 1000.0
	FOV   float64 = 90.0
)

type Renderer struct {
	Window           *sdl.Window
	renderer         *sdl.Renderer
	textures         [2]*Texture
	current          uint8
	screenWidth      int64
	screenHeight     int64
	ProjectionMatrix *primitives.Matrix4x4
	screenClipPlane  *primitives.Plane
	planes           []*primitives.Plane
	Texture          *primitives.Image
}

func InitRenderer(width int64, height int64) (*Renderer, error) {
	var renderer Renderer
	w, err := sdl.CreateWindow("Game", int32(width), int32(height), 0)
	if err != nil {
		return nil, err
	}
	renderer.Window = w
	r, err := sdl.CreateRenderer(renderer.Window, "")
	if err != nil {
		return nil, err
	}
	renderer.renderer = r

	renderer.screenWidth = width
	renderer.screenHeight = height

	renderer.textures[0] = InitTexture(renderer.renderer, int32(renderer.screenWidth), int32(renderer.screenHeight))
	renderer.textures[1] = InitTexture(renderer.renderer, int32(renderer.screenWidth), int32(renderer.screenHeight))
	renderer.current = 0
	renderer.textures[renderer.current].Lock()
	renderer.ProjectionMatrix = primitives.ProjectionMatrix(float64(renderer.screenHeight)/float64(renderer.screenWidth), fNear, fFar, FOV)

	renderer.screenClipPlane = primitives.NewPlane(primitives.NewVector3d(0, 0, 0.1), primitives.NewVector3d(0, 0, 1.0))
	top := primitives.NewPlane(primitives.NewVector3d(0, 0, 0), primitives.NewVector3d(0, 1, 0))
	bottom := primitives.NewPlane(primitives.NewVector3d(0, float64(renderer.screenHeight-1), 0), primitives.NewVector3d(0, -1, 0))
	left := primitives.NewPlane(primitives.NewVector3d(0, 0, 0), primitives.NewVector3d(1, 0, 0))
	right := primitives.NewPlane(primitives.NewVector3d(float64(renderer.screenWidth-1), 0, 0), primitives.NewVector3d(-1, 0, 0))

	renderer.planes = []*primitives.Plane{
		top, bottom, left, right,
	}
	return &renderer, nil
}

func (renderer *Renderer) Destroy() {
	sdl.DestroyRenderer(renderer.renderer)
	renderer.textures[0].Destroy()
	renderer.textures[1].Destroy()
	sdl.DestroyWindow(renderer.Window)
}

func (renderer *Renderer) ClipTriangles(tris []*primitives.Triangle) []*primitives.Triangle {
	trisToRaster := make([]*primitives.Triangle, 0)

	for _, t := range tris {
		listTriangles := make([]*primitives.Triangle, 0)
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
		trisToRaster = append(trisToRaster, listTriangles...)

	}
	return trisToRaster
}

func (renderer *Renderer) PushTriangles(tris []*primitives.Triangle, texture *primitives.Image) {
	for _, t := range tris {
		renderer.RasterizeTriangle(t, texture)
		// t.Render(renderer)
	}
}

func (renderer *Renderer) Render() {
	sdl.RenderClear(renderer.renderer)
	renderer.SwapTextures()
	sdl.RenderPresent(renderer.renderer)
}

func (renderer *Renderer) SwapTextures() {
	renderer.textures[renderer.current].Unlock(renderer.renderer)
	renderer.current = (renderer.current + 1) % 2
	renderer.textures[renderer.current].Lock()
}

func (renderer *Renderer) DrawPixel(x, y, w float64, color uint32) {
	dx := int64(x)
	dy := int64(y)

	isOutOfBounds := dx < 0 || dx >= renderer.screenWidth || dy < 0 || dy >= renderer.screenHeight
	isOutsideMemBuff := (renderer.screenWidth*dy + dx) >= renderer.screenWidth*renderer.screenHeight
	if renderer.textures[renderer.current].Buffer == nil || isOutOfBounds || isOutsideMemBuff {
		return
	}

	renderer.textures[renderer.current].DrawPixel(dx, dy, w, color)
}

func (renderer *Renderer) Project(o *primitives.Object, camera *Camera, trisToRender *Queue[*primitives.Triangle]) {
	var wg sync.WaitGroup
	worldMatrix := o.GetWorld()
	view := camera.GetView()
	screenClipPlane := renderer.screenClipPlane
	for _, t := range o.Mesh.Tris {
		wg.Add(1)
		go func() {
			var transformed, viewed primitives.Triangle

			// World Matrix Transform
			transformed.P[0] = worldMatrix.MulV(t.P[0])
			transformed.P[1] = worldMatrix.MulV(t.P[1])
			transformed.P[2] = worldMatrix.MulV(t.P[2])
			transformed.T[0] = t.T[0]
			transformed.T[1] = t.T[1]
			transformed.T[2] = t.T[2]

			// Get surface normal
			normal := transformed.Normal()
			cameraRay := transformed.P[0].Sub(camera.GetPosition())

			// check if triangle if visible
			if normal.DotProduct(cameraRay) < 0.0 {
				lightDir := primitives.NewVector3d(0, 0, -1.0)
				lightDir = lightDir.Normalize()

				// luminance := normal.DotProduct(lightDir)

				viewed.P[0] = view.MulV(transformed.P[0])
				viewed.P[1] = view.MulV(transformed.P[1])
				viewed.P[2] = view.MulV(transformed.P[2])
				viewed.T[0] = transformed.T[0]
				viewed.T[1] = transformed.T[1]
				viewed.T[2] = transformed.T[2]

				//Clip Triangle
				clipped := screenClipPlane.Clip(&viewed)

				for _, c := range clipped {
					projected := primitives.EmptyTriangle()
					// Project triangles from 3D --> 2D
					projected.P[0] = renderer.ProjectionMatrix.MulV(c.P[0])
					projected.P[1] = renderer.ProjectionMatrix.MulV(c.P[1])
					projected.P[2] = renderer.ProjectionMatrix.MulV(c.P[2])

					w1 := projected.P[0].W
					w2 := projected.P[1].W
					w3 := projected.P[2].W

					projected.P[0] = projected.P[0].Div(w1)
					projected.P[1] = projected.P[1].Div(w2)
					projected.P[2] = projected.P[2].Div(w3)
					projected.T[0] = c.T[0].Div(w1)
					projected.T[1] = c.T[1].Div(w2)
					projected.T[2] = c.T[2].Div(w3)

					// projected.T[0].U = projected.T[0].U / projected.P[0].W
					// projected.T[1].U = projected.T[1].U / projected.P[1].W
					// projected.T[2].U = projected.T[2].U / projected.P[2].W
					// projected.T[0].V = projected.T[0].V / projected.P[0].W
					// projected.T[1].V = projected.T[1].V / projected.P[1].W
					// projected.T[2].V = projected.T[2].V / projected.P[2].W
					// projected.T[0].W = 1.0 / projected.P[0].W
					// projected.T[1].W = 1.0 / projected.P[1].W
					// projected.T[2].W = 1.0 / projected.P[2].W

					offset := primitives.NewVector3d(1.0, 1.0, 0)
					projected.P[0] = projected.P[0].Add(offset)
					projected.P[1] = projected.P[1].Add(offset)
					projected.P[2] = projected.P[2].Add(offset)

					projected.P[0].X *= 0.5 * float64(renderer.screenWidth)
					projected.P[0].Y *= 0.5 * float64(renderer.screenHeight)
					projected.P[1].X *= 0.5 * float64(renderer.screenWidth)
					projected.P[1].Y *= 0.5 * float64(renderer.screenHeight)
					projected.P[2].X *= 0.5 * float64(renderer.screenWidth)
					projected.P[2].Y *= 0.5 * float64(renderer.screenHeight)

					// projected.Color = color.RGBA{R: 255, G: 255, B: 255, A: uint8(luminance * 255)}

					trisToRender.Push(projected)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	trisToRender.Done()
}

func (r *Renderer) RasterizeTriangle(tri *primitives.Triangle, texture *primitives.Image) {
	x1 := tri.P[0].X
	y1 := tri.P[0].Y
	x2 := tri.P[1].X
	y2 := tri.P[1].Y
	x3 := tri.P[2].X
	y3 := tri.P[2].Y

	u1 := tri.T[0].U
	v1 := tri.T[0].V
	w1 := tri.T[0].W
	u2 := tri.T[1].U
	v2 := tri.T[1].V
	w2 := tri.T[1].W
	u3 := tri.T[2].U
	v3 := tri.T[2].V
	w3 := tri.T[2].W

	if y1 > y2 {
		Swap(&y1, &y2)
		Swap(&x1, &x2)
		Swap(&u1, &u2)
		Swap(&v1, &v2)
		Swap(&w1, &w2)
	}
	if y1 > y3 {
		Swap(&y1, &y3)
		Swap(&x1, &x3)
		Swap(&u1, &u3)
		Swap(&v1, &v3)
		Swap(&w1, &w3)
	}
	if y2 > y3 {
		Swap(&y2, &y3)
		Swap(&x2, &x3)
		Swap(&u2, &u3)
		Swap(&v2, &v3)
		Swap(&w2, &w3)
	}

	dy1 := y2 - y1
	dx1 := x2 - x1
	dv1 := v2 - v1
	du1 := u2 - u1
	dw1 := w2 - w1

	dy2 := y3 - y1
	dx2 := x3 - x1
	dv2 := v3 - v1
	du2 := u3 - u1
	dw2 := w3 - w1

	daxStep, dbxStep := 0.0, 0.0
	du1Step, du2Step := 0.0, 0.0
	dv1Step, dv2Step := 0.0, 0.0
	dw1Step, dw2Step := 0.0, 0.0

	if dy1 != 0 {
		daxStep = dx1 / math.Abs(dy1)
		du1Step = du1 / math.Abs(dy1)
		dv1Step = dv1 / math.Abs(dy1)
		dw1Step = dw1 / math.Abs(dy1)
	}
	if dy2 != 0 {
		dbxStep = dx2 / math.Abs(dy2)
		du2Step = du2 / math.Abs(dy2)
		dv2Step = dv2 / math.Abs(dy2)
		dw2Step = dw2 / math.Abs(dy2)
	}

	if dy1 != 0 {
		for i := y1; i <= y2; i++ {
			ax := x1 + (i-y1)*daxStep
			bx := x1 + (i-y1)*dbxStep

			texSU := u1 + (i-y1)*du1Step
			texSV := v1 + (i-y1)*dv1Step
			texSW := w1 + (i-y1)*dw1Step

			texEU := u1 + (i-y1)*du2Step
			texEV := v1 + (i-y1)*dv2Step
			texEW := w1 + (i-y1)*dw2Step

			if ax > bx {
				Swap(&ax, &bx)
				Swap(&texSU, &texEU)
				Swap(&texSV, &texEV)
				Swap(&texSW, &texEW)
			}

			texU := texSU
			texV := texSV
			texW := texSW

			tStep := 1.0 / (bx - ax)
			t := 0.0

			for j := ax; j < bx; j++ {
				texU = (1.0-t)*texSU + t*texEU
				texV = (1.0-t)*texSV + t*texEV
				texW = (1.0-t)*texSW + t*texEW
				r.DrawPixel(j, i, texW, texture.GetPixel(texU/texW, texV/texW))
				t += tStep
			}
		}
	}

	dy1 = y3 - y2
	dx1 = x3 - x2
	dv1 = v3 - v2
	du1 = u3 - u2
	dw1 = w3 - w2

	du1Step, dv1Step = 0, 0
	if dy1 != 0.0 {
		daxStep = dx1 / math.Abs(dy1)
		du1Step = du1 / math.Abs(dy1)
		dv1Step = dv1 / math.Abs(dy1)
		dw1Step = dw1 / math.Abs(dy1)
	}
	if dy2 != 0.0 {
		dbxStep = dx2 / math.Abs(dy2)
	}

	if dy1 != 0.0 {
		for i := y2; i <= y3; i++ {
			ax := x2 + (i-y2)*daxStep
			bx := x1 + (i-y1)*dbxStep

			texSU := u2 + (i-y2)*du1Step
			texSV := v2 + (i-y2)*dv1Step
			texSW := w2 + (i-y2)*dw1Step

			texEU := u1 + (i-y1)*du2Step
			texEV := v1 + (i-y1)*dv2Step
			texEW := w1 + (i-y1)*dw2Step

			if ax > bx {
				Swap(&ax, &bx)
				Swap(&texSU, &texEU)
				Swap(&texSV, &texEV)
				Swap(&texSW, &texEW)
			}

			texU := texSU
			texV := texSV
			texW := texSW

			tStep := 1.0 / (bx - ax)
			t := 0.0

			for j := ax; j < bx; j++ {
				texU = (1.0-t)*texSU + t*texEU
				texV = (1.0-t)*texSV + t*texEV
				texW = (1.0-t)*texSW + t*texEW

				r.DrawPixel(j, i, texW, texture.GetPixel(texU/texW, texV/texW))
				t += tStep
			}
		}
	}
}
