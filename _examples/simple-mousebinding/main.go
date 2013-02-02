// Example simple-mousebinding shows how to grab buttons on the root window and
// respond to them via callback functions. It also shows how to remove such
// callbacks so that they no longer respond to the button events.
// Note that more documentation can be found in the mousebind package.
package main

import (
	"log"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
)

func main() {
	// Connect to the X server using the DISPLAY environment variable.
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	// Anytime the mousebind (keybind) package is used, mousebind.Initialize
	// *should* be called once. In the case of the mousebind package, this
	// isn't strictly necessary, but the 'Drag' features of the mousebind
	// package won't work without it.
	mousebind.Initialize(X)

	// Before attaching callbacks, wrap them in a callback function type.
	// The mouse package exposes two such callback types:
	// mousebind.ButtonPressFun and mousebind.ButtonReleaseFun.
	cb1 := mousebind.ButtonPressFun(
		func(X *xgbutil.XUtil, e xevent.ButtonPressEvent) {
			log.Println("Button press!")
		})

	// We can now attach the callback to a particular window and button
	// combination. This particular example grabs a button on the root window,
	// which makes it a global mouse binding.
	// Also, "Mod4-1" typically corresponds to pressing down the "Super" or
	// "Windows" key on your keyboard, and then pressing the left mouse button.
	// The last two parameters are whether to make a synchronous grab and
	// whether to actually issue a grab, respectively.
	// (The parameters used here are the common case.)
	// See the documentation for the Connect method for more details.
	err = cb1.Connect(X, X.RootWin(), "Mod4-1", false, true)

	// A mouse binding can fail if the mouse string could not be parsed, or if
	// you're trying to bind a button that has already been grabbed by another
	// client.
	if err != nil {
		log.Fatal(err)
	}

	// We can even attach multiple callbacks to the same button.
	err = mousebind.ButtonPressFun(
		func(X *xgbutil.XUtil, e xevent.ButtonPressEvent) {
			log.Println("A second handler always happens after the first.")
		}).Connect(X, X.RootWin(), "Mod4-1", false, true)
	if err != nil {
		log.Fatal(err)
	}

	// Finally, if we want this client to stop responding to mouse events, we
	// can attach another handler that, when run, detaches all previous
	// handlers.
	// This time, we'll show an example of a ButtonRelease binding.
	err = mousebind.ButtonReleaseFun(
		func(X *xgbutil.XUtil, e xevent.ButtonReleaseEvent) {
			// Use mousebind.Detach to detach the root window
			// from all ButtonPress *and* ButtonRelease handlers.
			mousebind.Detach(X, X.RootWin())
			mousebind.Detach(X, X.RootWin())

			log.Printf("Detached all Button{Press,Release}Events from the "+
				"root window (%d).", X.RootWin())
		}).Connect(X, X.RootWin(), "Mod4-Shift-1", false, true)
	if err != nil {
		log.Fatal(err)
	}

	// Finally, start the main event loop. This will route any appropriate
	// ButtonPressEvents to your callback function.
	log.Println("Program initialized. Start pressing mouse buttons!")
	xevent.Main(X)
}
