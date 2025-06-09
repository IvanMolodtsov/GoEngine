package engine

import (
	"github.com/IvanMolodtsov/GoEngine/object"
	"github.com/IvanMolodtsov/GoEngine/primitives"
)

type Camera struct {
	object.UObject
}

func InitCamera() *Camera {
	var c Camera
	c.UObject = *object.New()
	c.Set("Position", primitives.NewVector3d(0.0, 0.0, 0.0))
	c.Set("Rotation", primitives.NewVector3d(0.0, 0.0, 0.0))
	c.Set("Direction", primitives.NewVector3d(0.0, 0.0, 1.0))
	return &c
}

// func (c *Camera) getRotationMatrix() *primitives.Matrix4x4 {
// 	rotations := c.GetRotation()
// 	xRot := primitives.XRotationMatrix(rotations.X)
// 	yRot := primitives.YRotationMatrix(rotations.Y)
// 	zRot := primitives.ZRotationMatrix(rotations.Z)
// 	return xRot.MulM(yRot.MulM(zRot))
// }

func (c *Camera) GetDirection() *primitives.Vector3d {
	return c.Get("Direction").(*primitives.Vector3d)
}

func (camera *Camera) GetView() *primitives.Matrix4x4 {
	up := primitives.NewVector3d(0, 1, 0)
	target := primitives.NewVector3d(0, 0, 1)
	cameraRotation := camera.GetRotationMatrix()
	camera.Set("Direction", cameraRotation.MulV(target))
	direction := camera.GetDirection()

	target = camera.GetPosition().Add(direction)

	pointAt := primitives.PointAtMatrix(camera.GetPosition(), target, up)
	view := pointAt.Inverse()
	return view
}
