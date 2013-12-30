package hex

import (
	"math"
)

//  Hexagons have faces in directions NW, N, NE, SE, S, SW
//  and vertices in directions W, NW, NE, E, SE, SW.
type Direction int

const (
	N Direction = iota
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
	hexDirections    = []Direction{S, SE, E, NE, N, NW, W, SW}
	vertexDirections = []Direction{SE, E, NE, NW, W, SW}
	edgeDirections   = []Direction{S, SE, NE, N, NW, SW}
)

func copyDirections(ds []Direction) []Direction {
	var dsCopy = make([]Direction, len(ds))
	copy(dsCopy, ds)
	return dsCopy
}
func Directions() []Direction {
	return copyDirections(hexDirections)
}
func VertexDirections() []Direction {
	return copyDirections(vertexDirections)
}
func EdgeDirections() []Direction {
	return copyDirections(edgeDirections)
}

var hexDirectionInverse = []Direction{
	N:  S,
	NE: SW,
	E:  W,
	SE: NW,
	S:  N,
	SW: NE,
	W:  E,
	NW: SE,
}

func (dir Direction) Inverse() Direction {
	if int(dir) >= len(hexDirectionInverse) {
		return NilDirection
	}
	return hexDirectionInverse[dir]
}

type Vertex int

//  Return the vertex k in direction dir from a hex tile's center.
//  Returns -1 if dir is NilDirection, N, or S.
func (dir Direction) Vertex() Vertex {
	switch dir {
	case SW:
		return V_SW
	case SE:
		return V_SE
	case E:
		return V_E
	case NE:
		return V_NE
	case NW:
		return V_NW
	case W:
		return V_W
	}
	return -1
}

const (
	V_SW Vertex = iota
	V_SE
	V_E
	V_NE
	V_NW
	V_W
	_V_INVALID
)

var vertexDirection = []Direction{
	V_SW: SW,
	V_SE: SE,
	V_E:  E,
	V_NE: NE,
	V_NW: NW,
	V_W:  W,
}

func (v Vertex) Clockwise() Vertex {
	return (v + 1) % _V_INVALID
}

func (v Vertex) CounterClockwise() Vertex {
	return (v + 5) % _V_INVALID
}

func (v Vertex) Direction() Direction {
	if 0 <= v && v < _V_INVALID {
		return vertexDirection[v]
	}
	return NilDirection
}

func (v Vertex) Int() int {
	return int(v)
}

// An edge starts at a vertex and terminate at the counter-clockwise vertex
type Edge Vertex

//  Return the edge k in direction dir from a hex tile's center.
//  Returns -1 if dir is NilDirection, N, or S.
func (dir Direction) Edge() Edge {
	switch dir {
	case S:
		return E_S
	case SE:
		return E_SE
	case NE:
		return E_NE
	case N:
		return E_N
	case NW:
		return E_NW
	case SW:
		return E_SW
	}
	return -1
}

const (
	E_S  = Edge(V_SW)
	E_SE = Edge(V_SE)
	E_NE = Edge(V_E)
	E_N  = Edge(V_NE)
	E_NW = Edge(V_NW)
	E_SW = Edge(V_W)
	_E_INVALID
)

var edgeDirection = []Direction{
	E_S:  S,
	E_SE: SE,
	E_NE: NE,
	E_N:  N,
	E_NW: NW,
	E_SW: SW,
}

func (e Edge) Clockwise() Edge {
	return Edge(Vertex(e).Clockwise())
}

func (e Edge) CounterClockwise() Edge {
	return Edge(Vertex(e).CounterClockwise())
}

func (e Edge) Orig() Vertex {
	return Vertex(e)
}

func (e Edge) Term() Vertex {
	return e.Orig().CounterClockwise()
}

func (e Edge) Direction() Direction {
	if 0 <= e && e < _E_INVALID {
		return edgeDirection[e]
	}
	return NilDirection
}

//  Get the index of the vertex clockwise of vertex k.
func VertexIndexClockwise(k int) int {
	return (k + 5) % 6
}

//  Get the index of the vertex counter-clockwise of vertex k.
func VertexIndexCounterClockwise(k int) int {
	return (k + 1) % 6
}

//  Return the direction of vertex k relative to the center of a hexagon.
//  Returns NilDirection if k is not in the range [0,5].
func VertexDirection(k int) Direction {
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

//	DEPRECATED
//  Return the vertex k in direction dir from a hex tile's center.
//  Returns -1 if dir is NilDirection, N, or S.
func VertexIndex(dir Direction) int {
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
	TriangleAngle = math.Pi / 6
	RotateAngle   = math.Pi / 3
)

var (
	SideRadiusRatio = math.Tan(TriangleAngle)
)
