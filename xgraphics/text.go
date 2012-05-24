package xgraphics

import (
	"image"
	"image/color"
	"io"
	"io/ioutil"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
)

// DrawText takes an image and, using the freetype package, writes text in the
// position specified on to the image. A color.Color, a font size and a font  
// must also be specified.
// Finally, the (x, y) coordinate advanced by the text extents is returned.
func (im *Image) Text(x, y int, clr color.Color, fontSize float64,
	font *truetype.Font, text string) (int, int, error) {

	// Create a solid color image
	textClr := image.NewUniform(clr)

	// Set up the freetype context... mostly boiler plate
	c := ftContext(font, fontSize)
	c.SetClip(im.Bounds())
	c.SetDst(im)
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

// Returns the max width and height extents of a string given a font.
// This is calculated by determining the number of pixels in an "em" unit
// for the given font, and multiplying by the number of characters in 'text'.
// Since a particular character may be smaller than one "em" unit, this has
// a tendency to overestimate the extents.
// It is provided because I do not know how to calculate the precise extents
// using freetype-go.
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
func ParseFont(fontReader io.Reader) (*truetype.Font, error) {
	fontBytes, err := ioutil.ReadAll(fontReader)
	if err != nil {
		return nil, err
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	return font, nil
}
