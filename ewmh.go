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
    be "EwmhActiveWindowReqExtra". (If no "ReqExtra" method exists, then the
    "Req" method covers all available parameters.)

    For properties that store more than just a simple integer, name or list
    of integers, structs have been created and exposed to organized the
    information returned in a sensible manner. For example, the 
    "_NET_DESKTOP_GEOMETRY" property would typically return a slice of integers
    of length 2, where the first integer is the width and the second is the
    height. Xgbutil will wrap this in a struct with the obvious members. These
    structs are documented.

    Finally, methods ending in "*Set" are typically only used when setting
    properties on clients *you've* created or when the window manager sets
    properties. Thus, it's unlikely that you should use them. Stick to the
    get methods and the "*Req" methods.

    N.B. Not all properties have "*Req" methods.
*/
package xgbutil

import (
    "code.google.com/p/x-go-binding/xgb"
)

// EwmhClientEvent is a convenience function that sends ClientMessage events
// to the root window as specified by the EWMH spec.
func (xu *XUtil) EwmhClientEvent(window xgb.Id, message_type string,
                                 data ...interface{}) {
    evMask := (xgb.EventMaskSubstructureNotify |
               xgb.EventMaskSubstructureRedirect)
    cm := NewClientMessage(32, window, xu.Atm(message_type), data...)
    xu.SendRootEvent(cm, uint32(evMask))
}

// _NET_ACTIVE_WINDOW get
func (xu *XUtil) EwmhActiveWindow() xgb.Id {
    return PropValId(xu.GetProperty(xu.root, "_NET_ACTIVE_WINDOW"))
}

// _NET_ACTIVE_WINDOW set
func (xu *XUtil) EwmhActiveWindowSet(win xgb.Id) {
    xu.ChangeProperty32(xu.root, "_NET_ACTIVE_WINDOW", "WINDOW", uint32(win))
}

// _NET_ACTIVE_WINDOW req
func (xu *XUtil) EwmhActiveWindowReq(win xgb.Id) {
    xu.EwmhActiveWindowReqExtra(win, 2, 0, 0)
}

// _NET_ACTIVE_WINDOW req extra
func (xu *XUtil) EwmhActiveWindowReqExtra(win xgb.Id, source uint32,
                                          time xgb.Timestamp,
                                          current_active xgb.Id) {
    xu.EwmhClientEvent(win, "_NET_ACTIVE_WINDOW", source, uint32(time),
                       uint32(current_active))
}

// _NET_CURRENT_DESKTOP get
func (xu *XUtil) EwmhCurrentDesktop() uint32 {
    return PropValNum(xu.GetProperty(xu.root, "_NET_CURRENT_DESKTOP"))
}

// _NET_CURRENT_DESKTOP set
func (xu *XUtil) EwmhCurrentDesktopSet(desk uint32) {
    xu.ChangeProperty32(xu.root, "_NET_CURRENT_DESKTOP", "CARDINAL", desk)
}

// _NET_CURRENT_DESKTOP req
func (xu *XUtil) EwmhCurrentDesktopReq(desk uint32) {
    xu.EwmhClientEvent(xu.root, "_NET_CURRENT_DESKTOP", desk)
}

// _NET_CURRENT_DESKTOP req extra
func (xu *XUtil) EwmhCurrentDesktopReqExtra(desk uint32, time xgb.Timestamp) {
    xu.EwmhClientEvent(xu.root, "_NET_CURRENT_DESKTOP", desk, time)
}

// _NET_CLIENT_LIST get
func (xu *XUtil) EwmhClientList() []xgb.Id {
    return PropValIds(xu.GetProperty(xu.root, "_NET_CLIENT_LIST"))
}

// _NET_CLIENT_LIST set
func (xu *XUtil) EwmhClientListSet(wins []xgb.Id) {
    xu.ChangeProperty32(xu.root, "_NET_CLIENT_LIST", "WINDOW",
                        IdTo32(wins)...)
}

