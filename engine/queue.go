package engine

type Queue[T any] struct {
	Values  chan T
	counter uint32
}

func NewQueue[T any]() *Queue[T] {
	var q Queue[T]
	q.counter = 0
	q.Values = make(chan T)
	return &q
}

func (q *Queue[T]) Push(elem T) {
	q.Values <- elem
}

func (q *Queue[T]) Add() {
	q.counter++
}

func (q *Queue[T]) Done() {
	q.counter--
	if q.counter == 0 {
		close(q.Values)
	}
}

func (q *Queue[T]) Collect() []T {
	collector := make([]T, 0)
	for v := range q.Values {
		collector = append(collector, v)
	}
	return collector
}

func (q *Queue[T]) ForEach(operand func(T)) {
	for q.counter != 0 {
		select {
		case v := <-q.Values:
			operand(v)
		default:
		}
	}
}
