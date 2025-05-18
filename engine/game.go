package engine

import (
	"time"

	"github.com/jupiterrider/purego-sdl3/sdl"
)

type Game struct {
	Width      int64
	Height     int64
	window     *sdl.Window
	Renderer   *Renderer
	Camera     *Camera
	IsRunning  bool
	frameStart time.Time
	DeltaTime  time.Duration
}

func InitGame(width int64, height int64) (*Game, error) {
	var game Game
	game.IsRunning = true
	game.Height = height
	game.Width = width
	w := sdl.CreateWindow("Game", int32(width), int32(height), 0)
	if w == nil {
		return nil, GetError()
	}
	game.window = w
	renderer, err := InitRenderer(game.window, width, height)
	if err != nil {
		return nil, err
	}
	game.Renderer = renderer
	game.Camera = InitCamera()

	return &game, nil
}

func (game *Game) destroy() {
	sdl.DestroyWindow(game.window)
	game.Renderer.destroy()
}

func (game *Game) FrameStart() {
	game.frameStart = time.Now()
}

func (game *Game) FrameEnd() {
	game.DeltaTime = time.Since(game.frameStart)

}

func (game *Game) Loop(entities []Object) {

	for game.IsRunning {
		renderQueue := make(chan Triangle)
		// close := make(chan interface{})
		// wg := new(sync.WaitGroup)
		ReadEvents(game)
		game.FrameStart()

		for _, o := range entities {
			// wg.Add(1)
			go game.Renderer.Project(o, game.Camera, renderQueue)
		}

		render := make([]Triangle, 0)
		// wg.Wait()
		for t := range renderQueue {
			render = append(render, t)
		}

		// wg.Wait()
		// close(renderQueue)

		render = Sort(render, func(a, b Triangle) bool {
			z1 := (a.Points[0].Z + a.Points[1].Z + a.Points[2].Z) / 3.0
			z2 := (b.Points[0].Z + b.Points[1].Z + b.Points[2].Z) / 3.0
			return z1 > z2
		})

		game.Renderer.Render(render)

		game.FrameEnd()
	}

}
