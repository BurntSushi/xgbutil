// Example workarea-struts shows how to find the all of the struts set by
// top-level clients and apply them to each active monitor to get the true
// workarea for each monitor.
//
// For example, consider a monitor with a 1920x1080 resolution and a panel
// on the bottom of the desktop 50 pixels tall. The actual workarea for the
// desktop then, is (0, 0) 1920x1030. (When there is more than one monitor,
// the x and y coordinates will be adjusted too.)
package main

import (
	"fmt"
	"log"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xinerama"
	"github.com/BurntSushi/xgbutil/xrect"
	"github.com/BurntSushi/xgbutil/xwindow"
)

func main() {
	// Connect to the X server using the DISPLAY environment variable.
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	// Wrap the root window in a nice Window type.
	root := xwindow.New(X, X.RootWin())

	// Get the geometry of the root window.
	rgeom, err := root.Geometry()
	if err != nil {
		log.Fatal(err)
	}

	// Get the rectangles for each of the active physical heads.
	// These are returned sorted in order from left to right and then top
	// to bottom.
	// But first check if Xinerama is enabled. If not, use root geometry.
	var heads xinerama.Heads
	if X.ExtInitialized("XINERAMA") {
		heads, err = xinerama.PhysicalHeads(X)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		heads = xinerama.Heads{rgeom}
	}

	// Fetch the list of top-level client window ids currently managed
	// by the running window manager.
	clients, err := ewmh.ClientListGet(X)
	if err != nil {
		log.Fatal(err)
	}

	// Output the head geometry before modifying them.
	fmt.Println("Workarea for each head:")
	for i, head := range heads {
		fmt.Printf("\tHead #%d: %s\n", i+1, head)
	}

	// For each client, check to see if it has struts, and if so, apply
	// them to our list of head rectangles.
	for _, clientid := range clients {
		strut, err := ewmh.WmStrutPartialGet(X, clientid)
		if err != nil { // no struts for this client
			continue
		}

		// Apply the struts to our heads.
		// This modifies 'heads' in place.
		xrect.ApplyStrut(heads, uint(rgeom.Width()), uint(rgeom.Height()),
			strut.Left, strut.Right, strut.Top, strut.Bottom,
			strut.LeftStartY, strut.LeftEndY,
			strut.RightStartY, strut.RightEndY,
			strut.TopStartX, strut.TopEndX,
			strut.BottomStartX, strut.BottomEndX)
	}

	// Now output the head geometry again after modification.
	fmt.Println("Workarea for each head after applying struts:")
	for i, head := range heads {
		fmt.Printf("\tHead #%d: %s\n", i+1, head)
	}
}
