package main

import (
	"log"
	"time"

	"github.com/BurntSushi/xgb"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
)

func main() {
	sleepy := time.Millisecond
	X, _ := xgbutil.Dial("")
	conn := X.Conn()

	xwindow.Listen(X, X.RootWin(), xgb.EventMaskPropertyChange)
	xevent.PropertyNotifyFun(
		func(X *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
			for i := 0; i < 1; i++ {
				log.Println("PROPERTY CHANGE")
			}
		}).Connect(X, X.RootWin())

	go func() {
		for {
			reply, err := conn.InternAtom(true, "_NET_WM_DESKTOP")
			if err != nil {
				log.Fatal(err)
			}

			log.Println(reply)
			time.Sleep(sleepy)
		}
	}()

	go func() {
		for {
			reply, err := conn.InternAtom(true, "_NET_ACTIVE_WINDOW")
			if err != nil {
				log.Fatal(err)
			}

			log.Println(reply)
			time.Sleep(sleepy)
		}
	}()

	go func() {
		for {
			reply, err := conn.GetGeometry(0x1)
			if err != nil {
				log.Println("0x1:", err)
			} else {
				log.Println("0x1:", reply)
			}
			time.Sleep(sleepy)
		}
	}()

	go func() {
		for {
			reply, err := conn.GetGeometry(0x2)
			if err != nil {
				log.Println("0x2:", err)
			} else {
				log.Println("0x2:", reply)
			}
			time.Sleep(sleepy)
		}
	}()

	go xevent.Main(X)
	select {}
}
