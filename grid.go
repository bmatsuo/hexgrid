/*
*  File: grid.go
*  Author: Bryan Matsuo [bryan.matsuo@gmail.com]
*  Created: Tue Jun 28 03:40:52 PDT 2011
 */

package hexgrid

import (
	"github.com/bmatsuo/hexgrid/hex"
	"github.com/bmatsuo/hexgrid/point"
	"github.com/bmatsuo/hexgrid/hexcoords"

	"fmt"
	"math"
	//"log"
)

type Value interface{}

//  For each coordinate in a Grid there is one unique HexTile.
type Tile struct {
	Hex hexcoords.Hex
	Pos    point.Point
	Value  Value
}
type TileInitializer func(hexcoords.Hex) Value

//  A HexVertex represents the corner a HexTile. A HexVertex can be shared
//  by at most 3 HexTiles and can be the junction of between 2 and three
//  HexEdge objects. A HexVertex can be 'between' fewer than three tiles if
//  its tiles are on the edge of the grid. It will be the endpoint of two
//  edges only it belongs to one tile (and is on the edge of the grid).
type Vertex struct {
	Hex hexcoords.Vertex
	Pos    point.Point
	Value  Value
}
type VertexInitializer func(hexcoords.Vertex) Value

//  A HexEdge represents an edge between two HexVertex objects. It is
//  part of the boundary of a HexTile. A HexEdge can be 'Between' only
//  one tile if its tile is on the edge of the grid.
type Edge struct {
	Hex hexcoords.Edge
	Value  Value
}
type EdgeInitializer func(coords hexcoords.Edge, v1, v2 *Vertex) Value

//  A grid of hexagons in a discrete coordinate system (u,v) where u
//  indexes the column of the grid, and v the row.
type Grid struct {
	radius   float64
	n        int
	m        int
	p        []point.Point
	v        []Vertex
	e        []Edge
	t        []Tile
	hexes    [][]*hex.HexPoints
	tiles    [][]*Tile
	vertices [][][]*Vertex
	edges    [][][][]*Edge
}

//  Create an nxm grid of hexagons with radius r. Where n is the number of
//  columns and m is the number of rows. The integers n and m must be odd.
//  The *Default arguments dictate the initialized Value field of each Tile,
//  Vertex and Edge object. If the value of a default is a function taking
//  the proper arguments and returning a Value object then that function is
//  called to generate each objects initial value. See also, TileInitializer,
//  VertexInitializer, and EdgeInitializer.
func NewGrid(n, m int, r float64, tileDefault, vertexDefault, edgeDefault interface{}) *Grid {
	if n&1 == 0 || m&1 == 0 {
		panic("evensize")
	}
	if n < 0 || n < 0 {
		panic("negsize")
	}
	if r < 0 {
		panic("negradius")
	}
	var h = new(Grid)
	h.radius = r
	h.n = n
	h.m = m
	h.genHexagons()
	h.genTiles(tileDefault)
	h.genVertices(vertexDefault)
	h.genEdges(edgeDefault) // Must come after genVertices.
	return h
}

//  Retrieve a Tile object specified by its coordinates.
func (h *Grid) GetTile(c hexcoords.Hex) *Tile {
	if !h.WithinBounds(c) {
		return nil
	}
	i, j := h.hexIndex(c)
	return h.tiles[i][j]
}

//  Retrieve a Vertex object specified by its coordinates.
func (h *Grid) GetVertex(vert hexcoords.Vertex) *Vertex {
	var inbounds = h.getVCWithinBounds(vert)
	if !h.WithinBounds(inbounds.Hex()) {
		return nil
	}
	i, j := h.hexIndex(inbounds.Hex())
	return h.vertices[i][j][inbounds.K]
}

//  Retrieve an Edge object specified by its coordinates.
func (h *Grid) GetEdge(e hexcoords.Edge) *Edge {
	var c = e.Hex()
	if !h.WithinBounds(c) {
		return nil
	}
	i, j := h.hexIndex(c)
	return h.edges[i][j][e.K][e.L]
}
func (h *Grid) GetEdges(coords hexcoords.Hex) []*Edge {
	if !h.WithinBounds(coords) {
		return nil
	}
	var edges = make([]*Edge, 6)
	for k, ec := range coords.Edges(hex.NilDirection) {
		edges[k] = h.GetEdge(ec)
	}
	return edges
}

