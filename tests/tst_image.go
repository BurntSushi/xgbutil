package main

import (
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
)

func showIcon(X *xgbutil.XUtil, wid xproto.Window, name string) {
	hints, err := icccm.WmHintsGet(X, wid)
	if err != nil {
		log.Fatal(err)
	}

	ximg, err := xgraphics.NewIcccmIcon(X, hints.IconPixmap, hints.IconMask)
	if err != nil {
		log.Fatal(err)
	}

	err = ximg.SavePng(fmt.Sprintf("%s.png", name))
	if err != nil {
		log.Fatal(err)
	}

	ximg.XShow()
}

func main() {
	// winw, winh := 128, 128 
	time.Sleep(time.Nanosecond)

	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	libre := xproto.Window(0x4a0001d)
	xclock := xproto.Window(0x480000a)

	showIcon(X, libre, "libre")
	showIcon(X, xclock, "xclock")

	xevent.Main(X)
}
