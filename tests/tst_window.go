package main

import "fmt"
import "time"
import "github.com/BurntSushi/xgbutil"

func main() {
    fmt.Println("Step 1")
    X, _ := xgbutil.Dial("")
    fmt.Println("Step 2")

    active, _ := X.EwmhActiveWindow()
    fmt.Println("Step 3")
    geom, err := X.GetGeometry(active)
    fmt.Println("Step 4")
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(geom)
    }
    fmt.Println("Step 5")

    X.EwmhWmStateReqExtra(active, xgbutil.EwmhStateToggle,
                          "_NET_WM_STATE_MAXIMIZED_VERT",
                          "_NET_WM_STATE_MAXIMIZED_HORZ", 2)
    time.Sleep(time.Millisecond)
    err = X.MoveResize(active, geom.X, geom.Y,
                        geom.Width - 100, geom.Height)
    fmt.Println(err)
    fmt.Printf("\n")

    X.Conn().Close()
}

