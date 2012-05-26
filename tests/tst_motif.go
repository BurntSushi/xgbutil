package main

import "fmt"

import "github.com/BurntSushi/xgbutil"
import "github.com/BurntSushi/xgbutil/ewmh"
import "github.com/BurntSushi/xgbutil/motif"

func main() {
	X, _ := xgbutil.NewConn()

	active, _ := ewmh.ActiveWindowGet(X)

	mh, err := motif.WmHintsGet(X, active)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(mh)
		fmt.Println("Does Chrome want decorations?",
			motif.Decor(X, active, mh))
	}
}
