package hexgrid
/* 
*  File: hexgrid_test.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Tue Jun 28 17:55:54 PDT 2011
 */
import (
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
    hfield                         = NewGrid(23, 53, 80, nil, nil, nil)
    hexLabelFontSize       float64 = 24
    hexLabelOffset                 = Point{-63, -10}
    hexCornerLabelFontSize float64 = 12
    hexCornerLabelOffset           = Point{-3, -6}
)

func pointInList(p Point, points []Point) bool {
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
        hex   = h.GetHex(i, j)
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

func TestGridNumPoints(T *testing.T) {
    if len(hfield.p) != cap(hfield.p) {
        T.Errorf("Allocation performed %d %d", len(hfield.p), cap(hfield.p))
    } else {
        T.Logf("No wasted slice space (%d points)", len(hfield.p))
    }
}

func TestGridNumVertices(T *testing.T) {
    if len(hfield.v) != cap(hfield.v) {
        T.Errorf("Allocation performed %d %d", len(hfield.v), cap(hfield.v))
    } else {
        T.Logf("No wasted slice space (%d points)", len(hfield.v))
    }
}

func TestGridNumEdges(T *testing.T) {
    if len(hfield.e) != cap(hfield.e) {
        T.Errorf("Allocation performed %d %d", len(hfield.e), cap(hfield.e))
    } else {
        T.Logf("No wasted slice space (%d points)", len(hfield.e))
    }
}
