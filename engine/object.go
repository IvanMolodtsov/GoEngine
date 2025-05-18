package engine

type Object struct {
	Mesh        Mesh
	Rotation    Matrix4x4
	Translation Matrix4x4
}

func (obj *Object) GetWorld() *Matrix4x4 {
	worldMatrix := obj.Rotation.MulM(obj.Translation)
	return &worldMatrix
}

func NewObject(m Mesh, translation Vector3d, rotation float64) Object {
	var o Object
	o.Mesh = m
	o.Translation = TranslationMatrix(translation.X, translation.Y, translation.Z)
	o.Rotation = XRotationMatrix(rotation)

	return o
}
