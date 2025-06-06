package engine

import (
	"image/color"
	"math"

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

func (renderer *Renderer) Render() {
	sdl.RenderClear(renderer.renderer)
	renderer.SwapTextures()
	sdl.RenderPresent(renderer.renderer)
}

func (renderer *Renderer) SwapTextures() {
	newCurrent := (renderer.current + 1) % 2
	renderer.textures[newCurrent].ClearBuffer()
	renderer.textures[newCurrent].Lock()
	renderer.textures[newCurrent].Clear()

	renderer.textures[renderer.current].Unlock(renderer.renderer)
	renderer.current = newCurrent

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
	worldMatrix := o.GetWorld()
	view := camera.GetView()
	screenClipPlane := renderer.screenClipPlane
	for _, T := range o.Mesh.Tris {
		func(t *primitives.Triangle) {
			transformed := primitives.EmptyTriangle()

			// World Matrix Transform
			transformed.P[0] = worldMatrix.MulV(t.P[0])
			transformed.P[1] = worldMatrix.MulV(t.P[1])
			transformed.P[2] = worldMatrix.MulV(t.P[2])

			// Get surface normal
			normal := transformed.Normal()
			cameraRay := transformed.P[0].Sub(camera.GetPosition())

			// check if triangle if visible
			if normal.DotProduct(cameraRay) < 0.0 {
				// lightDir := primitives.NewVector3d(0, 0, -1.0)
				// lightDir = lightDir.Normalize()

				// luminance := normal.DotProduct(lightDir)

				viewed := primitives.EmptyTriangle()

				viewed.P[0] = view.MulV(transformed.P[0])
				viewed.P[1] = view.MulV(transformed.P[1])
				viewed.P[2] = view.MulV(transformed.P[2])
				viewed.T[0] = t.T[0].Copy()
				viewed.T[1] = t.T[1].Copy()
				viewed.T[2] = t.T[2].Copy()

				//Clip Triangle
				clipped := screenClipPlane.Clip(viewed)

				for _, c := range clipped {
					// Project triangles from 3D --> 2D
					projected := primitives.EmptyTriangle()
					projected.P[0] = renderer.ProjectionMatrix.MulV(c.P[0])
					projected.P[1] = renderer.ProjectionMatrix.MulV(c.P[1])
					projected.P[2] = renderer.ProjectionMatrix.MulV(c.P[2])
					projected.T[0] = c.T[0].Copy()
					projected.T[1] = c.T[1].Copy()
					projected.T[2] = c.T[2].Copy()

					projected.T[0] = projected.T[0].Div(projected.P[0].W)
					projected.T[1] = projected.T[1].Div(projected.P[1].W)
					projected.T[2] = projected.T[2].Div(projected.P[2].W)
					projected.P[0] = projected.P[0].Div(projected.P[0].W)
					projected.P[1] = projected.P[1].Div(projected.P[1].W)
					projected.P[2] = projected.P[2].Div(projected.P[2].W)

					// X/Y are inverted so put them back
					projected.P[0].X *= -1.0
					projected.P[1].X *= -1.0
					projected.P[2].X *= -1.0
					projected.P[0].Y *= -1.0
					projected.P[1].Y *= -1.0
					projected.P[2].Y *= -1.0

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

					projected.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}

					trisToRender.Push(projected)
				}
			}
		}(T)
	}
	trisToRender.Done()
}

func (r *Renderer) DrawTriangleWireframe(t *primitives.Triangle, color color.Color) {
	r.DrawLine(t.P[0].X, t.P[0].Y, t.P[1].X, t.P[1].Y, primitives.ToHex(color))
	r.DrawLine(t.P[0].X, t.P[0].Y, t.P[2].X, t.P[2].Y, primitives.ToHex(color))
	r.DrawLine(t.P[2].X, t.P[2].Y, t.P[1].X, t.P[1].Y, primitives.ToHex(color))
}

func (r *Renderer) DrawTriangle(tri *primitives.Triangle, texture *primitives.Image) {
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

			var texU = texSU
			var texV = texSV
			var texW = texSW

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

			var texU = texSU
			var texV = texSV
			var texW = texSW

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
		renderer.DrawPixel(x, y, 255, color)
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
			renderer.DrawPixel(x, y, 255, color)
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
		renderer.DrawPixel(x, y, 255, color)
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
			renderer.DrawPixel(x, y, 255, color)
		}
	}
}
