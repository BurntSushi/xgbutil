/*
    Package xgraphics makes drawing graphics to windows a bit easier.
    It uses the method of paiting a background pixmap.

    This packages requires the freetype and graphics packages from Google.

    This package is probably incomplete. I admit that I designed it with
    my window manager as a use case.
*/
package xgraphics

import (
    "image"
    "image/color"
    "image/draw"
    "io/ioutil"
)

import "code.google.com/p/freetype-go/freetype"
import "code.google.com/p/freetype-go/freetype/truetype"

import "code.google.com/p/jamslam-x-go-binding/xgb"

import "github.com/BurntSushi/xgbutil"
import "github.com/BurntSushi/xgbutil/ewmh"

// DrawText takes an image and, using the freetype package, writes text in the
// position specified on to the image. A color.Color, a font size and a font  
// must also be specified. For example, /usr/share/fonts/TTF/DejaVuSans-Bold.ttf
func DrawText(img draw.Image, x int, y int, clr color.Color, fontSize float64,
              fontFile string, text string) error {
    // get our truetype.Font
    font, err := parseFont(fontFile)
    if err != nil {
        return err
    }

    // Create a solid color image
    textClr := image.NewUniform(clr)

    // Set up the freetype context... mostly boiler plate
    c := ftContext(font, fontSize)
    c.SetClip(img.Bounds())
    c.SetDst(img)
    c.SetSrc(textClr)

    // Now let's actually draw the text...
    pt := freetype.Pt(x, y + c.FUnitToPixelRU(font.UnitsPerEm()))
    _, err = c.DrawString(text, pt)
    if err != nil {
        return err
    }

    return nil
}

// Returns the width and height extents of a string given a font.
// TODO: This does not currently account for multiple lines. It may never do so.
func TextExtents(fontFile string, fontSize float64,
                 text string) (width int, height int, err error) {
    // get our truetype.Font
    font, err := parseFont(fontFile)
    if err != nil {
        return 0, 0, err
    }

    // We need a context to calculate the extents
    c := ftContext(font, fontSize)

    emSquarePix := c.FUnitToPixelRU(font.UnitsPerEm())
    return len(text) * emSquarePix, emSquarePix, nil
}

// ftContext does the boiler plate to create a freetype context
func ftContext(font *truetype.Font, fontSize float64) *freetype.Context {
    c := freetype.NewContext()
    c.SetDPI(72)
    c.SetFont(font)
    c.SetFontSize(fontSize)

    return c
}

// parseFont reads a font file and creates a freetype.Font type
func parseFont(fontFile string) (*truetype.Font, error) {
    fontBytes, err := ioutil.ReadFile(fontFile)
    if err != nil {
        return nil, err
    }

    font, err := freetype.ParseFont(fontBytes)
    if err != nil {
        return nil, err
    }

    return font, nil
}

// CreateImageWindow automatically creates a window with the same size as
// the image given, positions it according to the x,y coordinates given,
// paints the image onto the background of the image, and returns the window
// id. It does *not* map the window for you though. You'll need to that
// with `X.Conn().MapWindow(window_id)`.
// XXX: This will likely change to include the window masks and vals as
// parameters.
func CreateImageWindow(xu *xgbutil.XUtil, img image.Image,
                       x int16, y int16) xgb.Id {
    win := xu.Conn().NewId()
    scrn := xu.Screen()
    width, height := getDim(img)

    winMask := uint32(xgb.CWBackPixmap | xgb.CWOverrideRedirect)
    winVals := []uint32{xgb.BackPixmapParentRelative, 1}
    xu.Conn().CreateWindow(scrn.RootDepth, win, xu.RootWin(), x, y,
                           uint16(width), uint16(height),
                           0, xgb.WindowClassInputOutput, scrn.RootVisual,
                           winMask, winVals)

    PaintImg(xu, win, img)

    return win
}

