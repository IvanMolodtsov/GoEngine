package primitives

type Plane struct {
	P *Vector3d
	N *Vector3d
}

func NewPlane(p, n *Vector3d) *Plane {
	var plane Plane
	plane.P = p
	plane.N = n.Normalize()
	return &plane
}

func (p *Plane) Intersection(lineStart, lineEnd *Vector3d) (*Vector3d, float64) {
	planeD := -1.0 * p.N.DotProduct(p.P)
	ad := lineStart.DotProduct(p.N)
	bd := lineEnd.DotProduct(p.N)
	t := (-planeD - ad) / (bd - ad)
	lineStartToEnd := lineEnd.Sub(lineStart)
	lineToIntersect := lineStartToEnd.Mul(t)
	return lineStart.Add(lineToIntersect), t
}

func (plane *Plane) Clip(tri *Triangle) []*Triangle {

	var result = make([]*Triangle, 0)
	var insidePoints, outsidePoints []*Vector3d
	var insideTPoints, outsideTPoints []*Vector2d
	d0 := plane.Dist(tri.P[0])
	d1 := plane.Dist(tri.P[1])
	d2 := plane.Dist(tri.P[2])

	if d0 >= 0 {
		insidePoints = append(insidePoints, tri.P[0])
		insideTPoints = append(insideTPoints, tri.T[0])
	} else {
		outsidePoints = append(outsidePoints, tri.P[0])
		outsideTPoints = append(outsideTPoints, tri.T[0])
	}
	if d1 >= 0 {
		insidePoints = append(insidePoints, tri.P[1])
		insideTPoints = append(insideTPoints, tri.T[1])
	} else {
		outsidePoints = append(outsidePoints, tri.P[1])
		outsideTPoints = append(outsideTPoints, tri.T[1])
	}
	if d2 >= 0 {
		insidePoints = append(insidePoints, tri.P[2])
		insideTPoints = append(insideTPoints, tri.T[2])
	} else {
		outsidePoints = append(outsidePoints, tri.P[2])
		outsideTPoints = append(outsideTPoints, tri.T[2])
	}

	if len(insidePoints) == 0 {
		return result
	}

	if len(insidePoints) == 3 {
		return append(result, tri)
	}

	if (len(insidePoints) == 1) && (len(outsidePoints) == 2) {
		var outT Triangle

		outT.Color = tri.Color

		outT.P[0] = insidePoints[0]
		outT.T[0] = insideTPoints[0]

		var t float64
		outT.P[1], t = plane.Intersection(insidePoints[0], outsidePoints[0])
		outT.T[1] = outsideTPoints[0].Sub(insideTPoints[0]).Mul(t).Add(insideTPoints[0])

		outT.P[2], t = plane.Intersection(insidePoints[0], outsidePoints[1])
		outT.T[2] = outsideTPoints[1].Sub(insideTPoints[0]).Mul(t).Add(insideTPoints[0])

		return append(result, &outT)
	}

	if (len(insidePoints) == 2) && (len(outsidePoints) == 1) {

		var outT1, outT2 Triangle
		var t float64

		outT1.Color = tri.Color
		outT2.Color = tri.Color
		outT1.P[0] = insidePoints[0]
		outT1.T[0] = insideTPoints[0]
		outT1.P[1] = insidePoints[1]
		outT1.T[1] = insideTPoints[1]
		outT1.P[2], t = plane.Intersection(insidePoints[0], outsidePoints[0])

		outT1.T[2] = outsideTPoints[0].Sub(insideTPoints[0]).Mul(t).Add(insideTPoints[0])

		outT2.P[0] = insidePoints[1]
		outT2.T[0] = insideTPoints[1]
		outT2.P[1] = outT1.P[2]
		outT2.T[1] = outT1.T[2]
		outT2.P[2], t = plane.Intersection(insidePoints[1], outsidePoints[0])

		outT2.T[2] = outsideTPoints[0].Sub(insideTPoints[1]).Mul(t).Add(insideTPoints[1])
		return append(result, &outT1, &outT2)
	}

	return result
}

func (plane *Plane) Dist(p *Vector3d) float64 {
	// n := v.Normalize()
	return plane.N.X*p.X + plane.N.Y*p.Y + plane.N.Z*p.Z - plane.N.DotProduct(plane.P)
}
