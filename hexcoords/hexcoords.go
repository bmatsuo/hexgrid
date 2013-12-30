/*
File: hexcoords.go
Created: Sat Jul  2 00:54:20 PDT 2011
*/

// A coordinate system for a hexogonal field.
package hexcoords

import (
	"github.com/bmatsuo/hexgrid/hex"
	"math"
	//"log"
)

//  Discrete hex coordinates consist of a horizontal U axis and a vertical
//  V axis. Each axis has range (-inf,inf) in theory. In practice, Grid
//  objects limit the accessible hex tiles.
type Hex struct{ U, V int }

// The coordinates have the same U and V fields.
func (coords Hex) Equals(other Hex) bool {
	return coords.U == other.U && coords.V == other.V
}

// Vertices farthest in the specified direction. If d is hex.NilDirection all
// coordinates vertices are returned, otherwise either one or two vertices
// are returned. The vertices are returned in increasing order (0, ..., 5).
func (c Hex) Vertices(d hex.Direction) []Vertex {
	var vertices []int
	if d < hex.NilDirection {
		vertices = directionVertices[d]
	} else {
		vertices = directionAllVertices
	}

	vcs := make([]Vertex, len(vertices))
	for i := range vertices {
		vcs[i] = Vertex{c.U, c.V, vertices[i]}
	}

	return vcs
}

var directionAllVertices = []int{0, 1, 2, 3, 4, 5}

var directionVertices = [][]int{
	hex.S: {
		hex.VertexIndex(hex.SW),
		hex.VertexIndex(hex.SE),
	},
	hex.SE: {
		hex.VertexIndex(hex.SE),
	},
	hex.E: {
		hex.VertexIndex(hex.SE),
		hex.VertexIndex(hex.NE),
	},
	hex.NE: {
		hex.VertexIndex(hex.NE),
	},
	hex.N: {
		hex.VertexIndex(hex.NE),
		hex.VertexIndex(hex.NW),
	},
	hex.NW: {
		hex.VertexIndex(hex.NW),
	},
	hex.W: {
		hex.VertexIndex(hex.NW),
		hex.VertexIndex(hex.SW),
	},
	hex.SW: {
		hex.VertexIndex(hex.SW),
	},
}

// See Vertices.
func (c Hex) Edges(d hex.Direction) []Edge {
	var coords []vertexPair
	if d < hex.NilDirection {
		coords = directionEdges[d]
	} else {
		coords = directionAllEdges
	}

	var edges = make([]Edge, 0, 6)
	for _, e := range coords {
		edges = append(edges, Edge{c.U, c.V, e.v1, e.v2})
	}
	return edges
}

type vertexPair struct{ v1, v2 int }

var directionAllEdges = []vertexPair{
	{0, 1},
	{1, 2},
	{2, 3},
	{3, 4},
	{4, 5},
	{5, 0},
}
var directionEdges = [][]vertexPair{
	hex.S: {
		{hex.VertexIndex(hex.SW), hex.VertexIndex(hex.SE)},
	},
	hex.SE: {
		{hex.VertexIndex(hex.SE), hex.VertexIndex(hex.E)},
	},
	hex.E: {
		{hex.VertexIndex(hex.SE), hex.VertexIndex(hex.E)},
		{hex.VertexIndex(hex.E), hex.VertexIndex(hex.NE)},
	},
	hex.NE: {
		{hex.VertexIndex(hex.E), hex.VertexIndex(hex.NE)},
	},
	hex.N: {
		{hex.VertexIndex(hex.NE), hex.VertexIndex(hex.NW)},
	},
	hex.NW: {
		{hex.VertexIndex(hex.NW), hex.VertexIndex(hex.W)},
	},
	hex.W: {
		{hex.VertexIndex(hex.NW), hex.VertexIndex(hex.W)},
		{hex.VertexIndex(hex.W), hex.VertexIndex(hex.SW)},
	},
	hex.SW: {
		{hex.VertexIndex(hex.W), hex.VertexIndex(hex.SW)},
	},
}

//  Vertices in the grid are indexed by hex coordinates paired with a
//  vertex index K. Vertex indices range from 0 to 5 and begin in the
//  south-west corner of the vertex. See also, hex.Direction.
type Vertex struct {
	U, V, K int
}

//  Create a Hex object from a Vertex object.
func (vc Vertex) Hex() Hex {
	return Hex{vc.U, vc.V}
}

