/* 
File: point_test.go
Created: Wed Jun 29 01:07:05 PDT 2011
*/

package point

import (
    "testing"
    "math"
)

func TestDotProduct(T *testing.T) {
    var (
        p1  = Point{-1.23456, 7.890987}
        p2  = Point{7.890987, 1.23456}
    )
    if math.Abs(p1.Dot(p2)-p2.Dot(p1)) > 1.0e-12 {
        T.Errorf("Symmetric dot product failure %e", math.Abs(p1.Dot(p2)-p2.Dot(p1)))
    }
    if math.Abs(p1.Dot(p2)) > 1.0e-12 {
        T.Errorf("Symmetric dot product failure %e", math.Abs(p1.Dot(p2)-p2.Dot(p1)))
    }
}

func TestNorm(T *testing.T) {
    var (
        e1  = Point{1, 0}
        e2  = Point{0, -1}
    )
    if math.Abs(e1.Norm()-1) > 1.0e-15 {
        T.Errorf("Unit vector e1 norm is not one %f", e1.Norm())
    }
    if math.Abs(e2.Norm()-1) > 1.0e-15 {
        T.Errorf("Unit vector -e2 norm is not one %f", e2.Norm())
    }
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
    if dot1 > 1.0e-12 {
        T.Errorf("Orthogonal Dot product %e", dot1)
    }
    if dot2 > 1.0e-12 {
        T.Errorf("180 degree Dot product %e", dot2)
    }
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
    if dot1 > 1.0e-12 {
        T.Errorf("Orthogonal Dot product %e", dot1)
    }
    if dot2 > 1.0e-12 {
        T.Errorf("180 degree Dot product %e", dot2)
    }
}
