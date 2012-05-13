package keybind

/*
keybind/xutil.go contains a collection of functions that modify the
Keybinds and Keygrabs state of an XUtil value.

They could have been placed inside the core xgbutil package, but they would
have to be exported for use by the keybind package. In which case, the API
would become cluttered with functions that should not be used.
*/

import (
	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
)

// attachKeyBindCallback associates an (event, window, mods, keycode)
// with a callback.
// This is exported for use in the keybind package. It should not be used.
func attachKeyBindCallback(xu *xgbutil.XUtil, evtype int, win xproto.Window,
	mods uint16, keycode xproto.Keycode, fun xgbutil.KeyBindCallback) {

	xu.KeybindsLck.Lock()
	defer xu.KeybindsLck.Unlock()

	// Create key
	key := xgbutil.KeyBindKey{evtype, win, mods, keycode}

	// Do we need to allocate?
	if _, ok := xu.Keybinds[key]; !ok {
		xu.Keybinds[key] = make([]xgbutil.KeyBindCallback, 0)
	}

	xu.Keybinds[key] = append(xu.Keybinds[key], fun)
	xu.Keygrabs[key] += 1
}

// keyBindKeys returns a copy of all the keys in the 'keybinds' map.
// This is exported for use in the keybind package. It should not be used.
func keyBindKeys(xu *xgbutil.XUtil) []xgbutil.KeyBindKey {
	xu.KeybindsLck.RLock()
	defer xu.KeybindsLck.RUnlock()

	keys := make([]xgbutil.KeyBindKey, len(xu.Keybinds))
	i := 0
	for key, _ := range xu.Keybinds {
		keys[i] = key
		i++
	}
	return keys
}

// updateKeyBindKey takes a key bind key and a new key code.
// It will then remove the old key from keybinds and keygrabs,
// and add the new key with the old key's data into keybinds and keygrabs.
// Its primary purpose is to facilitate the renewal of state when xgbutil
// receives a new keyboard mapping.
// This is exported for use in the keybind package. It should not be used.
func updateKeyBindKey(xu *xgbutil.XUtil, key xgbutil.KeyBindKey,
	newKc xproto.Keycode) {

	xu.KeybindsLck.Lock()
	defer xu.KeybindsLck.Unlock()

	newKey := xgbutil.KeyBindKey{key.Evtype, key.Win, key.Mod, newKc}

	// Save old info
	oldCallbacks := xu.Keybinds[key]
	oldGrabs := xu.Keygrabs[key]

	// Delete old keys
	delete(xu.Keybinds, key)
	delete(xu.Keygrabs, key)

	// Add new keys with old info
	xu.Keybinds[newKey] = oldCallbacks
	xu.Keygrabs[newKey] = oldGrabs
}

// runKeyBindCallbacks executes every callback corresponding to a
// particular event/window/mod/key tuple.
// This is exported for use in the keybind package. It should not be used.
func runKeyBindCallbacks(xu *xgbutil.XUtil, event interface{}, evtype int,
	win xproto.Window, mods uint16, keycode xproto.Keycode) {

	key := xgbutil.KeyBindKey{evtype, win, mods, keycode}
	for _, cb := range keyBindCallbacks(xu, key) {
		cb.Run(xu, event)
	}
}

// keyBindCallbacks returns a slice of callbacks for a particular key.
func keyBindCallbacks(xu *xgbutil.XUtil,
	key xgbutil.KeyBindKey) []xgbutil.KeyBindCallback {

	xu.KeybindsLck.RLock()
	defer xu.KeybindsLck.RUnlock()

	cbs := make([]xgbutil.KeyBindCallback, len(xu.Keybinds[key]))
	for i, cb := range xu.Keybinds[key] {
		cbs[i] = cb
	}
	return cbs
}

// ConnectedKeyBind checks to see if there are any key binds for a particular
// event type already in play.
func connectedKeyBind(xu *xgbutil.XUtil, evtype int, win xproto.Window) bool {
	xu.KeybindsLck.RLock()
	defer xu.KeybindsLck.RUnlock()

	// Since we can't create a full key, loop through all key binds
	// and check if evtype and window match.
	for key, _ := range xu.Keybinds {
		if key.Evtype == evtype && key.Win == win {
			return true
		}
	}
	return false
}

// detachKeyBindWindow removes all callbacks associated with a particular
// window and event type (either KeyPress or KeyRelease)
// Also decrements the counter in the corresponding 'keygrabs' map
// appropriately.
// This is exported for use in the keybind package. It should not be used.
// To detach a window from a key binding callbacks, please use keybind.Detach.
// (This method will issue an Ungrab requests, while keybind.Detach will.)
func detachKeyBindWindow(xu *xgbutil.XUtil, evtype int, win xproto.Window) {
	xu.KeybindsLck.Lock()
	defer xu.KeybindsLck.Unlock()

	// Since we can't create a full key, loop through all key binds
	// and check if evtype and window match.
	for key, _ := range xu.Keybinds {
		if key.Evtype == evtype && key.Win == win {
			xu.Keygrabs[key] -= len(xu.Keybinds[key])
			delete(xu.Keybinds, key)
		}
	}
}

// keyBindGrabs returns the number of grabs on a particular
// event/window/mods/keycode combination. Namely, this combination
// uniquely identifies a grab. If it's repeated, we get BadAccess.
// The idea is that if there are 0 grabs on a particular (modifiers, keycode)
// tuple, then we issue a grab request. Otherwise, we don't.
// This is exported for use in the keybind package. It should not be used.
func keyBindGrabs(xu *xgbutil.XUtil, evtype int, win xproto.Window, mods uint16,
	keycode xproto.Keycode) int {

	xu.KeybindsLck.RLock()
	defer xu.KeybindsLck.RUnlock()

	key := xgbutil.KeyBindKey{evtype, win, mods, keycode}
	return xu.Keygrabs[key] // returns 0 if key does not exist
}

// KeyMapGet accessor.
func KeyMapGet(xu *xgbutil.XUtil) *xgbutil.KeyboardMapping {
	return xu.Keymap
}

// KeyMapSet updates XUtil.keymap.
// This is exported for use in the keybind package. You probably shouldn't
// use this. (You may need to use this if you're rolling your own event loop,
// and still want to use the keybind package.)
func KeyMapSet(xu *xgbutil.XUtil, keyMapReply *xproto.GetKeyboardMappingReply) {
	xu.Keymap = &xgbutil.KeyboardMapping{keyMapReply}
}

// ModMapGet accessor.
func ModMapGet(xu *xgbutil.XUtil) *xgbutil.ModifierMapping {
	return xu.Modmap
}

// ModMapSet updates XUtil.modmap.
// This is exported for use in the keybind package. You probably shouldn't
// use this. (You may need to use this if you're rolling your own event loop,
// and still want to use the keybind package.)
func ModMapSet(xu *xgbutil.XUtil, modMapReply *xproto.GetModifierMappingReply) {
	xu.Modmap = &xgbutil.ModifierMapping{modMapReply}
}
