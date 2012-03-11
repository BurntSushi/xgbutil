package main

import "fmt"
import "time"
import "github.com/BurntSushi/xgbutil"
import "github.com/BurntSushi/xgbutil/ewmh"
import "github.com/BurntSushi/xgbutil/xwindow"

func main() {
    X, _ := xgbutil.Dial("")

    active, _ := ewmh.ActiveWindowGet(X)
    geom, err := xwindow.GetGeometry(X, active)
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(geom)
    }

    // ewmh.WmStateReqExtra(X, active, ewmh.StateToggle, 
                         // "_NET_WM_STATE_MAXIMIZED_VERT", 
                         // "_NET_WM_STATE_MAXIMIZED_HORZ", 2) 
    time.Sleep(time.Millisecond)
    // err = xwindow.MoveResize(X, active, geom.X, geom.Y, 
                             // geom.Width - 100, geom.Height) 
    fmt.Println(err)
    fmt.Printf("\n")

    rgeom, err := xwindow.RawGeometry(X, X.RootWin())
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(rgeom)
    }

    X.Conn().Close()
}

