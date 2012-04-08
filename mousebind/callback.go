/*
   Provides a callback interface, very similar to that found in
   xevent/callback.go --- but only for mouse bindings.
*/
package mousebind

import "code.google.com/p/jamslam-x-go-binding/xgb"

import (
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
)

// connect is essentially 'Connect' for either ButtonPress or
// ButtonRelease events.
func connect(xu *xgbutil.XUtil, callback xgbutil.MouseBindCallback,
	evtype int, win xgb.Id, buttonStr string, propagate bool, grab bool) {

	// Get the mods/button first
	mods, button := ParseString(xu, buttonStr)

	// Only do the grab if we haven't yet on this window.
	// And if we WANT a grab...
	if grab && xu.MouseBindGrabs(evtype, win, mods, button) == 0 {
		Grab(xu, win, mods, button, propagate)
	}

	// If we've never grabbed anything on this window before, we need to
	// make sure we can respond to it in the main event loop.
	var allCb xgbutil.Callback
	if evtype == xevent.ButtonPress {
		allCb = xevent.ButtonPressFun(RunButtonPressCallbacks)
	} else {
		allCb = xevent.ButtonReleaseFun(RunButtonReleaseCallbacks)
	}

	// If this is the first Button{Press|Release}Event on this window,
	// then we need to listen to Button{Press|Release} events in the main loop.
	if !xu.ConnectedMouseBind(evtype, win) {
		allCb.Connect(xu, win)
	}

	// Finally, attach the callback.
	xu.AttachMouseBindCallback(evtype, win, mods, button, callback)
}

func deduceButtonInfo(state uint16, detail byte) (uint16, byte) {
	mods, button := state, detail
	for _, m := range xgbutil.IgnoreMods {
		mods &= ^m
	}

	// We also need to mask out the button that has been pressed/released,
	// since it is also typically a modifier.
	modsTemp := int32(mods)
	switch button {
	case 1:
		modsTemp &= ^xgb.ButtonMask1
	case 2:
		modsTemp &= ^xgb.ButtonMask2
	case 3:
		modsTemp &= ^xgb.ButtonMask3
	case 4:
		modsTemp &= ^xgb.ButtonMask4
	case 5:
		modsTemp &= ^xgb.ButtonMask5
	}

	return uint16(modsTemp), button
}

type ButtonPressFun xevent.ButtonPressFun

func (callback ButtonPressFun) Connect(xu *xgbutil.XUtil, win xgb.Id,
	buttonStr string, propagate bool, grab bool) {

	connect(xu, callback, xevent.ButtonPress, win, buttonStr, propagate, grab)
}

func (callback ButtonPressFun) Run(xu *xgbutil.XUtil, event interface{}) {
	callback(xu, event.(xevent.ButtonPressEvent))
}

type ButtonReleaseFun xevent.ButtonReleaseFun

func (callback ButtonReleaseFun) Connect(xu *xgbutil.XUtil, win xgb.Id,
	buttonStr string, propagate bool, grab bool) {

	connect(xu, callback, xevent.ButtonRelease, win, buttonStr, propagate, grab)
}

func (callback ButtonReleaseFun) Run(xu *xgbutil.XUtil, event interface{}) {
	callback(xu, event.(xevent.ButtonReleaseEvent))
}

// RunButtonPressCallbacks infers the window, button and modifiers from a
// ButtonPressEvent and runs the corresponding callbacks.
func RunButtonPressCallbacks(xu *xgbutil.XUtil, ev xevent.ButtonPressEvent) {
	mods, button := deduceButtonInfo(ev.State, ev.Detail)

	xu.RunMouseBindCallbacks(ev, xevent.ButtonPress, ev.Event, mods, button)
}

// RunButtonReleaseCallbacks infers the window, keycode and modifiers from a
// ButtonPressEvent and runs the corresponding callbacks.
func RunButtonReleaseCallbacks(xu *xgbutil.XUtil,
	ev xevent.ButtonReleaseEvent) {

	mods, button := deduceButtonInfo(ev.State, ev.Detail)

	xu.RunMouseBindCallbacks(ev, xevent.ButtonRelease, ev.Event, mods, button)
}
