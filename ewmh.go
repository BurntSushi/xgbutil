/*
    An API to the entire EWMH spec.

    Since there are so many methods and they adhere to an existing spec, 
    this source file does not contain much documentation. Indeed, each
    method has only a single comment associated with it: the EWMH property name.

    See the EWMH spec for more info:
    http://standards.freedesktop.org/wm-spec/wm-spec-latest.html

    Here is the naming scheme using "_NET_ACTIVE_WINDOW" as an example.

    Methods "EwmhActiveWindow" and "EwmhActiveWindowSet" get and set the
    property, respectively. Both of these methods exist for most EWMH 
    properties.  Additionally, some EWMH properties support sending a client 
    message event to request the window manager perform some action. In the 
    case of "_NET_ACTIVE_WINDOW", this request is used to set the active 
    window.

    These sorts of methods end in "Req". So for "_NET_ACTIVE_WINDOW",
    the method name is "EwmhActiveWindowReq". Moreover, most requests include
    various parameters that don't need to be changed often (like the source
    indication). Thus, by default, methods ending in "Req" force these to
    sensible defaults. If you need access to all of the parameters, use the
    corresponding "ReqExtra" method. So for "_NET_ACTIVE_WINDOW", that would
    be "EwmhActiveWindowReqExtra". (If no "ReqExtra" method exists, then the
    "Req" method covers all available parameters.)

    This naming scheme has one exception: if a property's only use is through
    sending an event (like "_NET_CLOSE_WINDOW"), then the name will be
    "EwmhCloseWindow" for the short-hand version and "EwmhCloseWindowExtra"
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
package xgbutil

import "code.google.com/p/x-go-binding/xgb"

// EwmhClientEvent is a convenience function that sends ClientMessage events
// to the root window as specified by the EWMH spec.
func (xu *XUtil) EwmhClientEvent(window xgb.Id, message_type string,
                                 data ...interface{}) error {
    mstype, err := xu.Atm(message_type)
    if err != nil {
        return err
    }

    evMask := (xgb.EventMaskSubstructureNotify |
               xgb.EventMaskSubstructureRedirect)
    cm, err := NewClientMessage(32, window, mstype, data...)
    if err != nil {
        return err
    }

    xu.SendRootEvent(cm, uint32(evMask))
    return nil
}

// _NET_ACTIVE_WINDOW get
func (xu *XUtil) EwmhActiveWindow() (xgb.Id, error) {
    return PropValId(xu.GetProperty(xu.root, "_NET_ACTIVE_WINDOW"))
}

// _NET_ACTIVE_WINDOW set
func (xu *XUtil) EwmhActiveWindowSet(win xgb.Id) error {
    return xu.ChangeProperty32(xu.root, "_NET_ACTIVE_WINDOW", "WINDOW",
                               uint32(win))
}

// _NET_ACTIVE_WINDOW req
func (xu *XUtil) EwmhActiveWindowReq(win xgb.Id) error {
    return xu.EwmhActiveWindowReqExtra(win, 2, 0, 0)
}

// _NET_ACTIVE_WINDOW req extra
func (xu *XUtil) EwmhActiveWindowReqExtra(win xgb.Id, source uint32,
                                          time xgb.Timestamp,
                                          current_active xgb.Id) error {
    return xu.EwmhClientEvent(win, "_NET_ACTIVE_WINDOW", source, uint32(time),
                              uint32(current_active))
}

// _NET_CLIENT_LIST get
func (xu *XUtil) EwmhClientList() ([]xgb.Id, error) {
    return PropValIds(xu.GetProperty(xu.root, "_NET_CLIENT_LIST"))
}

// _NET_CLIENT_LIST set
func (xu *XUtil) EwmhClientListSet(wins []xgb.Id) error {
    return xu.ChangeProperty32(xu.root, "_NET_CLIENT_LIST", "WINDOW",
                               IdTo32(wins)...)
}

// _NET_CLIENT_LIST_STACKING get
func (xu *XUtil) EwmhClientListStacking() ([]xgb.Id, error) {
    return PropValIds(xu.GetProperty(xu.root, "_NET_CLIENT_LIST_STACKING"))
}

// _NET_CLIENT_LIST_STACKING set
func (xu *XUtil) EwmhClientListStackingSet(wins []xgb.Id) error {
    return xu.ChangeProperty32(xu.root, "_NET_CLIENT_LIST_STACKING", "WINDOW",
                               IdTo32(wins)...)
}

// _NET_CLOSE_WINDOW req
func (xu *XUtil) EwmhCloseWindow(win xgb.Id) error {
    return xu.EwmhCloseWindowExtra(win, 0, 2)
}

// _NET_CLOSE_WINDOW req extra
func (xu *XUtil) EwmhCloseWindowExtra(win xgb.Id, time xgb.Timestamp,
                                         source uint32) error {
    return xu.EwmhClientEvent(win, "_NET_CLOSE_WINDOW", uint32(time), source)
}

// _NET_CURRENT_DESKTOP get
func (xu *XUtil) EwmhCurrentDesktop() (uint32, error) {
    return PropValNum(xu.GetProperty(xu.root, "_NET_CURRENT_DESKTOP"))
}

// _NET_CURRENT_DESKTOP set
func (xu *XUtil) EwmhCurrentDesktopSet(desk uint32) error {
    return xu.ChangeProperty32(xu.root, "_NET_CURRENT_DESKTOP", "CARDINAL",
                               desk)
}

// _NET_CURRENT_DESKTOP req
func (xu *XUtil) EwmhCurrentDesktopReq(desk uint32) error {
    return xu.EwmhCurrentDesktopReqExtra(desk, 0)
}

// _NET_CURRENT_DESKTOP req extra
func (xu *XUtil) EwmhCurrentDesktopReqExtra(desk uint32,
                                            time xgb.Timestamp) error {
    return xu.EwmhClientEvent(xu.root, "_NET_CURRENT_DESKTOP", desk, time)
}

// _NET_DESKTOP_NAMES get
func (xu *XUtil) EwmhDesktopNames() ([]string, error) {
    return PropValStrs(xu.GetProperty(xu.root, "_NET_DESKTOP_NAMES"))
}

// _NET_DESKTOP_NAMES set
func (xu *XUtil) EwmhDesktopNamesSet(names []string) error {
    nullterm := make([]byte, 0)
    for _, name := range names {
        nullterm = append(nullterm, name...)
        nullterm = append(nullterm, 0)
    }
    return xu.ChangeProperty(xu.root, 8, "_NET_DESKTOP_NAMES", "UTF8_STRING",
                             nullterm)
}

// DesktopGeometry is a struct that houses the width and height of a
// _NET_DESKTOP_GEOMETRY property reply.
type DesktopGeometry struct {
    Width uint32
    Height uint32
}

// _NET_DESKTOP_GEOMETRY get
func (xu *XUtil) EwmhDesktopGeometry() (DesktopGeometry, error) {
    geom, err := PropValNums(xu.GetProperty(xu.root, "_NET_DESKTOP_GEOMETRY"))
    if err != nil {
        return DesktopGeometry{}, err
    }

    return DesktopGeometry{Width: geom[0], Height: geom[1]}, nil
}

// _NET_DESKTOP_GEOMETRY set
func (xu *XUtil) EwmhDesktopGeometrySet(dg DesktopGeometry) error {
    return xu.ChangeProperty32(xu.root, "_NET_DESKTOP_GEOMETRY", "CARDINAL",
                               dg.Width, dg.Height)
}

// _NET_DESKTOP_GEOMETRY req
func (xu *XUtil) EwmhDesktopGeometryReq(dg DesktopGeometry) error {
    return xu.EwmhClientEvent(xu.root, "_NET_DESKTOP_GEOMETRY", dg.Width,
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
    EwmhOrientHorz = iota
    EwmhOrientVert
)

// _NET_DESKTOP_LAYOUT constants for starting corner
const (
    EwmhTopLeft = iota
    EwmhTopRight
    EwmhBottomRight
    EwmhBottomLeft
)

// _NET_DESKTOP_LAYOUT get
func (xu *XUtil) EwmhDesktopLayout() (dl DesktopLayout, err error) {
    dlraw, err := PropValNums(xu.GetProperty(xu.root, "_NET_DESKTOP_LAYOUT"))
    if err != nil {
        return DesktopLayout{}, err
    }

    dl.Orientation = dlraw[0]
    dl.Columns = dlraw[1]
    dl.Rows = dlraw[2]

    if len(dlraw) > 3 {
        dl.StartingCorner = dlraw[3]
    } else {
        dl.StartingCorner = EwmhTopLeft
    }

    return dl, nil
}

// _NET_DESKTOP_LAYOUT set
func (xu *XUtil) EwmhDesktopLayoutSet(orientation, columns, rows,
                                      startingCorner uint32) error {
    return xu.ChangeProperty32(xu.root, "_NET_DESKTOP_LAYOUT", "CARDINAL",
                               orientation, columns, rows, startingCorner)
}

// DesktopViewport is a struct that contains a pairing of x,y coordinates
// representing the top-left corner of each desktop. (There will typically
// be one struct here for each desktop in existence.)
type DesktopViewport struct {
    X uint32
    Y uint32
}

// _NET_DESKTOP_VIEWPORT get
func (xu *XUtil) EwmhDesktopViewport() ([]DesktopViewport, error) {
    coords, err := PropValNums(xu.GetProperty(xu.root, "_NET_DESKTOP_VIEWPORT"))
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
func (xu *XUtil) EwmhDesktopViewportSet(viewports []DesktopViewport) error {
    coords := make([]uint32, len(viewports) * 2)
    for i, viewport := range viewports {
        coords[i * 2] = viewport.X
        coords[i * 2 + 1] = viewport.Y
    }

    return xu.ChangeProperty32(xu.root, "_NET_DESKTOP_VIEWPORT", "CARDINAL",
                               coords...)
}

// _NET_DESKTOP_VIEWPORT req
func (xu *XUtil) EwmhDesktopViewportReq(x uint32, y uint32) error {
    return xu.EwmhClientEvent(xu.root, "_NET_DESKTOP_VIEWPORT", x, y)
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
func (xu *XUtil) EwmhFrameExtents(win xgb.Id) (FrameExtents, error) {
    raw, err := PropValNums(xu.GetProperty(win, "_NET_FRAME_EXTENTS"))
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
func (xu *XUtil) EwmhFrameExtentsSet(win xgb.Id, extents FrameExtents) error {
    raw := make([]uint32, 4)
    raw[0] = extents.Left
    raw[1] = extents.Right
    raw[2] = extents.Top
    raw[3] = extents.Bottom

    return xu.ChangeProperty32(win, "_NET_FRAME_EXTENTS", "CARDINAL", raw...)
}

// _NET_MOVERESIZE_WINDOW req
// If 'w' or 'h' are 0, then they are not sent.
// If you need to resize a window without moving it, use the ReqExtra variant,
// or EwmhResize.
func (xu *XUtil) EwmhMoveresizeWindow(win xgb.Id, x, y, w, h uint32) error {
    return xu.EwmhMoveresizeWindowExtra(win, x, y, w, h, xgb.GravityBitForget,
                                        2, true, true)
}

// _NET_MOVERESIZE_WINDOW req resize only
func (xu *XUtil) EwmhResizeWindow(win xgb.Id, w, h uint32) error {
    return xu.EwmhMoveresizeWindowExtra(win, 0, 0, w, h, xgb.GravityBitForget,
                                        2, false, false)
}

// _NET_MOVERESIZE_WINDOW req move only
func (xu *XUtil) EwmhMoveWindow(win xgb.Id, x, y uint32) error {
    return xu.EwmhMoveresizeWindowExtra(win, x, y, 0, 0, xgb.GravityBitForget,
                                        2, true, true)
}

// _NET_MOVERESIZE_WINDOW req extra
// If 'w' or 'h' are 0, then they are not sent.
// To not set 'x' or 'y', 'usex' or 'usey' need to be set to false.
func (xu *XUtil) EwmhMoveresizeWindowExtra(win xgb.Id, x, y, w, h,
                                           gravity, source uint32,
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

    return xu.EwmhClientEvent(win, "_NET_MOVERESIZE_WINDOW", flags, x, y, w, h)
}

// _NET_NUMBER_OF_DESKTOPS get
func (xu *XUtil) EwmhNumberOfDesktops() (uint32, error) {
    return PropValNum(xu.GetProperty(xu.root, "_NET_NUMBER_OF_DESKTOPS"))
}

// _NET_NUMBER_OF_DESKTOPS set
func (xu *XUtil) EwmhNumberOfDesktopsSet(numDesks uint32) error {
    return xu.ChangeProperty32(xu.root, "_NET_NUMBER_OF_DESKTOPS", "CARDINAL",
                               numDesks)
}

// _NET_NUMBER_OF_DESKTOPS req
func (xu *XUtil) EwmhNumberOfDesktopsReq(numDesks uint32) error {
    return xu.EwmhClientEvent(xu.root, "_NET_NUMBER_OF_DESKTOPS", numDesks)
}

// _NET_REQUEST_FRAME_EXTENTS req
func (xu *XUtil) EwmhRequestFrameExtents(win xgb.Id) error {
    return xu.EwmhClientEvent(win, "_NET_REQUEST_FRAME_EXTENTS")
}

// _NET_RESTACK_WINDOW req
// The shortcut here is to just raise the window to the top of the window stack.
func (xu *XUtil) EwmhRestackWindow(win xgb.Id) error {
    return xu.EwmhRestackWindowExtra(win, xgb.StackModeAbove, 0, 2)
}

// _NET_RESTACK_WINDOW req extra
func (xu *XUtil) EwmhRestackWindowExtra(win xgb.Id, stack_mode uint32,
                                        sibling xgb.Id, source uint32) error {
    return xu.EwmhClientEvent(win, "_NET_RESTACK_WINDOW", source,
                              uint32(sibling), stack_mode)
}

// _NET_SHOWING_DESKTOP get
func (xu *XUtil) EwmhShowingDesktop() (bool, error) {
    reply, err := xu.GetProperty(xu.root, "_NET_SHOWING_DESKTOP")
    if err != nil {
        return false, err
    }

    val, err := PropValNum(reply, nil)
    if err != nil {
        return false, err
    }

    return val == 1, nil
}

// _NET_SHOWING_DESKTOP set
func (xu *XUtil) EwmhShowingDesktopSet(show bool) error {
    var showInt uint32
    if show {
        showInt = 1
    } else {
        showInt = 0
    }
    return xu.ChangeProperty32(xu.root, "_NET_SHOWING_DESKTOP", "CARDINAL",
                               showInt)
}

// _NET_SHOWING_DESKTOP req
func (xu *XUtil) EwmhShowingDesktopReq(show bool) error {
    var showInt uint32
    if show {
        showInt = 1
    } else {
        showInt = 0
    }
    return xu.EwmhClientEvent(xu.root, "_NET_SHOWING_DESKTOP", showInt)
}

// _NET_SUPPORTED get
func (xu *XUtil) EwmhSupported() ([]string, error) {
    return xu.PropValAtoms(xu.GetProperty(xu.root, "_NET_SUPPORTED"))
}

// _NET_SUPPORTED set
// This will create any atoms in the argument if they don't already exist.
func (xu *XUtil) EwmhSupportedSet(atomNames []string) error {
    atoms, err := xu.StrToAtoms(atomNames)
    if err != nil {
        return err
    }

    return xu.ChangeProperty32(xu.root, "_NET_SUPPORTED", "ATOM", atoms...)
}

// _NET_SUPPORTING_WM_CHECK get
func (xu *XUtil) EwmhSupportingWmCheck(win xgb.Id) (xgb.Id, error) {
    return PropValId(xu.GetProperty(win, "_NET_SUPPORTING_WM_CHECK"))
}

// _NET_SUPPORTING_WM_CHECK set
func (xu *XUtil) EwmhSupportingWmCheckSet(win xgb.Id, wm_win xgb.Id) error {
    return xu.ChangeProperty32(win, "_NET_SUPPORTING_WM_CHECK", "WINDOW",
                               uint32(wm_win))
}

// _NET_VIRTUAL_ROOTS get
func (xu *XUtil) EwmhVirtualRoots() ([]xgb.Id, error) {
    return PropValIds(xu.GetProperty(xu.root, "_NET_VIRTUAL_ROOTS"))
}

// _NET_VIRTUAL_ROOTS set
func (xu *XUtil) EwmhVirtualRootsSet(wins []xgb.Id) error {
    return xu.ChangeProperty32(xu.root, "_NET_VIRTUAL_ROOTS", "WINDOW",
                               IdTo32(wins)...)
}

// _NET_VISIBLE_DESKTOPS get
// This is not parted of the EWMH spec, but is a property of my own creation.
// It allows the window manager to report that it has multiple desktops
// viewable at the same time. (This conflicts with other EWMH properties,
// so I don't think this will ever be added to the official spec.)
func (xu *XUtil) EwmhVisibleDesktops() ([]uint32, error) {
    return PropValNums(xu.GetProperty(xu.root, "_NET_VISIBLE_DESKTOPS"))
}

// _NET_VISIBLE_DESKTOPS set
func (xu *XUtil) EwmhVisibleDesktopsSet(desktops []uint32) error {
    return xu.ChangeProperty32(xu.root, "_NET_VISIBLE_DESKTOPS", "CARDINAL",
                               desktops...)
}

// _NET_WM_ALLOWED_ACTIONS get
func (xu *XUtil) EwmhWmAllowedActions(win xgb.Id) ([]string, error) {
    return xu.PropValAtoms(xu.GetProperty(win, "_NET_WM_ALLOWED_ACTIONS"))
}

// _NET_WM_ALLOWED_ACTIONS set
func (xu *XUtil) EwmhWmAllowedActionsSet(win xgb.Id, atomNames []string) error {
    atoms, err := xu.StrToAtoms(atomNames)
    if err != nil {
        return err
    }

    return xu.ChangeProperty32(win, "_NET_WM_ALLOWED_ACTIONS", "ATOM", atoms...)
}

// _NET_WM_DESKTOP get
func (xu *XUtil) EwmhWmDesktop(win xgb.Id) (uint32, error) {
    return PropValNum(xu.GetProperty(win, "_NET_WM_DESKTOP"))
}

// _NET_WM_DESKTOP set
func (xu *XUtil) EwmhWmDesktopSet(win xgb.Id, desk uint32) error {
    return xu.ChangeProperty32(win, "_NET_WM_DESKTOP", "CARDINAL", desk)
}

// _NET_WM_DESKTOP req
func (xu *XUtil) EwmhWmDesktopReq(win xgb.Id, desk uint32) error {
    return xu.EwmhWmDesktopReqExtra(win, desk, 2)
}

// _NET_WM_DESKTOP req extra
func (xu *XUtil) EwmhWmDesktopReqExtra(win xgb.Id, desk uint32,
                                       source uint32) error {
    return xu.EwmhClientEvent(win, "_NET_WM_DESKTOP", desk, source)
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
func (xu *XUtil) EwmhWmFullscreenMonitors(win xgb.Id) (
                 WmFullscreenMonitors, error) {
    raw, err := PropValNums(xu.GetProperty(win, "_NET_WM_FULLSCREEN_MONITORS"))
    if err != nil {
        return WmFullscreenMonitors{}, err
    }

    return WmFullscreenMonitors{
        Top: raw[0],
        Bottom: raw[1],
        Left: raw[2],
        Right: raw[3],
    }, err
}

// _NET_WM_FULLSCREEN_MONITORS set
func (xu *XUtil) EwmhWmFullscreenMonitorsSet(win xgb.Id,
                                             edges WmFullscreenMonitors) error {
    raw := make([]uint32, 4)
    raw[0] = edges.Top
    raw[1] = edges.Bottom
    raw[2] = edges.Left
    raw[3] = edges.Right

    return xu.ChangeProperty32(win, "_NET_WM_FULLSCREEN_MONITORS", "CARDINAL",
                               raw...)
}

// _NET_WM_FULLSCREEN_MONITORS req
func (xu *XUtil) EwmhWmFullscreenMonitorsReq(win xgb.Id,
                                             edges WmFullscreenMonitors) error {
    return xu.EwmhWmFullscreenMonitorsReqExtra(win, edges, 2)
}

// _NET_WM_FULLSCREEN_MONITORS req extra
func (xu *XUtil) EwmhWmFullscreenMonitorsReqExtra(win xgb.Id,
                                                  edges WmFullscreenMonitors,
                                                  source uint32) error {
    return xu.EwmhClientEvent(win, "_NET_WM_FULLSCREEN_MONITORS",
                              edges.Top, edges.Bottom, edges.Left, edges.Right,
                              source)
}

// _NET_WM_HANDLED_ICONS get
func (xu *XUtil) EwmhWmHandledIcons(win xgb.Id) (bool, error) {
    reply, err := xu.GetProperty(win, "_NET_WM_HANDLED_ICONS")
    if err != nil {
        return false, err
    }

    val, err := PropValNum(reply, nil)
    if err != nil {
        return false, err
    }

    return val == 1, nil
}

// _NET_WM_HANDLED_ICONS set
func (xu *XUtil) EwmhWmHandledIconsSet(handle bool) error {
    var handled uint32
    if handle {
        handled = 1
    } else {
        handled = 0
    }
    return xu.ChangeProperty32(xu.root, "_NET_WM_HANDLED_ICONS", "CARDINAL",
                               handled)
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
func (xu *XUtil) EwmhWmIcon(win xgb.Id) ([]WmIcon, error) {
    icon, err := PropValNums(xu.GetProperty(win, "_NET_WM_ICON"))
    if err != nil {
        return nil, err
    }

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

    return wmicons, nil
}

// _NET_WM_ICON set
func (xu *XUtil) EwmhWmIconSet(win xgb.Id, icons []WmIcon) error {
    raw := make([]uint32, 0, 10000) // start big
    for _, icon := range icons {
        raw = append(raw, icon.Width, icon.Height)
        raw = append(raw, icon.Data...)
    }

    return xu.ChangeProperty32(win, "_NET_WM_ICON", "CARDINAL", raw...)
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
func (xu *XUtil) EwmhWmIconGeometry(win xgb.Id) (WmIconGeometry, error) {
    geom, err := PropValNums(xu.GetProperty(win, "_NET_WM_ICON_GEOMETRY"))
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
func (xu *XUtil) EwmhWmIconGeometrySet(win xgb.Id, geom WmIconGeometry) error {
    rawGeom := make([]uint32, 4)
    rawGeom[0] = geom.X
    rawGeom[1] = geom.Y
    rawGeom[2] = geom.Width
    rawGeom[3] = geom.Height

    return xu.ChangeProperty32(win, "_NET_WM_ICON_GEOMETRY", "CARDINAL",
                               rawGeom...)
}

// _NET_WM_ICON_NAME get
func (xu *XUtil) EwmhWmIconName(win xgb.Id) (string, error) {
    return PropValStr(xu.GetProperty(win, "_NET_WM_ICON_NAME"))
}

// _NET_WM_ICON_NAME set
func (xu *XUtil) EwmhWmIconNameSet(win xgb.Id, name string) error {
    return xu.ChangeProperty(win, 8, "_NET_WM_ICON_NAME", "UTF8_STRING",
                             []byte(name))
}

// _NET_WM_MOVERESIZE constants
const (
    EwmhSizeTopLeft = iota
    EwmhSizeTop
    EwmhSizeTopRight
    EwmhSizeRight
    EwmhSizeBottomRight
    EwmhSizeBottom
    EwmhSizeBottomLeft
    EwmhSizeLeft
    EwmhMove
    EwmhSizeKeyboard
    EwmhMoveKeyboard
    EwmhCancel
)

// _NET_WM_MOVERESIZE req
func (xu *XUtil) EwmhWmMoveresize(win xgb.Id, direction uint32) error {
    return xu.EwmhWmMoveresizeExtra(win, direction, 0, 0, 0, 2)
}

// _NET_WM_MOVERESIZE req extra
func (xu *XUtil) EwmhWmMoveresizeExtra(win xgb.Id, direction, x_root, y_root,
                                       button, source uint32) error {
    return xu.EwmhClientEvent(win, "_NET_WM_MOVERESIZE", x_root, y_root,
                              direction, button, source)
}

// _NET_WM_NAME get
func (xu *XUtil) EwmhWmName(win xgb.Id) (string, error) {
    return PropValStr(xu.GetProperty(win, "_NET_WM_NAME"))
}

// _NET_WM_NAME set
func (xu *XUtil) EwmhWmNameSet(win xgb.Id, name string) error {
    return xu.ChangeProperty(win, 8, "_NET_WM_NAME", "UTF8_STRING",
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
func (xu *XUtil) EwmhWmOpaqueRegion(win xgb.Id) ([]WmOpaqueRegion, error) {
    raw, err := PropValNums(xu.GetProperty(win, "_NET_WM_OPAQUE_REGION"))
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
func (xu *XUtil) EwmhWmOpaqueRegionSet(win xgb.Id,
                                       regions []WmOpaqueRegion) error {
    raw := make([]uint32, len(regions) * 4)

    for i, region := range(regions) {
        raw[i * 4 + 0] = region.X
        raw[i * 4 + 1] = region.Y
        raw[i * 4 + 2] = region.Width
        raw[i * 4 + 3] = region.Height
    }

    return xu.ChangeProperty32(win, "_NET_WM_OPAQUE_REGION", "CARDINAL", raw...)
}

// _NET_WM_PID get
func (xu *XUtil) EwmhWmPid(win xgb.Id) (uint32, error) {
    return PropValNum(xu.GetProperty(win, "_NET_WM_PID"))
}

// _NET_WM_PID set
func (xu *XUtil) EwmhWmPidSet(win xgb.Id, pid uint32) error {
    return xu.ChangeProperty32(win, "_NET_WM_PID", "CARDINAL", pid)
}

// _NET_WM_PING req
func (xu *XUtil) EwmhWmPing(win xgb.Id, response bool) error {
    return xu.EwmhWmPingExtra(win, response, 0)
}

// _NET_WM_PING req extra
func (xu *XUtil) EwmhWmPingExtra(win xgb.Id, response bool,
                                 time xgb.Timestamp) error {
    pingAtom, err := xu.Atm("_NET_WM_PING")
    if err != nil {
        return err
    }

    var evWindow xgb.Id
    if response {
        evWindow = xu.root
    } else {
        evWindow = win
    }

    return xu.EwmhClientEvent(evWindow, "WM_PROTOCOLS", uint32(pingAtom), time,
                              win)
}

// _NET_WM_STATE constants for state toggling
// These correspond to the "action" parameter.
const (
    EwmhStateRemove = iota
    EwmhStateAdd
    EwmhStateToggle
)

// _NET_WM_STATE get
func (xu *XUtil) EwmhWmState(win xgb.Id) ([]string, error) {
    return xu.PropValAtoms(xu.GetProperty(win, "_NET_WM_STATE"))
}

// _NET_WM_STATE set
func (xu *XUtil) EwmhWmStateSet(win xgb.Id, atomNames []string) error {
    atoms, err := xu.StrToAtoms(atomNames)
    if err != nil {
        return err
    }

    return xu.ChangeProperty32(win, "_NET_WM_STATE", "ATOM", atoms...)
}

// _NET_WM_STATE req
func (xu *XUtil) EwmhWmStateReq(win xgb.Id, action uint32,
                                atomName string) error {
    return xu.EwmhWmStateReqExtra(win, action, atomName, "", 2)
}

// _NET_WM_STATE req extra
func (xu *XUtil) EwmhWmStateReqExtra(win xgb.Id, action uint32,
                                     first string, second string,
                                     source uint32) (err error) {
    var atom1, atom2 xgb.Id

    atom1, err = xu.Atom(first, false)
    if err != nil {
        return err
    }

    if len(second) > 0 {
        atom2, err = xu.Atom(second, false)
        if err != nil {
            return err
        }
    } else {
        atom2 = 0
    }

    return xu.EwmhClientEvent(win, "_NET_WM_STATE", action, uint32(atom1),
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
func (xu *XUtil) EwmhWmStrut(win xgb.Id) (WmStrut, error) {
    struts, err := PropValNums(xu.GetProperty(win, "_NET_WM_STRUT"))
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
func (xu *XUtil) EwmhWmStrutSet(win xgb.Id, struts WmStrut) error {
    rawStruts := make([]uint32, 4)
    rawStruts[0] = struts.Left
    rawStruts[1] = struts.Right
    rawStruts[2] = struts.Top
    rawStruts[3] = struts.Bottom

    return xu.ChangeProperty32(win, "_NET_WM_STRUT", "CARDINAL", rawStruts...)
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
func (xu *XUtil) EwmhWmStrutPartial(win xgb.Id) (WmStrutPartial, error) {
    struts, err := PropValNums(xu.GetProperty(win, "_NET_WM_STRUT_PARTIAL"))
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
func (xu *XUtil) EwmhWmStrutPartialSet(win xgb.Id,
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

    return xu.ChangeProperty32(win, "_NET_WM_STRUT_PARTIAL", "CARDINAL",
                               rawStruts...)
}

// _NET_WM_SYNC_REQUEST req
func (xu *XUtil) EwmhWmSyncRequest(win xgb.Id, req_num uint64) error {
    return xu.EwmhWmSyncRequestExtra(win, req_num, 0)
}

// _NET_WM_SYNC_REQUEST req extra
func (xu *XUtil) EwmhWmSyncRequestExtra(win xgb.Id, req_num uint64,
                                        time xgb.Timestamp) error {
    syncReq, err := xu.Atm("_NET_WM_SYNC_REQUEST")
    if err != nil {
        return err
    }

    high := uint32(req_num >> 32)
    low := uint32(req_num << 32 ^ req_num)

    return xu.EwmhClientEvent(win, "WM_PROTOCOLS", syncReq, time, low, high)
}

// _NET_WM_SYNC_REQUEST_COUNTER get 
// I'm pretty sure this needs 64 bit integers, but I'm not quite sure
// how to go about that yet. Any ideas?
func (xu *XUtil) EwmhWmSyncRequestCounter(win xgb.Id) (uint32, error) {
    return PropValNum(xu.GetProperty(win, "_NET_WM_SYNC_REQUEST_COUNTER"))
}

// _NET_WM_SYNC_REQUEST_COUNTER set
// I'm pretty sure this needs 64 bit integers, but I'm not quite sure
// how to go about that yet. Any ideas?
func (xu *XUtil) EwmhWmSyncRequestCounterSet(win xgb.Id, counter uint32) error {
    return xu.ChangeProperty32(win, "_NET_WM_SYNC_REQUEST_COUNTER", "CARDINAL",
                               counter)
}

// _NET_WM_USER_TIME get
func (xu *XUtil) EwmhWmUserTime(win xgb.Id) (uint32, error) {
    return PropValNum(xu.GetProperty(win, "_NET_WM_USER_TIME"))
}

// _NET_WM_USER_TIME set
func (xu *XUtil) EwmhWmUserTimeSet(win xgb.Id, user_time uint32) error {
    return xu.ChangeProperty32(win, "_NET_WM_USER_TIME", "CARDINAL", user_time)
}

// _NET_WM_USER_TIME_WINDOW get
func (xu *XUtil) EwmhWmUserTimeWindow(win xgb.Id) (xgb.Id, error) {
    return PropValId(xu.GetProperty(win, "_NET_WM_USER_TIME_WINDOW"))
}

// _NET_WM_USER_TIME set
func (xu *XUtil) EwmhWmUserTimeWindowSet(win xgb.Id, time_win xgb.Id) error {
    return xu.ChangeProperty32(win, "_NET_WM_USER_TIME_WINDOW", "CARDINAL",
                               uint32(time_win))
}

// _NET_WM_VISIBLE_ICON_NAME get
func (xu *XUtil) EwmhWmVisibleIconName(win xgb.Id) (string, error) {
    return PropValStr(xu.GetProperty(win, "_NET_WM_VISIBLE_ICON_NAME"))
}

// _NET_WM_VISIBLE_ICON_NAME set
func (xu *XUtil) EwmhWmVisibleIconNameSet(win xgb.Id, name string) error {
    return xu.ChangeProperty(win, 8, "_NET_WM_VISIBLE_ICON_NAME", "UTF8_STRING",
                             []byte(name))
}

// _NET_WM_VISIBLE_NAME get
func (xu *XUtil) EwmhWmVisibleName(win xgb.Id) (string, error) {
    return PropValStr(xu.GetProperty(win, "_NET_WM_VISIBLE_NAME"))
}

// _NET_WM_VISIBLE_NAME set
func (xu *XUtil) EwmhWmVisibleNameSet(win xgb.Id, name string) error {
    return xu.ChangeProperty(win, 8, "_NET_WM_VISIBLE_NAME", "UTF8_STRING",
                             []byte(name))
}

// _NET_WM_WINDOW_OPACITY get
// This isn't part of the EWMH spec, but is widely used by drop in
// compositing managers (i.e., xcompmgr, cairo-compmgr, etc.).
// This property is typically set not on a client window, but the *parent*
// of a client window in reparenting window managers.
func (xu *XUtil) EwmhWmWindowOpacity(win xgb.Id) (float64, error) {
    intOpacity, err := PropValNum(xu.GetProperty(win, "_NET_WM_WINDOW_OPACITY"))
    if err != nil {
        return 0, err
    }

    return float64(intOpacity) / float64(0xffffffff), nil
}

// _NET_WM_WINDOW_OPACITY set
func (xu *XUtil) EwmhWmWindowOpacitySet(win xgb.Id, opacity float64) error {
    return xu.ChangeProperty32(win, "_NET_WM_WINDOW_OPACITY", "CARDINAL",
                               uint32(opacity * 0xffffffff))
}

// _NET_WM_WINDOW_TYPE get
func (xu *XUtil) EwmhWmWindowType(win xgb.Id) ([]string, error) {
    return xu.PropValAtoms(xu.GetProperty(win, "_NET_WM_WINDOW_TYPE"))
}

// _NET_WM_WINDOW_TYPE set
// This will create any atoms used in 'atomNames' if they don't already exist.
func (xu *XUtil) EwmhWmWindowTypeSet(win xgb.Id, atomNames []string) error {
    atoms, err := xu.StrToAtoms(atomNames)
    if err != nil {
        return err
    }

    return xu.ChangeProperty32(win, "_NET_WM_WINDOW_TYPE", "ATOM", atoms...)
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
func (xu *XUtil) EwmhWorkarea() ([]Workarea, error) {
    rects, err := PropValNums(xu.GetProperty(xu.root, "_NET_WORKAREA"))
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
func (xu *XUtil) EwmhWorkareaSet(workareas []Workarea) error {
    rects := make([]uint32, len(workareas) * 4)
    for i, workarea := range workareas {
        rects[i * 4] = workarea.X
        rects[i * 4 + 1] = workarea.Y
        rects[i * 4 + 2] = workarea.Width
        rects[i * 4 + 3] = workarea.Height
    }

    return xu.ChangeProperty32(xu.root, "_NET_WORKAREA", "CARDINAL", rects...)
}

