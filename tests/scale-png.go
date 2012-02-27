// This is simply testing Google's graphics library that scales images.
package main

import (
    "flag"
    "fmt"
    "image"
    "image/color"
    "image/draw"
    "image/png"
    "os"
    "strconv"

    "code.google.com/p/graphics-go/graphics"
)

func errMsg(msg string, vals ...interface{}) {
    fmt.Fprintf(os.Stderr, msg, vals...)
    fmt.Fprintf(os.Stderr, "\n")
    os.Exit(1)
}

func main() {
    flag.Parse()
    args := flag.Args()

    if flag.NArg() != 4 {
        errMsg("Usage: scale-png src-png new_width new_height dest-png")
    }

    pngsrc := args[0]
    pngdest := args[3]
    neww, errw := strconv.Atoi(args[1])
    newh, errh := strconv.Atoi(args[2])
    if errw != nil || errh != nil {
        errMsg("Could not convert (%s, %s) to integers.", args[1], args[2])
    }

    fmt.Printf("Resizing %s to (%d, %d) and saving as %s\n",
               pngsrc, neww, newh, pngdest)

    srcReader, err := os.Open(pngsrc)
    if err != nil {
        errMsg("%s is not readable.", pngsrc)
    }

    srcImg, err := png.Decode(srcReader)
    if err != nil {
        errMsg("Could not decode %s.", pngsrc)
    }

    destWriter, err := os.Create(pngdest)
    if err != nil {
        errMsg("Could not write %s", pngdest)
    }

    destImg := image.NewRGBA(image.Rect(0, 0, neww, newh))
    graphics.Scale(destImg, srcImg)

    // for x := 0; x < destImg.Bounds().Max.X; x++ { 
        // for y := 0; y < destImg.Bounds().Max.Y; y++ { 
            // c := destImg.At(x, y).(color.RGBA) 
            // blah := 0.3 * 0xff 
            // c.A = uint8(blah) 
            // destImg.SetRGBA(x, y, c) 
        // } 
    // } 

    // finalDest := image.NewRGBA(image.Rect(0, 0, neww, newh)) 
    // blue := color.RGBA{255, 255, 255, 255} 
    // draw.Draw(finalDest, finalDest.Bounds(), image.NewUniform(blue), 
              // image.ZP, draw.Src) 

    // Create a transparency mask
    mask := image.NewUniform(color.Alpha16{32767})

    // Now blend our scaled image 
    draw.DrawMask(destImg, destImg.Bounds(), destImg, image.ZP, mask,
                  image.ZP, draw.Src)

    png.Encode(destWriter, destImg)

    fmt.Printf("Type of src: %T\n", srcImg)
    fmt.Printf("Source opaque? %v\n", srcImg.(*image.RGBA).Opaque())
    fmt.Printf("Destination opaque? %v\n", destImg.Opaque())
}

