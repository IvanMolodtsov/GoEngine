package engine

import "math"

const epsilon = 1e-9

func Equal(a, b float32) bool {
	return math.Abs(float64(a-b)) <= epsilon
}

func CollectChan[T any](c chan T) []T {
	var s []T
	for v := range c {
		s = append(s, v)
	}
	return s
}
