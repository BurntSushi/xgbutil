// Example graceful-window-close shows how to create windows that can be closed
// without killing your X connection (and thereby destroying any other windows
// you may have open). This is actually achieved through an ICCCM standard.
// We add the WM_DELETE_WINDOW to the WM_PROTOCOLS property on our window.
// This indicates to well-behaving window managers that a certain kind of
// client message should be sent to our client when the window should be closed.
// If we *don't* add the WM_DELETE_WINDOW to the WM_PROTOCOLS property, the
// window manager will typically call KillClient---which will destroy the
// window and your X connection.
//
// If you click inside one of the windows created, a new window will be
// automatically created.
//
// The program will stop when all windows have been closed.
//
// For more information on the convention used please see
// http://tronche.com/gui/x/icccm/sec-4.html#s-4.1.2.7 and
// http://tronche.com/gui/x/icccm/sec-4.html#s-4.2.8
package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
)

// When counter reaches 0, exit.
var counter int32

// Just iniaitlize the RNG seed for generating random background colors.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// newWindow creates a new window with a random background color. It sets the
// WM_PROTOCOLS property to contain the WM_DELETE_WINDOW atom. It also sets
// up a ClientMessage event handler so that we know when to destroy the window.
// We also set up a mouse binding so that clicking inside a window will
// create another one.
func newWindow(X *xgbutil.XUtil) {
	counter++
	win, err := xwindow.Generate(X)
	if err != nil {
		log.Fatal(err)
	}

	// Get a random background color, create the window (ask to receive button
	// release events while we're at it) and map the window.
	bgColor := rand.Intn(0xffffff + 1)
	win.Create(X.RootWin(), 0, 0, 200, 200,
		xproto.CwBackPixel|xproto.CwEventMask,
		uint32(bgColor), xproto.EventMaskButtonRelease)

	// WMGracefulClose does all of the work for us. It sets the appropriate
	// values for WM_PROTOCOLS, and listens for ClientMessages that implement
	// the WM_DELETE_WINDOW protocol. When one is found, the provided callback
	// is executed.
	win.WMGracefulClose(
		func(w *xwindow.Window) {
			// Detach all event handlers.
			// This should always be done when a window can no longer
			// receive events.
			xevent.Detach(w.X, w.Id)
			mousebind.Detach(w.X, w.Id)
			w.Destroy()

			// Exit if there are no more windows left.
			counter--
			if counter == 0 {
				os.Exit(0)
			}
		})

	// It's important that the map comes after setting WMGracefulClose, since
	// the WM isn't obliged to watch updates to the WM_PROTOCOLS property.
	win.Map()

	// A mouse binding so that a left click will spawn a new window.
	// Note that we don't issue a grab here. Typically, window managers will
	// grab a button press on the client window (which usually activates the
	// window), so that we'd end up competing with the window manager if we
	// tried to grab it.
	// Instead, we set a ButtonRelease mask when creating the window and attach
	// a mouse binding *without* a grab.
	err = mousebind.ButtonReleaseFun(
		func(X *xgbutil.XUtil, ev xevent.ButtonReleaseEvent) {
			newWindow(X)
		}).Connect(X, win.Id, "1", false, false)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	// Anytime the mousebind (keybind) package is used, mousebind.Initialize
	// *should* be called once. In the case of the mousebind package, this
	// isn't strictly necessary (currently!), but the 'Drag' features of
	// the mousebind package won't work without it.
	mousebind.Initialize(X)

	// Create two windows to prove we can close one while keeping the
	// other alive.
	newWindow(X)
	newWindow(X)

	xevent.Main(X)
}
