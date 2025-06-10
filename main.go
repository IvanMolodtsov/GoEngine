package main

import (
	"runtime"

	"github.com/IvanMolodtsov/GoEngine/sdl"
)

const (
	WIDTH  int64 = 1024
	HEIGHT int64 = 768
	FPS    int64 = 120
)

func main() {
	// runtime.GOMAXPROCS(8)
	runtime.LockOSThread()
	// game, err := engine.Init(WIDTH, HEIGHT)
	// if err != nil {
	// 	panic(err)
	// }
	// defer game.Quit()

	// sdl.SetWindowRelativeMouseMode(game.Renderer.Window)

	// texture := primitives.LoadImage("Sprite.png")
	// var cube primitives.Mesh // = engine.ReadFile("./axis.obj")
	// cube.Tris = []*primitives.Triangle{
	// 	// South
	// 	primitives.NewTriangle(0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 0, 1, 0, 0, 1, 0),
	// 	primitives.NewTriangle(0.0, 0.0, 0.0, 1.0, 1.0, 0.0, 1.0, 0.0, 0.0, 0, 1, 1, 0, 1, 1),
	// 	// East
	// 	primitives.NewTriangle(1.0, 0.0, 0.0, 1.0, 1.0, 0.0, 1.0, 1.0, 1.0, 0, 1, 0, 0, 1, 0),
	// 	primitives.NewTriangle(1.0, 0.0, 0.0, 1.0, 1.0, 1.0, 1.0, 0.0, 1.0, 0, 1, 1, 0, 1, 1),
	// 	// North
	// 	primitives.NewTriangle(1.0, 0.0, 1.0, 1.0, 1.0, 1.0, 0.0, 1.0, 1.0, 0, 1, 0, 0, 1, 0),
	// 	primitives.NewTriangle(1.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, 0, 1, 1, 0, 1, 1),
	// 	// West
	// 	primitives.NewTriangle(0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0, 0.0, 0, 1, 0, 0, 1, 0),
	// 	primitives.NewTriangle(0.0, 0.0, 1.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0, 1, 1, 0, 1, 1),
	// 	// Top
	// 	primitives.NewTriangle(0.0, 1.0, 0.0, 0.0, 1.0, 1.0, 1.0, 1.0, 1.0, 0, 1, 0, 0, 1, 0),
	// 	primitives.NewTriangle(0.0, 1.0, 0.0, 1.0, 1.0, 1.0, 1.0, 1.0, 0.0, 0, 1, 1, 0, 1, 1),
	// 	// Bottom
	// 	primitives.NewTriangle(1.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0, 1, 0, 0, 1, 0),
	// 	primitives.NewTriangle(1.0, 0.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0, 1, 1, 0, 1, 1),
	// }
	// cube.Texture = texture
	// // r, g, b, a := cube.Texture.Data.At(16, 16).RGBA()
	// // println(uint8(r), " ", uint8(g), " ", uint8(b), " ", uint8(a))
	// o := object.NewEntity(&cube, primitives.NewVector3d(0, 0, 5), primitives.NewVector3d(90, 90, 90))
	// // game.Loop([]*engine.Object{o})

	// game.Run([]*object.UObject{o})

	err := sdl.Init(sdl.InitVideo)
	if err != nil {
		panic(err.Error())
	}

	window, err := sdl.CreateWindow("GPU test", int32(WIDTH), int32(HEIGHT), 0)
	if err != nil {
		panic(err.Error())
	}

	gpu, err := sdl.CreateGPUDevice(sdl.SPIRV, true, "")
	if err != nil {
		panic(err.Error())
	}

	err = sdl.ClaimWindowForGPUDevice(gpu, window)
	if err != nil {
		panic(err.Error())
	}

loop:
	for {
		var event sdl.Event
		for sdl.PollEvent(&event) {
			switch event.Type() {
			case sdl.EventQuit:
				break loop
			case sdl.EventKeyDown:
				if event.Key().Scancode == sdl.ScancodeEscape {
					break loop
				}
			}
		}

		cmdBuf, err := sdl.AcquireGPUCommandBuffer(gpu)
		if err != nil {
			panic(err.Error())
		}

		var swapChainTexture *sdl.GPUTexture

		err = sdl.WaitAndAcquireGPUSwapchainTexture(cmdBuf, window, &swapChainTexture)
		if err != nil {
			panic(err.Error())
		}

		colorTarget := sdl.GPUColorTargetInfo{
			Texture: swapChainTexture,
			LoadOp:  sdl.GPU_LOADOP_CLEAR,
			ClearColor: sdl.FColor{
				R: 0, G: 100, B: 100, A: 255,
			},
			StoreOp: sdl.GPU_STOREOP_STORE,
		}

		pass, err := sdl.BeginGPURenderPass(cmdBuf, colorTarget, 1, nil)
		if err != nil {
			panic(err.Error())
		}

		sdl.EndGPURenderPass(pass)
	}

	runtime.UnlockOSThread()
}
