/*
    Package mousebind provides an easy interface to bind and run callback
    functions for human readable mouse bindings.
*/
package mousebind

import "log"
import "strconv"
import "strings"

import "burntsushi.net/go/x-go-binding/xgb"

import (
    "burntsushi.net/go/xgbutil"
    "burntsushi.net/go/xgbutil/xevent"
)

var modifiers []uint16 = []uint16{ // order matters!
    xgb.ModMaskShift, xgb.ModMaskLock, xgb.ModMaskControl,
    xgb.ModMask1, xgb.ModMask2, xgb.ModMask3, xgb.ModMask4, xgb.ModMask5,
    xgb.ButtonMask1, xgb.ButtonMask2, xgb.ButtonMask3,
    xgb.ButtonMask4, xgb.ButtonMask5,
    xgb.ButtonMaskAny,
}

var pointerMasks uint16 = xgb.EventMaskPointerMotion |
                          xgb.EventMaskButtonRelease |
                          xgb.EventMaskButtonPress

// Initialize attaches the appropriate callbacks to make mouse bindings easier.
// i.e., prep the dummy window to handle mouse dragging events
func Initialize(xu *xgbutil.XUtil) {
    xevent.MotionNotifyFun(dragStep).Connect(xu, xu.Dummy())
    xevent.ButtonReleaseFun(dragEnd).Connect(xu, xu.Dummy())
}

// ParseString takes a string of the format '[Mod[-Mod[...]]-]-KEY',
// i.e., 'Mod4-1', and returns a modifiers/button combo.
// "Mod" could also be one of {button1, button2, button3, button4, button5}.
// (Actually, the parser is slightly more forgiving than what this comment
//  leads you to believe.)
func ParseString(xu *xgbutil.XUtil, str string) (uint16, byte) {
    mods, button := uint16(0), byte(0)
    for _, part := range strings.Split(str, "-") {
        switch(strings.ToLower(part)) {
        case "shift":
            mods |= xgb.ModMaskShift
        case "lock":
            mods |= xgb.ModMaskLock
        case "control":
            mods |= xgb.ModMaskControl
        case "mod1":
            mods |= xgb.ModMask1
        case "mod2":
            mods |= xgb.ModMask2
        case "mod3":
            mods |= xgb.ModMask3
        case "mod4":
            mods |= xgb.ModMask4
        case "mod5":
            mods |= xgb.ModMask5
        case "button1":
            mods |= xgb.ButtonMask1
        case "button2":
            mods |= xgb.ButtonMask2
        case "button3":
            mods |= xgb.ButtonMask3
        case "button4":
            mods |= xgb.ButtonMask4
        case "button5":
            mods |= xgb.ButtonMask5
        case "any":
            mods |= xgb.ButtonMaskAny
        default: // a button!
            if button == 0 { // only accept the first button we see
                possible, err := strconv.ParseUint(part, 10, 8)
                if err == nil {
                    button = byte(possible)
                } else {
                    log.Printf("We could not convert '%s' to a valid 8-bit " +
                               "integer. Assuming 0.\n", part)
                }
            }
        }
    }

    if button == 0 {
        log.Printf("We could not find a valid button in the string '%s'. " +
                   "Things probably will not work right.\n", str)
    }

    return mods, button
}

// Grabs a button with mods on a particular window.
// Will also grab all combinations of modifiers found in xgbutil.IgnoreMods
// If 'propagate' is True, then no further events can be processed until the
// grabbing client allows them to be. (Which is done via AllowEvents. Thus,
// if propagate is True, you *must* make some call to AllowEvents at some
// point, or else your client will lock.)
func Grab(xu *xgbutil.XUtil, win xgb.Id, mods uint16, button byte,
          propagate bool) {
    var pSync byte = xgb.GrabModeAsync
    if propagate {
        pSync = xgb.GrabModeSync
    }

    for _, m := range xgbutil.IgnoreMods {
        xu.Conn().GrabButton(true, win, pointerMasks, pSync,
                             xgb.GrabModeAsync, 0, 0, button, mods | m)
    }
}

// Ungrab undoes Grab. It will handle all combinations od modifiers found
// in xgbutil.IgnoreMods.
func Ungrab(xu *xgbutil.XUtil, win xgb.Id, mods uint16, button byte) {
    for _, m := range xgbutil.IgnoreMods {
        xu.Conn().UngrabButton(button, win, mods | m)
    }
}

// GrabPointer grabs the entire pointer.
// Returns whether GrabStatus is successful and an error if one is reported by 
// XGB. It is possible to not get an error and the grab to be unsuccessful.
// The purpose of 'win' is that after a grab is successful, ALL Button*Events 
// will be sent to that window. Make sure you have a callback attached :-)
func GrabPointer(xu *xgbutil.XUtil, win xgb.Id, confine xgb.Id,
                 cursor xgb.Id) (bool, error) {
    reply, err := xu.Conn().GrabPointer(false, win, pointerMasks,
                                        xgb.GrabModeAsync, xgb.GrabModeAsync,
                                        confine, cursor, 0)
    if err != nil {
        return false, xgbutil.Xerr(err, "GrabPointer",
                                   "Error grabbing pointer on window '%x'",
                                   win)
    }

    return reply.Status == xgb.GrabStatusSuccess, nil
}

// UngrabPointer undoes GrabPointer.
func UngrabPointer(xu *xgbutil.XUtil) {
    xu.Conn().UngrabPointer(0)
}

