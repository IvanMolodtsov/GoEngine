package engine

import "sync"

type Queue[T any] struct {
	values []T
	mut    sync.Mutex
}

func NewQueue[T any]() *Queue[T] {
	var q Queue[T]
	q.values = make([]T, 0)
	return &q
}

func (q *Queue[T]) Push(elem T) {
	q.mut.Lock()
	defer q.mut.Unlock()
	q.values = append(q.values, elem)
}

func (q *Queue[T]) Pop() (T, bool) {
	if len(q.values) == 0 {
		return *new(T), false
	}
	q.mut.Lock()
	defer q.mut.Unlock()
	res := q.values[len(q.values)-1]
	q.values = q.values[:len(q.values)-1]
	return res, true
}
