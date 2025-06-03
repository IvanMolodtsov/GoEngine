package command

import "github.com/IvanMolodtsov/GoEngine/primitives"

type RotateCommand struct {
	obj      primitives.Rotatable
	rotation *primitives.Vector3d
}

func NewRotateCommand(obj primitives.Rotatable, rotation *primitives.Vector3d) *RotateCommand {
	var cmd RotateCommand
	cmd.obj = obj
	cmd.rotation = rotation
	return &cmd
}

func (cmd *RotateCommand) Invoke() error {
	t := cmd.obj.GetRotation()
	result := t.Add(cmd.rotation)
	cmd.obj.SetRotation(result)
	return nil
}
