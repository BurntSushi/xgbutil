// Example fullscreen shows how to make a window showing the Go Gopher go
// fullscreen and back using keybindings and EWMH.
package main

import (
	"bytes"
	"image"
	_ "image/png"
	"log"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/gopher"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
)

func main() {
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	keybind.Initialize(X) // call once before using keybind package

	// Read an example gopher image into a regular png image.
	img, _, err := image.Decode(bytes.NewBuffer(gopher.GopherPng()))
	if err != nil {
		log.Fatal(err)
	}

	// Now convert it into an X image.
	ximg := xgraphics.NewConvert(X, img)

	// Now show it in a new window.
	// We set the window title and tell the program to quit gracefully when
	// the window is closed.
	// There is also a convenience method, XShow, that requires no parameters.
	win := showImage(ximg, "The Go Gopher!", true)

	// Listen for key press events.
	win.Listen(xproto.EventMaskKeyPress)

	err = keybind.KeyPressFun(
		func(X *xgbutil.XUtil, ev xevent.KeyPressEvent) {
			println("fullscreen!")
			err := ewmh.WmStateReq(X, win.Id, ewmh.StateToggle,
				"_NET_WM_STATE_FULLSCREEN")
			if err != nil {
				log.Fatal(err)
			}
		}).Connect(X, win.Id, "f", false)
	if err != nil {
		log.Fatal(err)
	}

	err = keybind.KeyPressFun(
		func(X *xgbutil.XUtil, ev xevent.KeyPressEvent) {
			println("quit fullscreen!")
			err := ewmh.WmStateReq(X, win.Id, ewmh.StateToggle,
				"_NET_WM_STATE_FULLSCREEN")
			if err != nil {
				log.Fatal(err)
			}
		}).Connect(X, win.Id, "Escape", false)
	if err != nil {
		log.Fatal(err)
	}

	// If we don't block, the program will end and the window will disappear.
	// We could use a 'select{}' here, but xevent.Main will emit errors if
	// something went wrong, so use that instead.
	xevent.Main(X)
}

// This is a slightly modified version of xgraphics.XShowExtra that does
// not set any resize constraints on the window (so that it can go
// fullscreen).
func showImage(im *xgraphics.Image, name string, quit bool) *xwindow.Window {
	if len(name) == 0 {
		name = "xgbutil Image Window"
	}
	w, h := im.Rect.Dx(), im.Rect.Dy()

	win, err := xwindow.Generate(im.X)
	if err != nil {
		xgbutil.Logger.Printf("Could not generate new window id: %s", err)
		return nil
	}

	// Create a very simple window with dimensions equal to the image.
	win.Create(im.X.RootWin(), 0, 0, w, h, 0)

	// Make this window close gracefully.
	win.WMGracefulClose(func(w *xwindow.Window) {
		xevent.Detach(w.X, w.Id)
		keybind.Detach(w.X, w.Id)
		mousebind.Detach(w.X, w.Id)
		w.Destroy()

		if quit {
			xevent.Quit(w.X)
		}
	})

	// Set WM_STATE so it is interpreted as a top-level window.
	err = icccm.WmStateSet(im.X, win.Id, &icccm.WmState{
		State: icccm.StateNormal,
	})
	if err != nil { // not a fatal error
		xgbutil.Logger.Printf("Could not set WM_STATE: %s", err)
	}

	// Set _NET_WM_NAME so it looks nice.
	err = ewmh.WmNameSet(im.X, win.Id, name)
	if err != nil { // not a fatal error
		xgbutil.Logger.Printf("Could not set _NET_WM_NAME: %s", err)
	}

	// Paint our image before mapping.
	im.XSurfaceSet(win.Id)
	im.XDraw()
	im.XPaint(win.Id)

	// Now we can map, since we've set all our properties.
	// (The initial map is when the window manager starts managing.)
	win.Map()

	return win
}
