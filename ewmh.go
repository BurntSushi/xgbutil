/*
    An implementation of the entire EWMH spec.

    Since there are so many methods and they adhere to an existing spec, 
    this source file does not contain much documentation. Indeed, each
    method has only a single comment associated with it: the EWMH property name.

    See the EWMH spec for more info:
    http://standards.freedesktop.org/wm-spec/wm-spec-latest.html

    Here is the naming scheme using "_NET_ACTIVE_WINDOW" as an example.

    Methods "Ewmh_active_window" and "Ewmh_active_window_set" get and set the
    property, respectively. Both of these methods exist for all EWMH properties.
    Additionally, some EWMH properties support sending a client message event
    to request the window manager perform some action. In the case of
    "_NET_ACTIVE_WINDOW", this request is used to set the active window.
    These sorts of methods end in "_request". So for "_NET_ACTIVE_WINDOW",
    the method name is "Ewmh_active_window_request".

    For properties that store more than just a simple integer, name or list
    of integers, structs have been created and exposed to organized the
    information returned in a sensible manner. For example, the 
    "_NET_DESKTOP_GEOMETRY" property would typically return a slice of integers
    of length 2, where the first integer is the width and the second is the
    height. Xgbutil will wrap this in a struct with the obvious members. These
    structs are documented.

    Finally, methods ending in "_set" are typically only used when setting
    properties on clients *you've* created or when the window manager sets
    properties. Thus, it's unlikely that you should use them. Stick to the
    get methods and the "*_request" methods.

    N.B. Not all properties have "_request" methods.
*/
package xgbutil

import (

    "code.google.com/p/x-go-binding/xgb"
)

// _NET_ACTIVE_WINDOW
func (c *XUtil) EwmhActiveWindow() xgb.Id {
    return PropValId(c.GetProperty(c.root, "_NET_ACTIVE_WINDOW"))
}

// _NET_ACTIVE_WINDOW
func (c *XUtil) EwmhActiveWindowRequest(win xgb.Id) {
    evMask := (xgb.EventMaskSubstructureNotify |
               xgb.EventMaskSubstructureRedirect)
    data := make([]byte, 32)

    data[0] = xgb.ClientMessage
    data[1] = 32
    put32(data[4:], uint32(win))
    put32(data[8:], uint32(c.Atm("_NET_ACTIVE_WINDOW")))
    put32(data[12:], 1)
    put32(data[16:], 0)
    put32(data[20:], 0)

    c.conn.SendEvent(false, c.root, uint32(evMask), data)
}

// _NET_CLIENT_LIST
func (c *XUtil) EwmhClientList() []xgb.Id {
    return PropValIds(c.GetProperty(c.root, "_NET_CLIENT_LIST"))
}

// _NET_CURRENT_DESKTOP
func (c *XUtil) EwmhCurrentDesktop() uint32 {
    return PropValNum(c.GetProperty(c.root, "_NET_CURRENT_DESKTOP"))
}

// _NET_CURRENT_DESKTOP
func (c *XUtil) EwmhCurrentDesktopSet(desk uint32) {
    data := make([]byte, 4)
    put32(data, desk)
    c.conn.ChangeProperty(xgb.PropModeReplace, c.root,
                          c.Atm("_NET_CURRENT_DESKTOP"), xgb.AtomCardinal,
                          32, data)
}

// _NET_DESKTOP_NAMES
func (c *XUtil) EwmhDesktopNames() []string {
    return PropValStrs(c.GetProperty(c.root, "_NET_DESKTOP_NAMES"))
}

// DesktopGeometry is a struct that houses the width and height of a
// _NET_DESKTOP_GEOMETRY property reply.
type DesktopGeometry struct {
    Width uint32
    Height uint32
}

// _NET_DESKTOP_GEOMETRY
func (c *XUtil) EwmhDesktopGeometry() DesktopGeometry {
    geom := PropValNums(c.GetProperty(c.root, "_NET_DESKTOP_GEOMETRY"))

    return DesktopGeometry{Width: geom[0], Height: geom[1]}
}

// _NET_WM_NAME
func (c *XUtil) EwmhWmName(win xgb.Id) string {
    return PropValStr(c.GetProperty(win, "_NET_WM_NAME"))
}

