package main

import (
    "flag"
    "fmt"
    "image"
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

    destImg := image.NewNRGBA(image.Rect(0, 0, neww, newh))
    graphics.Scale(destImg, srcImg)
    png.Encode(destWriter, destImg)
}

