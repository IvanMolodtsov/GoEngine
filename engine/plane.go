package engine

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

func (p *Plane) Intersection(lineStart, lineEnd *Vector3d) *Vector3d {
	planeD := -1.0 * p.N.DotProduct(p.P)
	ad := lineStart.DotProduct(p.N)
	bd := lineEnd.DotProduct(p.N)
	t := (-planeD - ad) / (bd - ad)
	lineStartToEnd := lineEnd.Sub(lineStart)
	lineToIntersect := lineStartToEnd.Mul(t)
	return lineStart.Add(lineToIntersect)
}

func (plane *Plane) Clip(t *Triangle) []*Triangle {
	var result = make([]*Triangle, 0)
	var insidePoints, outsidePoints []*Vector3d
	d0 := plane.Dist(t.Points[0])
	d1 := plane.Dist(t.Points[1])
	d2 := plane.Dist(t.Points[2])

	if d0 >= 0 {
		insidePoints = append(insidePoints, t.Points[0])
	} else {
		outsidePoints = append(outsidePoints, t.Points[0])
	}
	if d1 >= 0 {
		insidePoints = append(insidePoints, t.Points[1])
	} else {
		outsidePoints = append(outsidePoints, t.Points[1])
	}
	if d2 >= 0 {
		insidePoints = append(insidePoints, t.Points[2])
	} else {
		outsidePoints = append(outsidePoints, t.Points[2])
	}

	if len(insidePoints) == 0 {
		return result
	}

	if len(insidePoints) == 3 {
		return append(result, t)
	}

	if (len(insidePoints) == 1) && (len(outsidePoints) == 2) {
		var outT Triangle

		outT.Color = t.Color

		outT.Points[0] = insidePoints[0]

		outT.Points[1] = plane.Intersection(insidePoints[0], outsidePoints[0])
		outT.Points[2] = plane.Intersection(insidePoints[0], outsidePoints[1])

		return append(result, &outT)
	}

	if (len(insidePoints) == 2) && (len(outsidePoints) == 1) {
		var outT1, outT2 Triangle

		outT1.Color = t.Color
		outT2.Color = t.Color

		outT1.Points[0] = insidePoints[0]
		outT1.Points[1] = insidePoints[1]
		outT1.Points[2] = plane.Intersection(insidePoints[0], outsidePoints[0])

		outT2.Points[0] = insidePoints[1]
		outT2.Points[1] = outT1.Points[2]
		outT2.Points[2] = plane.Intersection(insidePoints[1], outsidePoints[0])

		return append(result, &outT1, &outT2)
	}

	return result
}

func (plane *Plane) Dist(p *Vector3d) float64 {
	// n := v.Normalize()
	return plane.N.X*p.X + plane.N.Y*p.Y + plane.N.Z*p.Z - plane.N.DotProduct(plane.P)
}
