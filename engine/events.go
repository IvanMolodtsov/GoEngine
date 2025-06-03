package engine

// import "github.com/jupiterrider/purego-sdl3/sdl"
import (
	"github.com/IvanMolodtsov/GoEngine/sdl"
)

func ReadEvents(game *Game) {
	var event sdl.Event
	for sdl.PollEvent(&event) {
		select {
		case game.eventQueue <- event:
		default:
		}
	}

}
