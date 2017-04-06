package keybind

/*
This file contains the keysym definitions from X.
Taken from X11/keysymdef.h

It also contains the "XFree86 vendor specific keysyms" taken from
X11/XF86keysym.h.

We store this as a map because we need to be able to do reverse lookups.

keysyms is a mapping from english strings to key symbols.
strKeysyms is a mapping from key symbols to english strings.
*/

import "github.com/BurntSushi/xgb/xproto"

func init() {
	strKeysyms = make(map[xproto.Keysym]string, len(Keysyms))
	for kstr, keysym := range Keysyms {
		// If we already have this keysym as a key, skip it.
		// (Prefer the first. This may be bad.)
		if _, ok := strKeysyms[keysym]; !ok {
			strKeysyms[keysym] = kstr
		}
	}
}

// weirdKeysyms is a feeble attempt to map english words to single
// characters. i.e., "bracketleft" -> [ and "exclam" to !
var weirdKeysyms = map[string]rune{
	"space":        ' ',
	"exclam":       '!',
	"at":           '@',
	"numbersign":   '#',
	"dollar":       '$',
	"percent":      '%',
	"asciicircum":  '^',
	"ampersand":    '&',
	"asterisk":     '*',
	"parenleft":    '(',
	"parenright":   ')',
	"bracketleft":  '[',
	"bracketright": ']',
	"braceleft":    '{',
	"braceright":   '}',
	"minus":        '-',
	"underscore":   '_',
	"equal":        '=',
	"plus":         '+',
	"backslash":    '\\',
	"bar":          '|',
	"semicolon":    ';',
	"colon":        ':',
	"apostrophe":   '\'',
	"quoteright":   '\'',
	"quotedbl":     '"',
	"less":         '<',
	"greater":      '>',
	"comma":        ',',
	"period":       '.',
	"slash":        '/',
	"question":     '?',
	"grave":        '`',
	"quoteleft":    '`',
	"asciitilde":   '~',
	"KP_Multiply":  '*',
	"KP_Divide":    '/',
	"KP_Subtract":  '-',
	"KP_Add":       '+',
	"KP_Decimal":   '.',
	"KP_0":         '0',
	"KP_1":         '1',
	"KP_2":         '2',
	"KP_3":         '3',
	"KP_4":         '4',
	"KP_5":         '5',
	"KP_6":         '6',
	"KP_7":         '7',
	"KP_8":         '8',
	"KP_9":         '9',
}

// strKeysyms is the reverse of keysyms. It is built upon initialization.
// TODO: Hard code the reverse map to be faster.
var strKeysyms map[xproto.Keysym]string

