package engine

type Queue[T any] struct {
	Values     chan T
	collection []T
	Counter    uint32
}

func NewQueue[T any]() *Queue[T] {
	var q Queue[T]
	q.Values = make(chan T)
	q.collection = make([]T, 0)
	q.Counter = 0
	return &q
}

func (q *Queue[T]) Start() {
	q.Counter += 1
}

func (q *Queue[T]) End() {
	q.Counter -= 1
	if q.Counter == 0 {
		close(q.Values)
	}
}

func (q *Queue[T]) Push(value T) {
	q.Values <- value
}

func (q *Queue[T]) Pop() T {
	return <-q.Values
}

func (q *Queue[T]) Close() {
	close(q.Values)
}
