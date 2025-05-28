package engine

type EventType string

const (
	Render     EventType = "Render"
	AddCommand EventType = "AddCommand"
	Quit       EventType = "Quit"
)

type Event interface {
	Type() EventType
}

type QuitEvent struct{}

func (c QuitEvent) Type() EventType {
	return Quit
}

func NewQuitEvent() *QuitEvent {
	return &QuitEvent{}
}
