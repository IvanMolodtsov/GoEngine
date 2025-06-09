package object

import "github.com/IvanMolodtsov/GoEngine/primitives"

type Rotatable interface {
	GetRotation() *primitives.Vector3d
	Rotate(rotation *primitives.Vector3d)
	GetRotationMatrix() *primitives.Matrix4x4
}

func (obj *UObject) GetRotation() *primitives.Vector3d {
	return obj.Get("Rotation").(*primitives.Vector3d)
}

func (obj *UObject) Rotate(rotation *primitives.Vector3d) {
	position := obj.GetRotation()
	obj.Set("Rotation", position.Add(rotation))
}

func (obj *UObject) GetRotationMatrix() *primitives.Matrix4x4 {
	rotations := obj.GetRotation()
	xRot := primitives.XRotationMatrix(rotations.X)
	yRot := primitives.YRotationMatrix(rotations.Y)
	zRot := primitives.ZRotationMatrix(rotations.Z)
	return xRot.MulM(yRot.MulM(zRot))
}
