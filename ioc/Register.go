package ioc

import "errors"

type RegisterArgs struct {
	Key        string
	Dependency Dependency[any, any]
}

var Register = NewKey[RegisterArgs, bool]("Register")

func RegisterDep(args RegisterArgs) (bool, error) {
	current := GetInstance().GetCurrent()
	mutable, ok := current.(MutableScope)
	if !ok {
		return false, errors.New("current scope is immutable")
	}
	mutable.Set(args.Key, args.Dependency)
	return true, nil
}

var Remove = NewKey[string, bool]("Remove")

func RemoveDep(args string) (bool, error) {
	current := GetInstance().GetCurrent()
	mutable, ok := current.(MutableScope)
	if !ok {
		return false, errors.New("current scope is immutable")
	}
	mutable.Remove(args)
	return true, nil
}
