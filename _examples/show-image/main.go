// Example show-image is a very simple example to show how to use the
// xgraphics package to show an image in a window.
package main

import (
	"bytes"
	"image"
	_ "image/png"
	"log"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/gopher"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
)

func main() {
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

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
	ximg.XShowExtra("The Go Gopher!", true)

	// If we don't block, the program will end and the window will disappear.
	// We could use a 'select{}' here, but xevent.Main will emit errors if
	// something went wrong, so use that instead.
	xevent.Main(X)
}
