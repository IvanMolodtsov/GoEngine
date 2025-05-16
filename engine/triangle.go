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

func (t Triangle) Fill(r *Renderer) {
	x1 := int64(t.Points[0].X)
	y1 := int64(t.Points[0].Y)
	x2 := int64(t.Points[1].X)
	y2 := int64(t.Points[1].Y)
	x3 := int64(t.Points[2].X)
	y3 := int64(t.Points[2].Y)

	SWAP := func(x, y *int64) {
		t := *x
		*x = *y
		*y = t
	}

	DRAWLINE := func(sx, ex, ny int64) {
		for i := sx; i <= ex; i++ {
			r.DrawPixel(float64(i), float64(ny), toHex(t.Color))
		}
	}

	changed1, changed2 := false, false

	// Sort vertices
	if y1 > y2 {
		SWAP(&y1, &y2)
		SWAP(&x1, &x2)
	}
	if y1 > y3 {
		SWAP(&y1, &y3)
		SWAP(&x1, &x3)
	}
	if y2 > y3 {
		SWAP(&y2, &y3)
		SWAP(&x2, &x3)
	}

	// Starting points
	t1x, t2x := x1, x1
	y := y1

	var signx1 int64
	var dx1 int64 = x2 - x1
	if dx1 < 0 {
		dx1 = -dx1
		signx1 = -1
	} else {
		signx1 = 1
	}

	dy1 := y2 - y1

	var signx2 int64
	dx2 := x3 - x1
	if dx2 < 0 {
		dx2 = -dx2
		signx2 = -1
	} else {
		signx2 = 1
	}

	dy2 := y3 - y1

	if dy1 > dx1 { // swap values
		SWAP(&dx1, &dy1)
		changed1 = true
	}
	if dy2 > dx2 { // swap values
		SWAP(&dy2, &dx2)
		changed2 = true
	}
	var minx, maxx, t1xp, t2xp, e1, i int64

	e2 := dx2 >> 1
	// Flat top, just process the second half
	if y1 == y2 {
		goto NEXT
	}
	e1 = dx1 >> 1

	for i = 0; i < dx1; {
		t1xp = 0
		t2xp = 0
		if t1x < t2x {
			minx = t1x
			maxx = t2x
		} else {
			minx = t2x
			maxx = t1x
		}
		// process first line until y value is about to change
		for i < dx1 {
			i++
			e1 += dy1
			for e1 >= dx1 {
				e1 -= dx1
				if changed1 { //t1x += signx1;
					t1xp = signx1
				} else {
					goto next1
				}
			}
			if changed1 {
				break
			} else {
				t1x += signx1
			}
		}
		// Move line
	next1:
		// process second line until y value is about to change
		for {
			e2 += dy2
			for e2 >= dx2 {
				e2 -= dx2
				if changed2 {
					//t2x += signx2;
					t2xp = signx2
				} else {
					goto next2
				}
			}
			if changed2 {
				break
			} else {
				t2x += signx2
			}
		}
	next2:
		if minx > t1x {
			minx = t1x
			if minx > t2x {
				minx = t2x
			}
		}
		if maxx < t1x {
			maxx = t1x
			if maxx < t2x {
				maxx = t2x
			}
		}
		DRAWLINE(minx, maxx, y) // Draw line from min to max points found on the y
		// Now increase y
		if !changed1 {
			t1x += signx1
		}
		t1x += t1xp
		if !changed2 {
			t2x += signx2
		}
		t2x += t2xp
		y += 1
		if y == y2 {
			break
		}

	}
NEXT:
	// Second half
	dx1 = x3 - x2
	if dx1 < 0 {
		dx1 = -dx1
		signx1 = -1
	} else {
		signx1 = 1
	}
	dy1 = (y3 - y2)
	t1x = x2

	if dy1 > dx1 { // swap values
		SWAP(&dy1, &dx1)
		changed1 = true
	} else {
		changed1 = false
	}

	e1 = dx1 >> 1

	for i = 0; i <= dx1; i++ {
		t1xp = 0
		t2xp = 0
		if t1x < t2x {
			minx = t1x
			maxx = t2x
		} else {
			minx = t2x
			maxx = t1x
		}
		// process first line until y value is about to change
		for i < dx1 {
			e1 += dy1
			for e1 >= dx1 {
				e1 -= dx1
				if changed1 {
					t1xp = signx1
					break
					//t1x += signx1;
				} else {
					goto next3
				}
			}
			if changed1 {
				break
			} else {
				t1x += signx1
			}
			if i < dx1 {
				i++
			}
		}
	next3:
		// process second line until y value is about to change
		for t2x != x3 {
			e2 += dy2
			for e2 >= dx2 {
				e2 -= dx2
				if changed2 {
					t2xp = signx2
				} else {
					goto next4
				}
			}
			if changed2 {
				break
			} else {
				t2x += signx2
			}
		}
	next4:

		if minx > t1x {
			minx = t1x
			if minx > t2x {
				minx = t2x
			}
		}
		if maxx < t1x {
			maxx = t1x
			if maxx < t2x {
				maxx = t2x
			}
		}
		DRAWLINE(minx, maxx, y)
		if !changed1 {
			t1x += signx1
		}
		t1x += t1xp
		if !changed2 {
			t2x += signx2
		}
		t2x += t2xp
		y += 1
		if y > y3 {
			return
		}
	}
}
