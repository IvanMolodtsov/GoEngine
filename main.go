package main

import "github.com/IvanMolodtsov/GoEngine/ioc"

const (
	WIDTH  int64 = 1024
	HEIGHT int64 = 768
	FPS    int64 = 120
)

func test(a int) (int, error) {
	return a + 1, nil
}

var testKey = ioc.NewKey[int, int]("test")

func main() {
	res, err := ioc.Resolve(ioc.Register, ioc.RegisterArgs{
		Key:        "test",
		Dependency: ioc.ToDependency(test),
	})
	println(err)
	println(res)

	res2, err2 := ioc.Resolve(testKey, 1)
	println(err2)
	println(res2)

	res3, err3 := ioc.Resolve(ioc.Remove, testKey.Value)
	println(err3)
	println(res3)

	re4, err4 := ioc.Resolve(testKey, 1)
	println(err4)
	println(re4)
	// runtime.GOMAXPROCS(4)
	// runtime.LockOSThread()
	// game, err := engine.Init(WIDTH, HEIGHT)
	// if err != nil {
	// 	panic(err)
	// }
	// defer game.Quit()
	// println("hello")

	// // var v = engine.Vector2{X: 1, Y: 1}

	// var cube engine.Mesh // = engine.ReadFile("./axis.obj")
	// cube.Tris = []*engine.Triangle{
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
	// o := engine.NewObject(cube, engine.NewVector3d(0, 0, 5), engine.NewVector3d(90, 90, 90))

	// // game.Loop([]*engine.Object{o})

	// game.Run([]*engine.Object{o})

	// runtime.UnlockOSThread()
}
