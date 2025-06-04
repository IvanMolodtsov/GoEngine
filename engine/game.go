package engine

import (
	"math"
	"time"

	"github.com/IvanMolodtsov/GoEngine/command"
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
	commandQueue chan command.Command
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
	game.commandQueue = make(chan command.Command)

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
			forward := game.Camera.Direction.Mul(velocity)
			side := deg90.MulV(forward)
			scanCode := event.Key().Scancode
			switch scanCode {
			case sdl.ScancodeUp:
				game.commandQueue <- command.NewMoveCommand(game.Camera, primitives.NewVector3d(0, -velocity, 0))
			case sdl.ScancodeDown:
				game.commandQueue <- command.NewMoveCommand(game.Camera, primitives.NewVector3d(0, velocity, 0))
			case sdl.ScancodeW:
				game.commandQueue <- command.NewMoveCommand(game.Camera, forward)
			case sdl.ScancodeS:
				game.commandQueue <- command.NewMoveCommand(game.Camera, forward.Negative())
			case sdl.ScancodeA:
				game.commandQueue <- command.NewMoveCommand(game.Camera, side)
			case sdl.ScancodeD:
				game.commandQueue <- command.NewMoveCommand(game.Camera, side.Negative())
			}
		case sdl.EventMouseMotion:
			e := event.MouseMotion()
			game.commandQueue <- command.NewRotateCommand(game.Camera, primitives.NewVector3d(float64(e.YRel), float64(e.XRel), 0).Mul(0.001))
		}
	}
}

func (game *Game) GenerateFrames(entities []*primitives.Object) {
	for game.IsRunning {
		game.FrameStart()

		projected := NewQueue[*primitives.Triangle]()
		counter := 0
		for _, o := range entities {
			projected.Add()
			counter++
			go game.Renderer.Project(o, game.Camera, projected)
		}
		render := projected.Collect()
		// render = Sort(render, func(a, b *primitives.Triangle) bool {
		// 	z1 := (a.P[0].Z + a.P[1].Z + a.P[2].Z) / 3.0
		// 	z2 := (b.P[0].Z + b.P[1].Z + b.P[2].Z) / 3.0
		// 	return z1 > z2
		// })

		trisToRender := game.Renderer.ClipTriangles(render)

		game.Renderer.PushTriangles(trisToRender, entities[0].Mesh.Texture)
		game.FrameEnd()
		game.renderQueue <- Signal{}
	}
}

func (game *Game) RunCommands() {
	for game.IsRunning {
		cmd := <-game.commandQueue
		cmd.Invoke()

	}
}

func (game *Game) Run(entities []*primitives.Object) {
	go game.HandleEvents()
	go game.RunCommands()
	go game.GenerateFrames(entities)

	for game.IsRunning {
		ReadEvents(game)
		select {
		case <-game.renderQueue:
			game.Renderer.Render()
		case <-time.After(20 * time.Millisecond):
		}
	}

}
