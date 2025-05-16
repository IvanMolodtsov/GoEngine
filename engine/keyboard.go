package engine

import "github.com/jupiterrider/purego-sdl3/sdl"

func ReadEvents(game *Game) {
	var event sdl.Event

	sdl.PollEvent(&event)
	switch event.Type() {
	case sdl.EventKeyDown:
		if event.Key().Scancode == sdl.ScancodeEscape {
			game.IsRunning = false
		}
	}

}
