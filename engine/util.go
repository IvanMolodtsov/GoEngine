package engine

import "math"

const epsilon = 1e-9

func Equal(a, b float32) bool {
	return math.Abs(float64(a-b)) <= epsilon
}
