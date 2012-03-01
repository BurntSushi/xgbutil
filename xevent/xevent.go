/*
    event.go organizes functions related to X event processing.

    (Note: The 'evtypes.go' file is related. It defines methods/interfaces
           for dealing with common X events. event.go on the other hand,
           is more concerned with responding to events that X sends us.)
*/
package xevent

import "code.google.com/p/x-go-binding/xgb"
import "github.com/BurntSushi/xgbutil"

// ReplayPointer is a quick alias to AllowEvents with 'ReplayPointer' mode.
func ReplayPointer(xu *xgbutil.XUtil) {
    xu.Conn().AllowEvents(xgb.AllowReplayPointer, 0)
}

