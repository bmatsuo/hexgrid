package hexgrid
/* 
*  File: hexgrid.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Tue Jun 28 03:40:52 PDT 2011
 */
import (
    "fmt"
    "math"
    //"log"
)

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
        {{-1, 0, 2}, {-1, -1, 3}}}
)

type Value interface{}
//  For each coordinate in a Grid there is one unique HexTile.
type Tile struct {
    Coords Coords
    Pos    Point
    Value  Value
}
type TileInitializer    func (Coords) Value
//  A HexVertex represents the corner a HexTile. A HexVertex can be shared
//  by at most 3 HexTiles and can be the junction of between 2 and three
//  HexEdge objects. A HexVertex can be 'between' fewer than three tiles if
//  its tiles are on the edge of the grid. It will be the endpoint of two
//  edges only it belongs to one tile (and is on the edge of the grid).
type Vertex struct {
    Coords  VertexCoords
    Pos     Point
    Value   Value
}
type VertexInitializer  func (VertexCoords) Value
//  A HexEdge represents an edge between two HexVertex objects. It is
//  part of the boundary of a HexTile. A HexEdge can be 'Between' only
//  one tile if its tile is on the edge of the grid.
type Edge struct {
    Coords  EdgeCoords
    Value   Value
}
type EdgeInitializer    func (coords EdgeCoords, v1, v2 *Vertex) Value

