package mousebind

import "log"

import "code.google.com/p/jamslam-x-go-binding/xgb"

import (
    "github.com/BurntSushi/xgbutil"
    "github.com/BurntSushi/xgbutil/xevent"
)

// Drag is the public interface that will make the appropriate connections
// to register a drag event for three functions: the begin function, the 
// step function and the end function.
func Drag(xu *xgbutil.XUtil, win xgb.Id, buttonStr string,
          begin xgbutil.MouseDragFun,
          step xgbutil.MouseDragFun,
          end xgbutil.MouseDragFun) {
    ButtonPressFun(
        func(xu *xgbutil.XUtil, ev xevent.ButtonPressEvent) {
            dragBegin(xu, ev, begin, step, end)
    }).Connect(xu, win, buttonStr, false, true)
}

// dragGrab is a shortcut for grabbing the pointer for a drag.
func dragGrab(xu *xgbutil.XUtil) bool {
    status, err := GrabPointer(xu, xu.Dummy(), xu.RootWin(), 0)
    if err != nil {
        log.Println("Mouse dragging was unsuccessful because: %v", err)
        return false
    }
    if !status {
        log.Println("Mouse dragging was unsuccessful because " +
                    "we could not establish a pointer grab.")
        return false
    }

    xu.MouseDragSet(true)
    return true
}

// dragUngrab is a shortcut for ungrabbing the pointer for a drag.
func dragUngrab(xu *xgbutil.XUtil) {
    UngrabPointer(xu)
    xu.MouseDragSet(false)
}

// dragStart executes the "begin" function registered for the current drag.
// It also initiates the grab.
func dragBegin(xu *xgbutil.XUtil, ev xevent.ButtonPressEvent,
               begin xgbutil.MouseDragFun,
               step xgbutil.MouseDragFun,
               end xgbutil.MouseDragFun) {
    // don't start a drag if one is already in progress
    if xu.MouseDrag() {
        return
    }

    // if we couldn't establish a grab, quit
    if !dragGrab(xu) {
        return
    }

    // we're committed. set the drag state and start the 'begin' function
    xu.MouseDragStepSet(step)
    xu.MouseDragEndSet(end)
    begin(xu, ev.RootX, ev.RootY, ev.EventX, ev.EventY)
}

// dragStep executes the "step" function registered for the current drag.
func dragStep(xu *xgbutil.XUtil, ev xevent.MotionNotifyEvent) {
    // If for whatever reason we don't have any *piece* of a grab,
    // we've gotta back out.
    if !xu.MouseDrag() || xu.MouseDragStep() == nil ||
           xu.MouseDragEnd() == nil {
        dragUngrab(xu)
        xu.MouseDragStepSet(nil)
        xu.MouseDragEndSet(nil)
        return
    }

    // now actually run the step
    xu.MouseDragStep()(xu, ev.RootX, ev.RootY, ev.EventX, ev.EventY)
}

// dragEnd executes the "end" function registered for the current drag.
func dragEnd(xu *xgbutil.XUtil, ev xevent.ButtonReleaseEvent) {
    if xu.MouseDragEnd() != nil {
        xu.MouseDragEnd()(xu, ev.RootX, ev.RootY, ev.EventX, ev.EventY)
    }

    dragUngrab(xu)
    xu.MouseDragStepSet(nil)
    xu.MouseDragEndSet(nil)
}
