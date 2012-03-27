package main

import (
    "fmt"
    "log"
)

import (
    "github.com/BurntSushi/xgbutil"
    "github.com/BurntSushi/xgbutil/keybind"
    "github.com/BurntSushi/xgbutil/xevent"
)

func main() {
    X, err := xgbutil.Dial("")
    if err != nil {
        log.Fatalf("Could not connect to X: %v", err)
    }

    keybind.Initialize(X)
    keybind.KeyPressFun(
        func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
            fmt.Println("Key press!")
    }).Connect(X, X.RootWin(), "Mod4-j")

    xevent.Main(X)
}

