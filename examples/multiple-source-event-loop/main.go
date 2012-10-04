// Example multiple-source-event-loop shows how to use the xevent package to
// combine multiple sources in your main event loop. This is particularly
// useful if your application can respond to user input from other sources.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/xgb/xproto"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/BurntSushi/xgbutil/xwindow"
)

// otherSource serves as a placeholder from some other source of user input.
func otherSource() chan int {
	c := make(chan int, 0)
	go func() {
		defer close(c)

		i := 1
		for {
			c <- i
			i++
			time.Sleep(time.Second)
		}
	}()
	return c
}

// sendClientMessages is a goroutine that sends client messages to the root
// window. We then listen to them later as a demonstration of responding to
// X events. (They are sent with SubstructureNotify and SubstructureRedirect
// masks set. So in order to receive them, we'll have to explicitly listen
// to events of that type on the root window.)
func xSource(X *xgbutil.XUtil) {
	i := 1
	for {
		ewmh.ClientEvent(X, X.RootWin(), "NOOP", i)
		i++
		time.Sleep(200 * time.Millisecond)
	}
}

func main() {
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	// Start generating other source events.
	otherChan := otherSource()

	// Start generating X events (by sending client messages to root window).
	go xSource(X)

	// Listen to those X events.
	xwindow.New(X, X.RootWin()).Listen(xproto.EventMaskSubstructureNotify)

	// Respond to those X events.
	xevent.ClientMessageFun(
		func(X *xgbutil.XUtil, ev xevent.ClientMessageEvent) {
			atmName, err := xprop.AtomName(X, ev.Type)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("ClientMessage: %d. %s\n", ev.Data.Data32[0], atmName)
		}).Connect(X, X.RootWin())

	// Instead of using the usual xevent.Main, we use xevent.MainPing.
	// It runs the main event loop inside a goroutine and returns ping
	// channels, which are sent benign values right before an event is
	// dequeued and right after that event has finished running all callbacks
	// associated with it, respectively.
	pingBefore, pingAfter, pingQuit := xevent.MainPing(X)
	for {
		select {
		case <-pingBefore:
			// Wait for the event to finish processing.
			<-pingAfter
		case otherVal := <-otherChan:
			fmt.Printf("Processing other event: %d\n", otherVal)
		case <-pingQuit:
			fmt.Printf("xevent loop has quit")
			return
		}
	}
}
