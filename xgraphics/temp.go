package xgraphics

import (
	"image"
	"image/color"
	"image/draw"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
)

// A couple of temporary functions to facilitate transitioning Wingo from
// the old xgbutil to the new one.

func CreatePixmap(X *xgbutil.XUtil, img image.Image) xproto.Pixmap {
	ximg := NewConvert(X, img)
	ximg.CreatePixmap()
	return ximg.Pixmap
}

func PaintImg(xu *xgbutil.XUtil, win xproto.Window, img image.Image) {
	pix := CreatePixmap(xu, img)
	xproto.ChangeWindowAttributes(xu.Conn(), win, uint32(xproto.CwBackPixmap),
		[]uint32{uint32(pix)})
	xproto.ClearArea(xu.Conn(), false, win, 0, 0, 0, 0)
	FreePixmap(xu, pix)
}

// Blend "blends" img with mask into dest at position (x, y) with
// transparency alpha.
func BlendOld(dest draw.Image, img image.Image, mask draw.Image,
	alpha, x, y int) {

	transClr := uint8((float64(alpha) / 100.0) * 255.0)
	blendMask := image.NewUniform(color.Alpha{transClr})

	if mask != nil {
		draw.DrawMask(mask, mask.Bounds(), mask, image.ZP, blendMask, image.ZP,
			draw.Src)
	}

	width, height := GetDim(img)
	rect := image.Rect(x, y, width+x, height+y)
	if mask != nil {
		draw.DrawMask(dest, rect, img, image.ZP, mask, image.ZP, draw.Over)
	} else {
		draw.DrawMask(dest, rect, img, image.ZP, blendMask, image.ZP, draw.Over)
	}
}

// BlendBg "blends" img with mask into a background with color clr with
// transparency, where alpha is a number 0-100 where 0 is completely
// transparent and 100 is completely opaque.
// It is very possible that I'm doing more than I need to here, but this
// was the only way I could get it to work.
func BlendBg(img image.Image, mask draw.Image, alpha int,
	clr color.RGBA) *image.RGBA {
	dest := image.NewRGBA(img.Bounds())
	draw.Draw(dest, dest.Bounds(), image.NewUniform(clr), image.ZP, draw.Src)
	BlendOld(dest, img, mask, alpha, 0, 0)
	return dest
}

// DrawText takes an image and, using the freetype package, writes text in the
// position specified on to the image. A color.Color, a font size and a font  
// must also be specified.
// Finally, the (x, y) coordinate advanced by the text extents is returned.
func DrawText(img draw.Image, x int, y int, clr color.Color, fontSize float64,
	font *truetype.Font, text string) (int, int, error) {

	// Create a solid color image
	textClr := image.NewUniform(clr)

	// Set up the freetype context... mostly boiler plate
	c := ftContext(font, fontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(textClr)

	// Now let's actually draw the text...
	pt := freetype.Pt(x, y+c.FUnitToPixelRU(font.UnitsPerEm()))
	newpt, err := c.DrawString(text, pt)
	if err != nil {
		return 0, 0, err
	}

	// i think this is right...
	return int(newpt.X / 256), int(newpt.Y / 256), nil
}

// GetDim gets the width and height of an image
func GetDim(img image.Image) (int, int) {
	bounds := img.Bounds()
	return bounds.Max.X - bounds.Min.X, bounds.Max.Y - bounds.Min.Y
}
