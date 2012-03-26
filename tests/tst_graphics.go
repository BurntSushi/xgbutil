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
    // "image" 
    "image/color"
    // "image/draw" 
    // "time" 

    "github.com/BurntSushi/xgbutil"
    "github.com/BurntSushi/xgbutil/ewmh"
    "github.com/BurntSushi/xgbutil/xgraphics"
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

    active, _ := ewmh.ActiveWindowGet(X)
    icons, _ := ewmh.WmIconGet(X, active)

    var width, height int = 200, 200

    work := xgraphics.FindBestIcon(width, height, icons)
    if work != nil {
        fmt.Printf("Working with icon (%d, %d)\n", work.Width, work.Height)
    } else {
        fmt.Println("No good icon... :-(")
        return
    }

    eimg, emask := xgraphics.EwmhIconToImage(work)

    img, mask := xgraphics.Scale(eimg, emask, int(width), int(height))

    dest := xgraphics.BlendBg(img, mask, 100, color.RGBA{0, 0, 255, 255})

    // Let's try to write some text...
    // xgraphics.DrawText(dest, 50, 50, color.RGBA{255, 255, 255, 255}, 20, 
                       // fontFile, "Hello, world!") 

    tw, th, err := xgraphics.TextExtents(fontFile, 11, "Hiya")
    fmt.Println(tw, th, err)

    win := xgraphics.CreateImageWindow(X, dest, 3940, 400)
    X.Conn().MapWindow(win)

    // time.Sleep(20 * time.Second) 
    select {}
}

