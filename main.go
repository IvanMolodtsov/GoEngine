package main

import (
	"image/color"
	"sort"

	"github.com/IvanMolodtsov/GoEngine/engine"
)

const (
	WIDTH  int64 = 1024
	HEIGHT int64 = 768
	FPS    int64 = 120
)

func main() {

	game, err := engine.Init(WIDTH, HEIGHT)
	if err != nil {
		panic(err)
	}
	defer engine.Quit(game)
	println("hello")

	// var v = engine.Vector2{X: 1, Y: 1}

	var cube engine.Mesh = engine.ReadFile("./axis.obj")
	// cube.Tris = []engine.Triangle{
	// 	// South
	// 	engine.NewTriangle(0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0),
	// 	engine.NewTriangle(0.0, 0.0, 0.0, 1.0, 1.0, 0.0, 1.0, 0.0, 0.0),
	// 	// East
	// 	engine.NewTriangle(1.0, 0.0, 0.0, 1.0, 1.0, 0.0, 1.0, 1., 1.0),
	// 	engine.NewTriangle(1.0, 0.0, 0.0, 1.0, 1.0, 1.0, 1.0, 0.0, 1.0),
	// 	// North
	// 	engine.NewTriangle(1.0, 0.0, 1.0, 1.0, 1.0, 1.0, 0.0, 1.0, 1.0),
	// 	engine.NewTriangle(1.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0),
	// 	// West
	// 	engine.NewTriangle(0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0, 0.0),
	// 	engine.NewTriangle(0.0, 0.0, 1.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0),
	// 	// Top
	// 	engine.NewTriangle(0.0, 1.0, 0.0, 0.0, 1.0, 1.0, 1.0, 1.0, 1.0),
	// 	engine.NewTriangle(0.0, 1.0, 0.0, 1.0, 1.0, 1.0, 1.0, 1.0, 0.0),
	// 	// Bottom
	// 	engine.NewTriangle(1.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0),
	// 	engine.NewTriangle(1.0, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0),
	// }
	var fTheta = 0.0
	for game.IsRunning {
		engine.ReadEvents(game)
		game.FrameStart()

		// fTheta += float64(game.DeltaTime.Seconds())

		rotZ := engine.ZRotationMatrix(fTheta)
		rotX := engine.XRotationMatrix(fTheta)

		translation := engine.TranslationMatrix(0, 0, 3.0)

		worldMatrix := engine.IdentityMatrix()
		worldMatrix = rotZ.MulM(rotX)
		worldMatrix = worldMatrix.MulM(translation)

		up := engine.NewVector3d(0, 1, 0)
		target := engine.NewVector3d(0, 0, 1)
		cameraRotation := engine.YRotationMatrix(game.Camera.Yaw)
		game.Camera.Direction = cameraRotation.MulV(target)

		target = game.Camera.Position.Add(game.Camera.Direction)

		camera := engine.PointAtMatrix(game.Camera.Position, target, up)
		view := camera.Inverse()

		var buffer = make([]engine.Triangle, 0)

		screenClipPlane := engine.NewPlane(engine.NewVector3d(0, 0, 0.1), engine.NewVector3d(0, 0, 1.0))

		game.Renderer.Render(func() {
			for _, t := range cube.Tris {
				var transformed, viewed, projected engine.Triangle

				// World Matrix Transform
				transformed.Points[0] = worldMatrix.MulV(t.Points[0])
				transformed.Points[1] = worldMatrix.MulV(t.Points[1])
				transformed.Points[2] = worldMatrix.MulV(t.Points[2])

				// Get surface normal
				normal := transformed.Normal()
				cameraRay := transformed.Points[0].Sub(game.Camera.Position)

				// check if triangle if visible
				if normal.DotProduct(cameraRay) < 0.0 {
					lightDir := engine.NewVector3d(0, 0, -1.0)
					lightDir = lightDir.Normalize()

					luminance := normal.DotProduct(lightDir)

					viewed.Points[0] = view.MulV(transformed.Points[0])
					viewed.Points[1] = view.MulV(transformed.Points[1])
					viewed.Points[2] = view.MulV(transformed.Points[2])

					//Clip Triangle
					clipped := screenClipPlane.Clip(viewed)

					for _, c := range clipped {
						// Project triangles from 3D --> 2D
						projected.Points[0] = game.Renderer.ProjectionMatrix.MulV(c.Points[0])
						projected.Points[1] = game.Renderer.ProjectionMatrix.MulV(c.Points[1])
						projected.Points[2] = game.Renderer.ProjectionMatrix.MulV(c.Points[2])
						projected.Points[0] = projected.Points[0].Div(projected.Points[0].W)
						projected.Points[1] = projected.Points[1].Div(projected.Points[1].W)
						projected.Points[2] = projected.Points[2].Div(projected.Points[2].W)

						offset := engine.NewVector3d(1.0, 1.0, 0)
						projected.Points[0] = projected.Points[0].Add(offset)
						projected.Points[1] = projected.Points[1].Add(offset)
						projected.Points[2] = projected.Points[2].Add(offset)

						projected.Points[0].X *= 0.5 * float64(WIDTH)
						projected.Points[0].Y *= 0.5 * float64(HEIGHT)
						projected.Points[1].X *= 0.5 * float64(WIDTH)
						projected.Points[1].Y *= 0.5 * float64(HEIGHT)
						projected.Points[2].X *= 0.5 * float64(WIDTH)
						projected.Points[2].Y *= 0.5 * float64(HEIGHT)

						projected.Color = color.RGBA{R: 255, G: 255, B: 255, A: uint8(luminance * 255)}

						buffer = append(buffer, projected)
					}

				}

			}
			sort.Slice(buffer, func(i, j int) bool {
				z1 := (buffer[i].Points[0].Z + buffer[i].Points[1].Z + buffer[i].Points[2].Z) / 3.0
				z2 := (buffer[j].Points[0].Z + buffer[j].Points[1].Z + buffer[j].Points[2].Z) / 3.0

				return z1 > z2
			})

			top := engine.NewPlane(engine.NewVector3d(0, 0, 0), engine.NewVector3d(0, 1, 0))
			bottom := engine.NewPlane(engine.NewVector3d(0, float64(game.Height-1), 0), engine.NewVector3d(0, -1, 0))
			left := engine.NewPlane(engine.NewVector3d(0, 0, 0), engine.NewVector3d(1, 0, 0))
			right := engine.NewPlane(engine.NewVector3d(float64(game.Width-1), 0, 0), engine.NewVector3d(-1, 0, 0))

			planes := []engine.Plane{
				top, bottom, left, right,
			}

			trisToRaster := make([]engine.Triangle, 0)

			for _, t := range buffer {

				listTriangles := make([]engine.Triangle, 0)
				listTriangles = append(listTriangles, t)
				var newT = 1
				for _, p := range planes {
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

			for _, t := range trisToRaster {
				t.Fill(game.Renderer)
				t.Render(game.Renderer)
			}
		})

		game.FrameEnd()
	}
}
