package engine

// import "github.com/jupiterrider/purego-sdl3/sdl"
import "github.com/IvanMolodtsov/GoEngine/sdl"

func ReadEvents(game *Game) {
	var event sdl.Event

	sdl.PollEvent(&event)
	switch event.Type() {
	case sdl.EventKeyDown:
		if event.Key().Scancode == sdl.ScancodeEscape {
			game.IsRunning = false
		}

		velocity := 8.0 * game.DeltaTime.Seconds()

		if event.Key().Scancode == sdl.ScancodeUp {
			game.Camera.Position.Y += velocity
		} else if event.Key().Scancode == sdl.ScancodeDown {
			game.Camera.Position.Y -= velocity
		}

		if event.Key().Scancode == sdl.ScancodeLeft {
			game.Camera.Position.X += velocity
		} else if event.Key().Scancode == sdl.ScancodeRight {
			game.Camera.Position.X -= velocity
		}

		if event.Key().Scancode == sdl.ScancodeA {
			game.Camera.Yaw += 1.0 * game.DeltaTime.Seconds()
		} else if event.Key().Scancode == sdl.ScancodeD {
			game.Camera.Yaw -= 1.0 * game.DeltaTime.Seconds()
		}

		forward := game.Camera.Direction.Mul(velocity)

		if event.Key().Scancode == sdl.ScancodeW {
			game.Camera.Position = game.Camera.Position.Add(forward)
		} else if event.Key().Scancode == sdl.ScancodeS {
			game.Camera.Position = game.Camera.Position.Sub(forward)
		}

	}

}
