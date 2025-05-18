package engine

import "image/color"

type Triangle struct {
	Points [3]Vector3d
	Color  color.RGBA
}

func NewTriangle(x1, y1, z1, x2, y2, z2, x3, y3, z3 float64) Triangle {
	var t Triangle
	t.Points[0] = NewVector3d(x1, y1, z1)
	t.Points[1] = NewVector3d(x2, y2, z2)
	t.Points[2] = NewVector3d(x3, y3, z3)
	return t
}

func (t Triangle) Print() {
	t.Points[0].Print()
	t.Points[0].Print()
	t.Points[0].Print()
	println()
}

func (t Triangle) Render(renderer *Renderer) {
	renderer.DrawLine(t.Points[0].X, t.Points[0].Y, t.Points[1].X, t.Points[1].Y, uint32(0xff00ff00))
	renderer.DrawLine(t.Points[1].X, t.Points[1].Y, t.Points[2].X, t.Points[2].Y, uint32(0xff00ff00))
	renderer.DrawLine(t.Points[2].X, t.Points[2].Y, t.Points[0].X, t.Points[0].Y, uint32(0xff00ff00))
}

func (t Triangle) Normal() Vector3d {
	line1 := t.Points[1].Sub(t.Points[0])
	line2 := t.Points[2].Sub(t.Points[0])

	normal := line1.CrossProduct(line2)
	return normal.Normalize()
}
