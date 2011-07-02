package hexgrid
/* 
*  File: coords.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sat Jul  2 00:54:20 PDT 2011
*/
import (
    "math"
    //"log"
)

//  Discrete hex coordinates consist of a horizontal U axis and a vertical
//  V axis. Each axis has range (-inf,inf) in theory. In practice, Grid
//  objects limit the accessible hex tiles.
type Coords struct {
    U, V int
}
func (coords Coords) Equals(other Coords) bool {
    return coords.U == other.U && coords.V == other.V
}
type GridDimensions struct {
    U, V float64
}
//  Vertices in the grid are indexed by hex coordinates paired with a
//  vertex index K. Vertex indices range from 0 to 5 and begin in the
//  south-west corner of the vertex. See also, HexDirection.
type VertexCoords struct {
    U, V, K int
}
//  Create a Coords object from a VertexCoords object.
func (vc VertexCoords) Coords() Coords {
    return Coords{vc.U, vc.V}
}
//  Returns true if  only if the U, V and K fields of vc and other are equal.
func (vc VertexCoords) Equals(other VertexCoords) bool {
    return vc.U == other.U && vc.V == other.V && vc.K == other.K
}
//  Returns true if (u1,v1,k1) and (u2,v2,k2) reference the same vertex.
func (vc VertexCoords) Identical(other VertexCoords) bool {
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

//  Edges in the grid are index by hex coordinates along with a pair of
//  vertex indices K and L.
type EdgeCoords struct {
    U, V, K, L int
}
func (e EdgeCoords) Ends() (v1, v2 VertexCoords) {
    v1 = VertexCoords{e.U,e.V,e.K}
    v2 = VertexCoords{e.U,e.V,e.L}
    return v1, v2
}

func columnIsHigh(u int) bool {
    var uOdd = uint(math.Fabs(float64(u)))%2 == 1
    return uOdd
}
func sameTile(u1, v1, u2, v2 int) bool {
    return Coords{u1, v1}.Equals(Coords{u2, v2})
}


//  If hex tiles (u1,v1) and (u2,v2) are adjacent, the direction of (u2,v2)
//  from (u1,v1) is returned. Otherwise NilDirection is returned.
func (c Coords) Adjacency(adj Coords) HexDirection {
    var (
        deltaU = adj.U - c.U
        deltaV = adj.U - c.U
    )
    if c.U == adj.U {
        if deltaV == 1 {
            return N
        } else if deltaV == -1 {
            return S
        }
        return NilDirection
    }
    if deltaU == 1 {
        if columnIsHigh(c.U) {
            if c.U == adj.U {
                return SE
            }
            if adj.U == c.U+1 {
                return NE
            }
        } else {
            if adj.V == c.V-1 {
                return SE
            }
            if c.V == adj.V {
                return NE
            }
        }
    } else if deltaU == -1 {
        if columnIsHigh(c.U) {
            if adj.V == c.V+1 {
                return NW
            }
            if c.V == adj.V {
                return SW
            }
        } else {
            if c.V == adj.V {
                return NW
            }
            if adj.V == c.V-1 {
                return SW
            }
        }
    }
    return NilDirection
}

//  Returns true if and only if c is adjacent to adj
func (c Coords) Adjacent(adj Coords) bool {
    return c.Adjacency(adj) != NilDirection
}

//  Return a slice of the coordinates for adjacent hexagons
//  (not necessarily in the grid).
//  If E (or W) is supplied then the NE and SE (or NE and SW) coordinates
//  are returned in that order.
//  If NilDirection is suppied, then coordinates for all adjacent hexagons
//  are returned in the order N, NE, SE, S, SW, NW.
func (coords Coords) Adjacents(dir HexDirection) []Coords {
    var (
        u = coords.U
        v = coords.V
    )
    switch dir {
    case N:
        return []Coords{Coords{u, v + 1}}
    case S:
        return []Coords{Coords{u, v - 1}}
    case E:
        var adjE = make([]Coords, 2)
        if columnIsHigh(u) {
            adjE[0] = Coords{u - 1, v + 1}
            adjE[1] = Coords{u - 1, v}
        } else {
            adjE[0] = Coords{u - 1, v}
            adjE[1] = Coords{u - 1, v - 1}
        }
        return adjE
    case W:
        var adjW = make([]Coords, 2)
        if columnIsHigh(u) {
            adjW[0] = Coords{u + 1, v + 1}
            adjW[1] = Coords{u + 1, v}
        } else {
            adjW[0] = Coords{u + 1, v}
            adjW[1] = Coords{u + 1, v - 1}
        }
        return adjW
    case NE:
        if columnIsHigh(u) {
            return []Coords{Coords{u - 1, v + 1}}
        }
        return []Coords{Coords{u - 1, v}}
    case NW:
        if columnIsHigh(u) {
            return []Coords{Coords{u + 1, v + 1}}
        }
        return []Coords{Coords{u + 1, v}}
    case SE:
        if columnIsHigh(u) {
            return []Coords{Coords{u - 1, v}}
        }
        return []Coords{Coords{u - 1, v - 1}}
    case SW:
        if columnIsHigh(u) {
            return []Coords{Coords{u + 1, v}}
        }
        return []Coords{Coords{u + 1, v - 1}}
    default:
        var adjAll = make([]Coords, 6)
        if columnIsHigh(u) {
            adjAll[0] = Coords{u, v + 1}     // North
            adjAll[1] = Coords{u - 1, v + 1} // NorthEast
            adjAll[2] = Coords{u - 1, v}     // SouthEast
            adjAll[3] = Coords{u, v - 1}     // South
            adjAll[4] = Coords{u + 1, v}     // SouthWest
            adjAll[5] = Coords{u + 1, v + 1} // NorthWest
        } else {
            adjAll[0] = Coords{u, v + 1}
            adjAll[1] = Coords{u - 1, v}
            adjAll[2] = Coords{u - 1, v - 1}
            adjAll[3] = Coords{u, v - 1}
            adjAll[4] = Coords{u + 1, v - 1}
            adjAll[5] = Coords{u + 1, v}
        }
        return adjAll
    }
    return nil
}

func (vert VertexCoords) Incidents() []Coords {
    var (
        adjVC = vert.IdenticalVertices()
        adj = make([]Coords, 0, len(adjVC))
    )
    for _, coords := range adjVC {
        var (
            adjCoords = Coords{coords.U, coords.V}
        )
        adj = append(adj, adjCoords)
    }

    return adj
}

func (edge EdgeCoords) SharedByEnds() []Coords {
    vert1, vert2 := edge.Ends()
    return vert1.SharedByVertex(vert2)
}

func (vert VertexCoords) SharedByVertex(vert2 VertexCoords) []Coords {
    var (
        adjC1 = vert.Incidents()
        adjC2 = vert2.Incidents()
    )
    var shared = make([]Coords, 0, 1)
    for _, c1 := range adjC1 {
        for _, c2 := range adjC2 {
            if c1.U == c2.U && c1.V == c2.V {
                shared = append(shared, c1)
            }
        }
    }
    return shared
}

func (vert VertexCoords) EdgeShared(vert2 VertexCoords) EdgeCoords {
    if vert.Identical(vert2) {
        return EdgeCoords{}
    }
    if vert.Coords().Equals(vert2.Coords()) {
        return EdgeCoords{vert.U, vert.V, vert.K, vert2.K}
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
            if ident1.Coords().Equals(ident2.Coords()) {
                return EdgeCoords{ident1.U, ident1.V, ident1.K, ident2.K}
            }
            var edge = ident1.Coords().EdgeIndicesShared(ident2.Coords())
            if edge != nil {
                return EdgeCoords{ident1.U, ident1.V, edge[0], edge[1]}
            }
        }
    }
    return EdgeCoords{}
}


