package ioc

import "sync"

type scopes struct {
	current Scope
	root    *RootScope
}

func (s *scopes) GetCurrent() Scope {
	return s.current
}

func (s *scopes) SetCurrent(new Scope) Scope {
	current := s.current
	s.current = new
	return current
}

func (s *scopes) GetRoot() Scope {
	return s.root
}

type Scopes interface {
	GetCurrent() Scope
	SetCurrent(Scope) Scope

	GetRoot() Scope
}

var once sync.Once
var instance *scopes = nil

func GetInstance() Scopes {
	if instance == nil {
		once.Do(func() {
			instance = new(scopes)
			var root RootScope
			root.store = make(map[string]D)
			root.store["Register"] = ToDependency(RegisterDep)
			root.store["Remove"] = ToDependency(RemoveDep)
			instance.root = &root
			instance.current = &root
		})
	}

	return instance
}
