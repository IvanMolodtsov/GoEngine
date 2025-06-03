package engine

import "github.com/IvanMolodtsov/GoEngine/primitives"

type Camera struct {
	position  *primitives.Vector3d
	rotation  *primitives.Vector3d
	Direction *primitives.Vector3d
	Yaw       float64
}

func (c *Camera) GetPosition() *primitives.Vector3d {
	return c.position
}

func (c *Camera) Move(translation *primitives.Vector3d) {
	c.position = c.position.Add(translation)
}

func (c *Camera) GetRotation() *primitives.Vector3d {
	return c.rotation
}

func (c *Camera) SetRotation(rotation *primitives.Vector3d) {
	c.rotation = rotation
}

func InitCamera() *Camera {
	var c Camera
	c.position = primitives.NewVector3d(0.0, 0.0, 0.0)
	c.rotation = primitives.NewVector3d(0.0, 0, 0)
	c.Direction = primitives.NewVector3d(0.0, 0.0, 1.0)
	return &c
}

func (c *Camera) getRotationMatrix() *primitives.Matrix4x4 {
	rotations := c.rotation
	xRot := primitives.XRotationMatrix(rotations.X)
	yRot := primitives.YRotationMatrix(rotations.Y)
	zRot := primitives.ZRotationMatrix(rotations.Z)
	return xRot.MulM(yRot.MulM(zRot))
}

func (camera *Camera) GetView() *primitives.Matrix4x4 {
	up := primitives.NewVector3d(0, 1, 0)
	target := primitives.NewVector3d(0, 0, 1)
	cameraRotation := camera.getRotationMatrix()
	camera.Direction = cameraRotation.MulV(target)

	target = camera.position.Add(camera.Direction)

	pointAt := primitives.PointAtMatrix(camera.position, target, up)
	view := pointAt.Inverse()
	return view
}
