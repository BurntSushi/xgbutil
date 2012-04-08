package main

import (
	"fmt"
	"log"
)

import (
	"burntsushi.net/go/xgbutil"
	"burntsushi.net/go/xgbutil/keybind"
	"burntsushi.net/go/xgbutil/xevent"
)

func main() {
	X, err := xgbutil.Dial("")
	if err != nil {
		log.Fatalf("Could not connect to X: %v", err)
	}

	keybind.Initialize(X)
	keybind.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			fmt.Println("Key press!")
		}).Connect(X, X.RootWin(), "Mod4-j")

	xevent.Main(X)
}
