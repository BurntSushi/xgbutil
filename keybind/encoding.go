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

// LookupKeysym attempts to convert a (modifiers, keycode) to an english string.
// It essentially implements the rules described at http://goo.gl/qum9q
// Namely, the bulleted list that describes how key syms should be interpreted
// when various modifiers are pressed.
// Note that we ignore the logic that asks us to check if particular key codes
// are mapped to particular modifiers (i.e., "XK_Caps_Lock" to "Lock" modifier).
// We just check if the modifiers are activated. That's good enough for me.
// XXX: We don't support ShiftLock, only CapsLock
func LookupKeysym(xu *xgbutil.XUtil, mods uint16, keycode xproto.Keycode) xproto.Keysym {
	var modeMod, level3Mod, numlockMod uint16
	modMap := ModMapGet(xu)
	for i, kc := range modMap.Keycodes {
		if kc == 0 {
			continue
		}

		mod := Modifiers[byte(i)/modMap.KeycodesPerModifier]
		switch KeysymGet(xu, kc, 0) {
		case Keysyms["Mode_switch"]:
			modeMod = mod
		case Keysyms["ISO_Level3_Shift"]:
			level3Mod = mod
		case Keysyms["Num_Lock"]:
			numlockMod = mod
		}
	}

	k1, k2, k3, k4, k5, k6 := interpretSymList(xu, keycode)

	shift := mods&xproto.ModMaskShift > 0
	lock := mods&xproto.ModMaskLock > 0
	mode := mods&modeMod > 0
	level3 := mods&level3Mod > 0
	numpad := mods&numlockMod > 0

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

	if numpad && ((group[1] >= 0xFF80 && group[1] <= 0xFFBD) || (group[1] >= 0x11000000 && group[1] <= 0x1100FFFF)) {
		// TODO(dh): if Lock is on and is ShiftLock (as opposed to
		// CapsLock), it should be treated like Shift. Currently, this
		// function doesn't differentiate between CapsLock and
		// ShiftLock, so we're skipping that step. Luckily, ShiftLock
		// is very rare nowadays.
		if shift {
			return group[0]
		}
		return group[1]
	}
	switch {
	case !shift && !lock:
		return group[0]
	case !shift && lock:
		// The Shift modifier is off, and the Lock modifier is on and
		// is interpreted as CapsLock. In this case, the first KeySym
		// is used, but if that KeySym is lowercase alphabetic, then
		// the corresponding uppercase KeySym is used instead.

		lower, upper := convertCase(group[0])
		if lower == group[0] {
			// either group[0] is alphabetic and lower case, or lower
			// == upper == group[0]
			return upper
		} else {
			return group[0]
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
			return upper
		} else {
			return group[1]
		}
	case shift:
		return group[1]
	}

	return 0
}

