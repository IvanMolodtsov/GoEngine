package engine

type Triangle struct {
	Points [3]Vector3d
}

func (t Triangle) Render(renderer *Renderer) {
	renderer.DrawLine(t.Points[0].X, t.Points[0].Y, t.Points[1].X, t.Points[1].Y, uint32(0xff00ff00))
	renderer.DrawLine(t.Points[1].X, t.Points[1].Y, t.Points[2].X, t.Points[2].Y, uint32(0xff00ff00))
	renderer.DrawLine(t.Points[2].X, t.Points[2].Y, t.Points[0].X, t.Points[0].Y, uint32(0xff00ff00))
}

type Mesh struct {
	Tris []Triangle
}

func (m Mesh) Render(renderer *Renderer) {
	for _, t := range m.Tris {
		var triProjected Triangle

		t.Points[0].Z = t.Points[0].Z + 5.0
		t.Points[1].Z = t.Points[1].Z + 5.0
		t.Points[2].Z = t.Points[2].Z + 5.0

		println(t.Points[0].X, t.Points[0].Y, t.Points[0].Z)

		// triProjected.Points[0] = t.Points[0].Multiply(renderer.ProjectionMatrix)
		// triProjected.Points[1] = t.Points[1].Multiply(renderer.ProjectionMatrix)
		// triProjected.Points[2] = t.Points[2].Multiply(renderer.ProjectionMatrix)

		//scale into view
		triProjected.Points[0].X += 1
		triProjected.Points[0].Y += 1
		triProjected.Points[1].X += 1
		triProjected.Points[1].Y += 1
		triProjected.Points[2].X += 1
		triProjected.Points[2].Y += 1

		triProjected.Points[0].X *= 0.5 * float64(renderer.screenWidth)
		triProjected.Points[0].Y *= 0.5 * float64(renderer.screenHeight)
		triProjected.Points[1].X *= 0.5 * float64(renderer.screenWidth)
		triProjected.Points[1].Y *= 0.5 * float64(renderer.screenHeight)
		triProjected.Points[2].X *= 0.5 * float64(renderer.screenWidth)
		triProjected.Points[2].Y *= 0.5 * float64(renderer.screenHeight)

		triProjected.Render(renderer)
	}
}
