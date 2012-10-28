// Example xmodmap shows how one might implement a rudimentary version
// of xmodmap's modifier listing.
// (xmodmap is a program that shows all modifier keys, and which
// keysyms activate each modifier. xmodmap can also modify the modifier
// mapping, which this doesn't do.)
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
)

func main() {
	// Connect to the X server using the DISPLAY environment variable.
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	// Nice names for the modifier keys (same ones used by xmodmap).
	// This slice corresponds to the keybind.Modifiers slice (except for
	// the last 'Any' modifier element).
	nice := []string{
		"shift", "lock", "control", "mod1", "mod2", "mod3", "mod4", "mod5",
	}

	// Whenever one uses the keybind package, Initialize should always be
	// called first. In this case, it initializes the key and modifier maps.
	keybind.Initialize(X)

	// Get the current modifier map.
	// The map is actually a table, where rows correspond to modifiers, and
	// columns correspond to the keys that activate that modifier.
	modMap := keybind.ModMapGet(X)

	// The number of keycodes allowed per modifier (i.e., the number of
	// columns in the modifier map).
	kPerMod := int(modMap.KeycodesPerModifier)

	// Get the number of allowable keysyms per keycode.
	// This is used to search for a valid keysym for a particular keycode.
	symsPerKc := int(keybind.KeyMapGet(X).KeysymsPerKeycode)

	// Imitate everything...
	fmt.Printf("xmodmap:  up to %d keys per modifier, "+
		"(keycodes in parentheses):\n\n", kPerMod)

	// Iterate through all keyboard modifiers defined in xgb/xproto
	// except the 'Any' modifier (which is last).
	for mmi := range keybind.Modifiers[:len(keybind.Modifiers)-1] {
		niceName := nice[mmi]
		keys := make([]string, 0, kPerMod)

		// row is the row for the 'mmi' modifier in the modifier mapping table.
		row := mmi * kPerMod

		// Iterate over each keycode in the modifier map for this modifier.
		for _, kc := range modMap.Keycodes[row : row+kPerMod] {
			// If this entry doesn't have a keycode (i.e., it's zero), we
			// have to skip it.
			if kc == 0 {
				continue
			}

			// Look for the first valid keysym in the keyboard map corresponding
			// to this keycode. If one can't be found, output "BadKey."
			var ksym xproto.Keysym = 0
			for column := 0; column < symsPerKc; column++ {
				ksym = keybind.KeysymGet(X, kc, byte(column))
				if ksym != 0 {
					break
				}
			}

			if ksym == 0 {
				keys = append(keys, fmt.Sprintf("BadKey (0x%x)", kc))
			} else {
				keys = append(keys,
					fmt.Sprintf("%s (0x%x)", keybind.KeysymToStr(ksym), kc))
			}
		}

		fmt.Printf("%-12s%s\n", niceName, strings.Join(keys, ",  "))
	}
	fmt.Println("")
}