//  Returns the width and height of the Grid wrapped in a
//  GridDimensions object.
func (h *Grid) Size() hexcoords.Hex {
	return hexcoords.Hex{h.n, h.m}
}

//  Total number of distinct hexagon vertices in the field.
func (h *Grid) expectedNumVertices() int {
	return 2 * (h.n*h.m + h.n + h.m)
}
func (h *Grid) NumVertices() int {
	return len(h.v)
}
func (h *Grid) expectedNumEdges() int {
	return 3*h.n*h.m + 2*h.n + 2*h.m - 1
}
func (h *Grid) NumEdges() int {
	return len(h.e)
}
func (h *Grid) expectedNumTiles() int {
	return h.n * h.m
}

//  Number of hex tiles in the field (n^2).
func (h *Grid) NumTiles() int {
	return len(h.t)
}
func (h *Grid) NumCols() int {
	return h.n
}
func (h *Grid) NumRows() int {
	return h.m
}
func (h *Grid) horizontalIndexOffset() int {
	return int(math.Floor(float64(h.NumCols()) / 2))
}
func (h *Grid) verticalIndexOffset() int {
	return int(math.Floor(float64(h.NumRows()) / 2))
}

//  Minimum value of the row coordinate v.
func (h *Grid) RowMin() int { return -h.verticalIndexOffset() }

//  Maximum value of the row coordinate v.
func (h *Grid) RowMax() int { return h.verticalIndexOffset() }

//  Minimum value of the column coordinate u.
func (h *Grid) ColMin() int { return -h.horizontalIndexOffset() }

//  Maximum value of the column coordinate u.
func (h *Grid) ColMax() int { return h.horizontalIndexOffset() }

/* Some coordinate <-> index internal methods. */
func (h *Grid) hexCoords(i, j int) hexcoords.Hex {
	return hexcoords.Hex{i + h.ColMin(), j + h.RowMin()}
}
func (h *Grid) hexIndex(c hexcoords.Hex) (int, int) {
	return c.U - h.ColMin(), c.V - h.RowMin()
}

/* Internal bounds checking method. */
func (h *Grid) indexWithinBounds(i, j int) bool {
	return h.WithinBounds(h.hexCoords(i, j))
}

//  Returns true if the hex at coordinates (u,v) is in the hex field.
func (h *Grid) WithinBounds(c hexcoords.Hex) bool {
	if c.U < h.ColMin() || c.U > h.ColMax() {
		return false
	} else if c.V < h.RowMin() || c.V > h.RowMax() {
		return false
	}
	return true
}

//  Generate points for the hexagon at row i, column j.
//  Returns nil when the position (i,j) is not within the bounds of the board.
func (h *Grid) GetHex(c hexcoords.Hex) *hex.HexPoints {
	if !h.WithinBounds(c) {
		return nil
	}
	var (
		i, j = h.hexIndex(c)
	)
	if h.hexes[i][j] != nil {
		var newh = new(hex.HexPoints)
		*newh = *(h.hexes[i][j])
		return newh
	}
	return hex.NewHex(h.TileCenter(c), h.radius)
}

func (h *Grid) getVCWithinBounds(vc hexcoords.Vertex) hexcoords.Vertex {
	if h.WithinBounds(vc.Hex()) {
		return vc
	}
	var idents = vc.IdenticalVertices()
	for _, id := range idents[1:] {
		if h.WithinBounds(id.Hex()) {
			return id
		}
	}
	return vc
}

//  Get a pointer to the kth corner point of the hex tile at (u,v).
//  Returns point.Inf() when no vertex identical to vc is within the
//  bounds of h.
func (h *Grid) GetVertexPoint(vc hexcoords.Vertex) point.Point {
	var inbounds = h.getVCWithinBounds(vc)
	var hex = h.GetHex(inbounds.Hex())
	if hex == nil {
		return point.Inf()
	}
	return hex[inbounds.K]
}

