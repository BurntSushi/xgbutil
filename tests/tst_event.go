package main

import (
    "fmt"
    // "os" 
)

import (
    "code.google.com/p/x-go-binding/xgb"

    "github.com/BurntSushi/xgbutil"
    "github.com/BurntSushi/xgbutil/xevent"
    "github.com/BurntSushi/xgbutil/xprop"
    "github.com/BurntSushi/xgbutil/xwindow"
)

func MyCallback(X *xgbutil.XUtil, e xevent.PropertyNotifyEvent) {
    atomName, err := xprop.AtomName(X, e.Atom)
    if err != nil {
        panic(err)
    } else {
        fmt.Printf("property %s, state %v\n", atomName, e.State)
    }
}

func MyCallback2(X *xgbutil.XUtil, e xevent.MappingNotifyEvent) {
    fmt.Printf("MappingNotify | Request = %v, FirstKeycode = %v, Count = %v\n",
               e.Request, e.FirstKeycode, e.Count)
}

func main() {
    fmt.Printf("Starting...\n")
    X, _ := xgbutil.Dial("")

    xwindow.Listen(X, X.RootWin(), xgb.EventMaskPropertyChange)

    cb := xevent.PropertyNotifyFun(MyCallback)
    cb.Connect(X, X.RootWin())

    cb2 := xevent.MappingNotifyFun(MyCallback2)
    cb2.Connect(X, xgbutil.NoWindow)

    xevent.Main(X)

    // testEvent := xevent.PropertyNotifyEvent{ 
        // &xgb.PropertyNotifyEvent{1, 6, 0, 1}} 
//  
    // cb := xevent.PropertyNotifyFun(MyCallback) 
    // cb.Run(X, testEvent) 

    // for { 
        // reply, err := X.Conn().WaitForEvent() 
        // if err != nil { 
            // fmt.Printf("ERROR: %v\n", err) 
            // os.Exit(1) 
        // } 
//  
        // fmt.Printf("EVENT: %T - ", reply) 
        // switch event := reply.(type) { 
        // case xgb.PropertyNotifyEvent: 
            // xuEvent := xevent.PropertyNotifyEvent{&event} 
            // cb.Run(X, xuEvent) 
        // default: 
            // fmt.Printf("ERROR: UNSUPPORTED EVENT TYPE") 
            // os.Exit(1) 
        // } 
    // } 
}

