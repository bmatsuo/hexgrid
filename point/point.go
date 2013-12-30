/*
File: point.go
Created: Wed Jun 29 01:07:05 PDT 2011
*/

package point

import (
	"image"
	"math"
)

var PointApproximationGap = 1.0e-12

//  A generic 2-dimensional point.
type Point struct{ X, Y float64 }

func E0() Point {
	return Point{1, 0}
}

func E1() Point {
	return Point{0, 1}
}

// The point (0, 0)
func Zero() Point {
	return Point{}
}

//  A point at (infinity, infinity). See also, PointIsInf.
func Inf() Point {
	var inf = math.Inf(1)
	return Point{inf, inf}
}

//  Test a point to see if it is infinity. A point is infinity if either
//  of its components are infinity (positive or negative).
func (p Point) IsInf() bool {
	if math.IsInf(p.X, 0) || math.IsInf(p.Y, 0) {
		return true
	}
	return false
}
func (p Point) ApproxEqual(p2 Point) bool {
	return p.Sub(p2).Norm() < PointApproximationGap
}
func (p Point) ImageCoords(rect image.Rectangle) Point {
	return Point{p.X, float64(rect.Max.Y) - p.Y}
}
func (p Point) ImagePoint(rect image.Rectangle) image.Point {
	var ip = p.ImageCoords(rect)
	return image.Point{int(ip.X), int(ip.Y)}
}
func (p Point) Norm() float64 {
	return math.Sqrt(p.Dot(p))
}
func (p Point) Dot(p2 Point) float64 {
	return p.X*p2.X + p.Y*p2.Y
}
func (p Point) Scale(a float64) Point {
	return Point{a * p.X, a * p.Y}
}
func (p Point) Sub(p2 Point) Point {
	return Point{p.X - p2.X, p.Y - p2.Y}
}
func (p Point) Add(p2 Point) Point {
	return Point{p.X + p2.X, p.Y + p2.Y}
}
func (p Point) Rot(theta float64) Point {
	return Point{
		p.X*math.Cos(theta) - p.Y*math.Sin(theta),
		p.Y*math.Cos(theta) + p.X*math.Sin(theta),
	}
}
func (p Point) RotAround(theta float64, center Point) Point {
	var (
		mv     = p.Sub(center)
		rot    = mv.Rot(theta)
		mvBack = rot.Add(center)
	)
	return mvBack
}
