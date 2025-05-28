package ioc

type Key[P any, R any] struct {
	Value string
}

func NewKey[P, R any](value string) Key[P, R] {
	var key Key[P, R]
	key.Value = value
	return key
}
