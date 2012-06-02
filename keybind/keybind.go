package keybind

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
)

var Modifiers []uint16 = []uint16{ // order matters!
	xproto.ModMaskShift, xproto.ModMaskLock, xproto.ModMaskControl,
	xproto.ModMask1, xproto.ModMask2, xproto.ModMask3,
	xproto.ModMask4, xproto.ModMask5,
	xproto.ModMaskAny,
}

// Initialize attaches the appropriate callbacks to make key bindings easier.
// i.e., update state of the world on a MappingNotify.
func Initialize(xu *xgbutil.XUtil) {
	// Listen to mapping notify events
	xevent.MappingNotifyFun(updateMaps).Connect(xu, xevent.NoWindow)

	// Give us an initial mapping state...
	keyMap, modMap := MapsGet(xu)
	KeyMapSet(xu, keyMap)
	ModMapSet(xu, modMap)
}

// updateMaps runs in response to MappingNotify events.
// It is responsible for making sure our view of the world's keyboard
// and modifier maps is correct. (Pointer mappings should be handled in
// a similar callback in the mousebind package.)
func updateMaps(xu *xgbutil.XUtil, e xevent.MappingNotifyEvent) {
	keyMap, modMap := MapsGet(xu)

	// Hold up... If this is MappingKeyboard, then we may have some keycode
	// changes. This is GROSS. We basically need to go through each keycode
	// in the map, look up the keysym using the new map and use that keysym
	// to look up the keycode in our current map. If the current keycode
	// does not equal the old keycode, then we log the change in a map of
	// old keycode -> new keycode.
	// Once the map is constructed, we look through all of our keybindings
	// and updated appropriately. *puke*
	// I am only somewhat confident that this is correct.
	if e.Request == xproto.MappingKeyboard {
		changes := make(map[xproto.Keycode]xproto.Keycode, 0)
		xuKeyMap := &xgbutil.KeyboardMapping{keyMap}

		min, max := minMaxKeycodeGet(xu)

		// let's not do too much allocation in our loop, shall we?
		var newSym, oldSym xproto.Keysym
		var oldKc xproto.Keycode
		var column byte

		// wrap 'int(..)' around bytes min and max to avoid overflow. Hideous.
		for newKc := int(min); newKc <= int(max); newKc++ {
			for column = 0; byte(column) < keyMap.KeysymsPerKeycode; column++ {
				// use new key map
				newSym = KeysymGetWithMap(xu, xuKeyMap, xproto.Keycode(newKc),
					column)

				// uses old key map
				oldKc = keycodeGet(xu, newSym)
				oldSym = KeysymGet(xu, xproto.Keycode(newKc), column)

				// If the old and new keysyms are the same, ignore!
				// Also ignore if either keysym is VoidSymbol
				if oldSym == newSym || oldSym == 0 || newSym == 0 {
					continue
				}

				// these should match if there are NO changes
				if oldKc != xproto.Keycode(newKc) {
					changes[oldKc] = xproto.Keycode(newKc)
				}
			}
		}

		// Now use 'changes' to do some regrabbing
		// Loop through all key bindings and check if we have any affected
		// key codes. (Note that each key binding may be associated with
		// multiple callbacks.)
		// We must ungrab everything first, in case two keys are being swapped.
		for _, key := range keyKeys(xu) {
			if _, ok := changes[key.Code]; ok {
				Ungrab(xu, key.Win, key.Mod, key.Code)
			}
		}
		// Okay, now grab.
		for _, key := range keyKeys(xu) {
			if newKc, ok := changes[key.Code]; ok {
				Grab(xu, key.Win, key.Mod, newKc)
				updateKeyBindKey(xu, key, newKc)
			}
		}
	}

	// We don't have to do something with MappingModifier like we do with
	// MappingKeyboard. This is due to us requiring that key strings use
	// modifier names built into X. (i.e., the names seen in the output of
	// `xmodmap`.) This means that the modifier mappings happen on the X server
	// side, so we don't *typically* have to care what key is actually being
	// pressed to trigger a modifier. (There are some exceptional cases, and
	// when that happens, we simply query on-demand which keys are modifiers.
	// See the RunKey{Press,Release}Callbacks functions in keybind/callback.go
	// for the deets.)

	// Finally update our view of the mappings.
	KeyMapSet(xu, keyMap)
	ModMapSet(xu, modMap)
}

// minMaxKeycodeGet a simple accessor to the X setup info to return the
// minimum and maximum keycodes. They are typically 8 and 255, respectively.
func minMaxKeycodeGet(xu *xgbutil.XUtil) (xproto.Keycode, xproto.Keycode) {
	return xu.Setup().MinKeycode, xu.Setup().MaxKeycode
}

