package ioc

type Dependency[P any, R any] = func(args P) (R, error)
type D Dependency[any, any]

func ToDependency[P any, R any](f func(args P) (R, error)) D {
	return func(args any) (any, error) {
		return f(args.(P))
	}
}
