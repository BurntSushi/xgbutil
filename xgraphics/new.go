package xgraphics

/*
xgraphics/new.go contains a few additional constructors for creating an
xgraphics.Image.
*/

import (
	"fmt"
	"image"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xwindow"
)

// NewConvert converts any image satisfying the image.Image interface to an
// xgraphics.Image type.
// NewConvert does not check if 'img' is an xgraphics.Image. Thus, NewConvert
// provides a convenient mechanism for copying xgraphic.Image values.
// Also, NewConvert attempts to optimize image conversion for some image
// formats. (Currently, only *image.YCbCr and *image.RGBA.)
func NewConvert(X *xgbutil.XUtil, img image.Image) *Image {
	ximg := New(X, img.Bounds())

	// I've attempted to optimize this loop.
	// It actually takes more time to convert an image than to send the bytes
	// over the wire. (I suspect 'copy' is super fast, which can be used in
	// XDraw, whereas computing each pixel is super slow.)
	// But how is image decoding so much faster than this? I'll have to
	// investigate... Maybe the Color interface being used here is the real
	// slow down.
	switch concrete := img.(type) {
	case *image.RGBA:
		convertRGBA(ximg, concrete)
	case *image.YCbCr:
		convertYCbCr(ximg, concrete)
	default:
		xgbutil.Logger.Printf("Converting image type %T the slow way. "+
			"Optimization for this image type hasn't been added yet.", img)
		convertImage(ximg, img)
	}
	return ximg
}

// NewEwmhIcon converts EWMH icon data (ARGB) to an xgraphics.Image type.
// You should probably use xgraphics.FindIcon instead of this directly.
func NewEwmhIcon(X *xgbutil.XUtil, icon *ewmh.WmIcon) *Image {
	ximg := New(X, image.Rect(0, 0, icon.Width, icon.Height))
	r := ximg.Rect
	width := r.Dx()

	var argb, x, y int
	for x = r.Min.X; x < r.Max.X; x++ {
		for y = r.Min.Y; y < r.Max.Y; y++ {
			argb = icon.Data[x+(y*width)]
			ximg.SetBGRA(x, y, BGRA{
				B: uint8(argb & 0x000000ff),
				G: uint8((argb & 0x0000ff00) >> 8),
				R: uint8((argb & 0x00ff0000) >> 16),
				A: uint8(argb >> 24),
			})
		}
	}
	return ximg
}

// NewIcccmIcon converts two pixmap ids (icon_pixmap and icon_mask in the
// WM_HINTS properts) to a single xgraphics.Image.
// It is okay for one of iconPixmap or iconMask to be 0, but not both.
// You should probably use xgraphics.FindIcon instead of this directly.
func NewIcccmIcon(X *xgbutil.XUtil, iconPixmap,
	iconMask xproto.Pixmap) (*Image, error) {

	if iconPixmap == 0 && iconMask == 0 {
		return nil, fmt.Errorf("NewIcccmIcon: At least one of iconPixmap or " +
			"iconMask must be non-zero, but both are 0.")
	}

	var pximg, mximg *Image
	var err error

	// Get the xgraphics.Image for iconPixmap.
	if iconPixmap != 0 {
		pximg, err = NewPixmap(X, iconPixmap)
		if err != nil {
			return nil, err
		}
	}

	// Now get the xgraphics.Image for iconMask.
	if iconMask != 0 {
		mximg, err = NewPixmap(X, iconMask)
		if err != nil {
			return nil, err
		}
	}

	// Now merge them together if both were specified.
	switch {
	case pximg != nil && mximg != nil:
		r := pximg.Bounds()

		var x, y int
		var bgra, maskBgra BGRA
		for x = r.Min.X; x < r.Max.X; x++ {
			for y = r.Min.Y; y < r.Max.Y; y++ {
				maskBgra = mximg.At(x, y).(BGRA)
				bgra = pximg.At(x, y).(BGRA)
				if maskBgra.A == 0 {
					pximg.SetBGRA(x, y, BGRA{
						B: bgra.B,
						G: bgra.G,
						R: bgra.R,
						A: 0,
					})
				}
			}
		}
		return pximg, nil
	case pximg != nil:
		return pximg, nil
	case mximg != nil:
		return mximg, nil
	}

	panic("unreachable")
}

// NewPixmap converts an X pixmap into a xgraphics.Image.
// This is used in NewIcccmIcon.
func NewPixmap(X *xgbutil.XUtil, iconPixmap xproto.Pixmap) (*Image, error) {
	// Get the geometry of the pixmap for use in the GetImage request.
	pgeom, err := xwindow.RawGeometry(X, xproto.Drawable(iconPixmap))
	if err != nil {
		return nil, err
	}

	// Get the image data for each pixmap.
	pixmapData, err := xproto.GetImage(X.Conn(), xproto.ImageFormatZPixmap,
		xproto.Drawable(iconPixmap),
		0, 0, uint16(pgeom.Width()), uint16(pgeom.Height()),
		(1<<32)-1).Reply()
	if err != nil {
		return nil, err
	}

	// Now create the xgraphics.Image and populate it with data from
	// pixmapData and maskData.
	ximg := New(X, image.Rect(0, 0, pgeom.Width(), pgeom.Height()))

	// We'll try to be a little flexible with the image format returned,
	// but not completely flexible.
	err = readPixmapData(X, ximg, iconPixmap, pixmapData,
		pgeom.Width(), pgeom.Height())
	if err != nil {
		return nil, err
	}

	return ximg, nil
}

