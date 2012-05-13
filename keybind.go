package xgbutil

/*
keybind.go contains several types used in the keybind package.

They are defined here because values of these types must be stored in an
XUtil value.
*/

import "github.com/BurntSushi/xgb/xproto"

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
