package engine

import "math"

type Vector2d struct {
	X, Y, W float64
}

type Vector3d struct {
	X, Y, Z, W float64
}

func NewVector3d(x, y, z float64) Vector3d {
	return Vector3d{
		X: x,
		Y: y,
		Z: z,
		W: 1.0,
	}
}

func (v1 Vector3d) Print() {
	print(v1.X, ", ", v1.Y, ", ", v1.Z, ", ", v1.W)
}

func (v1 Vector3d) Add(v2 Vector3d) Vector3d {
	return NewVector3d(
		v1.X+v2.X,
		v1.Y+v2.Y,
		v1.Z+v2.Z,
	)
}

func (v1 Vector3d) Sub(v2 Vector3d) Vector3d {
	return NewVector3d(
		v1.X-v2.X,
		v1.Y-v2.Y,
		v1.Z-v2.Z,
	)
}

func (v1 Vector3d) DotProduct(v2 Vector3d) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

func (v1 Vector3d) CrossProduct(v2 Vector3d) Vector3d {
	return NewVector3d(
		v1.Y*v2.Z-v1.Z*v2.Y,
		v1.Z*v2.X-v1.X*v2.Z,
		v1.X*v2.Y-v1.Y*v2.X,
	)
}

func (v Vector3d) Mul(k float64) Vector3d {
	return NewVector3d(
		v.X*k,
		v.Y*k,
		v.Z*k,
	)
}

func (v Vector3d) Div(k float64) Vector3d {
	return NewVector3d(
		v.X/k,
		v.Y/k,
		v.Z/k,
	)
}

func (v Vector3d) Length() float64 {
	return math.Sqrt(v.DotProduct(v))
}

func (v Vector3d) Normalize() Vector3d {
	l := v.Length()

	return NewVector3d(
		v.X/l,
		v.Y/l,
		v.Z/l,
	)
}
