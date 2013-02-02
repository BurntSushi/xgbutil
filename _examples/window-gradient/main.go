// Example window-gradient demonstrates how to create several windows and draw
// gradients as their background. Namely, it shows how to use the
// xgraphics.Image type as a canvas that can change size. This example also
// demonstrates how to compress ConfigureNotify events so that the gradient
// drawing does not lag behind the rate of incoming ConfigureNotify events.
package main

import (
	"image"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	// Create three gradient windows of varying size with random colors.
	// Waiting a little bit inbetween seems to increase the diversity of the
	// random colors.

	newGradientWindow(X, 200, 200, newRandomColor(), newRandomColor())

	time.Sleep(500 * time.Millisecond)
	newGradientWindow(X, 400, 400, newRandomColor(), newRandomColor())

	time.Sleep(500 * time.Millisecond)
	newGradientWindow(X, 600, 600, newRandomColor(), newRandomColor())

	xevent.Main(X)
}

// newGradientWindow creates a new X window, paints the initial gradient
// image, and listens for ConfigureNotify events. (A new gradient image must
// be painted in response to each ConfigureNotify event, since a
// ConfigureNotify event corresponds to a change in the window's geometry.)
func newGradientWindow(X *xgbutil.XUtil, width, height int,
	start, end color.RGBA) {

	// Generate a new window id.
	win, err := xwindow.Generate(X)
	if err != nil {
		log.Fatal(err)
	}

	// Create the window and die if it fails.
	err = win.CreateChecked(X.RootWin(), 0, 0, width, height, 0)
	if err != nil {
		log.Fatal(err)
	}

	// In order to get ConfigureNotify events, we must listen to the window
	// using the 'StructureNotify' mask.
	win.Listen(xproto.EventMaskStructureNotify)

	// Paint the initial gradient to the window and then map the window.
	paintGradient(X, win.Id, width, height, start, end)
	win.Map()

	xevent.ConfigureNotifyFun(
		func(X *xgbutil.XUtil, ev xevent.ConfigureNotifyEvent) {
			// If the width and height have not changed, skip this one.
			if int(ev.Width) == width && int(ev.Height) == height {
				return
			}

			// Compress ConfigureNotify events so that we don't lag when
			// drawing gradients in response.
			ev = compressConfigureNotify(X, ev)

			// Update the width and height and paint the gradient image.
			width, height = int(ev.Width), int(ev.Height)
			paintGradient(X, win.Id, width, height, start, end)
		}).Connect(X, win.Id)
}

// paintGradient creates a new xgraphics.Image value and draws a gradient
// starting at color 'start' and ending at color 'end'.
//
// Since xgraphics.Image values use pixmaps and pixmaps cannot be resized,
// a new pixmap must be allocated for each resize event.
func paintGradient(X *xgbutil.XUtil, wid xproto.Window, width, height int,
	start, end color.RGBA) {

	ximg := xgraphics.New(X, image.Rect(0, 0, width, height))

	// Now calculate the increment step between each RGB channel in
	// the start and end colors.
	rinc := (0xff * (int(end.R) - int(start.R))) / width
	ginc := (0xff * (int(end.G) - int(start.G))) / width
	binc := (0xff * (int(end.B) - int(start.B))) / width

	// Now apply the increment to each "column" in the image.
	// Using 'ForExp' allows us to bypass the creation of a color.BGRA value
	// for each pixel in the image.
	ximg.ForExp(func(x, y int) (uint8, uint8, uint8, uint8) {
		return uint8(int(start.B) + (binc*x)/0xff),
			uint8(int(start.G) + (ginc*x)/0xff),
			uint8(int(start.R) + (rinc*x)/0xff),
			0xff
	})

	// Set the surface to paint on for ximg.
	// (This simply sets the background pixmap of the window to the pixmap
	// used by ximg.)
	ximg.XSurfaceSet(wid)

	// XDraw will draw the contents of ximg to its corresponding pixmap.
	ximg.XDraw()

	// XPaint will "clear" the window provided so that it shows the updated
	// pixmap.
	ximg.XPaint(wid)

	// Since we will not reuse ximg, we must destroy its pixmap.
	ximg.Destroy()
}

// compressConfigureNotify "compresses" incoming ConfigureNotify events so that
// event processing never lags behind gradient drawing.
// This is necessary because drawing a gradient cannot keep up with the rate
// at which ConfigureNotify events are sent to us, thereby creating a "lag".
// Compression works by examining the "future" of the event queue, and skipping
// ahead to the most recent ConfigureNotify event and throwing away the rest.
//
// A more detailed treatment of event compression can be found in
// xgbutil/examples/compress-events.
func compressConfigureNotify(X *xgbutil.XUtil,
	ev xevent.ConfigureNotifyEvent) xevent.ConfigureNotifyEvent {

	// Catch up with all X events as much as we can.
	X.Sync()
	xevent.Read(X, false) // non-blocking

	laste := ev
	for i, ee := range xevent.Peek(X) {
		if ee.Err != nil {
			continue
		}
		if cn, ok := ee.Event.(xproto.ConfigureNotifyEvent); ok {
			// Only compress this ConfigureNotify if it matches the window
			// of the original event.
			if ev.Event == cn.Event && ev.Window == cn.Window {
				laste = xevent.ConfigureNotifyEvent{&cn}
				defer func(i int) { xevent.DequeueAt(X, i) }(i)
			}
		}
	}
	return laste
}

// newRandomColor creates a new RGBA color where each channel (except alpha)
// is randomly generated.
func newRandomColor() color.RGBA {
	return color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 0xff,
	}
}
