package engine

type RotateCommand struct {
	obj      Rotatable
	rotation *Vector3d
}

func NewRotateCommand(obj Rotatable, rotation *Vector3d) *RotateCommand {
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
