package xgbutil

/*
types.go contains several types used in the XUtil structure. In an ideal world,
they would be defined in their appropriate packages, but must be defined here
(and exported) for use in some sub-packages. (Namely, xevent, keybind and
mousebind.)
*/

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

// Callback is an interface that should be implemented by event callback 
// functions. Namely, to assign a function to a particular event/window
// combination, simply define a function with type 'SomeEventFun' (pre-defined
// in xevent/callback.go), and call the 'Connect' method.
// The 'Run' method is used inside the Main event loop, and shouldn't be used
// by the user.
// Also, it is perfectly legitimate to connect to events that don't specify
// a window (like MappingNotify and KeymapNotify). In this case, simply
// use 'xgbutil.NoWindow' as the window id.
//
// Example to respond to ConfigureNotify events on window 0x1
//
//     xevent.ConfigureNotifyFun(
//		func(X *xgbutil.XUtil, e xevent.ConfigureNotifyEvent) {
//			fmt.Printf("(%d, %d) %dx%d\n", e.X, e.Y, e.Width, e.Height)
//		}).Connect(X, 0x1)
type Callback interface {
	Connect(xu *XUtil, win xproto.Window)
	Run(xu *XUtil, ev interface{})
}

// KeyBindCallback works similarly to the more general Callback, but it adds
// parameters specific to key bindings.
type KeyBindCallback interface {
	// Connect modifies XUtil's state to attach an event handler to a
	// particular key press. If grab is true, connect will request a passive
	// grab.
	Connect(xu *XUtil, win xproto.Window, keyStr string, grab bool) error

	// Run is exported for use in the keybind package but should not be
	// used by the user. (It is used to run the callback function in the
	// main event loop.)
	Run(xu *XUtil, ev interface{})
}

// MouseBindCallback works similarly to the more general Callback, but it adds
// parameters specific to mouse bindings.
type MouseBindCallback interface {
	Connect(xu *XUtil, win xproto.Window, buttonStr string,
		propagate bool, grab bool) error
	Run(xu *XUtil, ev interface{})
}

// KeyBindKey is the type of the key in the map of keybindings.
// It essentially represents the tuple
// (event type, window id, modifier, keycode).
// It is exported for use in the keybind package. It should not be used.
type KeyBindKey struct {
	Evtype int
	Win    xproto.Window
	Mod    uint16
	Code   xproto.Keycode
}

// MouseBindKey is the type of the key in the map of mouse bindings.
// It essentially represents the tuple
// (event type, window id, modifier, button).
type MouseBindKey struct {
	Evtype int
	Win    xproto.Window
	Mod    uint16
	Button xproto.Button
}

// KeyboardMapping embeds a keyboard mapping reply from XGB.
// It should be retrieved using keybind.KeyMapGet, if necessary.
// xgbutil tries quite hard to absolve you from ever having to use this.
// A keyboard mapping is a table that maps keycodes to one or more keysyms.
type KeyboardMapping struct {
	*xproto.GetKeyboardMappingReply
}

// ModifierMapping embeds a modifier mapping reply from XGB.
// It should be retrieved using keybind.ModMapGet, if necessary.
// xgbutil tries quite hard to absolve you from ever having to use this.
// A modifier mapping is a table that maps modifiers to one or more keycodes.
type ModifierMapping struct {
	*xproto.GetModifierMappingReply
}

// ErrorHandlerFun is the type of function required to handle errors that
// come in through the main event loop.
// For example, to set a new error handler, use:
//
//	xevent.ErrorHandlerSet(xgbutil.ErrorHandlerFun(
//		func(err xgb.Error) {
//			// do something with err
//		}))
type ErrorHandlerFun func(err xgb.Error)

// EventOrError is a struct that contains either an event value or an error
// value. It is an error to contain both. Containing neither indicates an
// error too.
// This is exported for use in the xevent package. You shouldn't have any
// direct contact with values of this type, unless you need to inspect the
// queue directly with xevent.Peek.
type EventOrError struct {
	Event xgb.Event
	Err   xgb.Error
}
