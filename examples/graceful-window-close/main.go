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
// This seems like a lot of code to accomplish a relatively simple task, but a
// lot of it is boilerplate that you might already have in your program. Other
// portions of this example exist only to make the example workable (like
// creating new windows when clicking on one). The real magic in this example
// is in 'isDeleteRequest', setting WM_PROTOCOLS, and attaching a
// ClientMessage event handler.
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
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
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
	win.Map()

	// Tell the window manager that we support the WM_DELETE_WINDOW protocol.
	// If this is not set, the window manager will think you don't support
	// WM_DELETE_WINDOW and will KILL your client. Then your connection
	// will be lost. Try commenting out this and closing one of the windows.
	// You'll see :-)
	icccm.WmProtocolsSet(X, win.Id, []string{"WM_DELETE_WINDOW"})

	// A ClientMessage listener. We don't need to specify any event mask
	// to get these events, since they are sent with an empty event mask
	// (as specified by ICCCM).
	xevent.ClientMessageFun(
		func(X *xgbutil.XUtil, ev xevent.ClientMessageEvent) {
			if isDeleteRequest(X, ev) {
				// Make sure we detach all event handlers too.
				xevent.Detach(X, win.Id)
				mousebind.Detach(X, win.Id)
				xproto.DestroyWindow(X.Conn(), win.Id)
				counter--

				if counter == 0 {
					os.Exit(0)
				}
			}
		}).Connect(X, win.Id)

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

// isDeleteRequest checks whether a ClientMessage event satisfies the
// WM_DELETE_WINDOW protocol. Namely, the format must be 32, the type must
// be the WM_PROTOCOLS atom, and the first data item must be the atom
// WM_DELETE_WINDOW.
// This and other ICCCM protocols are certainly candidates to be included
// in xgbutil. I just haven't gotten there yet.
func isDeleteRequest(X *xgbutil.XUtil, ev xevent.ClientMessageEvent) bool {
	// Make sure the Format is 32. (Meaning that each data item is
	// 32 bits.)
	if ev.Format != 32 {
		return false
	}

	// Check to make sure the Type atom is WM_PROTOCOLS.
	typeName, err := xprop.AtomName(X, ev.Type)
	if err != nil || typeName != "WM_PROTOCOLS" { // not what we want
		return false
	}

	// Check to make sure the first data item is WM_DELETE_WINDOW.
	protocolType, err := xprop.AtomName(X,
		xproto.Atom(ev.Data.Data32[0]))
	if err != nil || protocolType != "WM_DELETE_WINDOW" {
		return false
	}

	return true
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
