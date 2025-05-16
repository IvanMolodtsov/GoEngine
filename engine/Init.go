package engine

import "github.com/jupiterrider/purego-sdl3/sdl"

func Init(width int64, height int64) (*Game, error) {
	if !sdl.Init(sdl.InitVideo) {
		return nil, GetError()
	}
	return InitGame(width, height)
}

func Quit(game *Game) {
	game.destroy()
	sdl.Quit()
}