// A convenience function to grab the KeyboardMapping and ModifierMapping
// from X. We need to do this on startup (see Initialize) and whenever we
// get a MappingNotify event.
func MapsGet(xu *xgbutil.XUtil) (*xproto.GetKeyboardMappingReply,
	*xproto.GetModifierMappingReply) {

	min, max := minMaxKeycodeGet(xu)
	newKeymap, keyErr := xproto.GetKeyboardMapping(xu.Conn(), min,
		byte(max-min+1)).Reply()
	newModmap, modErr := xproto.GetModifierMapping(xu.Conn()).Reply()

	// If there are errors, we really need to panic. We just can't do
	// any key binding without a mapping from the server.
	if keyErr != nil {
		panic(fmt.Sprintf("COULD NOT GET KEYBOARD MAPPING: %v\n"+
			"THIS IS AN UNRECOVERABLE ERROR.\n",
			keyErr))
	}
	if modErr != nil {
		panic(fmt.Sprintf("COULD NOT GET MODIFIER MAPPING: %v\n"+
			"THIS IS AN UNRECOVERABLE ERROR.\n",
			keyErr))
	}

	return newKeymap, newModmap
}

// ParseString takes a string of the format '[Mod[-Mod[...]]]-KEY',
// i.e., 'Mod4-j', and returns a modifiers/keycode combo.
// An error is returned if the string is malformed, or if no valid KEY can
// be found. 
// Valid values of KEY should include almost anything returned by pressing
// keys with the 'xev' program. Alternatively, you may reference the keys
// of the 'keysyms' map defined in keybind/keysymdef.go.
func ParseString(xu *xgbutil.XUtil, s string) (uint16, xproto.Keycode, error) {
	mods, kc := uint16(0), xproto.Keycode(0)
	for _, part := range strings.Split(s, "-") {
		switch strings.ToLower(part) {
		case "shift":
			mods |= xproto.ModMaskShift
		case "lock":
			mods |= xproto.ModMaskLock
		case "control":
			mods |= xproto.ModMaskControl
		case "mod1":
			mods |= xproto.ModMask1
		case "mod2":
			mods |= xproto.ModMask2
		case "mod3":
			mods |= xproto.ModMask3
		case "mod4":
			mods |= xproto.ModMask4
		case "mod5":
			mods |= xproto.ModMask5
		case "any":
			mods |= xproto.ModMaskAny
		default: // a key code!
			if kc == 0 { // only accept the first keycode we see
				kc = StrToKeycode(xu, part)
			}
		}
	}

	if kc == 0 {
		return 0, 0, fmt.Errorf("Could not find a valid keycode in the "+
			"string '%s'. Key binding failed.", s)
	}

	return mods, kc, nil
}

// StrToKeycode is a wrapper around keycodeGet meant to make our search
// a bit more flexible if needed. (i.e., case-insensitive)
func StrToKeycode(xu *xgbutil.XUtil, str string) xproto.Keycode {
	// Do some fancy case stuff before we give up.
	sym, ok := keysyms[str]
	if !ok {
		sym, ok = keysyms[strings.Title(str)]
	}
	if !ok {
		sym, ok = keysyms[strings.ToLower(str)]
	}
	if !ok {
		sym, ok = keysyms[strings.ToUpper(str)]
	}

	// If we don't know what 'str' is, return 0.
	// There will probably be a bad access. We should do better than that...
	if !ok {
		return xproto.Keycode(0)
	}
	return keycodeGet(xu, sym)
}

// keysymsPer gets the number of keysyms per keycode for the current key map.
func keysymsPer(xu *xgbutil.XUtil) int {
	return int(KeyMapGet(xu).KeysymsPerKeycode)
}

// Given a keysym, find the keycode mapped to it in the current X environment.
// keybind.Initialize MUST have been called before using this function.
func keycodeGet(xu *xgbutil.XUtil, keysym xproto.Keysym) xproto.Keycode {
	min, max := minMaxKeycodeGet(xu)
	keyMap := KeyMapGet(xu)
	if keyMap == nil {
		panic("keybind.Initialize must be called before using the keybind " +
			"package.")
	}

	var c byte
	for kc := int(min); kc <= int(max); kc++ {
		for c = 0; c < keyMap.KeysymsPerKeycode; c++ {
			if keysym == KeysymGet(xu, xproto.Keycode(kc), c) {
				return xproto.Keycode(kc)
			}
		}
	}
	return 0
}

// KeysymToStr converts a keysym to a string if one is available.
// If one is found, KeysymToStr also checks the 'weirdKeysyms' map, which
// contains a map from multi-character strings to single character
// representations (i.e., 'braceleft' to '{').
// If no match is found initially, an empty string is returned.
func KeysymToStr(keysym xproto.Keysym) string {
	symStr, ok := strKeysyms[keysym]
	if !ok {
		return ""
	}

	shortSymStr, ok := weirdKeysyms[symStr]
	if ok {
		return string(shortSymStr)
	}

	return symStr
}

// KeysymGet is a shortcut alias for 'KeysymGetWithMap' using the current
// keymap stored in XUtil.
// keybind.Initialize MUST have been called before using this function.
func KeysymGet(xu *xgbutil.XUtil, keycode xproto.Keycode,
	column byte) xproto.Keysym {

	return KeysymGetWithMap(xu, KeyMapGet(xu), keycode, column)
}

