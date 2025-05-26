package engine

// import "github.com/jupiterrider/purego-sdl3/sdl"
import "github.com/IvanMolodtsov/GoEngine/sdl"

func Init(width int64, height int64) (*Game, error) {
	err := sdl.Init(sdl.InitVideo)
	if err != nil {
		return nil, GetError()
	}
	return InitGame(width, height)
}

func Quit(game *Game) {
	game.destroy()
	sdl.Quit()
}
