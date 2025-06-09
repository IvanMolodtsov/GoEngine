package command

import (
	"github.com/IvanMolodtsov/GoEngine/object"
	"github.com/IvanMolodtsov/GoEngine/primitives"
)

type RotateCommand struct {
	obj      object.Rotatable
	rotation *primitives.Vector3d
}

func NewRotateCommand(obj object.Rotatable, rotation *primitives.Vector3d) *RotateCommand {
	var cmd RotateCommand
	cmd.obj = obj
	cmd.rotation = rotation
	return &cmd
}

func (cmd *RotateCommand) Invoke() error {
	cmd.obj.Rotate(cmd.rotation)
	return nil
}