// KeysymGetWithMap uses the given key map and finds a keysym associated
// with the given keycode in the current X environment.
func KeysymGetWithMap(xu *xgbutil.XUtil, keyMap *xgbutil.KeyboardMapping,
	keycode xproto.Keycode, column byte) xproto.Keysym {

	min, _ := minMaxKeycodeGet(xu)
	i := (int(keycode)-int(min))*int(keyMap.KeysymsPerKeycode) + int(column)

	return keyMap.Keysyms[i]
}

// ModGet finds the modifier currently associated with a given keycode.
// If a modifier doesn't exist for this keycode, then 0 is returned.
func ModGet(xu *xgbutil.XUtil, keycode xproto.Keycode) uint16 {
	modMap := ModMapGet(xu)

	var i byte
	for i = 0; int(i) < len(modMap.Keycodes); i++ {
		if modMap.Keycodes[i] == keycode {
			return Modifiers[i/modMap.KeycodesPerModifier]
		}
	}
	return 0
}

// Grab grabs a key with mods on a particular window.
// This will also grab all combinations of modifiers found in xevent.IgnoreMods.
func Grab(xu *xgbutil.XUtil, win xproto.Window,
	mods uint16, key xproto.Keycode) {

	for _, m := range xevent.IgnoreMods {
		xproto.GrabKey(xu.Conn(), true, win, mods|m, key,
			xproto.GrabModeAsync, xproto.GrabModeAsync)
	}
}

// GrabChecked Grabs a key with mods on a particular window.
// This is the same as Grab, except that it issue a checked request.
// Which means that an error could be returned and handled on the spot.
// (Checked requests are slower than unchecked requests.)
// This will also grab all combinations of modifiers found in xevent.IgnoreMods.
func GrabChecked(xu *xgbutil.XUtil, win xproto.Window,
	mods uint16, key xproto.Keycode) error {

	var err error
	for _, m := range xevent.IgnoreMods {
		err = xproto.GrabKeyChecked(xu.Conn(), true, win, mods|m, key,
			xproto.GrabModeAsync, xproto.GrabModeAsync).Check()
		if err != nil {
			return err
		}
	}
	return nil
}

// Ungrab undoes Grab. It will handle all combinations od modifiers found
// in xevent.IgnoreMods.
func Ungrab(xu *xgbutil.XUtil, win xproto.Window,
	mods uint16, key xproto.Keycode) {

	for _, m := range xevent.IgnoreMods {
		xproto.UngrabKey(xu.Conn(), key, win, mods|m)
	}
}

// GrabKeyboard grabs the entire keyboard.
// Returns whether GrabStatus is successful and an error if one is reported by 
// XGB. It is possible to not get an error and the grab to be unsuccessful.
// The purpose of 'win' is that after a grab is successful, ALL Key*Events will
// be sent to that window. Make sure you have a callback attached :-)
func GrabKeyboard(xu *xgbutil.XUtil, win xproto.Window) error {
	reply, err := xproto.GrabKeyboard(xu.Conn(), false, win, 0,
		xproto.GrabModeAsync, xproto.GrabModeAsync).Reply()
	if err != nil {
		return fmt.Errorf("GrabKeyboard: Error grabbing keyboard on "+
			"window '%x': %s", win, err)
	}

	switch reply.Status {
	case xproto.GrabStatusSuccess:
		// all is well
	case xproto.GrabStatusAlreadyGrabbed:
		return fmt.Errorf("GrabKeyboard: Could not grab keyboard. " +
			"Status: AlreadyGrabbed.")
	case xproto.GrabStatusInvalidTime:
		return fmt.Errorf("GrabKeyboard: Could not grab keyboard. " +
			"Status: InvalidTime.")
	case xproto.GrabStatusNotViewable:
		return fmt.Errorf("GrabKeyboard: Could not grab keyboard. " +
			"Status: NotViewable.")
	case xproto.GrabStatusFrozen:
		return fmt.Errorf("GrabKeyboard: Could not grab keyboard. " +
			"Status: Frozen.")
	}
	return nil
}

// UngrabKeyboard undoes GrabKeyboard.
func UngrabKeyboard(xu *xgbutil.XUtil) {
	xproto.UngrabKeyboard(xu.Conn(), 0)
}

// SmartGrab grabs the keyboard for the given window, and redirects all
// key events in the xevent main event loop to avoid races.
func SmartGrab(xu *xgbutil.XUtil, win xproto.Window) error {
	err := GrabKeyboard(xu, win)
	if err != nil {
		return fmt.Errorf("SmartGrab: %s", err)
	}

	// Now redirect all key events to the dummy window to prevent races
	xevent.RedirectKeyEvents(xu, win)

	return nil
}

// SmartUngrab reverses SmartGrab and stops redirecting all key events.
func SmartUngrab(xu *xgbutil.XUtil) {
	UngrabKeyboard(xu)

	// Stop redirecting all key events
	xevent.RedirectKeyEvents(xu, 0)
}

// DummyGrab grabs the keyboard and sends all key events to the dummy window.
func DummyGrab(xu *xgbutil.XUtil) error {
	return SmartGrab(xu, xu.Dummy())
}

// DummyUngrab ungrabs the keyboard from the dummy window.
func DummyUngrab(xu *xgbutil.XUtil) {
	SmartUngrab(xu)
}
