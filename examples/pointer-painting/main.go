// Example pointer-painting shows how to draw on a window, MS Paint style.
// This is an extremely involved example, but it showcases a lot of xgbutil
// and how pieces of it can be tied together.
//
// If you're just starting with xgbutil, I highly recommend checking out the
// other examples before attempting to digest this one.
package main

import (
	"bytes"
	"image"
	_ "image/png"
	"log"

	"code.google.com/p/graphics-go/graphics"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/gopher"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
)

var (
	// The color to use for the background.
	bg = xgraphics.BGRA{0x0, 0x0, 0x0, 0xff}

	// Different colors for drawing.
	// The keys represent the key sequences that must be pressed to
	// switch to the color value.
	pencils = map[string]xgraphics.BGRA{
		"1": xgraphics.BGRA{0xff, 0xff, 0xff, 0xff}, // white
		"2": xgraphics.BGRA{0xff, 0x0, 0x0, 0xff},   // blue
		"3": xgraphics.BGRA{0x0, 0xff, 0x0, 0xff},   // green
		"4": xgraphics.BGRA{0x0, 0x0, 0xff, 0xff},   // red
		"5": xgraphics.BGRA{0x0, 0x7f, 0xff, 0xff},  // orange
		"6": xgraphics.BGRA{0xaa, 0x0, 0xff, 0x55},  // transparent pink
	}

	// The current pencil color.
	pencil = xgraphics.BGRA{0xff, 0xff, 0xff, 0xff}

	// The size of the tip of the pencil, in pixels.
	pencilTip = 30

	// The width and height of the canvas.
	width, height = 1000, 1000

	// The key sequence to use to clear the canvas.
	clearKey = "c"

	// Easter egg! Right click to draw a gopher with the following dimensions.
	gopherWidth, gopherHeight = 100, 100
)

// drawPencil takes an (x, y) position (from a MotionNotify event) and draws
// a rectangle of size pencilTip on to canvas.
func drawPencil(canvas *xgraphics.Image, win *xwindow.Window, x, y int) {
	// Create a subimage at (x, y) with pencilTip width and height from canvas.
	// Creating subimages is very cheap---no pixels are copied.
	// Moreover, when subimages are drawn to the screen, only the pixels in
	// the sub-image are sent to X.
	tipRect := midRect(x, y, pencilTip, pencilTip, width, height)

	// If the rectangle contains no pixels, don't draw anything.
	if tipRect.Empty() {
		return
	}

	// Output a little message.
	log.Printf("Drawing pencil point at (%d, %d)", x, y)

	// Create the subimage of the canvas to draw to.
	tip := canvas.SubImage(tipRect)

	// Now color each pixel in tip with the pencil color.
	tip.For(func(x, y int) xgraphics.BGRA {
		return xgraphics.BlendBGRA(pencil, canvas.At(x, y).(xgraphics.BGRA))
	})

	// Now draw the changes to the pixmap.
	tip.XDraw()

	// And paint them to the window.
	tip.XPaint(win.Id)
}

// drawGopher draws the gopher image to the canvas.
func drawGopher(canvas *xgraphics.Image, gopher image.Image,
	win *xwindow.Window, x, y int) {

	// Find the rectangle of the canvas where we're going to draw the gopher.
	gopherRect := midRect(x, y, gopherWidth, gopherHeight, width, height)

	// If the rectangle contains no pixels, don't draw anything.
	if gopherRect.Empty() {
		return
	}

	// Output a little message.
	log.Printf("Drawing gopher at (%d, %d)", x, y)

	// Get a subimage of the gopher that's in sync with gopherRect.
	gopherPt := image.Pt(gopher.Bounds().Min.X, gopher.Bounds().Min.Y)
	if gopherRect.Min.X == 0 {
		gopherPt.X = gopherWidth - gopherRect.Dx()
	}
	if gopherRect.Min.Y == 0 {
		gopherPt.Y = gopherHeight - gopherRect.Dy()
	}

	// Create the canvas subimage.
	subCanvas := canvas.SubImage(gopherRect)

	// Blend the gopher image into the sub-canvas.
	// This does alpha blending.
	xgraphics.BlendImage(subCanvas, gopher, gopherPt)

	// Now draw the changes to the pixmap.
	subCanvas.XDraw()

	// And paint them to the window.
	subCanvas.XPaint(win.Id)
}