// LookupString is a convenience function that applies KeysymToStr to
// the result of LookupKeysym.
func LookupString(xu *xgbutil.XUtil, mods uint16,
	keycode xproto.Keycode) string {

	return KeysymToStr(LookupKeysym(xu, mods, keycode))
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
		if (sym >= Keysyms["A"]) && (sym <= Keysyms["Z"]) {
			lower += (Keysyms["a"] - Keysyms["A"])
		} else if (sym >= Keysyms["a"]) && (sym <= Keysyms["z"]) {
			upper -= (Keysyms["a"] - Keysyms["A"])
		} else if (sym >= Keysyms["Agrave"]) && (sym <= Keysyms["Odiaeresis"]) {
			lower += (Keysyms["agrave"] - Keysyms["Agrave"])
		} else if (sym >= Keysyms["agrave"]) && (sym <= Keysyms["odiaeresis"]) {
			upper -= (Keysyms["agrave"] - Keysyms["Agrave"])
		} else if (sym >= Keysyms["Ooblique"]) && (sym <= Keysyms["Thorn"]) {
			lower += (Keysyms["oslash"] - Keysyms["Ooblique"])
		} else if (sym >= Keysyms["oslash"]) && (sym <= Keysyms["thorn"]) {
			upper -= (Keysyms["oslash"] - Keysyms["Ooblique"])
		}
	case 1: /* Latin 2 */
		/* Assume the KeySym is a legal value (ignore discontinuities) */
		if sym == Keysyms["Aogonek"] {
			lower = Keysyms["aogonek"]
		} else if sym >= Keysyms["Lstroke"] && sym <= Keysyms["Sacute"] {
			lower += (Keysyms["lstroke"] - Keysyms["Lstroke"])
		} else if sym >= Keysyms["Scaron"] && sym <= Keysyms["Zacute"] {
			lower += (Keysyms["scaron"] - Keysyms["Scaron"])
		} else if sym >= Keysyms["Zcaron"] && sym <= Keysyms["Zabovedot"] {
			lower += (Keysyms["zcaron"] - Keysyms["Zcaron"])
		} else if sym == Keysyms["aogonek"] {
			upper = Keysyms["Aogonek"]
		} else if sym >= Keysyms["lstroke"] && sym <= Keysyms["sacute"] {
			upper -= (Keysyms["lstroke"] - Keysyms["Lstroke"])
		} else if sym >= Keysyms["scaron"] && sym <= Keysyms["zacute"] {
			upper -= (Keysyms["scaron"] - Keysyms["Scaron"])
		} else if sym >= Keysyms["zcaron"] && sym <= Keysyms["zabovedot"] {
			upper -= (Keysyms["zcaron"] - Keysyms["Zcaron"])
		} else if sym >= Keysyms["Racute"] && sym <= Keysyms["Tcedilla"] {
			lower += (Keysyms["racute"] - Keysyms["Racute"])
		} else if sym >= Keysyms["racute"] && sym <= Keysyms["tcedilla"] {
			upper -= (Keysyms["racute"] - Keysyms["Racute"])
		}
	case 2: /* Latin 3 */
		/* Assume the KeySym is a legal value (ignore discontinuities) */
		if sym >= Keysyms["Hstroke"] && sym <= Keysyms["Hcircumflex"] {
			lower += (Keysyms["hstroke"] - Keysyms["Hstroke"])
		} else if sym >= Keysyms["Gbreve"] && sym <= Keysyms["Jcircumflex"] {
			lower += (Keysyms["gbreve"] - Keysyms["Gbreve"])
		} else if sym >= Keysyms["hstroke"] && sym <= Keysyms["hcircumflex"] {
			upper -= (Keysyms["hstroke"] - Keysyms["Hstroke"])
		} else if sym >= Keysyms["gbreve"] && sym <= Keysyms["jcircumflex"] {
			upper -= (Keysyms["gbreve"] - Keysyms["Gbreve"])
		} else if sym >= Keysyms["Cabovedot"] && sym <= Keysyms["Scircumflex"] {
			lower += (Keysyms["cabovedot"] - Keysyms["Cabovedot"])
		} else if sym >= Keysyms["cabovedot"] && sym <= Keysyms["scircumflex"] {
			upper -= (Keysyms["cabovedot"] - Keysyms["Cabovedot"])
		}
	case 3: /* Latin 4 */
		/* Assume the KeySym is a legal value (ignore discontinuities) */
		if sym >= Keysyms["Rcedilla"] && sym <= Keysyms["Tslash"] {
			lower += (Keysyms["rcedilla"] - Keysyms["Rcedilla"])
		} else if sym >= Keysyms["rcedilla"] && sym <= Keysyms["tslash"] {
			upper -= (Keysyms["rcedilla"] - Keysyms["Rcedilla"])
		} else if sym == Keysyms["ENG"] {
			lower = Keysyms["eng"]
		} else if sym == Keysyms["eng"] {
			upper = Keysyms["ENG"]
		} else if sym >= Keysyms["Amacron"] && sym <= Keysyms["Umacron"] {
			lower += (Keysyms["amacron"] - Keysyms["Amacron"])
		} else if sym >= Keysyms["amacron"] && sym <= Keysyms["umacron"] {
			upper -= (Keysyms["amacron"] - Keysyms["Amacron"])
		}
	case 6: /* Cyrillic */
		/* Assume the KeySym is a legal value (ignore discontinuities) */
		if sym >= Keysyms["Serbian_DJE"] && sym <= Keysyms["Serbian_DZE"] {
			lower -= (Keysyms["Serbian_DJE"] - Keysyms["Serbian_dje"])
		} else if sym >= Keysyms["Serbian_dje"] && sym <= Keysyms["Serbian_dze"] {
			upper += (Keysyms["Serbian_DJE"] - Keysyms["Serbian_dje"])
		} else if sym >= Keysyms["Cyrillic_YU"] && sym <= Keysyms["Cyrillic_HARDSIGN"] {
			lower -= (Keysyms["Cyrillic_YU"] - Keysyms["Cyrillic_yu"])
		} else if sym >= Keysyms["Cyrillic_yu"] && sym <= Keysyms["Cyrillic_hardsign"] {
			upper += (Keysyms["Cyrillic_YU"] - Keysyms["Cyrillic_yu"])
		}
	case 7: /* Greek */
		/* Assume the KeySym is a legal value (ignore discontinuities) */
		if sym >= Keysyms["Greek_ALPHAaccent"] && sym <= Keysyms["Greek_OMEGAaccent"] {
			lower += (Keysyms["Greek_alphaaccent"] - Keysyms["Greek_ALPHAaccent"])
		} else if sym >= Keysyms["Greek_alphaaccent"] && sym <= Keysyms["Greek_omegaaccent"] &&
			sym != Keysyms["Greek_iotaaccentdieresis"] &&
			sym != Keysyms["Greek_upsilonaccentdieresis"] {
			upper -= (Keysyms["Greek_alphaaccent"] - Keysyms["Greek_ALPHAaccent"])
		} else if sym >= Keysyms["Greek_ALPHA"] && sym <= Keysyms["Greek_OMEGA"] {
			lower += (Keysyms["Greek_alpha"] - Keysyms["Greek_ALPHA"])
		} else if sym >= Keysyms["Greek_alpha"] && sym <= Keysyms["Greek_omega"] &&
			sym != Keysyms["Greek_finalsmallsigma"] {
			upper -= (Keysyms["Greek_alpha"] - Keysyms["Greek_ALPHA"])
		}
	case 0x14: /* Armenian */
		if sym >= Keysyms["Armenian_AYB"] && sym <= Keysyms["Armenian_fe"] {
			lower = sym | 1
			upper = sym & ^xproto.Keysym(1)
		}
	}

	return lower, upper
}
