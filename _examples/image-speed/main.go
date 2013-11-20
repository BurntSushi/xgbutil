//This demonstration shows drawing entire generated images to an x window
//and calculating the speed of the generation and the drawing operations
//It should be noted that redrawing the entire image is inefficient, this demo
//was made to show me how fast I could draw to the windows with this method.
package main

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
	"image"
	"log"
	"time"
)

func main() {
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Println(err)
		return
	}

	//Initialize the keybind package
	keybind.Initialize(X)

	//Create a window
	win, err := xwindow.Generate(X)
	if err != nil {
		log.Fatalf("Could not generate a new window X id: %s", err)
	}
	win.Create(X.RootWin(), 0, 0, 1024, 768, xproto.CwBackPixel, 0x606060ff)

	// Listen for Key{Press,Release} events.
	win.Listen(xproto.EventMaskKeyPress, xproto.EventMaskKeyRelease)

	// Map the window. This is what makes it on the screen
	win.Map()

	//Make a ...callback... for the events and connect
	xevent.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			modStr := keybind.ModifierString(e.State)
			keyStr := keybind.LookupString(X, e.State, e.Detail)
			if len(modStr) > 0 {
				log.Printf("Key: %s-%s\n", modStr, keyStr)
			} else {
				log.Println("Key:", keyStr)
			}

			if keybind.KeyMatch(X, "Escape", e.State, e.Detail) {
				if e.State&xproto.ModMaskControl > 0 {
					log.Println("Control-Escape detected. Quitting...")
					xevent.Quit(X)
				}
			}
		}).Connect(X, win.Id)

	//So here i'm going to try to make a image..'
	img := xgraphics.New(X, image.Rect(0, 0, 1024, 768))

	err = img.XSurfaceSet(win.Id)
	if err != nil {
		log.Printf("Error while setting window surface to image %d: %s\n",
			win, err)
	} else {
		log.Printf("Window %d surface set to image OK\n", win)
	}

	// I /think/ XDraw actually sends data to server?
	img.XDraw()
	// I /think/ XPaint tells the server to paint image to window
	img.XPaint(win.Id)

	//Start the routine that updates the window
	go updater(img, win)

	//This seems to start a main loop for listening to xevents
	xevent.Main(X)
}

func updater(img *xgraphics.Image, win *xwindow.Window) {
	//We keep track of times based on 1024 frames
	frame := 0
	start := time.Now()
	var genStart, drawStart time.Time
	var genTotal, drawTotal time.Duration

	for {
		frame = frame + 1
		if frame > 1024 {
			frame = 0
			log.Printf("Time elapsed: %s\n", time.Now().Sub(start))
			log.Printf("Time generate: %s\n", genTotal)
			log.Printf("Time drawing: %s\n", drawTotal)
			start = time.Now()
			drawTotal, genTotal = 0, 0
		}

		genStart = time.Now()
		var x, y, i int
		for y = 0; y < 768; y++ {
			//set last pixel back to black
			img.SetBGRA(frame-1, y, xgraphics.BGRA{0, 0, 0, 255})
			for i = 0; i < 10; i++ {
				x = i + frame
				if x > 1024 {
					x = 1024
				}
				img.SetBGRA(x, y, xgraphics.BGRA{0, 0, 255, 255})
			}
		}
		genTotal += time.Now().Sub(genStart)

		drawStart = time.Now()
		//hopefully using checked will block us from drawing again before x
		//draws although XDraw might block anyway, we can check for an error
		//here
		err := img.XDrawChecked()
		if err != nil {
			log.Println(err)
			return
		}
		//img.XDraw()

		img.XPaint(win.Id)
		drawTotal += time.Now().Sub(drawStart)
	}
}
