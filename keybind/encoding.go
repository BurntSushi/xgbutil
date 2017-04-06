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
// XXX: We don't support ShiftLock, only CapsLock
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

	// TODO(dh): do we need to handle ISO_Level3_Lock as well, or is
	// the X server doing that for us?
	var group []xproto.Keysym
	switch {
	case level3:
		group = []xproto.Keysym{k5, k6}
	case mode:
		group = []xproto.Keysym{k3, k4}
	default:
		group = []xproto.Keysym{k1, k2}
	}
	switch {
	case !shift && !lock:
		return KeysymToStr(group[0])
	case !shift && lock:
		// The Shift modifier is off, and the Lock modifier is on and
		// is interpreted as CapsLock. In this case, the first KeySym
		// is used, but if that KeySym is lowercase alphabetic, then
		// the corresponding uppercase KeySym is used instead.

		lower, upper := convertCase(group[0])
		if lower == group[0] {
			// either group[0] is alphabetic and lower case, or lower
			// == upper == group[0]
			return KeysymToStr(upper)
		} else {
			return KeysymToStr(group[0])
		}
	case shift && lock:
		// The Shift modifier is on, and the Lock modifier is on and
		// is interpreted as CapsLock. In this case, the second KeySym
		// is used, but if that KeySym is lowercase alphabetic, then
		// the corresponding uppercase KeySym is used instead.
		lower, upper := convertCase(group[1])
		if lower == group[1] {
			// either groups[1] is alphabetic and lower case, or lower
			// == upper == group[1]
			return KeysymToStr(upper)
		} else {
			return KeysymToStr(group[1])
		}
	case shift:
		return KeysymToStr(group[1])
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
	ks1, ks2, ks3, ks4, ks5, ks6 xproto.Keysym) {

	ks1 = KeysymGet(xu, keycode, 0)
	ks2 = KeysymGet(xu, keycode, 1)
	ks3 = KeysymGet(xu, keycode, 2)
	ks4 = KeysymGet(xu, keycode, 3)
	ks5 = KeysymGet(xu, keycode, 4)
	ks6 = KeysymGet(xu, keycode, 5)

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

	// follow the rules, fourth paragraph
	if ks2 == 0 {
		ks1, ks2 = convertCase(ks1)
	}
	if ks4 == 0 {
		ks3, ks4 = convertCase(ks3)
	}

	// TODO(dh): Again, no rules are specified for groups 3 and 4, so
	// we'll not do anything for now.

	return
}

