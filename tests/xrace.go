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
	X, _ := xgbutil.NewConn()
	conn := X.Conn()

	aDesktop := "_NET_WM_DESKTOP"
	aActive := "_NET_ACTIVE_WINDOW"

	xwindow.Listen(X, X.RootWin(), xgb.EventMaskPropertyChange)
	xevent.PropertyNotifyFun(
		func(X *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
			for i := 0; i < 1; i++ {
				log.Println("PROPERTY CHANGE")
			}
		}).Connect(X, X.RootWin())

	go func() {
		for {
			reply, err := conn.InternAtom(true, uint16(len(aDesktop)),
				aDesktop).Reply()
			if err != nil {
				log.Fatal(err)
			}

			log.Println("A1-299", reply.Sequence, reply.Atom)
			time.Sleep(sleepy)
		}
	}()

	go func() {
		for {
			reply, err := conn.InternAtom(true, uint16(len(aActive)),
				aActive).Reply()
			if err != nil {
				log.Fatal(err)
			}

			log.Println("A2-294", reply.Sequence, reply.Atom)
			time.Sleep(sleepy)
		}
	}()

	go func() {
		for {
			reply, err := conn.GetGeometry(0x1).Reply()
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
			reply, err := conn.GetGeometry(0x2).Reply()
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
