/*
Example compress-events shows how to manipulate the xevent package's event
queue to compress events that arrive more often than you'd like to process
them. This example in particular shows how to compress MotionNotify events,
but the same approach could be used to compress ConfigureNotify events.

Note that we show the difference between compressed and uncompressed
MotionNotify events by displaying two windows that listen for MotionNotify
events. The green window compresses them while the red window does not.
Hovering over each window will print the x and y positions in each
MotionNotify event received. You should notice that the red window
lags behind the pointer (particularly if you moved the pointer quickly in
and out of the window) while the green window always keeps up, regardless
of the speed of the pointer.

In each case, we simulate work by sleeping for some amount of time. (The
whole point of compressing events is that there is too much work to be done
for each event.)

Note that when compressing events, you should always make sure that the
event you're compressing *ought* to be compressed. For example, with
MotionNotify events, if the Event field changes, then it applies to a
different window and probably shouldn't be compressed with MotionNotify
events for other windows.

Finally, compressing events implicitly assumes that the event handler doing
the compression is the *only* event handler for a particular (event, window)
tuple. If there is more than one event handler for a single (event, window)
tuple and one of them does compression, the other will be left out in the
cold. (Since the main event loop is subverted and won't process the
compressed events in the usual way.)

N.B. This functionality isn't included in xgbutil because event compression
isn't something that is always desirable, and the conditions under which
compression happens can vary. In particular, compressing ConfigureRequest
events from the perspective of the window manager can be faulty, since
changes to other properties (like WM_NORMAL_HINTS) can change the semantics
of a ConfigureRequest event. (i.e., your compression would need to
specifically look for events that could change future ConfigureRequest
events.)
*/
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
)

// workTime is the amount of time to sleep to simulate "work" in response to
// MotionNotify events. Increasing this will exacerbate the difference
// between the green and red windows. But if you increase it too much,
// the red window starts to *really* lag, and you'll probably have to kill
// the program.
var workTime = 50 * time.Millisecond

// newWindow creates a new window that listens to MotionNotify events with
// the given backgroundcolor.
func newWindow(X *xgbutil.XUtil, color uint32) *xwindow.Window {
	win, err := xwindow.Generate(X)
	if err != nil {
		log.Fatal(err)
	}

	err = win.CreateChecked(X.RootWin(), 0, 0, 400, 400,
		xproto.CwBackPixel|xproto.CwEventMask,
		color, xproto.EventMaskPointerMotion)
	if err != nil {
		log.Fatal(err)
	}

	win.Map()
	return win
}

// compressMotionNotify takes a MotionNotify event, and inspects the event
// queue for any future MotionNotify events that can be received without
// blocking. The most recent MotionNotify event is then returned.
// Note that we need to make sure that the Event, Child, Detail, State, Root
// and SameScreen fields are the same to ensure the same window/action is
// generating events. That is, we are only compressing the RootX, RootY,
// EventX and EventY fields.
// This function is not thread safe, since Peek returns a *copy* of the
// event queue---which could be out of date by the time we dequeue events.
func compressMotionNotify(X *xgbutil.XUtil,
	ev xevent.MotionNotifyEvent) xevent.MotionNotifyEvent {

	// We force a round trip request so that we make sure to read all
	// available events.
	X.Sync()
	xevent.Read(X, false)

	// The most recent MotionNotify event that we'll end up returning.
	laste := ev

	// Look through each event in the queue. If it's an event and it matches
	// all the fields in 'ev' that are detailed above, then set it to 'laste'.
	// In which case, we'll also dequeue the event, otherwise it will be
	// processed twice!
	// N.B. If our only goal was to find the most recent relevant MotionNotify
	// event, we could traverse the event queue backwards and simply use
	// the first MotionNotify we see. However, this could potentially leave
	// other MotionNotify events in the queue, which we *don't* want to be
	// processed. So we stride along and just pick off MotionNotify events
	// until we don't see any more.
	for i, ee := range xevent.Peek(X) {
		if ee.Err != nil { // This is an error, skip it.
			continue
		}

		// Use type assertion to make sure this is a MotionNotify event.
		if mn, ok := ee.Event.(xproto.MotionNotifyEvent); ok {
			// Now make sure all appropriate fields are equivalent.
			if ev.Event == mn.Event && ev.Child == mn.Child &&
				ev.Detail == mn.Detail && ev.State == mn.State &&
				ev.Root == mn.Root && ev.SameScreen == mn.SameScreen {

				// Set the most recent/valid motion notify event.
				laste = xevent.MotionNotifyEvent{&mn}

				// We cheat and use the stack semantics of defer to dequeue
				// most recent motion notify events first, so that the indices
				// don't become invalid. (If we dequeued oldest first, we'd
				// have to account for all future events shifting to the left
				// by one.)
				defer func(i int) { xevent.DequeueAt(X, i) }(i)
			}
		}
	}

	// This isn't strictly necessary, but is correct. We should update
	// xgbutil's sense of time with the most recent event processed.
	// This is typically done in the main event loop, but since we are
	// subverting the main event loop, we should take care of it.
	X.TimeSet(laste.Time)

	return laste
}

func main() {
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	// Create window for receiving compressed MotionNotify events.
	cwin := newWindow(X, 0x00ff00)

	// Attach event handler for MotionNotify that compresses events.
	xevent.MotionNotifyFun(
		func(X *xgbutil.XUtil, ev xevent.MotionNotifyEvent) {
			ev = compressMotionNotify(X, ev)
			fmt.Printf("COMPRESSED: (EventX %d, EventY %d)\n",
				ev.EventX, ev.EventY)
			time.Sleep(workTime)
		}).Connect(X, cwin.Id)

	// Create window for receiving uncompressed MotionNotify events.
	uwin := newWindow(X, 0xff0000)

	// Attach event handler for MotionNotify that does not compress events.
	xevent.MotionNotifyFun(
		func(X *xgbutil.XUtil, ev xevent.MotionNotifyEvent) {
			fmt.Printf("UNCOMPRESSED: (EventX %d, EventY %d)\n",
				ev.EventX, ev.EventY)
			time.Sleep(workTime)
		}).Connect(X, uwin.Id)

	xevent.Main(X)
}
