package primitives

type Object struct {
	Mesh        Mesh
	rotation    *Vector3d
	translation *Vector3d
}

func (obj *Object) GetTranslation() *Vector3d {
	return obj.translation
}

func (obj *Object) SetTranslation(translation *Vector3d) {
	obj.translation = translation
}

func (obj *Object) GetWorld() *Matrix4x4 {
	t := TranslationMatrix(obj.translation)
	xRot := XRotationMatrix(obj.rotation.X)
	yRot := YRotationMatrix(obj.rotation.Y)
	zRot := ZRotationMatrix(obj.rotation.Z)
	worldMatrix := xRot.MulM(yRot.MulM(zRot)).MulM(t)
	return worldMatrix
}

func NewObject(m Mesh, translation, rotation *Vector3d) *Object {
	var o Object
	o.Mesh = m
	o.translation = translation
	o.rotation = rotation

	return &o
}
