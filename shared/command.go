package shared

type Command interface {
	Invoke() error
}