var Keysyms map[string]xproto.Keysym = map[string]xproto.Keysym{
	"VoidSymbol":                  0xffffff,
	"BackSpace":                   0xff08,
	"Tab":                         0xff09,
	"Linefeed":                    0xff0a,
	"Clear":                       0xff0b,
	"Return":                      0xff0d,
	"Pause":                       0xff13,
	"Scroll_Lock":                 0xff14,
	"Sys_Req":                     0xff15,
	"Escape":                      0xff1b,
	"Delete":                      0xffff,
	"Multi_key":                   0xff20,
	"Codeinput":                   0xff37,
	"SingleCandidate":             0xff3c,
	"MultipleCandidate":           0xff3d,
	"PreviousCandidate":           0xff3e,
	"Kanji":                       0xff21,
	"Muhenkan":                    0xff22,
	"Henkan_Mode":                 0xff23,
	"Henkan":                      0xff23,
	"Romaji":                      0xff24,
	"Hiragana":                    0xff25,
	"Katakana":                    0xff26,
	"Hiragana_Katakana":           0xff27,
	"Zenkaku":                     0xff28,
	"Hankaku":                     0xff29,
	"Zenkaku_Hankaku":             0xff2a,
	"Touroku":                     0xff2b,
	"Massyo":                      0xff2c,
	"Kana_Lock":                   0xff2d,
	"Kana_Shift":                  0xff2e,
	"Eisu_Shift":                  0xff2f,
	"Eisu_toggle":                 0xff30,
	"Kanji_Bangou":                0xff37,
	"Zen_Koho":                    0xff3d,
	"Mae_Koho":                    0xff3e,
	"Home":                        0xff50,
	"Left":                        0xff51,
	"Up":                          0xff52,
	"Right":                       0xff53,
	"Down":                        0xff54,
	"Prior":                       0xff55,
	"Page_Up":                     0xff55,
	"Next":                        0xff56,
	"Page_Down":                   0xff56,
	"End":                         0xff57,
	"Begin":                       0xff58,
	"Select":                      0xff60,
	"Print":                       0xff61,
	"Execute":                     0xff62,
	"Insert":                      0xff63,
	"Undo":                        0xff65,
	"Redo":                        0xff66,
	"Menu":                        0xff67,
	"Find":                        0xff68,
	"Cancel":                      0xff69,
	"Help":                        0xff6a,
	"Break":                       0xff6b,
	"Mode_switch":                 0xff7e,
	"script_switch":               0xff7e,
	"Num_Lock":                    0xff7f,
	"KP_Space":                    0xff80,
	"KP_Tab":                      0xff89,
	"KP_Enter":                    0xff8d,
	"KP_F1":                       0xff91,
	"KP_F2":                       0xff92,
	"KP_F3":                       0xff93,
	"KP_F4":                       0xff94,
	"KP_Home":                     0xff95,
	"KP_Left":                     0xff96,
	"KP_Up":                       0xff97,
	"KP_Right":                    0xff98,
	"KP_Down":                     0xff99,
	"KP_Prior":                    0xff9a,
	"KP_Page_Up":                  0xff9a,
	"KP_Next":                     0xff9b,
	"KP_Page_Down":                0xff9b,
	"KP_End":                      0xff9c,
	"KP_Begin":                    0xff9d,
	"KP_Insert":                   0xff9e,
	"KP_Delete":                   0xff9f,
	"KP_Equal":                    0xffbd,
	"KP_Multiply":                 0xffaa,
	"KP_Add":                      0xffab,
	"KP_Separator":                0xffac,
	"KP_Subtract":                 0xffad,
	"KP_Decimal":                  0xffae,
	"KP_Divide":                   0xffaf,
	"KP_0":                        0xffb0,
	"KP_1":                        0xffb1,
	"KP_2":                        0xffb2,
	"KP_3":                        0xffb3,
	"KP_4":                        0xffb4,
	"KP_5":                        0xffb5,
	"KP_6":                        0xffb6,
	"KP_7":                        0xffb7,
	"KP_8":                        0xffb8,
	"KP_9":                        0xffb9,
	"F1":                          0xffbe,
	"F2":                          0xffbf,
	"F3":                          0xffc0,
	"F4":                          0xffc1,
	"F5":                          0xffc2,
	"F6":                          0xffc3,
	"F7":                          0xffc4,
	"F8":                          0xffc5,
	"F9":                          0xffc6,
	"F10":                         0xffc7,
	"F11":                         0xffc8,
	"L1":                          0xffc8,
	"F12":                         0xffc9,
	"L2":                          0xffc9,
	"F13":                         0xffca,
	"L3":                          0xffca,
	"F14":                         0xffcb,
	"L4":                          0xffcb,
	"F15":                         0xffcc,
	"L5":                          0xffcc,
	"F16":                         0xffcd,
	"L6":                          0xffcd,
	"F17":                         0xffce,
	"L7":                          0xffce,
	"F18":                         0xffcf,
	"L8":                          0xffcf,
	"F19":                         0xffd0,
	"L9":                          0xffd0,
	"F20":                         0xffd1,
	"L10":                         0xffd1,
	"F21":                         0xffd2,
	"R1":                          0xffd2,
	"F22":                         0xffd3,
	"R2":                          0xffd3,
	"F23":                         0xffd4,
	"R3":                          0xffd4,
	"F24":                         0xffd5,
	"R4":                          0xffd5,
	"F25":                         0xffd6,
	"R5":                          0xffd6,
	"F26":                         0xffd7,
	"R6":                          0xffd7,
	"F27":                         0xffd8,
	"R7":                          0xffd8,
	"F28":                         0xffd9,
	"R8":                          0xffd9,
	"F29":                         0xffda,
	"R9":                          0xffda,
	"F30":                         0xffdb,
	"R10":                         0xffdb,
	"F31":                         0xffdc,
	"R11":                         0xffdc,
	"F32":                         0xffdd,
	"R12":                         0xffdd,
	"F33":                         0xffde,
	"R13":                         0xffde,
	"F34":                         0xffdf,
	"R14":                         0xffdf,
	"F35":                         0xffe0,
	"R15":                         0xffe0,
	"Shift_L":                     0xffe1,
	"Shift_R":                     0xffe2,
	"Control_L":                   0xffe3,
	"Control_R":                   0xffe4,
	"Caps_Lock":                   0xffe5,
	"Shift_Lock":                  0xffe6,
	"Meta_L":                      0xffe7,
	"Meta_R":                      0xffe8,
	"Alt_L":                       0xffe9,
	"Alt_R":                       0xffea,
	"Super_L":                     0xffeb,
	"Super_R":                     0xffec,
	"Hyper_L":                     0xffed,
	"Hyper_R":                     0xffee,
	"ISO_Lock":                    0xfe01,
	"ISO_Level2_Latch":            0xfe02,
	"ISO_Level3_Shift":            0xfe03,
	"ISO_Level3_Latch":            0xfe04,
	"ISO_Level3_Lock":             0xfe05,
	"ISO_Level5_Shift":            0xfe11,
	"ISO_Level5_Latch":            0xfe12,
	"ISO_Level5_Lock":             0xfe13,
	"ISO_Group_Shift":             0xff7e,
	"ISO_Group_Latch":             0xfe06,
	"ISO_Group_Lock":              0xfe07,
	"ISO_Next_Group":              0xfe08,
	"ISO_Next_Group_Lock":         0xfe09,
	"ISO_Prev_Group":              0xfe0a,
	"ISO_Prev_Group_Lock":         0xfe0b,
	"ISO_First_Group":             0xfe0c,
	"ISO_First_Group_Lock":        0xfe0d,
	"ISO_Last_Group":              0xfe0e,
	"ISO_Last_Group_Lock":         0xfe0f,
	"ISO_Left_Tab":                0xfe20,
	"ISO_Move_Line_Up":            0xfe21,
	"ISO_Move_Line_Down":          0xfe22,
	"ISO_Partial_Line_Up":         0xfe23,
	"ISO_Partial_Line_Down":       0xfe24,
	"ISO_Partial_Space_Left":      0xfe25,
	"ISO_Partial_Space_Right":     0xfe26,
	"ISO_Set_Margin_Left":         0xfe27,
	"ISO_Set_Margin_Right":        0xfe28,
	"ISO_Release_Margin_Left":     0xfe29,
	"ISO_Release_Margin_Right":    0xfe2a,
	"ISO_Release_Both_Margins":    0xfe2b,
	"ISO_Fast_Cursor_Left":        0xfe2c,
	"ISO_Fast_Cursor_Right":       0xfe2d,
	"ISO_Fast_Cursor_Up":          0xfe2e,
	"ISO_Fast_Cursor_Down":        0xfe2f,
	"ISO_Continuous_Underline":    0xfe30,
	"ISO_Discontinuous_Underline": 0xfe31,
	"ISO_Emphasize":               0xfe32,
	"ISO_Center_Object":           0xfe33,
	"ISO_Enter":                   0xfe34,
	"dead_grave":                  0xfe50,
	"dead_acute":                  0xfe51,
	"dead_circumflex":             0xfe52,
	"dead_tilde":                  0xfe53,
	"dead_perispomeni":            0xfe53,
	"dead_macron":                 0xfe54,
	"dead_breve":                  0xfe55,
	"dead_abovedot":               0xfe56,
	"dead_diaeresis":              0xfe57,
	"dead_abovering":              0xfe58,
	"dead_doubleacute":            0xfe59,
	"dead_caron":                  0xfe5a,
	"dead_cedilla":                0xfe5b,
	"dead_ogonek":                 0xfe5c,
	"dead_iota":                   0xfe5d,
	"dead_voiced_sound":           0xfe5e,
	"dead_semivoiced_sound":       0xfe5f,
	"dead_belowdot":               0xfe60,
	"dead_hook":                   0xfe61,
	"dead_horn":                   0xfe62,
	"dead_stroke":                 0xfe63,
	"dead_abovecomma":             0xfe64,
	"dead_psili":                  0xfe64,
	"dead_abovereversedcomma":     0xfe65,
	"dead_dasia":                  0xfe65,
	"dead_doublegrave":            0xfe66,
	"dead_belowring":              0xfe67,
	"dead_belowmacron":            0xfe68,
	"dead_belowcircumflex":        0xfe69,
	"dead_belowtilde":             0xfe6a,
	"dead_belowbreve":             0xfe6b,
	"dead_belowdiaeresis":         0xfe6c,
	"dead_invertedbreve":          0xfe6d,
	"dead_belowcomma":             0xfe6e,
	"dead_currency":               0xfe6f,
	"dead_a":                      0xfe80,
	"dead_A":                      0xfe81,
	"dead_e":                      0xfe82,
	"dead_E":                      0xfe83,
	"dead_i":                      0xfe84,
	"dead_I":                      0xfe85,
	"dead_o":                      0xfe86,
	"dead_O":                      0xfe87,
	"dead_u":                      0xfe88,
	"dead_U":                      0xfe89,
	"dead_small_schwa":            0xfe8a,
	"dead_capital_schwa":          0xfe8b,
	"First_Virtual_Screen":        0xfed0,
	"Prev_Virtual_Screen":         0xfed1,
	"Next_Virtual_Screen":         0xfed2,
	"Last_Virtual_Screen":         0xfed4,
	"Terminate_Server":            0xfed5,
	"AccessX_Enable":              0xfe70,
	"AccessX_Feedback_Enable":     0xfe71,
	"RepeatKeys_Enable":           0xfe72,
	"SlowKeys_Enable":             0xfe73,
	"BounceKeys_Enable":           0xfe74,
	"StickyKeys_Enable":           0xfe75,
	"MouseKeys_Enable":            0xfe76,
	"MouseKeys_Accel_Enable":      0xfe77,
	"Overlay1_Enable":             0xfe78,
	"Overlay2_Enable":             0xfe79,
	"AudibleBell_Enable":          0xfe7a,
	"Pointer_Left":                0xfee0,
	"Pointer_Right":               0xfee1,
	"Pointer_Up":                  0xfee2,
	"Pointer_Down":                0xfee3,
	"Pointer_UpLeft":              0xfee4,
	"Pointer_UpRight":             0xfee5,
	"Pointer_DownLeft":            0xfee6,
	"Pointer_DownRight":           0xfee7,
	"Pointer_Button_Dflt":         0xfee8,
	"Pointer_Button1":             0xfee9,
	"Pointer_Button2":             0xfeea,
	"Pointer_Button3":             0xfeeb,
	"Pointer_Button4":             0xfeec,
	"Pointer_Button5":             0xfeed,
	"Pointer_DblClick_Dflt":       0xfeee,
	"Pointer_DblClick1":           0xfeef,
	"Pointer_DblClick2":           0xfef0,
	"Pointer_DblClick3":           0xfef1,
	"Pointer_DblClick4":           0xfef2,
	"Pointer_DblClick5":           0xfef3,
	"Pointer_Drag_Dflt":           0xfef4,
	"Pointer_Drag1":               0xfef5,
	"Pointer_Drag2":               0xfef6,
	"Pointer_Drag3":               0xfef7,
	"Pointer_Drag4":               0xfef8,
	"Pointer_Drag5":               0xfefd,
	"Pointer_EnableKeys":          0xfef9,
	"Pointer_Accelerate":          0xfefa,
	"Pointer_DfltBtnNext":         0xfefb,
	"Pointer_DfltBtnPrev":         0xfefc,
	"3270_Duplicate":              0xfd01,
	"3270_FieldMark":              0xfd02,
	"3270_Right2":                 0xfd03,
	"3270_Left2":                  0xfd04,
	"3270_BackTab":                0xfd05,
	"3270_EraseEOF":               0xfd06,
	"3270_EraseInput":             0xfd07,
	"3270_Reset":                  0xfd08,
	"3270_Quit":                   0xfd09,
	"3270_PA1":                    0xfd0a,
	"3270_PA2":                    0xfd0b,
	"3270_PA3":                    0xfd0c,
	"3270_Test":                   0xfd0d,
	"3270_Attn":                   0xfd0e,
	"3270_CursorBlink":            0xfd0f,
	"3270_AltCursor":              0xfd10,
	"3270_KeyClick":               0xfd11,
	"3270_Jump":                   0xfd12,
	"3270_Ident":                  0xfd13,
	"3270_Rule":                   0xfd14,
	"3270_Copy":                   0xfd15,
	"3270_Play":                   0xfd16,
	"3270_Setup":                  0xfd17,
	"3270_Record":                 0xfd18,
	"3270_ChangeScreen":           0xfd19,
	"3270_DeleteWord":             0xfd1a,
	"3270_ExSelect":               0xfd1b,
	"3270_CursorSelect":           0xfd1c,
	"3270_PrintScreen":            0xfd1d,
	"3270_Enter":                  0xfd1e,
	"space":                       0x0020,
	"exclam":                      0x0021,
	"quotedbl":                    0x0022,
	"numbersign":                  0x0023,
	"dollar":                      0x0024,
	"percent":                     0x0025,
	"ampersand":                   0x0026,
	"apostrophe":                  0x0027,
	"quoteright":                  0x0027,
	"parenleft":                   0x0028,
	"parenright":                  0x0029,
	"asterisk":                    0x002a,
	"plus":                        0x002b,
	"comma":                       0x002c,
	"minus":                       0x002d,
	"period":                      0x002e,
	"slash":                       0x002f,
	"0":                           0x0030,
	"1":                           0x0031,
	"2":                           0x0032,
	"3":                           0x0033,
	"4":                           0x0034,
	"5":                           0x0035,
	"6":                           0x0036,
	"7":                           0x0037,
	"8":                           0x0038,
	"9":                           0x0039,
	"colon":                       0x003a,
	"semicolon":                   0x003b,
	"less":                        0x003c,
	"equal":                       0x003d,
	"greater":                     0x003e,
	"question":                    0x003f,
	"at":                          0x0040,
	"A":                           0x0041,
	"B":                           0x0042,
	"C":                           0x0043,
	"D":                           0x0044,
	"E":                           0x0045,
	"F":                           0x0046,
	"G":                           0x0047,
	"H":                           0x0048,
	"I":                           0x0049,
	"J":                           0x004a,
	"K":                           0x004b,
	"L":                           0x004c,
	"M":                           0x004d,
	"N":                           0x004e,
	"O":                           0x004f,
	"P":                           0x0050,
	"Q":                           0x0051,
	"R":                           0x0052,
	"S":                           0x0053,
	"T":                           0x0054,
	"U":                           0x0055,
	"V":                           0x0056,
	"W":                           0x0057,
	"X":                           0x0058,
	"Y":                           0x0059,
	"Z":                           0x005a,
	"bracketleft":                 0x005b,
	"backslash":                   0x005c,
	"bracketright":                0x005d,
	"asciicircum":                 0x005e,
	"underscore":                  0x005f,
	"grave":                       0x0060,
	"quoteleft":                   0x0060,
	"a":                           0x0061,
	"b":                           0x0062,
	"c":                           0x0063,
	"d":                           0x0064,
	"e":                           0x0065,
	"f":                           0x0066,
	"g":                           0x0067,
	"h":                           0x0068,
	"i":                           0x0069,
	"j":                           0x006a,
	"k":                           0x006b,
	"l":                           0x006c,
	"m":                           0x006d,
	"n":                           0x006e,
	"o":                           0x006f,
	"p":                           0x0070,
	"q":                           0x0071,
	"r":                           0x0072,
	"s":                           0x0073,
	"t":                           0x0074,
	"u":                           0x0075,
	"v":                           0x0076,
	"w":                           0x0077,
	"x":                           0x0078,
	"y":                           0x0079,
	"z":                           0x007a,
	"braceleft":                   0x007b,
	"bar":                         0x007c,
	"braceright":                  0x007d,
	"asciitilde":                  0x007e,
	"nobreakspace":                0x00a0,
	"exclamdown":                  0x00a1,
	"cent":                        0x00a2,
	"sterling":                    0x00a3,
	"currency":                    0x00a4,
	"yen":                         0x00a5,
	"brokenbar":                   0x00a6,
	"section":                     0x00a7,
	"diaeresis":                   0x00a8,
	"copyright":                   0x00a9,
	"ordfeminine":                 0x00aa,
	"guillemotleft":               0x00ab,
	"notsign":                     0x00ac,
	"hyphen":                      0x00ad,
	"registered":                  0x00ae,
	"macron":                      0x00af,
	"degree":                      0x00b0,
	"plusminus":                   0x00b1,
	"twosuperior":                 0x00b2,
	"threesuperior":               0x00b3,
	"acute":                       0x00b4,
	"mu":                          0x00b5,
	"paragraph":                   0x00b6,
	"periodcentered":              0x00b7,
	"cedilla":                     0x00b8,
	"onesuperior":                 0x00b9,
	"masculine":                   0x00ba,
	"guillemotright":              0x00bb,
	"onequarter":                  0x00bc,
	"onehalf":                     0x00bd,
	"threequarters":               0x00be,
	"questiondown":                0x00bf,
	"Agrave":                      0x00c0,
	"Aacute":                      0x00c1,
	"Acircumflex":                 0x00c2,
	"Atilde":                      0x00c3,
	"Adiaeresis":                  0x00c4,
	"Aring":                       0x00c5,
	"AE":                          0x00c6,
	"Ccedilla":                    0x00c7,
	"Egrave":                      0x00c8,
	"Eacute":                      0x00c9,
	"Ecircumflex":                 0x00ca,
	"Ediaeresis":                  0x00cb,
	"Igrave":                      0x00cc,
	"Iacute":                      0x00cd,
	"Icircumflex":                 0x00ce,
	"Idiaeresis":                  0x00cf,
	"ETH":                         0x00d0,
	"Eth":                         0x00d0,
	"Ntilde":                      0x00d1,
	"Ograve":                      0x00d2,
	"Oacute":                      0x00d3,
	"Ocircumflex":                 0x00d4,
	"Otilde":                      0x00d5,
	"Odiaeresis":                  0x00d6,
	"multiply":                    0x00d7,
	"Oslash":                      0x00d8,
	"Ooblique":                    0x00d8,
	"Ugrave":                      0x00d9,
	"Uacute":                      0x00da,
	"Ucircumflex":                 0x00db,
	"Udiaeresis":                  0x00dc,
	"Yacute":                      0x00dd,
	"THORN":                       0x00de,
	"Thorn":                       0x00de,
	"ssharp":                      0x00df,
	"agrave":                      0x00e0,
	"aacute":                      0x00e1,
	"acircumflex":                 0x00e2,
	"atilde":                      0x00e3,
	"adiaeresis":                  0x00e4,
	"aring":                       0x00e5,
	"ae":                          0x00e6,
	"ccedilla":                    0x00e7,
	"egrave":                      0x00e8,
	"eacute":                      0x00e9,
	"ecircumflex":                 0x00ea,
	"ediaeresis":                  0x00eb,
	"igrave":                      0x00ec,
	"iacute":                      0x00ed,
	"icircumflex":                 0x00ee,
	"idiaeresis":                  0x00ef,
	"eth":                         0x00f0,
	"ntilde":                      0x00f1,
	"ograve":                      0x00f2,
	"oacute":                      0x00f3,
	"ocircumflex":                 0x00f4,
	"otilde":                      0x00f5,
	"odiaeresis":                  0x00f6,
	"division":                    0x00f7,
	"oslash":                      0x00f8,
	"ooblique":                    0x00f8,
	"ugrave":                      0x00f9,
	"uacute":                      0x00fa,
	"ucircumflex":                 0x00fb,
	"udiaeresis":                  0x00fc,
	"yacute":                      0x00fd,
	"thorn":                       0x00fe,
	"ydiaeresis":                  0x00ff,
	"Aogonek":                     0x01a1,
	"breve":                       0x01a2,
	"Lstroke":                     0x01a3,
	"Lcaron":                      0x01a5,
	"Sacute":                      0x01a6,
	"Scaron":                      0x01a9,
	"Scedilla":                    0x01aa,
	"Tcaron":                      0x01ab,
	"Zacute":                      0x01ac,
	"Zcaron":                      0x01ae,
	"Zabovedot":                   0x01af,
	"aogonek":                     0x01b1,
	"ogonek":                      0x01b2,
	"lstroke":                     0x01b3,
	"lcaron":                      0x01b5,
	"sacute":                      0x01b6,
	"caron":                       0x01b7,
	"scaron":                      0x01b9,
	"scedilla":                    0x01ba,
	"tcaron":                      0x01bb,
	"zacute":                      0x01bc,
	"doubleacute":                 0x01bd,
	"zcaron":                      0x01be,
	"zabovedot":                   0x01bf,
	"Racute":                      0x01c0,
	"Abreve":                      0x01c3,
	"Lacute":                      0x01c5,
	"Cacute":                      0x01c6,
	"Ccaron":                      0x01c8,
	"Eogonek":                     0x01ca,
	"Ecaron":                      0x01cc,
	"Dcaron":                      0x01cf,
	"Dstroke":                     0x01d0,
	"Nacute":                      0x01d1,
	"Ncaron":                      0x01d2,
	"Odoubleacute":                0x01d5,
	"Rcaron":                      0x01d8,
	"Uring":                       0x01d9,
	"Udoubleacute":                0x01db,
	"Tcedilla":                    0x01de,
	"racute":                      0x01e0,
	"abreve":                      0x01e3,
	"lacute":                      0x01e5,
	"cacute":                      0x01e6,
	"ccaron":                      0x01e8,
	"eogonek":                     0x01ea,
	"ecaron":                      0x01ec,
	"dcaron":                      0x01ef,
	"dstroke":                     0x01f0,
	"nacute":                      0x01f1,
	"ncaron":                      0x01f2,
	"odoubleacute":                0x01f5,
	"udoubleacute":                0x01fb,
	"rcaron":                      0x01f8,
	"uring":                       0x01f9,
	"tcedilla":                    0x01fe,
	"abovedot":                    0x01ff,
	"Hstroke":                     0x02a1,
	"Hcircumflex":                 0x02a6,
	"Iabovedot":                   0x02a9,
	"Gbreve":                      0x02ab,
	"Jcircumflex":                 0x02ac,
	"hstroke":                     0x02b1,
	"hcircumflex":                 0x02b6,
	"idotless":                    0x02b9,
	"gbreve":                      0x02bb,
	"jcircumflex":                 0x02bc,
	"Cabovedot":                   0x02c5,
	"Ccircumflex":                 0x02c6,
	"Gabovedot":                   0x02d5,
	"Gcircumflex":                 0x02d8,
	"Ubreve":                      0x02dd,
	"Scircumflex":                 0x02de,
	"cabovedot":                   0x02e5,
	"ccircumflex":                 0x02e6,
	"gabovedot":                   0x02f5,
	"gcircumflex":                 0x02f8,
	"ubreve":                      0x02fd,
	"scircumflex":                 0x02fe,
	"kra":                         0x03a2,
	"kappa":                       0x03a2,
	"Rcedilla":                    0x03a3,
	"Itilde":                      0x03a5,
	"Lcedilla":                    0x03a6,
	"Emacron":                     0x03aa,
	"Gcedilla":                    0x03ab,
	"Tslash":                      0x03ac,
	"rcedilla":                    0x03b3,
	"itilde":                      0x03b5,
	"lcedilla":                    0x03b6,
	"emacron":                     0x03ba,
	"gcedilla":                    0x03bb,
	"tslash":                      0x03bc,
	"ENG":                         0x03bd,
	"eng":                         0x03bf,
	"Amacron":                     0x03c0,
	"Iogonek":                     0x03c7,
	"Eabovedot":                   0x03cc,
	"Imacron":                     0x03cf,
	"Ncedilla":                    0x03d1,
	"Omacron":                     0x03d2,
	"Kcedilla":                    0x03d3,
	"Uogonek":                     0x03d9,
	"Utilde":                      0x03dd,
	"Umacron":                     0x03de,
	"amacron":                     0x03e0,
	"iogonek":                     0x03e7,
	"eabovedot":                   0x03ec,
	"imacron":                     0x03ef,
	"ncedilla":                    0x03f1,
	"omacron":                     0x03f2,
	"kcedilla":                    0x03f3,
	"uogonek":                     0x03f9,
	"utilde":                      0x03fd,
	"umacron":                     0x03fe,
	"Babovedot":                   0x1001e02,
	"babovedot":                   0x1001e03,
	"Dabovedot":                   0x1001e0a,
	"Wgrave":                      0x1001e80,
	"Wacute":                      0x1001e82,
	"dabovedot":                   0x1001e0b,
	"Ygrave":                      0x1001ef2,
	"Fabovedot":                   0x1001e1e,
	"fabovedot":                   0x1001e1f,
	"Mabovedot":                   0x1001e40,
	"mabovedot":                   0x1001e41,
	"Pabovedot":                   0x1001e56,
	"wgrave":                      0x1001e81,
	"pabovedot":                   0x1001e57,
	"wacute":                      0x1001e83,
	"Sabovedot":                   0x1001e60,
	"ygrave":                      0x1001ef3,
	"Wdiaeresis":                  0x1001e84,
	"wdiaeresis":                  0x1001e85,
	"sabovedot":                   0x1001e61,
	"Wcircumflex":                 0x1000174,
	"Tabovedot":                   0x1001e6a,
	"Ycircumflex":                 0x1000176,
	"wcircumflex":                 0x1000175,
	"tabovedot":                   0x1001e6b,
	"ycircumflex":                 0x1000177,
	"OE":                          0x13bc,
	"oe":                          0x13bd,
	"Ydiaeresis":                  0x13be,
	"overline":                    0x047e,
	"kana_fullstop":               0x04a1,
	"kana_openingbracket":         0x04a2,
	"kana_closingbracket":         0x04a3,
	"kana_comma":                  0x04a4,
	"kana_conjunctive":            0x04a5,
	"kana_middledot":              0x04a5,
	"kana_WO":                     0x04a6,
	"kana_a":                      0x04a7,
	"kana_i":                      0x04a8,
	"kana_u":                      0x04a9,
	"kana_e":                      0x04aa,
	"kana_o":                      0x04ab,
	"kana_ya":                     0x04ac,
	"kana_yu":                     0x04ad,
	"kana_yo":                     0x04ae,
	"kana_tsu":                    0x04af,
	"kana_tu":                     0x04af,
	"prolongedsound":              0x04b0,
	"kana_A":                      0x04b1,
	"kana_I":                      0x04b2,
	"kana_U":                      0x04b3,
	"kana_E":                      0x04b4,
	"kana_O":                      0x04b5,
	"kana_KA":                     0x04b6,
	"kana_KI":                     0x04b7,
	"kana_KU":                     0x04b8,
	"kana_KE":                     0x04b9,
	"kana_KO":                     0x04ba,
	"kana_SA":                     0x04bb,
	"kana_SHI":                    0x04bc,
	"kana_SU":                     0x04bd,
	"kana_SE":                     0x04be,
	"kana_SO":                     0x04bf,
	"kana_TA":                     0x04c0,
	"kana_CHI":                    0x04c1,
	"kana_TI":                     0x04c1,
	"kana_TSU":                    0x04c2,
	"kana_TU":                     0x04c2,
	"kana_TE":                     0x04c3,
	"kana_TO":                     0x04c4,
	"kana_NA":                     0x04c5,
	"kana_NI":                     0x04c6,
	"kana_NU":                     0x04c7,
	"kana_NE":                     0x04c8,
	"kana_NO":                     0x04c9,
	"kana_HA":                     0x04ca,
	"kana_HI":                     0x04cb,
	"kana_FU":                     0x04cc,
	"kana_HU":                     0x04cc,
	"kana_HE":                     0x04cd,
	"kana_HO":                     0x04ce,
	"kana_MA":                     0x04cf,
	"kana_MI":                     0x04d0,
	"kana_MU":                     0x04d1,
	"kana_ME":                     0x04d2,
	"kana_MO":                     0x04d3,
	"kana_YA":                     0x04d4,
	"kana_YU":                     0x04d5,
	"kana_YO":                     0x04d6,
	"kana_RA":                     0x04d7,
	"kana_RI":                     0x04d8,
	"kana_RU":                     0x04d9,
	"kana_RE":                     0x04da,
	"kana_RO":                     0x04db,
	"kana_WA":                     0x04dc,
	"kana_N":                      0x04dd,
	"voicedsound":                 0x04de,
	"semivoicedsound":             0x04df,
	"kana_switch":                 0xff7e,
	"Farsi_0":                     0x10006f0,
	"Farsi_1":                     0x10006f1,
	"Farsi_2":                     0x10006f2,
	"Farsi_3":                     0x10006f3,
	"Farsi_4":                     0x10006f4,
	"Farsi_5":                     0x10006f5,
	"Farsi_6":                     0x10006f6,
	"Farsi_7":                     0x10006f7,
	"Farsi_8":                     0x10006f8,
	"Farsi_9":                     0x10006f9,
	"Arabic_percent":              0x100066a,
	"Arabic_superscript_alef":     0x1000670,
	"Arabic_tteh":                 0x1000679,
	"Arabic_peh":                  0x100067e,
	"Arabic_tcheh":                0x1000686,
	"Arabic_ddal":                 0x1000688,
	"Arabic_rreh":                 0x1000691,
	"Arabic_comma":                0x05ac,
	"Arabic_fullstop":             0x10006d4,
	"Arabic_0":                    0x1000660,
	"Arabic_1":                    0x1000661,
	"Arabic_2":                    0x1000662,
	"Arabic_3":                    0x1000663,
	"Arabic_4":                    0x1000664,
	"Arabic_5":                    0x1000665,
	"Arabic_6":                    0x1000666,
	"Arabic_7":                    0x1000667,
	"Arabic_8":                    0x1000668,
	"Arabic_9":                    0x1000669,
	"Arabic_semicolon":            0x05bb,
	"Arabic_question_mark":        0x05bf,
	"Arabic_hamza":                0x05c1,
	"Arabic_maddaonalef":          0x05c2,
	"Arabic_hamzaonalef":          0x05c3,
	"Arabic_hamzaonwaw":           0x05c4,
	"Arabic_hamzaunderalef":       0x05c5,
	"Arabic_hamzaonyeh":           0x05c6,
	"Arabic_alef":                 0x05c7,
	"Arabic_beh":                  0x05c8,
	"Arabic_tehmarbuta":           0x05c9,
	"Arabic_teh":                  0x05ca,
	"Arabic_theh":                 0x05cb,
	"Arabic_jeem":                 0x05cc,
	"Arabic_hah":                  0x05cd,
	"Arabic_khah":                 0x05ce,
	"Arabic_dal":                  0x05cf,
	"Arabic_thal":                 0x05d0,
	"Arabic_ra":                   0x05d1,
	"Arabic_zain":                 0x05d2,
	"Arabic_seen":                 0x05d3,
	"Arabic_sheen":                0x05d4,
	"Arabic_sad":                  0x05d5,
	"Arabic_dad":                  0x05d6,
	"Arabic_tah":                  0x05d7,
	"Arabic_zah":                  0x05d8,
	"Arabic_ain":                  0x05d9,
	"Arabic_ghain":                0x05da,
	"Arabic_tatweel":              0x05e0,
	"Arabic_feh":                  0x05e1,
	"Arabic_qaf":                  0x05e2,
	"Arabic_kaf":                  0x05e3,
	"Arabic_lam":                  0x05e4,
	"Arabic_meem":                 0x05e5,
	"Arabic_noon":                 0x05e6,
	"Arabic_ha":                   0x05e7,
	"Arabic_heh":                  0x05e7,
	"Arabic_waw":                  0x05e8,
	"Arabic_alefmaksura":          0x05e9,
	"Arabic_yeh":                  0x05ea,
	"Arabic_fathatan":             0x05eb,
	"Arabic_dammatan":             0x05ec,
	"Arabic_kasratan":             0x05ed,
	"Arabic_fatha":                0x05ee,
	"Arabic_damma":                0x05ef,
	"Arabic_kasra":                0x05f0,
	"Arabic_shadda":               0x05f1,
	"Arabic_sukun":                0x05f2,
	"Arabic_madda_above":          0x1000653,
	"Arabic_hamza_above":          0x1000654,
	"Arabic_hamza_below":          0x1000655,
	"Arabic_jeh":                  0x1000698,
	"Arabic_veh":                  0x10006a4,
	"Arabic_keheh":                0x10006a9,
	"Arabic_gaf":                  0x10006af,
	"Arabic_noon_ghunna":          0x10006ba,
	"Arabic_heh_doachashmee":      0x10006be,
	"Farsi_yeh":                   0x10006cc,
	"Arabic_farsi_yeh":            0x10006cc,
	"Arabic_yeh_baree":            0x10006d2,
	"Arabic_heh_goal":             0x10006c1,
	"Arabic_switch":               0xff7e,
	"Cyrillic_GHE_bar":            0x1000492,
	"Cyrillic_ghe_bar":            0x1000493,
	"Cyrillic_ZHE_descender":      0x1000496,
	"Cyrillic_zhe_descender":      0x1000497,
	"Cyrillic_KA_descender":       0x100049a,
	"Cyrillic_ka_descender":       0x100049b,
	"Cyrillic_KA_vertstroke":      0x100049c,
	"Cyrillic_ka_vertstroke":      0x100049d,
	"Cyrillic_EN_descender":       0x10004a2,
	"Cyrillic_en_descender":       0x10004a3,
	"Cyrillic_U_straight":         0x10004ae,
	"Cyrillic_u_straight":         0x10004af,
	"Cyrillic_U_straight_bar":     0x10004b0,
	"Cyrillic_u_straight_bar":     0x10004b1,
	"Cyrillic_HA_descender":       0x10004b2,
	"Cyrillic_ha_descender":       0x10004b3,
	"Cyrillic_CHE_descender":      0x10004b6,
	"Cyrillic_che_descender":      0x10004b7,
	"Cyrillic_CHE_vertstroke":     0x10004b8,
	"Cyrillic_che_vertstroke":     0x10004b9,
	"Cyrillic_SHHA":               0x10004ba,
	"Cyrillic_shha":               0x10004bb,
	"Cyrillic_SCHWA":              0x10004d8,
	"Cyrillic_schwa":              0x10004d9,
	"Cyrillic_I_macron":           0x10004e2,
	"Cyrillic_i_macron":           0x10004e3,
	"Cyrillic_O_bar":              0x10004e8,
	"Cyrillic_o_bar":              0x10004e9,
	"Cyrillic_U_macron":           0x10004ee,
	"Cyrillic_u_macron":           0x10004ef,
	"Serbian_dje":                 0x06a1,
	"Macedonia_gje":               0x06a2,
	"Cyrillic_io":                 0x06a3,
	"Ukrainian_ie":                0x06a4,
	"Ukranian_je":                 0x06a4,
	"Macedonia_dse":               0x06a5,
	"Ukrainian_i":                 0x06a6,
	"Ukranian_i":                  0x06a6,
	"Ukrainian_yi":                0x06a7,
	"Ukranian_yi":                 0x06a7,
	"Cyrillic_je":                 0x06a8,
	"Serbian_je":                  0x06a8,
	"Cyrillic_lje":                0x06a9,
	"Serbian_lje":                 0x06a9,
	"Cyrillic_nje":                0x06aa,
	"Serbian_nje":                 0x06aa,
	"Serbian_tshe":                0x06ab,
	"Macedonia_kje":               0x06ac,
	"Ukrainian_ghe_with_upturn":   0x06ad,
	"Byelorussian_shortu":         0x06ae,
	"Cyrillic_dzhe":               0x06af,
	"Serbian_dze":                 0x06af,
	"numerosign":                  0x06b0,
	"Serbian_DJE":                 0x06b1,
	"Macedonia_GJE":               0x06b2,
	"Cyrillic_IO":                 0x06b3,
	"Ukrainian_IE":                0x06b4,
	"Ukranian_JE":                 0x06b4,
	"Macedonia_DSE":               0x06b5,
	"Ukrainian_I":                 0x06b6,
	"Ukranian_I":                  0x06b6,
	"Ukrainian_YI":                0x06b7,
	"Ukranian_YI":                 0x06b7,
	"Cyrillic_JE":                 0x06b8,
	"Serbian_JE":                  0x06b8,
	"Cyrillic_LJE":                0x06b9,
	"Serbian_LJE":                 0x06b9,
	"Cyrillic_NJE":                0x06ba,
	"Serbian_NJE":                 0x06ba,
	"Serbian_TSHE":                0x06bb,
	"Macedonia_KJE":               0x06bc,
	"Ukrainian_GHE_WITH_UPTURN":   0x06bd,
	"Byelorussian_SHORTU":         0x06be,
	"Cyrillic_DZHE":               0x06bf,
	"Serbian_DZE":                 0x06bf,
	"Cyrillic_yu":                 0x06c0,
	"Cyrillic_a":                  0x06c1,
	"Cyrillic_be":                 0x06c2,
	"Cyrillic_tse":                0x06c3,
	"Cyrillic_de":                 0x06c4,
	"Cyrillic_ie":                 0x06c5,
	"Cyrillic_ef":                 0x06c6,
	"Cyrillic_ghe":                0x06c7,
	"Cyrillic_ha":                 0x06c8,
	"Cyrillic_i":                  0x06c9,
	"Cyrillic_shorti":             0x06ca,
	"Cyrillic_ka":                 0x06cb,
	"Cyrillic_el":                 0x06cc,
	"Cyrillic_em":                 0x06cd,
	"Cyrillic_en":                 0x06ce,
	"Cyrillic_o":                  0x06cf,
	"Cyrillic_pe":                 0x06d0,
	"Cyrillic_ya":                 0x06d1,
	"Cyrillic_er":                 0x06d2,
	"Cyrillic_es":                 0x06d3,
	"Cyrillic_te":                 0x06d4,
	"Cyrillic_u":                  0x06d5,
	"Cyrillic_zhe":                0x06d6,
	"Cyrillic_ve":                 0x06d7,
	"Cyrillic_softsign":           0x06d8,
	"Cyrillic_yeru":               0x06d9,
	"Cyrillic_ze":                 0x06da,
	"Cyrillic_sha":                0x06db,
	"Cyrillic_e":                  0x06dc,
	"Cyrillic_shcha":              0x06dd,
	"Cyrillic_che":                0x06de,
	"Cyrillic_hardsign":           0x06df,
	"Cyrillic_YU":                 0x06e0,
	"Cyrillic_A":                  0x06e1,
	"Cyrillic_BE":                 0x06e2,
	"Cyrillic_TSE":                0x06e3,
	"Cyrillic_DE":                 0x06e4,
	"Cyrillic_IE":                 0x06e5,
	"Cyrillic_EF":                 0x06e6,
	"Cyrillic_GHE":                0x06e7,
	"Cyrillic_HA":                 0x06e8,
	"Cyrillic_I":                  0x06e9,
	"Cyrillic_SHORTI":             0x06ea,
	"Cyrillic_KA":                 0x06eb,
	"Cyrillic_EL":                 0x06ec,
	"Cyrillic_EM":                 0x06ed,
	"Cyrillic_EN":                 0x06ee,
	"Cyrillic_O":                  0x06ef,
	"Cyrillic_PE":                 0x06f0,
	"Cyrillic_YA":                 0x06f1,
	"Cyrillic_ER":                 0x06f2,
	"Cyrillic_ES":                 0x06f3,
	"Cyrillic_TE":                 0x06f4,
	"Cyrillic_U":                  0x06f5,
	"Cyrillic_ZHE":                0x06f6,
	"Cyrillic_VE":                 0x06f7,
	"Cyrillic_SOFTSIGN":           0x06f8,
	"Cyrillic_YERU":               0x06f9,
	"Cyrillic_ZE":                 0x06fa,
	"Cyrillic_SHA":                0x06fb,
	"Cyrillic_E":                  0x06fc,
	"Cyrillic_SHCHA":              0x06fd,
	"Cyrillic_CHE":                0x06fe,
	"Cyrillic_HARDSIGN":           0x06ff,
	"Greek_ALPHAaccent":           0x07a1,
	"Greek_EPSILONaccent":         0x07a2,
	"Greek_ETAaccent":             0x07a3,
	"Greek_IOTAaccent":            0x07a4,
	"Greek_IOTAdieresis":          0x07a5,
	"Greek_IOTAdiaeresis":         0x07a5,
	"Greek_OMICRONaccent":         0x07a7,
	"Greek_UPSILONaccent":         0x07a8,
	"Greek_UPSILONdieresis":       0x07a9,
	"Greek_OMEGAaccent":           0x07ab,
	"Greek_accentdieresis":        0x07ae,
	"Greek_horizbar":              0x07af,
	"Greek_alphaaccent":           0x07b1,
	"Greek_epsilonaccent":         0x07b2,
	"Greek_etaaccent":             0x07b3,
	"Greek_iotaaccent":            0x07b4,
	"Greek_iotadieresis":          0x07b5,
	"Greek_iotaaccentdieresis":    0x07b6,
	"Greek_omicronaccent":         0x07b7,
	"Greek_upsilonaccent":         0x07b8,
	"Greek_upsilondieresis":       0x07b9,
	"Greek_upsilonaccentdieresis": 0x07ba,
	"Greek_omegaaccent":           0x07bb,
	"Greek_ALPHA":                 0x07c1,
	"Greek_BETA":                  0x07c2,
	"Greek_GAMMA":                 0x07c3,
	"Greek_DELTA":                 0x07c4,
	"Greek_EPSILON":               0x07c5,
	"Greek_ZETA":                  0x07c6,
	"Greek_ETA":                   0x07c7,
	"Greek_THETA":                 0x07c8,
	"Greek_IOTA":                  0x07c9,
	"Greek_KAPPA":                 0x07ca,
	"Greek_LAMDA":                 0x07cb,
	"Greek_LAMBDA":                0x07cb,
	"Greek_MU":                    0x07cc,
	"Greek_NU":                    0x07cd,
	"Greek_XI":                    0x07ce,
	"Greek_OMICRON":               0x07cf,
	"Greek_PI":                    0x07d0,
	"Greek_RHO":                   0x07d1,
	"Greek_SIGMA":                 0x07d2,
	"Greek_TAU":                   0x07d4,
	"Greek_UPSILON":               0x07d5,
	"Greek_PHI":                   0x07d6,
	"Greek_CHI":                   0x07d7,
	"Greek_PSI":                   0x07d8,
	"Greek_OMEGA":                 0x07d9,
	"Greek_alpha":                 0x07e1,
	"Greek_beta":                  0x07e2,
	"Greek_gamma":                 0x07e3,
	"Greek_delta":                 0x07e4,
	"Greek_epsilon":               0x07e5,
	"Greek_zeta":                  0x07e6,
	"Greek_eta":                   0x07e7,
	"Greek_theta":                 0x07e8,
	"Greek_iota":                  0x07e9,
	"Greek_kappa":                 0x07ea,
	"Greek_lamda":                 0x07eb,
	"Greek_lambda":                0x07eb,
	"Greek_mu":                    0x07ec,
	"Greek_nu":                    0x07ed,
	"Greek_xi":                    0x07ee,
	"Greek_omicron":               0x07ef,
	"Greek_pi":                    0x07f0,
	"Greek_rho":                   0x07f1,
	"Greek_sigma":                 0x07f2,
	"Greek_finalsmallsigma":       0x07f3,
	"Greek_tau":                   0x07f4,
	"Greek_upsilon":               0x07f5,
	"Greek_phi":                   0x07f6,
	"Greek_chi":                   0x07f7,
	"Greek_psi":                   0x07f8,
	"Greek_omega":                 0x07f9,
	"Greek_switch":                0xff7e,
	"leftradical":                 0x08a1,
	"topleftradical":              0x08a2,
	"horizconnector":              0x08a3,
	"topintegral":                 0x08a4,
	"botintegral":                 0x08a5,
	"vertconnector":               0x08a6,
	"topleftsqbracket":            0x08a7,
	"botleftsqbracket":            0x08a8,
	"toprightsqbracket":           0x08a9,
	"botrightsqbracket":           0x08aa,
	"topleftparens":               0x08ab,
	"botleftparens":               0x08ac,
	"toprightparens":              0x08ad,
	"botrightparens":              0x08ae,
	"leftmiddlecurlybrace":        0x08af,
	"rightmiddlecurlybrace":       0x08b0,
	"topleftsummation":            0x08b1,
	"botleftsummation":            0x08b2,
	"topvertsummationconnector":   0x08b3,
	"botvertsummationconnector":   0x08b4,
	"toprightsummation":           0x08b5,
	"botrightsummation":           0x08b6,
	"rightmiddlesummation":        0x08b7,
	"lessthanequal":               0x08bc,
	"notequal":                    0x08bd,
	"greaterthanequal":            0x08be,
	"integral":                    0x08bf,
	"therefore":                   0x08c0,
	"variation":                   0x08c1,
	"infinity":                    0x08c2,
	"nabla":                       0x08c5,
	"approximate":                 0x08c8,
	"similarequal":                0x08c9,
	"ifonlyif":                    0x08cd,
	"implies":                     0x08ce,
	"identical":                   0x08cf,
	"radical":                     0x08d6,
	"includedin":                  0x08da,
	"includes":                    0x08db,
	"intersection":                0x08dc,
	"union":                       0x08dd,
	"logicaland":                  0x08de,
	"logicalor":                   0x08df,
	"partialderivative":           0x08ef,
	"function":                    0x08f6,
	"leftarrow":                   0x08fb,
	"uparrow":                     0x08fc,
	"rightarrow":                  0x08fd,
	"downarrow":                   0x08fe,
	"blank":                       0x09df,
	"soliddiamond":                0x09e0,
	"checkerboard":                0x09e1,
	"ht":                          0x09e2,
	"ff":                          0x09e3,
	"cr":                          0x09e4,
	"lf":                          0x09e5,
	"nl":                          0x09e8,
	"vt":                          0x09e9,
	"lowrightcorner":              0x09ea,
	"uprightcorner":               0x09eb,
	"upleftcorner":                0x09ec,
	"lowleftcorner":               0x09ed,
	"crossinglines":               0x09ee,
	"horizlinescan1":              0x09ef,
	"horizlinescan3":              0x09f0,
	"horizlinescan5":              0x09f1,
	"horizlinescan7":              0x09f2,
	"horizlinescan9":              0x09f3,
	"leftt":                       0x09f4,
	"rightt":                      0x09f5,
	"bott":                        0x09f6,
	"topt":                        0x09f7,
	"vertbar":                     0x09f8,
	"emspace":                     0x0aa1,
	"enspace":                     0x0aa2,
	"em3space":                    0x0aa3,
	"em4space":                    0x0aa4,
	"digitspace":                  0x0aa5,
	"punctspace":                  0x0aa6,
	"thinspace":                   0x0aa7,
	"hairspace":                   0x0aa8,
	"emdash":                      0x0aa9,
	"endash":                      0x0aaa,
	"signifblank":                 0x0aac,
	"ellipsis":                    0x0aae,
	"doubbaselinedot":             0x0aaf,
	"onethird":                    0x0ab0,
	"twothirds":                   0x0ab1,
	"onefifth":                    0x0ab2,
	"twofifths":                   0x0ab3,
	"threefifths":                 0x0ab4,
	"fourfifths":                  0x0ab5,
	"onesixth":                    0x0ab6,
	"fivesixths":                  0x0ab7,
	"careof":                      0x0ab8,
	"figdash":                     0x0abb,
	"leftanglebracket":            0x0abc,
	"decimalpoint":                0x0abd,
	"rightanglebracket":           0x0abe,
	"marker":                      0x0abf,
	"oneeighth":                   0x0ac3,
	"threeeighths":                0x0ac4,
	"fiveeighths":                 0x0ac5,
	"seveneighths":                0x0ac6,
	"trademark":                   0x0ac9,
	"signaturemark":               0x0aca,
	"trademarkincircle":           0x0acb,
	"leftopentriangle":            0x0acc,
	"rightopentriangle":           0x0acd,
	"emopencircle":                0x0ace,
	"emopenrectangle":             0x0acf,
	"leftsinglequotemark":         0x0ad0,
	"rightsinglequotemark":        0x0ad1,
	"leftdoublequotemark":         0x0ad2,
	"rightdoublequotemark":        0x0ad3,
	"prescription":                0x0ad4,
	"minutes":                     0x0ad6,
	"seconds":                     0x0ad7,
	"latincross":                  0x0ad9,
	"hexagram":                    0x0ada,
	"filledrectbullet":            0x0adb,
	"filledlefttribullet":         0x0adc,
	"filledrighttribullet":        0x0add,
	"emfilledcircle":              0x0ade,
	"emfilledrect":                0x0adf,
	"enopencircbullet":            0x0ae0,
	"enopensquarebullet":          0x0ae1,
	"openrectbullet":              0x0ae2,
	"opentribulletup":             0x0ae3,
	"opentribulletdown":           0x0ae4,
	"openstar":                    0x0ae5,
	"enfilledcircbullet":          0x0ae6,
	"enfilledsqbullet":            0x0ae7,
	"filledtribulletup":           0x0ae8,
	"filledtribulletdown":         0x0ae9,
	"leftpointer":                 0x0aea,
	"rightpointer":                0x0aeb,
	"club":                        0x0aec,
	"diamond":                     0x0aed,
	"heart":                       0x0aee,
	"maltesecross":                0x0af0,
	"dagger":                      0x0af1,
	"doubledagger":                0x0af2,
	"checkmark":                   0x0af3,
	"ballotcross":                 0x0af4,
	"musicalsharp":                0x0af5,
	"musicalflat":                 0x0af6,
	"malesymbol":                  0x0af7,
	"femalesymbol":                0x0af8,
	"telephone":                   0x0af9,
	"telephonerecorder":           0x0afa,
	"phonographcopyright":         0x0afb,
	"caret":                       0x0afc,
	"singlelowquotemark":          0x0afd,
	"doublelowquotemark":          0x0afe,
	"cursor":                      0x0aff,
	"leftcaret":                   0x0ba3,
	"rightcaret":                  0x0ba6,
	"downcaret":                   0x0ba8,
	"upcaret":                     0x0ba9,
	"overbar":                     0x0bc0,
	"downtack":                    0x0bc2,
	"upshoe":                      0x0bc3,
	"downstile":                   0x0bc4,
	"underbar":                    0x0bc6,
	"jot":                         0x0bca,
	"quad":                        0x0bcc,
	"uptack":                      0x0bce,
	"circle":                      0x0bcf,
	"upstile":                     0x0bd3,
	"downshoe":                    0x0bd6,
	"rightshoe":                   0x0bd8,
	"leftshoe":                    0x0bda,
	"lefttack":                    0x0bdc,
	"righttack":                   0x0bfc,
	"hebrew_doublelowline":        0x0cdf,
	"hebrew_aleph":                0x0ce0,
	"hebrew_bet":                  0x0ce1,
	"hebrew_beth":                 0x0ce1,
	"hebrew_gimel":                0x0ce2,
	"hebrew_gimmel":               0x0ce2,
	"hebrew_dalet":                0x0ce3,
	"hebrew_daleth":               0x0ce3,
	"hebrew_he":                   0x0ce4,
	"hebrew_waw":                  0x0ce5,
	"hebrew_zain":                 0x0ce6,
	"hebrew_zayin":                0x0ce6,
	"hebrew_chet":                 0x0ce7,
	"hebrew_het":                  0x0ce7,
	"hebrew_tet":                  0x0ce8,
	"hebrew_teth":                 0x0ce8,
	"hebrew_yod":                  0x0ce9,
	"hebrew_finalkaph":            0x0cea,
	"hebrew_kaph":                 0x0ceb,
	"hebrew_lamed":                0x0cec,
	"hebrew_finalmem":             0x0ced,
	"hebrew_mem":                  0x0cee,
	"hebrew_finalnun":             0x0cef,
	"hebrew_nun":                  0x0cf0,
	"hebrew_samech":               0x0cf1,
	"hebrew_samekh":               0x0cf1,
	"hebrew_ayin":                 0x0cf2,
	"hebrew_finalpe":              0x0cf3,
	"hebrew_pe":                   0x0cf4,
	"hebrew_finalzade":            0x0cf5,
	"hebrew_finalzadi":            0x0cf5,
	"hebrew_zade":                 0x0cf6,
	"hebrew_zadi":                 0x0cf6,
	"hebrew_qoph":                 0x0cf7,
	"hebrew_kuf":                  0x0cf7,
	"hebrew_resh":                 0x0cf8,
	"hebrew_shin":                 0x0cf9,
	"hebrew_taw":                  0x0cfa,
	"hebrew_taf":                  0x0cfa,
	"Hebrew_switch":               0xff7e,
	"Thai_kokai":                  0x0da1,
	"Thai_khokhai":                0x0da2,
	"Thai_khokhuat":               0x0da3,
	"Thai_khokhwai":               0x0da4,
	"Thai_khokhon":                0x0da5,
	"Thai_khorakhang":             0x0da6,
	"Thai_ngongu":                 0x0da7,
	"Thai_chochan":                0x0da8,
	"Thai_choching":               0x0da9,
	"Thai_chochang":               0x0daa,
	"Thai_soso":                   0x0dab,
	"Thai_chochoe":                0x0dac,
	"Thai_yoying":                 0x0dad,
	"Thai_dochada":                0x0dae,
	"Thai_topatak":                0x0daf,
	"Thai_thothan":                0x0db0,
	"Thai_thonangmontho":          0x0db1,
	"Thai_thophuthao":             0x0db2,
	"Thai_nonen":                  0x0db3,
	"Thai_dodek":                  0x0db4,
	"Thai_totao":                  0x0db5,
	"Thai_thothung":               0x0db6,
	"Thai_thothahan":              0x0db7,
	"Thai_thothong":               0x0db8,
	"Thai_nonu":                   0x0db9,
	"Thai_bobaimai":               0x0dba,
	"Thai_popla":                  0x0dbb,
	"Thai_phophung":               0x0dbc,
	"Thai_fofa":                   0x0dbd,
	"Thai_phophan":                0x0dbe,
	"Thai_fofan":                  0x0dbf,
	"Thai_phosamphao":             0x0dc0,
	"Thai_moma":                   0x0dc1,
	"Thai_yoyak":                  0x0dc2,
	"Thai_rorua":                  0x0dc3,
	"Thai_ru":                     0x0dc4,
	"Thai_loling":                 0x0dc5,
	"Thai_lu":                     0x0dc6,
	"Thai_wowaen":                 0x0dc7,
	"Thai_sosala":                 0x0dc8,
	"Thai_sorusi":                 0x0dc9,
	"Thai_sosua":                  0x0dca,
	"Thai_hohip":                  0x0dcb,
	"Thai_lochula":                0x0dcc,
	"Thai_oang":                   0x0dcd,
	"Thai_honokhuk":               0x0dce,
	"Thai_paiyannoi":              0x0dcf,
	"Thai_saraa":                  0x0dd0,
	"Thai_maihanakat":             0x0dd1,
	"Thai_saraaa":                 0x0dd2,
	"Thai_saraam":                 0x0dd3,
	"Thai_sarai":                  0x0dd4,
	"Thai_saraii":                 0x0dd5,
	"Thai_saraue":                 0x0dd6,
	"Thai_sarauee":                0x0dd7,
	"Thai_sarau":                  0x0dd8,
	"Thai_sarauu":                 0x0dd9,
	"Thai_phinthu":                0x0dda,
	"Thai_maihanakat_maitho":      0x0dde,
	"Thai_baht":                   0x0ddf,
	"Thai_sarae":                  0x0de0,
	"Thai_saraae":                 0x0de1,
	"Thai_sarao":                  0x0de2,
	"Thai_saraaimaimuan":          0x0de3,
	"Thai_saraaimaimalai":         0x0de4,
	"Thai_lakkhangyao":            0x0de5,
	"Thai_maiyamok":               0x0de6,
	"Thai_maitaikhu":              0x0de7,
	"Thai_maiek":                  0x0de8,
	"Thai_maitho":                 0x0de9,
	"Thai_maitri":                 0x0dea,
	"Thai_maichattawa":            0x0deb,
	"Thai_thanthakhat":            0x0dec,
	"Thai_nikhahit":               0x0ded,
	"Thai_leksun":                 0x0df0,
	"Thai_leknung":                0x0df1,
	"Thai_leksong":                0x0df2,
	"Thai_leksam":                 0x0df3,
	"Thai_leksi":                  0x0df4,
	"Thai_lekha":                  0x0df5,
	"Thai_lekhok":                 0x0df6,
	"Thai_lekchet":                0x0df7,
	"Thai_lekpaet":                0x0df8,
	"Thai_lekkao":                 0x0df9,
	"Hangul":                      0xff31,
	"Hangul_Start":                0xff32,
	"Hangul_End":                  0xff33,
	"Hangul_Hanja":                0xff34,
	"Hangul_Jamo":                 0xff35,
	"Hangul_Romaja":               0xff36,
	"Hangul_Codeinput":            0xff37,
	"Hangul_Jeonja":               0xff38,
	"Hangul_Banja":                0xff39,
	"Hangul_PreHanja":             0xff3a,
	"Hangul_PostHanja":            0xff3b,
	"Hangul_SingleCandidate":      0xff3c,
	"Hangul_MultipleCandidate":    0xff3d,
	"Hangul_PreviousCandidate":    0xff3e,
	"Hangul_Special":              0xff3f,
	"Hangul_switch":               0xff7e,
	"Hangul_Kiyeog":               0x0ea1,
	"Hangul_SsangKiyeog":          0x0ea2,
	"Hangul_KiyeogSios":           0x0ea3,
	"Hangul_Nieun":                0x0ea4,
	"Hangul_NieunJieuj":           0x0ea5,
	"Hangul_NieunHieuh":           0x0ea6,
	"Hangul_Dikeud":               0x0ea7,
	"Hangul_SsangDikeud":          0x0ea8,
	"Hangul_Rieul":                0x0ea9,
	"Hangul_RieulKiyeog":          0x0eaa,
	"Hangul_RieulMieum":           0x0eab,
	"Hangul_RieulPieub":           0x0eac,
	"Hangul_RieulSios":            0x0ead,
	"Hangul_RieulTieut":           0x0eae,
	"Hangul_RieulPhieuf":          0x0eaf,
	"Hangul_RieulHieuh":           0x0eb0,
	"Hangul_Mieum":                0x0eb1,
	"Hangul_Pieub":                0x0eb2,
	"Hangul_SsangPieub":           0x0eb3,
	"Hangul_PieubSios":            0x0eb4,
	"Hangul_Sios":                 0x0eb5,
	"Hangul_SsangSios":            0x0eb6,
	"Hangul_Ieung":                0x0eb7,
	"Hangul_Jieuj":                0x0eb8,
	"Hangul_SsangJieuj":           0x0eb9,
	"Hangul_Cieuc":                0x0eba,
	"Hangul_Khieuq":               0x0ebb,
	"Hangul_Tieut":                0x0ebc,
	"Hangul_Phieuf":               0x0ebd,
	"Hangul_Hieuh":                0x0ebe,
	"Hangul_A":                    0x0ebf,
	"Hangul_AE":                   0x0ec0,
	"Hangul_YA":                   0x0ec1,
	"Hangul_YAE":                  0x0ec2,
	"Hangul_EO":                   0x0ec3,
	"Hangul_E":                    0x0ec4,
	"Hangul_YEO":                  0x0ec5,
	"Hangul_YE":                   0x0ec6,
	"Hangul_O":                    0x0ec7,
	"Hangul_WA":                   0x0ec8,
	"Hangul_WAE":                  0x0ec9,
	"Hangul_OE":                   0x0eca,
	"Hangul_YO":                   0x0ecb,
	"Hangul_U":                    0x0ecc,
	"Hangul_WEO":                  0x0ecd,
	"Hangul_WE":                   0x0ece,
	"Hangul_WI":                   0x0ecf,
	"Hangul_YU":                   0x0ed0,
	"Hangul_EU":                   0x0ed1,
	"Hangul_YI":                   0x0ed2,
	"Hangul_I":                    0x0ed3,
	"Hangul_J_Kiyeog":             0x0ed4,
	"Hangul_J_SsangKiyeog":        0x0ed5,
	"Hangul_J_KiyeogSios":         0x0ed6,
	"Hangul_J_Nieun":              0x0ed7,
	"Hangul_J_NieunJieuj":         0x0ed8,
	"Hangul_J_NieunHieuh":         0x0ed9,
	"Hangul_J_Dikeud":             0x0eda,
	"Hangul_J_Rieul":              0x0edb,
	"Hangul_J_RieulKiyeog":        0x0edc,
	"Hangul_J_RieulMieum":         0x0edd,
	"Hangul_J_RieulPieub":         0x0ede,
	"Hangul_J_RieulSios":          0x0edf,
	"Hangul_J_RieulTieut":         0x0ee0,
	"Hangul_J_RieulPhieuf":        0x0ee1,
	"Hangul_J_RieulHieuh":         0x0ee2,
	"Hangul_J_Mieum":              0x0ee3,
	"Hangul_J_Pieub":              0x0ee4,
	"Hangul_J_PieubSios":          0x0ee5,
	"Hangul_J_Sios":               0x0ee6,
	"Hangul_J_SsangSios":          0x0ee7,
	"Hangul_J_Ieung":              0x0ee8,
	"Hangul_J_Jieuj":              0x0ee9,
	"Hangul_J_Cieuc":              0x0eea,
	"Hangul_J_Khieuq":             0x0eeb,
	"Hangul_J_Tieut":              0x0eec,
	"Hangul_J_Phieuf":             0x0eed,
	"Hangul_J_Hieuh":              0x0eee,
	"Hangul_RieulYeorinHieuh":     0x0eef,
	"Hangul_SunkyeongeumMieum":    0x0ef0,
	"Hangul_SunkyeongeumPieub":    0x0ef1,
	"Hangul_PanSios":              0x0ef2,
	"Hangul_KkogjiDalrinIeung":    0x0ef3,
	"Hangul_SunkyeongeumPhieuf":   0x0ef4,
	"Hangul_YeorinHieuh":          0x0ef5,
	"Hangul_AraeA":                0x0ef6,
	"Hangul_AraeAE":               0x0ef7,
	"Hangul_J_PanSios":            0x0ef8,
	"Hangul_J_KkogjiDalrinIeung":  0x0ef9,
	"Hangul_J_YeorinHieuh":        0x0efa,
	"Korean_Won":                  0x0eff,
	"Armenian_ligature_ew":        0x1000587,
	"Armenian_full_stop":          0x1000589,
	"Armenian_verjaket":           0x1000589,
	"Armenian_separation_mark":    0x100055d,
	"Armenian_but":                0x100055d,
	"Armenian_hyphen":             0x100058a,
	"Armenian_yentamna":           0x100058a,
	"Armenian_exclam":             0x100055c,
	"Armenian_amanak":             0x100055c,
	"Armenian_accent":             0x100055b,
	"Armenian_shesht":             0x100055b,
	"Armenian_question":           0x100055e,
	"Armenian_paruyk":             0x100055e,
	"Armenian_AYB":                0x1000531,
	"Armenian_ayb":                0x1000561,
	"Armenian_BEN":                0x1000532,
	"Armenian_ben":                0x1000562,
	"Armenian_GIM":                0x1000533,
	"Armenian_gim":                0x1000563,
	"Armenian_DA":                 0x1000534,
	"Armenian_da":                 0x1000564,
	"Armenian_YECH":               0x1000535,
	"Armenian_yech":               0x1000565,
	"Armenian_ZA":                 0x1000536,
	"Armenian_za":                 0x1000566,
	"Armenian_E":                  0x1000537,
	"Armenian_e":                  0x1000567,
	"Armenian_AT":                 0x1000538,
	"Armenian_at":                 0x1000568,
	"Armenian_TO":                 0x1000539,
	"Armenian_to":                 0x1000569,
	"Armenian_ZHE":                0x100053a,
	"Armenian_zhe":                0x100056a,
	"Armenian_INI":                0x100053b,
	"Armenian_ini":                0x100056b,
	"Armenian_LYUN":               0x100053c,
	"Armenian_lyun":               0x100056c,
	"Armenian_KHE":                0x100053d,
	"Armenian_khe":                0x100056d,
	"Armenian_TSA":                0x100053e,
	"Armenian_tsa":                0x100056e,
	"Armenian_KEN":                0x100053f,
	"Armenian_ken":                0x100056f,
	"Armenian_HO":                 0x1000540,
	"Armenian_ho":                 0x1000570,
	"Armenian_DZA":                0x1000541,
	"Armenian_dza":                0x1000571,
	"Armenian_GHAT":               0x1000542,
	"Armenian_ghat":               0x1000572,
	"Armenian_TCHE":               0x1000543,
	"Armenian_tche":               0x1000573,
	"Armenian_MEN":                0x1000544,
	"Armenian_men":                0x1000574,
	"Armenian_HI":                 0x1000545,
	"Armenian_hi":                 0x1000575,
	"Armenian_NU":                 0x1000546,
	"Armenian_nu":                 0x1000576,
	"Armenian_SHA":                0x1000547,
	"Armenian_sha":                0x1000577,
	"Armenian_VO":                 0x1000548,
	"Armenian_vo":                 0x1000578,
	"Armenian_CHA":                0x1000549,
	"Armenian_cha":                0x1000579,
	"Armenian_PE":                 0x100054a,
	"Armenian_pe":                 0x100057a,
	"Armenian_JE":                 0x100054b,
	"Armenian_je":                 0x100057b,
	"Armenian_RA":                 0x100054c,
	"Armenian_ra":                 0x100057c,
	"Armenian_SE":                 0x100054d,
	"Armenian_se":                 0x100057d,
	"Armenian_VEV":                0x100054e,
	"Armenian_vev":                0x100057e,
	"Armenian_TYUN":               0x100054f,
	"Armenian_tyun":               0x100057f,
	"Armenian_RE":                 0x1000550,
	"Armenian_re":                 0x1000580,
	"Armenian_TSO":                0x1000551,
	"Armenian_tso":                0x1000581,
	"Armenian_VYUN":               0x1000552,
	"Armenian_vyun":               0x1000582,
	"Armenian_PYUR":               0x1000553,
	"Armenian_pyur":               0x1000583,
	"Armenian_KE":                 0x1000554,
	"Armenian_ke":                 0x1000584,
	"Armenian_O":                  0x1000555,
	"Armenian_o":                  0x1000585,
	"Armenian_FE":                 0x1000556,
	"Armenian_fe":                 0x1000586,
	"Armenian_apostrophe":         0x100055a,
	"Georgian_an":                 0x10010d0,
	"Georgian_ban":                0x10010d1,
	"Georgian_gan":                0x10010d2,
	"Georgian_don":                0x10010d3,
	"Georgian_en":                 0x10010d4,
	"Georgian_vin":                0x10010d5,
	"Georgian_zen":                0x10010d6,
	"Georgian_tan":                0x10010d7,
	"Georgian_in":                 0x10010d8,
	"Georgian_kan":                0x10010d9,
	"Georgian_las":                0x10010da,
	"Georgian_man":                0x10010db,
	"Georgian_nar":                0x10010dc,
	"Georgian_on":                 0x10010dd,
	"Georgian_par":                0x10010de,
	"Georgian_zhar":               0x10010df,
	"Georgian_rae":                0x10010e0,
	"Georgian_san":                0x10010e1,
	"Georgian_tar":                0x10010e2,
	"Georgian_un":                 0x10010e3,
	"Georgian_phar":               0x10010e4,
	"Georgian_khar":               0x10010e5,
	"Georgian_ghan":               0x10010e6,
	"Georgian_qar":                0x10010e7,
	"Georgian_shin":               0x10010e8,
	"Georgian_chin":               0x10010e9,
	"Georgian_can":                0x10010ea,
	"Georgian_jil":                0x10010eb,
	"Georgian_cil":                0x10010ec,
	"Georgian_char":               0x10010ed,
	"Georgian_xan":                0x10010ee,
	"Georgian_jhan":               0x10010ef,
	"Georgian_hae":                0x10010f0,
	"Georgian_he":                 0x10010f1,
	"Georgian_hie":                0x10010f2,
	"Georgian_we":                 0x10010f3,
	"Georgian_har":                0x10010f4,
	"Georgian_hoe":                0x10010f5,
	"Georgian_fi":                 0x10010f6,
	"Xabovedot":                   0x1001e8a,
	"Ibreve":                      0x100012c,
	"Zstroke":                     0x10001b5,
	"Gcaron":                      0x10001e6,
	"Ocaron":                      0x10001d1,
	"Obarred":                     0x100019f,
	"xabovedot":                   0x1001e8b,
	"ibreve":                      0x100012d,
	"zstroke":                     0x10001b6,
	"gcaron":                      0x10001e7,
	"ocaron":                      0x10001d2,
	"obarred":                     0x1000275,
	"SCHWA":                       0x100018f,
	"schwa":                       0x1000259,
	"Lbelowdot":                   0x1001e36,
	"lbelowdot":                   0x1001e37,
	"Abelowdot":                   0x1001ea0,
	"abelowdot":                   0x1001ea1,
	"Ahook":                       0x1001ea2,
	"ahook":                       0x1001ea3,
	"Acircumflexacute":            0x1001ea4,
	"acircumflexacute":            0x1001ea5,
	"Acircumflexgrave":            0x1001ea6,
	"acircumflexgrave":            0x1001ea7,
	"Acircumflexhook":             0x1001ea8,
	"acircumflexhook":             0x1001ea9,
	"Acircumflextilde":            0x1001eaa,
	"acircumflextilde":            0x1001eab,
	"Acircumflexbelowdot":         0x1001eac,
	"acircumflexbelowdot":         0x1001ead,
	"Abreveacute":                 0x1001eae,
	"abreveacute":                 0x1001eaf,
	"Abrevegrave":                 0x1001eb0,
	"abrevegrave":                 0x1001eb1,
	"Abrevehook":                  0x1001eb2,
	"abrevehook":                  0x1001eb3,
	"Abrevetilde":                 0x1001eb4,
	"abrevetilde":                 0x1001eb5,
	"Abrevebelowdot":              0x1001eb6,
	"abrevebelowdot":              0x1001eb7,
	"Ebelowdot":                   0x1001eb8,
	"ebelowdot":                   0x1001eb9,
	"Ehook":                       0x1001eba,
	"ehook":                       0x1001ebb,
	"Etilde":                      0x1001ebc,
	"etilde":                      0x1001ebd,
	"Ecircumflexacute":            0x1001ebe,
	"ecircumflexacute":            0x1001ebf,
	"Ecircumflexgrave":            0x1001ec0,
	"ecircumflexgrave":            0x1001ec1,
	"Ecircumflexhook":             0x1001ec2,
	"ecircumflexhook":             0x1001ec3,
	"Ecircumflextilde":            0x1001ec4,
	"ecircumflextilde":            0x1001ec5,
	"Ecircumflexbelowdot":         0x1001ec6,
	"ecircumflexbelowdot":         0x1001ec7,
	"Ihook":                       0x1001ec8,
	"ihook":                       0x1001ec9,
	"Ibelowdot":                   0x1001eca,
	"ibelowdot":                   0x1001ecb,
	"Obelowdot":                   0x1001ecc,
	"obelowdot":                   0x1001ecd,
	"Ohook":                       0x1001ece,
	"ohook":                       0x1001ecf,
	"Ocircumflexacute":            0x1001ed0,
	"ocircumflexacute":            0x1001ed1,
	"Ocircumflexgrave":            0x1001ed2,
	"ocircumflexgrave":            0x1001ed3,
	"Ocircumflexhook":             0x1001ed4,
	"ocircumflexhook":             0x1001ed5,
	"Ocircumflextilde":            0x1001ed6,
	"ocircumflextilde":            0x1001ed7,
	"Ocircumflexbelowdot":         0x1001ed8,
	"ocircumflexbelowdot":         0x1001ed9,
	"Ohornacute":                  0x1001eda,
	"ohornacute":                  0x1001edb,
	"Ohorngrave":                  0x1001edc,
	"ohorngrave":                  0x1001edd,
	"Ohornhook":                   0x1001ede,
	"ohornhook":                   0x1001edf,
	"Ohorntilde":                  0x1001ee0,
	"ohorntilde":                  0x1001ee1,
	"Ohornbelowdot":               0x1001ee2,
	"ohornbelowdot":               0x1001ee3,
	"Ubelowdot":                   0x1001ee4,
	"ubelowdot":                   0x1001ee5,
	"Uhook":                       0x1001ee6,
	"uhook":                       0x1001ee7,
	"Uhornacute":                  0x1001ee8,
	"uhornacute":                  0x1001ee9,
	"Uhorngrave":                  0x1001eea,
	"uhorngrave":                  0x1001eeb,
	"Uhornhook":                   0x1001eec,
	"uhornhook":                   0x1001eed,
	"Uhorntilde":                  0x1001eee,
	"uhorntilde":                  0x1001eef,
	"Uhornbelowdot":               0x1001ef0,
	"uhornbelowdot":               0x1001ef1,
	"Ybelowdot":                   0x1001ef4,
	"ybelowdot":                   0x1001ef5,
	"Yhook":                       0x1001ef6,
	"yhook":                       0x1001ef7,
	"Ytilde":                      0x1001ef8,
	"ytilde":                      0x1001ef9,
	"Ohorn":                       0x10001a0,
	"ohorn":                       0x10001a1,
	"Uhorn":                       0x10001af,
	"uhorn":                       0x10001b0,
	"EcuSign":                     0x10020a0,
	"ColonSign":                   0x10020a1,
	"CruzeiroSign":                0x10020a2,
	"FFrancSign":                  0x10020a3,
	"LiraSign":                    0x10020a4,
	"MillSign":                    0x10020a5,
	"NairaSign":                   0x10020a6,
	"PesetaSign":                  0x10020a7,
	"RupeeSign":                   0x10020a8,
	"WonSign":                     0x10020a9,
	"NewSheqelSign":               0x10020aa,
	"DongSign":                    0x10020ab,
	"EuroSign":                    0x20ac,
	"zerosuperior":                0x1002070,
	"foursuperior":                0x1002074,
	"fivesuperior":                0x1002075,
	"sixsuperior":                 0x1002076,
	"sevensuperior":               0x1002077,
	"eightsuperior":               0x1002078,
	"ninesuperior":                0x1002079,
	"zerosubscript":               0x1002080,
	"onesubscript":                0x1002081,
	"twosubscript":                0x1002082,
	"threesubscript":              0x1002083,
	"foursubscript":               0x1002084,
	"fivesubscript":               0x1002085,
	"sixsubscript":                0x1002086,
	"sevensubscript":              0x1002087,
	"eightsubscript":              0x1002088,
	"ninesubscript":               0x1002089,
	"partdifferential":            0x1002202,
	"emptyset":                    0x1002205,
	"elementof":                   0x1002208,
	"notelementof":                0x1002209,
	"containsas":                  0x100220B,
	"squareroot":                  0x100221A,
	"cuberoot":                    0x100221B,
	"fourthroot":                  0x100221C,
	"dintegral":                   0x100222C,
	"tintegral":                   0x100222D,
	"because":                     0x1002235,
	"approxeq":                    0x1002248,
	"notapproxeq":                 0x1002247,
	"notidentical":                0x1002262,
	"stricteq":                    0x1002263,
	"braille_dot_1":               0xfff1,
	"braille_dot_2":               0xfff2,
	"braille_dot_3":               0xfff3,
	"braille_dot_4":               0xfff4,
	"braille_dot_5":               0xfff5,
	"braille_dot_6":               0xfff6,
	"braille_dot_7":               0xfff7,
	"braille_dot_8":               0xfff8,
	"braille_dot_9":               0xfff9,
	"braille_dot_10":              0xfffa,
	"braille_blank":               0x1002800,
	"braille_dots_1":              0x1002801,
	"braille_dots_2":              0x1002802,
	"braille_dots_12":             0x1002803,
	"braille_dots_3":              0x1002804,
	"braille_dots_13":             0x1002805,
	"braille_dots_23":             0x1002806,
	"braille_dots_123":            0x1002807,
	"braille_dots_4":              0x1002808,
	"braille_dots_14":             0x1002809,
	"braille_dots_24":             0x100280a,
	"braille_dots_124":            0x100280b,
	"braille_dots_34":             0x100280c,
	"braille_dots_134":            0x100280d,
	"braille_dots_234":            0x100280e,
	"braille_dots_1234":           0x100280f,
	"braille_dots_5":              0x1002810,
	"braille_dots_15":             0x1002811,
	"braille_dots_25":             0x1002812,
	"braille_dots_125":            0x1002813,
	"braille_dots_35":             0x1002814,
	"braille_dots_135":            0x1002815,
	"braille_dots_235":            0x1002816,
	"braille_dots_1235":           0x1002817,
	"braille_dots_45":             0x1002818,
	"braille_dots_145":            0x1002819,
	"braille_dots_245":            0x100281a,
	"braille_dots_1245":           0x100281b,
	"braille_dots_345":            0x100281c,
	"braille_dots_1345":           0x100281d,
	"braille_dots_2345":           0x100281e,
	"braille_dots_12345":          0x100281f,
	"braille_dots_6":              0x1002820,
	"braille_dots_16":             0x1002821,
	"braille_dots_26":             0x1002822,
	"braille_dots_126":            0x1002823,
	"braille_dots_36":             0x1002824,
	"braille_dots_136":            0x1002825,
	"braille_dots_236":            0x1002826,
	"braille_dots_1236":           0x1002827,
	"braille_dots_46":             0x1002828,
	"braille_dots_146":            0x1002829,
	"braille_dots_246":            0x100282a,
	"braille_dots_1246":           0x100282b,
	"braille_dots_346":            0x100282c,
	"braille_dots_1346":           0x100282d,
	"braille_dots_2346":           0x100282e,
	"braille_dots_12346":          0x100282f,
	"braille_dots_56":             0x1002830,
	"braille_dots_156":            0x1002831,
	"braille_dots_256":            0x1002832,
	"braille_dots_1256":           0x1002833,
	"braille_dots_356":            0x1002834,
	"braille_dots_1356":           0x1002835,
	"braille_dots_2356":           0x1002836,
	"braille_dots_12356":          0x1002837,
	"braille_dots_456":            0x1002838,
	"braille_dots_1456":           0x1002839,
	"braille_dots_2456":           0x100283a,
	"braille_dots_12456":          0x100283b,
	"braille_dots_3456":           0x100283c,
	"braille_dots_13456":          0x100283d,
	"braille_dots_23456":          0x100283e,
	"braille_dots_123456":         0x100283f,
	"braille_dots_7":              0x1002840,
	"braille_dots_17":             0x1002841,
	"braille_dots_27":             0x1002842,
	"braille_dots_127":            0x1002843,
	"braille_dots_37":             0x1002844,
	"braille_dots_137":            0x1002845,
	"braille_dots_237":            0x1002846,
	"braille_dots_1237":           0x1002847,
	"braille_dots_47":             0x1002848,
	"braille_dots_147":            0x1002849,
	"braille_dots_247":            0x100284a,
	"braille_dots_1247":           0x100284b,
	"braille_dots_347":            0x100284c,
	"braille_dots_1347":           0x100284d,
	"braille_dots_2347":           0x100284e,
	"braille_dots_12347":          0x100284f,
	"braille_dots_57":             0x1002850,
	"braille_dots_157":            0x1002851,
	"braille_dots_257":            0x1002852,
	"braille_dots_1257":           0x1002853,
	"braille_dots_357":            0x1002854,
	"braille_dots_1357":           0x1002855,
	"braille_dots_2357":           0x1002856,
	"braille_dots_12357":          0x1002857,
	"braille_dots_457":            0x1002858,
	"braille_dots_1457":           0x1002859,
	"braille_dots_2457":           0x100285a,
	"braille_dots_12457":          0x100285b,
	"braille_dots_3457":           0x100285c,
	"braille_dots_13457":          0x100285d,
	"braille_dots_23457":          0x100285e,
	"braille_dots_123457":         0x100285f,
	"braille_dots_67":             0x1002860,
	"braille_dots_167":            0x1002861,
	"braille_dots_267":            0x1002862,
	"braille_dots_1267":           0x1002863,
	"braille_dots_367":            0x1002864,
	"braille_dots_1367":           0x1002865,
	"braille_dots_2367":           0x1002866,
	"braille_dots_12367":          0x1002867,
	"braille_dots_467":            0x1002868,
	"braille_dots_1467":           0x1002869,
	"braille_dots_2467":           0x100286a,
	"braille_dots_12467":          0x100286b,
	"braille_dots_3467":           0x100286c,
	"braille_dots_13467":          0x100286d,
	"braille_dots_23467":          0x100286e,
	"braille_dots_123467":         0x100286f,
	"braille_dots_567":            0x1002870,
	"braille_dots_1567":           0x1002871,
	"braille_dots_2567":           0x1002872,
	"braille_dots_12567":          0x1002873,
	"braille_dots_3567":           0x1002874,
	"braille_dots_13567":          0x1002875,
	"braille_dots_23567":          0x1002876,
	"braille_dots_123567":         0x1002877,
	"braille_dots_4567":           0x1002878,
	"braille_dots_14567":          0x1002879,
	"braille_dots_24567":          0x100287a,
	"braille_dots_124567":         0x100287b,
	"braille_dots_34567":          0x100287c,
	"braille_dots_134567":         0x100287d,
	"braille_dots_234567":         0x100287e,
	"braille_dots_1234567":        0x100287f,
	"braille_dots_8":              0x1002880,
	"braille_dots_18":             0x1002881,
	"braille_dots_28":             0x1002882,
	"braille_dots_128":            0x1002883,
	"braille_dots_38":             0x1002884,
	"braille_dots_138":            0x1002885,
	"braille_dots_238":            0x1002886,
	"braille_dots_1238":           0x1002887,
	"braille_dots_48":             0x1002888,
	"braille_dots_148":            0x1002889,
	"braille_dots_248":            0x100288a,
	"braille_dots_1248":           0x100288b,
	"braille_dots_348":            0x100288c,
	"braille_dots_1348":           0x100288d,
	"braille_dots_2348":           0x100288e,
	"braille_dots_12348":          0x100288f,
	"braille_dots_58":             0x1002890,
	"braille_dots_158":            0x1002891,
	"braille_dots_258":            0x1002892,
	"braille_dots_1258":           0x1002893,
	"braille_dots_358":            0x1002894,
	"braille_dots_1358":           0x1002895,
	"braille_dots_2358":           0x1002896,
	"braille_dots_12358":          0x1002897,
	"braille_dots_458":            0x1002898,
	"braille_dots_1458":           0x1002899,
	"braille_dots_2458":           0x100289a,
	"braille_dots_12458":          0x100289b,
	"braille_dots_3458":           0x100289c,
	"braille_dots_13458":          0x100289d,
	"braille_dots_23458":          0x100289e,
	"braille_dots_123458":         0x100289f,
	"braille_dots_68":             0x10028a0,
	"braille_dots_168":            0x10028a1,
	"braille_dots_268":            0x10028a2,
	"braille_dots_1268":           0x10028a3,
	"braille_dots_368":            0x10028a4,
	"braille_dots_1368":           0x10028a5,
	"braille_dots_2368":           0x10028a6,
	"braille_dots_12368":          0x10028a7,
	"braille_dots_468":            0x10028a8,
	"braille_dots_1468":           0x10028a9,
	"braille_dots_2468":           0x10028aa,
	"braille_dots_12468":          0x10028ab,
	"braille_dots_3468":           0x10028ac,
	"braille_dots_13468":          0x10028ad,
	"braille_dots_23468":          0x10028ae,
	"braille_dots_123468":         0x10028af,
	"braille_dots_568":            0x10028b0,
	"braille_dots_1568":           0x10028b1,
	"braille_dots_2568":           0x10028b2,
	"braille_dots_12568":          0x10028b3,
	"braille_dots_3568":           0x10028b4,
	"braille_dots_13568":          0x10028b5,
	"braille_dots_23568":          0x10028b6,
	"braille_dots_123568":         0x10028b7,
	"braille_dots_4568":           0x10028b8,
	"braille_dots_14568":          0x10028b9,
	"braille_dots_24568":          0x10028ba,
	"braille_dots_124568":         0x10028bb,
	"braille_dots_34568":          0x10028bc,
	"braille_dots_134568":         0x10028bd,
	"braille_dots_234568":         0x10028be,
	"braille_dots_1234568":        0x10028bf,
	"braille_dots_78":             0x10028c0,
	"braille_dots_178":            0x10028c1,
	"braille_dots_278":            0x10028c2,
	"braille_dots_1278":           0x10028c3,
	"braille_dots_378":            0x10028c4,
	"braille_dots_1378":           0x10028c5,
	"braille_dots_2378":           0x10028c6,
	"braille_dots_12378":          0x10028c7,
	"braille_dots_478":            0x10028c8,
	"braille_dots_1478":           0x10028c9,
	"braille_dots_2478":           0x10028ca,
	"braille_dots_12478":          0x10028cb,
	"braille_dots_3478":           0x10028cc,
	"braille_dots_13478":          0x10028cd,
	"braille_dots_23478":          0x10028ce,
	"braille_dots_123478":         0x10028cf,
	"braille_dots_578":            0x10028d0,
	"braille_dots_1578":           0x10028d1,
	"braille_dots_2578":           0x10028d2,
	"braille_dots_12578":          0x10028d3,
	"braille_dots_3578":           0x10028d4,
	"braille_dots_13578":          0x10028d5,
	"braille_dots_23578":          0x10028d6,
	"braille_dots_123578":         0x10028d7,
	"braille_dots_4578":           0x10028d8,
	"braille_dots_14578":          0x10028d9,
	"braille_dots_24578":          0x10028da,
	"braille_dots_124578":         0x10028db,
	"braille_dots_34578":          0x10028dc,
	"braille_dots_134578":         0x10028dd,
	"braille_dots_234578":         0x10028de,
	"braille_dots_1234578":        0x10028df,
	"braille_dots_678":            0x10028e0,
	"braille_dots_1678":           0x10028e1,
	"braille_dots_2678":           0x10028e2,
	"braille_dots_12678":          0x10028e3,
	"braille_dots_3678":           0x10028e4,
	"braille_dots_13678":          0x10028e5,
	"braille_dots_23678":          0x10028e6,
	"braille_dots_123678":         0x10028e7,
	"braille_dots_4678":           0x10028e8,
	"braille_dots_14678":          0x10028e9,
	"braille_dots_24678":          0x10028ea,
	"braille_dots_124678":         0x10028eb,
	"braille_dots_34678":          0x10028ec,
	"braille_dots_134678":         0x10028ed,
	"braille_dots_234678":         0x10028ee,
	"braille_dots_1234678":        0x10028ef,
	"braille_dots_5678":           0x10028f0,
	"braille_dots_15678":          0x10028f1,
	"braille_dots_25678":          0x10028f2,
	"braille_dots_125678":         0x10028f3,
	"braille_dots_35678":          0x10028f4,
	"braille_dots_135678":         0x10028f5,
	"braille_dots_235678":         0x10028f6,
	"braille_dots_1235678":        0x10028f7,
	"braille_dots_45678":          0x10028f8,
	"braille_dots_145678":         0x10028f9,
	"braille_dots_245678":         0x10028fa,
	"braille_dots_1245678":        0x10028fb,
	"braille_dots_345678":         0x10028fc,
	"braille_dots_1345678":        0x10028fd,
	"braille_dots_2345678":        0x10028fe,
	"braille_dots_12345678":       0x10028ff,

	"XF86ModeLock":          0x1008FF01,
	"XF86MonBrightnessUp":   0x1008FF02,
	"XF86MonBrightnessDown": 0x1008FF03,
	"XF86KbdLightOnOff":     0x1008FF04,
	"XF86KbdBrightnessUp":   0x1008FF05,
	"XF86KbdBrightnessDown": 0x1008FF06,
	"XF86Standby":           0x1008FF10,
	"XF86AudioLowerVolume":  0x1008FF11,
	"XF86AudioMute":         0x1008FF12,
	"XF86AudioRaiseVolume":  0x1008FF13,
	"XF86AudioPlay":         0x1008FF14,
	"XF86AudioStop":         0x1008FF15,
	"XF86AudioPrev":         0x1008FF16,
	"XF86AudioNext":         0x1008FF17,
	"XF86HomePage":          0x1008FF18,
	"XF86Mail":              0x1008FF19,
	"XF86Start":             0x1008FF1A,
	"XF86Search":            0x1008FF1B,
	"XF86AudioRecord":       0x1008FF1C,
	"XF86Calculator":        0x1008FF1D,
	"XF86Memo":              0x1008FF1E,
	"XF86ToDoList":          0x1008FF1F,
	"XF86Calendar":          0x1008FF20,
	"XF86PowerDown":         0x1008FF21,
	"XF86ContrastAdjust":    0x1008FF22,
	"XF86RockerUp":          0x1008FF23,
	"XF86RockerDown":        0x1008FF24,
	"XF86RockerEnter":       0x1008FF25,
	"XF86Back":              0x1008FF26,
	"XF86Forward":           0x1008FF27,
	"XF86Stop":              0x1008FF28,
	"XF86Refresh":           0x1008FF29,
	"XF86PowerOff":          0x1008FF2A,
	"XF86WakeUp":            0x1008FF2B,
	"XF86Eject":             0x1008FF2C,
	"XF86ScreenSaver":       0x1008FF2D,
	"XF86WWW":               0x1008FF2E,
	"XF86Sleep":             0x1008FF2F,
	"XF86Favorites":         0x1008FF30,
	"XF86AudioPause":        0x1008FF31,
	"XF86AudioMedia":        0x1008FF32,
	"XF86MyComputer":        0x1008FF33,
	"XF86VendorHome":        0x1008FF34,
	"XF86LightBulb":         0x1008FF35,
	"XF86Shop":              0x1008FF36,
	"XF86History":           0x1008FF37,
	"XF86OpenURL":           0x1008FF38,
	"XF86AddFavorite":       0x1008FF39,
	"XF86HotLinks":          0x1008FF3A,
	"XF86BrightnessAdjust":  0x1008FF3B,
	"XF86Finance":           0x1008FF3C,
	"XF86Community":         0x1008FF3D,
	"XF86AudioRewind":       0x1008FF3E,
	"XF86BackForward":       0x1008FF3F,
	"XF86Launch0":           0x1008FF40,
	"XF86Launch1":           0x1008FF41,
	"XF86Launch2":           0x1008FF42,
	"XF86Launch3":           0x1008FF43,
	"XF86Launch4":           0x1008FF44,
	"XF86Launch5":           0x1008FF45,
	"XF86Launch6":           0x1008FF46,
	"XF86Launch7":           0x1008FF47,
	"XF86Launch8":           0x1008FF48,
	"XF86Launch9":           0x1008FF49,
	"XF86LaunchA":           0x1008FF4A,
	"XF86LaunchB":           0x1008FF4B,
	"XF86LaunchC":           0x1008FF4C,
	"XF86LaunchD":           0x1008FF4D,
	"XF86LaunchE":           0x1008FF4E,
	"XF86LaunchF":           0x1008FF4F,
	"XF86ApplicationLeft":   0x1008FF50,
	"XF86ApplicationRight":  0x1008FF51,
	"XF86Book":              0x1008FF52,
	"XF86CD":                0x1008FF53,
	"XF86Calculater":        0x1008FF54,
	"XF86Clear":             0x1008FF55,
	"XF86Close":             0x1008FF56,
	"XF86Copy":              0x1008FF57,
	"XF86Cut":               0x1008FF58,
	"XF86Display":           0x1008FF59,
	"XF86DOS":               0x1008FF5A,
	"XF86Documents":         0x1008FF5B,
	"XF86Excel":             0x1008FF5C,
	"XF86Explorer":          0x1008FF5D,
	"XF86Game":              0x1008FF5E,
	"XF86Go":                0x1008FF5F,
	"XF86iTouch":            0x1008FF60,
	"XF86LogOff":            0x1008FF61,
	"XF86Market":            0x1008FF62,
	"XF86Meeting":           0x1008FF63,
	"XF86MenuKB":            0x1008FF65,
	"XF86MenuPB":            0x1008FF66,
	"XF86MySites":           0x1008FF67,
	"XF86New":               0x1008FF68,
	"XF86News":              0x1008FF69,
	"XF86OfficeHome":        0x1008FF6A,
	"XF86Open":              0x1008FF6B,
	"XF86Option":            0x1008FF6C,
	"XF86Paste":             0x1008FF6D,
	"XF86Phone":             0x1008FF6E,
	"XF86Q":                 0x1008FF70,
	"XF86Reply":             0x1008FF72,
	"XF86Reload":            0x1008FF73,
	"XF86RotateWindows":     0x1008FF74,
	"XF86RotationPB":        0x1008FF75,
	"XF86RotationKB":        0x1008FF76,
	"XF86Save":              0x1008FF77,
	"XF86ScrollUp":          0x1008FF78,
	"XF86ScrollDown":        0x1008FF79,
	"XF86ScrollClick":       0x1008FF7A,
	"XF86Send":              0x1008FF7B,
	"XF86Spell":             0x1008FF7C,
	"XF86SplitScreen":       0x1008FF7D,
	"XF86Support":           0x1008FF7E,
	"XF86TaskPane":          0x1008FF7F,
	"XF86Terminal":          0x1008FF80,
	"XF86Tools":             0x1008FF81,
	"XF86Travel":            0x1008FF82,
	"XF86UserPB":            0x1008FF84,
	"XF86User1KB":           0x1008FF85,
	"XF86User2KB":           0x1008FF86,
	"XF86Video":             0x1008FF87,
	"XF86WheelButton":       0x1008FF88,
	"XF86Word":              0x1008FF89,
	"XF86Xfer":              0x1008FF8A,
	"XF86ZoomIn":            0x1008FF8B,
	"XF86ZoomOut":           0x1008FF8C,
	"XF86Away":              0x1008FF8D,
	"XF86Messenger":         0x1008FF8E,
	"XF86WebCam":            0x1008FF8F,
	"XF86MailForward":       0x1008FF90,
	"XF86Pictures":          0x1008FF91,
	"XF86Music":             0x1008FF92,
	"XF86Battery":           0x1008FF93,
	"XF86Bluetooth":         0x1008FF94,
	"XF86WLAN":              0x1008FF95,
	"XF86UWB":               0x1008FF96,
	"XF86AudioForward":      0x1008FF97,
	"XF86AudioRepeat":       0x1008FF98,
	"XF86AudioRandomPlay":   0x1008FF99,
	"XF86Subtitle":          0x1008FF9A,
	"XF86AudioCycleTrack":   0x1008FF9B,
	"XF86CycleAngle":        0x1008FF9C,
	"XF86FrameBack":         0x1008FF9D,
	"XF86FrameForward":      0x1008FF9E,
	"XF86Time":              0x1008FF9F,
	"XF86Select":            0x1008FFA0,
	"XF86View":              0x1008FFA1,
	"XF86TopMenu":           0x1008FFA2,
	"XF86Red":               0x1008FFA3,
	"XF86Green":             0x1008FFA4,
	"XF86Yellow":            0x1008FFA5,
	"XF86Blue":              0x1008FFA6,
	"XF86Suspend":           0x1008FFA7,
	"XF86Hibernate":         0x1008FFA8,
	"XF86TouchpadToggle":    0x1008FFA9,
	"XF86TouchpadOn":        0x1008FFB0,
	"XF86TouchpadOff":       0x1008FFB1,
	"XF86AudioMicMute":      0x1008FFB2,
	"XF86Switch_VT_1":       0x1008FE01,
	"XF86Switch_VT_2":       0x1008FE02,
	"XF86Switch_VT_3":       0x1008FE03,
	"XF86Switch_VT_4":       0x1008FE04,
	"XF86Switch_VT_5":       0x1008FE05,
	"XF86Switch_VT_6":       0x1008FE06,
	"XF86Switch_VT_7":       0x1008FE07,
	"XF86Switch_VT_8":       0x1008FE08,
	"XF86Switch_VT_9":       0x1008FE09,
	"XF86Switch_VT_10":      0x1008FE0A,
	"XF86Switch_VT_11":      0x1008FE0B,
	"XF86Switch_VT_12":      0x1008FE0C,
	"XF86Ungrab":            0x1008FE20,
	"XF86ClearGrab":         0x1008FE21,
	"XF86Next_VMode":        0x1008FE22,
	"XF86Prev_VMode":        0x1008FE23,
	"XF86LogWindowTree":     0x1008FE24,
	"XF86LogGrabInfo":       0x1008FE25,
}

