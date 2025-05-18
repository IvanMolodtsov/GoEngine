package engine

type Camera struct {
	Position  Vector3d
	Direction Vector3d
	Yaw       float64
}

func InitCamera() *Camera {
	var c Camera
	c.Position = NewVector3d(0.0, 0.0, 0.0)
	c.Direction = NewVector3d(0.0, 0.0, 1.0)
	c.Yaw = 0.0
	return &c
}

func (camera *Camera) GetView() *Matrix4x4 {
	up := NewVector3d(0, 1, 0)
	target := NewVector3d(0, 0, 1)
	cameraRotation := YRotationMatrix(camera.Yaw)
	camera.Direction = cameraRotation.MulV(target)

	target = camera.Position.Add(camera.Direction)

	pointAt := PointAtMatrix(camera.Position, target, up)
	view := pointAt.Inverse()
	return &view
}
