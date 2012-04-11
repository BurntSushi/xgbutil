/*
   A few utility functions related to client windows. In 
   particular, getting an accurate geometry of a client window
   including the decorations (this can vary with the window
   manager). Also, a functon to move and/or resize a window
   accurately by the top-left corner. (Also can change based on
   the currently running window manager.) 

   This module also contains a function 'Listen' that must be used 
   in order to receive certain events from a window.

   {SHOW EXAMPLE}

   The idea here is to tell X that you want events that fall under
   the 'PropertyChange' category. Then you bind 'func' to the 
   particular event 'PropertyNotify'.

   Most of the methods here aren't useful for window manager developers.
   Particularly the 'GetGeometry' and move/resizing methods---they are
   designed to fool the currently running window manager to get desired
   results.

   Window manager developers may find 'ParentWindow' and 'Listen' useful.
*/
package xwindow

import "code.google.com/p/jamslam-x-go-binding/xgb"

import (
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xrect"
)

// Listen will tell X to report events corresponding to the event masks
// provided for the given window. If a call to Listen is omitted, you will
// not receive the events you desire.
func Listen(xu *xgbutil.XUtil, win xgb.Id, evMasks ...int) {
	evMask := 0
	for _, mask := range evMasks {
		evMask |= mask
	}

	xu.Conn().ChangeWindowAttributes(win, xgb.CWEventMask,
		[]uint32{uint32(evMask)})
}

// ParentWindow queries the QueryTree and finds the parent window.
func ParentWindow(xu *xgbutil.XUtil, win xgb.Id) (xgb.Id, error) {
	tree, err := xu.Conn().QueryTree(win)

	if err != nil {
		return 0, xgbutil.Xerr(err, "ParentWindow",
			"Error retrieving parent window for %x", win)
	}

	return tree.Parent, nil
}

// MoveResize is an accurate means of resizing a window, accounting for
// decorations. Usually, the x,y coordinates are fine---we just need to
// adjust the width and height.
func MoveResize(xu *xgbutil.XUtil, win xgb.Id, x, y, w, h int) error {
	neww, newh, err := adjustSize(xu, win, w, h)
	if err != nil {
		return err
	}

	return ewmh.MoveresizeWindowExtra(xu, win, x, y, neww, newh,
		xgb.GravityBitForget, 2, true, true)
}

// Move changes the position of a window without touching the size.
func Move(xu *xgbutil.XUtil, win xgb.Id, x, y int) error {
	return ewmh.MoveWindow(xu, win, x, y)
}

// Resize changes the size of a window without touching the position.
func Resize(xu *xgbutil.XUtil, win xgb.Id, w, h int) error {
	neww, newh, err := adjustSize(xu, win, w, h)
	if err != nil {
		return err
	}

	return ewmh.ResizeWindow(xu, win, neww, newh)
}

// adjustSize takes a client and dimensions, and adjust them so that they'll
// account for window decorations. For example, if you want a window to be
// 200 pixels wide, a window manager will typically determine that as
// you wanting the *client* to be 200 pixels wide. The end result is that
// the client plus decorations ends up being
// (200 + left decor width + right decor width) pixels wide. Which is probably
// not what you want. Therefore, transform 200 into
// 200 - decoration window width - client window width.
// Similarly for height.
func adjustSize(xu *xgbutil.XUtil, win xgb.Id, w, h int) (int, int, error) {
	cGeom, err := RawGeometry(xu, win) // raw client geometry
	if err != nil {
		return 0, 0, err
	}

	pGeom, err := GetGeometry(xu, win) // geometry with decorations
	if err != nil {
		return 0, 0, err
	}

	neww := w - (pGeom.Width() - cGeom.Width())
	newh := h - (pGeom.Height() - cGeom.Height())
	if neww < 1 {
		neww = 1
	}
	if newh < 1 {
		newh = 1
	}

	return neww, newh, nil
}

// GetGeometry retrieves the client's width and height *including* decorations.
// This can be tricky. In a non-parenting window manager, the width/height of
// a client can be found by inspecting the client directly. In a reparenting
// window manager like Openbox, the parent of the client reflects the true
// width/height. Still yet, in KWin, it's the parent of the parent of the
// client that reflects the true width/height.
// The idea then is to traverse up the tree until we hit the root window.
// Therefore, we're at a top-level window which should accurately reflect
// the width/height.
func GetGeometry(xu *xgbutil.XUtil, win xgb.Id) (xrect.Rect, error) {
	parent := win
	for {
		tempParent, err := ParentWindow(xu, parent)
		if err != nil || tempParent == xu.RootWin() {
			break
		}
		parent = tempParent
	}

	return RawGeometry(xu, parent)
}

// RawGeometry isn't smart. It just queries the window given for geometry.
func RawGeometry(xu *xgbutil.XUtil, win xgb.Id) (xrect.Rect, error) {
	xgeom, err := xu.Conn().GetGeometry(win)
	if err != nil {
		return nil, err
	}

	return xrect.New(int(xgeom.X), int(xgeom.Y),
		int(xgeom.Width), int(xgeom.Height)), nil
}
