package main

import (
	"time"

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
		{
			Points: [3]engine.Vector3d{
				engine.NewVector3d(0.0, 0.0, 0.0),
				engine.NewVector3d(0.0, 1.0, 0.0),
				engine.NewVector3d(1.0, 1.0, 0.0),
			},
		}, {
			Points: [3]engine.Vector3d{
				engine.NewVector3d(0.0, 0.0, 0.0),
				engine.NewVector3d(1.0, 1.0, 0.0),
				engine.NewVector3d(1.0, 0.0, 0.0),
			},
		}, {
			Points: [3]engine.Vector3d{
				engine.NewVector3d(1.0, 0.0, 0.0),
				engine.NewVector3d(1.0, 1.0, 0.0),
				engine.NewVector3d(1.0, 1.0, 1.0),
			},
		}, {
			Points: [3]engine.Vector3d{
				engine.NewVector3d(1.0, 0.0, 0.0),
				engine.NewVector3d(1.0, 1.0, 1.0),
				engine.NewVector3d(1.0, 0.0, 1.0),
			},
		}, {
			Points: [3]engine.Vector3d{
				engine.NewVector3d(1.0, 0.0, 1.0),
				engine.NewVector3d(1.0, 1.0, 1.0),
				engine.NewVector3d(0.0, 1.0, 1.0),
			},
		}, {
			Points: [3]engine.Vector3d{
				engine.NewVector3d(1.0, 0.0, 1.0),
				engine.NewVector3d(0.0, 1.0, 1.0),
				engine.NewVector3d(0.0, 0.0, 1.0),
			},
		}, {
			Points: [3]engine.Vector3d{
				engine.NewVector3d(0.0, 0.0, 1.0),
				engine.NewVector3d(0.0, 1.0, 1.0),
				engine.NewVector3d(0.0, 1.0, 0.0),
			},
		}, {
			Points: [3]engine.Vector3d{
				engine.NewVector3d(0.0, 0.0, 1.0),
				engine.NewVector3d(0.0, 1.0, 0.0),
				engine.NewVector3d(0.0, 0.0, 0.0),
			},
		},
		{
			Points: [3]engine.Vector3d{
				engine.NewVector3d(0.0, 1.0, 0.0),
				engine.NewVector3d(0.0, 1.0, 1.0),
				engine.NewVector3d(1.0, 1.0, 1.0),
			},
		}, {
			Points: [3]engine.Vector3d{
				engine.NewVector3d(0.0, 1.0, 0.0),
				engine.NewVector3d(1.0, 1.0, 1.0),
				engine.NewVector3d(1.0, 1.0, 0.0),
			},
		},
		{
			Points: [3]engine.Vector3d{
				engine.NewVector3d(1.0, 0.0, 1.0),
				engine.NewVector3d(0.0, 0.0, 0.1),
				engine.NewVector3d(0.0, 0.0, 0.0),
			},
		}, {
			Points: [3]engine.Vector3d{
				engine.NewVector3d(1.0, 0.0, 1.0),
				engine.NewVector3d(0.0, 0.0, 0.0),
				engine.NewVector3d(1.0, 0.0, 0.0),
			},
		},
	}
	var fTheta = 1.0
	for game.IsRunning {
		engine.ReadEvents(game)
		time.Sleep(50 * time.Millisecond)
		game.FrameStart()

		fTheta += (game.DeltaTime / 10000000)

		rotZ := engine.ZRotationMatrix(fTheta + 0.5)
		rotX := engine.XRotationMatrix(fTheta)

		translation := engine.TranslationMatrix(0, 0, 5.0)

		worldMatrix := engine.IdentityMatrix()
		worldMatrix = rotZ.MulM(rotX)
		worldMatrix = worldMatrix.MulM(translation)

		game.Renderer.Render(func() {
			for _, t := range cube.Tris {
				var transformed, projected engine.Triangle

				// //Rotate Z
				// transformed.Points[0] = rotZ.MulV(t.Points[0])
				// transformed.Points[1] = rotZ.MulV(t.Points[1])
				// transformed.Points[2] = rotZ.MulV(t.Points[2])

				// //Rotate X
				// transformed.Points[0] = rotX.MulV(transformed.Points[0])
				// transformed.Points[1] = rotX.MulV(transformed.Points[1])
				// transformed.Points[2] = rotX.MulV(transformed.Points[2])

				// //Translate
				// transformed.Points[0] = translation.MulV(transformed.Points[0])
				// transformed.Points[1] = translation.MulV(transformed.Points[1])
				// transformed.Points[2] = translation.MulV(transformed.Points[2])

				// World Matrix Transform
				transformed.Points[0] = worldMatrix.MulV(t.Points[0])
				transformed.Points[1] = worldMatrix.MulV(t.Points[1])
				transformed.Points[2] = worldMatrix.MulV(t.Points[2])

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
		})

		game.FrameEnd()
	}
}
