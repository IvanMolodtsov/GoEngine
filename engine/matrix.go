package engine

import "math"

type Matrix4x4 struct {
	M [4][4]float64
}

func (m Matrix4x4) MulV(v Vector3d) Vector3d {
	var res Vector3d

	res.X = v.X*m.M[0][0] + v.Y*m.M[1][0] + v.Z*m.M[2][0] + v.W*m.M[3][0]
	res.Y = v.X*m.M[0][1] + v.Y*m.M[1][1] + v.Z*m.M[2][1] + v.W*m.M[3][1]
	res.Z = v.X*m.M[0][2] + v.Y*m.M[1][2] + v.Z*m.M[2][2] + v.W*m.M[3][2]
	res.W = v.X*m.M[0][3] + v.Y*m.M[1][3] + v.Z*m.M[2][3] + v.W*m.M[3][3]

	return res
}

func (m1 Matrix4x4) MulM(m2 Matrix4x4) Matrix4x4 {
	var matrix Matrix4x4
	for c := range 4 {
		for r := range 4 {
			matrix.M[r][c] = m1.M[r][0]*m2.M[0][c] + m1.M[r][1]*m2.M[1][c] + m1.M[r][2]*m2.M[2][c] + m1.M[r][3]*m2.M[3][c]
		}
	}

	return matrix
}

func (m Matrix4x4) Inverse() Matrix4x4 {
	var matrix Matrix4x4

	matrix.M[0][0] = m.M[0][0]
	matrix.M[0][1] = m.M[1][0]
	matrix.M[0][2] = m.M[2][0]
	matrix.M[0][3] = 0.0
	matrix.M[1][0] = m.M[0][1]
	matrix.M[1][1] = m.M[1][1]
	matrix.M[1][2] = m.M[2][1]
	matrix.M[1][3] = 0.0
	matrix.M[2][0] = m.M[0][2]
	matrix.M[2][1] = m.M[1][2]
	matrix.M[2][2] = m.M[2][2]
	matrix.M[2][3] = 0.0
	matrix.M[3][0] = -(m.M[3][0]*matrix.M[0][0] + m.M[3][1]*matrix.M[1][0] + m.M[3][2]*matrix.M[2][0])
	matrix.M[3][1] = -(m.M[3][0]*matrix.M[0][1] + m.M[3][1]*matrix.M[1][1] + m.M[3][2]*matrix.M[2][1])
	matrix.M[3][2] = -(m.M[3][0]*matrix.M[0][2] + m.M[3][1]*matrix.M[1][2] + m.M[3][2]*matrix.M[2][2])
	matrix.M[3][3] = 1.0

	return matrix
}

func IdentityMatrix() Matrix4x4 {
	var m Matrix4x4

	m.M[0][0] = 1.0
	m.M[1][1] = 1.0
	m.M[2][2] = 1.0
	m.M[3][3] = 1.0

	return m
}

func XRotationMatrix(angle float64) Matrix4x4 {
	var matrix Matrix4x4
	matrix.M[0][0] = 1.0
	matrix.M[1][1] = math.Cos(angle)
	matrix.M[1][2] = math.Sin(angle)
	matrix.M[2][1] = -1.0 * math.Sin(angle)
	matrix.M[2][2] = math.Cos(angle)
	matrix.M[3][3] = 1.0
	return matrix
}

func YRotationMatrix(angle float64) Matrix4x4 {
	var matrix Matrix4x4
	matrix.M[0][0] = math.Cos(angle)
	matrix.M[1][1] = math.Sin(angle)
	matrix.M[1][2] = -1.0 * math.Sin(angle)
	matrix.M[2][1] = 1.0
	matrix.M[2][2] = math.Cos(angle)
	matrix.M[3][3] = 1.0
	return matrix
}

func ZRotationMatrix(angle float64) Matrix4x4 {
	var matrix Matrix4x4
	matrix.M[0][0] = math.Cos(angle)
	matrix.M[0][1] = math.Sin(angle)
	matrix.M[1][0] = -1.0 * math.Sin(angle)
	matrix.M[1][1] = math.Cos(angle)
	matrix.M[2][2] = 1.0
	matrix.M[3][3] = 1.0
	return matrix
}

func TranslationMatrix(x, y, z float64) Matrix4x4 {
	matrix := IdentityMatrix()

	matrix.M[3][0] = x
	matrix.M[3][1] = y
	matrix.M[3][2] = z

	return matrix
}

func ProjectionMatrix(aspectRatio, fNear, fFar float64, fov float64) Matrix4x4 {
	var matrix Matrix4x4
	rFOV := 1.0 / math.Tan(fov*0.5/180.0*math.Pi)
	matrix.M[0][0] = aspectRatio * rFOV
	matrix.M[1][1] = rFOV
	matrix.M[2][2] = fFar / (fFar - fNear)
	matrix.M[3][2] = (-fFar * fNear) / (fFar - fNear)
	matrix.M[2][3] = 1.0
	matrix.M[3][3] = 0.0

	return matrix
}

func (m Matrix4x4) Print() {
	println(m.M[0][0], " ", m.M[0][1], " ", m.M[0][2], " ", m.M[0][3])
	println(m.M[1][0], " ", m.M[1][1], " ", m.M[1][2], " ", m.M[1][3])
	println(m.M[2][0], " ", m.M[2][1], " ", m.M[2][2], " ", m.M[2][3])
	println(m.M[3][0], " ", m.M[3][1], " ", m.M[3][2], " ", m.M[3][3])
}
