package engine

// import "github.com/jupiterrider/purego-sdl3/sdl"
import (
	"math"

	"github.com/IvanMolodtsov/GoEngine/sdl"
)

var deg90 = YRotationMatrix(math.Pi / 2)

func ReadEvents(game *Game) {
	var event sdl.Event

	sdl.PollEvent(&event)
	switch event.Type() {
	case sdl.EventKeyDown:
		if event.Key().Scancode == sdl.ScancodeEscape {
			game.eventQueue <- NewQuitEvent()
		}

		velocity := 8.0 * game.DeltaTime.Seconds()
		forward := game.Camera.Direction.Mul(velocity)
		side := deg90.MulV(forward)
		scanCode := event.Key().Scancode
		switch scanCode {
		case sdl.ScancodeUp:
			game.commandQueue <- NewMoveCommand(game.Camera, NewVector3d(0, velocity, 0))
		case sdl.ScancodeDown:
			game.commandQueue <- NewMoveCommand(game.Camera, NewVector3d(0, -velocity, 0))
		case sdl.ScancodeW:
			game.eventQueue <- (NewAddCommandEvent(NewMoveCommand(game.Camera, forward)))
		case sdl.ScancodeS:
			game.eventQueue <- (NewAddCommandEvent(NewMoveCommand(game.Camera, forward.Negative())))
		case sdl.ScancodeA:
			game.eventQueue <- (NewAddCommandEvent(NewMoveCommand(game.Camera, side)))
		case sdl.ScancodeD:
			game.eventQueue <- (NewAddCommandEvent(NewMoveCommand(game.Camera, side.Negative())))
		}
	case sdl.EventMouseMotion:
		e := event.MouseMotion()
		game.commandQueue <- NewRotateCommand(game.Camera, NewVector3d(-float64(e.YRel), float64(e.XRel), 0).Mul(0.001))

	}

}
