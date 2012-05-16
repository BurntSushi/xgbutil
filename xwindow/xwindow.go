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

import (
	"fmt"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xrect"
)

// Window represents an X window. It contains an XUtilValue to simply the
// parameter lists for methods declared on the Window type.
// Geom is updated whenever Geometry is called, or when Move, Resize or
// MoveResize are called.
type Window struct {
	X    *xgbutil.XUtil
	Id   xproto.Window
	Geom xrect.Rect
}

// New creates a new window value from a window id and an XUtil type.
// Geom is initialize to nil. Use Window.Geometry to load it.
// Note that the geometry is the size of this particular window and nothing
// else. If you want the geometry of a client window including decorations,
// please use Window.DecorGeometry.
func New(xu *xgbutil.XUtil, win xproto.Window) *Window {
	return &Window{
		X:    xu,
		Id:   win,
		Geom: nil,
	}
}

// Generate is just like New, but generates a new X resource id for you.
// Geom is set to nil.
// It is possible for id generation to return an error, in which case, an
// error is returned here.
func Generate(xu *xgbutil.XUtil) (*Window, error) {
	wid, err := xproto.NewWindowId(xu.Conn())
	if err != nil {
		return nil, err
	}
	return &Window{
		X:    xu,
		Id:   wid,
		Geom: nil,
	}, nil
}

// Create issues a CreateWindow request for Window.
// Its purpose is to omit several boiler-plate parameters to CreateWindow
// and expose the commonly useful ones.
// The value mask describes which values are present in valueList.
// Value masks can be found in xgb/xproto with the prefix 'Cw'.
// The value list must contain values in the same order as the constants
// are defined in xgb/xproto.
//
// For example, the following creates a window positioned at (20, 50) with
// width 500 and height 700 with a background color of white.
//
//	w, err := xwindow.Generate(X)
//	if err != nil {
//		log.Fatalf("Could not generate a new resource identifier: %s", err)
//	}
//	w.Create(X.RootWin(), 20, 50, 500, 700,
//		xproto.CwBackPixel, 0xffffff)
func (w *Window) Create(parent xproto.Window, x, y, width, height,
	valueMask int, valueList ...uint32) {

	s := w.X.Screen()
	xproto.CreateWindow(w.X.Conn(), s.RootDepth, w.Id, parent,
		int16(x), int16(y), uint16(width), uint16(height), 0,
		xproto.WindowClassInputOutput, s.RootVisual,
		uint32(valueMask), valueList)
}

// CreateChecked issues a CreateWindow checked request for Window.
// A checked request is a synchronous request. Meaning that if the request
// fails, you can get the error returned to you. However, it also forced your
// program to block for a round trip to the X server, so it is slower.
// See the docs for Create for more info.
func (w *Window) CreateChecked(parent xproto.Window, x, y, width, height,
	valueMask int, valueList ...uint32) error {

	s := w.X.Screen()
	return xproto.CreateWindowChecked(w.X.Conn(), s.RootDepth, w.Id, parent,
		int16(x), int16(y), uint16(width), uint16(height), 0,
		xproto.WindowClassInputOutput, s.RootVisual,
		uint32(valueMask), valueList).Check()
}

// Change issues a ChangeWindowAttributes request with the provide mask
// and value list. Please see Window.Create for an example on how to use
// the mask and value list.
func (w *Window) Change(valueMask int, valueList ...uint32) {
	xproto.ChangeWindowAttributes(w.X.Conn(), w.Id,
		uint32(valueMask), valueList)
}

// Listen will tell X to report events corresponding to the event masks
// provided for the given window. If a call to Listen is omitted, you will
// not receive the events you desire.
// Event masks are constants declare in the xgb/xproto package starting with the
// EventMask prefix.
func (w *Window) Listen(evMasks ...int) {
	evMask := 0
	for _, mask := range evMasks {
		evMask |= mask
	}

	xproto.ChangeWindowAttributes(w.X.Conn(), w.Id, xproto.CwEventMask,
		[]uint32{uint32(evMask)})
}

// Geometry retrieves an up-to-date version of the this window's geometry.
// It also loads the geometry into the Geom member of Window.
func (w *Window) Geometry() (xrect.Rect, error) {
	geom, err := rawGeometry(w.X, xproto.Drawable(w.Id))
	if err != nil {
		return nil, err
	}
	w.Geom = geom
	return geom, err
}

// rawGeometry isn't smart. It just queries the window given for geometry.
func rawGeometry(xu *xgbutil.XUtil, win xproto.Drawable) (xrect.Rect, error) {
	xgeom, err := xproto.GetGeometry(xu.Conn(), win).Reply()
	if err != nil {
		return nil, err
	}
	return xrect.New(int(xgeom.X), int(xgeom.Y),
		int(xgeom.Width), int(xgeom.Height)), nil
}

// MoveResize issues a ConfigureRequest for this window with the provided
// x, y, width and height. Note that if width or height is 0, X will stomp
// all over you. Really hard. Don't do it.
// If you're trying to move/resize a top-level window in a window manager that
// supports EWMH, please use WMMoveResize instead.
func (w *Window) MoveResize(x, y, width, height int) {
	w.Geom.XSet(x)
	w.Geom.YSet(y)
	w.Geom.WidthSet(width)
	w.Geom.HeightSet(height)

	xproto.ConfigureWindow(w.X.Conn(), w.Id,
		xproto.ConfigWindowX|xproto.ConfigWindowY|
			xproto.ConfigWindowWidth|xproto.ConfigWindowHeight,
		[]uint32{uint32(x), uint32(y), uint32(width), uint32(height)})
}

