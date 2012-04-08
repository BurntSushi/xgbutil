package main

import (
	"fmt"
	// "os" 
)

import (
	"burntsushi.net/go/x-go-binding/xgb"

	"burntsushi.net/go/xgbutil"
	// "burntsushi.net/go/xgbutil/ewmh" 
	"burntsushi.net/go/xgbutil/keybind"
	"burntsushi.net/go/xgbutil/mousebind"
	"burntsushi.net/go/xgbutil/xevent"
	"burntsushi.net/go/xgbutil/xprop"
	"burntsushi.net/go/xgbutil/xwindow"
)

func MyCallback(X *xgbutil.XUtil, e xevent.PropertyNotifyEvent) {
	atomName, err := xprop.AtomName(X, e.Atom)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("property %s, state %v\n", atomName, e.State)
	}
}

func MyCallback2(X *xgbutil.XUtil, e xevent.MappingNotifyEvent) {
	fmt.Printf("MappingNotify | Request = %v, FirstKeycode = %v, Count = %v\n",
		e.Request, e.FirstKeycode, e.Count)
}

func BPressCB(X *xgbutil.XUtil, e xevent.ButtonPressEvent) {
	fmt.Printf("Button press callback!\n")
	// xevent.ReplayPointer(X) 
}

func BReleaseCB(X *xgbutil.XUtil, e xevent.ButtonReleaseEvent) {
	fmt.Printf("Button release callback!\n")
}

func KeyPressCallback(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
	fmt.Printf("Key press callback!\n")
}

func KeyReleaseCallback(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
	fmt.Printf("Key release callback!\n")
}

func main() {
	fmt.Printf("Starting...\n")
	X, _ := xgbutil.Dial("")

	// active, _ := ewmh.ActiveWindowGet(X) 

	xwindow.Listen(X, X.RootWin(), xgb.EventMaskPropertyChange)

	// cb := xevent.PropertyNotifyFun(MyCallback) 
	// cb.Connect(X, X.RootWin()) 

	keybind.Initialize(X)

	keycbPress := keybind.KeyPressFun(KeyPressCallback)
	keycbPress.Connect(X, X.RootWin(), "Shift-delete") // Mod4-j

	keybind.XModMap(X)

	keybind.KeyPressFun(
		func(X *xgbutil.XUtil, ev xevent.KeyPressEvent) {
			keybind.XModMap(X)
		}).Connect(X, X.RootWin(), "Mod4-z")

	mousebind.ButtonPressFun(BPressCB).Connect(X, X.RootWin(), "Mod4-8", false)
	mousebind.ButtonReleaseFun(BReleaseCB).Connect(X, X.RootWin(), "Mod4-9", false)

	// keycbRelease := keybind.KeyReleaseFun(KeyReleaseCallback) 
	// keycbRelease.Connect(X, X.RootWin(), "Mod4-j") // Mod4-j 

	// fmt.Println(keybind.ParseString(X, "F1")) 

	xevent.Main(X)
}
