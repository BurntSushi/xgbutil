// Example simple-keybinding shows how to grab keys on the root window and
// respond to them via callback functions. It also shows how to remove such
// callbacks so that they no longer respond to the key events.
// Note that more documentation can be found in the keybind package.
package main

import (
	"log"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
)

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
	keybind.Initialize(X)

	// Before attaching callbacks, wrap them in a callback function type.
	// The keybind package exposes two such callback types: keybind.KeyPressFun
	// and keybind.KeyReleaseFun.
	cb1 := keybind.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			log.Println("Key press!")
		})

	// We can now attach the callback to a particular window and key
	// combination. This particular example grabs a key on the root window,
	// which makes it a global keybinding.
	// Also, "Mod4-j" typically corresponds to pressing down the "Super" or
	// "Windows" key on your keyboard, and then pressing the letter "j".
	// N.B. This approach works by issuing a passive grab on the window
	// specified. To respond to Key{Press,Release} events without a grab, use
	// the xevent.Key{Press,Release}Fun callback function types instead.
	err = cb1.Connect(X, X.RootWin(), "Mod4-j", true)

	// A keybinding can fail if the key string could not be parsed, or if you're
	// trying to bind a key that has already been grabbed by another client.
	if err != nil {
		log.Fatal(err)
	}

	// We can even attach multiple callbacks to the same key.
	err = keybind.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			log.Println("A second handler always happens after the first.")
		}).Connect(X, X.RootWin(), "Mod4-j", true)
	if err != nil {
		log.Fatal(err)
	}

	// Finally, if we want this client to stop responding to key events, we
	// can attach another handler that, when run, detaches all previous
	// handlers.
	// This time, we'll show an example of a KeyRelease binding.
	err = keybind.KeyReleaseFun(
		func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
			// Use keybind.Detach to detach the root window
			// from all KeyPress *and* KeyRelease handlers.
			keybind.Detach(X, X.RootWin())

			log.Printf("Detached all Key{Press,Release}Events from the "+
				"root window (%d).", X.RootWin())
		}).Connect(X, X.RootWin(), "Mod4-Shift-q", true)
	if err != nil {
		log.Fatal(err)
	}

	// Finally, start the main event loop. This will route any appropriate
	// KeyPressEvents to your callback function.
	log.Println("Program initialized. Start pressing keys!")
	xevent.Main(X)
}
