package engine

type Movable interface {
	GetTranslation() *Vector3d
	SetTranslation(translation *Vector3d)
}
