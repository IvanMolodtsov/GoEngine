package engine

import "github.com/IvanMolodtsov/GoEngine/shared"

type AddCommandEvent struct {
	Command shared.Command
}

func (c AddCommandEvent) Type() EventType {
	return AddCommand
}

func NewAddCommandEvent(cmd shared.Command) *AddCommandEvent {
	var e AddCommandEvent
	e.Command = cmd
	return &e
}
