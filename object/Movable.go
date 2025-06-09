package object

import "github.com/IvanMolodtsov/GoEngine/primitives"

type Movable interface {
	GetPosition() *primitives.Vector3d
	Move(translation *primitives.Vector3d)
	GetTranslationMatrix() *primitives.Matrix4x4
}

func (obj *UObject) GetPosition() *primitives.Vector3d {
	return obj.Get("Position").(*primitives.Vector3d)
}

func (obj *UObject) Move(translation *primitives.Vector3d) {
	position := obj.GetPosition()
	obj.Set("Position", position.Add(translation))
}

func (obj *UObject) GetTranslationMatrix() *primitives.Matrix4x4 {
	pos := obj.GetPosition()
	return primitives.TranslationMatrix(pos)
}
