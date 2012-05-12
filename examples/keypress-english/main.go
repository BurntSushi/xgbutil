// Example keypress-english shows how to convert the State (modifiers) and
// Detail (keycode) members of Key{Press,Release} events to an english
// string representation.
package main

import (
	"log"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
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
	// because we aren't trying to make a grab.
	xevent.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			log.Println("Key:", keybind.LookupString(X, e.State, e.Detail))
		}).Connect(X, win.Id)

	// Finally, start the main event loop. This will route any appropriate
	// KeyPressEvents to your callback function.
	log.Println("Program initialized. Start pressing keys!")
	xevent.Main(X)
}
