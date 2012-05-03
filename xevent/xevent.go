/*
   event.go organizes functions related to X event processing.

   (Note: The 'evtypes.go' file is related. It defines methods/interfaces
          for dealing with common X events. event.go on the other hand,
          is more concerned with responding to events that X sends us.)
*/
package xevent

import (
	"io"
	"log"

	"github.com/BurntSushi/xgb"

	"github.com/BurntSushi/xgbutil"
)

// Read reads one or more events and queues them in XUtil.
// If 'block' is True, then call 'WaitForEvent' before sucking up
// all events that have been queued by XGB.
// Returns false if we've hit an unrecoverable error and need to quit.
func Read(xu *xgbutil.XUtil, block bool) bool {
	if block {
		ev, err := xu.Conn().WaitForEvent()
		if err != nil {
			if processEventError(xu, err) {
				return false
			}
		}
		xu.Enqueue(ev)
	}

	// Clean up anything that's in the queue
	for {
		ev, err := xu.Conn().PollForEvent()

		// No events left...
		if ev == nil && err == nil {
			break
		}

		// Error!
		if err != nil {
			if processEventError(xu, err) {
				return false
			} else {
				continue
			}
		}

		// We're good, queue it up
		xu.Enqueue(ev)
	}

	return true
}

// processEventError handles errors returned when calling
// WaitForEvent or PollForEvent.
// In particular, if the error is extreme (like EOF), then we unfortunately
// need to crash. If it's something like a BadValue or a BadWindow, we can
// live for another day...
// The 'bool' return value is true when the caller needs to QUIT.
func processEventError(xu *xgbutil.XUtil, err error) bool {
	if err == io.EOF {
		log.Println("EOF. Stopping everything. Sorry :-(")
		return true
	} else if xgbErr, ok := err.(xgb.Error); ok {
		if !xu.IgnoredWindow(xgbErr.BadId()) {
			log.Printf("ERROR: %v\n", err)
		}
		return false
	}

	log.Printf("UNKNOWN ERROR: %v\n", err)
	return true
}

