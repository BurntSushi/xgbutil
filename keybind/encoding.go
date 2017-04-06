package keybind

/*
This file contains the logic to implement X's Keyboard Encoding
described here: http://goo.gl/qum9q

Essentially, LookupString is analogous to Xlib's XLookupString. It's useful
in determining the english string representation of modifiers + keycode.

It is not for the faint of heart.
*/

import (
	"strings"
	"unicode"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
)

// LookupString attempts to convert a (modifiers, keycode) to an english string.
// It essentially implements the rules described at http://goo.gl/qum9q
// Namely, the bulleted list that describes how key syms should be interpreted
// when various modifiers are pressed.
// Note that we ignore the logic that asks us to check if particular key codes
// are mapped to particular modifiers (i.e., "XK_Caps_Lock" to "Lock" modifier).
// We just check if the modifiers are activated. That's good enough for me.
// XXX: We ignore num lock stuff.
func LookupString(xu *xgbutil.XUtil, mods uint16,
	keycode xproto.Keycode) string {

	modMap := ModMapGet(xu)
	symToMod := map[xproto.Keysym]uint16{}
	for i, kc := range modMap.Keycodes {
		if kc == 0 {
			continue
		}
		// map keysyms to modifiers. we're only really interested in
		// Mode_switch, ISO_Level3_Shift and Num_Lock, though.
		symToMod[KeysymGet(xu, kc, 0)] = Modifiers[byte(i)/modMap.KeycodesPerModifier]
	}

	k1, k2, k3, k4, k5, k6 := interpretSymList(xu, keycode)

	shift := mods&xproto.ModMaskShift > 0
	lock := mods&xproto.ModMaskLock > 0
	mode := mods&symToMod[keysyms["Mode_switch"]] > 0
	level3 := mods&symToMod[keysyms["ISO_Level3_Shift"]] > 0

	var group []string
	switch {
	case level3:
		group = []string{k5, k6}
	case mode:
		group = []string{k3, k4}
	default:
		group = []string{k1, k2}
	}
	// TODO(dh): do we need to handle ISO_Level3_Lock as well, or is
	// the X server doing that for us?
	switch {
	case !shift && !lock:
		return group[0]
	case !shift && lock:
		if len(group[0]) == 1 && unicode.IsLower(rune(group[0][0])) {
			return group[1]
		} else {
			return group[0]
		}
	case shift && lock:
		if len(group[1]) == 1 && unicode.IsLower(rune(group[1][0])) {
			return string(unicode.ToUpper(rune(group[1][0])))
		} else {
			return group[1]
		}
	case shift:
		return group[1]
	}

	return ""
}

// ModifierString takes in a keyboard state and returns a string of all
// modifiers in the state.
func ModifierString(mods uint16) string {
	modStrs := make([]string, 0, 3)
	for i, mod := range Modifiers {
		if mod&mods > 0 && len(NiceModifiers[i]) > 0 {
			modStrs = append(modStrs, NiceModifiers[i])
		}
	}
	return strings.Join(modStrs, "-")
}

// KeyMatch returns true if a string representation of a key can
// be matched (case insensitive) to the (modifiers, keycode) tuple provided.
// String representations can be found in keybind/keysymdef.go
func KeyMatch(xu *xgbutil.XUtil,
	keyStr string, mods uint16, keycode xproto.Keycode) bool {

	guess := LookupString(xu, mods, keycode)
	return strings.ToLower(guess) == strings.ToLower(keyStr)
}

// interpretSymList interprets the keysym list for a particular keycode as
// described in the third and fourth paragraphs of http://goo.gl/qum9q
func interpretSymList(xu *xgbutil.XUtil, keycode xproto.Keycode) (
	k1, k2, k3, k4, k5, k6 string) {

	ks1 := KeysymGet(xu, keycode, 0)
	ks2 := KeysymGet(xu, keycode, 1)
	ks3 := KeysymGet(xu, keycode, 2)
	ks4 := KeysymGet(xu, keycode, 3)
	ks5 := KeysymGet(xu, keycode, 4)
	ks6 := KeysymGet(xu, keycode, 5)

	// follow the rules, third paragraph
	switch {
	case ks2 == 0 && ks3 == 0 && ks4 == 0:
		// If the list (ignoring trailing NoSymbol entries) is a
		// single KeySym ``K'', then the list is treated as if it were
		// the list ``K NoSymbol K NoSymbol''.
		ks3 = ks1
	case ks3 == 0 && ks4 == 0:
		// If the list (ignoring trailing NoSymbol entries) is a pair
		// of KeySyms ``K1 K2'', then the list is treated as if it
		// were the list ``K1 K2 K1 K2''.
		ks3 = ks1
		ks4 = ks2
	case ks4 == 0:
		// If the list (ignoring trailing NoSymbol entries) is a
		// triple of KeySyms ``K1 K2 K3'', then the list is treated as
		// if it were the list ``K1 K2 K3 NoSymbol''.
		ks4 = 0
	}

	// TODO(dh): The document doesn't specify any rules for keysyms 5
	// and 6. For now, we'll do nothing.

	// Now convert keysyms to strings, so we can do alphabetic shit.
	k1 = KeysymToStr(ks1)
	k2 = KeysymToStr(ks2)
	k3 = KeysymToStr(ks3)
	k4 = KeysymToStr(ks4)
	k5 = KeysymToStr(ks5)
	k6 = KeysymToStr(ks6)

	// follow the rules, fourth paragraph
	if k2 == "" {
		if len(k1) == 1 && unicode.IsLetter(rune(k1[0])) {
			k1 = string(unicode.ToLower(rune(k1[0])))
			k2 = string(unicode.ToUpper(rune(k1[0])))
		} else {
			k2 = k1
		}
	}
	if k4 == "" {
		if len(k3) == 1 && unicode.IsLetter(rune(k3[0])) {
			k3 = string(unicode.ToLower(rune(k3[0])))
			k4 = string(unicode.ToUpper(rune(k4[0])))
		} else {
			k4 = k3
		}
	}

	// TODO(dh): Again, no rules are specified for groups 3 and 4, so
	// we'll not do anything for now.

	return
}
