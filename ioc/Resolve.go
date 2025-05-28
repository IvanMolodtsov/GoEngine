package ioc

import (
	"errors"
	"fmt"
)

func Resolve[P any, R any](key Key[P, R], args P) (R, error) {
	current := GetInstance().GetCurrent()
	d, err := current.Get(key.Value)
	if err != nil {
		return *new(R), err
	}
	result, err2 := d(any(args))
	if err2 != nil {
		return *new(R), fmt.Errorf("%s dependency invocation error: %v", key.Value, err)
	}
	res, ok := result.(R)
	if !ok {
		return *new(R), errors.New("type_mismatch")
	}
	return res, nil
}