// _NET_CLIENT_LIST_STACKING get
func (xu *XUtil) EwmhClientListStacking() []xgb.Id {
    return PropValIds(xu.GetProperty(xu.root, "_NET_CLIENT_LIST_STACKING"))
}

// _NET_CLIENT_LIST_STACKING set
func (xu *XUtil) EwmhClientListStackingSet(wins []xgb.Id) {
    xu.ChangeProperty32(xu.root, "_NET_CLIENT_LIST_STACKING", "WINDOW",
                        IdTo32(wins)...)
}

// _NET_DESKTOP_NAMES get
func (xu *XUtil) EwmhDesktopNames() []string {
    return PropValStrs(xu.GetProperty(xu.root, "_NET_DESKTOP_NAMES"))
}

// _NET_DESKTOP_NAMES set
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

// _NET_DESKTOP_GEOMETRY get
func (xu *XUtil) EwmhDesktopGeometry() DesktopGeometry {
    geom := PropValNums(xu.GetProperty(xu.root, "_NET_DESKTOP_GEOMETRY"))

    return DesktopGeometry{Width: geom[0], Height: geom[1]}
}

// _NET_DESKTOP_GEOMETRY set
func (xu *XUtil) EwmhDesktopGeometrySet(dg DesktopGeometry) {
    xu.ChangeProperty32(xu.root, "_NET_DESKTOP_GEOMETRY", "CARDINAL",
                        dg.Width, dg.Height)
}

// _NET_DESKTOP_GEOMETRY req
func (xu *XUtil) EwmhDesktopGeometryReq(dg DesktopGeometry) {
    xu.EwmhClientEvent(xu.root, "_NET_DESKTOP_GEOMETRY", dg.Width, dg.Height)
}

// DesktopViewport is a struct that contains a pairing of x,y coordinates
// representing the top-left corner of each desktop. (There will typically
// be one struct here for each desktop in existence.)
type DesktopViewport struct {
    X uint32
    Y uint32
}

// _NET_DESKTOP_VIEWPORT get
func (xu *XUtil) EwmhDesktopViewport() []DesktopViewport {
    coords := PropValNums(xu.GetProperty(xu.root, "_NET_DESKTOP_VIEWPORT"))
    viewports := make([]DesktopViewport, len(coords) / 2)

    for i, _ := range viewports {
        viewports[i] = DesktopViewport{
            X: coords[i * 2],
            Y: coords[i * 2 + 1],
        }
    }

    return viewports
}

// _NET_DESKTOP_VIEWPORT set
func (xu *XUtil) EwmhDesktopViewportSet(viewports []DesktopViewport) {
    coords := make([]uint32, len(viewports) * 2)
    for i, viewport := range viewports {
        coords[i * 2] = viewport.X
        coords[i * 2 + 1] = viewport.Y
    }

    xu.ChangeProperty32(xu.root, "_NET_DESKTOP_VIEWPORT", "CARDINAL", coords...)
}

// _NET_DESKTOP_VIEWPORT req
func (xu *XUtil) EwmhDesktopViewportReq(x uint32, y uint32) {
    xu.EwmhClientEvent(xu.root, "_NET_DESKTOP_VIEWPORT", x, y)
}

// _NET_NUMBER_OF_DESKTOPS get
func (xu *XUtil) EwmhNumberOfDesktops() uint32 {
    return PropValNum(xu.GetProperty(xu.root, "_NET_NUMBER_OF_DESKTOPS"))
}

// _NET_NUMBER_OF_DESKTOPS set
func (xu *XUtil) EwmhNumberOfDesktopsSet(numDesks uint32) {
    xu.ChangeProperty32(xu.root, "_NET_NUMBER_OF_DESKTOPS", "CARDINAL",
                        numDesks)
}

// _NET_NUMBER_OF_DESKTOPS req
func (xu *XUtil) EwmhNumberOfDesktopsReq(numDesks uint32) {
    xu.EwmhClientEvent(xu.root, "_NET_NUMBER_OF_DESKTOPS", numDesks)
}

