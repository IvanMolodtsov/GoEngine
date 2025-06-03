package primitives

type Rotatable interface {
	GetRotation() *Vector3d
	SetRotation(rotation *Vector3d)
}
