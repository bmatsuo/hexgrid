/*
*  File: hexagon.go
*  Author: Bryan Matsuo [bryan.matsuo@gmail.com]
*  Created: Wed Jun 29 13:56:22 PDT 2011
 */

package hexgrid

import (
	point "github.com/bmatsuo/hexgrid/point"

	"math"
	//"log"
)

//  Hexagons have faces in directions NW, N, NE, SE, S, SW
//  and vertices in directions W, NW, NE, E, SE, SW.
type HexDirection int

const (
	N HexDirection = iota
	NE
	NW
	S
	SE
	SW
	E
	W
	NilDirection
)

var (
	hexDirections    = []HexDirection{S, SE, E, NE, N, NW, W, SW}
	vertexDirections = []HexDirection{SE, E, NE, NW, W, SW}
	edgeDirections   = []HexDirection{S, SE, NE, N, NW, SW}
)

func copyDirections(ds []HexDirection) []HexDirection {
	var dsCopy = make([]HexDirection, len(ds))
	copy(dsCopy, ds)
	return dsCopy
}
func HexDirections() []HexDirection {
	return copyDirections(hexDirections)
}
func VertexDirections() []HexDirection {
	return copyDirections(vertexDirections)
}
func EdgeDirections() []HexDirection {
	return copyDirections(edgeDirections)
}

var hexDirectionInverse = []HexDirection{
	N:  S,
	NE: SW,
	E:  W,
	SE: NW,
	S:  N,
	SW: NE,
	W:  E,
	NW: SE,
}

func (dir HexDirection) Inverse() HexDirection {
	if int(dir) >= len(hexDirectionInverse) {
		return NilDirection
	}
	return hexDirectionInverse[dir]
}

//  Get the index of the vertex clockwise of vertex k.
func HexVertexIndexClockwise(k int) int {
	return (k + 5) % 6
}

//  Get the index of the vertex counter-clockwise of vertex k.
func HexVertexIndexCounterClockwise(k int) int {
	return (k + 1) % 6
}

//  Return the direction of vertex k relative to the center of a hexagon.
//  Returns NilDirection if k is not in the range [0,5].
func HexVertexDirection(k int) HexDirection {
	switch k {
	case 0:
		return SW
	case 1:
		return SE
	case 2:
		return E
	case 3:
		return NE
	case 4:
		return NW
	case 5:
		return W
	}
	return NilDirection
}

//  Return the vertex k in direction dir from a hex tile's center.
//  Returns -1 if dir is NilDirection, N, or S.
func HexVertexIndex(dir HexDirection) int {
	switch dir {
	case SW:
		return 0
	case SE:
		return 1
	case E:
		return 2
	case NE:
		return 3
	case NW:
		return 4
	case W:
		return 5
	}
	return -1
}

const (
	hexTriangleAngle = math.Pi / 6
	hexRotateAngle   = math.Pi / 3
)

var (
	hexSideRadiusRatio = math.Tan(hexTriangleAngle)
)

//  A simple hexagon type thinly wrapping a Point array.
type HexPoints [6]point.Point

func (hex *HexPoints) Point(k int) point.Point {
	if k < 0 || k >= len(hex) {
		panic("Point index out of bounds")
	}
	return hex[k]
}
func (hex *HexPoints) Points() []point.Point {
	var points = make([]point.Point, 6)
	copy(points, hex[:])
	return points
}
func HexEdgeDirection(k, ell int) HexDirection {
	if k > ell {
		var tmp = k
		k = ell
		ell = tmp
	}
	if k == 0 && ell == 1 {
		return S
	} else if k == 1 && ell == 2 {
		return SE
	} else if k == 2 && ell == 3 {
		return NE
	} else if k == 3 && ell == 4 {
		return N
	} else if k == 4 && ell == 5 {
		return NW
	} else if k == 0 && ell == 5 {
		return SW
	}
	return NilDirection
}
func HexEdgeIndices(dir HexDirection) []int {
	switch dir {
	case S:
		return []int{0, 1}
	case SE:
		return []int{1, 2}
	case NE:
		return []int{2, 3}
	case N:
		return []int{3, 4}
	case NW:
		return []int{4, 5}
	case SW:
		return []int{5, 0}
	}
	return nil
}
func (hex *HexPoints) EdgePoints(dir HexDirection) []point.Point {
	var edgeIndices = HexEdgeIndices(dir)
	if edgeIndices == nil {
		return nil
	}
	var (
		p1 = hex[edgeIndices[0]]
		p2 = hex[edgeIndices[1]]
	)
	return []point.Point{p1, p2}
}

//  Generate a hexagon at a given point.
func NewHex(p point.Point, r float64) *HexPoints {
	var (
		hex  = new(HexPoints)
		side = point.Point{r, 0}.Scale(hexSideRadiusRatio)
	)
	hex[0] = point.Point{0, -r}.Sub(side)
	for i := 1; i < 6; i++ {
		hex[i] = hex[i-1].Rot(hexRotateAngle)
	}
	for i := 0; i < 6; i++ {
		hex[i] = hex[i].Add(p)
	}
	return hex
}
