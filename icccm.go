/*
    Provides an API for most of the ICCCM spec[1].

    This API follows most of the same conventions as the corresponding EWMH
    API, except that all methods are prefixed with 'Icccm' instead of 'Ewmh'.

    Also, I believe there are a few client root events specified by the ICCCM,
    but I haven't put them in here (yet). I am not sure that they are even
    used any more.

    [1] - http://tronche.com/gui/x/icccm/
*/
package xgbutil

import "code.google.com/p/x-go-binding/xgb"

const (
    HintInput = (1 << iota)
    HintState
    HintIconPixmap
    HintIconWindow
    HintIconPosition
    HintIconMask
    HintWindowGroup
    HintMessage
    HintUrgency
)

const (
    SizeHintUSPosition = (1 << iota)
    SizeHintUSSize
    SizeHintPPosition
    SizeHintPMinSize
    SizeHintPMaxSize
    SizeHintPResizeInc
    SizeHintPAspect
    SizeHintPBaseSize
    SizeHintPWinGravity
)

const (
    StateWithdrawn = iota
    StateNormal
    StateZoomed
    StateIconic
    StateInactive
)

// WM_NAME get
func (xu *XUtil) IcccmWmName(win xgb.Id) (string, error) {
    return PropValStr(xu.GetProperty(win, "WM_NAME"))
}

// WM_NAME set
func (xu *XUtil) IcccmWmNameSet(win xgb.Id, name string) error {
    return xu.ChangeProperty(win, 8, "WM_NAME", "STRING", ([]byte)(name))
}

// WM_ICON_NAME get
func (xu *XUtil) IcccmWmIconName(win xgb.Id) (string, error) {
    return PropValStr(xu.GetProperty(win, "WM_ICON_NAME"))
}

// WM_ICON_NAME set
func (xu *XUtil) IcccmWmIconNameSet(win xgb.Id, name string) error {
    return xu.ChangeProperty(win, 8, "WM_ICON_NAME", "STRING", ([]byte)(name))
}

// NormalHints is a struct that organizes the information related to the
// WM_NORMAL_HINTS property. Please see the ICCCM spec for more details.
type NormalHints struct {
    Flags uint32
    X, Y, Width, Height, MinWidth, MinHeight, MaxWidth, MaxHeight uint32
    WidthInc, HeightInc uint32
    MinAspectNum, MinAspectDen, MaxAspectNum, MaxAspectDen uint32
    BaseWidth, BaseHeight, WinGravity uint32
}

// WM_NORMAL_HINTS get
func (xu *XUtil) IcccmWmNormalHints(win xgb.Id) (nh NormalHints, err error) {
    lenExpect := 18
    hints, err := PropValNums(xu.GetProperty(win, "WM_NORMAL_HINTS"))
    if err != nil {
        return NormalHints{}, err
    }
    if len(hints) != lenExpect {
        return NormalHints{}, xuerr("IcccmWmNormalHints",
                                    "There are %d fields in " +
                                    "WM_NORMAL_HINTS, but xgbutil expects %d.",
                                    len(hints), lenExpect)
    }

    nh.Flags = hints[0]
    nh.X = hints[1]
    nh.Y = hints[2]
    nh.Width = hints[3]
    nh.Height = hints[4]
    nh.MinWidth = hints[5]
    nh.MinHeight = hints[6]
    nh.MaxWidth = hints[7]
    nh.MaxHeight = hints[8]
    nh.WidthInc = hints[9]
    nh.HeightInc = hints[10]
    nh.MinAspectNum = hints[11]
    nh.MinAspectDen = hints[12]
    nh.MaxAspectNum = hints[13]
    nh.MaxAspectDen = hints[14]
    nh.BaseWidth = hints[15]
    nh.BaseHeight = hints[16]
    nh.WinGravity = hints[17]

    if nh.WinGravity <= 0 {
        nh.WinGravity = xgb.GravityNorthWest
    }

    return nh, nil
}

