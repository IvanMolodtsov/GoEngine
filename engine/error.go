package engine

import (
	"errors"

	"github.com/jupiterrider/purego-sdl3/sdl"
)

func GetError() error {
	return errors.New(sdl.GetError())
}
