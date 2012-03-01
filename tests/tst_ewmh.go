package main

import (
    "fmt"
    "log"
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
        // if xuError, ok := r.(*xgbutil.XError); ok { 
            // if xuError.XGBError != nil { 
                // if xuError.XGBError.Detail == xgb.BadValue { 
                    // log.Println("BadValue") 
                // } else { 
                    // log.Println("WOOYAA") 
                // } 
            // } else { 
                // log.Println(r) 
            // } 
        // } else { 
            // log.Println(r) 
        // } 

        switch err := r.(type) {
        case *xgb.Error:
            log.Println("XGB ERROR:", err)
        case *xgbutil.XError:
            log.Println("XGB UTIL ERROR:", err)
        default: // not our problem, produce stack trace
            panic(err)
        }

        // os.Exit(1) 
    }
}

func main() {
    // defer Recovery() 

    X, Xerr = xgbutil.Dial("")
    if Xerr != nil {
        panic(Xerr)
    }

    fmt.Println(X)

    showDesk, _ := X.EwmhShowingDesktop()
    fmt.Printf("Showing desktop? %v\n", showDesk)

    wmName, err := X.GetEwmhWM()
    if err != nil {
        fmt.Printf("No conforming window manager found... :-(\n")
        fmt.Println(err)
    } else {
        fmt.Printf("Window manager: %s\n", wmName)
    }

    pager := xgb.Id(0x160001e)
    middle := xgb.Id(0x3200016)
    geom, _ := X.EwmhDesktopGeometry()
    active, _ := X.EwmhActiveWindow()
    desktops, _ := X.EwmhDesktopNames()
    curdesk, _ := X.EwmhCurrentDesktop()
    clients, _ := X.EwmhClientList()
    activeName, _ := X.EwmhWmName(active)

    fmt.Printf("Active window: %x\n", active)
    fmt.Printf("Current desktop: %d\n", curdesk)
    fmt.Printf("Client list: %v\n", clients)
    fmt.Printf("Desktop geometry: (width: %d, height: %d)\n",
               geom.Width, geom.Height)
    fmt.Printf("Active window name: %s\n", activeName)
    fmt.Printf("Desktop names: %s\n", desktops)

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

    fmt.Printf("Setting active win to %x\n", middle)
    X.EwmhActiveWindowReq(middle)

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
    newName, _ := X.EwmhWmName(active)
    fmt.Printf("New name: %s\n", newName)

    // deskNames := X.EwmhDesktopNames() 
    // fmt.Printf("Desktop names: %s\n", deskNames) 
    // deskNames[len(deskNames) - 1] = "xgbutil" 
    // X.EwmhDesktopNamesSet(deskNames) 
    // fmt.Printf("Desktop names: %s\n", X.EwmhDesktopNames()) 

    supported, _ := X.EwmhSupported()
    fmt.Printf("Supported hints: %v\n", supported)
    fmt.Printf("Setting supported hints...\n")
    X.EwmhSupportedSet([]string{"_NET_CLIENT_LIST", "_NET_WM_NAME",
                                "_NET_WM_DESKTOP"})

    numDesks, _ := X.EwmhNumberOfDesktops()
    fmt.Printf("Number of desktops: %d\n", numDesks)
    // X.EwmhNumberOfDesktopsReq(X.EwmhNumberOfDesktops() + 1) 
    // time.Sleep(time.Second) 
    // fmt.Printf("Number of desktops: %d\n", X.EwmhNumberOfDesktops()) 

    viewports, _ := X.EwmhDesktopViewport()
    fmt.Printf("Viewports (%d): %v\n", len(viewports), viewports)

    // viewports[2].X = 50
    // viewports[2].Y = 293 
    // X.EwmhDesktopViewportSet(viewports) 
    // time.Sleep(time.Second) 
//  
    // viewports = X.EwmhDesktopViewport() 
    // fmt.Printf("Viewports (%d): %v\n", len(viewports), viewports) 

    // X.EwmhCurrentDesktopReq(3) 

    visDesks, _ := X.EwmhVisibleDesktops()
    workarea, _ := X.EwmhWorkarea()
    fmt.Printf("Visible desktops: %v\n", visDesks)
    fmt.Printf("Workareas: %v\n", workarea)
    // fmt.Printf("Virtual roots: %v\n", X.EwmhVirtualRoots()) 
    // fmt.Printf("Desktop layout: %v\n", X.EwmhDesktopLayout()) 
    fmt.Printf("Closing window %x\n", 0x2e004c5)
    X.EwmhCloseWindow(0x2e004c5)

    fmt.Printf("Moving/resizing window: %x\n", 0x2e004d0)
    X.EwmhMoveresizeWindow(0x2e004d0, 1920, 30, 500, 500)

    // fmt.Printf("Trying to initiate a moveresize...\n") 
    // X.EwmhWmMoveresize(0x2e004db, xgbutil.EwmhMove) 
    // time.Sleep(5 * time.Second) 
    // X.EwmhWmMoveresize(0x2e004db, xgbutil.EwmhCancel) 

    // fmt.Printf("Stacking window %x...\n", 0x2e00509) 
    // X.EwmhRestackWindow(0x2e00509) 

    fmt.Printf("Requesting frame extents for active window...\n")
    X.EwmhRequestFrameExtents(active)

    actOpacity, _ := X.EwmhWmWindowOpacity(X.ParentWindow(active))
    // actOpacity2 := X.EwmhWmWindowOpacity(X.ParentWindow(X.EwmhActiveWindow())) 
    fmt.Printf("Opacity for active window: %f\n", actOpacity)
    // fmt.Printf("Opacity for real active window: %f\n", actOpacity2) 
    // X.EwmhWmWindowOpacitySet(X.ParentWindow(active), 0.5) 

    activeDesk, _ := X.EwmhWmDesktop(active)
    activeType, _ := X.EwmhWmWindowType(active)
    fmt.Printf("Active window's desktop: %d\n", activeDesk)
    fmt.Printf("Active's types: %v\n", activeType)
    // fmt.Printf("Pager's types: %v\n", X.EwmhWmWindowType(0x180001e)) 

    // fmt.Printf("Pager's state: %v\n", X.EwmhWmState(0x180001e)) 

    // X.EwmhWmStateReq(active, xgbutil.EwmhStateToggle, "_NET_WM_STATE_HIDDEN") 
    // X.EwmhWmStateReqExtra(active, xgbutil.EwmhStateToggle, 
                          // "_NET_WM_STATE_MAXIMIZED_VERT", 
                          // "_NET_WM_STATE_MAXIMIZED_HORZ", 2) 

    activeAllowed, _ := X.EwmhWmAllowedActions(active)
    fmt.Printf("Allowed actions on active: %v\n", activeAllowed)

    struts, err := X.EwmhWmStrut(pager)
    if err != nil {
        fmt.Printf("Pager struts: %v\n", err)
    } else {
        fmt.Printf("Pager struts: %v\n", struts)
    }

    pstruts, err := X.EwmhWmStrutPartial(pager)
    if err != nil {
        fmt.Printf("Pager struts partial: %v - %v\n", pstruts, err)
    } else {
        fmt.Printf("Pager struts partial: %v\n", pstruts.BottomStartX)
    }

    // fmt.Printf("Icon geometry for active: %v\n", X.EwmhWmIconGeometry(active)) 

    icons, _ := X.EwmhWmIcon(active)
    fmt.Printf("Active window's (%x) icon data: (length: %v)\n",
               active, len(icons))
    for _, icon := range icons {
        fmt.Printf("\t(%d, %d)", icon.Width, icon.Height)
        fmt.Printf(" :: %d == %d\n", icon.Width * icon.Height, len(icon.Data))
    }
    // fmt.Printf("Now set them again...\n") 
    // X.EwmhWmIconSet(active, icons[:len(icons) - 1]) 
}

