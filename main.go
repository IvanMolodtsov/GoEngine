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

	var cube engine.Mesh // = engine.ReadFile("./axis.obj")
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
	o := engine.NewObject(cube, engine.NewVector3d(0, 0, 5), 0.0)
	// o.Rotation = engine.XRotationMatrix(50)

	game.Loop([]engine.Object{o})
}
