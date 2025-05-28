package engine

type Camera struct {
	position  *Vector3d
	rotation  *Vector3d
	Direction *Vector3d
	Yaw       float64
}

func (c *Camera) GetTranslation() *Vector3d {
	return c.position
}

func (c *Camera) SetTranslation(translation *Vector3d) {
	c.position = translation
}

func (c *Camera) GetRotation() *Vector3d {
	return c.rotation
}

func (c *Camera) SetRotation(rotation *Vector3d) {
	c.rotation = rotation
}

func InitCamera() *Camera {
	var c Camera
	c.position = NewVector3d(0.0, 0.0, 0.0)
	c.rotation = NewVector3d(0.0, 0, 0)
	c.Direction = NewVector3d(0.0, 0.0, 1.0)
	return &c
}

func (c *Camera) getRotationMatrix() *Matrix4x4 {
	rotations := c.rotation
	xRot := XRotationMatrix(rotations.X)
	yRot := YRotationMatrix(rotations.Y)
	zRot := ZRotationMatrix(rotations.Z)
	return xRot.MulM(yRot.MulM(zRot))
}

func (camera *Camera) GetView() *Matrix4x4 {
	up := NewVector3d(0, 1, 0)
	target := NewVector3d(0, 0, 1)
	cameraRotation := camera.getRotationMatrix()
	camera.Direction = cameraRotation.MulV(target)

	target = camera.position.Add(camera.Direction)

	pointAt := PointAtMatrix(camera.position, target, up)
	view := pointAt.Inverse()
	return view
}
