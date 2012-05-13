package xevent

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"

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
// N.B. If you have multiple X connections in the same program, you should be
// able to run this in different goroutines concurrently. However, only
// *one* of these should run for *each* connection.
func Main(xu *xgbutil.XUtil) {
	mainEventLoop(xu, nil)
}

// MainPing starts the main X event loop, and returns a ping channel.
// A benign value will be sent to the ping channel every time an event/error is
// dequeued.
// This is useful if your event loop needs to draw from other sources. e.g.,
//
//	ping := xevent.MainPing()
//	for {
//		select {
//		case <-ping:
//		case val <- someOtherChannel:
//			// do some work with val
//		}
//	}
//
// Note that an unbuffered channel is returned, which implies that any work
// done in 'val' will delay further X event processing.
// N.B. If you have multiple X connections in the same program, you should be
// able to run this in different goroutines concurrently. However, only
// *one* of these should run for *each* connection.
func MainPing(xu *xgbutil.XUtil) chan struct{} {
	ping := make(chan struct{}, 0)
	go func() {
		mainEventLoop(xu, ping)
	}()
	return ping
}

// mainEventLoop runs the main event loop with an optional ping channel.
func mainEventLoop(xu *xgbutil.XUtil, ping chan struct{}) {
	for {
		if xu.Quitting() {
			break
		}

		// Gobble up as many events as possible (into the queue).
		// If there are no events, we block.
		Read(xu, true)

		// Now process every event/error in the queue.
		processEventQueue(xu, ping)
	}
}

