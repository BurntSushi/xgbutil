/*
    event.go organizes functions related to X event processing.

    (Note: The 'evtypes.go' file is related. It defines methods/interfaces
           for dealing with common X events. event.go on the other hand,
           is more concerned with responding to events that X sends us.)
*/
package xgbutil

import "code.google.com/p/x-go-binding/xgb"

// ReplayPointer is a quick alias to AllowEvents with 'ReplayPointer' mode.
func (xu *XUtil) ReplayPointer() {
    xu.conn.AllowEvents(xgb.AllowReplayPointer, 0)
}

