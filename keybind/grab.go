package keybind

import "log"

import "code.google.com/p/jamslam-x-go-binding/xgb"

import (
    "github.com/BurntSushi/xgbutil"
    "github.com/BurntSushi/xgbutil/xevent"
)

// Grabber is the public interface that will make the appropriate connections
// to register a lasting keyboard grab for three functions: begin,
// step and end.
// This is analogous to mousebind/drag.go
func Grabber(xu *xgbutil.XUtil, win xgb.Id, keyStr string,
             begin xgbutil.KeyGrabberFun,
             step xgbutil.KeyGrabberFun,
             end xgbutil.KeyGrabberFun) {
    KeyPressFun(
        func(xu *xgbutil.XUtil, ev xevent.KeyPressEvent) {
            grabberBegin(xu, ev, win, begin, step, end)
    }).Connect(xu, win, keyStr)
    xevent.KeyPressFun(grabberStep).Connect(xu, xu.Dummy())
    xevent.KeyReleaseFun(grabberEnd).Connect(xu, xu.Dummy())
}

// grabberGrab is a shortcut for grabbing the keyboard for a grabber.
func grabberGrab(xu *xgbutil.XUtil, win xgb.Id) bool {
    status, err := GrabKeyboard(xu, xu.Dummy())
    if err != nil {
        log.Printf("Key grabbering was unsuccessful because: %v", err)
        return false
    }
    if !status {
        log.Println("Key grabbering was unsuccessful because " +
                    "we could not establish a keyboard grab.")
        return false
    }

    xu.KeyGrabberSet(true)
    return true
}

// grabberUngrab is a shortcut for ungrabbing the keyboard for a grabber.
func grabberUngrab(xu *xgbutil.XUtil) {
    UngrabKeyboard(xu)
    xu.KeyGrabberSet(false)
}

// grabberBegin executes the "begin" function registered for the current drag.
// It also initiates the grab.
func grabberBegin(xu *xgbutil.XUtil, ev xevent.KeyPressEvent, win xgb.Id,
                  begin xgbutil.KeyGrabberFun,
                  step xgbutil.KeyGrabberFun,
                  end xgbutil.KeyGrabberFun) {
    // don't start a drag if one is already in progress
    if xu.KeyGrabber() {
        return
    }

    // Run begin first. It may tell us to cancel the grab.
    grab := begin(xu)

    // if we couldn't establish a grab, quit
    // Or quit if 'begin' tells us to.
    if !grab || !grabberGrab(xu, win) {
        return
    }

    xu.KeyGrabberStepSet(step)
    xu.KeyGrabberEndSet(end)
}

// grabberStep executes the "step" function registered for the current grabber.
func grabberStep(xu *xgbutil.XUtil, ev xevent.KeyPressEvent) {
    // If for whatever reason we don't have any *piece* of a grab,
    // we've gotta back out.
    if !xu.KeyGrabber() || xu.KeyGrabberStep() == nil ||
       xu.KeyGrabberEnd() == nil {
        grabberUngrab(xu)
        xu.KeyGrabberStepSet(nil)
        xu.KeyGrabberEndSet(nil)
        return
    }

    // now actually run the step
    xu.KeyGrabberStep()(xu)
}

// grabberEnd executes the "end" function registered for the current grabber.
func grabberEnd(xu *xgbutil.XUtil, ev xevent.KeyReleaseEvent) {
    if xu.KeyGrabberEnd() != nil {
        xu.KeyGrabberEnd()(xu)
    }

    grabberUngrab(xu)
    xu.KeyGrabberStepSet(nil)
    xu.KeyGrabberEndSet(nil)
}
