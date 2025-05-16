package main

import (
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

	var cube engine.Mesh
	cube.Tris = []engine.Triangle{
		// South
		engine.NewTriangle(0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0),
		engine.NewTriangle(0.0, 0.0, 0.0, 1.0, 1.0, 0.0, 1.0, 0.0, 0.0),
		// East
		engine.NewTriangle(1.0, 0.0, 0.0, 1.0, 1.0, 0.0, 1.0, 1., 1.0),
		engine.NewTriangle(1.0, 0.0, 0.0, 1.0, 1.0, 1.0, 1.0, 0.0, 1.0),
		// North
		engine.NewTriangle(1.0, 0.0, 1.0, 1.0, 1.0, 1.0, 0.0, 1.0, 1.0),
		engine.NewTriangle(1.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0),
		// West
		engine.NewTriangle(0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0, 0.0),
		engine.NewTriangle(0.0, 0.0, 1.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0),
		// Top
		engine.NewTriangle(0.0, 1.0, 0.0, 0.0, 1.0, 1.0, 1.0, 1.0, 1.0),
		engine.NewTriangle(0.0, 1.0, 0.0, 1.0, 1.0, 1.0, 1.0, 1.0, 0.0),
		// Bottom
		engine.NewTriangle(1.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0),
		engine.NewTriangle(1.0, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0),
	}
	var fTheta = 1.0
	for game.IsRunning {
		engine.ReadEvents(game)
		game.FrameStart()

		fTheta += float64(game.DeltaTime.Seconds())

		rotZ := engine.ZRotationMatrix(fTheta + 0.5)
		rotX := engine.XRotationMatrix(fTheta)

		translation := engine.TranslationMatrix(0, 0, 5.0)

		worldMatrix := engine.IdentityMatrix()
		worldMatrix = rotZ.MulM(rotX)
		worldMatrix = worldMatrix.MulM(translation)

		game.Renderer.Render(func() {
			for _, t := range cube.Tris {
				var transformed, projected engine.Triangle

				// World Matrix Transform
				transformed.Points[0] = worldMatrix.MulV(t.Points[0])
				transformed.Points[1] = worldMatrix.MulV(t.Points[1])
				transformed.Points[2] = worldMatrix.MulV(t.Points[2])

				// Get surface normal
				normal := transformed.Normal()
				cameraRay := transformed.Points[0].Sub(game.Renderer.Camera)

				// check if triangle if visible
				if normal.DotProduct(cameraRay) < 0.0 {

					// Project triangles from 3D --> 2D
					projected.Points[0] = game.Renderer.ProjectionMatrix.MulV(transformed.Points[0])
					projected.Points[1] = game.Renderer.ProjectionMatrix.MulV(transformed.Points[1])
					projected.Points[2] = game.Renderer.ProjectionMatrix.MulV(transformed.Points[2])
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

					projected.Render(game.Renderer)

				}

			}
		})

		game.FrameEnd()
	}
}
