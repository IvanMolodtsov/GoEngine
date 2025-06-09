package primitives

import "image/color"

type Triangle struct {
	P         [3]*Vector3d
	T         [3]*Vector2d
	Color     color.Color
	Texture   *Image
	Luminance float64
}

func NewTriangle(x1, y1, z1, x2, y2, z2, x3, y3, z3, tx1, ty1, tx2, ty2, tx3, ty3 float64) *Triangle {
	var t Triangle
	t.P[0] = NewVector3d(x1, y1, z1)
	t.P[1] = NewVector3d(x2, y2, z2)
	t.P[2] = NewVector3d(x3, y3, z3)
	t.T[0] = NewVector2d(tx1, ty1)
	t.T[1] = NewVector2d(tx2, ty2)
	t.T[2] = NewVector2d(tx3, ty3)
	t.Luminance = 0
	return &t
}

func EmptyTriangle() *Triangle {
	var t Triangle
	t.P[0] = NewVector3d(0, 0, 0)
	t.P[1] = NewVector3d(0, 0, 0)
	t.P[2] = NewVector3d(0, 0, 0)
	t.T[0] = NewVector2d(0, 0)
	t.T[1] = NewVector2d(0, 0)
	t.T[2] = NewVector2d(0, 0)
	t.Luminance = 0
	return &t
}

func (t *Triangle) Print() {
	t.P[0].Print()
	t.P[1].Print()
	t.P[2].Print()
	println()
}

func (t *Triangle) Normal() *Vector3d {
	line1 := t.P[1].Sub(t.P[0])
	line2 := t.P[2].Sub(t.P[0])

	normal := line1.CrossProduct(line2)
	return normal.Normalize()
}
