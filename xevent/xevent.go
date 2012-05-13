package xevent

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
)

// Sometimes we need to specify NO WINDOW when a window is typically
// expected. (Like connecting to MappingNotify or KeymapNotify events.)
// Use this value to do that.
var NoWindow xproto.Window = 0

// IgnoreMods is a list of X modifiers that we don't want interfering
// with our mouse or key bindings. In particular, for each mouse or key binding 
// issued, there is a seperate mouse or key binding made for each of the 
// modifiers specified.
var IgnoreMods []uint16 = []uint16{
	0,
	xproto.ModMaskLock,                   // Num lock
	xproto.ModMask2,                      // Caps lock
	xproto.ModMaskLock | xproto.ModMask2, // Caps and Num lock
}

// Enqueue queues up an event read from X.
// Note that an event read may return an error, in which case, this queue
// entry will be an error and not an event.
//
//	ev, err := XUtilValue.Conn().WaitForEvent()
//	xevent.Enqueue(XUtilValue, ev, err)
//
// You probably shouldn't have to enqueue events yourself. This is done
// automatically if you're using xevent.Main{Ping} and/or xevent.Read.
func Enqueue(xu *xgbutil.XUtil, ev xgb.Event, err xgb.Error) {
	xu.EvqueueLck.Lock()
	defer xu.EvqueueLck.Unlock()

	xu.Evqueue = append(xu.Evqueue, xgbutil.EventOrError{
		Event: ev,
		Err:   err,
	})
}

// Dequeue pops an event/error from the queue and returns it.
// The queue item is unwrapped and returned as multiple return values.
// Only one of the return values can be nil.
func Dequeue(xu *xgbutil.XUtil) (xgb.Event, xgb.Error) {
	xu.EvqueueLck.Lock()
	defer xu.EvqueueLck.Unlock()

	everr := xu.Evqueue[0]
	xu.Evqueue = xu.Evqueue[1:]
	return everr.Event, everr.Err
}

// DequeueAt removes a particular item from the queue.
// This is primarily useful when attempting to compress events.
func DequeueAt(xu *xgbutil.XUtil, i int) {
	xu.EvqueueLck.Lock()
	defer xu.EvqueueLck.Unlock()

	xu.Evqueue = append(xu.Evqueue[:i], xu.Evqueue[i+1:]...)
}

// QueueEmpty returns whether the event queue is empty or not.
func Empty(xu *xgbutil.XUtil) bool {
	xu.EvqueueLck.RLock()
	defer xu.EvqueueLck.RUnlock()

	return len(xu.Evqueue) == 0
}

// QueuePeek returns a *copy* of the current queue so we can examine it.
// This can be useful when trying to determine if a particular kind of
// event will be processed in the future.
func Peek(xu *xgbutil.XUtil) []xgbutil.EventOrError {
	xu.EvqueueLck.RLock()
	defer xu.EvqueueLck.RUnlock()

	cpy := make([]xgbutil.EventOrError, len(xu.Evqueue))
	for i, everr := range xu.Evqueue {
		cpy[i] = everr
	}
	return cpy
}

// ErrorHandlerSet sets the default error handler for errors that come
// into the main event loop. (This may be removed in the future in favor
// of a particular callback interface like events, but these sorts of errors
// aren't handled often in practice, so maybe not.)
// This is only called for errors returned from unchecked (asynchronous error
// handling) requests.
// The default error handler just emits them to stderr.
func ErrorHandlerSet(xu *xgbutil.XUtil, fun xgbutil.ErrorHandlerFun) {
	xu.ErrorHandler = fun
}

// ErrorHandlerGet retrieves the default error handler.
func ErrorHandlerGet(xu *xgbutil.XUtil) xgbutil.ErrorHandlerFun {
	return xu.ErrorHandler
}

// RedirectKeyEvents, when set to a window id (greater than 0), will force
// *all* Key{Press,Release} to callbacks attached to the specified window.
// This is close to emulating a Keyboard grab without the racing.
// To stop redirecting key events, use window identifier '0'.
func RedirectKeyEvents(xu *xgbutil.XUtil, wid xproto.Window) {
	xu.KeyRedirect = wid
}

// RedirectKeyGet gets the window that key events are being redirected to.
// If 0, then no redirection occurs.
func RedirectKeyGet(xu *xgbutil.XUtil) xproto.Window {
	return xu.KeyRedirect
}

// Quit elegantly exits out of the main event loop.
// "Elegantly" in this case means that it finishes processing the current
// event, and breaks out of the loop afterwards.
// There is no particular reason to use this instead of something like os.Exit
// other than you might have code to run after the main event loop exits to
// "clean up."
func Quit(xu *xgbutil.XUtil) {
	xu.Quit = true
}

// Quitting returns whether it's time to quit.
// This is only used in the main event loop in xevent.
func Quitting(xu *xgbutil.XUtil) bool {
	return xu.Quit
}

// attachCallback associates a (event, window) tuple with an event.
func attachCallback(xu *xgbutil.XUtil, evtype int, win xproto.Window,
	fun xgbutil.Callback) {

	xu.CallbacksLck.Lock()
	defer xu.CallbacksLck.Unlock()

	if _, ok := xu.Callbacks[evtype]; !ok {
		xu.Callbacks[evtype] = make(map[xproto.Window][]xgbutil.Callback, 20)
	}
	if _, ok := xu.Callbacks[evtype][win]; !ok {
		xu.Callbacks[evtype][win] = make([]xgbutil.Callback, 0)
	}
	xu.Callbacks[evtype][win] = append(xu.Callbacks[evtype][win], fun)
}

// RunCallbacks executes every callback corresponding to a
// particular event/window tuple.
func runCallbacks(xu *xgbutil.XUtil, event interface{}, evtype int,
	win xproto.Window) {

	xu.CallbacksLck.RLock()
	defer xu.CallbacksLck.RUnlock()

	for _, cb := range xu.Callbacks[evtype][win] {
		cb.Run(xu, event)
	}
}

// Detach removes all callbacks associated with a particular window.
// Note that if you're also using the keybind and mousebind packages, a complete
// detachment should look like:
//
//	keybind.Detach(XUtilValue, window-id)
//	mousebind.Detach(XUtilValue, window-id)
//	xevent.Detach(XUtilValue, window-id)
//
// If a window is no longer receiving events, these methods should be called.
// Otherwise, the memory used to store the handler info for that window will
// never be released.
func Detach(xu *xgbutil.XUtil, win xproto.Window) {
	xu.CallbacksLck.Lock()
	defer xu.CallbacksLck.Unlock()

	for evtype, _ := range xu.Callbacks {
		delete(xu.Callbacks[evtype], win)
	}
}

// SendRootEvent takes a type implementing the xgb.Event interface, converts it
// to raw X bytes, and sends it to the root window using the SendEvent request.
func SendRootEvent(xu *xgbutil.XUtil, ev xgb.Event, evMask uint32) {
	xproto.SendEvent(xu.Conn(), false, xu.RootWin(), evMask, string(ev.Bytes()))
}

// ReplayPointer is a quick alias to AllowEvents with 'ReplayPointer' mode.
func ReplayPointer(xu *xgbutil.XUtil) {
	xproto.AllowEvents(xu.Conn(), xproto.AllowReplayPointer, 0)
}
