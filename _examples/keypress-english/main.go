// Example keypress-english shows how to convert the State (modifiers) and
// Detail (keycode) members of Key{Press,Release} events to an english
// string representation.
package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
)

var flagRoot = false

func init() {
	log.SetFlags(0)
	flag.BoolVar(&flagRoot, "root", flagRoot,
		"When set, the keyboard will be grabbed on the root window. "+
			"Make sure you have a way to kill the window created with "+
			"the mouse.")
	flag.Parse()
}

func main() {

	// Connect to the X server using the DISPLAY environment variable.
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	// Anytime the keybind (mousebind) package is used, keybind.Initialize
	// *should* be called once. It isn't strictly necessary, but allows your
	// keybindings to persist even if the keyboard mapping is changed during
	// run-time. (Assuming you're using the xevent package's event loop.)
	// It also handles the case when your modifier map is changed.
	keybind.Initialize(X)

	// Create a new window. We will listen for key presses and translate them
	// only when this window is in focus. (Similar to how `xev` works.)
	win, err := xwindow.Generate(X)
	if err != nil {
		log.Fatalf("Could not generate a new window X id: %s", err)
	}
	win.Create(X.RootWin(), 0, 0, 500, 500, xproto.CwBackPixel, 0xffffffff)

	// Listen for Key{Press,Release} events.
	win.Listen(xproto.EventMaskKeyPress, xproto.EventMaskKeyRelease)

	// Map the window.
	win.Map()

	// Notice that we use xevent.KeyPressFun instead of keybind.KeyPressFun,
	// because we aren't trying to make a grab *and* because we want to listen
	// to *all* key press events, rather than just a particular key sequence
	// that has been pressed.
	wid := win.Id
	if flagRoot {
		wid = X.RootWin()
	}
	xevent.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			// keybind.LookupString does the magic of implementing parts of
			// the X Keyboard Encoding to determine an english representation
			// of the modifiers/keycode tuple.
			// N.B. It's working for me, but probably isn't 100% correct in
			// all environments yet.
			modStr := keybind.ModifierString(e.State)
			keyStr := keybind.LookupString(X, e.State, e.Detail)
			if len(modStr) > 0 {
				log.Printf("Key: %s-%s\n", modStr, keyStr)
			} else {
				log.Println("Key:", keyStr)
			}

			if keybind.KeyMatch(X, "Escape", e.State, e.Detail) {
				if e.State&xproto.ModMaskControl > 0 {
					log.Println("Control-Escape detected. Quitting...")
					xevent.Quit(X)
				}
			}
		}).Connect(X, wid)

	// If we want root, then we take over the entire keyboard.
	if flagRoot {
		if err := keybind.GrabKeyboard(X, X.RootWin()); err != nil {
			log.Fatalf("Could not grab keyboard: %s", err)
		}
		log.Println("WARNING: We are taking *complete* control of the root " +
			"window. The only way out is to press 'Control + Escape' or to " +
			"close the window with the mouse.")
	}

	// Finally, start the main event loop. This will route any appropriate
	// KeyPressEvents to your callback function.
	log.Println("Program initialized. Start pressing keys!")
	xevent.Main(X)
}
