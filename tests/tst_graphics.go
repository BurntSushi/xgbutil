package main

import (
    "fmt"
    "image"
    "image/color"
    "image/draw"
    "time"

    "io/ioutil"
    "code.google.com/p/freetype-go/freetype"

    "code.google.com/p/x-go-binding/xgb"
    "github.com/BurntSushi/xgbutil"
)

// If we want to save the image as png output,
// these imports need to be added and the writer
// needs to be uncommented below.
// import ( 
    // "image/png" 
    // "os" 
// ) 

var X *xgbutil.XUtil
var Xerr error

func Recovery() {
    if r := recover(); r != nil {
        fmt.Println("ERROR:", r)
        // os.Exit(1) 
    }
}

func WriteText(img draw.Image) {
    // w, h := img.Bounds().Max.X, img.Bounds().Max.Y 
    text := "Hello, world!"
    fontFile := "/usr/share/fonts/TTF/DejaVuSans-Bold.ttf"

    fontBytes, err := ioutil.ReadFile(fontFile)
    if err != nil {
        fmt.Println(err)
        return
    }

    font, err := freetype.ParseFont(fontBytes)
    if err != nil {
        fmt.Println(err)
        return
    }

    fg := image.NewUniform(color.RGBA{0xff, 0xff, 0xff, 0xff})
    // ruler := color.RGBA{0xdd, 0xdd, 0xdd, 0xff} 
    // rgba := image.NewRGBA(image.Rect(0, 0, w, h)) 
    c := freetype.NewContext()
    c.SetDPI(72)
    c.SetFont(font)
    c.SetFontSize(14)
    c.SetClip(img.Bounds())
    c.SetDst(img)
    c.SetSrc(fg)

    // Draw the guidelines
    // for i := 0; i < 200; i++ { 
        // img.Set(10, 10 + i, ruler) 
        // img.Set(10 + i, 10, ruler) 
    // } 

    // Draw the text
    pt := freetype.Pt(14, 75 + c.FUnitToPixelRU(font.UnitsPerEm()))
    // pt := freetype.Pt(8, 127) 
    _, err = c.DrawString(text, pt)
    if err != nil {
        fmt.Println(err)
        return
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
    allBlue := image.NewUniform(color.RGBA{127, 127, 127, 255})
    draw.Draw(dest, dest.Bounds(), allBlue, image.ZP, draw.Src)
    draw.DrawMask(dest, dest.Bounds(), img, image.ZP, mask, image.ZP, draw.Over)

    // Let's try to write some text...
    WriteText(dest)

    // destWriter, err := os.Create("someicon.png") 
    // if err != nil { 
        // fmt.Print("could not create someicon.png") 
        // os.Exit(1) 
    // } 

    // png.Encode(destWriter, dest) 

    // Let's see if we can paint the image we generated above to a window.

    win := X.Conn().NewId()
    gc := X.Conn().NewId()
    scrn := X.Conn().DefaultScreen()

    winMask := uint32(xgb.CWBackPixmap | xgb.CWOverrideRedirect |
                      xgb.CWBackPixel)
    winVals := []uint32{xgb.BackPixmapParentRelative, scrn.BlackPixel, 1}
    X.Conn().CreateWindow(scrn.RootDepth, win, X.RootWin(), 100, 100,
                          uint16(width), uint16(height),
                          0, xgb.WindowClassInputOutput, scrn.RootVisual,
                          winMask, winVals)
    X.Conn().CreateGC(gc, X.RootWin(), xgb.GCForeground,
                      []uint32{scrn.WhitePixel})
    X.Conn().MapWindow(win)

    // try paitning the image we created above...
    // First we have to transform the image into X format.
    // Then we have to allocate resources for the pixmap.
    // Then we can paint the pixmap.
    // Finally, we attach that pixmap as the "BackPixmap" of our window above.
    // (And free pixmap right thereafter, of course.)
    imgData := make([]byte, width * height * 4)
    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            r, g, b, a := dest.At(x, y).RGBA()
            i := 4 * (x + (y * height))
            imgData[i + 0] = byte(b)
            imgData[i + 1] = byte(g)
            imgData[i + 2] = byte(r)
            imgData[i + 3] = byte(a)
        }
    }
    pix := X.Conn().NewId()
    X.Conn().CreatePixmap(scrn.RootDepth, pix, X.RootWin(),
                          uint16(width), uint16(height))
    X.Conn().PutImage(xgb.ImageFormatZPixmap, pix, gc,
                      uint16(width), uint16(height), 0, 0, 0, 24, imgData)
    X.Conn().ChangeWindowAttributes(win, uint32(xgb.CWBackPixmap),
                                    []uint32{uint32(pix)})
    X.Conn().ClearArea(false, win, 0, 0, 0, 0)
    X.Conn().FreePixmap(pix)

    time.Sleep(20 * time.Second)
}

