package keybind

/*
This file contains the logic to implement X's Keyboard Encoding
described here: http://goo.gl/qum9q

Essentially, LookupString is analogous to Xlib's XLookupString. It's useful
in determining the english string representation of modifiers + keycode.

It is not for the faint of heart.
*/

import (
	"unicode"

	"code.google.com/p/jamslam-x-go-binding/xgb"

	"github.com/BurntSushi/xgbutil"
)

// interpretSymList interprets the keysym list for a particular keycode as
// described in the third and fourth paragraphs of http://goo.gl/qum9q
func interpretSymList(xu *xgbutil.XUtil, keycode byte) (
	k1 rune, k2 rune, k3 rune, k4 rune) {

	ks1 := keysymGet(xu, keycode, 0)
	ks2 := keysymGet(xu, keycode, 1)
	ks3 := keysymGet(xu, keycode, 2)
	ks4 := keysymGet(xu, keycode, 3)

	// follow the rules, third paragraph
	switch {
	case ks2 == 0 && ks3 == 0 && ks4 == 0:
		ks3 = ks1
	case ks3 == 0 && ks4 == 0:
		ks3 = ks1
		ks4 = ks2
	case ks4 == 0:
		ks4 = 0
	}

	// Now convert keysyms to runes, so we can do alphabetic shit.
	k1 = keysymToRune(ks1)
	k2 = keysymToRune(ks2)
	k3 = keysymToRune(ks3)
	k4 = keysymToRune(ks4)

	// follow the rules, fourth paragraph
	if k2 == 0 {
		if unicode.IsLetter(k1) {
			k1 = unicode.ToLower(k1)
			k2 = unicode.ToUpper(k1)
		} else {
			k2 = k1
		}
	}
	if k4 == 0 {
		if unicode.IsLetter(k3) {
			k3 = unicode.ToLower(k3)
			k4 = unicode.ToUpper(k4)
		} else {
			k4 = k3
		}
	}

	return
}

// LookupString attempts to convert a (modifiers, keycode) to an english string.
// It essentially implements the rules described at http://goo.gl/qum9q
// Namely, the bulleted list that describes how key syms should be interpreted
// when various modifiers are pressed.
// Note that we ignore the logic that asks us to check if particular key codes
// are mapped to particular modifiers (i.e., "XK_Caps_Lock" to "Lock" modifier).
// We just check if the modifiers are activated. That's good enough for me.
// XXX: We ignore num lock stuff.
// XXX: We ignore MODE SWITCH stuff. (i.e., we don't use group 2 key syms.)
func LookupString(xu *xgbutil.XUtil, mods uint16, keycode byte) rune {
	k1, k2, _, _ := interpretSymList(xu, keycode)

	shift := mods&xgb.ModMaskShift > 0
	lock := mods&xgb.ModMaskLock > 0
	switch {
	case !shift && !lock:
		return k1
	case !shift && lock:
		if unicode.IsLower(k1) {
			return k2
		} else {
			return k1
		}
	case shift && lock:
		if unicode.IsLower(k2) {
			return unicode.ToUpper(k2)
		} else {
			return k2
		}
	case shift:
		return k2
	}

	return 0
}
