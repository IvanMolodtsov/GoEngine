package command

type Command interface {
	Invoke() error
}
