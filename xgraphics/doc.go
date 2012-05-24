/*
Package xgraphics defines an X image type and provides convenience functions 
reading and writing X pixmaps and bitmaps. It is a work-in-progress, and while 
it works for some common configurations, it does not work for all 
configurations. Package xgraphics also provides some support for drawing text 
on to images using freetype-go, scaling images using graphics-go, simple alpha 
blending, finding EWMH and ICCCM window icons and efficiently drawing any image 
into an X pixmap. (Which "efficient" means being able to specify sub-regions of 
images to draw, so that the entire image isn't sent to X.)

In general, xgraphics paints pixmaps to windows using using the BackPixmap 
approach. (Setting the background pixmap of the window to the pixmap containing 
your image, and clearing the window's background when the pixmap is updated.) 
It also provides experimental support for another mechanism: copying the 
contents of your image's pixmap directly to the window. (This requires 
responding to expose events to redraw the pixmap.) The former approach requires 
less book-keeping, but supposedly has some issues with some video cards. The 
latter approach is probably more reliable, but requires more book-keeping.

A quick example

This is a simple example the converts any value satisfying the image.Image 
interface into an *xgraphics.Image value, and creates a new window with that 
image painted in the window. (The XShow function probably doesn't have any 
practical applications outside serving as an example, but can be useful for 
debugging what an image looks like.)

	imgFile, err := os.Open(imgPath)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Fatal(err)
	}

	ximg := xgraphics.NewConvert(XUtilValue, img)
	ximg.XShow()

A complete working example named 'show-image' that's similar to this can be 
found in the examples directory of the xgbutil package.

Portability

The xgraphics package *assumes* a particular kind of X server configuration. 
Namely, this configuration specifies bits per pixel, image byte order, bitmap 
bit order, scanline padding and unit length, image depth and so on. Handling 
all of the possible values for each configuration option will greatly inflate 
the code, but is on the TODO list.

I am undecided (perhaps because I haven't thought about it too much) about 
whether to hide these configuration details behind multiple xgraphics.Image 
types or hiding everything inside one xgraphics.Image type. I lean toward the 
latter because the former requires a large number of types (and therefore a lot 
of code duplication).

If your X server is not configured to what the xgraphics package expects, 
messages will be emitted to stderr when a new xgraphics.Image value is created. 
If you see any of these messages, please report them to xgbutil's project page:
https://github.com/BurntSushi/xgbutil.
*/
package xgraphics
