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
    "image/png"
    "image/draw"
    "os"
    "time"

    "code.google.com/p/graphics-go/graphics"

    "github.com/BurntSushi/xgbutil"
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

    srcReader, err := os.Open("close.png")
    if err != nil {
        fmt.Println("%s is not readable.", "old.png")
    }

    simg, err := png.Decode(srcReader)
    if err != nil {
        fmt.Println("Could not decode %s.", "old.png")
    }

    img := image.NewRGBA(image.Rect(0, 0, 100, 100))
    graphics.Scale(img, simg)

    dimg := image.NewRGBA(img.Bounds())
    draw.Draw(dimg, dimg.Bounds(), img, image.ZP, draw.Src)
    dmask := image.NewRGBA(img.Bounds())
    draw.Draw(dmask, img.Bounds(), image.NewUniform(color.Alpha{255}),
              image.ZP, draw.Src)
    dest := xgraphics.BlendBg(img, dmask, 100, color.RGBA{255, 255, 255, 255})

    // Let's try to write some text...
    // xgraphics.DrawText(dest, 5, 5, color.RGBA{255, 255, 255, 255}, 10, 
                       // fontFile, "Hello, world!") 
//  
    // tw, th, err := xgraphics.TextExtents(fontFile, 11, "Hiya") 
    // fmt.Println(tw, th, err) 

    win := xgraphics.CreateImageWindow(X, dest, 3940, 400)
    X.Conn().MapWindow(win)

    time.Sleep(20 * time.Second)
}

