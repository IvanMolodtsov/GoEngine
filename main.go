package main

import (
	"runtime"

	"github.com/IvanMolodtsov/GoEngine/engine"
	"github.com/IvanMolodtsov/GoEngine/sdl"
)

const (
	WIDTH  int64 = 1024
	HEIGHT int64 = 768
	FPS    int64 = 120
)

func main() {
	runtime.GOMAXPROCS(4)
	runtime.LockOSThread()
	game, err := engine.Init(WIDTH, HEIGHT)
	if err != nil {
		panic(err)
	}
	defer game.Quit()

	sdl.SetWindowRelativeMouseMode(game.Window)

	var cube engine.Mesh // = engine.ReadFile("./axis.obj")
	cube.Tris = []*engine.Triangle{
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
	o := engine.NewObject(cube, engine.NewVector3d(0, 0, 5), engine.NewVector3d(90, 90, 90))

	// game.Loop([]*engine.Object{o})

	game.Run([]*engine.Object{o})

	runtime.UnlockOSThread()
}
