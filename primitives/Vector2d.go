package primitives

type Vector2d struct {
	U, V, W float64
}

func NewVector2d(u, v float64) *Vector2d {
	var vector Vector2d
	vector.U = u
	vector.V = v
	vector.W = 1.0
	return &vector
}

func (v1 *Vector2d) Add(v2 *Vector2d) *Vector2d {
	var result Vector2d
	result.U = v1.U + v2.U
	result.V = v1.V + v2.V
	result.W = v1.U + v2.W

	return &result
}

func (v1 *Vector2d) Sub(v2 *Vector2d) *Vector2d {
	var result Vector2d
	result.U = v1.U - v2.U
	result.V = v1.V - v2.V
	result.W = v1.U - v2.W

	return &result
}

func (v1 *Vector2d) Mul(k float64) *Vector2d {
	var result Vector2d
	result.U = v1.U * k
	result.V = v1.V * k
	result.W = v1.U * k

	return &result
}

func (v1 *Vector2d) Div(k float64) *Vector2d {
	var result Vector2d
	result.U = v1.U / k
	result.V = v1.V / k
	result.W = v1.U / k

	return &result
}
