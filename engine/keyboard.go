package engine

// import "github.com/jupiterrider/purego-sdl3/sdl"
import "github.com/IvanMolodtsov/GoEngine/sdl"

func ReadEvents(game *Game) {
	var event sdl.Event

	sdl.PollEvent(&event)
	switch event.Type() {
	case sdl.EventKeyDown:
		if event.Key().Scancode == sdl.ScancodeEscape {
			game.eventQueue <- NewQuitEvent()
		}

		velocity := 8.0 * game.DeltaTime.Seconds()

		if event.Key().Scancode == sdl.ScancodeUp {
			game.commandQueue <- NewMoveCommand(game.Camera, NewVector3d(0, velocity, 0))
			// game.Camera.Position.Y += velocity
		} else if event.Key().Scancode == sdl.ScancodeDown {
			game.commandQueue <- NewMoveCommand(game.Camera, NewVector3d(0, -velocity, 0))
		}

		// if event.Key().Scancode == sdl.ScancodeLeft {
		// 	game.Camera.Position.X += velocity
		// } else if event.Key().Scancode == sdl.ScancodeRight {
		// 	game.Camera.Position.X -= velocity
		// }

		if event.Key().Scancode == sdl.ScancodeA {
			game.commandQueue <- NewRotateCommand(game.Camera, NewVector3d(0, 1.0*game.DeltaTime.Seconds(), 0))
		} else if event.Key().Scancode == sdl.ScancodeD {
			game.commandQueue <- NewRotateCommand(game.Camera, NewVector3d(0, -1.0*game.DeltaTime.Seconds(), 0))
		}

		forward := game.Camera.Direction.Mul(velocity)

		if event.Key().Scancode == sdl.ScancodeW {
			game.eventQueue <- (NewAddCommandEvent(NewMoveCommand(game.Camera, forward)))
		} else if event.Key().Scancode == sdl.ScancodeS {
			game.eventQueue <- (NewAddCommandEvent(NewMoveCommand(game.Camera, forward.Negative())))
		}

	}

}
