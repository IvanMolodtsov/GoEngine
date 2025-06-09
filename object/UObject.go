package object

import "sync"

type UObject struct {
	store map[string]any
	mut   sync.RWMutex
}

func New() *UObject {
	var obj UObject
	obj.store = make(map[string]any)
	return &obj
}

func (obj *UObject) Get(key string) any {
	obj.mut.RLock()
	defer obj.mut.RUnlock()
	return obj.store[key]
}

func (obj *UObject) Set(key string, value any) {
	obj.mut.Lock()
	defer obj.mut.Unlock()
	obj.store[key] = value
}
