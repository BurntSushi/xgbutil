/*
    All types and functions related to the mousebind package.
    They live here because they rely on state in XUtil.
*/
package xgbutil

import "burntsushi.net/go/x-go-binding/xgb"

// MouseBindCallback operates in the spirit of Callback, except that it works
// specifically on mouse bindings.
type MouseBindCallback interface {
    Connect(xu *XUtil, win xgb.Id, buttonStr string, propagate bool, grab bool)
    Run(xu *XUtil, ev interface{})
}

// MouseBindKey is the type of the key in the map of mouse bindings.
// It essentially represents the tuple
// (event type, window id, modifier, button).
type MouseBindKey struct {
    Evtype int
    Win xgb.Id
    Mod uint16
    Button byte
}

// AttackMouseBindCallback associates an (event, window, mods, button)
// with a callback.
func (xu *XUtil) AttachMouseBindCallback(evtype int, win xgb.Id,
                                         mods uint16, button byte,
                                         fun MouseBindCallback) {
    // Create key
    key := MouseBindKey{evtype, win, mods, button}

    // Do we need to allocate?
    if _, ok := xu.mousebinds[key]; !ok {
        xu.mousebinds[key] = make([]MouseBindCallback, 0)
    }

    xu.mousebinds[key] = append(xu.mousebinds[key], fun)
    xu.mousegrabs[key] += 1
}

// MouseBindKeys returns a copy of all the keys in the 'mousebinds' map.
func (xu *XUtil) MouseBindKeys() []MouseBindKey {
    keys := make([]MouseBindKey, len(xu.mousebinds))
    i := 0
    for key, _ := range xu.mousebinds {
        keys[i] = key
        i++
    }
    return keys
}

// RunMouseBindCallbacks executes every callback corresponding to a
// particular event/window/mod/button tuple.
func (xu *XUtil) RunMouseBindCallbacks(event interface{}, evtype int,
                                       win xgb.Id, mods uint16, button byte) {
    // Create key
    key := MouseBindKey{evtype, win, mods, button}

    for _, cb := range xu.mousebinds[key] {
        cb.Run(xu, event)
    }
}

// ConnectedMouseBind checks to see if there are any key binds for a particular
// event type already in play. This is to work around comparing function
// pointers (not allowed in Go), which would be used in 'Connected'.
func (xu *XUtil) ConnectedMouseBind(evtype int, win xgb.Id) bool {
    // Since we can't create a full key, loop through all mouse binds
    // and check if evtype and window match.
    for key, _ := range xu.mousebinds {
        if key.Evtype == evtype && key.Win == win {
            return true
        }
    }

    return false
}

// DetachMouseBindWindow removes all callbacks associated with a particular
// window and event type (either ButtonPress or ButtonRelease)
// Also decrements the counter in the corresponding 'mousegrabs' map
// appropriately.
func (xu *XUtil) DetachMouseBindWindow(evtype int, win xgb.Id) {
    // Since we can't create a full key, loop through all mouse binds
    // and check if evtype and window match.
    for key, _ := range xu.mousebinds {
        if key.Evtype == evtype && key.Win == win {
            xu.mousegrabs[key] -= len(xu.mousebinds[key])
            delete(xu.mousebinds, key)
        }
    }
}

// MouseBindGrabs returns the number of grabs on a particular
// event/window/mods/button combination. Namely, this combination
// uniquely identifies a grab. If it's repeated, we get BadAccess.
func (xu *XUtil) MouseBindGrabs(evtype int, win xgb.Id, mods uint16,
                                button byte) int {
    key := MouseBindKey{evtype, win, mods, button}
    return xu.mousegrabs[key] // returns 0 if key does not exist
}

// MouseDrag true when a mouse drag is in progress.
func (xu *XUtil) MouseDrag() bool {
    return xu.mouseDrag
}

// MouseDragSet sets whether a mouse drag is in progress.
func (xu *XUtil) MouseDragSet(dragging bool) {
    xu.mouseDrag = dragging
}

// MouseDragStep returns the function currently associated with each
// step of a mouse drag.
func (xu *XUtil) MouseDragStep() MouseDragFun {
    return xu.mouseDragStep
}

// MouseDragStepSet sets the function associated with the step of a drag.
func (xu *XUtil) MouseDragStepSet(f MouseDragFun) {
    xu.mouseDragStep = f
}

// MouseDragEnd returns the function currently associated with the
// end of a mouse drag.
func (xu *XUtil) MouseDragEnd() MouseDragFun {
    return xu.mouseDragEnd
}

// MouseDragEndSet sets the function associated with the end of a drag.
func (xu *XUtil) MouseDragEndSet(f MouseDragFun) {
    xu.mouseDragEnd = f
}

