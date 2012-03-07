/*
    Package keybind provides an easy interface to bind and run callback
    functions for human readable keybindings.
*/
package keybind

import "code.google.com/p/jamslam-x-go-binding/xgb"
import "github.com/BurntSushi/xgbutil"
import "github.com/BurntSushi/xgbutil/xevent"

func Grab(xu *xgbutil.XUtil, win xgb.Id, mods uint16, key byte) {
    for _, m := range xgbutil.IgnoreMods {
        xu.Conn().GrabKey(true, win, mods | m, key,
                          xgb.GrabModeAsync, xgb.GrabModeAsync)
    }
}

type KeyPressFun xevent.KeyPressFun

func (callback KeyPressFun) Connect(xu *xgbutil.XUtil, win xgb.Id,
                                    mods uint16, keycode byte) {
    xu.AttachKeyBindCallback(xevent.KeyPress, win, mods, keycode, callback)
    Grab(xu, win, mods, keycode)
}

func (callback KeyPressFun) Run(xu *xgbutil.XUtil, event interface{}) {
    callback(xu, event.(xevent.KeyPressEvent))
}

type KeyReleaseFun xevent.KeyReleaseFun

func (callback KeyReleaseFun) Connect(xu *xgbutil.XUtil, win xgb.Id,
                                      mods uint16, keycode byte) {
    xu.AttachKeyBindCallback(xevent.KeyRelease, win, mods, keycode, callback)
    Grab(xu, win, mods, keycode)
}

func (callback KeyReleaseFun) Run(xu *xgbutil.XUtil, event interface{}) {
    callback(xu, event.(xevent.KeyReleaseEvent))
}