// readPixmapData uses Format information to read data from an X pixmap
// into an xgraphics.Image.
// readPixmapData does not take into account all information possible to read
// X pixmaps and bitmaps. Of particular note is bit order/byte order.
func readPixmapData(X *xgbutil.XUtil, ximg *Image, pixid xproto.Pixmap,
	imgData *xproto.GetImageReply, width, height int) error {

	format := getFormat(X, imgData.Depth)
	if format == nil {
		return fmt.Errorf("Could not find valid format for pixmap %d "+
			"with depth %d", pixid, imgData.Depth)
	}

	switch format.Depth {
	case 1: // We read bitmaps in as alpha masks.
		if format.BitsPerPixel != 1 {
			return fmt.Errorf("The image returned for pixmap id %d with "+
				"depth %d has an unsupported value for bits-per-pixel: %d",
				pixid, format.Depth, format.BitsPerPixel)
		}

		// Calculate the padded width of our image data.
		pad := int(X.Setup().BitmapFormatScanlinePad)
		paddedWidth := width
		if width%pad != 0 {
			paddedWidth = width + pad - (width % pad)
		}

		// Process one scanline at a time. Each 'y' represents a
		// single scanline.
		for y := 0; y < height; y++ {
			// Each scanline has length 'width' padded to
			// BitmapFormatScanlinePad, which is found in the X setup info.
			// 'i' is the index to the starting byte of the yth scanline.
			i := y * paddedWidth / 8
			for x := 0; x < width; x++ {
				b := imgData.Data[i+x/8] >> uint(x%8)
				if b&1 > 0 { // opaque
					ximg.Set(x, y, BGRA{0x0, 0x0, 0x0, 0xff})
				} else { // transparent
					ximg.Set(x, y, BGRA{0xff, 0xff, 0xff, 0x0})
				}
			}
		}
	case 24:
		if format.BitsPerPixel != 24 && format.BitsPerPixel != 32 {
			return fmt.Errorf("The image returned for pixmap id %d has "+
				"an unsupported value for bits-per-pixel: %d",
				pixid, format.BitsPerPixel)
		}
		bytesPer := int(format.BitsPerPixel) / 8
		var i int
		ximg.For(func(x, y int) BGRA {
			i = y*width*bytesPer + x*bytesPer
			return BGRA{
				B: imgData.Data[i],
				G: imgData.Data[i+1],
				R: imgData.Data[i+2],
				A: 0xff,
			}
		})
	default:
		return fmt.Errorf("The image returned for pixmap id %d has an "+
			"unsupported value for depth: %d", pixid, format.Depth)
	}

	return nil
}

// getFormat searches SetupInfo for a Format matching the depth provided.
func getFormat(X *xgbutil.XUtil, depth byte) *xproto.Format {
	for _, pixForm := range X.Setup().PixmapFormats {
		if pixForm.Depth == depth {
			return &pixForm
		}
	}
	return nil
}

// getVisualInfo searches SetupInfo for a VisualInfo value matching
// the depth provided.
// XXX: This isn't used (yet).
func getVisualInfo(X *xgbutil.XUtil, depth byte,
	visualid xproto.Visualid) *xproto.VisualInfo {

	for _, depthInfo := range X.Screen().AllowedDepths {
		fmt.Printf("%#v\n", depthInfo)
		// fmt.Printf("%#v\n", depthInfo.Visuals) 
		fmt.Println("------------")
		if depthInfo.Depth == depth {
			for _, visual := range depthInfo.Visuals {
				if visual.VisualId == visualid {
					return &visual
				}
			}
		}
	}
	return nil
}

// checkCompatibility reads info in the X setup info struct and emits
// messages to stderr if they don't correspond to values that xgraphics
// supports.
// The idea is that in the future, we'll support more values.
// The real reason for checkCompatibility is to make debugging easier. Without
// it, if the values weren't what we'd expect, we'd see garbled images in the
// best case, and probably BadLength errors in the worst case.
func checkCompatibility(X *xgbutil.XUtil) {
	s := X.Setup()
	scrn := X.Screen()
	lg := xgbutil.Logger
	if s.ImageByteOrder != xproto.ImageOrderLSBFirst {
		lg.Printf("Your X server uses MSB image byte order. Unfortunately, " +
			"xgraphics currently requires LSB image byte order. You may see " +
			"weird things. Please report this.")
	}
	if s.BitmapFormatBitOrder != xproto.ImageOrderLSBFirst {
		lg.Printf("Your X server uses MSB bitmap bit order. Unfortunately, " +
			"xgraphics currently requires LSB bitmap bit order. If you " +
			"aren't using X bitmaps, you should be able to proceed normally. " +
			"Please report this.")
	}
	if s.BitmapFormatScanlineUnit != 32 {
		lg.Printf("xgraphics expects that the scanline unit is set to 32, but "+
			"your X server has it set to '%d'. "+
			"Namely, xgraphics hasn't been tested on other values. Things "+
			"may still work. Particularly, if you aren't using X bitmaps, "+
			"you should be completely unaffected. Please report this.",
			s.BitmapFormatScanlineUnit)
	}
	if scrn.RootDepth != 24 {
		lg.Printf("xgraphics expects that the root window has a depth of 24, "+
			"but yours has depth '%d'. Its possible things will still work "+
			"if your value is 32, but will be unlikely to work with values "+
			"less than 24. Please report this.", scrn.RootDepth)
	}

	// Look for the default format for pixmaps and make sure bits per pixel
	// is 32.
	format := getFormat(X, scrn.RootDepth)
	if format.BitsPerPixel != 32 {
		lg.Printf("xgraphics expects that the bits per pixel for the root "+
			"window depth is 32. On your system, the root depth is %d and "+
			"the bits per pixel is %d. Things will most certainly not work. "+
			"Please report this.",
			scrn.RootDepth, format.BitsPerPixel)
	}
}