//  Function for determining the vertex indices of an edge in
//  the hex tile at (u1,v1) that is alse in tile (u2,v2).
//  Returns nil if the hex coordinates are not adjacent.
func (coord Coords) EdgeIndicesShared(other Coords) []int {
    var adjDir = coord.Adjacency(other)
    if adjDir == NilDirection {
        return nil
    }
    return HexEdgeIndices(adjDir)
}

//  This method needs testing.
func (vert VertexCoords) IncidentEdges() []EdgeCoords {
    var (
        adjVCs = vert.Adjacents()
        edges = make([]EdgeCoords, len(adjVCs))
    )
    for i, other := range adjVCs {
        edges[i] = vert.EdgeShared(other)
    }
    return edges
}

//  This is untested.
func (vert VertexCoords) AdjacentByEdge(edge EdgeCoords) VertexCoords {
    v1, v2 := edge.Ends()
    if vert.Identical(v1) {
        return v2
    } else if vert.Identical(v2) {
        return v1
    }
    return VertexCoords{}
}

//  Get coordinates of hex vertices in the field incident to vertex (u,v,k).
//  Returns a slice of vertex coordinates (slices of 3 ints), the first of
//  which being []int{u, v, k}. See also, (vc VertexCoords) Identical.
func (vert VertexCoords) IdenticalVertices() []VertexCoords {
    var adjC = make([]VertexCoords, 1, 3)
    adjC[0] = vert

    var adjOffsets [][]int
    if columnIsHigh(vert.U) {
        adjOffsets = hexHighVertexIncidenceOffset[vert.K]
    } else {
        adjOffsets = hexLowVertexIncidenceOffset[vert.K]
    }
    for _, offset := range adjOffsets {
        var (
            du         = offset[0]
            dv         = offset[1]
            kAdj       = offset[2]
            uNew       = vert.U + du
            vNew       = vert.V + dv
        )
        var offsetCp = VertexCoords{uNew, vNew, kAdj}
        adjC = append(adjC, offsetCp)
    }

    return adjC
}

//  Get a list of unique vertices adjacent to (u,v,k).
//  See also, (VertexCoords) Identical.
func (vert VertexCoords) Adjacents() []VertexCoords {
    var identVerts = vert.IdenticalVertices()
    if identVerts == nil {
        return nil
    }
    var adjVerts = make([]VertexCoords, len(identVerts))
    for i, vert := range identVerts {
        adjVerts[i] = VertexCoords{vert.U, vert.V, HexVertexIndexClockwise(vert.K)}
    }
    return adjVerts
}