//  Returns true if  only if the U, V and K fields of vc and other are equal.
func (vc Vertex) Equals(other Vertex) bool {
	return vc.U == other.U && vc.V == other.V && vc.K == other.K
}

//  Returns true if (u1,v1,k1) and (u2,v2,k2) reference the same vertex.
func (vc Vertex) IsIdentical(other Vertex) bool {
	var identVertices = vc.IdenticalVertices()
	if identVertices == nil {
		panic("nilident")
	}
	for _, ident := range identVertices {
		if ident.Equals(other) {
			return true
		}
	}
	return false
}
func (vc Vertex) Clockwise() Vertex {
	return Vertex{vc.U, vc.V, hex.VertexIndexClockwise(vc.K)}
}
func (vc Vertex) CounterClockwise() Vertex {
	return Vertex{vc.U, vc.V, hex.VertexIndexCounterClockwise(vc.K)}
}

//  Edges in the grid are index by hex coordinates along with a pair of
//  vertex indices K and L.
type Edge struct {
	U, V, K, L int
}

var (
	nilEdge = Edge{}
)

//  The zero value of Edge, which does not represent a real edge.
func NilEdge() Edge {
	return nilEdge
}
func (e Edge) Hex() Hex {
	return Hex{e.U, e.V}
}

//  Returns true if and only if e and e2 have exactly the same fields.
func (e Edge) Equals(e2 Edge) bool {
	return e.U == e2.U && e.V == e2.V && e.K == e2.K && e.L == e2.L
}
func (e Edge) reverse() Edge {
	return Edge{U: e.U, V: e.V, K: e.L, L: e.L}
}

//  Returns true if and only if e and e2 reference the same edge.
func (e Edge) IsIdentical(e2 Edge) bool {
	var (
		adjc  = e.Incidents()
		adjc2 = e2.Incidents()
	)
	if adjc[0].Equals(adjc2[1]) && adjc[1].Equals(adjc[0]) {
		return true
	} else if adjc[1].Equals(adjc2[0]) && adjc[0].Equals(adjc[1]) {
		return true
	}
	return false
}

//	Returns coordinates of edges sharing one endpoint with e.
func (e Edge) Adjacents() []Edge {
	if e.IsNil() {
		return nil
	}
	var (
		v1, v2 = e.Ends()
		adj    = make([]Edge, 0, 4)
	)
	for _, e1 := range v1.Edges() {
		if !e.IsIdentical(e1) {
			adj = append(adj, e1)
		}
	}
	for _, e2 := range v2.Edges() {
		if !e.IsIdentical(e2) {
			adj = append(adj, e2)
		}
	}
	return adj
}

//  Test if an edge is NilEdge. Synonymn for e.Equals(NilEdge())
func (e Edge) IsNil() bool {
	return e.Equals(nilEdge)
}

//  Retrieve the coordinates of e's incident vertices.
func (e Edge) Ends() (v1, v2 Vertex) {
	v1 = Vertex{e.U, e.V, e.K}
	v2 = Vertex{e.U, e.V, e.L}
	return v1, v2
}

func ColumnIsHigh(u int) bool {
	var uOdd = uint(math.Abs(float64(u)))%2 == 1
	return uOdd
}
func sameTile(u1, v1, u2, v2 int) bool {
	return Hex{u1, v1}.Equals(Hex{u2, v2})
}

//  If hex tiles (u1,v1) and (u2,v2) are adjacent, the direction of (u2,v2)
//  from (u1,v1) is returned. Otherwise hex.NilDirection is returned.
func (c Hex) Adjacency(adj Hex) hex.Direction {
	var (
		deltaU = adj.U - c.U
		deltaV = adj.U - c.U
	)
	if c.U == adj.U {
		if deltaV == 1 {
			return hex.N
		} else if deltaV == -1 {
			return hex.S
		}
		return hex.NilDirection
	}
	if deltaU == 1 {
		if ColumnIsHigh(c.U) {
			if c.U == adj.U {
				return hex.SE
			}
			if adj.U == c.U+1 {
				return hex.NE
			}
		} else {
			if adj.V == c.V-1 {
				return hex.SE
			}
			if c.V == adj.V {
				return hex.NE
			}
		}
	} else if deltaU == -1 {
		if ColumnIsHigh(c.U) {
			if adj.V == c.V+1 {
				return hex.NW
			}
			if c.V == adj.V {
				return hex.SW
			}
		} else {
			if c.V == adj.V {
				return hex.NW
			}
			if adj.V == c.V-1 {
				return hex.SW
			}
		}
	}
	return hex.NilDirection
}

