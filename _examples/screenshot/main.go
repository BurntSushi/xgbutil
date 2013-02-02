// Example screenshot shows how to take a screenshot of the current desktop
// and show it in a window. In a comment, it also shows how to save it as
// a png.
//
// It works by getting the image of the root window, which automatically
// includes all child windows.
package main

import (
	"log"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
)

func main() {
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	// Use the "NewDrawable" constructor to create an xgraphics.Image value
	// from a drawable. (Usually this is done with pixmaps, but drawables
	// can also be windows.)
	ximg, err := xgraphics.NewDrawable(X, xproto.Drawable(X.RootWin()))
	if err != nil {
		log.Fatal(err)
	}

	// Shows the screenshot in a window.
	ximg.XShowExtra("Screenshot", true)

	// If you'd like to save it as a png, use:
	// err = ximg.SavePng("screenshot.png")
	// if err != nil {
	// log.Fatal(err)
	// }

	xevent.Main(X)
}
