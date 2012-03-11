/*
    Package motif has a few functions to allow easy access to Motif related
    properties.

    The main purpose here is that some applications communicate "no window
    decorations" to the window manager using _MOTIF_WM_HINTS. (Like Google
    Chrome.) I haven't seen Motif stuff used for other purposes in the wild for 
    a long time.

    As a result, only the useful bits are implemented here. More may be added
    on an on-demand basis, but don't count on it.

    Try not to bash your head against your desk too hard:
    http://www.opengroup.org/openmotif/hardcopydocs.html
*/
package motif

import "code.google.com/p/jamslam-x-go-binding/xgb"
import "github.com/BurntSushi/xgbutil"
import "github.com/BurntSushi/xgbutil/xprop"

const (
    HintFunctions = (1 << iota)
    HintDecorations
    HintInputMode
    HintStatus
)

const (
    FunctionAll = (1 << iota)
    FunctionResize
    FunctionMove
    FunctionMinimize
    FunctionMaximize
    FunctionClose
    FunctionNone = 0
)

const (
    DecorationAll = (1 << iota)
    DecorationBorder
    DecorationResizeH
    DecorationTitle
    DecorationMenu
    DecorationMinimize
    DecorationMaximize
    DecorationNone = 0
)

const (
    InputPrimaryApplicationModal = (1 << iota)
    InputSystemModal
    InputFullApplicationModal
    InputModeless = 0
)

const StatusTearoffWindow = 1

// Hints is a struct that organizes the information related to the
// WM_NORMAL_HINTS property.
type Hints struct {
    Flags uint32
    Function, Decoration, Input, Status uint32
}

// _MOTIF_WM_HINTS get
func WmHintsGet(xu *xgbutil.XUtil, win xgb.Id) (mh Hints, err error) {
    lenExpect := 5
    hints, err := xprop.PropValNums(xprop.GetProperty(xu, win,
                                                      "_MOTIF_WM_HINTS"))
    if err != nil {
        return Hints{}, err
    }
    if len(hints) != lenExpect {
        return Hints{},
               xgbutil.Xuerr("motif.WmHintsGet",
                             "There are %d fields in _MOTIF_WM_HINTS, " +
                             "but xgbutil expects %d.", len(hints), lenExpect)
    }

    mh.Flags = hints[0]
    mh.Function = hints[1]
    mh.Decoration = hints[2]
    mh.Input = hints[3]
    mh.Status = hints[4]

    return
}

// _MOTIF_WM_HINTS set
func WmHintsSet(xu *xgbutil.XUtil, win xgb.Id, mh Hints) error {
    raw := []uint32{mh.Flags, mh.Function, mh.Decoration, mh.Input, mh.Status}
    return xprop.ChangeProp32(xu, win, "_MOTIF_WM_HINTS", "_MOTIF_WM_HINTS",
                              raw...)
}

