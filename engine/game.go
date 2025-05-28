package engine

import (
	"time"

	"github.com/IvanMolodtsov/GoEngine/sdl"
	"github.com/IvanMolodtsov/GoEngine/shared"
)

type Game struct {
	Width        int64
	Height       int64
	Window       *sdl.Window
	Renderer     *Renderer
	Camera       *Camera
	IsRunning    bool
	frameStart   time.Time
	DeltaTime    time.Duration
	eventQueue   chan Event
	commandQueue chan shared.Command
	renderQueue  chan []*Triangle
}

func InitGame(width int64, height int64) (*Game, error) {
	var game Game
	game.IsRunning = true
	game.Height = height
	game.Width = width
	w, err := sdl.CreateWindow("Game", int32(width), int32(height), 0)
	if err != nil {
		return nil, err
	}
	game.Window = w
	renderer, err := InitRenderer(game.Window, width, height)
	if err != nil {
		return nil, err
	}
	game.Renderer = renderer
	game.Camera = InitCamera()
	game.eventQueue = make(chan Event)
	game.renderQueue = make(chan []*Triangle)
	game.commandQueue = make(chan shared.Command)

	return &game, nil
}

func (game *Game) destroy() {
	sdl.DestroyWindow(game.Window)
	game.Renderer.destroy()
}

func (game *Game) FrameStart() {
	game.frameStart = time.Now()
}

func (game *Game) FrameEnd() {
	game.DeltaTime = time.Since(game.frameStart)
}

func (game *Game) HandleEvents() {
	for game.IsRunning {
		select {
		case event := <-game.eventQueue:
			println(event.Type())
			switch event.Type() {
			case Quit:
				game.IsRunning = false
			case AddCommand:
				game.commandQueue <- event.(*AddCommandEvent).Command
			}
		default:
		}
	}
}

func (game *Game) GenerateFrames(entities []*Object) {
	for game.IsRunning {
		game.FrameStart()

		projected := NewQueue[*Triangle]()
		counter := 0
		for _, o := range entities {
			projected.Add()
			counter++
			go game.Renderer.Project(o, game.Camera, projected)
		}
		render := projected.Collect()
		render = Sort(render, func(a, b *Triangle) bool {
			z1 := (a.Points[0].Z + a.Points[1].Z + a.Points[2].Z) / 3.0
			z2 := (b.Points[0].Z + b.Points[1].Z + b.Points[2].Z) / 3.0
			return z1 > z2
		})

		trisToRender := game.Renderer.ClipTriangles(render)

		game.renderQueue <- trisToRender
		game.FrameEnd()
		println("generated frames")
	}
}

func (game *Game) RunCommands() {
	for game.IsRunning {
		select {
		case cmd := <-game.commandQueue:
			cmd.Invoke()
		default:
		}
	}
}

func (game *Game) Run(entities []*Object) {

	go game.HandleEvents()
	go game.RunCommands()
	go game.GenerateFrames(entities)

	for game.IsRunning {
		ReadEvents(game)

		select {
		case tris := <-game.renderQueue:
			println("render: ", len(tris))
			game.Renderer.Render(tris)
		default:
		}
	}

}