// _NET_SUPPORTED get
func (xu *XUtil) EwmhSupported() []string {
    return xu.PropValAtoms(xu.GetProperty(xu.root, "_NET_SUPPORTED"))
}

// _NET_SUPPORTED set
// This will create any atoms in the argument if they don't already exist.
func (xu *XUtil) EwmhSupportedSet(atomNames []string) {
    atoms := make([]uint32, len(atomNames))
    for i, atomName := range atomNames {
        atoms[i] = uint32(xu.Atom(atomName, false))
    }

    xu.ChangeProperty32(xu.root, "_NET_SUPPORTED", "ATOM", atoms...)
}

// _NET_SUPPORTING_WM_CHECK get
func (xu *XUtil) EwmhSupportingWmCheck(win xgb.Id) xgb.Id {
    return PropValId(xu.GetProperty(win, "_NET_SUPPORTING_WM_CHECK"))
}

// _NET_SUPPORTING_WM_CHECK set
func (xu *XUtil) EwmhSupportingWmCheckSet(win xgb.Id, wm_win xgb.Id) {
    xu.ChangeProperty32(win, "_NET_SUPPORTING_WM_CHECK", "WINDOW",
                        uint32(wm_win))
}

// _NET_VISIBLE_DESKTOPS get
// This is not parted of the EWMH spec, but is a property of my own creation.
// It allows the window manager to report that it has multiple desktops
// viewable at the same time. (This conflicts with other EWMH properties,
// so I don't think this will ever be added to the official spec.)
func (xu *XUtil) EwmhVisibleDesktops() []uint32 {
    return PropValNums(xu.GetProperty(xu.root, "_NET_VISIBLE_DESKTOPS"))
}

// _NET_VISIBLE_DESKTOPS set
func (xu *XUtil) EwmhVisibleDesktopsSet(desktops []uint32) {
    xu.ChangeProperty32(xu.root, "_NET_VISIBLE_DESKTOPS", "CARDINAL",
                        desktops...)
}

// WmIcon is a struct that contains data for a single icon.
// The EwmhWmIcon method will return a list of these, since a single
// client can specify multiple icons of varying sizes.
type WmIcon struct {
    Width uint32
    Height uint32
    Data []uint32
}

// _NET_WM_ICON get
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

// _NET_WM_NAME get
func (xu *XUtil) EwmhWmName(win xgb.Id) string {
    return PropValStr(xu.GetProperty(win, "_NET_WM_NAME"))
}

// _NET_WM_NAME set
func (xu *XUtil) EwmhWmNameSet(win xgb.Id, name string) {
    xu.ChangeProperty(win, 8, "_NET_WM_NAME", "UTF8_STRING", []byte(name))
}

// Workarea is a struct that represents a rectangle as a bounding box of
// a single desktop. So there should be as many Workarea structs as there
// are desktops.
type Workarea struct {
    X uint32
    Y uint32
    Width uint32
    Height uint32
}

// _NET_WORKAREA get
func (xu *XUtil) EwmhWorkarea() []Workarea {
    rects := PropValNums(xu.GetProperty(xu.root, "_NET_WORKAREA"))
    workareas := make([]Workarea, len(rects) / 4)

    for i, _ := range workareas {
        workareas[i] = Workarea {
            X: rects[i * 4],
            Y: rects[i * 4 + 1],
            Width: rects[i * 4 + 2],
            Height: rects[i * 4 + 3],
        }
    }

    return workareas
}

// _NET_WORKAREA set
func (xu *XUtil) EwmhWorkareaSet(workareas []Workarea) {
    rects := make([]uint32, len(workareas) * 4)
    for i, workarea := range workareas {
        rects[i * 4] = workarea.X
        rects[i * 4 + 1] = workarea.Y
        rects[i * 4 + 2] = workarea.Width
        rects[i * 4 + 3] = workarea.Height
    }

    xu.ChangeProperty32(xu.root, "_NET_WORKAREA", "CARDINAL", rects...)
}

