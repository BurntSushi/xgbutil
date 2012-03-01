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
package xgbutil

import "code.google.com/p/x-go-binding/xgb"

// Geometry is a struct representing a window rectangle.
// The top left corner is the origin. Window coordinates could be
// negative, so watch out!
type Geometry struct {
    X, Y int32
    Width, Height uint32
}

// ParentWindow queries the QueryTree and finds the parent window.
func (xu *XUtil) ParentWindow(win xgb.Id) (xgb.Id, error) {
    tree, err := xu.conn.QueryTree(win)

    if err != nil {
        return 0, xerr(err, "ParentWindow",
                       "Error retrieving parent window for %x", win)
    }

    return tree.Parent, nil
}

// MoveResize is an accurate means of resizing a window, accounting for
// decorations. Usually, the x,y coordinates are fine---we just need to
// adjust the width and height.
func (xu *XUtil) MoveResize(win xgb.Id, x, y int32, w, h uint32) error {
    neww, newh, err := xu.adjustSize(win, w, h)
    if err != nil {
        return err
    }

    return xu.EwmhMoveresizeWindowExtra(win, uint32(x), uint32(y), neww, newh,
                                        xgb.GravityBitForget, 2, true, true)
}

// Move changes the position of a window without touching the size.
func (xu *XUtil) Move(win xgb.Id, x, y int32) error {
    return xu.EwmhMoveWindow(win, uint32(x), uint32(y))
}

// Resize changes the size of a window without touching the position.
func (xu *XUtil) Resize(win xgb.Id, w, h uint32) error {
    neww, newh, err := xu.adjustSize(win, w, h)
    if err != nil {
        return err
    }

    return xu.EwmhResizeWindow(win, neww, newh)
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
func (xu *XUtil) adjustSize(win xgb.Id, w, h uint32) (uint32, uint32, error) {
    cGeom, err := xu.RawGeometry(win) // raw client geometry
    if err != nil {
        return 0, 0, err
    }

    pGeom, err := xu.GetGeometry(win) // geometry with decorations
    if err != nil {
        return 0, 0, err
    }

    neww := w - (pGeom.Width - cGeom.Width)
    newh := h - (pGeom.Height - cGeom.Height)
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
func (xu *XUtil) GetGeometry(win xgb.Id) (*Geometry, error) {
    parent, err := xu.ParentWindow(win)
    for {
        tempParent, err := xu.ParentWindow(parent)
        if err != nil || tempParent == xu.root {
            break
        }
        parent = tempParent
    }
    if err != nil {
        return nil, err
    }

    return xu.RawGeometry(parent)
}

// RawGeometry isn't smart. It just queries the window given for geometry.
func (xu *XUtil) RawGeometry(win xgb.Id) (*Geometry, error) {
    xgeom, err := xu.conn.GetGeometry(win)
    if err != nil {
        return nil, err
    }

    return &Geometry{
        X: int32(xgeom.X), Y: int32(xgeom.Y),
        Width: uint32(xgeom.Width), Height: uint32(xgeom.Height),
    }, nil
}

