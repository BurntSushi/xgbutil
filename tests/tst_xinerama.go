package main

import "fmt"
import "github.com/BurntSushi/xgbutil"
import "github.com/BurntSushi/xgbutil/ewmh"
import "github.com/BurntSushi/xgbutil/xinerama"

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

