/*
File: point_test.go
Created: Wed Jun 29 01:07:05 PDT 2011
*/

package point

import (
	"math"
	"testing"
)

var E = 1.0e-15
var e = 1.0e-12

func approx(t *testing.T, desc string, expect, value, epsilon float64) {
	if math.Abs(expect-value) > epsilon {
		if expect == 0 {
			t.Errorf("%s is not zero (%g)", desc, expect, value)
		} else {
			t.Errorf("%s is not %g (%g)", desc, expect, value)
		}
	}
}

func TestDotProduct(T *testing.T) {
	var (
		p1 = Point{-1.23456, 7.890987}
		p2 = Point{7.890987, 1.23456}
	)
	approx(T, "orthogonal dot product", p1.Dot(p2), 0, E)
	approx(T, "symmetric dot product", p2.Dot(p1), 0, E)
	approx(T, "parallel dot product", p1.Dot(p1), math.Pow(p1.Norm(), 2), e)
}

func TestNorm(T *testing.T) {
	var (
		e1 = Point{1, 0}
		e2 = Point{0, -1}
	)
	approx(T, "e1 norm", e1.Norm(), 1, E)
	approx(T, "-e2 norm", e2.Norm(), 1, E)
}

func TestSimpleRotation(T *testing.T) {
	var (
		x    = float64(43)
		y    = float64(-15)
		p    = Point{x, y}
		r1   = p.Rot(math.Pi / 2)
		dot1 = p.Dot(r1)
		r2   = p.Rot(math.Pi)
		dot2 = math.Abs(p.Dot(r2) + p.Dot(p))
	)
	if !p.ApproxEqual(p.Rot(0)) {
		T.Errorf("0 degree rotation changed point %v %v", p, p.Rot(0))
	}
	approx(T, "90° rotation dot product", dot1, 0, e)
	approx(T, "180° rotation dot product", dot2, 0, E)
}

func TestAdvRotation(T *testing.T) {
	var (
		x         = float64(43)
		y         = float64(-15)
		p         = Point{x, y}
		center    = Point{10, 10}
		pCentered = p.Sub(center)
		r1        = p.RotAround(math.Pi/2, center)
		dot1      = pCentered.Dot(r1.Sub(center))
		r2        = p.RotAround(math.Pi, center)
		dot2      = math.Abs(pCentered.Dot(r2.Sub(center)) + pCentered.Dot(pCentered))
	)
	if !p.ApproxEqual(p.RotAround(0, center)) {
		T.Errorf("0 degree rotation changed point %v %v", p, p.RotAround(0, center))
	}
	approx(T, "90° rotation dot product", dot1, 0, e)
	approx(T, "180° rotation dot product", dot2, 0, E)
}
