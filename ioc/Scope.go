package ioc

import (
	"errors"
)

type Scope interface {
	Get(key string) (D, error)
}

type MutableScope interface {
	Set(key string, dependency D)
	Remove(key string)
}

type BasicScope struct {
	parent Scope
	store  map[string]D
}

func (scope *BasicScope) Get(key string) (Dependency[any, any], error) {
	result, ok := scope.store[key]
	if !ok {
		return scope.parent.Get(key)
	}
	return result, nil
}

func (scope *BasicScope) Set(key string, dependency Dependency[any, any]) {
	scope.store[key] = dependency
}

func (scope *BasicScope) Remove(key string) {
	delete(scope.store, key)
}

type RootScope struct {
	store map[string]D
}

func (scope *RootScope) Get(key string) (D, error) {
	println("Get ", key)
	result, ok := scope.store[key]
	println("result ", result, " ", ok)
	if !ok {
		return nil, errors.New("not_found")
	}
	return result, nil
}