// WM_NORMAL_HINTS set
// Make sure to set the flags in the NormalHints struct correctly!
func (xu *XUtil) IcccmWmNormalHintsSet(win xgb.Id, nh NormalHints) error {
    raw := []uint32{
        nh.Flags, nh.X, nh.Y, nh.Width, nh.Height,
        nh.MinWidth, nh.MinHeight, nh.MaxWidth, nh.MaxHeight,
        nh.WidthInc, nh.HeightInc,
        nh.MinAspectNum, nh.MinAspectDen, nh.MaxAspectNum, nh.MaxAspectDen,
        nh.BaseWidth, nh.BaseHeight, nh.WinGravity,
    }
    return xu.ChangeProperty32(win, "WM_NORMAL_HINTS", "WM_SIZE_HINTS", raw...)
}

// Hints is a struct that organizes information related to the WM_HINTS
// property. Once again, I refer you to the ICCCM spec for documentation.
type Hints struct {
    Flags uint32
    Input, InitialState, IconX, IconY, WindowGroup uint32
    IconPixmap, IconWindow, IconMask xgb.Id
}

// WM_HINTS get
func (xu *XUtil) IcccmWmHints(win xgb.Id) (hints Hints, err error) {
    lenExpect := 9
    raw, err := PropValNums(xu.GetProperty(win, "WM_HINTS"))
    if err != nil {
        return Hints{}, err
    }
    if len(raw) != lenExpect {
        return Hints{}, xuerr("IcccmWmHints",
                              "There are %d fields in " +
                              "WM_HINTS, but xgbutil expects %d.",
                              len(raw), lenExpect)
    }

    hints.Flags = raw[0]
    hints.Input = raw[1]
    hints.InitialState = raw[2]
    hints.IconPixmap = xgb.Id(raw[3])
    hints.IconWindow = xgb.Id(raw[4])
    hints.IconX = raw[5]
    hints.IconY = raw[6]
    hints.IconMask = xgb.Id(raw[7])
    hints.WindowGroup = raw[8]

    return hints, nil
}

// WM_HINTS set
// Make sure to set the flags in the Hints struct correctly!
func (xu *XUtil) IcccmWmHintsSet(win xgb.Id, hints Hints) error {
    raw := []uint32{
        hints.Flags, hints.Input, hints.InitialState,
        uint32(hints.IconPixmap), uint32(hints.IconWindow),
        hints.IconX, hints.IconY,
        uint32(hints.IconMask),
        hints.WindowGroup,
    }
    return xu.ChangeProperty32(win, "WM_HINTS", "WM_HINTS", raw...)
}

// WmClass struct contains two data points:
// the instance and a class of a window.
type WmClass struct {
    Instance, Class string
}

// WM_CLASS get
func (xu *XUtil) IcccmWmClass(win xgb.Id) (WmClass, error) {
    raw, err := PropValStrs(xu.GetProperty(win, "WM_CLASS"))
    if err != nil {
        return WmClass{}, err
    }
    if len(raw) != 2 {
        return WmClass{}, xuerr("IcccmWmClass",
                                "Two string make up WM_CLASS, but " +
                                "xgbutil found %d in '%v'.", len(raw), raw)
    }

    return WmClass {
        Instance: raw[0],
        Class: raw[1],
    }, nil
}

// WM_CLASS set
func (xu *XUtil) IcccmWmClassSet(win xgb.Id, class WmClass) error {
    raw := make([]byte, len(class.Instance) + len(class.Class) + 2)
    copy(raw, class.Instance)
    copy(raw[(len(class.Instance) + 1):], class.Class)

    return xu.ChangeProperty(win, 8, "WM_CLASS", "STRING", raw)
}

// WM_TRANSIENT_FOR get
func (xu *XUtil) IcccmWmTransientFor(win xgb.Id) (xgb.Id, error) {
    return PropValId(xu.GetProperty(win, "WM_TRANSIENT_FOR"))
}

// WM_TRANSIENT_FOR set
func (xu *XUtil) IcccmWmTransientForSet(win xgb.Id, transient xgb.Id) error {
    return xu.ChangeProperty32(win, "WM_TRANSIENT_FOR", "WINDOW",
                               uint32(transient))
}

