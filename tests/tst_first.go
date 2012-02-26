package main

import (
    "fmt"
    "math/rand"
    // "os" 
    "time"

    "code.google.com/p/x-go-binding/xgb"
    "github.com/BurntSushi/xgbutil"
)

var X *xgbutil.XUtil
var Xerr error

func Recovery() {
    if r := recover(); r != nil {
        fmt.Println("ERROR:", r)
        // os.Exit(1) 
    }
}

func main() {
    defer Recovery()

    X, Xerr = xgbutil.Dial("")
    if Xerr != nil {
        panic(Xerr)
    }

    fmt.Println(X)

    geom := X.EwmhDesktopGeometry()
    active := X.EwmhActiveWindow()
    desktops := X.EwmhDesktopNames()
    curdesk := X.EwmhCurrentDesktop()

    fmt.Printf("Active window: %x\n", active)
    fmt.Printf("Current desktop: %d\n", X.EwmhCurrentDesktop())
    fmt.Printf("Client list: %v\n", X.EwmhClientList())
    fmt.Printf("Desktop geometry: (width: %d, height: %d)\n",
               geom.Width, geom.Height)
    fmt.Printf("Active window name: %s\n", X.EwmhWmName(active))
    fmt.Printf("Desktop names: %s\n", X.EwmhDesktopNames())

    var desk string
    if curdesk < uint32(len(desktops)) {
        desk = desktops[curdesk]
    } else {
        desk = string(curdesk)
    }
    fmt.Printf("Current desktop: %s\n", desk)

    // fmt.Printf("\nChanging current desktop to 25 from %d\n", curdesk) 
    X.EwmhCurrentDesktopSet(curdesk)
    // fmt.Printf("Current desktop is now: %d\n", X.EwmhCurrentDesktop()) 

    var newactive xgb.Id = 0x2e00016
    fmt.Printf("Setting active win to %x\n", newactive)
    X.EwmhActiveWindowReq(newactive)

    rand.Seed(int64(time.Now().Nanosecond()))
    randStr := make([]byte, 20)
    for i, _ := range randStr {
        if rf := rand.Float32(); rf < 0.33 {
            randStr[i] = byte('a' + rand.Intn('z' - 'a'))
        } else if rf < 0.66 {
            randStr[i] = byte('A' + rand.Intn('Z' - 'A'))
        } else {
            randStr[i] = ' '
        }
    }

    X.EwmhWmNameSet(active, string(randStr))
    fmt.Printf("New name: %s\n", X.EwmhWmName(active))

    // deskNames := X.EwmhDesktopNames() 
    // fmt.Printf("Desktop names: %s\n", deskNames) 
    // deskNames[len(deskNames) - 1] = "xgbutil" 
    // X.EwmhDesktopNamesSet(deskNames) 
    // fmt.Printf("Desktop names: %s\n", X.EwmhDesktopNames()) 

    icons := X.EwmhWmIcon(active)
    fmt.Printf("Active window's (%x) icon data: (length: %v)\n", 
               active, len(icons))
    for _, icon := range icons {
        fmt.Printf("\t(%d, %d)", icon.Width, icon.Height)
        fmt.Printf(" :: %d == %d\n", icon.Width * icon.Height, len(icon.Data))
    }
}

