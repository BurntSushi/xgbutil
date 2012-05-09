package main

import "fmt"
import "github.com/BurntSushi/xgbutil"
import "github.com/BurntSushi/xgbutil/ewmh"
import "github.com/BurntSushi/xgbutil/xinerama"
import "github.com/BurntSushi/xgbutil/xrect"
import "github.com/BurntSushi/xgbutil/xwindow"

func main() {
	X, _ := xgbutil.NewConn()

	heads, err := xinerama.PhysicalHeads(X)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

	// Test intersection
	r1 := xrect.New(0, 0, 100, 100)
	r2 := xrect.New(100, 100, 100, 100)
	fmt.Println(xrect.IntersectArea(r1, r2))

	// Test largest overlap
	window := xrect.New(1800, 0, 300, 200)
	fmt.Println(xrect.LargestOverlap(window, heads))

	// Test ApplyStrut
	rgeom, _ := xwindow.RawGeometry(X, X.RootWin())
	fmt.Println("---------------------------")
	for i, head := range heads {
		fmt.Printf("%d - %v\n", i, head)
	}

	// Let's actually apply struts to the current environment
	clients, _ := ewmh.ClientListGet(X)
	for _, client := range clients {
		strut, err := ewmh.WmStrutPartialGet(X, client)
		if err == nil {
			xrect.ApplyStrut(heads, rgeom.Width(), rgeom.Height(),
				strut.Left, strut.Right, strut.Top, strut.Bottom,
				strut.LeftStartY, strut.LeftEndY,
				strut.RightStartY, strut.RightEndY,
				strut.TopStartX, strut.TopEndX,
				strut.BottomStartX, strut.BottomEndX)
		}
	}

	fmt.Println("---------------------------")
	fmt.Println("After applying struts...")
	for i, head := range heads {
		fmt.Printf("%d - %v\n", i, head)
	}
}