//  Returns true if and only if c is adjacent to adj
func (c Hex) IsAdjacent(adj Hex) bool {
	return c.Adjacency(adj) != hex.NilDirection
}

//  Return a slice of the coordinates for adjacent hexagons
//  (not necessarily in the grid).
//  If E (or W) is supplied then the NE and SE (or NE and SW) coordinates
//  are returned in that order.
//  If NilDirection is suppied, then coordinates for all adjacent hexagons
//  are returned in the order N, NE, SE, S, SW, NW.
func (coords Hex) Adjacents(dir hex.Direction) []Hex {
	var (
		u = coords.U
		v = coords.V
	)
	switch dir {
	case hex.N:
		return []Hex{Hex{u, v + 1}}
	case hex.S:
		return []Hex{Hex{u, v - 1}}
	case hex.E:
		var adjE = make([]Hex, 2)
		if ColumnIsHigh(u) {
			adjE[0] = Hex{u - 1, v + 1}
			adjE[1] = Hex{u - 1, v}
		} else {
			adjE[0] = Hex{u - 1, v}
			adjE[1] = Hex{u - 1, v - 1}
		}
		return adjE
	case hex.W:
		var adjW = make([]Hex, 2)
		if ColumnIsHigh(u) {
			adjW[0] = Hex{u + 1, v + 1}
			adjW[1] = Hex{u + 1, v}
		} else {
			adjW[0] = Hex{u + 1, v}
			adjW[1] = Hex{u + 1, v - 1}
		}
		return adjW
	case hex.NE:
		if ColumnIsHigh(u) {
			return []Hex{Hex{u - 1, v + 1}}
		}
		return []Hex{Hex{u - 1, v}}
	case hex.NW:
		if ColumnIsHigh(u) {
			return []Hex{Hex{u + 1, v + 1}}
		}
		return []Hex{Hex{u + 1, v}}
	case hex.SE:
		if ColumnIsHigh(u) {
			return []Hex{Hex{u - 1, v}}
		}
		return []Hex{Hex{u - 1, v - 1}}
	case hex.SW:
		if ColumnIsHigh(u) {
			return []Hex{Hex{u + 1, v}}
		}
		return []Hex{Hex{u + 1, v - 1}}
	default:
		var adjAll = make([]Hex, 6)
		if ColumnIsHigh(u) {
			adjAll[0] = Hex{u, v + 1}     // North
			adjAll[1] = Hex{u - 1, v + 1} // NorthEast
			adjAll[2] = Hex{u - 1, v}     // SouthEast
			adjAll[3] = Hex{u, v - 1}     // South
			adjAll[4] = Hex{u + 1, v}     // SouthWest
			adjAll[5] = Hex{u + 1, v + 1} // NorthWest
		} else {
			adjAll[0] = Hex{u, v + 1}
			adjAll[1] = Hex{u - 1, v}
			adjAll[2] = Hex{u - 1, v - 1}
			adjAll[3] = Hex{u, v - 1}
			adjAll[4] = Hex{u + 1, v - 1}
			adjAll[5] = Hex{u + 1, v}
		}
		return adjAll
	}
	return nil
}

func (vert Vertex) Incidents() []Hex {
	var (
		adjVC = vert.IdenticalVertices()
		adj   = make([]Hex, 0, len(adjVC))
	)
	for _, coords := range adjVC {
		var (
			adjHex = Hex{coords.U, coords.V}
		)
		adj = append(adj, adjHex)
	}

	return adj
}

//  The Hex of tiles which share edge e. There are exactly two such
//  Hex for any (real) edge. Returns nil if e is NilEdge()
func (e Edge) Incidents() []Hex {
	if e.IsNil() {
		return nil
	}
	vert1, vert2 := e.Ends()
	return vert1.HexShared(vert2)
}

func (vert Vertex) HexShared(vert2 Vertex) []Hex {
	var (
		adjC1 = vert.Incidents()
		adjC2 = vert2.Incidents()
	)
	var shared = make([]Hex, 0, 1)
	for _, c1 := range adjC1 {
		for _, c2 := range adjC2 {
			if c1.U == c2.U && c1.V == c2.V {
				shared = append(shared, c1)
			}
		}
	}
	return shared
}

