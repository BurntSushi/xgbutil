// This test file is a giant soup that had one goal: act as a testing ground
// to get various graphical routines working in X. Namely:
//
// Blending one image on top of another (with alpha).
// Painting an arbitrary image into a window.
// Drawing text on to an image.
//
// All of this is done here successfully. Most of this file will be split up
// into nicer pieces in my window manager.
package main

import (
    "fmt"
    "image"
    "image/color"
    // "image/png" 
    // "image/draw" 
    // "os" 
    "time"

    "code.google.com/p/graphics-go/graphics"

    "burntsushi.net/go/xgbutil"
    "burntsushi.net/go/xgbutil/xgraphics"
)

var X *xgbutil.XUtil
var Xerr error

func Recovery() {
    if r := recover(); r != nil {
        fmt.Println("ERROR:", r)
        // os.Exit(1) 
    }
}

var fontFile string = "/usr/share/fonts/TTF/DejaVuSans-Bold.ttf"

func main() {
    defer Recovery()

    X, Xerr = xgbutil.Dial("")
    if Xerr != nil {
        panic(Xerr)
    }

    simg, err := xgraphics.LoadPngFromFile("openbox.png")
    if err != nil {
        fmt.Println(err)
    }

    img := image.NewRGBA(image.Rect(0, 0, 50, 101))
    graphics.Scale(img, simg)

    dest := xgraphics.BlendBg(img, nil, 100, color.RGBA{255, 255, 255, 255})

    win := xgraphics.CreateImageWindow(X, dest, 3940, 400)
    X.Conn().MapWindow(win)

    time.Sleep(20 * time.Second)
}

