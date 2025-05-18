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
