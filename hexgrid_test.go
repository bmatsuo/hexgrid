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
    idesc = imageDesc{800,800,32}
    hfield = NewHexGrid(15, 80)
    hexLabelFontSize float64 = 24
    hexLabelOffset = Point{-63, -10}
    hexCornerLabelFontSize float64 = 12
    hexCornerLabelOffset = Point{-3,-6}
)

func pointInList(p Point, points []Point) bool {
    var found = false
    for _,p2 := range points {
        found = p.ApproxEqual(p2)
        if found {
            break
        }
    }
    return found
}

/*
func writeToPng(m image.Image, pngPath string) {
	f, err := os.Create(pngPath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	b := bufio.NewWriter(f)
	err = png.Encode(b, m)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
*/

func hexAsString(h *HexGrid, i, j int, T *testing.T) string {
    var (
        sbuff = new([6]string)
        hex = h.GetHex(i,j)
    )
    if hex == nil {
        T.Errorf("Nil tile: n: %d i: %d j:%d", h.Size(), i, j)
        return ""
    }
    for k := 0 ; k < 6 ; k++ {
        sbuff[k] = fmt.Sprintf("[%.02f %.02f]", hex[k].X, hex[k].Y)
    }
    return strings.Join(sbuff[:], ", ")
}

/*
func drawHex(gc draw2d.GraphicContext, bounds image.Rectangle, h *HexGrid, i, j int, color image.Color, T *testing.T) {
    var hex = h.GetHex(i,j)
    if hex == nil {
        T.Errorf("Nil tile: n: %d i: %d j:%d", h.Size(), i, j)
        return
    }
    var path = PolygonDrawPath(hex, bounds)
    gc.SetLineWidth(5)
    gc.SetStrokeColor(color)
    gc.Stroke(path)
}

func drawHexPosition(gc draw2d.GraphicContext, bounds image.Rectangle, h *HexGrid, i,j int, T *testing.T) {
    if !h.WithinBounds(i,j) {
        T.Errorf("Internal error: (%d,%d) is not within HexGrid bounds", i,j)
    }

    var (
        center = h.tileCenter(i,j)
        offCenter = float64((bounds.Max.X-bounds.Min.X)/2)
        placement = center.Add(Point{offCenter,offCenter}).Add(hexLabelOffset)
    )

    gc.MoveTo(placement.ImageCoords(bounds).X, placement.ImageCoords(bounds).Y)
    gc.SetFontSize(hexLabelFontSize)
    gc.SetFontData(draw2d.FontData{"luxi", draw2d.FontFamilyMono, draw2d.FontStyleBold})
    gc.FillString(fmt.Sprintf("(%2d %2d)", i,j))
}

func drawHexCornerNames(gc draw2d.GraphicContext, bounds image.Rectangle, h *HexGrid, i,j int, T *testing.T) {
    var hex = h.GetHex(i,j)

    if hex == nil {
        T.Errorf("Nil tile: n: %d i: %d j:%d", h.Size(), i, j)
        return
    }

    var (
        center = h.tileCenter(i,j)
        centerOffset = float64((bounds.Max.X-bounds.Min.X)/2)
        offCenter = Point{centerOffset, centerOffset}
        offText = hexCornerLabelOffset
    )


    gc.SetFontSize(hexCornerLabelFontSize)
    gc.SetFontData(draw2d.FontData{"luxi", draw2d.FontFamilyMono, draw2d.FontStyleBold})
    for k := 0 ; k < 6 ; k++ {
        var (
            toCorner = hex[k].Sub(center).Scale(0.8)
            placement = center.Add(toCorner).Add(offCenter).Add(offText)
        )
        gc.MoveTo(placement.ImageCoords(bounds).X, placement.ImageCoords(bounds).Y)
        gc.FillString(fmt.Sprintf("%d", k))
    }
}
*/

func TestHexGridMissingPoints(T *testing.T) {
    var numPointsMissing = 0
    for i, hexSlice := range hfield.hexes {
        for j, hex := range hexSlice {
            for k, p := range hex[:] {
                if !pointInList(p, hfield.p) {
                    T.Errorf("Couldn't find point (%3d,%3d,%3d) %v in point list", i,j,k,p)
                    numPointsMissing++
                }
            }
        }
    }
    T.Logf("%d missing points", numPointsMissing)
}

func TestHexGridDuplicatedPoints(T *testing.T) {
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

func TestHexGridNumPoints(T *testing.T) {
    if len(hfield.p) != cap(hfield.p) {
        T.Errorf("Allocation performed %d %d", len(hfield.p), cap(hfield.p))
    }
    T.Logf("No wasted slice space (%d points)", len(hfield.p))
}

/* THIS NEEDS TO BE MOVED. THE REST IS BAD
func TestDrawHexGrid(T *testing.T) {
    var (
        output = "img/test.png"
        width = 400
        height = 400
        img = image.NewRGBA(2*width, 2*height)
        gc  = draw2d.NewGraphicContext(img)
        h = hfield
        colors = []image.Color{
            image.RGBAColor{0xFF, 0x00, 0x00, 0xFF}, image.RGBAColor{0x00, 0xFF, 0x00, 0xFF},
            image.RGBAColor{0x00, 0x00, 0xFF, 0xFF}, image.RGBAColor{0x00, 0xFF, 0xFF, 0xFF},
            image.RGBAColor{0xFF, 0x00, 0xFF, 0xFF}, image.RGBAColor{0xFF, 0xFF, 0x00, 0xFF} }
        t1, t2 int64
    )

    T.Log("Drawing hexagons...")
    t1 = time.Nanoseconds()
    for i := h.RowMin() ; i <= h.RowMax() ; i++ {
        for j := h.ColMin() ; j <= h.ColMax() ; j++ {
            drawHex(gc, img.Bounds(), h, i, j, colors[int(math.Fabs(float64(i+j)))%len(colors)], T)
            drawHexPosition(gc, img.Bounds(), h, i, j, T)
            drawHexCornerNames(gc, img.Bounds(), h, i, j, T)
        }
    }
    t2 = time.Nanoseconds()
    T.Logf("Completed (%d hexagons) in %.03fs", h.n*h.n,float64(t2-t1)*1.0e-9)


    T.Logf("Writing image %s...", output)
    t1 = time.Nanoseconds()
    writeToPng(img,output)
    t2 = time.Nanoseconds()
    T.Logf("Completed (PNG %d pixels) in %.03fs", width*height, float64(t2-t1)*1.0e-9)
}

func TestHexGenerate(T *testing.T) {
    var (
        h = hfield
        center = []int{0,0}
        around = [][]int{
                        []int{1,0},
            []int{0,-1},            []int{0,1},
            []int{-1,-1},           []int{-1,1},
                        []int{-1,0}, }
    )

    T.Logf("C (%2d,%2d) %v", center[0], center[1], h.TileCenter(center[0], center[1]))
    for _, adj := range around {
        T.Logf("A (%2d,%2d) %v", adj[0], adj[1], h.TileCenter(adj[0], adj[1]))
    }

    T.Logf("C (%2d,%2d) %s", center[0], center[1],hexAsString(h, center[0], center[1],T))
    for _, adj := range around {
        T.Logf("A (%2d,%2d) %s", adj[0], adj[1], hexAsString(h, adj[0], adj[1],T))
    }
}
*/

