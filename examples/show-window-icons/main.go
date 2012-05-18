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
		xicon, err := xgraphics.FindIcon(X, wid, 300, 300)
		if err != nil {
			log.Printf("Could not find icon for window %d.", wid)
			continue
		}

		xgraphics.BlendBgColor(xicon, color.RGBA{0xff, 0x0, 0xff, 0xff})
		xicon.XShow()
	}

	// All we really need to do is block, so a 'select{}' would be sufficient.
	// But running the event loop will emit errors if anything went wrong.
	xevent.Main(X)
}
