package main

import (
	"image"
	"log"
	"time"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
)

func main() {
	winw, winh := 128, 128
	time.Sleep(time.Nanosecond)

	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	active, _ := ewmh.ActiveWindowGet(X)
	icons, _ := ewmh.WmIconGet(X, active)
	icon := xgraphics.FindBestIcon(256, 256, icons)

	ximg, err := xgraphics.NewEwmhIcon(X, icon)
	if err != nil {
		log.Fatal(err)
	}

	ximg, err = ximg.Scale(winw, winh)
	if err != nil {
		log.Fatal(err)
	}

	ximg.BlendBgColor(xgraphics.BGRA{0x0, 0x0, 0xff, 0xff})

	win, err := xwindow.Generate(X)
	if err != nil {
		log.Fatal(err)
	}
	win.Create(X.RootWin(), 0, 0, winw, winh, xproto.CwBackPixel, 0xffffff)
	// win.Map() 
	xproto.MapWindowChecked(X.Conn(), win.Id).Check()
	ximg.XSurfaceSet(win.Id)
	ximg.XDraw()
	ximg.XPaint(win.Id)

	ximg.XShow()

	time.Sleep(time.Second)

	subimg := ximg.SubImage(image.Rect(20, 20, 50, 50))
	subimg.For(func(x, y int) xgraphics.BGRA {
		return xgraphics.BGRA{0xff, 0x0, 0x0, 0xff}
	})

	subimg.XDraw()
	subimg.XPaint(win.Id)

	ximg.XShow()

	subimg.For(func(x, y int) xgraphics.BGRA {
		return xgraphics.BGRA{0x00, 0xff, 0x0, 0xff}
	})
	subimg.XDraw()
	subimg.XPaint(win.Id)

	err = ximg.SavePng("a.png")
	if err != nil {
		log.Fatal(err)
	}

	xevent.Main(X)
}
