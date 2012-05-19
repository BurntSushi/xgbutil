// Example show-window-icons shows how to get a list of all top-level client
// windows managed by the currently running window manager, and show the icon
// for each window. (Each icon is shown by opening its own window.)
package main

import (
	"image/color"
	"log"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
)

var (
	// The icon width and height to use.
	// _NET_WM_ICON will be searched for an icon closest to these values.
	// The icon closest in size to what's specified here will be used.
	// The resulting icon will be scaled to this size.
	// (Set both to 0 to avoid scaling and use the biggest possible icon.)
	iconWidth, iconHeight = 300, 300
)

func main() {
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	// Get the list of window ids managed by the window manager.
	clients, err := ewmh.ClientListGet(X)
	if err != nil {
		log.Fatal(err)
	}

	// For each client, try to find its icon. If we find one, blend it with
	// a nice background color and show it in its own window.
	// Otherwise, skip it.
	for _, wid := range clients {
		// FindIcon will find an icon closest to the size specified.
		// If one can't be found, the resulting image will be scaled
		// automatically.
		// To avoid scaling the icon, specify '0' for both the width and height.
		// In this case, the largest icon found will be returned.
		xicon, err := xgraphics.FindIcon(X, wid, iconWidth, iconHeight)
		if err != nil {
			log.Printf("Could not find icon for window %d.", wid)
			continue
		}

		// Get the name of this client. (It will be set as the icon window's
		// name.)
		name, err := ewmh.WmNameGet(X, wid)
		if err != nil { // not a fatal error
			log.Println(err)
			name = ""
		}

		// Blend a pink background color so its easy to see that alpha blending
		// works.
		xgraphics.BlendBgColor(xicon, color.RGBA{0xff, 0x0, 0xff, 0xff})
		xicon.XShowExtra(name, false)
	}

	// All we really need to do is block, so a 'select{}' would be sufficient.
	// But running the event loop will emit errors if anything went wrong.
	xevent.Main(X)
}
