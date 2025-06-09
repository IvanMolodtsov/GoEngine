package object

import "github.com/IvanMolodtsov/GoEngine/primitives"

type Entity interface {
	GetMesh() *primitives.Mesh
	GetWorld() *primitives.Matrix4x4
}

func (obj *UObject) GetWorld() *primitives.Matrix4x4 {
	t := obj.GetTranslationMatrix()
	rot := obj.GetRotationMatrix()
	return rot.MulM(t)
}

func (obj *UObject) GetMesh() *primitives.Mesh {
	return obj.Get("Mesh").(*primitives.Mesh)
}

func NewEntity(m *primitives.Mesh, translation, rotation *primitives.Vector3d) *UObject {
	obj := New()

	obj.Set("Mesh", m)
	obj.Set("Position", translation)
	obj.Set("Rotation", rotation)

	return obj
}
