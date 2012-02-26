package main

import (
    "fmt"
    "image"
    "image/color"
    "image/draw"
    "image/png"
    "os"

    // "code.google.com/p/x-go-binding/xgb" 
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

    active := X.EwmhActiveWindow()
    icons := X.EwmhWmIcon(active)
    fmt.Printf("Active window's (%x) icon data: (length: %v)\n", 
               active, len(icons))
    for _, icon := range icons {
        fmt.Printf("\t(%d, %d)", icon.Width, icon.Height)
        fmt.Printf(" :: %d == %d\n", icon.Width * icon.Height, len(icon.Data))
    }

    work := icons[3]
    fmt.Printf("Working with (%d, %d)\n", work.Width, work.Height)

    width, height := int(work.Width), int(work.Height)
    img := image.NewRGBA(image.Rect(0, 0, width, height))
    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            argb := work.Data[x + (y * height)]
            alpha := argb >> 24
            red := ((alpha << 24) ^ argb) >> 16
            green := (((alpha << 24) + (red << 16)) ^ argb) >> 8
            blue := (((alpha << 24) + (red << 16) + (green << 8)) ^ argb) >> 0

            c := color.RGBA{
                R: uint8(red),
                G: uint8(green),
                B: uint8(blue),
                A: uint8(alpha),
            }
            img.SetRGBA(x, y, c)
        }
    }

    var mask image.Image
    // mask = image.NewUniform(color.Alpha{127}) 
    dest := image.NewRGBA(image.Rect(0, 0, width, height))
    allBlue := image.NewUniform(color.RGBA{0, 0, 255, 255})
    draw.Draw(dest, dest.Bounds(), allBlue, image.ZP, draw.Src)
    draw.DrawMask(dest, dest.Bounds(), img, image.ZP, mask, image.ZP, draw.Over)

    destWriter, err := os.Create("someicon.png")
    if err != nil {
        fmt.Print("could not create someicon.png")
        os.Exit(1)
    }

    png.Encode(destWriter, dest)
}