// convertCase converts sym to its lower- and uppercase versions,
// assuming sym is alphabetic. Otherwise, it returns lower == upper ==
// sym.
func convertCase(sym xproto.Keysym) (lower, upper xproto.Keysym) {
	// This function is modelled after xcb_convert_case.
	lower = sym
	upper = sym

	switch sym >> 8 {
	case 0: /* Latin 1 */
		if (sym >= keysyms["A"]) && (sym <= keysyms["Z"]) {
			lower += (keysyms["a"] - keysyms["A"])
		} else if (sym >= keysyms["a"]) && (sym <= keysyms["z"]) {
			upper -= (keysyms["a"] - keysyms["A"])
		} else if (sym >= keysyms["Agrave"]) && (sym <= keysyms["Odiaeresis"]) {
			lower += (keysyms["agrave"] - keysyms["Agrave"])
		} else if (sym >= keysyms["agrave"]) && (sym <= keysyms["odiaeresis"]) {
			upper -= (keysyms["agrave"] - keysyms["Agrave"])
		} else if (sym >= keysyms["Ooblique"]) && (sym <= keysyms["Thorn"]) {
			lower += (keysyms["oslash"] - keysyms["Ooblique"])
		} else if (sym >= keysyms["oslash"]) && (sym <= keysyms["thorn"]) {
			upper -= (keysyms["oslash"] - keysyms["Ooblique"])
		}
	case 1: /* Latin 2 */
		/* Assume the KeySym is a legal value (ignore discontinuities) */
		if sym == keysyms["Aogonek"] {
			lower = keysyms["aogonek"]
		} else if sym >= keysyms["Lstroke"] && sym <= keysyms["Sacute"] {
			lower += (keysyms["lstroke"] - keysyms["Lstroke"])
		} else if sym >= keysyms["Scaron"] && sym <= keysyms["Zacute"] {
			lower += (keysyms["scaron"] - keysyms["Scaron"])
		} else if sym >= keysyms["Zcaron"] && sym <= keysyms["Zabovedot"] {
			lower += (keysyms["zcaron"] - keysyms["Zcaron"])
		} else if sym == keysyms["aogonek"] {
			upper = keysyms["Aogonek"]
		} else if sym >= keysyms["lstroke"] && sym <= keysyms["sacute"] {
			upper -= (keysyms["lstroke"] - keysyms["Lstroke"])
		} else if sym >= keysyms["scaron"] && sym <= keysyms["zacute"] {
			upper -= (keysyms["scaron"] - keysyms["Scaron"])
		} else if sym >= keysyms["zcaron"] && sym <= keysyms["zabovedot"] {
			upper -= (keysyms["zcaron"] - keysyms["Zcaron"])
		} else if sym >= keysyms["Racute"] && sym <= keysyms["Tcedilla"] {
			lower += (keysyms["racute"] - keysyms["Racute"])
		} else if sym >= keysyms["racute"] && sym <= keysyms["tcedilla"] {
			upper -= (keysyms["racute"] - keysyms["Racute"])
		}
	case 2: /* Latin 3 */
		/* Assume the KeySym is a legal value (ignore discontinuities) */
		if sym >= keysyms["Hstroke"] && sym <= keysyms["Hcircumflex"] {
			lower += (keysyms["hstroke"] - keysyms["Hstroke"])
		} else if sym >= keysyms["Gbreve"] && sym <= keysyms["Jcircumflex"] {
			lower += (keysyms["gbreve"] - keysyms["Gbreve"])
		} else if sym >= keysyms["hstroke"] && sym <= keysyms["hcircumflex"] {
			upper -= (keysyms["hstroke"] - keysyms["Hstroke"])
		} else if sym >= keysyms["gbreve"] && sym <= keysyms["jcircumflex"] {
			upper -= (keysyms["gbreve"] - keysyms["Gbreve"])
		} else if sym >= keysyms["Cabovedot"] && sym <= keysyms["Scircumflex"] {
			lower += (keysyms["cabovedot"] - keysyms["Cabovedot"])
		} else if sym >= keysyms["cabovedot"] && sym <= keysyms["scircumflex"] {
			upper -= (keysyms["cabovedot"] - keysyms["Cabovedot"])
		}
	case 3: /* Latin 4 */
		/* Assume the KeySym is a legal value (ignore discontinuities) */
		if sym >= keysyms["Rcedilla"] && sym <= keysyms["Tslash"] {
			lower += (keysyms["rcedilla"] - keysyms["Rcedilla"])
		} else if sym >= keysyms["rcedilla"] && sym <= keysyms["tslash"] {
			upper -= (keysyms["rcedilla"] - keysyms["Rcedilla"])
		} else if sym == keysyms["ENG"] {
			lower = keysyms["eng"]
		} else if sym == keysyms["eng"] {
			upper = keysyms["ENG"]
		} else if sym >= keysyms["Amacron"] && sym <= keysyms["Umacron"] {
			lower += (keysyms["amacron"] - keysyms["Amacron"])
		} else if sym >= keysyms["amacron"] && sym <= keysyms["umacron"] {
			upper -= (keysyms["amacron"] - keysyms["Amacron"])
		}
	case 6: /* Cyrillic */
		/* Assume the KeySym is a legal value (ignore discontinuities) */
		if sym >= keysyms["Serbian_DJE"] && sym <= keysyms["Serbian_DZE"] {
			lower -= (keysyms["Serbian_DJE"] - keysyms["Serbian_dje"])
		} else if sym >= keysyms["Serbian_dje"] && sym <= keysyms["Serbian_dze"] {
			upper += (keysyms["Serbian_DJE"] - keysyms["Serbian_dje"])
		} else if sym >= keysyms["Cyrillic_YU"] && sym <= keysyms["Cyrillic_HARDSIGN"] {
			lower -= (keysyms["Cyrillic_YU"] - keysyms["Cyrillic_yu"])
		} else if sym >= keysyms["Cyrillic_yu"] && sym <= keysyms["Cyrillic_hardsign"] {
			upper += (keysyms["Cyrillic_YU"] - keysyms["Cyrillic_yu"])
		}
	case 7: /* Greek */
		/* Assume the KeySym is a legal value (ignore discontinuities) */
		if sym >= keysyms["Greek_ALPHAaccent"] && sym <= keysyms["Greek_OMEGAaccent"] {
			lower += (keysyms["Greek_alphaaccent"] - keysyms["Greek_ALPHAaccent"])
		} else if sym >= keysyms["Greek_alphaaccent"] && sym <= keysyms["Greek_omegaaccent"] &&
			sym != keysyms["Greek_iotaaccentdieresis"] &&
			sym != keysyms["Greek_upsilonaccentdieresis"] {
			upper -= (keysyms["Greek_alphaaccent"] - keysyms["Greek_ALPHAaccent"])
		} else if sym >= keysyms["Greek_ALPHA"] && sym <= keysyms["Greek_OMEGA"] {
			lower += (keysyms["Greek_alpha"] - keysyms["Greek_ALPHA"])
		} else if sym >= keysyms["Greek_alpha"] && sym <= keysyms["Greek_omega"] &&
			sym != keysyms["Greek_finalsmallsigma"] {
			upper -= (keysyms["Greek_alpha"] - keysyms["Greek_ALPHA"])
		}
	case 0x14: /* Armenian */
		if sym >= keysyms["Armenian_AYB"] && sym <= keysyms["Armenian_fe"] {
			lower = sym | 1
			upper = sym & ^xproto.Keysym(1)
		}
	}

	return lower, upper
}
