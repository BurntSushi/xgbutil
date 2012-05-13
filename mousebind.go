/*
   All types and functions related to the mousebind package.
   They live here because they rely on state in XUtil.
*/
package xgbutil

import "github.com/BurntSushi/xgb/xproto"

// AttackMouseBindCallback associates an (event, window, mods, button)
// with a callback.
func (xu *XUtil) AttachMouseBindCallback(evtype int, win xproto.Window,
	mods uint16, button xproto.Button, fun MouseBindCallback) {

	xu.MousebindsLck.Lock()
	defer xu.MousebindsLck.Unlock()

	// Create key
	key := MouseBindKey{evtype, win, mods, button}

	// Do we need to allocate?
	if _, ok := xu.Mousebinds[key]; !ok {
		xu.Mousebinds[key] = make([]MouseBindCallback, 0)
	}

	xu.Mousebinds[key] = append(xu.Mousebinds[key], fun)
	xu.Mousegrabs[key] += 1
}

// MouseBindKeys returns a copy of all the keys in the 'Mousebinds' map.
func (xu *XUtil) MouseBindKeys() []MouseBindKey {
	xu.MousebindsLck.RLock()
	defer xu.MousebindsLck.RUnlock()

	keys := make([]MouseBindKey, len(xu.Mousebinds))
	i := 0
	for key, _ := range xu.Mousebinds {
		keys[i] = key
		i++
	}
	return keys
}

// MouseBindCallbacks returns a slice of callbacks for a particular key.
func (xu *XUtil) MouseBindCallbacks(key MouseBindKey) []MouseBindCallback {
	xu.MousebindsLck.RLock()
	defer xu.MousebindsLck.RUnlock()

	cbs := make([]MouseBindCallback, len(xu.Mousebinds[key]))
	for i, cb := range xu.Mousebinds[key] {
		cbs[i] = cb
	}
	return cbs
}

// RunMouseBindCallbacks executes every callback corresponding to a
// particular event/window/mod/button tuple.
func (xu *XUtil) RunMouseBindCallbacks(event interface{}, evtype int,
	win xproto.Window, mods uint16, button xproto.Button) {

	key := MouseBindKey{evtype, win, mods, button}
	for _, cb := range xu.MouseBindCallbacks(key) {
		cb.Run(xu, event)
	}
}

// ConnectedMouseBind checks to see if there are any key binds for a particular
// event type already in play. This is to work around comparing function
// pointers (not allowed in Go), which would be used in 'Connected'.
func (xu *XUtil) ConnectedMouseBind(evtype int, win xproto.Window) bool {
	xu.MousebindsLck.RLock()
	defer xu.MousebindsLck.RUnlock()

	// Since we can't create a full key, loop through all mouse binds
	// and check if evtype and window match.
	for key, _ := range xu.Mousebinds {
		if key.Evtype == evtype && key.Win == win {
			return true
		}
	}

	return false
}

// DetachMouseBindWindow removes all callbacks associated with a particular
// window and event type (either ButtonPress or ButtonRelease)
// Also decrements the counter in the corresponding 'Mousegrabs' map
// appropriately.
func (xu *XUtil) DetachMouseBindWindow(evtype int, win xproto.Window) {
	xu.MousebindsLck.Lock()
	defer xu.MousebindsLck.Unlock()

	// Since we can't create a full key, loop through all mouse binds
	// and check if evtype and window match.
	for key, _ := range xu.Mousebinds {
		if key.Evtype == evtype && key.Win == win {
			xu.Mousegrabs[key] -= len(xu.Mousebinds[key])
			delete(xu.Mousebinds, key)
		}
	}
}

// MouseBindGrabs returns the number of grabs on a particular
// event/window/mods/button combination. Namely, this combination
// uniquely identifies a grab. If it's repeated, we get BadAccess.
func (xu *XUtil) MouseBindGrabs(evtype int, win xproto.Window, mods uint16,
	button xproto.Button) int {

	xu.MousebindsLck.RLock()
	defer xu.MousebindsLck.RUnlock()

	key := MouseBindKey{evtype, win, mods, button}
	return xu.Mousegrabs[key] // returns 0 if key does not exist
}

// MouseDragFun is the kind of function used on each dragging step
// and at the end of a drag.
type MouseDragFun func(xu *XUtil, rootX, rootY, eventX, eventY int)

// MouseDragBeginFun is the kind of function used to initialize a drag.
// The difference between this and MouseDragFun is that the begin function
// returns a bool (of whether or not to cancel the drag) and an X resource
// identifier corresponding to a cursor.
type MouseDragBeginFun func(xu *XUtil, rootX, rootY,
	eventX, eventY int) (bool, xproto.Cursor)

// MouseDrag true when a mouse drag is in progress.
func (xu *XUtil) MouseDrag() bool {
	return xu.InMouseDrag
}

// MouseDragSet sets whether a mouse drag is in progress.
func (xu *XUtil) MouseDragSet(dragging bool) {
	xu.InMouseDrag = dragging
}

// MouseDragStep returns the function currently associated with each
// step of a mouse drag.
func (xu *XUtil) MouseDragStep() MouseDragFun {
	return xu.MouseDragStepFun
}

// MouseDragStepSet sets the function associated with the step of a drag.
func (xu *XUtil) MouseDragStepSet(f MouseDragFun) {
	xu.MouseDragStepFun = f
}

// MouseDragEnd returns the function currently associated with the
// end of a mouse drag.
func (xu *XUtil) MouseDragEnd() MouseDragFun {
	return xu.MouseDragEndFun
}

// MouseDragEndSet sets the function associated with the end of a drag.
func (xu *XUtil) MouseDragEndSet(f MouseDragFun) {
	xu.MouseDragEndFun = f
}
