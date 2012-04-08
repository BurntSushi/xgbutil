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

import "fmt"

// Define a base and simple Rect interface.
type Rect interface {
	X() int
	Y() int
	Width() int
	Height() int
	XSet(x int)
	YSet(y int)
	WidthSet(width int)
	HeightSet(height int)
	Pieces() (int, int, int, int)
}

// RectPieces just returns a four-tuple of x, y, width and height
func RectPieces(xr Rect) (int, int, int, int) {
	return xr.X(), xr.Y(), xr.Width(), xr.Height()
}
func Pieces(xr Rect) (int, int, int, int) {
	return RectPieces(xr)
}

// Provide a simple implementation of a rect.
// Maybe this will be all we need?
type XRect struct {
	x, y          int
	width, height int
}

// Provide the ability to construct an XRect.
func Make(x, y, w, h int) *XRect {
	return &XRect{x, y, w, h}
}

func (r *XRect) String() string {
	return fmt.Sprintf("[(%d, %d) %dx%d]", r.x, r.y, r.width, r.height)
}

// Satisfy the Rect interface
func (r *XRect) X() int {
	return r.x
}

func (r *XRect) Y() int {
	return r.y
}

func (r *XRect) Width() int {
	return r.width
}

func (r *XRect) Height() int {
	return r.height
}

func (r *XRect) XSet(x int) {
	r.x = x
}

func (r *XRect) YSet(y int) {
	r.y = y
}

func (r *XRect) WidthSet(width int) {
	r.width = width
}

func (r *XRect) HeightSet(height int) {
	r.height = height
}

// Pieces just returns a four-tuple of x, y, width and height
func (r *XRect) Pieces() (int, int, int, int) {
	return r.X(), r.Y(), r.Width(), r.Height()
}

// IntersectArea takes two rectangles satisfying the Rect interface and
// returns the area of their intersection. If there is no intersection, return
// 0 area.
func IntersectArea(r1 Rect, r2 Rect) int {
	x1, y1, w1, h1 := RectPieces(r1)
	x2, y2, w2, h2 := RectPieces(r2)
	if x2 < x1+w1 && x2+w2 > x1 && y2 < y1+h1 && y2+h2 > y1 {
		iw := Min(x1+w1-1, x2+w2-1) - Max(x1, x2) + 1
		ih := Min(y1+h1-1, y2+h2-1) - Max(y1, y2) + 1
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

// ApplyStrut takes a list of Rects (typically the rectangles that represent
// each physical head in this case) and a set of parameters representing a
// strut, and modifies the list of Rects to account for struts.
// That is, it shrinks each rect.
// Note that if struts overlap, the *most restrictive* one is used. This seems
// like the most sensible response to a weird scenario.
// (If you don't have a partial strut, just use '0' for the extra fields.)
// See tests/tst_rect.go for an example of how to use this to get accurate
// workarea for each physical head.
func ApplyStrut(rects []Rect, rootWidth, rootHeight int,
	left, right, top, bottom,
	left_start_y, left_end_y, right_start_y, right_end_y,
	top_start_x, top_end_x, bottom_start_x, bottom_end_x int) {
	var nx, ny int // 'n*' are new values that may or may not be used
	var nw, nh int
	var x, y, w, h int
	var bt, tp, lt, rt bool
	rWidth, rHeight := rootWidth, rootHeight

	// The essential idea of struts, and particularly partial struts, is that
	// one piece of a border of the screen can be "reserved" for some
	// special windows like docks, panels, taskbars and system trays.
	// Since we assume that one window can only reserve one piece of a border
	// (either top, left, right or bottom), we iterate through each rect
	// in our list and check if that rect is affected by the given strut.
	// If it is, we modify the current rect appropriately.
	// TODO: Fix this so old school _NET_WM_STRUT can work too. It actually
	// should be pretty simple: change conditions like 'if tp' to
	// 'if tp || (top_start_x == 0 && top_end_x == 0 && top != 0)'.
	// Thus, we would end up changing every rect, which is what old school
	// struts should do. We may also make a conscious choice to ignore them
	// when 'rects' has more than one rect, since the old school struts will
	// typically result in undesirable behavior.
	for _, rect := range rects {
		x, y, w, h = RectPieces(rect)

		bt = bottom_start_x != bottom_end_x &&
			(xInRect(bottom_start_x, rect) || xInRect(bottom_end_x, rect))
		tp = top_start_x != top_end_x &&
			(xInRect(top_start_x, rect) || xInRect(top_end_x, rect))
		lt = left_start_y != left_end_y &&
			(yInRect(left_start_y, rect) || yInRect(left_end_y, rect))
		rt = right_start_y != right_end_y &&
			(yInRect(right_start_y, rect) || yInRect(right_end_y, rect))

		if bt {
			nh = h - (bottom - ((rHeight - h) - y))
			if nh < rect.Height() {
				rect.HeightSet(nh)
			}
		} else if tp {
			nh = h - (top - y)
			if nh < rect.Height() {
				rect.HeightSet(nh)
			}

			ny = top
			if ny > rect.Y() {
				rect.YSet(ny)
			}
		} else if rt {
			nw = w - (right - ((rWidth - w) - x))
			if nw < rect.Width() {
				rect.WidthSet(nw)
			}
		} else if lt {
			nw = w - (left - x)
			if nw < rect.Width() {
				rect.WidthSet(nw)
			}

			nx = left
			if nx > rect.X() {
				rect.XSet(nx)
			}
		}
	}
}

// xInRect is whether a particular x-coordinate is vertically constrained by
// a rectangle.
func xInRect(xtest int, rect Rect) bool {
	x, _, w, _ := RectPieces(rect)
	return xtest >= x && xtest < (x+w)
}

// yInRect is whether a particular y-coordinate is horizontally constrained by
// a rectangle.
func yInRect(ytest int, rect Rect) bool {
	_, y, _, h := RectPieces(rect)
	return ytest >= y && ytest < (y+h)
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
