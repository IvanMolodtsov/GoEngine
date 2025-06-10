package engine

import (
	"math"
	"time"

	"github.com/IvanMolodtsov/GoEngine/command"
	"github.com/IvanMolodtsov/GoEngine/object"
	"github.com/IvanMolodtsov/GoEngine/primitives"
	"github.com/IvanMolodtsov/GoEngine/sdl"
)

type Signal struct{}

type Game struct {
	Width        int64
	Height       int64
	Renderer     *Renderer
	Camera       *Camera
	IsRunning    bool
	frameStart   time.Time
	DeltaTime    time.Duration
	eventQueue   chan sdl.Event
	commandQueue *Queue[command.Command]
	renderQueue  chan Signal
}

func InitGame(width int64, height int64) (*Game, error) {
	var game Game
	game.IsRunning = true
	game.Height = height
	game.Width = width
	renderer, err := InitRenderer(width, height)
	if err != nil {
		return nil, err
	}
	game.Renderer = renderer
	game.Camera = InitCamera()
	game.eventQueue = make(chan sdl.Event)
	game.renderQueue = make(chan Signal)
	game.commandQueue = NewQueue[command.Command]()

	return &game, nil
}

func (game *Game) destroy() {
	game.Renderer.Destroy()
}

func (game *Game) FrameStart() {
	game.frameStart = time.Now()
}

func (game *Game) FrameEnd() {
	game.DeltaTime = time.Since(game.frameStart)
}

func (game *Game) HandleEvents() {
	var deg90 = primitives.YRotationMatrix(math.Pi / 2)
	for game.IsRunning {
		event := <-game.eventQueue
		switch event.Type() {
		case sdl.EventKeyDown:
			if event.Key().Scancode == sdl.ScancodeEscape {
				game.IsRunning = false
			}

			velocity := 8.0 * game.DeltaTime.Seconds()
			forward := game.Camera.GetDirection().Mul(velocity)
			side := deg90.MulV(forward)
			scanCode := event.Key().Scancode
			switch scanCode {
			case sdl.ScancodeUp:
				game.commandQueue.Push(command.NewMoveCommand(game.Camera, primitives.NewVector3d(0, velocity, 0)))
			case sdl.ScancodeDown:
				game.commandQueue.Push(command.NewMoveCommand(game.Camera, primitives.NewVector3d(0, -velocity, 0)))
			case sdl.ScancodeW:
				game.commandQueue.Push(command.NewMoveCommand(game.Camera, forward))
			case sdl.ScancodeS:
				game.commandQueue.Push(command.NewMoveCommand(game.Camera, forward.Negative()))
			case sdl.ScancodeA:
				game.commandQueue.Push(command.NewMoveCommand(game.Camera, side.Negative()))
			case sdl.ScancodeD:
				game.commandQueue.Push(command.NewMoveCommand(game.Camera, side))
			}
		case sdl.EventMouseMotion:
			e := event.MouseMotion()
			game.commandQueue.Push(command.NewRotateCommand(game.Camera, primitives.NewVector3d(float64(e.YRel), float64(e.XRel), 0).Mul(0.001)))
		}
	}
}

func (game *Game) RunCommands() {
	for game.IsRunning {
		cmd, ok := game.commandQueue.Pop()
		if ok {
			cmd.Invoke()
		}
	}
}

func (game *Game) Run(entities []*object.UObject) {
	go game.HandleEvents()
	go game.RunCommands()

	pipe := NewPipeline(game.Camera, game.Renderer, true)
	// command.NewRotateCommand(entities[0], primitives.NewVector3d(1, 1, 1)).Invoke()
	for game.IsRunning {
		game.FrameStart()
		ReadEvents(game)
		pipe.Render(entities)
		game.Renderer.Render()
		game.FrameEnd()
	}

}
