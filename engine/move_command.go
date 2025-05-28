package engine

type MoveCommand struct {
	obj         Movable
	translation *Vector3d
}

func NewMoveCommand(obj Movable, translation *Vector3d) *MoveCommand {
	var cmd MoveCommand
	cmd.obj = obj
	cmd.translation = translation
	return &cmd
}

func (cmd *MoveCommand) Invoke() error {
	t := cmd.obj.GetTranslation()
	result := t.Add(cmd.translation)
	cmd.obj.SetTranslation(result)
	println("moved")
	return nil
}
