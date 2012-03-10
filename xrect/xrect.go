/*
    xrect defines several utility functions that perform math on X rectangles.
    X rectangles are defined as a 4-tuple (x, y, w, h) where x and y are the
    top-left coordinates in the x,y plane (with origin at the top left corner).
    w and h are the width and height of the rectangle.

    One additional constraint is that x and y are signed 16 bit integers
    (int16). In particular, they may be negative!

    w and h are unsigned 16 bit integers (uint16). If they are negative, X
    will yell at you and then stomp on you.
*/
package xrect

// import "github.com/BurntSushi/xgbutil" 

// Define a base and simple Rect interface.
type Rect interface {
    X() int16
    Y() int16
    Width() uint16
    Height() uint16
}

// Turn all elements of a Rect interface into integers
func Intify(xr Rect) (int, int, int, int) {
    return int(xr.X()), int(xr.Y()), int(xr.Width()), int(xr.Height())
}

// Provide a simple implementation of a rect.
// Maybe this will be all we need?
type XRect struct {
    x, y int16
    width, height uint16
}

// Provide the ability to construct an XRect.
func Make(x, y int16, w, h uint16) *XRect {
    return &XRect{x, y, w, h}
}

// Satisfy the Rect interface
func (r *XRect) X() int16 {
    return r.x
}

func (r *XRect) Y() int16 {
    return r.y
}

func (r *XRect) Width() uint16 {
    return r.width
}

func (r *XRect) Height() uint16 {
    return r.height
}

// IntersectArea takes two rectangles satisfying the Rect interface and
// returns the area of their intersection. If there is no intersection, return
// 0 area.
func IntersectArea(r1 Rect, r2 Rect) int {
    x1, y1, w1, h1 := Intify(r1)
    x2, y2, w2, h2 := Intify(r2)
    if x2 < x1 + w1 && x2 + w2 > x1 && y2 < y1 + h1 && y2 + h2 > y1 {
        iw := Min(x1 + w1 - 1, x2 + w2 - 1) - Max(x1, x2) + 1
        ih := Min(y1 + h1 - 1, y2 + h2 - 1) - Max(y1, y2) + 1
        return iw * ih
    }

    return 0
}

// LargestOverlap returns the rectangle in 'haystack' that has the largest
// overlap with the rectangle 'needle'. This is commonly used to find which
// monitor a window should belong on. (Since it can technically be partially
// displayed on more than one monitor at a time.)
func LargestOverlap(needle Rect, haystack []Rect) (result Rect) {
    biggestArea := 0

    var area int
    for _, possible := range haystack {
        area = IntersectArea(needle, possible)
        if area > biggestArea {
            biggestArea = area
            result = possible
        }
    }
    return
}

// Min should be in Go's standard library... but not for floats.
func Min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

// Max should be in Go's standard library... but not for floats.
func Max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