//  This methods should be replaced.
func (h *Grid) GetVertices(coords hexcoords.Hex) []*Vertex {
	if !h.WithinBounds(coords) {
		return nil
	}
	var vertices = make([]*Vertex, 6)
	for _, v := range coords.Vertices(hex.NilDirection) {
		vertices[v.K] = h.GetVertex(v)
	}
	return vertices
}

/* Internal methods for computing hexagon positions. */
func (h *Grid) horizontalSpacing() float64 {
	return 2 * h.radius * math.Cos(hex.TriangleAngle)
}
func (h *Grid) verticalSpacing() float64 {
	return 2 * h.radius
}
func (h *Grid) verticalOffset(u int) float64 {
	if hexcoords.ColumnIsHigh(u) {
		return 2 * h.radius * math.Sin(hex.TriangleAngle)
	}
	return 0
}

func (h *Grid) TileCenter(c hexcoords.Hex) point.Point {
	var (
		centerX = float64(c.U) * h.horizontalSpacing()
		centerY = float64(c.V) * h.verticalSpacing()
	)
	centerY += h.verticalOffset(c.U)
	return point.Point{centerX, centerY}
}

func (h *Grid) genTiles(defaultValue Value) {
	h.t = make([]Tile, 0, h.expectedNumTiles())
	// Generate all tiles.
	h.tiles = make([][]*Tile, h.n)
	for i := 0; i < h.n; i++ {
		h.tiles[i] = make([]*Tile, h.m)
		for j := 0; j < h.m; j++ {
			var (
				coords = h.hexCoords(i, j)
				center = h.TileCenter(coords)
				value  Value
			)
			switch defaultValue.(type) {
			case func(hexcoords.Hex) Value:
				var f = defaultValue.(func(hexcoords.Hex) Value)
				value = f(coords)
			default:
				value = defaultValue
			}
			h.t = append(h.t, Tile{Hex: coords, Pos: center, Value: value})
			h.tiles[i][j] = &(h.t[len(h.t)-1])
		}
	}
}
func (h *Grid) genVertices(defaultValue Value) {
	// Make space for vertices/pointers.
	h.v = make([]Vertex, 0, h.expectedNumVertices())
	h.vertices = make([][][]*Vertex, h.n)
	for i := 0; i < h.n; i++ {
		h.vertices[i] = make([][]*Vertex, h.m)
		for j := 0; j < h.m; j++ {
			h.vertices[i][j] = make([]*Vertex, 6)
		}
	}
	// Generate vertices
	for i := 0; i < h.n; i++ {
		for j := 0; j < h.m; j++ {
			var (
				c   = h.hexCoords(i, j)
				hex = h.GetHex(c)
			)
			for k := 0; k < 6; k++ {
				if h.vertices[i][j][k] == nil {
					var (
						vert          = hexcoords.Vertex{c.U, c.V, k}
						identVertices = vert.IdenticalVertices()
						coords        = vert
						value         Value
					)
					if identVertices == nil {
						panic("outofbounds")
					}
					switch defaultValue.(type) {
					case func(hexcoords.Vertex) Value:
						var f = defaultValue.(func(hexcoords.Vertex) Value)
						value = f(coords)
					default:
						value = defaultValue
					}
					h.v = append(h.v, Vertex{Hex: coords, Pos: hex[k], Value: value})
					for _, ident := range identVertices {
						var (
							c              = ident.Hex()
							iIdent, jIdent = h.hexIndex(c)
						)
						if h.WithinBounds(c) {
							h.vertices[iIdent][jIdent][ident.K] = &(h.v[len(h.v)-1])
						}
					}
				}
			}
		}
	}
}
func (h *Grid) genEdges(defaultValue Value) {
	// Make space for edges/pointers.
	h.e = make([]Edge, 0, h.expectedNumEdges())
	h.edges = make([][][][]*Edge, h.n)
	for i := 0; i < h.n; i++ {
		h.edges[i] = make([][][]*Edge, h.m)
		for j := 0; j < h.m; j++ {
			h.edges[i][j] = make([][]*Edge, 6)
			for k := 0; k < 6; k++ {
				h.edges[i][j][k] = make([]*Edge, 6)
			}
		}
	}
	// Generate all edges.
	for i := 0; i < h.n; i++ {
		for j := 0; j < h.m; j++ { // BEGIN (u,v) TILE ANALYSIS
			var (
				c = h.hexCoords(i, j)
			)
			for k := 0; k < 6; k++ {
				for ell := 0; ell < 6; ell++ { // BEGIN (k,ell) EDGE ANALYSIS
					// Ensure an edge between k and ell exists.
					var edgeDir = hex.EdgeDirection(k, ell)
					if edgeDir != hex.NilDirection && h.edges[i][j][k][ell] == nil {
						var (
							coords = hexcoords.Edge{c.U, c.V, k, ell}
							value  Value
							v1     = h.vertices[i][j][k]
							v2     = h.vertices[i][j][ell]
						)
						switch defaultValue.(type) {
						case func(hexcoords.Edge, *Vertex, *Vertex) Value:
							var f = defaultValue.(func(hexcoords.Edge, *Vertex, *Vertex) Value)
							value = f(coords, v1, v2)
						default:
							value = defaultValue
						}
						// Create the edge, compute the other incident tile.
						h.e = append(h.e, Edge{Hex: coords, Value: value})
						var (
							edgePtr        = &(h.e[len(h.e)-1])
							adjEdgeIndices = hex.HexEdgeIndices(edgeDir.Inverse())
						)
						if adjEdgeIndices == nil {
							panic("niladjindices")
						}
						var (
							adjK   = adjEdgeIndices[0]
							adjEll = adjEdgeIndices[1]
						)
						if edgeDir == hex.NilDirection {
							panic("niledgedirection")
						}
						var adjHexSlice = c.Adjacents(edgeDir)
						if adjHexSlice == nil {
							panic("niladjcoords")
						}
						// Store the edge pointer is its various configurations.
						var (
							adjHex = adjHexSlice[0]
							//adjU       = adjHex.U
							//adjV       = adjHex.V
							adjI, adjJ = h.hexIndex(adjHex)
						)
						if h.WithinBounds(adjHex) {
							h.edges[adjI][adjJ][adjK][adjEll] = edgePtr
							h.edges[adjI][adjJ][adjEll][adjK] = edgePtr
						}
						h.edges[i][j][k][ell] = edgePtr
						h.edges[i][j][ell][k] = edgePtr
					}
				} // END (k,ell) EDGE ANALYSIS
			}
		} // END (u,v) TILE ANALYSIS
	}
}
func (h *Grid) genHexagons() {
	h.p = make([]point.Point, 0, h.expectedNumVertices())
	h.hexes = make([][]*hex.HexPoints, h.n)
	// Generate all hexagons.
	for i := 0; i < h.n; i++ {
		h.hexes[i] = make([]*hex.HexPoints, h.m)
		for j := 0; j < h.m; j++ {
			var (
				c = h.hexCoords(i, j)
			)
			h.hexes[i][j] = h.GetHex(c)
			if h.hexes[i][j] == nil {
				panic(fmt.Sprintf("OutOfBounds(%d,%d)", i, j))
			}
		}
	}

	// Collect all points, sharing common points belonging to adjacent hexagons.
	for i := 0; i < h.n; i++ {
		for j := 0; j < h.m; j++ {
			var (
				toAdd = [6]bool{true, true, true, true, true, true}
				c     = h.hexCoords(i, j)
				hex   = h.hexes[i][j]
			)
			for k := 0; k < 6; k++ {
				var idents = hexcoords.Vertex{c.U, c.V, k}.IdenticalVertices()
				for _, id := range idents[1:] {
					var joinVertices = func() {
						var (
							adjHex = h.GetHex(id.Hex())
						)
						if adjHex == nil {
							return
						}
						hex[k] = adjHex[id.K]
						//log.Printf("Joined (%d %d %d) and (%d %d %d)", u, v, k, id.U, id.V, id.K)
						toAdd[k] = false
					}
					if id.V < c.V {
						joinVertices()
					} else if id.V == c.V {
						if id.U < c.U {
							joinVertices()
						}
					}
				}
			}

			// Account for points in the hex tile not already accounted for.
			for k, shouldAdd := range toAdd {
				if !shouldAdd {
					continue
				}
				h.p = append(h.p, hex[k])
			}
		}
	}
}