var keysymUnicode = map[xproto.Keysym]string{
	0x01a1: "\u0104", // Aogonek = LATIN CAPITAL LETTER A WITH OGONEK
	0x01a2: "\u02d8", // breve = BREVE
	0x01a3: "\u0141", // Lstroke = LATIN CAPITAL LETTER L WITH STROKE
	0x01a5: "\u013d", // Lcaron = LATIN CAPITAL LETTER L WITH CARON
	0x01a6: "\u015a", // Sacute = LATIN CAPITAL LETTER S WITH ACUTE
	0x01a9: "\u0160", // Scaron = LATIN CAPITAL LETTER S WITH CARON
	0x01aa: "\u015e", // Scedilla = LATIN CAPITAL LETTER S WITH CEDILLA
	0x01ab: "\u0164", // Tcaron = LATIN CAPITAL LETTER T WITH CARON
	0x01ac: "\u0179", // Zacute = LATIN CAPITAL LETTER Z WITH ACUTE
	0x01ae: "\u017d", // Zcaron = LATIN CAPITAL LETTER Z WITH CARON
	0x01af: "\u017b", // Zabovedot = LATIN CAPITAL LETTER Z WITH DOT ABOVE
	0x01b1: "\u0105", // aogonek = LATIN SMALL LETTER A WITH OGONEK
	0x01b2: "\u02db", // ogonek = OGONEK
	0x01b3: "\u0142", // lstroke = LATIN SMALL LETTER L WITH STROKE
	0x01b5: "\u013e", // lcaron = LATIN SMALL LETTER L WITH CARON
	0x01b6: "\u015b", // sacute = LATIN SMALL LETTER S WITH ACUTE
	0x01b7: "\u02c7", // caron = CARON
	0x01b9: "\u0161", // scaron = LATIN SMALL LETTER S WITH CARON
	0x01ba: "\u015f", // scedilla = LATIN SMALL LETTER S WITH CEDILLA
	0x01bb: "\u0165", // tcaron = LATIN SMALL LETTER T WITH CARON
	0x01bc: "\u017a", // zacute = LATIN SMALL LETTER Z WITH ACUTE
	0x01bd: "\u02dd", // doubleacute = DOUBLE ACUTE ACCENT
	0x01be: "\u017e", // zcaron = LATIN SMALL LETTER Z WITH CARON
	0x01bf: "\u017c", // zabovedot = LATIN SMALL LETTER Z WITH DOT ABOVE
	0x01c0: "\u0154", // Racute = LATIN CAPITAL LETTER R WITH ACUTE
	0x01c3: "\u0102", // Abreve = LATIN CAPITAL LETTER A WITH BREVE
	0x01c5: "\u0139", // Lacute = LATIN CAPITAL LETTER L WITH ACUTE
	0x01c6: "\u0106", // Cacute = LATIN CAPITAL LETTER C WITH ACUTE
	0x01c8: "\u010c", // Ccaron = LATIN CAPITAL LETTER C WITH CARON
	0x01ca: "\u0118", // Eogonek = LATIN CAPITAL LETTER E WITH OGONEK
	0x01cc: "\u011a", // Ecaron = LATIN CAPITAL LETTER E WITH CARON
	0x01cf: "\u010e", // Dcaron = LATIN CAPITAL LETTER D WITH CARON
	0x01d0: "\u0110", // Dstroke = LATIN CAPITAL LETTER D WITH STROKE
	0x01d1: "\u0143", // Nacute = LATIN CAPITAL LETTER N WITH ACUTE
	0x01d2: "\u0147", // Ncaron = LATIN CAPITAL LETTER N WITH CARON
	0x01d5: "\u0150", // Odoubleacute = LATIN CAPITAL LETTER O WITH DOUBLE ACUTE
	0x01d8: "\u0158", // Rcaron = LATIN CAPITAL LETTER R WITH CARON
	0x01d9: "\u016e", // Uring = LATIN CAPITAL LETTER U WITH RING ABOVE
	0x01db: "\u0170", // Udoubleacute = LATIN CAPITAL LETTER U WITH DOUBLE ACUTE
	0x01de: "\u0162", // Tcedilla = LATIN CAPITAL LETTER T WITH CEDILLA
	0x01e0: "\u0155", // racute = LATIN SMALL LETTER R WITH ACUTE
	0x01e3: "\u0103", // abreve = LATIN SMALL LETTER A WITH BREVE
	0x01e5: "\u013a", // lacute = LATIN SMALL LETTER L WITH ACUTE
	0x01e6: "\u0107", // cacute = LATIN SMALL LETTER C WITH ACUTE
	0x01e8: "\u010d", // ccaron = LATIN SMALL LETTER C WITH CARON
	0x01ea: "\u0119", // eogonek = LATIN SMALL LETTER E WITH OGONEK
	0x01ec: "\u011b", // ecaron = LATIN SMALL LETTER E WITH CARON
	0x01ef: "\u010f", // dcaron = LATIN SMALL LETTER D WITH CARON
	0x01f0: "\u0111", // dstroke = LATIN SMALL LETTER D WITH STROKE
	0x01f1: "\u0144", // nacute = LATIN SMALL LETTER N WITH ACUTE
	0x01f2: "\u0148", // ncaron = LATIN SMALL LETTER N WITH CARON
	0x01f5: "\u0151", // odoubleacute = LATIN SMALL LETTER O WITH DOUBLE ACUTE
	0x01f8: "\u0159", // rcaron = LATIN SMALL LETTER R WITH CARON
	0x01f9: "\u016f", // uring = LATIN SMALL LETTER U WITH RING ABOVE
	0x01fb: "\u0171", // udoubleacute = LATIN SMALL LETTER U WITH DOUBLE ACUTE
	0x01fe: "\u0163", // tcedilla = LATIN SMALL LETTER T WITH CEDILLA
	0x01ff: "\u02d9", // abovedot = DOT ABOVE
	0x02a1: "\u0126", // Hstroke = LATIN CAPITAL LETTER H WITH STROKE
	0x02a6: "\u0124", // Hcircumflex = LATIN CAPITAL LETTER H WITH CIRCUMFLEX
	0x02a9: "\u0130", // Iabovedot = LATIN CAPITAL LETTER I WITH DOT ABOVE
	0x02ab: "\u011e", // Gbreve = LATIN CAPITAL LETTER G WITH BREVE
	0x02ac: "\u0134", // Jcircumflex = LATIN CAPITAL LETTER J WITH CIRCUMFLEX
	0x02b1: "\u0127", // hstroke = LATIN SMALL LETTER H WITH STROKE
	0x02b6: "\u0125", // hcircumflex = LATIN SMALL LETTER H WITH CIRCUMFLEX
	0x02b9: "\u0131", // idotless = LATIN SMALL LETTER DOTLESS I
	0x02bb: "\u011f", // gbreve = LATIN SMALL LETTER G WITH BREVE
	0x02bc: "\u0135", // jcircumflex = LATIN SMALL LETTER J WITH CIRCUMFLEX
	0x02c5: "\u010a", // Cabovedot = LATIN CAPITAL LETTER C WITH DOT ABOVE
	0x02c6: "\u0108", // Ccircumflex = LATIN CAPITAL LETTER C WITH CIRCUMFLEX
	0x02d5: "\u0120", // Gabovedot = LATIN CAPITAL LETTER G WITH DOT ABOVE
	0x02d8: "\u011c", // Gcircumflex = LATIN CAPITAL LETTER G WITH CIRCUMFLEX
	0x02dd: "\u016c", // Ubreve = LATIN CAPITAL LETTER U WITH BREVE
	0x02de: "\u015c", // Scircumflex = LATIN CAPITAL LETTER S WITH CIRCUMFLEX
	0x02e5: "\u010b", // cabovedot = LATIN SMALL LETTER C WITH DOT ABOVE
	0x02e6: "\u0109", // ccircumflex = LATIN SMALL LETTER C WITH CIRCUMFLEX
	0x02f5: "\u0121", // gabovedot = LATIN SMALL LETTER G WITH DOT ABOVE
	0x02f8: "\u011d", // gcircumflex = LATIN SMALL LETTER G WITH CIRCUMFLEX
	0x02fd: "\u016d", // ubreve = LATIN SMALL LETTER U WITH BREVE
	0x02fe: "\u015d", // scircumflex = LATIN SMALL LETTER S WITH CIRCUMFLEX
	0x03a2: "\u0138", // kra = LATIN SMALL LETTER KRA
	0x03a3: "\u0156", // Rcedilla = LATIN CAPITAL LETTER R WITH CEDILLA
	0x03a5: "\u0128", // Itilde = LATIN CAPITAL LETTER I WITH TILDE
	0x03a6: "\u013b", // Lcedilla = LATIN CAPITAL LETTER L WITH CEDILLA
	0x03aa: "\u0112", // Emacron = LATIN CAPITAL LETTER E WITH MACRON
	0x03ab: "\u0122", // Gcedilla = LATIN CAPITAL LETTER G WITH CEDILLA
	0x03ac: "\u0166", // Tslash = LATIN CAPITAL LETTER T WITH STROKE
	0x03b3: "\u0157", // rcedilla = LATIN SMALL LETTER R WITH CEDILLA
	0x03b5: "\u0129", // itilde = LATIN SMALL LETTER I WITH TILDE
	0x03b6: "\u013c", // lcedilla = LATIN SMALL LETTER L WITH CEDILLA
	0x03ba: "\u0113", // emacron = LATIN SMALL LETTER E WITH MACRON
	0x03bb: "\u0123", // gcedilla = LATIN SMALL LETTER G WITH CEDILLA
	0x03bc: "\u0167", // tslash = LATIN SMALL LETTER T WITH STROKE
	0x03bd: "\u014a", // ENG = LATIN CAPITAL LETTER ENG
	0x03bf: "\u014b", // eng = LATIN SMALL LETTER ENG
	0x03c0: "\u0100", // Amacron = LATIN CAPITAL LETTER A WITH MACRON
	0x03c7: "\u012e", // Iogonek = LATIN CAPITAL LETTER I WITH OGONEK
	0x03cc: "\u0116", // Eabovedot = LATIN CAPITAL LETTER E WITH DOT ABOVE
	0x03cf: "\u012a", // Imacron = LATIN CAPITAL LETTER I WITH MACRON
	0x03d1: "\u0145", // Ncedilla = LATIN CAPITAL LETTER N WITH CEDILLA
	0x03d2: "\u014c", // Omacron = LATIN CAPITAL LETTER O WITH MACRON
	0x03d3: "\u0136", // Kcedilla = LATIN CAPITAL LETTER K WITH CEDILLA
	0x03d9: "\u0172", // Uogonek = LATIN CAPITAL LETTER U WITH OGONEK
	0x03dd: "\u0168", // Utilde = LATIN CAPITAL LETTER U WITH TILDE
	0x03de: "\u016a", // Umacron = LATIN CAPITAL LETTER U WITH MACRON
	0x03e0: "\u0101", // amacron = LATIN SMALL LETTER A WITH MACRON
	0x03e7: "\u012f", // iogonek = LATIN SMALL LETTER I WITH OGONEK
	0x03ec: "\u0117", // eabovedot = LATIN SMALL LETTER E WITH DOT ABOVE
	0x03ef: "\u012b", // imacron = LATIN SMALL LETTER I WITH MACRON
	0x03f1: "\u0146", // ncedilla = LATIN SMALL LETTER N WITH CEDILLA
	0x03f2: "\u014d", // omacron = LATIN SMALL LETTER O WITH MACRON
	0x03f3: "\u0137", // kcedilla = LATIN SMALL LETTER K WITH CEDILLA
	0x03f9: "\u0173", // uogonek = LATIN SMALL LETTER U WITH OGONEK
	0x03fd: "\u0169", // utilde = LATIN SMALL LETTER U WITH TILDE
	0x03fe: "\u016b", // umacron = LATIN SMALL LETTER U WITH MACRON
	0x047e: "\u203e", // overline = OVERLINE
	0x04a1: "\u3002", // kana_fullstop = IDEOGRAPHIC FULL STOP
	0x04a2: "\u300c", // kana_openingbracket = LEFT CORNER BRACKET
	0x04a3: "\u300d", // kana_closingbracket = RIGHT CORNER BRACKET
	0x04a4: "\u3001", // kana_comma = IDEOGRAPHIC COMMA
	0x04a5: "\u30fb", // kana_conjunctive = KATAKANA MIDDLE DOT
	0x04a6: "\u30f2", // kana_WO = KATAKANA LETTER WO
	0x04a7: "\u30a1", // kana_a = KATAKANA LETTER SMALL A
	0x04a8: "\u30a3", // kana_i = KATAKANA LETTER SMALL I
	0x04a9: "\u30a5", // kana_u = KATAKANA LETTER SMALL U
	0x04aa: "\u30a7", // kana_e = KATAKANA LETTER SMALL E
	0x04ab: "\u30a9", // kana_o = KATAKANA LETTER SMALL O
	0x04ac: "\u30e3", // kana_ya = KATAKANA LETTER SMALL YA
	0x04ad: "\u30e5", // kana_yu = KATAKANA LETTER SMALL YU
	0x04ae: "\u30e7", // kana_yo = KATAKANA LETTER SMALL YO
	0x04af: "\u30c3", // kana_tsu = KATAKANA LETTER SMALL TU
	0x04b0: "\u30fc", // prolongedsound = KATAKANA-HIRAGANA PROLONGED SOUND MARK
	0x04b1: "\u30a2", // kana_A = KATAKANA LETTER A
	0x04b2: "\u30a4", // kana_I = KATAKANA LETTER I
	0x04b3: "\u30a6", // kana_U = KATAKANA LETTER U
	0x04b4: "\u30a8", // kana_E = KATAKANA LETTER E
	0x04b5: "\u30aa", // kana_O = KATAKANA LETTER O
	0x04b6: "\u30ab", // kana_KA = KATAKANA LETTER KA
	0x04b7: "\u30ad", // kana_KI = KATAKANA LETTER KI
	0x04b8: "\u30af", // kana_KU = KATAKANA LETTER KU
	0x04b9: "\u30b1", // kana_KE = KATAKANA LETTER KE
	0x04ba: "\u30b3", // kana_KO = KATAKANA LETTER KO
	0x04bb: "\u30b5", // kana_SA = KATAKANA LETTER SA
	0x04bc: "\u30b7", // kana_SHI = KATAKANA LETTER SI
	0x04bd: "\u30b9", // kana_SU = KATAKANA LETTER SU
	0x04be: "\u30bb", // kana_SE = KATAKANA LETTER SE
	0x04bf: "\u30bd", // kana_SO = KATAKANA LETTER SO
	0x04c0: "\u30bf", // kana_TA = KATAKANA LETTER TA
	0x04c1: "\u30c1", // kana_CHI = KATAKANA LETTER TI
	0x04c2: "\u30c4", // kana_TSU = KATAKANA LETTER TU
	0x04c3: "\u30c6", // kana_TE = KATAKANA LETTER TE
	0x04c4: "\u30c8", // kana_TO = KATAKANA LETTER TO
	0x04c5: "\u30ca", // kana_NA = KATAKANA LETTER NA
	0x04c6: "\u30cb", // kana_NI = KATAKANA LETTER NI
	0x04c7: "\u30cc", // kana_NU = KATAKANA LETTER NU
	0x04c8: "\u30cd", // kana_NE = KATAKANA LETTER NE
	0x04c9: "\u30ce", // kana_NO = KATAKANA LETTER NO
	0x04ca: "\u30cf", // kana_HA = KATAKANA LETTER HA
	0x04cb: "\u30d2", // kana_HI = KATAKANA LETTER HI
	0x04cc: "\u30d5", // kana_FU = KATAKANA LETTER HU
	0x04cd: "\u30d8", // kana_HE = KATAKANA LETTER HE
	0x04ce: "\u30db", // kana_HO = KATAKANA LETTER HO
	0x04cf: "\u30de", // kana_MA = KATAKANA LETTER MA
	0x04d0: "\u30df", // kana_MI = KATAKANA LETTER MI
	0x04d1: "\u30e0", // kana_MU = KATAKANA LETTER MU
	0x04d2: "\u30e1", // kana_ME = KATAKANA LETTER ME
	0x04d3: "\u30e2", // kana_MO = KATAKANA LETTER MO
	0x04d4: "\u30e4", // kana_YA = KATAKANA LETTER YA
	0x04d5: "\u30e6", // kana_YU = KATAKANA LETTER YU
	0x04d6: "\u30e8", // kana_YO = KATAKANA LETTER YO
	0x04d7: "\u30e9", // kana_RA = KATAKANA LETTER RA
	0x04d8: "\u30ea", // kana_RI = KATAKANA LETTER RI
	0x04d9: "\u30eb", // kana_RU = KATAKANA LETTER RU
	0x04da: "\u30ec", // kana_RE = KATAKANA LETTER RE
	0x04db: "\u30ed", // kana_RO = KATAKANA LETTER RO
	0x04dc: "\u30ef", // kana_WA = KATAKANA LETTER WA
	0x04dd: "\u30f3", // kana_N = KATAKANA LETTER N
	0x04de: "\u309b", // voicedsound = KATAKANA-HIRAGANA VOICED SOUND MARK
	0x04df: "\u309c", // semivoicedsound = KATAKANA-HIRAGANA SEMI-VOICED SOUND MARK
	0x05ac: "\u060c", // Arabic_comma = ARABIC COMMA
	0x05bb: "\u061b", // Arabic_semicolon = ARABIC SEMICOLON
	0x05bf: "\u061f", // Arabic_question_mark = ARABIC QUESTION MARK
	0x05c1: "\u0621", // Arabic_hamza = ARABIC LETTER HAMZA
	0x05c2: "\u0622", // Arabic_maddaonalef = ARABIC LETTER ALEF WITH MADDA ABOVE
	0x05c3: "\u0623", // Arabic_hamzaonalef = ARABIC LETTER ALEF WITH HAMZA ABOVE
	0x05c4: "\u0624", // Arabic_hamzaonwaw = ARABIC LETTER WAW WITH HAMZA ABOVE
	0x05c5: "\u0625", // Arabic_hamzaunderalef = ARABIC LETTER ALEF WITH HAMZA BELOW
	0x05c6: "\u0626", // Arabic_hamzaonyeh = ARABIC LETTER YEH WITH HAMZA ABOVE
	0x05c7: "\u0627", // Arabic_alef = ARABIC LETTER ALEF
	0x05c8: "\u0628", // Arabic_beh = ARABIC LETTER BEH
	0x05c9: "\u0629", // Arabic_tehmarbuta = ARABIC LETTER TEH MARBUTA
	0x05ca: "\u062a", // Arabic_teh = ARABIC LETTER TEH
	0x05cb: "\u062b", // Arabic_theh = ARABIC LETTER THEH
	0x05cc: "\u062c", // Arabic_jeem = ARABIC LETTER JEEM
	0x05cd: "\u062d", // Arabic_hah = ARABIC LETTER HAH
	0x05ce: "\u062e", // Arabic_khah = ARABIC LETTER KHAH
	0x05cf: "\u062f", // Arabic_dal = ARABIC LETTER DAL
	0x05d0: "\u0630", // Arabic_thal = ARABIC LETTER THAL
	0x05d1: "\u0631", // Arabic_ra = ARABIC LETTER REH
	0x05d2: "\u0632", // Arabic_zain = ARABIC LETTER ZAIN
	0x05d3: "\u0633", // Arabic_seen = ARABIC LETTER SEEN
	0x05d4: "\u0634", // Arabic_sheen = ARABIC LETTER SHEEN
	0x05d5: "\u0635", // Arabic_sad = ARABIC LETTER SAD
	0x05d6: "\u0636", // Arabic_dad = ARABIC LETTER DAD
	0x05d7: "\u0637", // Arabic_tah = ARABIC LETTER TAH
	0x05d8: "\u0638", // Arabic_zah = ARABIC LETTER ZAH
	0x05d9: "\u0639", // Arabic_ain = ARABIC LETTER AIN
	0x05da: "\u063a", // Arabic_ghain = ARABIC LETTER GHAIN
	0x05e0: "\u0640", // Arabic_tatweel = ARABIC TATWEEL
	0x05e1: "\u0641", // Arabic_feh = ARABIC LETTER FEH
	0x05e2: "\u0642", // Arabic_qaf = ARABIC LETTER QAF
	0x05e3: "\u0643", // Arabic_kaf = ARABIC LETTER KAF
	0x05e4: "\u0644", // Arabic_lam = ARABIC LETTER LAM
	0x05e5: "\u0645", // Arabic_meem = ARABIC LETTER MEEM
	0x05e6: "\u0646", // Arabic_noon = ARABIC LETTER NOON
	0x05e7: "\u0647", // Arabic_ha = ARABIC LETTER HEH
	0x05e8: "\u0648", // Arabic_waw = ARABIC LETTER WAW
	0x05e9: "\u0649", // Arabic_alefmaksura = ARABIC LETTER ALEF MAKSURA
	0x05ea: "\u064a", // Arabic_yeh = ARABIC LETTER YEH
	0x05eb: "\u064b", // Arabic_fathatan = ARABIC FATHATAN
	0x05ec: "\u064c", // Arabic_dammatan = ARABIC DAMMATAN
	0x05ed: "\u064d", // Arabic_kasratan = ARABIC KASRATAN
	0x05ee: "\u064e", // Arabic_fatha = ARABIC FATHA
	0x05ef: "\u064f", // Arabic_damma = ARABIC DAMMA
	0x05f0: "\u0650", // Arabic_kasra = ARABIC KASRA
	0x05f1: "\u0651", // Arabic_shadda = ARABIC SHADDA
	0x05f2: "\u0652", // Arabic_sukun = ARABIC SUKUN
	0x06a1: "\u0452", // Serbian_dje = CYRILLIC SMALL LETTER DJE
	0x06a2: "\u0453", // Macedonia_gje = CYRILLIC SMALL LETTER GJE
	0x06a3: "\u0451", // Cyrillic_io = CYRILLIC SMALL LETTER IO
	0x06a4: "\u0454", // Ukrainian_ie = CYRILLIC SMALL LETTER UKRAINIAN IE
	0x06a5: "\u0455", // Macedonia_dse = CYRILLIC SMALL LETTER DZE
	0x06a6: "\u0456", // Ukrainian_i = CYRILLIC SMALL LETTER BYELORUSSIAN-UKRAINIAN I
	0x06a7: "\u0457", // Ukrainian_yi = CYRILLIC SMALL LETTER YI
	0x06a8: "\u0458", // Cyrillic_je = CYRILLIC SMALL LETTER JE
	0x06a9: "\u0459", // Cyrillic_lje = CYRILLIC SMALL LETTER LJE
	0x06aa: "\u045a", // Cyrillic_nje = CYRILLIC SMALL LETTER NJE
	0x06ab: "\u045b", // Serbian_tshe = CYRILLIC SMALL LETTER TSHE
	0x06ac: "\u045c", // Macedonia_kje = CYRILLIC SMALL LETTER KJE
	0x06ae: "\u045e", // Byelorussian_shortu = CYRILLIC SMALL LETTER SHORT U
	0x06af: "\u045f", // Cyrillic_dzhe = CYRILLIC SMALL LETTER DZHE
	0x06b0: "\u2116", // numerosign = NUMERO SIGN
	0x06b1: "\u0402", // Serbian_DJE = CYRILLIC CAPITAL LETTER DJE
	0x06b2: "\u0403", // Macedonia_GJE = CYRILLIC CAPITAL LETTER GJE
	0x06b3: "\u0401", // Cyrillic_IO = CYRILLIC CAPITAL LETTER IO
	0x06b4: "\u0404", // Ukrainian_IE = CYRILLIC CAPITAL LETTER UKRAINIAN IE
	0x06b5: "\u0405", // Macedonia_DSE = CYRILLIC CAPITAL LETTER DZE
	0x06b6: "\u0406", // Ukrainian_I = CYRILLIC CAPITAL LETTER BYELORUSSIAN-UKRAINIAN I
	0x06b7: "\u0407", // Ukrainian_YI = CYRILLIC CAPITAL LETTER YI
	0x06b8: "\u0408", // Cyrillic_JE = CYRILLIC CAPITAL LETTER JE
	0x06b9: "\u0409", // Cyrillic_LJE = CYRILLIC CAPITAL LETTER LJE
	0x06ba: "\u040a", // Cyrillic_NJE = CYRILLIC CAPITAL LETTER NJE
	0x06bb: "\u040b", // Serbian_TSHE = CYRILLIC CAPITAL LETTER TSHE
	0x06bc: "\u040c", // Macedonia_KJE = CYRILLIC CAPITAL LETTER KJE
	0x06be: "\u040e", // Byelorussian_SHORTU = CYRILLIC CAPITAL LETTER SHORT U
	0x06bf: "\u040f", // Cyrillic_DZHE = CYRILLIC CAPITAL LETTER DZHE
	0x06c0: "\u044e", // Cyrillic_yu = CYRILLIC SMALL LETTER YU
	0x06c1: "\u0430", // Cyrillic_a = CYRILLIC SMALL LETTER A
	0x06c2: "\u0431", // Cyrillic_be = CYRILLIC SMALL LETTER BE
	0x06c3: "\u0446", // Cyrillic_tse = CYRILLIC SMALL LETTER TSE
	0x06c4: "\u0434", // Cyrillic_de = CYRILLIC SMALL LETTER DE
	0x06c5: "\u0435", // Cyrillic_ie = CYRILLIC SMALL LETTER IE
	0x06c6: "\u0444", // Cyrillic_ef = CYRILLIC SMALL LETTER EF
	0x06c7: "\u0433", // Cyrillic_ghe = CYRILLIC SMALL LETTER GHE
	0x06c8: "\u0445", // Cyrillic_ha = CYRILLIC SMALL LETTER HA
	0x06c9: "\u0438", // Cyrillic_i = CYRILLIC SMALL LETTER I
	0x06ca: "\u0439", // Cyrillic_shorti = CYRILLIC SMALL LETTER SHORT I
	0x06cb: "\u043a", // Cyrillic_ka = CYRILLIC SMALL LETTER KA
	0x06cc: "\u043b", // Cyrillic_el = CYRILLIC SMALL LETTER EL
	0x06cd: "\u043c", // Cyrillic_em = CYRILLIC SMALL LETTER EM
	0x06ce: "\u043d", // Cyrillic_en = CYRILLIC SMALL LETTER EN
	0x06cf: "\u043e", // Cyrillic_o = CYRILLIC SMALL LETTER O
	0x06d0: "\u043f", // Cyrillic_pe = CYRILLIC SMALL LETTER PE
	0x06d1: "\u044f", // Cyrillic_ya = CYRILLIC SMALL LETTER YA
	0x06d2: "\u0440", // Cyrillic_er = CYRILLIC SMALL LETTER ER
	0x06d3: "\u0441", // Cyrillic_es = CYRILLIC SMALL LETTER ES
	0x06d4: "\u0442", // Cyrillic_te = CYRILLIC SMALL LETTER TE
	0x06d5: "\u0443", // Cyrillic_u = CYRILLIC SMALL LETTER U
	0x06d6: "\u0436", // Cyrillic_zhe = CYRILLIC SMALL LETTER ZHE
	0x06d7: "\u0432", // Cyrillic_ve = CYRILLIC SMALL LETTER VE
	0x06d8: "\u044c", // Cyrillic_softsign = CYRILLIC SMALL LETTER SOFT SIGN
	0x06d9: "\u044b", // Cyrillic_yeru = CYRILLIC SMALL LETTER YERU
	0x06da: "\u0437", // Cyrillic_ze = CYRILLIC SMALL LETTER ZE
	0x06db: "\u0448", // Cyrillic_sha = CYRILLIC SMALL LETTER SHA
	0x06dc: "\u044d", // Cyrillic_e = CYRILLIC SMALL LETTER E
	0x06dd: "\u0449", // Cyrillic_shcha = CYRILLIC SMALL LETTER SHCHA
	0x06de: "\u0447", // Cyrillic_che = CYRILLIC SMALL LETTER CHE
	0x06df: "\u044a", // Cyrillic_hardsign = CYRILLIC SMALL LETTER HARD SIGN
	0x06e0: "\u042e", // Cyrillic_YU = CYRILLIC CAPITAL LETTER YU
	0x06e1: "\u0410", // Cyrillic_A = CYRILLIC CAPITAL LETTER A
	0x06e2: "\u0411", // Cyrillic_BE = CYRILLIC CAPITAL LETTER BE
	0x06e3: "\u0426", // Cyrillic_TSE = CYRILLIC CAPITAL LETTER TSE
	0x06e4: "\u0414", // Cyrillic_DE = CYRILLIC CAPITAL LETTER DE
	0x06e5: "\u0415", // Cyrillic_IE = CYRILLIC CAPITAL LETTER IE
	0x06e6: "\u0424", // Cyrillic_EF = CYRILLIC CAPITAL LETTER EF
	0x06e7: "\u0413", // Cyrillic_GHE = CYRILLIC CAPITAL LETTER GHE
	0x06e8: "\u0425", // Cyrillic_HA = CYRILLIC CAPITAL LETTER HA
	0x06e9: "\u0418", // Cyrillic_I = CYRILLIC CAPITAL LETTER I
	0x06ea: "\u0419", // Cyrillic_SHORTI = CYRILLIC CAPITAL LETTER SHORT I
	0x06eb: "\u041a", // Cyrillic_KA = CYRILLIC CAPITAL LETTER KA
	0x06ec: "\u041b", // Cyrillic_EL = CYRILLIC CAPITAL LETTER EL
	0x06ed: "\u041c", // Cyrillic_EM = CYRILLIC CAPITAL LETTER EM
	0x06ee: "\u041d", // Cyrillic_EN = CYRILLIC CAPITAL LETTER EN
	0x06ef: "\u041e", // Cyrillic_O = CYRILLIC CAPITAL LETTER O
	0x06f0: "\u041f", // Cyrillic_PE = CYRILLIC CAPITAL LETTER PE
	0x06f1: "\u042f", // Cyrillic_YA = CYRILLIC CAPITAL LETTER YA
	0x06f2: "\u0420", // Cyrillic_ER = CYRILLIC CAPITAL LETTER ER
	0x06f3: "\u0421", // Cyrillic_ES = CYRILLIC CAPITAL LETTER ES
	0x06f4: "\u0422", // Cyrillic_TE = CYRILLIC CAPITAL LETTER TE
	0x06f5: "\u0423", // Cyrillic_U = CYRILLIC CAPITAL LETTER U
	0x06f6: "\u0416", // Cyrillic_ZHE = CYRILLIC CAPITAL LETTER ZHE
	0x06f7: "\u0412", // Cyrillic_VE = CYRILLIC CAPITAL LETTER VE
	0x06f8: "\u042c", // Cyrillic_SOFTSIGN = CYRILLIC CAPITAL LETTER SOFT SIGN
	0x06f9: "\u042b", // Cyrillic_YERU = CYRILLIC CAPITAL LETTER YERU
	0x06fa: "\u0417", // Cyrillic_ZE = CYRILLIC CAPITAL LETTER ZE
	0x06fb: "\u0428", // Cyrillic_SHA = CYRILLIC CAPITAL LETTER SHA
	0x06fc: "\u042d", // Cyrillic_E = CYRILLIC CAPITAL LETTER E
	0x06fd: "\u0429", // Cyrillic_SHCHA = CYRILLIC CAPITAL LETTER SHCHA
	0x06fe: "\u0427", // Cyrillic_CHE = CYRILLIC CAPITAL LETTER CHE
	0x06ff: "\u042a", // Cyrillic_HARDSIGN = CYRILLIC CAPITAL LETTER HARD SIGN
	0x07a1: "\u0386", // Greek_ALPHAaccent = GREEK CAPITAL LETTER ALPHA WITH TONOS
	0x07a2: "\u0388", // Greek_EPSILONaccent = GREEK CAPITAL LETTER EPSILON WITH TONOS
	0x07a3: "\u0389", // Greek_ETAaccent = GREEK CAPITAL LETTER ETA WITH TONOS
	0x07a4: "\u038a", // Greek_IOTAaccent = GREEK CAPITAL LETTER IOTA WITH TONOS
	0x07a5: "\u03aa", // Greek_IOTAdiaeresis = GREEK CAPITAL LETTER IOTA WITH DIALYTIKA
	0x07a7: "\u038c", // Greek_OMICRONaccent = GREEK CAPITAL LETTER OMICRON WITH TONOS
	0x07a8: "\u038e", // Greek_UPSILONaccent = GREEK CAPITAL LETTER UPSILON WITH TONOS
	0x07a9: "\u03ab", // Greek_UPSILONdieresis = GREEK CAPITAL LETTER UPSILON WITH DIALYTIKA
	0x07ab: "\u038f", // Greek_OMEGAaccent = GREEK CAPITAL LETTER OMEGA WITH TONOS
	0x07ae: "\u0385", // Greek_accentdieresis = GREEK DIALYTIKA TONOS
	0x07af: "\u2015", // Greek_horizbar = HORIZONTAL BAR
	0x07b1: "\u03ac", // Greek_alphaaccent = GREEK SMALL LETTER ALPHA WITH TONOS
	0x07b2: "\u03ad", // Greek_epsilonaccent = GREEK SMALL LETTER EPSILON WITH TONOS
	0x07b3: "\u03ae", // Greek_etaaccent = GREEK SMALL LETTER ETA WITH TONOS
	0x07b4: "\u03af", // Greek_iotaaccent = GREEK SMALL LETTER IOTA WITH TONOS
	0x07b5: "\u03ca", // Greek_iotadieresis = GREEK SMALL LETTER IOTA WITH DIALYTIKA
	0x07b6: "\u0390", // Greek_iotaaccentdieresis = GREEK SMALL LETTER IOTA WITH DIALYTIKA AND TONOS
	0x07b7: "\u03cc", // Greek_omicronaccent = GREEK SMALL LETTER OMICRON WITH TONOS
	0x07b8: "\u03cd", // Greek_upsilonaccent = GREEK SMALL LETTER UPSILON WITH TONOS
	0x07b9: "\u03cb", // Greek_upsilondieresis = GREEK SMALL LETTER UPSILON WITH DIALYTIKA
	0x07ba: "\u03b0", // Greek_upsilonaccentdieresis = GREEK SMALL LETTER UPSILON WITH DIALYTIKA AND TONOS
	0x07bb: "\u03ce", // Greek_omegaaccent = GREEK SMALL LETTER OMEGA WITH TONOS
	0x07c1: "\u0391", // Greek_ALPHA = GREEK CAPITAL LETTER ALPHA
	0x07c2: "\u0392", // Greek_BETA = GREEK CAPITAL LETTER BETA
	0x07c3: "\u0393", // Greek_GAMMA = GREEK CAPITAL LETTER GAMMA
	0x07c4: "\u0394", // Greek_DELTA = GREEK CAPITAL LETTER DELTA
	0x07c5: "\u0395", // Greek_EPSILON = GREEK CAPITAL LETTER EPSILON
	0x07c6: "\u0396", // Greek_ZETA = GREEK CAPITAL LETTER ZETA
	0x07c7: "\u0397", // Greek_ETA = GREEK CAPITAL LETTER ETA
	0x07c8: "\u0398", // Greek_THETA = GREEK CAPITAL LETTER THETA
	0x07c9: "\u0399", // Greek_IOTA = GREEK CAPITAL LETTER IOTA
	0x07ca: "\u039a", // Greek_KAPPA = GREEK CAPITAL LETTER KAPPA
	0x07cb: "\u039b", // Greek_LAMBDA = GREEK CAPITAL LETTER LAMDA
	0x07cc: "\u039c", // Greek_MU = GREEK CAPITAL LETTER MU
	0x07cd: "\u039d", // Greek_NU = GREEK CAPITAL LETTER NU
	0x07ce: "\u039e", // Greek_XI = GREEK CAPITAL LETTER XI
	0x07cf: "\u039f", // Greek_OMICRON = GREEK CAPITAL LETTER OMICRON
	0x07d0: "\u03a0", // Greek_PI = GREEK CAPITAL LETTER PI
	0x07d1: "\u03a1", // Greek_RHO = GREEK CAPITAL LETTER RHO
	0x07d2: "\u03a3", // Greek_SIGMA = GREEK CAPITAL LETTER SIGMA
	0x07d4: "\u03a4", // Greek_TAU = GREEK CAPITAL LETTER TAU
	0x07d5: "\u03a5", // Greek_UPSILON = GREEK CAPITAL LETTER UPSILON
	0x07d6: "\u03a6", // Greek_PHI = GREEK CAPITAL LETTER PHI
	0x07d7: "\u03a7", // Greek_CHI = GREEK CAPITAL LETTER CHI
	0x07d8: "\u03a8", // Greek_PSI = GREEK CAPITAL LETTER PSI
	0x07d9: "\u03a9", // Greek_OMEGA = GREEK CAPITAL LETTER OMEGA
	0x07e1: "\u03b1", // Greek_alpha = GREEK SMALL LETTER ALPHA
	0x07e2: "\u03b2", // Greek_beta = GREEK SMALL LETTER BETA
	0x07e3: "\u03b3", // Greek_gamma = GREEK SMALL LETTER GAMMA
	0x07e4: "\u03b4", // Greek_delta = GREEK SMALL LETTER DELTA
	0x07e5: "\u03b5", // Greek_epsilon = GREEK SMALL LETTER EPSILON
	0x07e6: "\u03b6", // Greek_zeta = GREEK SMALL LETTER ZETA
	0x07e7: "\u03b7", // Greek_eta = GREEK SMALL LETTER ETA
	0x07e8: "\u03b8", // Greek_theta = GREEK SMALL LETTER THETA
	0x07e9: "\u03b9", // Greek_iota = GREEK SMALL LETTER IOTA
	0x07ea: "\u03ba", // Greek_kappa = GREEK SMALL LETTER KAPPA
	0x07eb: "\u03bb", // Greek_lambda = GREEK SMALL LETTER LAMDA
	0x07ec: "\u03bc", // Greek_mu = GREEK SMALL LETTER MU
	0x07ed: "\u03bd", // Greek_nu = GREEK SMALL LETTER NU
	0x07ee: "\u03be", // Greek_xi = GREEK SMALL LETTER XI
	0x07ef: "\u03bf", // Greek_omicron = GREEK SMALL LETTER OMICRON
	0x07f0: "\u03c0", // Greek_pi = GREEK SMALL LETTER PI
	0x07f1: "\u03c1", // Greek_rho = GREEK SMALL LETTER RHO
	0x07f2: "\u03c3", // Greek_sigma = GREEK SMALL LETTER SIGMA
	0x07f3: "\u03c2", // Greek_finalsmallsigma = GREEK SMALL LETTER FINAL SIGMA
	0x07f4: "\u03c4", // Greek_tau = GREEK SMALL LETTER TAU
	0x07f5: "\u03c5", // Greek_upsilon = GREEK SMALL LETTER UPSILON
	0x07f6: "\u03c6", // Greek_phi = GREEK SMALL LETTER PHI
	0x07f7: "\u03c7", // Greek_chi = GREEK SMALL LETTER CHI
	0x07f8: "\u03c8", // Greek_psi = GREEK SMALL LETTER PSI
	0x07f9: "\u03c9", // Greek_omega = GREEK SMALL LETTER OMEGA
	0x08a1: "\u23b7", // leftradical = ???
	0x08a2: "\u250c", // topleftradical = BOX DRAWINGS LIGHT DOWN AND RIGHT
	0x08a3: "\u2500", // horizconnector = BOX DRAWINGS LIGHT HORIZONTAL
	0x08a4: "\u2320", // topintegral = TOP HALF INTEGRAL
	0x08a5: "\u2321", // botintegral = BOTTOM HALF INTEGRAL
	0x08a6: "\u2502", // vertconnector = BOX DRAWINGS LIGHT VERTICAL
	0x08a7: "\u23a1", // topleftsqbracket = ???
	0x08a8: "\u23a3", // botleftsqbracket = ???
	0x08a9: "\u23a4", // toprightsqbracket = ???
	0x08aa: "\u23a6", // botrightsqbracket = ???
	0x08ab: "\u239b", // topleftparens = ???
	0x08ac: "\u239d", // botleftparens = ???
	0x08ad: "\u239e", // toprightparens = ???
	0x08ae: "\u23a0", // botrightparens = ???
	0x08af: "\u23a8", // leftmiddlecurlybrace = ???
	0x08b0: "\u23ac", // rightmiddlecurlybrace = ???
	0x08bc: "\u2264", // lessthanequal = LESS-THAN OR EQUAL TO
	0x08bd: "\u2260", // notequal = NOT EQUAL TO
	0x08be: "\u2265", // greaterthanequal = GREATER-THAN OR EQUAL TO
	0x08bf: "\u222b", // integral = INTEGRAL
	0x08c0: "\u2234", // therefore = THEREFORE
	0x08c1: "\u221d", // variation = PROPORTIONAL TO
	0x08c2: "\u221e", // infinity = INFINITY
	0x08c5: "\u2207", // nabla = NABLA
	0x08c8: "\u223c", // approximate = TILDE OPERATOR
	0x08c9: "\u2243", // similarequal = ASYMPTOTICALLY EQUAL TO
	0x08cd: "\u21d4", // ifonlyif = LEFT RIGHT DOUBLE ARROW
	0x08ce: "\u21d2", // implies = RIGHTWARDS DOUBLE ARROW
	0x08cf: "\u2261", // identical = IDENTICAL TO
	0x08d6: "\u221a", // radical = SQUARE ROOT
	0x08da: "\u2282", // includedin = SUBSET OF
	0x08db: "\u2283", // includes = SUPERSET OF
	0x08dc: "\u2229", // intersection = INTERSECTION
	0x08dd: "\u222a", // union = UNION
	0x08de: "\u2227", // logicaland = LOGICAL AND
	0x08df: "\u2228", // logicalor = LOGICAL OR
	0x08ef: "\u2202", // partialderivative = PARTIAL DIFFERENTIAL
	0x08f6: "\u0192", // function = LATIN SMALL LETTER F WITH HOOK
	0x08fb: "\u2190", // leftarrow = LEFTWARDS ARROW
	0x08fc: "\u2191", // uparrow = UPWARDS ARROW
	0x08fd: "\u2192", // rightarrow = RIGHTWARDS ARROW
	0x08fe: "\u2193", // downarrow = DOWNWARDS ARROW
	0x09e0: "\u25c6", // soliddiamond = BLACK DIAMOND
	0x09e1: "\u2592", // checkerboard = MEDIUM SHADE
	0x09e2: "\u2409", // ht = SYMBOL FOR HORIZONTAL TABULATION
	0x09e3: "\u240c", // ff = SYMBOL FOR FORM FEED
	0x09e4: "\u240d", // cr = SYMBOL FOR CARRIAGE RETURN
	0x09e5: "\u240a", // lf = SYMBOL FOR LINE FEED
	0x09e8: "\u2424", // nl = SYMBOL FOR NEWLINE
	0x09e9: "\u240b", // vt = SYMBOL FOR VERTICAL TABULATION
	0x09ea: "\u2518", // lowrightcorner = BOX DRAWINGS LIGHT UP AND LEFT
	0x09eb: "\u2510", // uprightcorner = BOX DRAWINGS LIGHT DOWN AND LEFT
	0x09ec: "\u250c", // upleftcorner = BOX DRAWINGS LIGHT DOWN AND RIGHT
	0x09ed: "\u2514", // lowleftcorner = BOX DRAWINGS LIGHT UP AND RIGHT
	0x09ee: "\u253c", // crossinglines = BOX DRAWINGS LIGHT VERTICAL AND HORIZONTAL
	0x09ef: "\u23ba", // horizlinescan1 = HORIZONTAL SCAN LINE-1 (Unicode 3.2 draft)
	0x09f0: "\u23bb", // horizlinescan3 = HORIZONTAL SCAN LINE-3 (Unicode 3.2 draft)
	0x09f1: "\u2500", // horizlinescan5 = BOX DRAWINGS LIGHT HORIZONTAL
	0x09f2: "\u23bc", // horizlinescan7 = HORIZONTAL SCAN LINE-7 (Unicode 3.2 draft)
	0x09f3: "\u23bd", // horizlinescan9 = HORIZONTAL SCAN LINE-9 (Unicode 3.2 draft)
	0x09f4: "\u251c", // leftt = BOX DRAWINGS LIGHT VERTICAL AND RIGHT
	0x09f5: "\u2524", // rightt = BOX DRAWINGS LIGHT VERTICAL AND LEFT
	0x09f6: "\u2534", // bott = BOX DRAWINGS LIGHT UP AND HORIZONTAL
	0x09f7: "\u252c", // topt = BOX DRAWINGS LIGHT DOWN AND HORIZONTAL
	0x09f8: "\u2502", // vertbar = BOX DRAWINGS LIGHT VERTICAL
	0x0aa1: "\u2003", // emspace = EM SPACE
	0x0aa2: "\u2002", // enspace = EN SPACE
	0x0aa3: "\u2004", // em3space = THREE-PER-EM SPACE
	0x0aa4: "\u2005", // em4space = FOUR-PER-EM SPACE
	0x0aa5: "\u2007", // digitspace = FIGURE SPACE
	0x0aa6: "\u2008", // punctspace = PUNCTUATION SPACE
	0x0aa7: "\u2009", // thinspace = THIN SPACE
	0x0aa8: "\u200a", // hairspace = HAIR SPACE
	0x0aa9: "\u2014", // emdash = EM DASH
	0x0aaa: "\u2013", // endash = EN DASH
	0x0aae: "\u2026", // ellipsis = HORIZONTAL ELLIPSIS
	0x0aaf: "\u2025", // doubbaselinedot = TWO DOT LEADER
	0x0ab0: "\u2153", // onethird = VULGAR FRACTION ONE THIRD
	0x0ab1: "\u2154", // twothirds = VULGAR FRACTION TWO THIRDS
	0x0ab2: "\u2155", // onefifth = VULGAR FRACTION ONE FIFTH
	0x0ab3: "\u2156", // twofifths = VULGAR FRACTION TWO FIFTHS
	0x0ab4: "\u2157", // threefifths = VULGAR FRACTION THREE FIFTHS
	0x0ab5: "\u2158", // fourfifths = VULGAR FRACTION FOUR FIFTHS
	0x0ab6: "\u2159", // onesixth = VULGAR FRACTION ONE SIXTH
	0x0ab7: "\u215a", // fivesixths = VULGAR FRACTION FIVE SIXTHS
	0x0ab8: "\u2105", // careof = CARE OF
	0x0abb: "\u2012", // figdash = FIGURE DASH
	0x0abc: "\u2329", // leftanglebracket = LEFT-POINTING ANGLE BRACKET
	0x0abe: "\u232a", // rightanglebracket = RIGHT-POINTING ANGLE BRACKET
	0x0ac3: "\u215b", // oneeighth = VULGAR FRACTION ONE EIGHTH
	0x0ac4: "\u215c", // threeeighths = VULGAR FRACTION THREE EIGHTHS
	0x0ac5: "\u215d", // fiveeighths = VULGAR FRACTION FIVE EIGHTHS
	0x0ac6: "\u215e", // seveneighths = VULGAR FRACTION SEVEN EIGHTHS
	0x0ac9: "\u2122", // trademark = TRADE MARK SIGN
	0x0aca: "\u2613", // signaturemark = SALTIRE
	0x0acc: "\u25c1", // leftopentriangle = WHITE LEFT-POINTING TRIANGLE
	0x0acd: "\u25b7", // rightopentriangle = WHITE RIGHT-POINTING TRIANGLE
	0x0ace: "\u25cb", // emopencircle = WHITE CIRCLE
	0x0acf: "\u25af", // emopenrectangle = WHITE VERTICAL RECTANGLE
	0x0ad0: "\u2018", // leftsinglequotemark = LEFT SINGLE QUOTATION MARK
	0x0ad1: "\u2019", // rightsinglequotemark = RIGHT SINGLE QUOTATION MARK
	0x0ad2: "\u201c", // leftdoublequotemark = LEFT DOUBLE QUOTATION MARK
	0x0ad3: "\u201d", // rightdoublequotemark = RIGHT DOUBLE QUOTATION MARK
	0x0ad4: "\u211e", // prescription = PRESCRIPTION TAKE
	0x0ad6: "\u2032", // minutes = PRIME
	0x0ad7: "\u2033", // seconds = DOUBLE PRIME
	0x0ad9: "\u271d", // latincross = LATIN CROSS
	0x0adb: "\u25ac", // filledrectbullet = BLACK RECTANGLE
	0x0adc: "\u25c0", // filledlefttribullet = BLACK LEFT-POINTING TRIANGLE
	0x0add: "\u25b6", // filledrighttribullet = BLACK RIGHT-POINTING TRIANGLE
	0x0ade: "\u25cf", // emfilledcircle = BLACK CIRCLE
	0x0adf: "\u25ae", // emfilledrect = BLACK VERTICAL RECTANGLE
	0x0ae0: "\u25e6", // enopencircbullet = WHITE BULLET
	0x0ae1: "\u25ab", // enopensquarebullet = WHITE SMALL SQUARE
	0x0ae2: "\u25ad", // openrectbullet = WHITE RECTANGLE
	0x0ae3: "\u25b3", // opentribulletup = WHITE UP-POINTING TRIANGLE
	0x0ae4: "\u25bd", // opentribulletdown = WHITE DOWN-POINTING TRIANGLE
	0x0ae5: "\u2606", // openstar = WHITE STAR
	0x0ae6: "\u2022", // enfilledcircbullet = BULLET
	0x0ae7: "\u25aa", // enfilledsqbullet = BLACK SMALL SQUARE
	0x0ae8: "\u25b2", // filledtribulletup = BLACK UP-POINTING TRIANGLE
	0x0ae9: "\u25bc", // filledtribulletdown = BLACK DOWN-POINTING TRIANGLE
	0x0aea: "\u261c", // leftpointer = WHITE LEFT POINTING INDEX
	0x0aeb: "\u261e", // rightpointer = WHITE RIGHT POINTING INDEX
	0x0aec: "\u2663", // club = BLACK CLUB SUIT
	0x0aed: "\u2666", // diamond = BLACK DIAMOND SUIT
	0x0aee: "\u2665", // heart = BLACK HEART SUIT
	0x0af0: "\u2720", // maltesecross = MALTESE CROSS
	0x0af1: "\u2020", // dagger = DAGGER
	0x0af2: "\u2021", // doubledagger = DOUBLE DAGGER
	0x0af3: "\u2713", // checkmark = CHECK MARK
	0x0af4: "\u2717", // ballotcross = BALLOT X
	0x0af5: "\u266f", // musicalsharp = MUSIC SHARP SIGN
	0x0af6: "\u266d", // musicalflat = MUSIC FLAT SIGN
	0x0af7: "\u2642", // malesymbol = MALE SIGN
	0x0af8: "\u2640", // femalesymbol = FEMALE SIGN
	0x0af9: "\u260e", // telephone = BLACK TELEPHONE
	0x0afa: "\u2315", // telephonerecorder = TELEPHONE RECORDER
	0x0afb: "\u2117", // phonographcopyright = SOUND RECORDING COPYRIGHT
	0x0afc: "\u2038", // caret = CARET
	0x0afd: "\u201a", // singlelowquotemark = SINGLE LOW-9 QUOTATION MARK
	0x0afe: "\u201e", // doublelowquotemark = DOUBLE LOW-9 QUOTATION MARK
	0x0ba3: "\u003c", // leftcaret = LESS-THAN SIGN
	0x0ba6: "\u003e", // rightcaret = GREATER-THAN SIGN
	0x0ba8: "\u2228", // downcaret = LOGICAL OR
	0x0ba9: "\u2227", // upcaret = LOGICAL AND
	0x0bc0: "\u00af", // overbar = MACRON
	0x0bc2: "\u22a5", // downtack = UP TACK
	0x0bc3: "\u2229", // upshoe = INTERSECTION
	0x0bc4: "\u230a", // downstile = LEFT FLOOR
	0x0bc6: "\u005f", // underbar = LOW LINE
	0x0bca: "\u2218", // jot = RING OPERATOR
	0x0bcc: "\u2395", // quad = APL FUNCTIONAL SYMBOL QUAD
	0x0bce: "\u22a4", // uptack = DOWN TACK
	0x0bcf: "\u25cb", // circle = WHITE CIRCLE
	0x0bd3: "\u2308", // upstile = LEFT CEILING
	0x0bd6: "\u222a", // downshoe = UNION
	0x0bd8: "\u2283", // rightshoe = SUPERSET OF
	0x0bda: "\u2282", // leftshoe = SUBSET OF
	0x0bdc: "\u22a2", // lefttack = RIGHT TACK
	0x0bfc: "\u22a3", // righttack = LEFT TACK
	0x0cdf: "\u2017", // hebrew_doublelowline = DOUBLE LOW LINE
	0x0ce0: "\u05d0", // hebrew_aleph = HEBREW LETTER ALEF
	0x0ce1: "\u05d1", // hebrew_bet = HEBREW LETTER BET
	0x0ce2: "\u05d2", // hebrew_gimel = HEBREW LETTER GIMEL
	0x0ce3: "\u05d3", // hebrew_dalet = HEBREW LETTER DALET
	0x0ce4: "\u05d4", // hebrew_he = HEBREW LETTER HE
	0x0ce5: "\u05d5", // hebrew_waw = HEBREW LETTER VAV
	0x0ce6: "\u05d6", // hebrew_zain = HEBREW LETTER ZAYIN
	0x0ce7: "\u05d7", // hebrew_chet = HEBREW LETTER HET
	0x0ce8: "\u05d8", // hebrew_tet = HEBREW LETTER TET
	0x0ce9: "\u05d9", // hebrew_yod = HEBREW LETTER YOD
	0x0cea: "\u05da", // hebrew_finalkaph = HEBREW LETTER FINAL KAF
	0x0ceb: "\u05db", // hebrew_kaph = HEBREW LETTER KAF
	0x0cec: "\u05dc", // hebrew_lamed = HEBREW LETTER LAMED
	0x0ced: "\u05dd", // hebrew_finalmem = HEBREW LETTER FINAL MEM
	0x0cee: "\u05de", // hebrew_mem = HEBREW LETTER MEM
	0x0cef: "\u05df", // hebrew_finalnun = HEBREW LETTER FINAL NUN
	0x0cf0: "\u05e0", // hebrew_nun = HEBREW LETTER NUN
	0x0cf1: "\u05e1", // hebrew_samech = HEBREW LETTER SAMEKH
	0x0cf2: "\u05e2", // hebrew_ayin = HEBREW LETTER AYIN
	0x0cf3: "\u05e3", // hebrew_finalpe = HEBREW LETTER FINAL PE
	0x0cf4: "\u05e4", // hebrew_pe = HEBREW LETTER PE
	0x0cf5: "\u05e5", // hebrew_finalzade = HEBREW LETTER FINAL TSADI
	0x0cf6: "\u05e6", // hebrew_zade = HEBREW LETTER TSADI
	0x0cf7: "\u05e7", // hebrew_qoph = HEBREW LETTER QOF
	0x0cf8: "\u05e8", // hebrew_resh = HEBREW LETTER RESH
	0x0cf9: "\u05e9", // hebrew_shin = HEBREW LETTER SHIN
	0x0cfa: "\u05ea", // hebrew_taw = HEBREW LETTER TAV
	0x0da1: "\u0e01", // Thai_kokai = THAI CHARACTER KO KAI
	0x0da2: "\u0e02", // Thai_khokhai = THAI CHARACTER KHO KHAI
	0x0da3: "\u0e03", // Thai_khokhuat = THAI CHARACTER KHO KHUAT
	0x0da4: "\u0e04", // Thai_khokhwai = THAI CHARACTER KHO KHWAI
	0x0da5: "\u0e05", // Thai_khokhon = THAI CHARACTER KHO KHON
	0x0da6: "\u0e06", // Thai_khorakhang = THAI CHARACTER KHO RAKHANG
	0x0da7: "\u0e07", // Thai_ngongu = THAI CHARACTER NGO NGU
	0x0da8: "\u0e08", // Thai_chochan = THAI CHARACTER CHO CHAN
	0x0da9: "\u0e09", // Thai_choching = THAI CHARACTER CHO CHING
	0x0daa: "\u0e0a", // Thai_chochang = THAI CHARACTER CHO CHANG
	0x0dab: "\u0e0b", // Thai_soso = THAI CHARACTER SO SO
	0x0dac: "\u0e0c", // Thai_chochoe = THAI CHARACTER CHO CHOE
	0x0dad: "\u0e0d", // Thai_yoying = THAI CHARACTER YO YING
	0x0dae: "\u0e0e", // Thai_dochada = THAI CHARACTER DO CHADA
	0x0daf: "\u0e0f", // Thai_topatak = THAI CHARACTER TO PATAK
	0x0db0: "\u0e10", // Thai_thothan = THAI CHARACTER THO THAN
	0x0db1: "\u0e11", // Thai_thonangmontho = THAI CHARACTER THO NANGMONTHO
	0x0db2: "\u0e12", // Thai_thophuthao = THAI CHARACTER THO PHUTHAO
	0x0db3: "\u0e13", // Thai_nonen = THAI CHARACTER NO NEN
	0x0db4: "\u0e14", // Thai_dodek = THAI CHARACTER DO DEK
	0x0db5: "\u0e15", // Thai_totao = THAI CHARACTER TO TAO
	0x0db6: "\u0e16", // Thai_thothung = THAI CHARACTER THO THUNG
	0x0db7: "\u0e17", // Thai_thothahan = THAI CHARACTER THO THAHAN
	0x0db8: "\u0e18", // Thai_thothong = THAI CHARACTER THO THONG
	0x0db9: "\u0e19", // Thai_nonu = THAI CHARACTER NO NU
	0x0dba: "\u0e1a", // Thai_bobaimai = THAI CHARACTER BO BAIMAI
	0x0dbb: "\u0e1b", // Thai_popla = THAI CHARACTER PO PLA
	0x0dbc: "\u0e1c", // Thai_phophung = THAI CHARACTER PHO PHUNG
	0x0dbd: "\u0e1d", // Thai_fofa = THAI CHARACTER FO FA
	0x0dbe: "\u0e1e", // Thai_phophan = THAI CHARACTER PHO PHAN
	0x0dbf: "\u0e1f", // Thai_fofan = THAI CHARACTER FO FAN
	0x0dc0: "\u0e20", // Thai_phosamphao = THAI CHARACTER PHO SAMPHAO
	0x0dc1: "\u0e21", // Thai_moma = THAI CHARACTER MO MA
	0x0dc2: "\u0e22", // Thai_yoyak = THAI CHARACTER YO YAK
	0x0dc3: "\u0e23", // Thai_rorua = THAI CHARACTER RO RUA
	0x0dc4: "\u0e24", // Thai_ru = THAI CHARACTER RU
	0x0dc5: "\u0e25", // Thai_loling = THAI CHARACTER LO LING
	0x0dc6: "\u0e26", // Thai_lu = THAI CHARACTER LU
	0x0dc7: "\u0e27", // Thai_wowaen = THAI CHARACTER WO WAEN
	0x0dc8: "\u0e28", // Thai_sosala = THAI CHARACTER SO SALA
	0x0dc9: "\u0e29", // Thai_sorusi = THAI CHARACTER SO RUSI
	0x0dca: "\u0e2a", // Thai_sosua = THAI CHARACTER SO SUA
	0x0dcb: "\u0e2b", // Thai_hohip = THAI CHARACTER HO HIP
	0x0dcc: "\u0e2c", // Thai_lochula = THAI CHARACTER LO CHULA
	0x0dcd: "\u0e2d", // Thai_oang = THAI CHARACTER O ANG
	0x0dce: "\u0e2e", // Thai_honokhuk = THAI CHARACTER HO NOKHUK
	0x0dcf: "\u0e2f", // Thai_paiyannoi = THAI CHARACTER PAIYANNOI
	0x0dd0: "\u0e30", // Thai_saraa = THAI CHARACTER SARA A
	0x0dd1: "\u0e31", // Thai_maihanakat = THAI CHARACTER MAI HAN-AKAT
	0x0dd2: "\u0e32", // Thai_saraaa = THAI CHARACTER SARA AA
	0x0dd3: "\u0e33", // Thai_saraam = THAI CHARACTER SARA AM
	0x0dd4: "\u0e34", // Thai_sarai = THAI CHARACTER SARA I
	0x0dd5: "\u0e35", // Thai_saraii = THAI CHARACTER SARA II
	0x0dd6: "\u0e36", // Thai_saraue = THAI CHARACTER SARA UE
	0x0dd7: "\u0e37", // Thai_sarauee = THAI CHARACTER SARA UEE
	0x0dd8: "\u0e38", // Thai_sarau = THAI CHARACTER SARA U
	0x0dd9: "\u0e39", // Thai_sarauu = THAI CHARACTER SARA UU
	0x0dda: "\u0e3a", // Thai_phinthu = THAI CHARACTER PHINTHU
	0x0ddf: "\u0e3f", // Thai_baht = THAI CURRENCY SYMBOL BAHT
	0x0de0: "\u0e40", // Thai_sarae = THAI CHARACTER SARA E
	0x0de1: "\u0e41", // Thai_saraae = THAI CHARACTER SARA AE
	0x0de2: "\u0e42", // Thai_sarao = THAI CHARACTER SARA O
	0x0de3: "\u0e43", // Thai_saraaimaimuan = THAI CHARACTER SARA AI MAIMUAN
	0x0de4: "\u0e44", // Thai_saraaimaimalai = THAI CHARACTER SARA AI MAIMALAI
	0x0de5: "\u0e45", // Thai_lakkhangyao = THAI CHARACTER LAKKHANGYAO
	0x0de6: "\u0e46", // Thai_maiyamok = THAI CHARACTER MAIYAMOK
	0x0de7: "\u0e47", // Thai_maitaikhu = THAI CHARACTER MAITAIKHU
	0x0de8: "\u0e48", // Thai_maiek = THAI CHARACTER MAI EK
	0x0de9: "\u0e49", // Thai_maitho = THAI CHARACTER MAI THO
	0x0dea: "\u0e4a", // Thai_maitri = THAI CHARACTER MAI TRI
	0x0deb: "\u0e4b", // Thai_maichattawa = THAI CHARACTER MAI CHATTAWA
	0x0dec: "\u0e4c", // Thai_thanthakhat = THAI CHARACTER THANTHAKHAT
	0x0ded: "\u0e4d", // Thai_nikhahit = THAI CHARACTER NIKHAHIT
	0x0df0: "\u0e50", // Thai_leksun = THAI DIGIT ZERO
	0x0df1: "\u0e51", // Thai_leknung = THAI DIGIT ONE
	0x0df2: "\u0e52", // Thai_leksong = THAI DIGIT TWO
	0x0df3: "\u0e53", // Thai_leksam = THAI DIGIT THREE
	0x0df4: "\u0e54", // Thai_leksi = THAI DIGIT FOUR
	0x0df5: "\u0e55", // Thai_lekha = THAI DIGIT FIVE
	0x0df6: "\u0e56", // Thai_lekhok = THAI DIGIT SIX
	0x0df7: "\u0e57", // Thai_lekchet = THAI DIGIT SEVEN
	0x0df8: "\u0e58", // Thai_lekpaet = THAI DIGIT EIGHT
	0x0df9: "\u0e59", // Thai_lekkao = THAI DIGIT NINE
	0x0ea1: "\u3131", // Hangul_Kiyeog = HANGUL LETTER KIYEOK
	0x0ea2: "\u3132", // Hangul_SsangKiyeog = HANGUL LETTER SSANGKIYEOK
	0x0ea3: "\u3133", // Hangul_KiyeogSios = HANGUL LETTER KIYEOK-SIOS
	0x0ea4: "\u3134", // Hangul_Nieun = HANGUL LETTER NIEUN
	0x0ea5: "\u3135", // Hangul_NieunJieuj = HANGUL LETTER NIEUN-CIEUC
	0x0ea6: "\u3136", // Hangul_NieunHieuh = HANGUL LETTER NIEUN-HIEUH
	0x0ea7: "\u3137", // Hangul_Dikeud = HANGUL LETTER TIKEUT
	0x0ea8: "\u3138", // Hangul_SsangDikeud = HANGUL LETTER SSANGTIKEUT
	0x0ea9: "\u3139", // Hangul_Rieul = HANGUL LETTER RIEUL
	0x0eaa: "\u313a", // Hangul_RieulKiyeog = HANGUL LETTER RIEUL-KIYEOK
	0x0eab: "\u313b", // Hangul_RieulMieum = HANGUL LETTER RIEUL-MIEUM
	0x0eac: "\u313c", // Hangul_RieulPieub = HANGUL LETTER RIEUL-PIEUP
	0x0ead: "\u313d", // Hangul_RieulSios = HANGUL LETTER RIEUL-SIOS
	0x0eae: "\u313e", // Hangul_RieulTieut = HANGUL LETTER RIEUL-THIEUTH
	0x0eaf: "\u313f", // Hangul_RieulPhieuf = HANGUL LETTER RIEUL-PHIEUPH
	0x0eb0: "\u3140", // Hangul_RieulHieuh = HANGUL LETTER RIEUL-HIEUH
	0x0eb1: "\u3141", // Hangul_Mieum = HANGUL LETTER MIEUM
	0x0eb2: "\u3142", // Hangul_Pieub = HANGUL LETTER PIEUP
	0x0eb3: "\u3143", // Hangul_SsangPieub = HANGUL LETTER SSANGPIEUP
	0x0eb4: "\u3144", // Hangul_PieubSios = HANGUL LETTER PIEUP-SIOS
	0x0eb5: "\u3145", // Hangul_Sios = HANGUL LETTER SIOS
	0x0eb6: "\u3146", // Hangul_SsangSios = HANGUL LETTER SSANGSIOS
	0x0eb7: "\u3147", // Hangul_Ieung = HANGUL LETTER IEUNG
	0x0eb8: "\u3148", // Hangul_Jieuj = HANGUL LETTER CIEUC
	0x0eb9: "\u3149", // Hangul_SsangJieuj = HANGUL LETTER SSANGCIEUC
	0x0eba: "\u314a", // Hangul_Cieuc = HANGUL LETTER CHIEUCH
	0x0ebb: "\u314b", // Hangul_Khieuq = HANGUL LETTER KHIEUKH
	0x0ebc: "\u314c", // Hangul_Tieut = HANGUL LETTER THIEUTH
	0x0ebd: "\u314d", // Hangul_Phieuf = HANGUL LETTER PHIEUPH
	0x0ebe: "\u314e", // Hangul_Hieuh = HANGUL LETTER HIEUH
	0x0ebf: "\u314f", // Hangul_A = HANGUL LETTER A
	0x0ec0: "\u3150", // Hangul_AE = HANGUL LETTER AE
	0x0ec1: "\u3151", // Hangul_YA = HANGUL LETTER YA
	0x0ec2: "\u3152", // Hangul_YAE = HANGUL LETTER YAE
	0x0ec3: "\u3153", // Hangul_EO = HANGUL LETTER EO
	0x0ec4: "\u3154", // Hangul_E = HANGUL LETTER E
	0x0ec5: "\u3155", // Hangul_YEO = HANGUL LETTER YEO
	0x0ec6: "\u3156", // Hangul_YE = HANGUL LETTER YE
	0x0ec7: "\u3157", // Hangul_O = HANGUL LETTER O
	0x0ec8: "\u3158", // Hangul_WA = HANGUL LETTER WA
	0x0ec9: "\u3159", // Hangul_WAE = HANGUL LETTER WAE
	0x0eca: "\u315a", // Hangul_OE = HANGUL LETTER OE
	0x0ecb: "\u315b", // Hangul_YO = HANGUL LETTER YO
	0x0ecc: "\u315c", // Hangul_U = HANGUL LETTER U
	0x0ecd: "\u315d", // Hangul_WEO = HANGUL LETTER WEO
	0x0ece: "\u315e", // Hangul_WE = HANGUL LETTER WE
	0x0ecf: "\u315f", // Hangul_WI = HANGUL LETTER WI
	0x0ed0: "\u3160", // Hangul_YU = HANGUL LETTER YU
	0x0ed1: "\u3161", // Hangul_EU = HANGUL LETTER EU
	0x0ed2: "\u3162", // Hangul_YI = HANGUL LETTER YI
	0x0ed3: "\u3163", // Hangul_I = HANGUL LETTER I
	0x0ed4: "\u11a8", // Hangul_J_Kiyeog = HANGUL JONGSEONG KIYEOK
	0x0ed5: "\u11a9", // Hangul_J_SsangKiyeog = HANGUL JONGSEONG SSANGKIYEOK
	0x0ed6: "\u11aa", // Hangul_J_KiyeogSios = HANGUL JONGSEONG KIYEOK-SIOS
	0x0ed7: "\u11ab", // Hangul_J_Nieun = HANGUL JONGSEONG NIEUN
	0x0ed8: "\u11ac", // Hangul_J_NieunJieuj = HANGUL JONGSEONG NIEUN-CIEUC
	0x0ed9: "\u11ad", // Hangul_J_NieunHieuh = HANGUL JONGSEONG NIEUN-HIEUH
	0x0eda: "\u11ae", // Hangul_J_Dikeud = HANGUL JONGSEONG TIKEUT
	0x0edb: "\u11af", // Hangul_J_Rieul = HANGUL JONGSEONG RIEUL
	0x0edc: "\u11b0", // Hangul_J_RieulKiyeog = HANGUL JONGSEONG RIEUL-KIYEOK
	0x0edd: "\u11b1", // Hangul_J_RieulMieum = HANGUL JONGSEONG RIEUL-MIEUM
	0x0ede: "\u11b2", // Hangul_J_RieulPieub = HANGUL JONGSEONG RIEUL-PIEUP
	0x0edf: "\u11b3", // Hangul_J_RieulSios = HANGUL JONGSEONG RIEUL-SIOS
	0x0ee0: "\u11b4", // Hangul_J_RieulTieut = HANGUL JONGSEONG RIEUL-THIEUTH
	0x0ee1: "\u11b5", // Hangul_J_RieulPhieuf = HANGUL JONGSEONG RIEUL-PHIEUPH
	0x0ee2: "\u11b6", // Hangul_J_RieulHieuh = HANGUL JONGSEONG RIEUL-HIEUH
	0x0ee3: "\u11b7", // Hangul_J_Mieum = HANGUL JONGSEONG MIEUM
	0x0ee4: "\u11b8", // Hangul_J_Pieub = HANGUL JONGSEONG PIEUP
	0x0ee5: "\u11b9", // Hangul_J_PieubSios = HANGUL JONGSEONG PIEUP-SIOS
	0x0ee6: "\u11ba", // Hangul_J_Sios = HANGUL JONGSEONG SIOS
	0x0ee7: "\u11bb", // Hangul_J_SsangSios = HANGUL JONGSEONG SSANGSIOS
	0x0ee8: "\u11bc", // Hangul_J_Ieung = HANGUL JONGSEONG IEUNG
	0x0ee9: "\u11bd", // Hangul_J_Jieuj = HANGUL JONGSEONG CIEUC
	0x0eea: "\u11be", // Hangul_J_Cieuc = HANGUL JONGSEONG CHIEUCH
	0x0eeb: "\u11bf", // Hangul_J_Khieuq = HANGUL JONGSEONG KHIEUKH
	0x0eec: "\u11c0", // Hangul_J_Tieut = HANGUL JONGSEONG THIEUTH
	0x0eed: "\u11c1", // Hangul_J_Phieuf = HANGUL JONGSEONG PHIEUPH
	0x0eee: "\u11c2", // Hangul_J_Hieuh = HANGUL JONGSEONG HIEUH
	0x0eef: "\u316d", // Hangul_RieulYeorinHieuh = HANGUL LETTER RIEUL-YEORINHIEUH
	0x0ef0: "\u3171", // Hangul_SunkyeongeumMieum = HANGUL LETTER KAPYEOUNMIEUM
	0x0ef1: "\u3178", // Hangul_SunkyeongeumPieub = HANGUL LETTER KAPYEOUNPIEUP
	0x0ef2: "\u317f", // Hangul_PanSios = HANGUL LETTER PANSIOS
	0x0ef3: "\u3181", // Hangul_KkogjiDalrinIeung = HANGUL LETTER YESIEUNG
	0x0ef4: "\u3184", // Hangul_SunkyeongeumPhieuf = HANGUL LETTER KAPYEOUNPHIEUPH
	0x0ef5: "\u3186", // Hangul_YeorinHieuh = HANGUL LETTER YEORINHIEUH
	0x0ef6: "\u318d", // Hangul_AraeA = HANGUL LETTER ARAEA
	0x0ef7: "\u318e", // Hangul_AraeAE = HANGUL LETTER ARAEAE
	0x0ef8: "\u11eb", // Hangul_J_PanSios = HANGUL JONGSEONG PANSIOS
	0x0ef9: "\u11f0", // Hangul_J_KkogjiDalrinIeung = HANGUL JONGSEONG YESIEUNG
	0x0efa: "\u11f9", // Hangul_J_YeorinHieuh = HANGUL JONGSEONG YEORINHIEUH
	0x0eff: "\u20a9", // Korean_Won = WON SIGN
	0x13a4: "\u20ac", // Euro = EURO SIGN
	0x13bc: "\u0152", // OE = LATIN CAPITAL LIGATURE OE
	0x13bd: "\u0153", // oe = LATIN SMALL LIGATURE OE
	0x13be: "\u0178", // Ydiaeresis = LATIN CAPITAL LETTER Y WITH DIAERESIS
	0x20ac: "\u20ac", // EuroSign = EURO SIGN */
}
