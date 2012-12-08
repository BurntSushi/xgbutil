package main

import (
	"log"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xgraphics"
)

func main() {
	log.SetFlags(0)

	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatalln(err)
	}
	checkCompatibility(X)
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
	failed := false

	if s.ImageByteOrder != xproto.ImageOrderLSBFirst {
		log.Printf("Your X server uses MSB image byte order. Unfortunately, " +
			"xgraphics currently requires LSB image byte order. You may see " +
			"weird things. Please report this.")
		failed = true
	}
	if s.BitmapFormatBitOrder != xproto.ImageOrderLSBFirst {
		log.Printf("Your X server uses MSB bitmap bit order. Unfortunately, " +
			"xgraphics currently requires LSB bitmap bit order. If you " +
			"aren't using X bitmaps, you should be able to proceed normally. " +
			"Please report this.")
		failed = true
	}
	if s.BitmapFormatScanlineUnit != 32 {
		log.Printf("xgraphics expects that the scanline unit is set to 32, "+
			"but your X server has it set to '%d'. "+
			"Namely, xgraphics hasn't been tested on other values. Things "+
			"may still work. Particularly, if you aren't using X bitmaps, "+
			"you should be completely unaffected. Please report this.",
			s.BitmapFormatScanlineUnit)
		failed = true
	}
	if scrn.RootDepth != 24 {
		log.Printf("xgraphics expects that the root window has a depth of 24, "+
			"but yours has depth '%d'. Its possible things will still work "+
			"if your value is 32, but will be unlikely to work with values "+
			"less than 24. Please report this.", scrn.RootDepth)
		failed = true
	}

	// Look for the default format for pixmaps and make sure bits per pixel
	// is 32.
	format := xgraphics.GetFormat(X, scrn.RootDepth)
	if format.BitsPerPixel != 32 {
		log.Printf("xgraphics expects that the bits per pixel for the root "+
			"window depth is 32. On your system, the root depth is %d and "+
			"the bits per pixel is %d. Things will most certainly not work. "+
			"Please report this.",
			scrn.RootDepth, format.BitsPerPixel)
		failed = true
	}

	// Give instructions on reporting the issue.
	if failed {
		log.Printf("Please report the aforementioned error message(s) at " +
			"https://github.com/BurntSushi/xgbutil. Please also include the " +
			"entire output of the `xdpyinfo` command in your report. Thanks!")
	} else {
		log.Printf("No compatibility issues detected.")
	}
}
