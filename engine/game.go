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
	IsRunning  bool
	frameStart time.Time
	DeltaTime  float64
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
	elapsed := time.Since(game.frameStart)
	game.DeltaTime = float64(elapsed.Nanoseconds())
}