//  A grid of hexagons in a discrete coordinate system (u,v) where u
//  indexes the column of the grid, and v the row.
type Grid struct {
    radius   float64
    n        int
    m        int
    p        []Point
    v        []Vertex
    e        []Edge
    t        []Tile
    hexes    [][]*HexPoints
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
func (h *Grid) GetTile(u, v int) *Tile {
    if !h.WithinBounds(u, v) {
        return nil
    }
    i, j := h.hexIndex(u, v)
    return h.tiles[i][j]
}

//  Retrieve a Vertex object specified by its coordinates.
func (h *Grid) GetVertex(u, v, k int) *Vertex {
    if !h.WithinBounds(u, v) {
        return nil
    }
    if k < 0 {
        return nil
    }
    i, j := h.hexIndex(u, v)
    return h.vertices[i][j][k%6]
}

//  Retrieve an Edge object specified by its coordinates.
func (h *Grid) GetEdge(u, v, k, ell int) *Edge {
    if !h.WithinBounds(u, v) {
        return nil
    }
    if k < 0 || ell < 0 {
        return nil
    }
    i, j := h.hexIndex(u, v)
    return h.edges[i][j][k%6][ell%6]
}
func (h *Grid) GetEdges(coords Coords) []*Edge {
    if !h.WithinBounds(coords.U, coords.V) {
        return nil
    }
    var edges = make([]*Edge, 6)
    for k := 0 ; k < 6 ; k++ {
        edges[k] = h.GetEdge(coords.U, coords.V, k, (k+1)%6)
    }
    return edges
}
//  Returns the width and height of the Grid wrapped in a
//  GridDimensions object.
func (h *Grid) Size() GridDimensions {
    return GridDimensions{float64(h.n), float64(h.m)}
}

//  Total number of distinct hexagon vertices in the field.
func (h *Grid) expectedNumVertices() int {
    return 2*(h.n*h.m + h.n + h.m)
}
func (h *Grid) NumVertices() int {
    return len(h.v)
}
func (h *Grid) expectedNumEdges() int {
    return 3*h.n*h.m + 2*h.n + 2*h.m-1
}
func (h *Grid) NumEdges() int {
    return len(h.e)
}
func (h *Grid) expectedNumTiles() int {
    return h.n*h.m
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
func (h *Grid) hexCoords(i, j int) (int, int) {
    return i + h.ColMin(), j + h.RowMin()
}
func (h *Grid) hexIndex(u, v int) (int, int) {
    return u - h.ColMin(), v - h.RowMin()
}

/* Internal bounds checking method. */
func (h *Grid) indexWithinBounds(i, j int) bool {
    u, v := h.hexCoords(i, j)
    return h.WithinBounds(u, v)
}
//  Returns true if the hex at coordinates (u,v) is in the hex field.
func (h *Grid) WithinBounds(u, v int) bool {
    if int(math.Fabs(float64(u))) > h.ColMax() {
        return false
    }
    if int(math.Fabs(float64(v))) > h.RowMax() {
        return false
    }
    return true
}

//  Generate points for the hexagon at row i, column j.
//  Returns nil when the position (i,j) is not within the bounds of the board.
func (h *Grid) GetHex(u, v int) *HexPoints {
    if !h.WithinBounds(u, v) {
        return nil
    }
    var (
        i, j = h.hexIndex(u, v)
    )
    if h.hexes[i][j] != nil {
        var newh = new(HexPoints)
        *newh = *(h.hexes[i][j])
        return newh
    }
    return NewHex(h.TileCenter(u, v), h.radius)
}

//  Get a pointer to the kth corner point of the hex tile at (u,v).
//  Returns PointInf() when (u,v) is not within the bounds of h.
func (h *Grid) GetVertexPoint(vert VertexCoords) Point {
    var hex = h.GetHex(vert.U, vert.V)
    if hex == nil {
        return PointInf()
    }
    return hex[vert.K]
}

//  Get the index of the vertex clockwise of vertex k.
func HexVertexIndexClockwise(k int) int {
    return (k + 5) % 6
}
//  Get the index of the vertex counter-clockwise of vertex k.
func HexVertexIndexCounterClockwise(k int) int {
    return (k + 1) % 6
}
func (h *Grid) GetVertexAdjacentByEdge(vert VertexCoords, edge EdgeCoords) *Vertex {
    var coords = vert.AdjacentByEdge(edge)
    return h.GetVertex(coords.U, coords.V, coords.K)
}
func (h *Grid) GetVertices(coords Coords) []*Vertex {
    if !h.WithinBounds(coords.U, coords.V) {
        return nil
    }
    var vertices = make([]*Vertex, 6)
    for k := 0 ; k < 6 ; k++ {
        vertices[k] = h.GetVertex(coords.U, coords.V, k)
    }
    return vertices
}

//  Get hex tiles incident with the kth corner point of hex at (u,v).
//  Returns nil when (u,v) is not within the bounds of h.
//  Otherwise, a slice of *HexPoints is returned w/ hex tile (u,v) at index 0.
func (h *Grid) GetHexIncident(vert VertexCoords) []*HexPoints {
    if !h.WithinBounds(vert.U, vert.V) {
        return nil
    }
    var adjC = vert.IncidentCoords()
    var adj = make([]*HexPoints, 0, len(adjC))
    for _, coords := range adjC {
        var (
            hexAdj = h.GetHex(coords.U, coords.V)
        )
        if hexAdj == nil {
            panic("coordoutofbounds")
        }
        adj = append(adj, hexAdj)
    }

    return adj
}
func (h *Grid) GetTilesIncident(vert VertexCoords) []*Tile {
    var adjC = vert.IncidentCoords()
    var adj = make([]*Tile, 0, len(adjC))
    for _, coords := range adjC {
        var (
            tileAdj = h.GetTile(coords.U, coords.V)
        )
        if tileAdj == nil {
            continue
        }
        adj = append(adj, tileAdj)
    }

    return adj
}
func (h *Grid) GetTilesSharedByCoords(vert1, vert2 VertexCoords) []*Tile {
    var (
        shared = vert1.SharedByVertex(vert2)
        tiles = make([]*Tile, 0, len(shared))
    )
    for _, coord := range shared {
        var tile = h.GetTile(coord.U, coord.V)
        if tile == nil {
            continue
        }
        tiles = append(tiles, tile)
    }
    if len(tiles) == 0 {
        return nil
    }
    return tiles
}

/* Internal methods for computing hexagon positions. */
func (h *Grid) horizontalSpacing() float64 {
    return 2 * h.radius * math.Cos(hexTriangleAngle)
}
func (h *Grid) verticalSpacing() float64 {
    return 2 * h.radius
}
func (h *Grid) verticalOffset(u int) float64 {
    if columnIsHigh(u) {
        return 2 * h.radius * math.Sin(hexTriangleAngle)
    }
    return 0
}

func (h *Grid) TileCenter(u, v int) Point {
    var (
        centerX = float64(u) * h.horizontalSpacing()
        centerY = float64(v) * h.verticalSpacing()
    )
    centerY += h.verticalOffset(u)
    return Point{centerX, centerY}
}

//  Return a slice of hexagons adjacent to the hex tile at coordinates (u, v).
//  Only hex tiles in the Grid are returned.
//  If (u,v) is not in the Grid, a nil slice is returned.
func (h *Grid) GetHexAdjacent(u, v int, dir HexDirection) []*HexPoints {
    if !h.WithinBounds(u, v) {
        return nil
    }
    var adjC = Coords{u, v}.AdjacentCoords(dir)
    if adjC == nil {
        panic("niladjacency")
    }
    var adj = make([]*HexPoints, 0, len(adjC))
    for _, coords := range adjC {
        var (
            uPrime = coords.U
            vPrime = coords.V
        )
        if !h.WithinBounds(uPrime, vPrime) {
            continue
        }
        adj = append(adj, h.GetHex(uPrime, vPrime))
    }
    return adj
}

func (h *Grid) GetEdgeSharedByVertices(vert1, vert2 VertexCoords) *Edge {
    var coords = vert1.EdgeCoordsSharedByVertex(vert2)
    return h.GetEdge(coords.U, coords.V, coords.K, coords.L)
}

//  Function for determining the edge container if any,
//  between the hex tile at (u1,v1) that is alse in tile
//  (u2,v2). Returns nil if the hex coordinates are not
//  adjacent.
func (h *Grid) GetEdgeShared(coord1, coord2 Coords) *Edge {
    var indices = coord1.EdgeIndicesShared(coord2)
    if indices == nil {
        return nil
    }
    return h.GetEdge(coord1.U, coord1.V, indices[0], indices[1])
}
//  Function for determining the actual points determining any shared edge
//  between hex tiles (u1,v1) and (u2,v2). Returns nil if either
//  coordinates are outside of the hex field. Returns nil if the hex
//  coordinates are not adjacent.
func (h *Grid) GetEdgePointsShared(coord1, coord2 Coords) []Point {
    if !h.WithinBounds(coord1.U, coord1.V) {
        return nil
    }
    if !h.WithinBounds(coord2.U, coord2.V) {
        return nil
    }
    var sharedIndices = coord1.EdgeIndicesShared(coord2)
    if sharedIndices == nil {
        return nil
    }
    var h1 = h.GetHex(coord1.U, coord2.V)
    return []Point{h1[sharedIndices[0]], h1[sharedIndices[1]]}
}

func (h *Grid) genTiles(defaultValue Value) {
    h.t = make([]Tile, 0, h.expectedNumTiles())
    // Generate all tiles.
    h.tiles = make([][]*Tile, h.n)
    for i := 0; i < h.n; i++ {
        h.tiles[i] = make([]*Tile, h.m)
        for j := 0; j < h.m; j++ {
            var (
                u, v    = h.hexCoords(i, j)
                coords  = Coords{u, v}
                center  = h.TileCenter(u, v)
                value   Value
            )
            switch defaultValue.(type) {
            case func (Coords) Value:
                var f = defaultValue.(func (Coords) Value)
                value = f(coords)
            default:
                value = defaultValue
            }
            h.t = append(h.t, Tile{Coords:coords, Pos:center, Value:value})
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
                u, v = h.hexCoords(i, j)
                hex  = h.GetHex(u, v)
            )
            for k := 0; k < 6; k++ {
                if h.vertices[i][j][k] == nil {
                    var (
                        vert            = VertexCoords{u, v, k}
                        identVertices   = vert.IdenticalVertices()
                        coords          = vert
                        value           Value
                    )
                    if identVertices == nil {
                        panic("outofbounds")
                    }
                    switch defaultValue.(type) {
                    case func (VertexCoords) Value:
                        var f = defaultValue.(func (VertexCoords) Value)
                        value = f(coords)
                    default:
                        value = defaultValue
                    }
                    h.v = append(h.v, Vertex{Coords:coords, Pos:hex[k], Value:value})
                    for _, ident := range identVertices {
                        var (
                            uIdent         = ident.U
                            vIdent         = ident.V
                            kIdent         = ident.K
                            iIdent, jIdent = h.hexIndex(uIdent, vIdent)
                        )
                        if h.WithinBounds(uIdent, vIdent) {
                            h.vertices[iIdent][jIdent][kIdent] = &(h.v[len(h.v)-1])
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
                u, v = h.hexCoords(i, j)
            )
            for k := 0; k < 6; k++ {
                for ell := 0; ell < 6; ell++ { // BEGIN (k,ell) EDGE ANALYSIS
                    // Ensure an edge between k and ell exists.
                    var edgeDir = HexEdgeDirection(k, ell)
                    if edgeDir != NilDirection && h.edges[i][j][k][ell] == nil {
                        var (
                            coords  = EdgeCoords{u, v, k, ell}
                            value   Value
                            v1      = h.vertices[i][j][k]
                            v2      = h.vertices[i][j][ell]
                        )
                        switch defaultValue.(type) {
                        case func (EdgeCoords, *Vertex, *Vertex) Value:
                            var f = defaultValue.(func (EdgeCoords, *Vertex, *Vertex) Value)
                            value = f(coords, v1, v2)
                        default:
                            value = defaultValue
                        }
                        // Create the edge, compute the other incident tile.
                        h.e = append(h.e, Edge{Coords:coords, Value: value})
                        var (
                            edgePtr        = &(h.e[len(h.e)-1])
                            adjEdgeIndices = HexEdgeIndices(edgeDir.Inverse())
                        )
                        if adjEdgeIndices == nil {
                            panic("niladjindices")
                        }
                        var (
                            adjK   = adjEdgeIndices[0]
                            adjEll = adjEdgeIndices[1]
                        )
                        if edgeDir == NilDirection {
                            panic("niledgedirection")
                        }
                        var adjCoordsSlice = Coords{u, v}.AdjacentCoords(edgeDir)
                        if adjCoordsSlice == nil {
                            panic("niladjcoords")
                        }
                        // Store the edge pointer is its various configurations.
                        var (
                            adjCoords  = adjCoordsSlice[0]
                            adjU       = adjCoords.U
                            adjV       = adjCoords.V
                            adjI, adjJ = h.hexIndex(adjU, adjV)
                        )
                        if h.WithinBounds(adjU, adjV) {
                            h.edges[adjI][adjJ][adjK][adjEll] = edgePtr
                            h.edges[adjI][adjJ][adjEll][adjK] = edgePtr
                        }
                        h.edges[i][j][k][ell] = edgePtr
                        h.edges[i][j][ell][k] = edgePtr
                    }
                }   // END (k,ell) EDGE ANALYSIS
            }
        }   // END (u,v) TILE ANALYSIS
    }
}
func (h *Grid) genHexagons() {
    h.p = make([]Point, 0, h.expectedNumVertices())
    h.hexes = make([][]*HexPoints, h.n)
    // Generate all hexagons.
    for i := 0; i < h.n; i++ {
        h.hexes[i] = make([]*HexPoints, h.m)
        for j := 0; j < h.m; j++ {
            var (
                u, v = h.hexCoords(i, j)
            )
            h.hexes[i][j] = h.GetHex(u, v)
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
                u, v  = h.hexCoords(i,j)
                hex   = h.hexes[i][j]
            )
            for k := 0 ; k < 6 ; k++ {
                var idents = VertexCoords{u, v, k}.IdenticalVertices()
                for _, id := range idents[1:] {
                    if h.WithinBounds(id.U, id.V) {
                        var joinVertices = func() {
                            var (
                                iPrime, jPrime = h.hexIndex(id.U, id.V)
                                adjHex = h.hexes[iPrime][jPrime]
                            )
                            hex[k] = adjHex[id.K]
                            toAdd[k] = false
                        }
                        if id.V < v {
                            joinVertices()
                        } else if id.V == v {
                            if id.U < u {
                                joinVertices()
                            }
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
