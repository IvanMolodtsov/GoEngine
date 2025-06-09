package command

import (
	"github.com/IvanMolodtsov/GoEngine/object"
	"github.com/IvanMolodtsov/GoEngine/primitives"
)

type MoveCommand struct {
	obj         object.Movable
	translation *primitives.Vector3d
}

func NewMoveCommand(obj object.Movable, translation *primitives.Vector3d) *MoveCommand {
	var cmd MoveCommand
	cmd.obj = obj
	cmd.translation = translation
	return &cmd
}

func (cmd *MoveCommand) Invoke() error {
	cmd.obj.Move(cmd.translation)
	return nil
}
