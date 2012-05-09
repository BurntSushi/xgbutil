package main

import (
	"fmt"
	"math/rand"
	// "os" 
	"github.com/BurntSushi/xgb"
	"time"
)

import (
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xwindow"
)

var X *xgbutil.XUtil
var Xerr error

func main() {
	X, Xerr = xgbutil.NewConn()
	if Xerr != nil {
		panic(Xerr)
	}

	fmt.Println(X)

	showDesk, _ := ewmh.ShowingDesktopGet(X)
	fmt.Printf("Showing desktop? %v\n", showDesk)

	wmName, err := ewmh.GetEwmhWM(X)
	if err != nil {
		fmt.Printf("No conforming window manager found... :-(\n")
		fmt.Println(err)
	} else {
		fmt.Printf("Window manager: %s\n", wmName)
	}

	pager := xgb.Id(0x160001e)
	middle := xgb.Id(0x3200016)
	geom, _ := ewmh.DesktopGeometryGet(X)
	active, _ := ewmh.ActiveWindowGet(X)
	desktops, _ := ewmh.DesktopNamesGet(X)
	curdesk, _ := ewmh.CurrentDesktopGet(X)
	clients, _ := ewmh.ClientListGet(X)
	activeName, _ := ewmh.WmNameGet(X, active)

	fmt.Printf("Active window: %x\n", active)
	fmt.Printf("Current desktop: %d\n", curdesk)
	fmt.Printf("Client list: %v\n", clients)
	fmt.Printf("Desktop geometry: (width: %d, height: %d)\n",
		geom.Width, geom.Height)
	fmt.Printf("Active window name: %s\n", activeName)
	fmt.Printf("Desktop names: %s\n", desktops)

	var desk string
	if curdesk < len(desktops) {
		desk = desktops[curdesk]
	} else {
		desk = string(curdesk)
	}
	fmt.Printf("Current desktop: %s\n", desk)

	// fmt.Printf("\nChanging current desktop to 25 from %d\n", curdesk) 
	ewmh.CurrentDesktopSet(X, curdesk)
	// fmt.Printf("Current desktop is now: %d\n", ewmh.CurrentDesktop(X)) 

	fmt.Printf("Setting active win to %x\n", middle)
	ewmh.ActiveWindowReq(X, middle)

	rand.Seed(int64(time.Now().Nanosecond()))
	randStr := make([]byte, 20)
	for i, _ := range randStr {
		if rf := rand.Float32(); rf < 0.40 {
			randStr[i] = byte('a' + rand.Intn('z'-'a'))
		} else if rf < 0.80 {
			randStr[i] = byte('A' + rand.Intn('Z'-'A'))
		} else {
			randStr[i] = ' '
		}
	}

	ewmh.WmNameSet(X, active, string(randStr))
	newName, _ := ewmh.WmNameGet(X, active)
	fmt.Printf("New name: %s\n", newName)

	// deskNames := ewmh.DesktopNamesGet(X) 
	// fmt.Printf("Desktop names: %s\n", deskNames) 
	// deskNames[len(deskNames) - 1] = "xgbutil" 
	// ewmh.DesktopNamesSet(X, deskNames) 
	// fmt.Printf("Desktop names: %s\n", ewmh.DesktopNamesGet(X)) 

	supported, _ := ewmh.SupportedGet(X)
	fmt.Printf("Supported hints: %v\n", supported)
	fmt.Printf("Setting supported hints...\n")
	ewmh.SupportedSet(X, []string{"_NET_CLIENT_LIST", "_NET_WM_NAME",
		"_NET_WM_DESKTOP"})

	numDesks, _ := ewmh.NumberOfDesktopsGet(X)
	fmt.Printf("Number of desktops: %d\n", numDesks)
	// ewmh.NumberOfDesktopsReq(X.EwmhNumberOfDesktops(X) + 1) 
	// time.Sleep(time.Second) 
	// fmt.Printf("Number of desktops: %d\n", ewmh.NumberOfDesktops(X)) 

	viewports, _ := ewmh.DesktopViewportGet(X)
	fmt.Printf("Viewports (%d): %v\n", len(viewports), viewports)

	// viewports[2].X = 50
	// viewports[2].Y = 293 
	// ewmh.DesktopViewportSet(X, viewports) 
	// time.Sleep(time.Second) 
	//  
	// viewports = ewmh.DesktopViewport(X) 
	// fmt.Printf("Viewports (%d): %v\n", len(viewports), viewports) 

	// ewmh.CurrentDesktopReq(X, 3) 

	visDesks, _ := ewmh.VisibleDesktopsGet(X)
	workarea, _ := ewmh.WorkareaGet(X)
	fmt.Printf("Visible desktops: %v\n", visDesks)
	fmt.Printf("Workareas: %v\n", workarea)
	// fmt.Printf("Virtual roots: %v\n", ewmh.VirtualRoots(X)) 
	// fmt.Printf("Desktop layout: %v\n", ewmh.DesktopLayout(X)) 
	fmt.Printf("Closing window %x\n", 0x2e004c5)
	ewmh.CloseWindow(X, 0x2e004c5)

	fmt.Printf("Moving/resizing window: %x\n", 0x2e004d0)
	ewmh.MoveresizeWindow(X, 0x2e004d0, 1920, 30, 500, 500)

	// fmt.Printf("Trying to initiate a moveresize...\n") 
	// ewmh.WmMoveresize(X, 0x2e004db, xgbutil.EwmhMove) 
	// time.Sleep(5 * time.Second) 
	// ewmh.WmMoveresize(X, 0x2e004db, xgbutil.EwmhCancel) 

	// fmt.Printf("Stacking window %x...\n", 0x2e00509) 
	// ewmh.RestackWindow(X, 0x2e00509) 

	fmt.Printf("Requesting frame extents for active window...\n")
	ewmh.RequestFrameExtents(X, active)

	parent, _ := xwindow.ParentWindow(X, active)
	actOpacity, _ := ewmh.WmWindowOpacityGet(X, parent)
	// actOpacity2 := ewmh.WmWindowOpacityGet(
	// X.ParentWindow(X.EwmhActiveWindow(X))) 
	fmt.Printf("Opacity for active window: %f\n", actOpacity)
	// fmt.Printf("Opacity for real active window: %f\n", actOpacity2) 
	// ewmh.WmWindowOpacitySet(X.ParentWindow(X, active), 0.5) 

	activeDesk, _ := ewmh.WmDesktopGet(X, active)
	activeType, _ := ewmh.WmWindowTypeGet(X, active)
	fmt.Printf("Active window's desktop: %d\n", activeDesk)
	fmt.Printf("Active's types: %v\n", activeType)
	// fmt.Printf("Pager's types: %v\n", ewmh.WmWindowType(X, 0x180001e)) 

	// fmt.Printf("Pager's state: %v\n", ewmh.WmState(X, 0x180001e)) 

	// ewmh.WmStateReq(X, active, xgbutil.EwmhStateToggle,
	// "_NET_WM_STATE_HIDDEN") 
	// ewmh.WmStateReqExtra(X, active, xgbutil.EwmhStateToggle, 
	// "_NET_WM_STATE_MAXIMIZED_VERT", 
	// "_NET_WM_STATE_MAXIMIZED_HORZ", 2) 

	activeAllowed, _ := ewmh.WmAllowedActionsGet(X, active)
	fmt.Printf("Allowed actions on active: %v\n", activeAllowed)

	struts, err := ewmh.WmStrutGet(X, pager)
	if err != nil {
		fmt.Printf("Pager struts: %v\n", err)
	} else {
		fmt.Printf("Pager struts: %v\n", struts)
	}

	pstruts, err := ewmh.WmStrutPartialGet(X, pager)
	if err != nil {
		fmt.Printf("Pager struts partial: %v - %v\n", pstruts, err)
	} else {
		fmt.Printf("Pager struts partial: %v\n", pstruts.BottomStartX)
	}

	// fmt.Printf("Icon geometry for active: %v\n",
	// ewmh.WmIconGeometry(X, active)) 

	icons, _ := ewmh.WmIconGet(X, active)
	fmt.Printf("Active window's (%x) icon data: (length: %v)\n",
		active, len(icons))
	for _, icon := range icons {
		fmt.Printf("\t(%d, %d)", icon.Width, icon.Height)
		fmt.Printf(" :: %d == %d\n", icon.Width*icon.Height, len(icon.Data))
	}
	// fmt.Printf("Now set them again...\n") 
	// ewmh.WmIconSet(X, active, icons[:len(icons) - 1]) 
}
