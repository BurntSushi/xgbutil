package xgraphics

/*
xgraphics/image.go contains an implementation of the draw.Image interface.

RGBA could feasibly be used, but the representation of image data is dependent
upon the configuration of the X server.

For the time being, I'm hard-coding a lot of that configuration for the common
case. Namely:

Byte order: least significant byte first
Depth: 24
Bits per pixel: 32

This will have to be fixed for this to be truly compatible with any X server.

Most of the code is based heavily on the implementation of common images in
the Go standard library.
*/

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"os"

	"code.google.com/p/graphics-go/graphics"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
)

// Model for the BGRA color type.
var BGRAModel color.Model = color.ModelFunc(bgraModel)

type Image struct {
	// X images must be tied to an X connection.
	X *xgbutil.XUtil

	// X images must also be tied to a pixmap (its drawing surface).
	// Calls to 'Draw' will draw data to this pixmap.
	// Calls to 'Paint' will tell X to show the pixmap on some window.
	Pixmap xproto.Pixmap

	// Pix holds the image's pixels in BGRA order, so that they don't need
	// to be swapped for every PutImage request.
	Pix []uint8

	// Stride corresponds to the number of elements in Pix between two pixels
	// that are vertically adjacent.
	Stride int

	// The geometry of the image.
	Rect image.Rectangle
}

// New returns a new instance of Image with colors initialized to black
// for the geometry given.
// New will also create an X pixmap. When you are no longer using this
// image, you should call Destroy so that the X pixmap can be freed on the
// X server.
// (Generating a pixmap id can cause an error, so this call could return
// an error.)
func New(X *xgbutil.XUtil, r image.Rectangle) (*Image, error) {
	w, h := r.Dx(), r.Dy()

	pid, err := xproto.NewPixmapId(X.Conn())
	if err != nil {
		return nil, err
	}

	// Now actually create the pixmap.
	err = xproto.CreatePixmapChecked(X.Conn(), X.Screen().RootDepth, pid,
		xproto.Drawable(X.RootWin()), uint16(w), uint16(h)).Check()
	if err != nil {
		return nil, err
	}

	buf := make([]uint8, 4*w*h)
	return &Image{
		X:      X,
		Pixmap: pid,
		Pix:    buf,
		Stride: 4 * w,
		Rect:   r,
	}, nil
}

// NewConvert converts any image satisfying the image.Image interface to an
// xgraphics.Image type.
func NewConvert(X *xgbutil.XUtil, img image.Image) (*Image, error) {
	if ximg, ok := img.(*Image); ok {
		return ximg, nil // wtf?
	}

	ximg, err := New(X, img.Bounds())
	if err != nil {
		return nil, err
	}
	for x := 0; x < ximg.Rect.Dx(); x++ {
		for y := 0; y < ximg.Rect.Dy(); y++ {
			ximg.Set(x, y, img.At(x, y))
		}
	}
	return ximg, nil
}

// NewEwmhIcon converts EWMH icon data (ARGB) to an xgraphics.Image type.
func NewEwmhIcon(X *xgbutil.XUtil, icon *ewmh.WmIcon) (*Image, error) {
	ximg, err := New(X, image.Rect(0, 0, icon.Width, icon.Height))
	if err != nil {
		return nil, err
	}
	for x := 0; x < ximg.Rect.Dx(); x++ {
		for y := 0; y < ximg.Rect.Dy(); y++ {
			argb := icon.Data[x+(y*ximg.Rect.Dx())]
			clr := BGRA{
				B: uint8(argb & 0x000000ff),
				G: uint8((argb & 0x0000ff00) >> 8),
				R: uint8((argb & 0x00ff0000) >> 16),
				A: uint8(argb >> 24),
			}
			ximg.Set(x, y, clr)
		}
	}
	return ximg, nil
}

// Destroy frees the pixmap resource being used by this image.
// It should be called whenever the image will no longer be drawn or painted.
func (im *Image) Destroy() {
	xproto.FreePixmap(im.X.Conn(), im.Pixmap)
}

// Scale will scale the image to the size provided.
// Note that this will destroy the current pixmap associated with this image
// and create a new one (since pixmaps cannot be resized).
// After scaling, XSurfaceSet will need to be called for each window that
// this image is painted to. (And obviously, XDraw and XPaint.)
// This function may return an error if a new pixmap cannot be allocated.
func (im *Image) Scale(width, height int) (*Image, error) {
	dimg, err := New(im.X, image.Rect(0, 0, width, height))
	if err != nil {
		return nil, err
	}

	graphics.Scale(dimg, im)
	im.Destroy()

	return dimg, nil
}

// WritePng encodes the image to w as a png.
func (im *Image) WritePng(w io.Writer) error {
	return png.Encode(w, im)
}

// SavePng writes the Image to a file with name as a png.
func (im *Image) SavePng(name string) error {
	file, err := os.Create(name)
	if err != nil {
		return err
	}
	return im.WritePng(file)
}

// ColorModel returns the color.Model used by the Image struct.
func (im *Image) ColorModel() color.Model {
	return BGRAModel
}

// Bounds returns the rectangle representing the geometry of Image.
func (im *Image) Bounds() image.Rectangle {
	return im.Rect
}

// At returns the color at the specified pixel.
func (im *Image) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(im.Rect)) {
		return BGRA{}
	}
	i := im.PixOffset(x, y)
	return BGRA{
		B: im.Pix[i],
		G: im.Pix[i+1],
		R: im.Pix[i+2],
		A: im.Pix[i+3],
	}
}