// Move issues a ConfigureRequest for this window with the provided
// x and y positions.
// If you're trying to move a top-level window in a window manager that
// supports EWMH, please use WMMove instead.
func (w *Window) Move(x, y int) {
	w.Geom.XSet(x)
	w.Geom.YSet(y)

	xproto.ConfigureWindow(w.X.Conn(), w.Id,
		xproto.ConfigWindowX|xproto.ConfigWindowY,
		[]uint32{uint32(x), uint32(y)})
}

// Resize issues a ConfigureRequest for this window with the provided
// width and height. Note that if width or height is 0, X will stomp
// all over you. Really hard. Don't do it.
// If you're trying to resize a top-level window in a window manager that
// supports EWMH, please use WMResize instead.
func (w *Window) Resize(width, height int) {
	w.Geom.WidthSet(width)
	w.Geom.HeightSet(height)

	xproto.ConfigureWindow(w.X.Conn(), w.Id,
		xproto.ConfigWindowWidth|xproto.ConfigWindowHeight,
		[]uint32{uint32(width), uint32(height)})
}

// Stack issues a configure request to change the stack mode of Window.
// If you're using a window manager that supports EWMH, you may want to try
// and use ewmh.RestackWindow instead. Although this should still work.
// 'mode' values can be found as constants in xgb/xproto with the prefix
// StackMode.
// A value of xproto.StackModeAbove will put the window to the top of the stack,
// while a value of xproto.StackMoveBelow will put the window to the
// bottom of the stack.
// Remember that stacking is at the discretion of the window manager, and
// therefore may not always work as one would expect.
func (w *Window) Stack(mode byte) {
	xproto.ConfigureWindow(w.X.Conn(), w.Id,
		xproto.ConfigWindowStackMode, []uint32{uint32(mode)})
}

// StackSibling issues a configure request to change the sibling and stack mode
// of Window.
// If you're using a window manager that supports EWMH, you may want to try
// and use ewmh.RestackWindowExtra instead. Although this should still work.
// 'mode' values can be found as constants in xgb/xproto with the prefix
// StackMode.
// 'sibling' refers to the sibling window in the stacking order through which
// 'mode' is interpreted. Note that 'sibling' should be taken literally. A
// window can only be stacked with respect to a *sibling* in the window tree.
// This means that a client window that has been wrapped in decorations cannot
// be stacked with respect to another client window. (This is why you should
// use ewmh.RestackWindowExtra instead.)
func (w *Window) StackSibling(sibling xproto.Window, mode byte) {
	xproto.ConfigureWindow(w.X.Conn(), w.Id,
		xproto.ConfigWindowSibling|xproto.ConfigWindowStackMode,
		[]uint32{uint32(sibling), uint32(mode)})
}

// Map is a simple alias to map the window.
func (w *Window) Map() {
	xproto.MapWindow(w.X.Conn(), w.Id)
}

// Unamp is a simple alias to unmap the window.
func (w *Window) Unmap() {
	xproto.UnmapWindow(w.X.Conn(), w.Id)
}

// Destroy is a simple alias to destroy a window. You should use this when
// you no longer intend to use this window. (It will free the resource
// identifier for use in other places.)
func (w *Window) Destroy() {
	xproto.DestroyWindow(w.X.Conn(), w.Id)
}

// Focus tries to issue a SetInputFocus to get the focus.
// If you're trying to change the top-level active window, please use
// ewmh.ActiveWindowReq instead.
func (w *Window) Focus() {
	xproto.SetInputFocus(w.X.Conn(), xproto.InputFocusPointerRoot, w.Id, 0)
}

// Kill forcefully destroys a client. It is almost never what you want, and if
// you do it to one your clients, you'll lose your connection.
// (This is typically used in a special client like `xkill` or in a window
// manager.)
func (w *Window) Kill() {
	xproto.KillClient(w.X.Conn(), uint32(w.Id))
}

// Clear paints the region of the window specified with the corresponding
// background pixmap. If the window doesn't have a background pixmap,
// this has no effect.
// If width/height is 0, then it is set to the width/height of the background
// pixmap minus x/y.
func (w *Window) Clear(x, y, width, height int) {
	xproto.ClearArea(w.X.Conn(), false, w.Id,
		int16(x), int16(y), uint16(width), uint16(height))
}

// ClearAll is the same as Clear, but does it for the entire background pixmap.
func (w *Window) ClearAll() {
	xproto.ClearArea(w.X.Conn(), false, w.Id, 0, 0, 0, 0)
}

// Parent queries the QueryTree and finds the parent window.
func (w *Window) Parent() (*Window, error) {
	tree, err := xproto.QueryTree(w.X.Conn(), w.Id).Reply()
	if err != nil {
		return nil, fmt.Errorf("ParentWindow: Error retrieving parent window "+
			"for %x: %s", w.Id, err)
	}
	return New(w.X, tree.Parent), nil
}
