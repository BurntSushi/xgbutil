package xgraphics

import (
	"image"
	"image/color"
	"image/draw"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
)

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

// Returns the width and height extents of a string given a font.
// TODO: This does not currently account for multiple lines. It may never do so.
func TextMaxExtents(font *truetype.Font, fontSize float64,
	text string) (width int, height int, err error) {

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

// ParseFont reads a font file and creates a freetype.Font type
func ParseFont(fontBytes []byte) (*truetype.Font, error) {
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	return font, nil
}
