/*
    event.go organizes functions related to X event processing.

    (Note: The 'evtypes.go' file is related. It defines methods/interfaces
           for dealing with common X events. event.go on the other hand,
           is more concerned with responding to events that X sends us.)
*/
package xevent

import "log"
import "code.google.com/p/jamslam-x-go-binding/xgb"
import "github.com/BurntSushi/xgbutil"

// Main starts the main X event loop. It will read events and call appropriate
// callback functions. Note that xgbutil builds in a few callbacks of its own,
// particularly the MappingNotify event so that the key mapping and
// modifier mapping can be automatically updated.
// XXX: This only supports using one X connection. It should at least allow
//      some arbitrary number of connections. Not sure what the best approach
//      is, but I'm sure it needs to use channels (and multiple calls to
//      Main using select, by the user).
func Main(xu *xgbutil.XUtil) error {
    for {
        // XXX: Use PollForEvent to gobble up lots of events.
        reply, err := xu.Conn().WaitForEvent()
        if err != nil {
            log.Printf("ERROR: %v\n", err)
            continue
        }

        // We have to look for xgb events here. But we re-wrap them in our
        // own event types.
        switch event := reply.(type) {
        case xgb.KeyPressEvent:
            e := KeyPressEvent{&event}
            xu.RunCallbacks(e, KeyPress, e.Event)
        case xgb.KeyReleaseEvent:
            e := KeyReleaseEvent{&event}
            xu.RunCallbacks(e, KeyRelease, e.Event)
        case xgb.ButtonPressEvent:
            e := ButtonPressEvent{&event}
            xu.RunCallbacks(e, ButtonPress, e.Event)
        case xgb.ButtonReleaseEvent:
            e := ButtonReleaseEvent{&event}
            xu.RunCallbacks(e, ButtonRelease, e.Event)
        case xgb.MotionNotifyEvent:
            e := MotionNotifyEvent{&event}
            xu.RunCallbacks(e, MotionNotify, e.Event)
        case xgb.EnterNotifyEvent:
            e := EnterNotifyEvent{&event}
            xu.RunCallbacks(e, EnterNotify, e.Event)
        case xgb.LeaveNotifyEvent:
            e := LeaveNotifyEvent{&event}
            xu.RunCallbacks(e, LeaveNotify, e.Event)
        case xgb.FocusInEvent:
            e := FocusInEvent{&event}
            xu.RunCallbacks(e, FocusIn, e.Event)
        case xgb.FocusOutEvent:
            e := FocusOutEvent{&event}
            xu.RunCallbacks(e, FocusOut, e.Event)
        case xgb.KeymapNotifyEvent:
            e := KeymapNotifyEvent{&event}
            xu.RunCallbacks(e, KeymapNotify, xgbutil.NoWindow)
        case xgb.ExposeEvent:
            e := ExposeEvent{&event}
            xu.RunCallbacks(e, Expose, e.Window)
        case xgb.GraphicsExposureEvent:
            e := GraphicsExposureEvent{&event}
            xu.RunCallbacks(e, GraphicsExposure, e.Drawable)
        case xgb.NoExposureEvent:
            e := NoExposureEvent{&event}
            xu.RunCallbacks(e, NoExposure, e.Drawable)
        case xgb.VisibilityNotifyEvent:
            e := VisibilityNotifyEvent{&event}
            xu.RunCallbacks(e, VisibilityNotify, e.Window)
        case xgb.CreateNotifyEvent:
            e := CreateNotifyEvent{&event}
            xu.RunCallbacks(e, CreateNotify, e.Window)
        case xgb.DestroyNotifyEvent:
            e := DestroyNotifyEvent{&event}
            xu.RunCallbacks(e, DestroyNotify, e.Window)
        case xgb.UnmapNotifyEvent:
            e := UnmapNotifyEvent{&event}
            xu.RunCallbacks(e, UnmapNotify, e.Window)
        case xgb.MapNotifyEvent:
            e := MapNotifyEvent{&event}
            xu.RunCallbacks(e, MapNotify, e.Window)
        case xgb.MapRequestEvent:
            e := MapRequestEvent{&event}
            xu.RunCallbacks(e, MapRequest, e.Window)
        case xgb.ReparentNotifyEvent:
            e := ReparentNotifyEvent{&event}
            xu.RunCallbacks(e, ReparentNotify, e.Window)
        case xgb.ConfigureNotifyEvent:
            e := ConfigureNotifyEvent{&event}
            xu.RunCallbacks(e, ConfigureNotify, e.Window)
        case xgb.ConfigureRequestEvent:
            e := ConfigureRequestEvent{&event}
            xu.RunCallbacks(e, ConfigureRequest, e.Window)
        case xgb.GravityNotifyEvent:
            e := GravityNotifyEvent{&event}
            xu.RunCallbacks(e, GravityNotify, e.Window)
        case xgb.ResizeRequestEvent:
            e := ResizeRequestEvent{&event}
            xu.RunCallbacks(e, ResizeRequest, e.Window)
        case xgb.CirculateNotifyEvent:
            e := CirculateNotifyEvent{&event}
            xu.RunCallbacks(e, CirculateNotify, e.Window)
        case xgb.CirculateRequestEvent:
            e := CirculateRequestEvent{&event}
            xu.RunCallbacks(e, CirculateRequest, e.Window)
        case xgb.PropertyNotifyEvent:
            e := PropertyNotifyEvent{&event}
            xu.RunCallbacks(e, PropertyNotify, e.Window)
        case xgb.SelectionClearEvent:
            e := SelectionClearEvent{&event}
            xu.RunCallbacks(e, SelectionClear, e.Owner)
        case xgb.SelectionRequestEvent:
            e := SelectionRequestEvent{&event}
            xu.RunCallbacks(e, SelectionRequest, e.Requestor)
        case xgb.SelectionNotifyEvent:
            e := SelectionNotifyEvent{&event}
            xu.RunCallbacks(e, SelectionNotify, e.Requestor)
        case xgb.ColormapNotifyEvent:
            e := ColormapNotifyEvent{&event}
            xu.RunCallbacks(e, ColormapNotify, e.Window)
        case xgb.ClientMessageEvent:
            e := ClientMessageEvent{&event}
            xu.RunCallbacks(e, ClientMessage, e.Window)
        case xgb.MappingNotifyEvent:
            e := MappingNotifyEvent{&event}
            xu.RunCallbacks(e, MappingNotify, xgbutil.NoWindow)
        default:
            log.Printf("ERROR: UNSUPPORTED EVENT TYPE: %T\n", event)
            continue
        }
    }

    return nil
}

// SendRootEvent takes a type implementing the XEvent interface, converts it
// to raw X bytes, and sends it off using the SendEvent request.
func SendRootEvent(xu *xgbutil.XUtil, ev XEvent, evMask uint32) {
    xu.Conn().SendEvent(false, xu.RootWin(), evMask, ev.Bytes())
}

// ReplayPointer is a quick alias to AllowEvents with 'ReplayPointer' mode.
func ReplayPointer(xu *xgbutil.XUtil) {
    xu.Conn().AllowEvents(xgb.AllowReplayPointer, 0)
}

