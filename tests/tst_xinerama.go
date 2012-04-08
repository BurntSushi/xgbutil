package main

import "fmt"
import "burntsushi.net/go/xgbutil"
import "burntsushi.net/go/xgbutil/ewmh"
import "burntsushi.net/go/xgbutil/xinerama"

func main() {
	X, _ := xgbutil.Dial("")

	heads, err := xinerama.PhysicalHeads(X)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		for i, head := range heads {
			fmt.Printf("%d - %v\n", i, head)
		}
	}

	wmName, err := ewmh.GetEwmhWM(X)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Printf("Window manager: %s\n", wmName)
	}
}
