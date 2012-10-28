// Example window-names fetches a list of all top-level client windows managed
// by the currently running window manager, and prints the name and size
// of each window.
//
// This example demonstrates how to use some aspects of the ewmh and icccm
// packages. It also shows how to use the xwindow package to find the
// geometry of a client window. In particular, finding the geometry is
// intelligent, as it includes the geometry of the decorations if they exist.
package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/xwindow"
)

func main() {
	// Connect to the X server using the DISPLAY environment variable.
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	// Get a list of all client ids.
	clientids, err := ewmh.ClientListGet(X)
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through each client, find its name and find its size.
	for _, clientid := range clientids {
		name, err := ewmh.WmNameGet(X, clientid)

		// If there was a problem getting _NET_WM_NAME or if its empty,
		// try the old-school version.
		if err != nil || len(name) == 0 {
			name, err = icccm.WmNameGet(X, clientid)

			// If we still can't find anything, give up.
			if err != nil || len(name) == 0 {
				name = "N/A"
			}
		}

		// Now find the geometry, including decorations, of the client window.
		// Note that DecorGeometry actually traverses the window tree by
		// issuing QueryTree requests until a top-level window (i.e., its
		// parent is the root window) is found. The geometry of *that* window
		// is then returned.
		dgeom, err := xwindow.New(X, clientid).DecorGeometry()
		if err != nil {
			log.Printf("Could not get geometry for %s (0x%X) because: %s",
				name, clientid, err)
			continue
		}

		fmt.Printf("%s (0x%x)\n", name, clientid)
		fmt.Printf("\tGeometry: %s\n", dgeom)
	}
}
