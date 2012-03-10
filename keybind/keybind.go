/*
    Package keybind provides an easy interface to bind and run callback
    functions for human readable keybindings.
*/
package keybind

import "fmt"
import "log"
import "strings"

import "code.google.com/p/jamslam-x-go-binding/xgb"
import "github.com/BurntSushi/xgbutil"
import "github.com/BurntSushi/xgbutil/xevent"

var modifiers []uint16 = []uint16{ // order matters!
    xgb.ModMaskShift, xgb.ModMaskLock, xgb.ModMaskControl,
    xgb.ModMask1, xgb.ModMask2, xgb.ModMask3, xgb.ModMask4, xgb.ModMask5,
    xgb.ModMaskAny,
}

// Initialize attaches the appropriate callbacks to make key bindings easier.
// i.e., update state of the world on a MappingNotify.
func Initialize(xu *xgbutil.XUtil) {
    // Listen to mapping notify events
    xevent.MappingNotifyFun(updateMaps).Connect(xu, xgbutil.NoWindow)

    // Give us an initial mapping state...
    keyMap, modMap := mapsGet(xu)
    xu.KeyMapSet(keyMap)
    xu.ModMapSet(modMap)
}

// updateMaps runs in response to MappingNotify events.
// It is responsible for making sure our view of the world's keyboard
// and modifier maps is correct. (Pointer mappings should be handled in
// a similar callback in the mousebind package.)
func updateMaps(xu *xgbutil.XUtil, e xevent.MappingNotifyEvent) {
    keyMap, modMap := mapsGet(xu)

    // Hold up... If this is MappingKeyboard, then we may have some keycode
    // changes. This is GROSS. We basically need to go through each keycode
    // in the map, look up the keysym using the new map and use that keysym
    // to look up the keycode in our current map. If the current keycode
    // does not equal the old keycode, then we log the change in a map of
    // old keycode -> new keycode.
    // Once the map is constructed, we look through all of our keybindings
    // and updated appropriately. *puke*
    // I am only somewhat confident that this is correct.
    if e.Request == xgb.MappingKeyboard {
        changes := make(map[byte]byte, 0)
        xuKeyMap := &xgbutil.KeyboardMapping{keyMap}

        min, max := minMaxKeycodeGet(xu)

        // let's not do too much allocation in our loop, shall we?
        var newSym, oldSym xgb.Keysym
        var column, oldKc byte

        // wrap 'int(..)' around bytes min and max to avoid overflow. Hideous.
        for newKc := int(min); newKc <= int(max); newKc++ {
            for column = 0; column < keyMap.KeysymsPerKeycode; column++ {
                // use new key map
                newSym = keysymGetWithMap(xu, xuKeyMap, byte(newKc), column)

                // uses old key map
                oldKc = keycodeGet(xu, newSym)
                oldSym = keysymGet(xu, byte(newKc), column)

                // If the old and new keysyms are the same, ignore!
                // Also ignore if either keysym is VoidSymbol
                if oldSym == newSym || oldSym == 0 || newSym == 0 {
                    continue
                }

                // these should match if there are NO changes
                if oldKc != byte(newKc) {
                    changes[oldKc] = byte(newKc)
                }
            }
        }

        // Now use 'changes' to do some regrabbing
        // Loop through all key bindings and check if we have any affected
        // key codes. (Note that each key binding may be associated with
        // multiple callbacks.)
        // We must ungrab everything first, in case two keys are being swapped.
        for _, key := range xu.KeyBindKeys() {
            if _, ok := changes[key.Code]; ok {
                Ungrab(xu, key.Win, key.Mod, key.Code)
            }
        }
        // Okay, now grab.
        for _, key := range xu.KeyBindKeys() {
            if newKc, ok := changes[key.Code]; ok {
                Grab(xu, key.Win, key.Mod, newKc)
                xu.UpdateKeyBindKey(key, newKc)
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
    xu.KeyMapSet(keyMap)
    xu.ModMapSet(modMap)
}

// minMaxKeycodeGet a simple accessor to the X setup info to return the
// minimum and maximum keycodes. They are typically 8 and 255, respectively.
func minMaxKeycodeGet(xu *xgbutil.XUtil) (byte, byte) {
    return xu.Conn().Setup.MinKeycode, xu.Conn().Setup.MaxKeycode
}

// A convenience function to grab the KeyboardMapping and ModifierMapping
// from X. We need to do this on startup (see Initialize) and whenever we
// get a MappingNotify event.
func mapsGet(xu *xgbutil.XUtil) (*xgb.GetKeyboardMappingReply,
                                 *xgb.GetModifierMappingReply) {
    min, max := minMaxKeycodeGet(xu)
    newKeymap, keyErr := xu.Conn().GetKeyboardMapping(min, max - min + 1)
    newModmap, modErr := xu.Conn().GetModifierMapping()

    // If there are errors, we really need to panic. We just can't do
    // any key binding without a mapping from the server.
    if keyErr != nil {
        panic(fmt.Sprintf("COULD NOT GET KEYBOARD MAPPING: %v\n" +
                          "THIS IS AN UNRECOVERABLE ERROR.\n",
                          keyErr))
    }
    if modErr != nil {
        panic(fmt.Sprintf("COULD NOT GET MODIFIER MAPPING: %v\n" +
                          "THIS IS AN UNRECOVERABLE ERROR.\n",
                          keyErr))
    }

    return newKeymap, newModmap
}

// ParseString takes a string of the format '[Mod[-Mod[...]]-]-KEY',
// i.e., 'Mod4-j', and returns a modifiers/keycode combo.
// (Actually, the parser is slightly more forgiving than what this comment
//  leads you to believe.)
func ParseString(xu *xgbutil.XUtil, str string) (uint16, byte) {
    mods, kc := uint16(0), byte(0)
    for _, part := range strings.Split(str, "-") {
        switch(strings.ToLower(part)) {
        case "shift":
            mods |= xgb.ModMaskShift
        case "lock":
            mods |= xgb.ModMaskLock
        case "control":
            mods |= xgb.ModMaskControl
        case "mod1":
            mods |= xgb.ModMask1
        case "mod2":
            mods |= xgb.ModMask2
        case "mod3":
            mods |= xgb.ModMask3
        case "mod4":
            mods |= xgb.ModMask4
        case "mod5":
            mods |= xgb.ModMask5
        case "any":
            mods |= xgb.ModMaskAny
        default: // a key code!
            if kc == 0 { // only accept the first keycode we see
                kc = lookupString(xu, part)
            }
        }
    }

    if kc == 0 {
        log.Printf("We could not find a valid keycode in the string '%s'. " +
                   "Things probably will not work right.\n", str)
    }

    return mods, kc
}

// lookupString is a wrapper around keycodeGet meant to make our search
// a bit more flexible if needed. (i.e., case-insensitive)
func lookupString(xu *xgbutil.XUtil, str string) byte {
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
        return byte(0)
    }

    return keycodeGet(xu, sym)
}

// Given a keysym, find the keycode mapped to it in the current X environment.
// keybind.Initialize MUST have been called before using this function.
func keycodeGet(xu *xgbutil.XUtil, keysym xgb.Keysym) byte {
    min, max := minMaxKeycodeGet(xu)
    keyMap := xu.KeyMapGet()

    var c byte
    for kc := int(min); kc <= int(max); kc++ {
        for c = 0; c < keyMap.KeysymsPerKeycode; c++ {
            if keysym == keysymGet(xu, byte(kc), c) {
                return byte(kc)
            }
        }
    }
    return 0
}

// keysymGet is a shortcut alias for 'keysymGetWithMap' using the current
// keymap stored in XUtil.
// keybind.Initialize MUST have been called before using this function.
func keysymGet(xu *xgbutil.XUtil, keycode byte, column byte) xgb.Keysym {
    return keysymGetWithMap(xu, xu.KeyMapGet(), keycode, column)
}

// keysymGetWithMap uses the given key map and finds a keysym associated
// with the given keycode in the current X environment.
func keysymGetWithMap(xu *xgbutil.XUtil, keyMap *xgbutil.KeyboardMapping,
                      keycode byte, column byte) xgb.Keysym {
    min, _ := minMaxKeycodeGet(xu)
    i := (int(keycode) - int(min)) * int(keyMap.KeysymsPerKeycode) + int(column)

    return keyMap.Keysyms[i]
}

// modGet finds the modifier currently associated with a given keycode.
// If a modifier doesn't exist for this keycode, then 0 is returned.
func modGet(xu *xgbutil.XUtil, keycode byte) uint16 {
    modMap := xu.ModMapGet()

    var i byte
    for i = 0; int(i) < len(modMap.Keycodes); i++ {
        if modMap.Keycodes[i] == keycode {
            return modifiers[i / modMap.KeycodesPerModifier]
        }
    }

    return 0
}

// XModMap should replicate the output of 'xmodmap'.
// This is mainly a sanity check, and may serve as an example of how to
// use modifier mappings.
func XModMap(xu *xgbutil.XUtil) {
    fmt.Println("Replicating `xmodmap`...")
    modMap := xu.ModMapGet()
    kPerMod := int(modMap.KeycodesPerModifier)

    // some nice names for the modifiers like xmodmap
    nice := []string{
        "shift", "lock", "control", "mod1", "mod2", "mod3", "mod4", "mod5",
    }

    var row int
    var comma string
    for mmi, _ := range modifiers[:len(modifiers) - 1] { // skip 'ModMaskAny'
        row = mmi * kPerMod
        comma = ""

        fmt.Printf("%s\t\t", nice[mmi])
        for _, kc := range modMap.Keycodes[row:row + kPerMod] {
            if kc != 0 {
                // This trickery is where things get really complicated.
                // We throw our hands up in the air if the first two columns
                // in our key map give us nothing.
                // But how do we know which column is the right one? I'm not
                // sure. This is what makes going from key code -> english
                // so difficult. The answer is probably buried somewhere
                // in the implementation of XLookupString in xlib. *shiver*
                ksym := keysymGet(xu, kc, 0)
                if ksym == 0 {
                    ksym = keysymGet(xu, kc, 1)
                }
                fmt.Printf("%s %s (0x%X)", comma, strKeysyms[ksym], kc)
                comma = ","
            }
        }
        fmt.Printf("\n")
    }
}

// Grabs a key with mods on a particular window.
// Will also grab all combinations of modifiers found in xgbutil.IgnoreMods
func Grab(xu *xgbutil.XUtil, win xgb.Id, mods uint16, key byte) {
    for _, m := range xgbutil.IgnoreMods {
        xu.Conn().GrabKey(true, win, mods | m, key,
                          xgb.GrabModeAsync, xgb.GrabModeAsync)
    }
}

// Ungrab undoes Grab. It will handle all combinations od modifiers found
// in xgbutil.IgnoreMods.
func Ungrab(xu *xgbutil.XUtil, win xgb.Id, mods uint16, key byte) {
    for _, m := range xgbutil.IgnoreMods {
        xu.Conn().UngrabKey(key, win, mods | m)
    }
}

// GrabKeyboard grabs the entire keyboard.
// Returns whether GrabStatus is successful and an error if one is reported by 
// XGB. It is possible to not get an error and the grab to be unsuccessful.
// The purpose of 'win' is that after a grab is successful, ALL Key*Events will
// be sent to that window. Make sure you have a callback attached :-)
func GrabKeyboard(xu *xgbutil.XUtil, win xgb.Id) (bool, error) {
    reply, err := xu.Conn().GrabKeyboard(false, win, 0,
                                         xgb.GrabModeAsync, xgb.GrabModeAsync)
    if err != nil {
        return false, xgbutil.Xerr(err, "GrabKeyboard",
                                   "Error grabbing keyboard on window '%x'",
                                   win)
    }

    return reply.Status == xgb.GrabStatusSuccess, nil
}

// UngrabKeyboard undoes GrabKeyboard.
func UngrabKeyboard(xu *xgbutil.XUtil) {
    xu.Conn().UngrabKeyboard(0)
}