// processEventQueue processes every item in the event/error queue.
func processEventQueue(xu *xgbutil.XUtil, ping chan struct{}) {
	for !xu.QueueEmpty() {
		if xu.Quitting() {
			return
		}

		// We technically send the ping *before* the next event is dequeued.
		// This is so the queue doesn't present a misrepresentation of which
		// events haven't been processed yet.
		if ping != nil {
			ping <- struct{}{}
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
		case xproto.KeyPressEvent:
			e := KeyPressEvent{&event}

			// If we're redirecting key events, this is the place to do it!
			if wid := xu.RedirectKeyGet(); wid > 0 {
				e.Event = wid
			}

			xu.TimeSet(e.Time)
			xu.RunCallbacks(e, KeyPress, e.Event)
		case xproto.KeyReleaseEvent:
			e := KeyReleaseEvent{&event}

			// If we're redirecting key events, this is the place to do it!
			if wid := xu.RedirectKeyGet(); wid > 0 {
				e.Event = wid
			}

			xu.TimeSet(e.Time)
			xu.RunCallbacks(e, KeyRelease, e.Event)
		case xproto.ButtonPressEvent:
			e := ButtonPressEvent{&event}
			xu.TimeSet(e.Time)
			xu.RunCallbacks(e, ButtonPress, e.Event)
		case xproto.ButtonReleaseEvent:
			e := ButtonReleaseEvent{&event}
			xu.TimeSet(e.Time)
			xu.RunCallbacks(e, ButtonRelease, e.Event)
		case xproto.MotionNotifyEvent:
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
			var laste xproto.MotionNotifyEvent
			for {
				xu.Sync()
				Read(xu, false)

				found := false
				for i, ee := range xu.QueuePeek() {
					if ee.Err != nil {
						continue
					}
					if mn, ok := ee.Event.(xproto.MotionNotifyEvent); ok {
						if mn.Event == e.Event {
							laste = mn
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
		case xproto.EnterNotifyEvent:
			e := EnterNotifyEvent{&event}
			xu.TimeSet(e.Time)
			xu.RunCallbacks(e, EnterNotify, e.Event)
		case xproto.LeaveNotifyEvent:
			e := LeaveNotifyEvent{&event}
			xu.TimeSet(e.Time)
			xu.RunCallbacks(e, LeaveNotify, e.Event)
		case xproto.FocusInEvent:
			e := FocusInEvent{&event}
			xu.RunCallbacks(e, FocusIn, e.Event)
		case xproto.FocusOutEvent:
			e := FocusOutEvent{&event}
			xu.RunCallbacks(e, FocusOut, e.Event)
		case xproto.KeymapNotifyEvent:
			e := KeymapNotifyEvent{&event}
			xu.RunCallbacks(e, KeymapNotify, xgbutil.NoWindow)
		case xproto.ExposeEvent:
			e := ExposeEvent{&event}
			xu.RunCallbacks(e, Expose, e.Window)
		case xproto.GraphicsExposureEvent:
			e := GraphicsExposureEvent{&event}
			xu.RunCallbacks(e, GraphicsExposure, xproto.Window(e.Drawable))
		case xproto.NoExposureEvent:
			e := NoExposureEvent{&event}
			xu.RunCallbacks(e, NoExposure, xproto.Window(e.Drawable))
		case xproto.VisibilityNotifyEvent:
			e := VisibilityNotifyEvent{&event}
			xu.RunCallbacks(e, VisibilityNotify, e.Window)
		case xproto.CreateNotifyEvent:
			e := CreateNotifyEvent{&event}
			xu.RunCallbacks(e, CreateNotify, e.Window)
		case xproto.DestroyNotifyEvent:
			e := DestroyNotifyEvent{&event}
			xu.RunCallbacks(e, DestroyNotify, e.Window)
		case xproto.UnmapNotifyEvent:
			e := UnmapNotifyEvent{&event}
			xu.RunCallbacks(e, UnmapNotify, e.Window)
		case xproto.MapNotifyEvent:
			e := MapNotifyEvent{&event}
			xu.RunCallbacks(e, MapNotify, e.Window)
		case xproto.MapRequestEvent:
			e := MapRequestEvent{&event}
			xu.RunCallbacks(e, MapRequest, e.Window)
			xu.RunCallbacks(e, MapRequest, e.Parent)
		case xproto.ReparentNotifyEvent:
			e := ReparentNotifyEvent{&event}
			xu.RunCallbacks(e, ReparentNotify, e.Window)
		case xproto.ConfigureNotifyEvent:
			e := ConfigureNotifyEvent{&event}
			xu.RunCallbacks(e, ConfigureNotify, e.Window)
		case xproto.ConfigureRequestEvent:
			e := ConfigureRequestEvent{&event}
			xu.RunCallbacks(e, ConfigureRequest, e.Window)
			xu.RunCallbacks(e, ConfigureRequest, e.Parent)
		case xproto.GravityNotifyEvent:
			e := GravityNotifyEvent{&event}
			xu.RunCallbacks(e, GravityNotify, e.Window)
		case xproto.ResizeRequestEvent:
			e := ResizeRequestEvent{&event}
			xu.RunCallbacks(e, ResizeRequest, e.Window)
		case xproto.CirculateNotifyEvent:
			e := CirculateNotifyEvent{&event}
			xu.RunCallbacks(e, CirculateNotify, e.Window)
		case xproto.CirculateRequestEvent:
			e := CirculateRequestEvent{&event}
			xu.RunCallbacks(e, CirculateRequest, e.Window)
		case xproto.PropertyNotifyEvent:
			e := PropertyNotifyEvent{&event}
			xu.TimeSet(e.Time)
			xu.RunCallbacks(e, PropertyNotify, e.Window)
		case xproto.SelectionClearEvent:
			e := SelectionClearEvent{&event}
			xu.TimeSet(e.Time)
			xu.RunCallbacks(e, SelectionClear, e.Owner)
		case xproto.SelectionRequestEvent:
			e := SelectionRequestEvent{&event}
			xu.TimeSet(e.Time)
			xu.RunCallbacks(e, SelectionRequest, e.Requestor)
		case xproto.SelectionNotifyEvent:
			e := SelectionNotifyEvent{&event}
			xu.TimeSet(e.Time)
			xu.RunCallbacks(e, SelectionNotify, e.Requestor)
		case xproto.ColormapNotifyEvent:
			e := ColormapNotifyEvent{&event}
			xu.RunCallbacks(e, ColormapNotify, e.Window)
		case xproto.ClientMessageEvent:
			e := ClientMessageEvent{&event}
			xu.RunCallbacks(e, ClientMessage, e.Window)
		case xproto.MappingNotifyEvent:
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

// SendRootEvent takes a type implementing the xgb.Event interface, converts it
// to raw X bytes, and sends it off using the SendEvent request.
func SendRootEvent(xu *xgbutil.XUtil, ev xgb.Event, evMask uint32) {
	xproto.SendEvent(xu.Conn(), false, xu.RootWin(), evMask, string(ev.Bytes()))
}

// ReplayPointer is a quick alias to AllowEvents with 'ReplayPointer' mode.
func ReplayPointer(xu *xgbutil.XUtil) {
	xproto.AllowEvents(xu.Conn(), xproto.AllowReplayPointer, 0)
}
