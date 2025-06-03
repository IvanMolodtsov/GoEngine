package primitives

type Movable interface {
	GetPosition() *Vector3d
	Move(translation *Vector3d)
}
