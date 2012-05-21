package xwindow

import (
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/xevent"
)

// WMGracefulClose will do all the necessary setup to implement the
// WM_DELETE_WINDOW protocol. This will prevent well-behaving window managers
// from killing your client whenever one of your windows is closed. (Killing
// a client is bad because it will destroy your X connection and any other
// clients you have open.)
// You must provide a callback function that is called when the window manager
// asks you to close your window. (You may provide some means of confirmation
// to the user, i.e., "Do you really want to quit?", but you should probably
// just wrap things up and call DestroyWindow.)
func (w *Window) WMGracefulClose(cb func(w *Window)) {
	// Get the current protocols so we don't overwrite anything.
	prots, _ := icccm.WmProtocolsGet(w.X, w.Id)

	// If WM_DELETE_WINDOW isn't here, add it. Otherwise, move on.
	wmdelete := false
	for _, prot := range prots {
		if prot == "WM_DELETE_WINDOW" {
			wmdelete = true
			break
		}
	}
	if !wmdelete {
		icccm.WmProtocolsSet(w.X, w.Id, append(prots, "WM_DELETE_WINDOW"))
	}

	// Attach a ClientMessage event handler. It will determine whether the
	// ClientMessage is a 'close' request, and if so, run the callback 'cb'
	// provided.
	xevent.ClientMessageFun(
		func(X *xgbutil.XUtil, ev xevent.ClientMessageEvent) {
			if icccm.IsDeleteRequest(X, ev) {
				cb(w)
			}
		}).Connect(w.X, w.Id)
}