func (vert Vertex) EdgeShared(vert2 Vertex) Edge {
	if vert.IsIdentical(vert2) {
		return Edge{}
	}
	if vert.Hex().Equals(vert2.Hex()) {
		return Edge{vert.U, vert.V, vert.K, vert2.K}
	}
	var (
		identVerts1 = vert.IdenticalVertices()
		identVerts2 = vert2.IdenticalVertices()
	)
	if identVerts1 == nil || identVerts2 == nil {
		panic("nilident")
	}
	for _, ident1 := range identVerts1 {
		for _, ident2 := range identVerts2 {
			if ident1.Hex().Equals(ident2.Hex()) {
				return Edge{ident1.U, ident1.V, ident1.K, ident2.K}
			}
			var ec = ident1.Hex().EdgeShared(ident2.Hex())
			if ec.IsNil() {
				return ec
			}
		}
	}
	return Edge{}
}

//  Function for determining the vertex indices of an edge in
//  the hex tile at (u1,v1) that is alse in tile (u2,v2).
//  Returns nil if the hex coordinates are not adjacent.
func (coord Hex) EdgeShared(other Hex) Edge {
	var adjDir = coord.Adjacency(other)
	if adjDir == hex.NilDirection {
		return nilEdge
	}
	edge := adjDir.Edge()
	return Edge{coord.U, coord.V, edge.Orig().Int(), edge.Term().Int()}
}

//  This method needs testing.
func (vert Vertex) Edges() []Edge {
	var (
		adjVCs = vert.Adjacents()
		edges  = make([]Edge, len(adjVCs))
	)
	for i, other := range adjVCs {
		edges[i] = vert.EdgeShared(other)
	}
	return edges
}

//  This is untested.
func (vert Vertex) AdjacentByEdge(edge Edge) Vertex {
	v1, v2 := edge.Ends()
	if vert.IsIdentical(v1) {
		return v2
	} else if vert.IsIdentical(v2) {
		return v1
	}
	return Vertex{}
}

//  Get coordinates of hex vertices in the field incident to vertex (u,v,k).
//  Returns a slice of vertex coordinates (slices of 3 ints), the first of
//  which being []int{u, v, k}. See also, (vc Vertex) IsIdentical.
func (vert Vertex) IdenticalVertices() []Vertex {
	var adjC = make([]Vertex, 1, 3)
	adjC[0] = vert

	var adjOffsets [][]int
	if ColumnIsHigh(vert.U) {
		adjOffsets = hexHighVertexIncidenceOffset[vert.K]
	} else {
		adjOffsets = hexLowVertexIncidenceOffset[vert.K]
	}
	for _, offset := range adjOffsets {
		var (
			du   = offset[0]
			dv   = offset[1]
			kAdj = offset[2]
			uNew = vert.U + du
			vNew = vert.V + dv
		)
		var offsetCp = Vertex{uNew, vNew, kAdj}
		adjC = append(adjC, offsetCp)
	}

	return adjC
}

//  Get a list of unique vertices adjacent to (u,v,k).
//  See also, (Vertex) Identical.
func (vert Vertex) Adjacents() []Vertex {
	var identVerts = vert.IdenticalVertices()
	var adjVerts = make([]Vertex, len(identVerts))
	for i, vert := range identVerts {
		adjVerts[i] = Vertex{vert.U, vert.V, hex.VertexIndexClockwise(vert.K)}
	}
	return adjVerts
}

func (vert Vertex) IsAdjacent(other Vertex) bool {
	for _, adj := range vert.Adjacents() {
		if adj.IsIdentical(other) {
			return true
		}
	}
	return false
}

var (
	hexHighVertexIncidenceOffset = [][][]int{
		{{-1, 0, 2}, {0, -1, 4}},
		{{0, -1, 3}, {1, 0, 5}},
		{{1, 0, 4}, {1, 1, 0}},
		{{1, 1, 5}, {0, 1, 1}},
		{{0, 1, 0}, {-1, 1, 2}},
		{{-1, 1, 2}, {-1, 0, 3}}}
	hexLowVertexIncidenceOffset = [][][]int{
		{{-1, -1, 2}, {0, -1, 4}},
		{{0, -1, 3}, {1, -1, 5}},
		{{1, -1, 4}, {1, 0, 0}},
		{{1, 0, 5}, {0, 1, 1}},
		{{0, 1, 0}, {-1, 0, 2}},
		{{-1, 0, 1}, {-1, -1, 3}}}
)