// Main starts the main X event loop. It will read events and call appropriate
// callback functions. 
// XXX: This only supports using one X connection. It should at least allow
//      some arbitrary number of connections. Not sure what the best approach
//      is, but I'm sure it needs to use channels (and multiple calls to
//      Main using select, by the user).
func Main(xu *xgbutil.XUtil) error {
	for {
		if xu.Quitting() {
			break
		}

		if !Read(xu, true) {
			break
		}

		// We have to look for xgb events here. But we re-wrap them in our
		// own event types.
		for !xu.QueueEmpty() {
			if xu.Quitting() {
				return nil
			}

			reply := xu.Dequeue()

			switch event := reply.(type) {
			case xgb.KeyPressEvent:
				e := KeyPressEvent{&event}

				// If we're redirecting key events, this is the place to do it!
				if wid := xu.RedirectKeyGet(); wid > 0 {
					e.Event = wid
				}

				xu.SetTime(e.Time)
				xu.RunCallbacks(e, KeyPress, e.Event)
			case xgb.KeyReleaseEvent:
				e := KeyReleaseEvent{&event}

				// If we're redirecting key events, this is the place to do it!
				if wid := xu.RedirectKeyGet(); wid > 0 {
					e.Event = wid
				}

				xu.SetTime(e.Time)
				xu.RunCallbacks(e, KeyRelease, e.Event)
			case xgb.ButtonPressEvent:
				e := ButtonPressEvent{&event}
				xu.SetTime(e.Time)
				xu.RunCallbacks(e, ButtonPress, e.Event)
			case xgb.ButtonReleaseEvent:
				e := ButtonReleaseEvent{&event}
				xu.SetTime(e.Time)
				xu.RunCallbacks(e, ButtonRelease, e.Event)
			case xgb.MotionNotifyEvent:
				e := MotionNotifyEvent{&event}

				// Peek at the next events, if it's just another
				// MotionNotify, let's compress!
				// This is actually pretty nasty. The key here is to flush
				// the buffer so we have an updated list of events.
				// Then we read those events into our queue, but don't block
				// while we do. Finally, we look through the queue and start
				// popping off motion notifies that match 'e'. If we pop one
				// off, restart the process of finding a motion notify.
				// Otherwise, we're done and we move on with the current
				// motion notify.
				var laste xgb.MotionNotifyEvent
				for {
					xu.Flush()
					Read(xu, false)

					found := false
					for i, ev := range xu.QueuePeek() {
						if motNot, ok := ev.(xgb.MotionNotifyEvent); ok {
							if motNot.Event == e.Event {
								laste = motNot
								xu.DequeueAt(i)
								found = true
								break
							}
						}
					}
					if !found {
						break
					}
				}

				if laste.Root != 0 {
					e.Time = laste.Time
					e.RootX = laste.RootX
					e.RootY = laste.RootY
					e.EventX = laste.EventX
					e.EventY = laste.EventY
				}

				xu.SetTime(e.Time)
				xu.RunCallbacks(e, MotionNotify, e.Event)
			case xgb.EnterNotifyEvent:
				e := EnterNotifyEvent{&event}
				xu.SetTime(e.Time)
				xu.RunCallbacks(e, EnterNotify, e.Event)
			case xgb.LeaveNotifyEvent:
				e := LeaveNotifyEvent{&event}
				xu.SetTime(e.Time)
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
				xu.RunCallbacks(e, GraphicsExposure, xgb.Id(e.Drawable))
			case xgb.NoExposureEvent:
				e := NoExposureEvent{&event}
				xu.RunCallbacks(e, NoExposure, xgb.Id(e.Drawable))
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
				xu.RunCallbacks(e, MapRequest, e.Parent)
			case xgb.ReparentNotifyEvent:
				e := ReparentNotifyEvent{&event}
				xu.RunCallbacks(e, ReparentNotify, e.Window)
			case xgb.ConfigureNotifyEvent:
				e := ConfigureNotifyEvent{&event}
				xu.RunCallbacks(e, ConfigureNotify, e.Window)
			case xgb.ConfigureRequestEvent:
				e := ConfigureRequestEvent{&event}
				xu.RunCallbacks(e, ConfigureRequest, e.Window)
				xu.RunCallbacks(e, ConfigureRequest, e.Parent)
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
				xu.SetTime(e.Time)
				xu.RunCallbacks(e, PropertyNotify, e.Window)
			case xgb.SelectionClearEvent:
				e := SelectionClearEvent{&event}
				xu.SetTime(e.Time)
				xu.RunCallbacks(e, SelectionClear, e.Owner)
			case xgb.SelectionRequestEvent:
				e := SelectionRequestEvent{&event}
				xu.SetTime(e.Time)
				xu.RunCallbacks(e, SelectionRequest, e.Requestor)
			case xgb.SelectionNotifyEvent:
				e := SelectionNotifyEvent{&event}
				xu.SetTime(e.Time)
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
				if event != nil {
					log.Printf("ERROR: UNSUPPORTED EVENT TYPE: %T\n", event)
				}
				continue
			}
		}
	}

	return nil
}

// SendRootEvent takes a type implementing the XEvent interface, converts it
// to raw X bytes, and sends it off using the SendEvent request.
func SendRootEvent(xu *xgbutil.XUtil, ev XEvent, evMask uint32) {
	xu.Conn().SendEvent(false, xu.RootWin(), evMask, string(ev.Bytes()))
}

// ReplayPointer is a quick alias to AllowEvents with 'ReplayPointer' mode.
func ReplayPointer(xu *xgbutil.XUtil) {
	xu.Conn().AllowEvents(xgb.AllowReplayPointer, 0)
}

// Detach removes *everything* associated with a particular
// window, including key and mouse bindings.
func Detach(xu *xgbutil.XUtil, win xgb.Id) {
	xu.DetachWindow(win)
	xu.DetachKeyBindWindow(KeyPress, win)
	xu.DetachKeyBindWindow(KeyRelease, win)
	xu.DetachMouseBindWindow(ButtonPress, win)
	xu.DetachMouseBindWindow(ButtonRelease, win)
}
