package main

import "fmt"
import "github.com/BurntSushi/xgbutil"
import "github.com/BurntSushi/xgbutil/ewmh"
import "github.com/BurntSushi/xgbutil/xinerama"
import "github.com/BurntSushi/xgbutil/xrect"

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

    // Test intersection
    r1 := xrect.Make(0, 0, 100, 100)
    r2 := xrect.Make(100, 100, 100, 100)
    fmt.Println(xrect.IntersectArea(r1, r2))

    // Test largest overlap
    window := xrect.Make(1800, 0, 200, 200)
    fmt.Println(xrect.LargestOverlap(window, heads))
}

