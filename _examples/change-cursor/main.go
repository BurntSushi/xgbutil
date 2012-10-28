// Example change-cursor shows how to use the cursor package to change the
// X cursor in a particular window. To see the new cursor, move your cursor
// into the window created by this program.
// Note that this only shows how to use one of the pre-defined cursors built
// into X using the "cursor" font. Creating your own cursor with your own
// image is a bit more complex, and probably not an instructive example.
//
// While this example shows how to set a cursor in an entire window, the cursor
// value returned from xcursor.CreateCursor[Extra] can be used in pointer
// grab requests too. (So that the cursor changes during the grab.)
package main

import (
	"log"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xcursor"
	"github.com/BurntSushi/xgbutil/xwindow"
)

func main() {
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	// Create the cursor. You can find a list of available cursors in
	// xcursor/cursordef.go.
	// We'll make an umbrella here, with an orange foreground and a blue
	// background. (The background it typically the outline of the cursor.)
	// Note that each component of the RGB color is a 16 bit color. I think
	// using the most significant byte to specify each component is good
	// enough.
	cursor, err := xcursor.CreateCursorExtra(X, xcursor.Umbrella,
		0xff00, 0x5500, 0x0000,
		0x3300, 0x6600, 0xff00)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new window. In the create window request, we'll set the
	// background color and set the cursor we created above.
	// This results in changing the cursor only when it moves into this window.
	win, err := xwindow.Generate(X)
	if err != nil {
		log.Fatal(err)
	}
	win.Create(X.RootWin(), 0, 0, 500, 500,
		xproto.CwBackPixel|xproto.CwCursor,
		0xffffffff, uint32(cursor))
	win.Map()

	// We can free the cursor now that we've set it.
	// If you plan on using this cursor again, then it shouldn't be freed.
	// (i.e., if you try to free this before setting it as the cursor in a
	// window, you'll get a BadCursor error when trying to use it.)
	xproto.FreeCursor(X.Conn(), cursor)

	// Block. No need to process any events.
	select {}
}
