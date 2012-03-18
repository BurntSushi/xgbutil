/*
    An API to the entire EWMH spec.

    Since there are so many methods and they adhere to an existing spec, 
    this source file does not contain much documentation. Indeed, each
    method has only a single comment associated with it: the EWMH property name.

    See the EWMH spec for more info:
    http://standards.freedesktop.org/wm-spec/wm-spec-latest.html

    Here is the naming scheme using "_NET_ACTIVE_WINDOW" as an example.

    Methods "ActiveWindowGet" and "ActiveWindowSet" get and set the
    property, respectively. Both of these methods exist for most EWMH 
    properties.  Additionally, some EWMH properties support sending a client 
    message event to request the window manager perform some action. In the 
    case of "_NET_ACTIVE_WINDOW", this request is used to set the active 
    window.

    These sorts of methods end in "Req". So for "_NET_ACTIVE_WINDOW",
    the method name is "ActiveWindowReq". Moreover, most requests include
    various parameters that don't need to be changed often (like the source
    indication). Thus, by default, methods ending in "Req" force these to
    sensible defaults. If you need access to all of the parameters, use the
    corresponding "ReqExtra" method. So for "_NET_ACTIVE_WINDOW", that would
    be "ActiveWindowReqExtra". (If no "ReqExtra" method exists, then the
    "Req" method covers all available parameters.)

    This naming scheme has one exception: if a property's only use is through
    sending an event (like "_NET_CLOSE_WINDOW"), then the name will be
    "CloseWindow" for the short-hand version and "CloseWindowExtra"
    for access to all of the parameters.

    For properties that store more than just a simple integer, name or list
    of integers, structs have been created and exposed to organize the
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
package ewmh

import "code.google.com/p/jamslam-x-go-binding/xgb"
import "github.com/BurntSushi/xgbutil"
import "github.com/BurntSushi/xgbutil/xevent"
import "github.com/BurntSushi/xgbutil/xprop"

// ClientEvent is a convenience function that sends ClientMessage events
// to the root window as specified by the EWMH spec.
func ClientEvent(xu *xgbutil.XUtil, window xgb.Id, message_type string,
                 data ...interface{}) error {
    mstype, err := xprop.Atm(xu, message_type)
    if err != nil {
        return err
    }

    evMask := (xgb.EventMaskSubstructureNotify |
               xgb.EventMaskSubstructureRedirect)
    cm, err := xevent.NewClientMessage(32, window, mstype, data...)
    if err != nil {
        return err
    }

    xevent.SendRootEvent(xu, cm, uint32(evMask))
    return nil
}

// _NET_ACTIVE_WINDOW get
func ActiveWindowGet(xu *xgbutil.XUtil) (xgb.Id, error) {
    return xprop.PropValId(xprop.GetProperty(xu, xu.RootWin(),
                                             "_NET_ACTIVE_WINDOW"))
}

// _NET_ACTIVE_WINDOW set
func ActiveWindowSet(xu *xgbutil.XUtil, win xgb.Id) error {
    return xprop.ChangeProp32(xu, xu.RootWin(), "_NET_ACTIVE_WINDOW", "WINDOW",
                              uint32(win))
}

// _NET_ACTIVE_WINDOW req
func ActiveWindowReq(xu *xgbutil.XUtil, win xgb.Id) error {
    return ActiveWindowReqExtra(xu, win, 2, 0, 0)
}

// _NET_ACTIVE_WINDOW req extra
func ActiveWindowReqExtra(xu *xgbutil.XUtil, win xgb.Id, source uint32,
                          time xgb.Timestamp, current_active xgb.Id) error {
    return ClientEvent(xu, win, "_NET_ACTIVE_WINDOW", source, uint32(time),
                       uint32(current_active))
}

// _NET_CLIENT_LIST get
func ClientListGet(xu *xgbutil.XUtil) ([]xgb.Id, error) {
    return xprop.PropValIds(xprop.GetProperty(xu, xu.RootWin(),
                                              "_NET_CLIENT_LIST"))
}

// _NET_CLIENT_LIST set
func ClientListSet(xu *xgbutil.XUtil, wins []xgb.Id) error {
    return xprop.ChangeProp32(xu, xu.RootWin(), "_NET_CLIENT_LIST", "WINDOW",
                              xprop.IdTo32(wins)...)
}

// _NET_CLIENT_LIST_STACKING get
func ClientListStackingGet(xu *xgbutil.XUtil) ([]xgb.Id, error) {
    return xprop.PropValIds(xprop.GetProperty(xu, xu.RootWin(),
                                              "_NET_CLIENT_LIST_STACKING"))
}

// _NET_CLIENT_LIST_STACKING set
func ClientListStackingSet(xu *xgbutil.XUtil, wins []xgb.Id) error {
    return xprop.ChangeProp32(xu, xu.RootWin(), "_NET_CLIENT_LIST_STACKING",
                              "WINDOW", xprop.IdTo32(wins)...)
}

// _NET_CLOSE_WINDOW req
func CloseWindow(xu *xgbutil.XUtil, win xgb.Id) error {
    return CloseWindowExtra(xu, win, 0, 2)
}

// _NET_CLOSE_WINDOW req extra
func CloseWindowExtra(xu *xgbutil.XUtil, win xgb.Id, time xgb.Timestamp,
                      source uint32) error {
    return ClientEvent(xu, win, "_NET_CLOSE_WINDOW", uint32(time), source)
}

// _NET_CURRENT_DESKTOP get
func CurrentDesktopGet(xu *xgbutil.XUtil) (uint32, error) {
    return xprop.PropValNum(xprop.GetProperty(xu, xu.RootWin(),
                                              "_NET_CURRENT_DESKTOP"))
}

// _NET_CURRENT_DESKTOP set
func CurrentDesktopSet(xu *xgbutil.XUtil, desk uint32) error {
    return xprop.ChangeProp32(xu, xu.RootWin(), "_NET_CURRENT_DESKTOP",
                              "CARDINAL", desk)
}

// _NET_CURRENT_DESKTOP req
func CurrentDesktopReq(xu *xgbutil.XUtil, desk uint32) error {
    return CurrentDesktopReqExtra(xu, desk, 0)
}

// _NET_CURRENT_DESKTOP req extra
func CurrentDesktopReqExtra(xu *xgbutil.XUtil, desk uint32,
                            time xgb.Timestamp) error {
    return ClientEvent(xu, xu.RootWin(), "_NET_CURRENT_DESKTOP", desk, time)
}

// _NET_DESKTOP_NAMES get
func DesktopNamesGet(xu *xgbutil.XUtil) ([]string, error) {
    return xprop.PropValStrs(xprop.GetProperty(xu, xu.RootWin(),
                                               "_NET_DESKTOP_NAMES"))
}

// _NET_DESKTOP_NAMES set
func DesktopNamesSet(xu *xgbutil.XUtil, names []string) error {
    nullterm := make([]byte, 0)
    for _, name := range names {
        nullterm = append(nullterm, name...)
        nullterm = append(nullterm, 0)
    }
    return xprop.ChangeProp(xu, xu.RootWin(), 8, "_NET_DESKTOP_NAMES",
                            "UTF8_STRING", nullterm)
}

// DesktopGeometry is a struct that houses the width and height of a
// _NET_DESKTOP_GEOMETRY property reply.
type DesktopGeometry struct {
    Width uint32
    Height uint32
}

// _NET_DESKTOP_GEOMETRY get
func DesktopGeometryGet(xu *xgbutil.XUtil) (DesktopGeometry, error) {
    geom, err := xprop.PropValNums(xprop.GetProperty(xu, xu.RootWin(),
                                                     "_NET_DESKTOP_GEOMETRY"))
    if err != nil {
        return DesktopGeometry{}, err
    }

    return DesktopGeometry{Width: geom[0], Height: geom[1]}, nil
}

// _NET_DESKTOP_GEOMETRY set
func DesktopGeometrySet(xu *xgbutil.XUtil, dg DesktopGeometry) error {
    return xprop.ChangeProp32(xu, xu.RootWin(), "_NET_DESKTOP_GEOMETRY",
                              "CARDINAL", dg.Width, dg.Height)
}

// _NET_DESKTOP_GEOMETRY req
func DesktopGeometryReq(xu *xgbutil.XUtil, dg DesktopGeometry) error {
    return ClientEvent(xu, xu.RootWin(), "_NET_DESKTOP_GEOMETRY", dg.Width,
                       dg.Height)
}

// DesktopLayout is a struct that organizes information pertaining to
// the _NET_DESKTOP_LAYOUT property. Namely, the orientation, the number
// of columns, the number of rows, and the starting corner.
type DesktopLayout struct {
    Orientation uint32
    Columns uint32
    Rows uint32
    StartingCorner uint32
}

// _NET_DESKTOP_LAYOUT constants for orientation
const (
    OrientHorz = iota
    OrientVert
)

// _NET_DESKTOP_LAYOUT constants for starting corner
const (
    TopLeft = iota
    TopRight
    BottomRight
    BottomLeft
)

// _NET_DESKTOP_LAYOUT get
func DesktopLayoutGet(xu *xgbutil.XUtil) (dl DesktopLayout, err error) {
    dlraw, err := xprop.PropValNums(xprop.GetProperty(xu, xu.RootWin(),
                                                      "_NET_DESKTOP_LAYOUT"))
    if err != nil {
        return DesktopLayout{}, err
    }

    dl.Orientation = dlraw[0]
    dl.Columns = dlraw[1]
    dl.Rows = dlraw[2]

    if len(dlraw) > 3 {
        dl.StartingCorner = dlraw[3]
    } else {
        dl.StartingCorner = TopLeft
    }

    return dl, nil
}

// _NET_DESKTOP_LAYOUT set
func DesktopLayoutSet(xu *xgbutil.XUtil, orientation, columns, rows,
                      startingCorner uint32) error {
    return xprop.ChangeProp32(xu, xu.RootWin(), "_NET_DESKTOP_LAYOUT",
                              "CARDINAL", orientation, columns, rows,
                              startingCorner)
}

// DesktopViewport is a struct that contains a pairing of x,y coordinates
// representing the top-left corner of each desktop. (There will typically
// be one struct here for each desktop in existence.)
type DesktopViewport struct {
    X uint32
    Y uint32
}

// _NET_DESKTOP_VIEWPORT get
func DesktopViewportGet(xu *xgbutil.XUtil) ([]DesktopViewport, error) {
    coords, err := xprop.PropValNums(xprop.GetProperty(xu, xu.RootWin(),
                                                       "_NET_DESKTOP_VIEWPORT"))
    if err != nil {
        return nil, err
    }

    viewports := make([]DesktopViewport, len(coords) / 2)
    for i, _ := range viewports {
        viewports[i] = DesktopViewport{
            X: coords[i * 2],
            Y: coords[i * 2 + 1],
        }
    }
    return viewports, nil
}

// _NET_DESKTOP_VIEWPORT set
func DesktopViewportSet(xu *xgbutil.XUtil, viewports []DesktopViewport) error {
    coords := make([]uint32, len(viewports) * 2)
    for i, viewport := range viewports {
        coords[i * 2] = viewport.X
        coords[i * 2 + 1] = viewport.Y
    }

    return xprop.ChangeProp32(xu, xu.RootWin(), "_NET_DESKTOP_VIEWPORT",
                              "CARDINAL", coords...)
}

// _NET_DESKTOP_VIEWPORT req
func DesktopViewportReq(xu *xgbutil.XUtil, x uint32, y uint32) error {
    return ClientEvent(xu, xu.RootWin(), "_NET_DESKTOP_VIEWPORT", x, y)
}

// FrameExtents is a struct that organizes information associated with
// the _NET_FRAME_EXTENTS property. Namely, the left, right, top and bottom
// decoration sizes.
type FrameExtents struct {
    Left uint32
    Right uint32
    Top uint32
    Bottom uint32
}

// _NET_FRAME_EXTENTS get
func FrameExtentsGet(xu *xgbutil.XUtil, win xgb.Id) (FrameExtents, error) {
    raw, err := xprop.PropValNums(xprop.GetProperty(xu, win,
                                                    "_NET_FRAME_EXTENTS"))
    if err != nil {
        return FrameExtents{}, nil
    }

    return FrameExtents{
        Left: raw[0],
        Right: raw[1],
        Top: raw[2],
        Bottom: raw[3],
    }, nil
}

// _NET_FRAME_EXTENTS set
func FrameExtentsSet(xu *xgbutil.XUtil, win xgb.Id, extents FrameExtents) error {
    raw := make([]uint32, 4)
    raw[0] = extents.Left
    raw[1] = extents.Right
    raw[2] = extents.Top
    raw[3] = extents.Bottom

    return xprop.ChangeProp32(xu, win, "_NET_FRAME_EXTENTS", "CARDINAL", raw...)
}

// _NET_MOVERESIZE_WINDOW req
// If 'w' or 'h' are 0, then they are not sent.
// If you need to resize a window without moving it, use the ReqExtra variant,
// or Resize.
func MoveresizeWindow(xu *xgbutil.XUtil, win xgb.Id, x, y int16,
                      w, h uint16) error {
    return MoveresizeWindowExtra(xu, win, x, y, w, h, xgb.GravityBitForget,
                                 2, true, true)
}

// _NET_MOVERESIZE_WINDOW req resize only
func ResizeWindow(xu *xgbutil.XUtil, win xgb.Id, w, h uint16) error {
    return MoveresizeWindowExtra(xu, win, 0, 0, w, h, xgb.GravityBitForget,
                                 2, false, false)
}

// _NET_MOVERESIZE_WINDOW req move only
func MoveWindow(xu *xgbutil.XUtil, win xgb.Id, x, y int16) error {
    return MoveresizeWindowExtra(xu, win, x, y, 0, 0, xgb.GravityBitForget,
                                 2, true, true)
}

// _NET_MOVERESIZE_WINDOW req extra
// If 'w' or 'h' are 0, then they are not sent.
// To not set 'x' or 'y', 'usex' or 'usey' need to be set to false.
func MoveresizeWindowExtra(xu *xgbutil.XUtil, win xgb.Id, x, y int16,
                           w, h uint16, gravity, source uint32,
                           usex, usey bool) error {
    flags := gravity
    flags |= source << 12
    if usex {
        flags |= 1 << 8
    }
    if usey {
        flags |= 1 << 9
    }
    if w > 0 {
        flags |= 1 << 10
    }
    if h > 0 {
        flags |= 1 << 11
    }

    return ClientEvent(xu, win, "_NET_MOVERESIZE_WINDOW", flags,
                       uint32(x), uint32(y), uint32(w), uint32(h))
}

// _NET_NUMBER_OF_DESKTOPS get
func NumberOfDesktopsGet(xu *xgbutil.XUtil) (uint32, error) {
    return xprop.PropValNum(xprop.GetProperty(xu, xu.RootWin(),
                                              "_NET_NUMBER_OF_DESKTOPS"))
}

// _NET_NUMBER_OF_DESKTOPS set
func NumberOfDesktopsSet(xu *xgbutil.XUtil, numDesks uint32) error {
    return xprop.ChangeProp32(xu, xu.RootWin(), "_NET_NUMBER_OF_DESKTOPS",
                              "CARDINAL", numDesks)
}

// _NET_NUMBER_OF_DESKTOPS req
func NumberOfDesktopsReq(xu *xgbutil.XUtil, numDesks uint32) error {
    return ClientEvent(xu, xu.RootWin(), "_NET_NUMBER_OF_DESKTOPS", numDesks)
}

// _NET_REQUEST_FRAME_EXTENTS req
func RequestFrameExtents(xu *xgbutil.XUtil, win xgb.Id) error {
    return ClientEvent(xu, win, "_NET_REQUEST_FRAME_EXTENTS")
}

// _NET_RESTACK_WINDOW req
// The shortcut here is to just raise the window to the top of the window stack.
func RestackWindow(xu *xgbutil.XUtil, win xgb.Id) error {
    return RestackWindowExtra(xu, win, xgb.StackModeAbove, 0, 2)
}

// _NET_RESTACK_WINDOW req extra
func RestackWindowExtra(xu *xgbutil.XUtil, win xgb.Id, stack_mode uint32,
                        sibling xgb.Id, source uint32) error {
    return ClientEvent(xu, win, "_NET_RESTACK_WINDOW", source, uint32(sibling),
                       stack_mode)
}

// _NET_SHOWING_DESKTOP get
func ShowingDesktopGet(xu *xgbutil.XUtil) (bool, error) {
    reply, err := xprop.GetProperty(xu, xu.RootWin(), "_NET_SHOWING_DESKTOP")
    if err != nil {
        return false, err
    }

    val, err := xprop.PropValNum(reply, nil)
    if err != nil {
        return false, err
    }

    return val == 1, nil
}

// _NET_SHOWING_DESKTOP set
func ShowingDesktopSet(xu *xgbutil.XUtil, show bool) error {
    var showInt uint32
    if show {
        showInt = 1
    } else {
        showInt = 0
    }
    return xprop.ChangeProp32(xu, xu.RootWin(), "_NET_SHOWING_DESKTOP",
                              "CARDINAL", showInt)
}

// _NET_SHOWING_DESKTOP req
func ShowingDesktopReq(xu *xgbutil.XUtil, show bool) error {
    var showInt uint32
    if show {
        showInt = 1
    } else {
        showInt = 0
    }
    return ClientEvent(xu, xu.RootWin(), "_NET_SHOWING_DESKTOP", showInt)
}

// _NET_SUPPORTED get
func SupportedGet(xu *xgbutil.XUtil) ([]string, error) {
    reply, err := xprop.GetProperty(xu, xu.RootWin(), "_NET_SUPPORTED")
    return xprop.PropValAtoms(xu, reply, err)
}

// _NET_SUPPORTED set
// This will create any atoms in the argument if they don't already exist.
func SupportedSet(xu *xgbutil.XUtil, atomNames []string) error {
    atoms, err := xprop.StrToAtoms(xu, atomNames)
    if err != nil {
        return err
    }

    return xprop.ChangeProp32(xu, xu.RootWin(), "_NET_SUPPORTED", "ATOM",
                              atoms...)
}

// _NET_SUPPORTING_WM_CHECK get
func SupportingWmCheckGet(xu *xgbutil.XUtil, win xgb.Id) (xgb.Id, error) {
    return xprop.PropValId(xprop.GetProperty(xu, win,
                                             "_NET_SUPPORTING_WM_CHECK"))
}

// _NET_SUPPORTING_WM_CHECK set
func SupportingWmCheckSet(xu *xgbutil.XUtil, win xgb.Id, wm_win xgb.Id) error {
    return xprop.ChangeProp32(xu, win, "_NET_SUPPORTING_WM_CHECK", "WINDOW",
                              uint32(wm_win))
}

// _NET_VIRTUAL_ROOTS get
func VirtualRootsGet(xu *xgbutil.XUtil) ([]xgb.Id, error) {
    return xprop.PropValIds(xprop.GetProperty(xu, xu.RootWin(),
                                              "_NET_VIRTUAL_ROOTS"))
}

// _NET_VIRTUAL_ROOTS set
func VirtualRootsSet(xu *xgbutil.XUtil, wins []xgb.Id) error {
    return xprop.ChangeProp32(xu, xu.RootWin(), "_NET_VIRTUAL_ROOTS", "WINDOW",
                              xprop.IdTo32(wins)...)
}

// _NET_VISIBLE_DESKTOPS get
// This is not parted of the EWMH spec, but is a property of my own creation.
// It allows the window manager to report that it has multiple desktops
// viewable at the same time. (This conflicts with other EWMH properties,
// so I don't think this will ever be added to the official spec.)
func VisibleDesktopsGet(xu *xgbutil.XUtil) ([]uint32, error) {
    return xprop.PropValNums(xprop.GetProperty(xu, xu.RootWin(),
                                               "_NET_VISIBLE_DESKTOPS"))
}

// _NET_VISIBLE_DESKTOPS set
func VisibleDesktopsSet(xu *xgbutil.XUtil, desktops []uint32) error {
    return xprop.ChangeProp32(xu, xu.RootWin(), "_NET_VISIBLE_DESKTOPS",
                              "CARDINAL", desktops...)
}

// _NET_WM_ALLOWED_ACTIONS get
func WmAllowedActionsGet(xu *xgbutil.XUtil, win xgb.Id) ([]string, error) {
    raw, err := xprop.GetProperty(xu, win, "_NET_WM_ALLOWED_ACTIONS")
    return xprop.PropValAtoms(xu, raw, err)
}

// _NET_WM_ALLOWED_ACTIONS set
func WmAllowedActionsSet(xu *xgbutil.XUtil, win xgb.Id,
                         atomNames []string) error {
    atoms, err := xprop.StrToAtoms(xu, atomNames)
    if err != nil {
        return err
    }

    return xprop.ChangeProp32(xu, win, "_NET_WM_ALLOWED_ACTIONS", "ATOM",
                              atoms...)
}

// _NET_WM_DESKTOP get
func WmDesktopGet(xu *xgbutil.XUtil, win xgb.Id) (uint32, error) {
    return xprop.PropValNum(xprop.GetProperty(xu, win, "_NET_WM_DESKTOP"))
}

// _NET_WM_DESKTOP set
func WmDesktopSet(xu *xgbutil.XUtil, win xgb.Id, desk uint32) error {
    return xprop.ChangeProp32(xu, win, "_NET_WM_DESKTOP", "CARDINAL", desk)
}

// _NET_WM_DESKTOP req
func WmDesktopReq(xu *xgbutil.XUtil, win xgb.Id, desk uint32) error {
    return WmDesktopReqExtra(xu, win, desk, 2)
}

// _NET_WM_DESKTOP req extra
func WmDesktopReqExtra(xu *xgbutil.XUtil, win xgb.Id, desk uint32,
                       source uint32) error {
    return ClientEvent(xu, win, "_NET_WM_DESKTOP", desk, source)
}

// WmFullscreenMonitors is a struct that organizes information related to the
// _NET_WM_FULLSCREEN_MONITORS property. Namely, the top, bottom, left and
// right monitor edges for a particular window.
type WmFullscreenMonitors struct {
    Top uint32
    Bottom uint32
    Left uint32
    Right uint32
}

// _NET_WM_FULLSCREEN_MONITORS get
func WmFullscreenMonitorsGet(xu *xgbutil.XUtil, win xgb.Id) (
     WmFullscreenMonitors, error) {
    raw, err := xprop.PropValNums(
                    xprop.GetProperty(xu, win, "_NET_WM_FULLSCREEN_MONITORS"))
    if err != nil {
        return WmFullscreenMonitors{}, err
    }

    return WmFullscreenMonitors{
        Top: raw[0],
        Bottom: raw[1],
        Left: raw[2],
        Right: raw[3],
    }, nil
}

// _NET_WM_FULLSCREEN_MONITORS set
func WmFullscreenMonitorsSet(xu *xgbutil.XUtil, win xgb.Id,
                             edges WmFullscreenMonitors) error {
    raw := make([]uint32, 4)
    raw[0] = edges.Top
    raw[1] = edges.Bottom
    raw[2] = edges.Left
    raw[3] = edges.Right

    return xprop.ChangeProp32(xu, win, "_NET_WM_FULLSCREEN_MONITORS",
                              "CARDINAL", raw...)
}

// _NET_WM_FULLSCREEN_MONITORS req
func WmFullscreenMonitorsReq(xu *xgbutil.XUtil, win xgb.Id,
                             edges WmFullscreenMonitors) error {
    return WmFullscreenMonitorsReqExtra(xu, win, edges, 2)
}

// _NET_WM_FULLSCREEN_MONITORS req extra
func WmFullscreenMonitorsReqExtra(xu *xgbutil.XUtil, win xgb.Id,
                                  edges WmFullscreenMonitors,
                                  source uint32) error {
    return ClientEvent(xu, win, "_NET_WM_FULLSCREEN_MONITORS",
                       edges.Top, edges.Bottom, edges.Left, edges.Right, source)
}

// _NET_WM_HANDLED_ICONS get
func WmHandledIconsGet(xu *xgbutil.XUtil, win xgb.Id) (bool, error) {
    reply, err := xprop.GetProperty(xu, win, "_NET_WM_HANDLED_ICONS")
    if err != nil {
        return false, err
    }

    val, err := xprop.PropValNum(reply, nil)
    if err != nil {
        return false, err
    }

    return val == 1, nil
}

// _NET_WM_HANDLED_ICONS set
func WmHandledIconsSet(xu *xgbutil.XUtil, handle bool) error {
    var handled uint32
    if handle {
        handled = 1
    } else {
        handled = 0
    }
    return xprop.ChangeProp32(xu, xu.RootWin(), "_NET_WM_HANDLED_ICONS",
                              "CARDINAL", handled)
}

// WmIcon is a struct that contains data for a single icon.
// The WmIcon method will return a list of these, since a single
// client can specify multiple icons of varying sizes.
type WmIcon struct {
    Width uint32
    Height uint32
    Data []uint32
}

// _NET_WM_ICON get
func WmIconGet(xu *xgbutil.XUtil, win xgb.Id) ([]*WmIcon, error) {
    icon, err := xprop.PropValNums(xprop.GetProperty(xu, win, "_NET_WM_ICON"))
    if err != nil {
        return nil, err
    }

    wmicons := make([]*WmIcon, 0)
    start := uint32(0)
    for int(start) < len(icon) {
        w, h := icon[start], icon[start + 1]
        upto := w * h

        wmicon := &WmIcon{
            Width: w,
            Height: h,
            Data: icon[(start + 2):(start + upto + 2)],
        }
        wmicons = append(wmicons, wmicon)

        start += upto + 2
    }

    return wmicons, nil
}

// _NET_WM_ICON set
func WmIconSet(xu *xgbutil.XUtil, win xgb.Id, icons []*WmIcon) error {
    raw := make([]uint32, 0, 10000) // start big
    for _, icon := range icons {
        raw = append(raw, icon.Width, icon.Height)
        raw = append(raw, icon.Data...)
    }

    return xprop.ChangeProp32(xu, win, "_NET_WM_ICON", "CARDINAL", raw...)
}

// WmIconGeometry struct organizes the information pertaining to the
// _NET_WM_ICON_GEOMETRY property. Namely, x, y, width and height.
type WmIconGeometry struct {
    X uint32
    Y uint32
    Width uint32
    Height uint32
}

// _NET_WM_ICON_GEOMETRY get
func WmIconGeometryGet(xu *xgbutil.XUtil, win xgb.Id) (WmIconGeometry, error) {
    geom, err := xprop.PropValNums(xprop.GetProperty(xu, win,
                                                     "_NET_WM_ICON_GEOMETRY"))
    if err != nil {
        return WmIconGeometry{}, err
    }

    return WmIconGeometry{
        X: geom[0],
        Y: geom[1],
        Width: geom[2],
        Height: geom[3],
    }, nil
}

// _NET_WM_ICON_GEOMETRY set
func WmIconGeometrySet(xu *xgbutil.XUtil, win xgb.Id, geom WmIconGeometry) error {
    rawGeom := make([]uint32, 4)
    rawGeom[0] = geom.X
    rawGeom[1] = geom.Y
    rawGeom[2] = geom.Width
    rawGeom[3] = geom.Height

    return xprop.ChangeProp32(xu, win, "_NET_WM_ICON_GEOMETRY", "CARDINAL",
                              rawGeom...)
}

// _NET_WM_ICON_NAME get
func WmIconNameGet(xu *xgbutil.XUtil, win xgb.Id) (string, error) {
    return xprop.PropValStr(xprop.GetProperty(xu, win, "_NET_WM_ICON_NAME"))
}

// _NET_WM_ICON_NAME set
func WmIconNameSet(xu *xgbutil.XUtil, win xgb.Id, name string) error {
    return xprop.ChangeProp(xu, win, 8, "_NET_WM_ICON_NAME", "UTF8_STRING",
                            []byte(name))
}

// _NET_WM_MOVERESIZE constants
const (
    SizeTopLeft = iota
    SizeTop
    SizeTopRight
    SizeRight
    SizeBottomRight
    SizeBottom
    SizeBottomLeft
    SizeLeft
    Move
    SizeKeyboard
    MoveKeyboard
    Cancel
    Infer // special for Wingo. DO NOT USE.
)

// _NET_WM_MOVERESIZE req
func WmMoveresize(xu *xgbutil.XUtil, win xgb.Id, direction uint32) error {
    return WmMoveresizeExtra(xu, win, direction, 0, 0, 0, 2)
}

// _NET_WM_MOVERESIZE req extra
func WmMoveresizeExtra(xu *xgbutil.XUtil, win xgb.Id, direction, x_root, y_root,
                       button, source uint32) error {
    return ClientEvent(xu, win, "_NET_WM_MOVERESIZE", x_root, y_root,
                       direction, button, source)
}

// _NET_WM_NAME get
func WmNameGet(xu *xgbutil.XUtil, win xgb.Id) (string, error) {
    return xprop.PropValStr(xprop.GetProperty(xu, win, "_NET_WM_NAME"))
}

// _NET_WM_NAME set
func WmNameSet(xu *xgbutil.XUtil, win xgb.Id, name string) error {
    return xprop.ChangeProp(xu, win, 8, "_NET_WM_NAME", "UTF8_STRING",
                            []byte(name))
}

// WmOpaqueRegion organizes information related to the _NET_WM_OPAQUE_REGION
// property. Namely, the x, y, width and height of an opaque rectangle
// relative to the client window.
type WmOpaqueRegion struct {
    X uint32
    Y uint32
    Width uint32
    Height uint32
}

// _NET_WM_OPAQUE_REGION get
func WmOpaqueRegionGet(xu *xgbutil.XUtil, win xgb.Id) (
     []WmOpaqueRegion, error) {
    raw, err := xprop.PropValNums(xprop.GetProperty(xu, win,
                                                    "_NET_WM_OPAQUE_REGION"))
    if err != nil {
        return nil, err
    }

    regions := make([]WmOpaqueRegion, len(raw) / 4)
    for i, _ := range(regions) {
        regions[i] = WmOpaqueRegion{
            X: raw[i * 4 + 0],
            Y: raw[i * 4 + 1],
            Width: raw[i * 4 + 2],
            Height: raw[i * 4 + 3],
        }
    }
    return regions, nil
}

// _NET_WM_OPAQUE_REGION set
func WmOpaqueRegionSet(xu *xgbutil.XUtil, win xgb.Id,
                       regions []WmOpaqueRegion) error {
    raw := make([]uint32, len(regions) * 4)

    for i, region := range(regions) {
        raw[i * 4 + 0] = region.X
        raw[i * 4 + 1] = region.Y
        raw[i * 4 + 2] = region.Width
        raw[i * 4 + 3] = region.Height
    }

    return xprop.ChangeProp32(xu, win, "_NET_WM_OPAQUE_REGION", "CARDINAL",
                              raw...)
}

// _NET_WM_PID get
func WmPidGet(xu *xgbutil.XUtil, win xgb.Id) (uint32, error) {
    return xprop.PropValNum(xprop.GetProperty(xu, win, "_NET_WM_PID"))
}

// _NET_WM_PID set
func WmPidSet(xu *xgbutil.XUtil, win xgb.Id, pid uint32) error {
    return xprop.ChangeProp32(xu, win, "_NET_WM_PID", "CARDINAL", pid)
}

// _NET_WM_PING req
func WmPing(xu *xgbutil.XUtil, win xgb.Id, response bool) error {
    return WmPingExtra(xu, win, response, 0)
}

// _NET_WM_PING req extra
func WmPingExtra(xu *xgbutil.XUtil, win xgb.Id, response bool,
                                 time xgb.Timestamp) error {
    pingAtom, err := xprop.Atm(xu, "_NET_WM_PING")
    if err != nil {
        return err
    }

    var evWindow xgb.Id
    if response {
        evWindow = xu.RootWin()
    } else {
        evWindow = win
    }

    return ClientEvent(xu, evWindow, "WM_PROTOCOLS", uint32(pingAtom), time,
                       win)
}

// _NET_WM_STATE constants for state toggling
// These correspond to the "action" parameter.
const (
    StateRemove = iota
    StateAdd
    StateToggle
)

// _NET_WM_STATE get
func WmStateGet(xu *xgbutil.XUtil, win xgb.Id) ([]string, error) {
    raw, err := xprop.GetProperty(xu, win, "_NET_WM_STATE")
    return xprop.PropValAtoms(xu, raw, err)
}

// _NET_WM_STATE set
func WmStateSet(xu *xgbutil.XUtil, win xgb.Id, atomNames []string) error {
    atoms, err := xprop.StrToAtoms(xu, atomNames)
    if err != nil {
        return err
    }

    return xprop.ChangeProp32(xu, win, "_NET_WM_STATE", "ATOM", atoms...)
}

// _NET_WM_STATE req
func WmStateReq(xu *xgbutil.XUtil, win xgb.Id, action uint32,
                                atomName string) error {
    return WmStateReqExtra(xu, win, action, atomName, "", 2)
}

// _NET_WM_STATE req extra
func WmStateReqExtra(xu *xgbutil.XUtil, win xgb.Id, action uint32,
                     first string, second string, source uint32) (err error) {
    var atom1, atom2 xgb.Id

    atom1, err = xprop.Atom(xu, first, false)
    if err != nil {
        return err
    }

    if len(second) > 0 {
        atom2, err = xprop.Atom(xu, second, false)
        if err != nil {
            return err
        }
    } else {
        atom2 = 0
    }

    return ClientEvent(xu, win, "_NET_WM_STATE", action, uint32(atom1),
                       uint32(atom2), source)
}

// WmStrut struct organizes information for the _NET_WM_STRUT property.
// Namely, it encapsulates its four values: left, right, top and bottom.
type WmStrut struct {
    Left uint32
    Right uint32
    Top uint32
    Bottom uint32
}

// _NET_WM_STRUT get
func WmStrutGet(xu *xgbutil.XUtil, win xgb.Id) (WmStrut, error) {
    struts, err := xprop.PropValNums(xprop.GetProperty(xu, win,
                                                       "_NET_WM_STRUT"))
    if err != nil {
        return WmStrut{}, err
    }

    return WmStrut {
        Left: struts[0],
        Right: struts[1],
        Top: struts[2],
        Bottom: struts[3],
    }, nil
}

// _NET_WM_STRUT set
func WmStrutSet(xu *xgbutil.XUtil, win xgb.Id, struts WmStrut) error {
    rawStruts := make([]uint32, 4)
    rawStruts[0] = struts.Left
    rawStruts[1] = struts.Right
    rawStruts[2] = struts.Top
    rawStruts[3] = struts.Bottom

    return xprop.ChangeProp32(xu, win, "_NET_WM_STRUT", "CARDINAL",
                              rawStruts...)
}

// WmStrutPartial struct organizes information for the _NET_WM_STRUT_PARTIAL
// property. Namely, it encapsulates its twelve values: left, right, top,
// bottom, left_start_y, left_end_y, right_start_y, right_end_y,
// top_start_x, top_end_x, bottom_start_x, and bottom_end_x.
type WmStrutPartial struct {
    Left, Right, Top, Bottom uint32
    LeftStartY, LeftEndY, RightStartY, RightEndY uint32
    TopStartX, TopEndX, BottomStartX, BottomEndX uint32
}

// _NET_WM_STRUT_PARTIAL get
func WmStrutPartialGet(xu *xgbutil.XUtil, win xgb.Id) (WmStrutPartial, error) {
    struts, err := xprop.PropValNums(xprop.GetProperty(xu, win,
                                                       "_NET_WM_STRUT_PARTIAL"))
    if err != nil {
        return WmStrutPartial{}, err
    }

    return WmStrutPartial {
        Left: struts[0], Right: struts[1], Top: struts[2], Bottom: struts[3],
        LeftStartY: struts[4], LeftEndY: struts[5],
        RightStartY: struts[6], RightEndY: struts[7],
        TopStartX: struts[8], TopEndX: struts[9],
        BottomStartX: struts[10], BottomEndX: struts[11],
    }, nil
}

// _NET_WM_STRUT_PARTIAL set
func WmStrutPartialSet(xu *xgbutil.XUtil, win xgb.Id,
                       struts WmStrutPartial) error {
    rawStruts := make([]uint32, 4)
    rawStruts[0] = struts.Left
    rawStruts[1] = struts.Right
    rawStruts[2] = struts.Top
    rawStruts[3] = struts.Bottom
    rawStruts[4] = struts.LeftStartY
    rawStruts[5] = struts.LeftEndY
    rawStruts[6] = struts.RightStartY
    rawStruts[7] = struts.RightEndY
    rawStruts[8] = struts.TopStartX
    rawStruts[9] = struts.TopEndX
    rawStruts[10] = struts.BottomStartX
    rawStruts[11] = struts.BottomEndX

    return xprop.ChangeProp32(xu, win, "_NET_WM_STRUT_PARTIAL", "CARDINAL",
                              rawStruts...)
}

// _NET_WM_SYNC_REQUEST req
func WmSyncRequest(xu *xgbutil.XUtil, win xgb.Id, req_num uint64) error {
    return WmSyncRequestExtra(xu, win, req_num, 0)
}

// _NET_WM_SYNC_REQUEST req extra
func WmSyncRequestExtra(xu *xgbutil.XUtil, win xgb.Id, req_num uint64,
                                        time xgb.Timestamp) error {
    syncReq, err := xprop.Atm(xu, "_NET_WM_SYNC_REQUEST")
    if err != nil {
        return err
    }

    high := uint32(req_num >> 32)
    low := uint32(req_num << 32 ^ req_num)

    return ClientEvent(xu, win, "WM_PROTOCOLS", syncReq, time, low, high)
}

// _NET_WM_SYNC_REQUEST_COUNTER get 
// I'm pretty sure this needs 64 bit integers, but I'm not quite sure
// how to go about that yet. Any ideas?
func WmSyncRequestCounter(xu *xgbutil.XUtil, win xgb.Id) (uint32, error) {
    return xprop.PropValNum(xprop.GetProperty(xu, win,
                                              "_NET_WM_SYNC_REQUEST_COUNTER"))
}

// _NET_WM_SYNC_REQUEST_COUNTER set
// I'm pretty sure this needs 64 bit integers, but I'm not quite sure
// how to go about that yet. Any ideas?
func WmSyncRequestCounterSet(xu *xgbutil.XUtil, win xgb.Id,
                             counter uint32) error {
    return xprop.ChangeProp32(xu, win, "_NET_WM_SYNC_REQUEST_COUNTER",
                              "CARDINAL", counter)
}

// _NET_WM_USER_TIME get
func WmUserTimeGet(xu *xgbutil.XUtil, win xgb.Id) (uint32, error) {
    return xprop.PropValNum(xprop.GetProperty(xu, win, "_NET_WM_USER_TIME"))
}

// _NET_WM_USER_TIME set
func WmUserTimeSet(xu *xgbutil.XUtil, win xgb.Id, user_time uint32) error {
    return xprop.ChangeProp32(xu, win, "_NET_WM_USER_TIME", "CARDINAL",
                              user_time)
}

// _NET_WM_USER_TIME_WINDOW get
func WmUserTimeWindowGet(xu *xgbutil.XUtil, win xgb.Id) (xgb.Id, error) {
    return xprop.PropValId(xprop.GetProperty(xu, win,
                                             "_NET_WM_USER_TIME_WINDOW"))
}

// _NET_WM_USER_TIME set
func WmUserTimeWindowSet(xu *xgbutil.XUtil, win xgb.Id, time_win xgb.Id) error {
    return xprop.ChangeProp32(xu, win, "_NET_WM_USER_TIME_WINDOW", "CARDINAL",
                              uint32(time_win))
}

// _NET_WM_VISIBLE_ICON_NAME get
func WmVisibleIconNameGet(xu *xgbutil.XUtil, win xgb.Id) (string, error) {
    return xprop.PropValStr(xprop.GetProperty(xu, win,
                                              "_NET_WM_VISIBLE_ICON_NAME"))
}

// _NET_WM_VISIBLE_ICON_NAME set
func WmVisibleIconNameSet(xu *xgbutil.XUtil, win xgb.Id, name string) error {
    return xprop.ChangeProp(xu, win, 8, "_NET_WM_VISIBLE_ICON_NAME",
                            "UTF8_STRING", []byte(name))
}

// _NET_WM_VISIBLE_NAME get
func WmVisibleNameGet(xu *xgbutil.XUtil, win xgb.Id) (string, error) {
    return xprop.PropValStr(xprop.GetProperty(xu, win, "_NET_WM_VISIBLE_NAME"))
}

// _NET_WM_VISIBLE_NAME set
func WmVisibleNameSet(xu *xgbutil.XUtil, win xgb.Id, name string) error {
    return xprop.ChangeProp(xu, win, 8, "_NET_WM_VISIBLE_NAME", "UTF8_STRING",
                            []byte(name))
}

// _NET_WM_WINDOW_OPACITY get
// This isn't part of the EWMH spec, but is widely used by drop in
// compositing managers (i.e., xcompmgr, cairo-compmgr, etc.).
// This property is typically set not on a client window, but the *parent*
// of a client window in reparenting window managers.
func WmWindowOpacityGet(xu *xgbutil.XUtil, win xgb.Id) (float64, error) {
    intOpacity, err := xprop.PropValNum(
                           xprop.GetProperty(xu, win, "_NET_WM_WINDOW_OPACITY"))
    if err != nil {
        return 0, err
    }

    return float64(intOpacity) / float64(0xffffffff), nil
}

// _NET_WM_WINDOW_OPACITY set
func WmWindowOpacitySet(xu *xgbutil.XUtil, win xgb.Id, opacity float64) error {
    return xprop.ChangeProp32(xu, win, "_NET_WM_WINDOW_OPACITY", "CARDINAL",
                              uint32(opacity * 0xffffffff))
}

// _NET_WM_WINDOW_TYPE get
func WmWindowTypeGet(xu *xgbutil.XUtil, win xgb.Id) ([]string, error) {
    raw, err := xprop.GetProperty(xu, win, "_NET_WM_WINDOW_TYPE")
    return xprop.PropValAtoms(xu, raw, err)
}

// _NET_WM_WINDOW_TYPE set
// This will create any atoms used in 'atomNames' if they don't already exist.
func WmWindowTypeSet(xu *xgbutil.XUtil, win xgb.Id, atomNames []string) error {
    atoms, err := xprop.StrToAtoms(xu, atomNames)
    if err != nil {
        return err
    }

    return xprop.ChangeProp32(xu, win, "_NET_WM_WINDOW_TYPE", "ATOM", atoms...)
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
func WorkareaGet(xu *xgbutil.XUtil) ([]Workarea, error) {
    rects, err := xprop.PropValNums(xprop.GetProperty(xu, xu.RootWin(),
                                                      "_NET_WORKAREA"))
    if err != nil {
        return nil, err
    }

    workareas := make([]Workarea, len(rects) / 4)
    for i, _ := range workareas {
        workareas[i] = Workarea {
            X: rects[i * 4],
            Y: rects[i * 4 + 1],
            Width: rects[i * 4 + 2],
            Height: rects[i * 4 + 3],
        }
    }
    return workareas, nil
}

// _NET_WORKAREA set
func WorkareaSet(xu *xgbutil.XUtil, workareas []Workarea) error {
    rects := make([]uint32, len(workareas) * 4)
    for i, workarea := range workareas {
        rects[i * 4] = workarea.X
        rects[i * 4 + 1] = workarea.Y
        rects[i * 4 + 2] = workarea.Width
        rects[i * 4 + 3] = workarea.Height
    }

    return xprop.ChangeProp32(xu, xu.RootWin(), "_NET_WORKAREA", "CARDINAL",
                              rects...)
}