// Set satisfies the draw.Image interface by allowing the color of a pixel
// at (x, y) to be changed.
func (im *Image) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(im.Rect)) {
		return
	}

	i := im.PixOffset(x, y)
	cc := BGRAModel.Convert(c).(BGRA)
	im.Pix[i] = cc.B
	im.Pix[i+1] = cc.G
	im.Pix[i+2] = cc.R
	im.Pix[i+3] = cc.A
}

// SetBGRA is like set, but without the type assertion.
func (im *Image) SetBGRA(x, y int, c BGRA) {
	if !(image.Point{x, y}.In(im.Rect)) {
		return
	}

	i := im.PixOffset(x, y)
	im.Pix[i] = c.B
	im.Pix[i+1] = c.G
	im.Pix[i+2] = c.R
	im.Pix[i+3] = c.A
}

// For transforms every pixel color to the color returned by 'each' given
// an (x, y) position.
func (im *Image) For(each func(x, y int) BGRA) {
	for x := im.Rect.Min.X; x < im.Rect.Max.X; x++ {
		for y := im.Rect.Min.Y; y < im.Rect.Max.Y; y++ {
			im.SetBGRA(x, y, each(x, y))
		}
	}
}

// BlendImage alpha blends the src image (starting at the spt Point) into the
// dest image.
// If you're blending into a solid background color, use BlendBgColor
// instead. (It's more efficient.)
// BlendImage does not (currently) blend with the destination's alpha channel,
// only the source's alpha channel.
func BlendImage(dest draw.Image, src image.Image, sp image.Point) {
	rsrc, dsrc := src.Bounds(), dest.Bounds()
	_, smxx, _, smxy := rsrc.Min.X, rsrc.Max.X, rsrc.Min.Y, rsrc.Max.Y
	dmnx, dmxx, dmny, dmxy := dsrc.Min.X, dsrc.Max.X, dsrc.Min.Y, dsrc.Max.Y

	for sx, dx := sp.X, dmnx; sx < smxx && dx < dmxx; sx, dx = sx+1, dx+1 {
		for sy, dy := sp.Y, dmny; sy < smxy && dy < dmxy; sy, dy = sy+1, dy+1 {
			sr, sg, sb, sa := src.At(sx, sy).RGBA()
			dr, dg, db, _ := dest.At(dx, dy).RGBA()
			alpha := float64(sa) / 255.0

			dest.Set(dx, dy, color.RGBA{
				blend(uint8(sr), uint8(dr), alpha),
				blend(uint8(sg), uint8(dg), alpha),
				blend(uint8(sb), uint8(db), alpha),
				0xff,
			})
		}
	}
}

// BlendBgColor blends the Image (receiver) into the background color
// specified. This is more efficient than creating a background image and
// blending with Blend.
func (im *Image) BlendBgColor(c color.Color) {
	bgra := BGRAModel.Convert(c).(BGRA)
	im.For(func(x, y int) BGRA {
		c := im.At(x, y).(BGRA)
		alpha := float64(c.A) / 255.0
		return BGRA{
			B: blend(c.B, bgra.B, alpha),
			G: blend(c.G, bgra.G, alpha),
			R: blend(c.R, bgra.R, alpha),
			A: 0xff,
		}
	})
}

// Blend returns the blended alpha color for src and dest colors.
// This assumes that the destination has alpha = 1.
func BlendBGRA(src, dest BGRA) BGRA {
	alpha := float64(src.A) / 255.0
	return BGRA{
		B: blend(src.B, dest.B, alpha),
		G: blend(src.G, dest.G, alpha),
		R: blend(src.R, dest.R, alpha),
		A: 0xff,
	}
}

// blend calculates the value of a color given some alpha value in [0, 1]
// and a source and destination color. Note that this assumes that the
// destination is fully opaque (has an alpha value of 1).
func blend(s, d uint8, alpha float64) uint8 {
	return uint8(float64(s)*alpha + float64(d)*(1-alpha))
}

// SubImage provides a sub image of Image without copying image data.
// N.B. The standard library defines a similar function, but returns an
// image.Image. Here, we return xgraphics.Image so that we can use the extra
// methods defined by xgraphics on it.
func (im *Image) SubImage(r image.Rectangle) *Image {
	r = r.Intersect(im.Rect)
	if r.Empty() {
		return nil
	}

	i := im.PixOffset(r.Min.X, r.Min.Y)
	return &Image{
		X:      im.X,
		Pixmap: im.Pixmap,
		Pix:    im.Pix[i:],
		Stride: im.Stride,
		Rect:   r,
	}
}

// PixOffset returns the index of the frst element of the Pix data that
// corresponds to the pixel at (x, y).
func (im *Image) PixOffset(x, y int) int {
	return (y-im.Rect.Min.Y)*im.Stride + (x-im.Rect.Min.X)*4
}

// BGRA is the representation of color for each pixel in an X pixmap.
// BUG(burntsushi): This is hard-coded when it shouldn't be.
type BGRA struct {
	B, G, R, A uint8
}

// RGBA satisfies the color.Color interface.
func (c BGRA) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8

	g = uint32(c.G)
	g |= g << 8

	b = uint32(c.B)
	b |= b << 8

	a = uint32(c.A)
	a |= a << 8

	return
}

// bgraModel converts from any color to a BGRA color type.
func bgraModel(c color.Color) color.Color {
	if _, ok := c.(BGRA); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	return BGRA{
		B: uint8(b >> 8),
		G: uint8(g >> 8),
		R: uint8(r >> 8),
		A: uint8(a >> 8),
	}
}