// WM_PROTOCOLS get
func (xu *XUtil) IcccmWmProtocols(win xgb.Id) ([]string, error) {
    return xu.PropValAtoms(xu.GetProperty(win, "WM_PROTOCOLS"))
}

// WM_PROTOCOLS set
func (xu *XUtil) IcccmWmProtocolsSet(win xgb.Id, atomNames []string) error {
    atoms, err := xu.StrToAtoms(atomNames)
    if err != nil {
        return err
    }

    return xu.ChangeProperty32(win, "WM_PROTOCOLS", "ATOM", atoms...)
}

// WM_COLORMAP_WINDOWS get
func (xu *XUtil) IcccmWmColormapWindows(win xgb.Id) ([]xgb.Id, error) {
    return PropValIds(xu.GetProperty(win, "WM_COLORMAP_WINDOWS"))
}

// WM_COLORMAP_WINDOWS set
func (xu *XUtil) IcccmWmColormapWindowsSet(win xgb.Id, windows []xgb.Id) error {
    return xu.ChangeProperty32(win, "WM_COLORMAP_WINDOWS", "WINDOW",
                               IdTo32(windows)...)
}

// WM_CLIENT_MACHINE get
func (xu *XUtil) IcccmWmClientMachine(win xgb.Id) (string, error) {
    return PropValStr(xu.GetProperty(win, "WM_CLIENT_MACHINE"))
}

// WM_CLIENT_MACHINE set
func (xu *XUtil) IcccmWmClientMachineSet(win xgb.Id, client string) error {
    return xu.ChangeProperty(win, 8, "WM_CLIENT_MACHINE", "STRING",
                             ([]byte)(client))
}

// WmState is a struct that organizes information related to the WM_STATE
// property. Namely, the state (corresponding to a State* constant in this file)
// and the icon window (probably not used).
type WmState struct {
    State uint32
    Icon xgb.Id
}

// WM_STATE get
func (xu *XUtil) IcccmWmState(win xgb.Id) (WmState, error) {
    raw, err := PropValNums(xu.GetProperty(win, "WM_STATE"))
    if err != nil {
        return WmState{}, err
    }
    if len(raw) != 2 {
        return WmState{}, xuerr("IcccmWmState",
                                "Expected two integers in WM_STATE property " +
                                "but xgbutil found %d in '%v'.", len(raw), raw)
    }

    return WmState{
        State: raw[0],
        Icon: xgb.Id(raw[1]),
    }, nil
}

// WM_STATE set
func (xu *XUtil) IcccmWmStateSet(win xgb.Id, state WmState) error {
    raw := []uint32{
        state.State,
        uint32(state.Icon),
    }

    return xu.ChangeProperty32(win, "WM_STATE", "WM_STATE", raw...)
}

// IconSize is a struct the organizes information related to the WM_ICON_SIZE
// property. Mostly info about its dimensions.
type IconSize struct {
    MinWidth, MinHeight, MaxWidth, MaxHeight, WidthInc, HeightInc uint32
}

// WM_ICON_SIZE get
func (xu *XUtil) IcccmWmIconSize(win xgb.Id) (IconSize, error) {
    raw, err := PropValNums(xu.GetProperty(win, "WM_ICON_SIZE"))
    if err != nil {
        return IconSize{}, err
    }
    if len(raw) != 6 {
        return IconSize{}, xuerr("IcccmWmIconSize",
                                 "Expected six integers in WM_ICON_SIZE " +
                                 "property, but xgbutil found " +
                                 "%d in '%v'.", len(raw), raw)
    }

    return IconSize{
        MinWidth: raw[0], MinHeight: raw[1],
        MaxWidth: raw[2], MaxHeight: raw[3],
        WidthInc: raw[4], HeightInc: raw[5],
    }, nil
}

// WM_ICON_SIZE set
func (xu *XUtil) IcccmWmIconSizeSet(win xgb.Id, icondim IconSize) error {
    raw := []uint32{
        icondim.MinWidth, icondim.MinHeight,
        icondim.MaxWidth, icondim.MaxHeight,
        icondim.WidthInc, icondim.HeightInc,
    }

    return xu.ChangeProperty32(win, "WM_ICON_SIZE", "WM_ICON_SIZE", raw...)
}

