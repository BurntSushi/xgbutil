/*
    An implementation of the entire EWMH spec.

    Since there are so many methods and they adhere to an existing spec, 
    this source file does not contain much documentation. Indeed, each
    method has only a single comment associated with it: the EWMH property name.

    See the EWMH spec for more info:
    http://standards.freedesktop.org/wm-spec/wm-spec-latest.html

    Here is the naming scheme using "_NET_ACTIVE_WINDOW" as an example.

    Methods "EwmhActiveWindow" and "EwmhActiveWindowSet" get and set the
    property, respectively. Both of these methods exist for all EWMH properties.
    Additionally, some EWMH properties support sending a client message event
    to request the window manager perform some action. In the case of
    "_NET_ACTIVE_WINDOW", this request is used to set the active window.
    These sorts of methods end in "Req". So for "_NET_ACTIVE_WINDOW",
    the method name is "EwmhActiveWindowReq". Moreover, most requests include
    various parameters that don't need to be changed often (like the source
    indication). Thus, by default, methods ending in "Req" force these to
    sensible defaults. If you need access to all of the parameters, use the
    corresponding "ReqExtra" method. So for "_NET_ACTIVE_WINDOW", that would
    be "EwmhActiveWindowReqExtra".

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
func (xu *XUtil) EwmhActiveWindow() xgb.Id {
    return PropValId(xu.GetProperty(xu.root, "_NET_ACTIVE_WINDOW"))
}

// _NET_ACTIVE_WINDOW
func (xu *XUtil) EwmhActiveWindowReq(win xgb.Id) {
    xu.EwmhActiveWindowReqExtra(win, 2, 0, 0)
}

// _NET_ACTIVE_WINDOW
func (xu *XUtil) EwmhActiveWindowReqExtra(win xgb.Id, source uint32,
                                          time xgb.Timestamp,
                                          current_active xgb.Id) {
    cm := NewClientMessage(32, win, xu.Atm("_NET_ACTIVE_WINDOW"), source,
                           uint32(time), uint32(current_active))
    xu.SendRootEvent(cm)
}

// _NET_CLIENT_LIST
func (xu *XUtil) EwmhClientList() []xgb.Id {
    return PropValIds(xu.GetProperty(xu.root, "_NET_CLIENT_LIST"))
}

// _NET_CURRENT_DESKTOP
func (xu *XUtil) EwmhCurrentDesktop() uint32 {
    return PropValNum(xu.GetProperty(xu.root, "_NET_CURRENT_DESKTOP"))
}

// _NET_CURRENT_DESKTOP
func (xu *XUtil) EwmhCurrentDesktopSet(desk uint32) {
    xu.ChangeProperty32(xu.root, "_NET_CURRENT_DESKTOP", "CARDINAL", desk)
    // data := make([]byte, 4) 
    // put32(data, desk) 
    // xu.conn.ChangeProperty(xgb.PropModeReplace, xu.root, 
                           // xu.Atm("_NET_CURRENT_DESKTOP"), xgb.AtomCardinal, 
                           // 32, data) 
}

// _NET_DESKTOP_NAMES
func (xu *XUtil) EwmhDesktopNames() []string {
    return PropValStrs(xu.GetProperty(xu.root, "_NET_DESKTOP_NAMES"))
}

// _NET_DESKTOP_NAMES
func (xu *XUtil) EwmhDesktopNamesSet(names []string) {
    nullterm := make([]byte, 0)
    for _, name := range names {
        nullterm = append(nullterm, name...)
        nullterm = append(nullterm, 0)
    }
    xu.ChangeProperty(xu.root, 8, "_NET_DESKTOP_NAMES", "UTF8_STRING", nullterm)
}

// DesktopGeometry is a struct that houses the width and height of a
// _NET_DESKTOP_GEOMETRY property reply.
type DesktopGeometry struct {
    Width uint32
    Height uint32
}

// _NET_DESKTOP_GEOMETRY
func (xu *XUtil) EwmhDesktopGeometry() DesktopGeometry {
    geom := PropValNums(xu.GetProperty(xu.root, "_NET_DESKTOP_GEOMETRY"))

    return DesktopGeometry{Width: geom[0], Height: geom[1]}
}

// WmIcon is a struct that contains data for a single icon.
// The EwmhWmIcon method will return a list of these, since a single
// client can specify multiple icons of varying sizes.
type WmIcon struct {
    Width uint32
    Height uint32
    Data []uint32
}

// _NET_WM_ICON
func (xu *XUtil) EwmhWmIcon(win xgb.Id) []WmIcon {
    icon := PropValNums(xu.GetProperty(win, "_NET_WM_ICON"))

    wmicons := make([]WmIcon, 0)
    start := uint32(0)
    for int(start) < len(icon) {
        w, h := icon[start], icon[start + 1]
        upto := w * h

        wmicon := WmIcon{
            Width: w,
            Height: h,
            Data: icon[(start + 2):(start + upto + 2)],
        }
        wmicons = append(wmicons, wmicon)

        start += upto + 2
    }

    return wmicons
}

// _NET_WM_NAME
func (xu *XUtil) EwmhWmName(win xgb.Id) string {
    return PropValStr(xu.GetProperty(win, "_NET_WM_NAME"))
}

// _NET_WM_NAME
func (xu *XUtil) EwmhWmNameSet(win xgb.Id, name string) {
    xu.ChangeProperty(win, 8, "_NET_WM_NAME", "UTF8_STRING", []byte(name))
}

