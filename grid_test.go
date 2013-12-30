package hexgrid
/* 
*  File: hexgrid_test.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Tue Jun 28 17:55:54 PDT 2011
 */
import (
	point "github.com/bmatsuo/hexgrid/point"

    //"time"
    "fmt"
    "strings"
    //"math"
    "testing"
    //"os"
    //"log"
    //"bufio"
    //"image"
    //"image/png"
    //"draw2d.googlecode.com/hg/draw2d"
)


type imageDesc struct {
    W, H, D int
}

var (
    idesc                          = imageDesc{800, 800, 32}
    hfield                         = NewGrid(23, 33, 80, nil, nil, nil)
    hexLabelFontSize       float64 = 24
    hexLabelOffset                 = point.Point{-63, -10}
    hexCornerLabelFontSize float64 = 12
    hexCornerLabelOffset           = point.Point{-3, -6}
)

func pointInList(p point.Point, points []point.Point) bool {
    var found = false
    for _, p2 := range points {
        found = p.ApproxEqual(p2)
        if found {
            break
        }
    }
    return found
}

func hexAsString(h *Grid, i, j int, T *testing.T) string {
    var (
        sbuff = new([6]string)
        hex   = h.GetHex(Coords{i, j})
    )
    if hex == nil {
        T.Errorf("Nil tile: n: %d i: %d j:%d", h.Size(), i, j)
        return ""
    }
    for k := 0; k < 6; k++ {
        sbuff[k] = fmt.Sprintf("[%.02f %.02f]", hex[k].X, hex[k].Y)
    }
    return strings.Join(sbuff[:], ", ")
}

func TestGridMissingPoints(T *testing.T) {
    var numPointsMissing = 0
    for i, hexSlice := range hfield.hexes {
        for j, hex := range hexSlice {
            for k, p := range hex[:] {
                if !pointInList(p, hfield.p) {
                    T.Errorf("Couldn't find point (%3d,%3d,%3d) %v in point list", i, j, k, p)
                    numPointsMissing++
                }
            }
        }
    }
    T.Logf("%d missing points", numPointsMissing)
}

func TestGridDuplicatedPoints(T *testing.T) {
    var numPointsDuplicated = 0
    for t, p := range hfield.p {
        if pointInList(p, hfield.p[:t]) {
            T.Errorf("Duplicate point found %v (%d)", p, t)
            numPointsDuplicated++
            continue
        }
        if pointInList(p, hfield.p[t+1:]) {
            T.Errorf("Duplicate point found %v (%d)", p, t)
            numPointsDuplicated++
            continue
        }
    }
    T.Logf("%d duplicate points", numPointsDuplicated)
}

func testAllocation(expected, length, capacity int, T *testing.T, ) {
    if length != capacity {
        T.Errorf("Wasted slice space %d %d", length, capacity)
    } else {
        T.Log("No wasted slice space")
    }
    if length > expected {
        T.Errorf("More than expected %d %d", length, expected)
    } else if length < expected {
        T.Errorf("Less than expected %d %d", length, expected)
    } else {
        T.Logf("Found the expected number %d", expected)
    }
}
func TestHexPoints(T *testing.T) {
    for u := hfield.ColMin() ; u <= hfield.ColMax() ; u++ {
        for v := hfield.RowMin() ; v <= hfield.RowMax() ; v++ {
            var (
                c   = Coords{u, v}
                hex = hfield.GetHex(c)
            )
            if hex == nil {
                T.Errorf("Nil *HexPoints encountered at %d %d", u, v)
                continue
            }
            for k := 0 ; k < 6 ; k++ {
                if pointInList(hex[k], hex[:k]) {
                    T.Errorf("Duplicate point found, index %d, (u,v)=(%d,%d)", k, u, v)
                }
                if pointInList(hex[k], hex[k+1:]) {
                    T.Errorf("Duplicate point found, index %d, (u,v)=(%d,%d)", k, u, v)
                }
            }
        }
    }
}

func TestGridNumTiles(T *testing.T) {
    var (
        expected = hfield.expectedNumTiles()
        s = hfield.t
        length = len(s)
        capacity = cap(s)
    )
    testAllocation(expected, length, capacity, T)
}
func TestGridNumPoints(T *testing.T) {
    var (
        expected = hfield.expectedNumVertices()
        s = hfield.p
        length = len(s)
        capacity = cap(s)
    )
    testAllocation(expected, length, capacity, T)
}
func TestGridNumVertices(T *testing.T) {
    var (
        expected = hfield.expectedNumVertices()
        s = hfield.v
        length = len(s)
        capacity = cap(s)
    )
    testAllocation(expected, length, capacity, T)
}
func TestGridNumEdges(T *testing.T) {
    var (
        expected = hfield.expectedNumEdges()
        s = hfield.e
        length = len(s)
        capacity = cap(s)
    )
    testAllocation(expected, length, capacity, T)
}
