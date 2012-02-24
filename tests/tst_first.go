package main

import (
    "fmt"
    "github.com/BurntSushi/xgbutil"
)

func main() {
    conn, err := xgbutil.Dial("")
    if err != nil {
        panic(err)
    }

    fmt.Println(conn)

    geom := conn.EwmhDesktopGeometry()
    active := conn.EwmhActiveWindow()
    desktops := conn.EwmhDesktopNames()
    curdesk := conn.EwmhCurrentDesktop()

    fmt.Printf("Active window: %x\n", active)
    fmt.Printf("Current desktop: %d\n", conn.EwmhCurrentDesktop())
    fmt.Printf("Client list: %v\n", conn.EwmhClientList())
    fmt.Printf("Desktop geometry: (width: %d, height: %d)\n",
               geom.Width, geom.Height)
    fmt.Printf("Active window name: %s\n", conn.EwmhWmName(active))
    fmt.Printf("Desktop names: %s\n", conn.EwmhDesktopNames())
    fmt.Printf("Current desktop: %s\n", desktops[curdesk])

    fmt.Printf("\nChanging current desktop to 25 from %d\n", curdesk)
    conn.EwmhCurrentDesktopSet(25)
    fmt.Printf("Current desktop is now: %d\n", conn.EwmhCurrentDesktop())
}

