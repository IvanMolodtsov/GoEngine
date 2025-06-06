package engine

import (
	"image/color"

	"github.com/IvanMolodtsov/GoEngine/primitives"
)

type PipeLine struct {
	Camera   *Camera
	renderer *Renderer
	IsDebug  bool
}

func NewPipeline(camera *Camera, renderer *Renderer, isDebug bool) *PipeLine {
	var pipe PipeLine
	pipe.renderer = renderer
	pipe.Camera = camera
	pipe.IsDebug = isDebug
	return &pipe
}

func (pipe PipeLine) Render(objects []*primitives.Object) {
	for _, o := range objects {
		projected := pipe.Project(o)
		clipped := pipe.ClipTriangles(projected)
		pipe.RasterizeTriangles(clipped, o.Mesh.Texture)
	}
}

func (pipe PipeLine) Project(object *primitives.Object) []*primitives.Triangle {
	worldMatrix := object.GetWorld()
	view := pipe.Camera.GetView()
	screenClipPlane := pipe.renderer.screenClipPlane

	result := make([]*primitives.Triangle, 0)

	for _, t := range object.Mesh.Tris {
		transformed := primitives.EmptyTriangle()

		// World Matrix Transform
		transformed.P[0] = worldMatrix.MulV(t.P[0])
		transformed.P[1] = worldMatrix.MulV(t.P[1])
		transformed.P[2] = worldMatrix.MulV(t.P[2])

		// Get surface normal
		normal := transformed.Normal()
		cameraRay := transformed.P[0].Sub(pipe.Camera.GetPosition())

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
				projected.P[0] = pipe.renderer.ProjectionMatrix.MulV(c.P[0])
				projected.P[1] = pipe.renderer.ProjectionMatrix.MulV(c.P[1])
				projected.P[2] = pipe.renderer.ProjectionMatrix.MulV(c.P[2])
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

				projected.P[0].X *= 0.5 * float64(pipe.renderer.screenWidth)
				projected.P[0].Y *= 0.5 * float64(pipe.renderer.screenHeight)
				projected.P[1].X *= 0.5 * float64(pipe.renderer.screenWidth)
				projected.P[1].Y *= 0.5 * float64(pipe.renderer.screenHeight)
				projected.P[2].X *= 0.5 * float64(pipe.renderer.screenWidth)
				projected.P[2].Y *= 0.5 * float64(pipe.renderer.screenHeight)

				projected.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}

				result = append(result, projected)
			}
		}
	}
	return result
}

func (pipe *PipeLine) ClipTriangles(tris []*primitives.Triangle) []*primitives.Triangle {
	result := make([]*primitives.Triangle, 0)

	for _, t := range tris {
		listTriangles := make([]*primitives.Triangle, 0)
		listTriangles = append(listTriangles, t)
		var newT = 1
		for _, p := range pipe.renderer.planes {
			for newT > 0 {
				test := listTriangles[0]
				listTriangles = listTriangles[1:]
				newT -= 1

				clipped := p.Clip(test)
				listTriangles = append(listTriangles, clipped...)

			}
			newT = len(listTriangles)
		}
		result = append(result, listTriangles...)

	}
	return result
}

func (pipe *PipeLine) RasterizeTriangles(tris []*primitives.Triangle, texture *primitives.Image) {
	for _, t := range tris {
		pipe.renderer.DrawTriangle(t, texture)

		if pipe.IsDebug {
			pipe.renderer.DrawTriangleWireframe(t, color.RGBA{R: 150, G: 150, B: 150, A: 255})
		}
	}
}
