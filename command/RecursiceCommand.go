package command

import "time"

type RecursiveCommand struct {
	CMD   Command
	Queue chan Command
}

func (cmd *RecursiveCommand) Invoke() error {
	err := cmd.CMD.Invoke()
	if err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	cmd.Queue <- cmd
	return nil
}
