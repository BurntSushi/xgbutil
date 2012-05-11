/*
   Provides a callback interface, very similar to that found in
   xevent/callback.go --- but only for key bindings.
*/
package keybind

import (
	"fmt"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
)

type KeyPressFun xevent.KeyPressFun

// connect is essentially 'Connect' for either KeyPress or KeyRelease events.
func connect(xu *xgbutil.XUtil, callback xgbutil.KeyBindCallback,
	evtype int, win xproto.Window, keyStr string) error {

	// Get the mods/key first
	mods, keycode, err := ParseString(xu, keyStr)
	if err != nil {
		return err
	}

	// Only do the grab if we haven't yet on this window.
	if xu.KeyBindGrabs(evtype, win, mods, keycode) == 0 {
		err := GrabChecked(xu, win, mods, keycode)
		if err != nil {
			// If a bad access, let's be nice and give a good error message.
			switch err.(type) {
			case xproto.AccessError:
				return fmt.Errorf("Got a bad access error when trying to bind "+
					"'%s'. This usually means another client has already "+
					"grabbed this keybinding.", keyStr)
			default:
				return fmt.Errorf("Could not bind '%s' because: %s",
					keyStr, err)
			}
		}
	}

	// If we've never grabbed anything on this window before, we need to
	// make sure we can respond to it in the main event loop.
	var allCb xgbutil.Callback
	if evtype == xevent.KeyPress {
		allCb = xevent.KeyPressFun(RunKeyPressCallbacks)
	} else {
		allCb = xevent.KeyReleaseFun(RunKeyReleaseCallbacks)
	}

	// If this is the first Key{Press|Release}Event on this window,
	// then we need to listen to Key{Press|Release} events in the main loop.
	if !xu.ConnectedKeyBind(evtype, win) {
		allCb.Connect(xu, win)
	}

	// Finally, attach the callback.
	xu.AttachKeyBindCallback(evtype, win, mods, keycode, callback)

	return nil
}

func DeduceKeyInfo(state uint16,
	detail xproto.Keycode) (uint16, xproto.Keycode) {

	mods, kc := state, detail
	for _, m := range xgbutil.IgnoreMods {
		mods &= ^m
	}
	return mods, kc
}

func (callback KeyPressFun) Connect(xu *xgbutil.XUtil, win xproto.Window,
	keyStr string) error {

	return connect(xu, callback, xevent.KeyPress, win, keyStr)
}

func (callback KeyPressFun) Run(xu *xgbutil.XUtil, event interface{}) {
	callback(xu, event.(xevent.KeyPressEvent))
}

type KeyReleaseFun xevent.KeyReleaseFun

func (callback KeyReleaseFun) Connect(xu *xgbutil.XUtil, win xproto.Window,
	keyStr string) error {

	return connect(xu, callback, xevent.KeyRelease, win, keyStr)
}

func (callback KeyReleaseFun) Run(xu *xgbutil.XUtil, event interface{}) {
	callback(xu, event.(xevent.KeyReleaseEvent))
}

// RunKeyPressCallbacks infers the window, keycode and modifiers from a
// KeyPressEvent and runs the corresponding callbacks.
func RunKeyPressCallbacks(xu *xgbutil.XUtil, ev xevent.KeyPressEvent) {
	mods, kc := DeduceKeyInfo(ev.State, ev.Detail)

	xu.RunKeyBindCallbacks(ev, xevent.KeyPress, ev.Event, mods, kc)
}

// RunKeyReleaseCallbacks infers the window, keycode and modifiers from a
// KeyPressEvent and runs the corresponding callbacks.
func RunKeyReleaseCallbacks(xu *xgbutil.XUtil, ev xevent.KeyReleaseEvent) {
	mods, kc := DeduceKeyInfo(ev.State, ev.Detail)

	xu.RunKeyBindCallbacks(ev, xevent.KeyRelease, ev.Event, mods, kc)
}
