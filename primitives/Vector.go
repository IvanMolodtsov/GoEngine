package primitives

type Vector interface {
	Negative() Vector
	Add(v2 Vector) Vector
	Sub(v2 Vector) Vector
	DotProduct(v2 Vector) float64
	CrossProduct(v2 Vector) Vector
	Mul(k float64) Vector
	Div(k float64) Vector
	Length() float64
	Normalize() Vector
}
