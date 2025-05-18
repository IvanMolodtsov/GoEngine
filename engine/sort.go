package engine

import (
	"sync"
)

func MergeChan[T interface{}](left chan T, right chan T, c chan T, sort func(a, b T) bool) {
	defer close(c)
	val, ok := <-left
	val2, ok2 := <-right
	for ok && ok2 {
		if sort(val, val2) {
			c <- val
			val, ok = <-left
		} else {
			c <- val2
			val2, ok2 = <-right
		}
	}
	for ok {
		c <- val
		val, ok = <-left
	}
	for ok2 {
		c <- val2
		val2, ok2 = <-right
	}
}
func MergeSort[T interface{}](arr []T, ch chan T, sort func(a, b T) bool) {
	if len(arr) < 2 {
		ch <- arr[0]
		defer close(ch)
		return
	}
	left := make(chan T)
	right := make(chan T)
	go MergeSort(arr[len(arr)/2:], left, sort)
	go MergeSort(arr[:len(arr)/2], right, sort)
	go MergeChan(left, right, ch, sort)
}
func Sort[T interface{}](a []T, sort func(a, b T) bool) []T {
	c := make(chan T)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		MergeSort(a, c, sort)
	}()
	wg.Wait()
	var s []T
	for v := range c {
		s = append(s, v)
	}
	return s
}
