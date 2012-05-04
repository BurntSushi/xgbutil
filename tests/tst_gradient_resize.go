package main

import (
	"image"
	"image/color"
	"log"

	"github.com/BurntSushi/xgb"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
)

var X *xgbutil.XUtil

func createWindow() xgb.Id {
	wid, err := X.Conn().NewId()
	if err != nil {
		log.Fatal(err)
	}
	scrn := X.Screen()

	X.Conn().CreateWindow(scrn.RootDepth, wid, X.RootWin(), 0, 0, 400, 400, 0,
		xgb.WindowClassInputOutput, scrn.RootVisual, 0, []uint32{})

	return wid
}

func gradient(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	start, end := 0, 255
	inc := float64(end-start) / float64(width)

	for x := 0; x < width; x++ {
		clr := uint8(float64(start) + inc*float64(x))
		for y := 0; y < height; y++ {
			img.SetRGBA(x, y, color.RGBA{clr, clr, clr, 255})
		}
	}
	return img
}

func newGradientWindow(width, height int) {
	win := createWindow()
	X.Conn().ConfigureWindow(
		win, xgb.ConfigWindowWidth|xgb.ConfigWindowHeight,
		[]uint32{uint32(width), uint32(height)})
	xwindow.Listen(X, win, xgb.EventMaskStructureNotify)

	X.Conn().MapWindow(win)

	xgraphics.PaintImg(X, win, gradient(width, height))

	xevent.ConfigureNotifyFun(
		func(X *xgbutil.XUtil, ev xevent.ConfigureNotifyEvent) {
			img := gradient(int(ev.Width), int(ev.Height))
			log.Printf("Painting new image (%d, %d)", ev.Width, ev.Height)
			xgraphics.PaintImg(X, win, img)
		}).Connect(X, win)
}

func main() {
	X, _ = xgbutil.Dial("")

	go newGradientWindow(200, 200)
	go newGradientWindow(400, 400)

	xevent.Main(X)
}
