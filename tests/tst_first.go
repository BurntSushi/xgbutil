package main

import (
    "fmt"
    "github.com/BurntSushi/xgbutil"
)

func main() {
    conn, _ := xgbutil.Dial("")
    fmt.Println(conn)

    geom := conn.Ewmh_desktop_geometry()
    active := conn.Ewmh_active_window()
    desktops := conn.Ewmh_desktop_names()
    curdesk := conn.Ewmh_current_desktop()

    fmt.Printf("Active window: %x\n", active)
    fmt.Printf("Current desktop: %d\n", conn.Ewmh_current_desktop())
    fmt.Printf("Client list: %v\n", conn.Ewmh_client_list())
    fmt.Printf("Desktop geometry: (width: %d, height: %d)\n", 
               geom.Width, geom.Height)
    fmt.Printf("Active window name: %s\n", conn.Ewmh_wm_name(active))
    fmt.Printf("Desktop names: %s\n", conn.Ewmh_desktop_names())
    fmt.Printf("Current desktop: %s\n", desktops[curdesk])

    fmt.Printf("\nChanging current desktop to 25 from %d\n", curdesk)
    conn.Ewmh_current_desktop_set(25)
    fmt.Printf("Current desktop is now: %d\n", conn.Ewmh_current_desktop())
}