// PaintImg will slap the given image as a background pixmap into the given
// window.
// TODO: There is currently a limitation in XGB (not xgbutil) that prevents
// requests from being bigger than (2^16 * 4) bytes. (This is caused by silly
// X nonsense.) To fix this, XGB needs to work around it, but it isn't quite
// clear how that should be done yet.
// Therefore, try to keep images less than 250x250, otherwise X will stomp
// on you. And it will hurt.
func PaintImg(xu *xgbutil.XUtil, win xgb.Id, img image.Image) {
    // gather up image data in the form X wants it... so picky
    width, height := getDim(img)
    imgData := make([]byte, width * height * 4)
    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            r, g, b, a := img.At(x, y).RGBA()
            i := 4 * (x + (y * height))
            imgData[i + 0] = byte(b)
            imgData[i + 1] = byte(g)
            imgData[i + 2] = byte(r)
            imgData[i + 3] = byte(a)
        }
    }

    // The hard part is over, boiler-plate time!
    pix := xu.Conn().NewId()
    xu.Conn().CreatePixmap(xu.Screen().RootDepth, pix, xu.RootWin(),
                           uint16(width), uint16(height))
    xu.Conn().PutImage(xgb.ImageFormatZPixmap, pix, xu.GC(),
                       uint16(width), uint16(height), 0, 0, 0, 24, imgData)
    xu.Conn().ChangeWindowAttributes(win, uint32(xgb.CWBackPixmap),
                                     []uint32{uint32(pix)})
    xu.Conn().ClearArea(false, win, 0, 0, 0, 0)
    xu.Conn().FreePixmap(pix)
}

// getDim gets the width and height of an image
func getDim(img image.Image) (int, int) {
    bounds := img.Bounds()
    return bounds.Max.X - bounds.Min.X, bounds.Max.Y - bounds.Min.Y
}

// BlendBg "blends" img with mask into a background with color clr with
// transparency, where transparency is a number 0-100 where 0 is completely
// transparent and 100 is completely opaque.
// It is very possible that I'm doing more than I need to here, but this
// was the only way I could get it to work.
func BlendBg(img image.Image, mask draw.Image, transparency int,
             clr color.RGBA) (dest *image.RGBA) {
    transClr := uint8((float64(transparency) / 100.0) * 255.0)
    blendMask := image.NewUniform(color.Alpha{transClr})
    draw.DrawMask(mask, mask.Bounds(), mask, image.ZP, blendMask, image.ZP,
                  draw.Src)

    dest = image.NewRGBA(img.Bounds())
    draw.Draw(dest, dest.Bounds(), image.NewUniform(clr), image.ZP, draw.Src)
    draw.DrawMask(dest, dest.Bounds(), img, image.ZP, mask, image.ZP, draw.Over)

    return
}

// EwmhIconToImage takes a ewmh.WmIcon and converts it to an image and
// an alpha mask. A ewmh.WmIcon is in ARGB order, and the image package wants
// things in RGBA order. (What makes things is worse is when it comes time
// to paint the image to the screen, X wants it in BGR order. *facepalm*.)
func EwmhIconToImage(icon *ewmh.WmIcon) (img *image.RGBA, mask *image.RGBA) {
    width, height := int(icon.Width), int(icon.Height)
    img = image.NewRGBA(image.Rect(0, 0, width, height))
    mask = image.NewRGBA(image.Rect(0, 0, width, height))

    for x := 0; x < width; x++ {
        for y := 0; y < height; y++ {
            argb := icon.Data[x + (y * height)]
            alpha := argb >> 24
            red := ((alpha << 24) ^ argb) >> 16
            green := (((alpha << 24) + (red << 16)) ^ argb) >> 8
            blue := (((alpha << 24) + (red << 16) + (green << 8)) ^ argb) >> 0

            c := color.RGBA{uint8(red), uint8(green), uint8(blue), uint8(alpha)}

            img.SetRGBA(x, y, c)
            mask.Set(x, y, color.Alpha{uint8(alpha)})
        }
    }

    return
}

