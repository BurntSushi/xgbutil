/*
   event.go organizes functions related to X event processing.

   (Note: The 'evtypes.go' file is related. It defines methods/interfaces
          for dealing with common X events. event.go on the other hand,
          is more concerned with responding to events that X sends us.)
*/
package xevent

import (
	"github.com/BurntSushi/xgb"

	"github.com/BurntSushi/xgbutil"
)

// Read reads one or more events and queues them in XUtil.
// If 'block' is True, then call 'WaitForEvent' before sucking up
// all events that have been queued by XGB.
func Read(xu *xgbutil.XUtil, block bool) {
	if block {
		ev, err := xu.Conn().WaitForEvent()
		if ev == nil && err == nil {
			xgbutil.Logger.Fatal("BUG: Could not read an event or an error.")
		}
		xu.Enqueue(xgbutil.NewEventOrError(ev, err))
	}

	// Clean up anything that's in the queue
	for {
		ev, err := xu.Conn().PollForEvent()

		// No events left...
		if ev == nil && err == nil {
			break
		}

		// We're good, queue it up
		xu.Enqueue(xgbutil.NewEventOrError(ev, err))
	}
}

// Main starts the main X event loop. It will read events and call appropriate
// callback functions. 
// N.B. If you have multiple X connections in the same program, you could be
// able to run this in different goroutines concurrently. However, only
// *one* of these should run for *each* connection.
func Main(xu *xgbutil.XUtil) error {
	for {
		if xu.Quitting() {
			break
		}

		// Gobble up as many events as possible (into the queue).
		// If there are no events, we block.
		Read(xu, true)

		for !xu.QueueEmpty() {
			if xu.Quitting() {
				return nil
			}

			everr := xu.Dequeue()

			// If we gobbled up an error, send it to the error event handler
			// and move on the next event/error.
			if everr.Err != nil {
				xu.ErrorHandlerGet()(everr.Err)
				continue
			}

			// We know there isn't an error. If there isn't an event either,
			// then there's a bug somewhere.
			if everr.Event == nil {
				xgbutil.Logger.Fatal("BUG: Expected an event but got nil.")
			}

			switch event := everr.Event.(type) {
			case xgb.KeyPressEvent:
				e := KeyPressEvent{&event}

				// If we're redirecting key events, this is the place to do it!
				if wid := xu.RedirectKeyGet(); wid > 0 {
					e.Event = wid
				}

				xu.TimeSet(e.Time)
				xu.RunCallbacks(e, KeyPress, e.Event)
			case xgb.KeyReleaseEvent:
				e := KeyReleaseEvent{&event}

				// If we're redirecting key events, this is the place to do it!
				if wid := xu.RedirectKeyGet(); wid > 0 {
					e.Event = wid
				}

				xu.TimeSet(e.Time)
				xu.RunCallbacks(e, KeyRelease, e.Event)
			case xgb.ButtonPressEvent:
				e := ButtonPressEvent{&event}
				xu.TimeSet(e.Time)
				xu.RunCallbacks(e, ButtonPress, e.Event)
			case xgb.ButtonReleaseEvent:
				e := ButtonReleaseEvent{&event}
				xu.TimeSet(e.Time)
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
					xu.Sync()
					Read(xu, false)

					found := false
					for i, ee := range xu.QueuePeek() {
						if ee.Err != nil {
							continue
						}
						if motNot, ok := ee.Event.(xgb.MotionNotifyEvent); ok {
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

				xu.TimeSet(e.Time)
				xu.RunCallbacks(e, MotionNotify, e.Event)
			case xgb.EnterNotifyEvent:
				e := EnterNotifyEvent{&event}
				xu.TimeSet(e.Time)
				xu.RunCallbacks(e, EnterNotify, e.Event)
			case xgb.LeaveNotifyEvent:
				e := LeaveNotifyEvent{&event}
				xu.TimeSet(e.Time)
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
				xu.TimeSet(e.Time)
				xu.RunCallbacks(e, PropertyNotify, e.Window)
			case xgb.SelectionClearEvent:
				e := SelectionClearEvent{&event}
				xu.TimeSet(e.Time)
				xu.RunCallbacks(e, SelectionClear, e.Owner)
			case xgb.SelectionRequestEvent:
				e := SelectionRequestEvent{&event}
				xu.TimeSet(e.Time)
				xu.RunCallbacks(e, SelectionRequest, e.Requestor)
			case xgb.SelectionNotifyEvent:
				e := SelectionNotifyEvent{&event}
				xu.TimeSet(e.Time)
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
					xgbutil.Logger.Printf("ERROR: UNSUPPORTED EVENT TYPE: %T",
						event)
				}
				continue
			}
		}
	}

	return nil
}

// SendRootEvent takes a type implementing the xgb.Event interface, converts it
// to raw X bytes, and sends it off using the SendEvent request.
func SendRootEvent(xu *xgbutil.XUtil, ev xgb.Event, evMask uint32) {
	xu.Conn().SendEvent(false, xu.RootWin(), evMask, string(ev.Bytes()))
}

// ReplayPointer is a quick alias to AllowEvents with 'ReplayPointer' mode.
func ReplayPointer(xu *xgbutil.XUtil) {
	xu.Conn().AllowEvents(xgb.AllowReplayPointer, 0)
}

// Detach removes *everything* associated with a particular
// window, including key and mouse bindings.
// This should be used on a window that can no longer receive events. (i.e.,
// it was destroyed.)
func Detach(xu *xgbutil.XUtil, win xgb.Id) {
	xu.DetachWindow(win)
	xu.DetachKeyBindWindow(KeyPress, win)
	xu.DetachKeyBindWindow(KeyRelease, win)
	xu.DetachMouseBindWindow(ButtonPress, win)
	xu.DetachMouseBindWindow(ButtonRelease, win)
}