// clearCanvas erases all your pencil marks.
func clearCanvas(canvas *xgraphics.Image, win *xwindow.Window) {
	log.Println("Clearing canvas...")
	canvas.For(func(x, y int) xgraphics.BGRA {
		return bg
	})

	canvas.XDraw()
	canvas.XPaint(win.Id)
}

func main() {
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	// Whenever the mousebind package is used, you must call Initialize.
	// Similarly for the keybind package.
	keybind.Initialize(X)
	mousebind.Initialize(X)

	// Easter egg! Use a right click to draw a gopher.
	gopherPng, _, err := image.Decode(bytes.NewBuffer(gopher.GopherPng()))
	if err != nil {
		log.Fatal(err)
	}

	// Now scale it to a reasonable size.
	gopher := image.NewRGBA(image.Rect(0, 0, gopherWidth, gopherHeight))
	graphics.Scale(gopher, gopherPng)

	// Create a new xgraphics.Image. It automatically creates an X pixmap for
	// you, and handles drawing to windows in the XDraw, XPaint and
	// XSurfaceSet functions.
	// N.B. An error is possible since X pixmap allocation can fail.
	canvas, err := xgraphics.New(X, image.Rect(0, 0, width, height))

	// Color in the background color.
	canvas.For(func(x, y int) xgraphics.BGRA {
		return bg
	})

	// Use the convenience function XShow to create and map the canvas window.
	// XShow will also set the surface window of canvas for us.
	win := canvas.XShow()

	// Listen for pointer motion events and key press events.
	win.Listen(xproto.EventMaskButtonPress | xproto.EventMaskButtonRelease |
		xproto.EventMaskKeyPress)

	mousebind.Drag(X, win.Id, win.Id, "1", false,
		func(X *xgbutil.XUtil, rx, ry, ex, ey int) (bool, xproto.Cursor) {
			drawPencil(canvas, win, ex, ey)
			return true, 0
		},
		func(X *xgbutil.XUtil, rx, ry, ex, ey int) {
			drawPencil(canvas, win, ex, ey)
		},
		func(X *xgbutil.XUtil, rx, ry, ex, ey int) {})

	mousebind.Drag(X, win.Id, win.Id, "3", false,
		func(X *xgbutil.XUtil, rx, ry, ex, ey int) (bool, xproto.Cursor) {
			drawGopher(canvas, gopher, win, ex, ey)
			return true, 0
		},
		func(X *xgbutil.XUtil, rx, ry, ex, ey int) {
			drawGopher(canvas, gopher, win, ex, ey)
		},
		func(X *xgbutil.XUtil, rx, ry, ex, ey int) {})

	// Bind to the clear key specified, and just redraw the bg color.
	keybind.KeyPressFun(
		func(X *xgbutil.XUtil, ev xevent.KeyPressEvent) {
			clearCanvas(canvas, win)
		}).Connect(X, win.Id, clearKey, false)

	// Bind a callback to each key specified in the 'pencils' color map.
	// The response is to simply switch the pencil color.
	for key, clr := range pencils {
		c := clr
		keybind.KeyPressFun(
			func(X *xgbutil.XUtil, ev xevent.KeyPressEvent) {
				log.Printf("Changing pencil color to: %#v", c)
				pencil = c
			}).Connect(X, win.Id, key, false)
	}

	xevent.Main(X)
}

// midRect takes an (x, y) position where the pointer was clicked, along with
// the width and height of the thing being drawn and the width and height of
// the canvas, and returns a Rectangle
// whose midpoint (roughly) is (x, y) and whose width and height match the
// parameters when the rectangle doesn't extend past the border of the canvas.
// Make sure to check if the rectange is empty or not before using it!
func midRect(x, y, width, height, canWidth, canHeight int) image.Rectangle {
	return image.Rect(
		max(0, min(canWidth, x-width/2)),   // top left x
		max(0, min(canHeight, y-height/2)), // top left y
		max(0, min(canWidth, x+width/2)),   // bottom right x
		max(0, min(canHeight, y+height/2)), // bottom right y
	)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
