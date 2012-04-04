/*
    A file with a couple of functions related to the Xinerama extension.
    Namely, just grab a list of rectangles representing all active heads.

    The cool thing about Xinerama is that both RandR and TwinView provide
    Xinerama data even when the Xinerama extension isn't the thing powering
    your displays. Considering this is the only extension I have working with
    XGB, this is good news.
*/
package xinerama

import "sort"

import "burntsushi.net/go/xgbutil"
import "burntsushi.net/go/xgbutil/xrect"

// Alias so we use it as a receiver to satisfy sort.Interface
type Heads []xrect.Rect

// Len satisfies 'Len' in sort.Interface.
func (hds Heads) Len() int {
    return len(hds)
}

// Less satisfies 'Less' in sort.Interface.
func (hds Heads) Less(i int, j int) bool {
    return hds[i].X() < hds[j].X() || (hds[i].X() == hds[j].X() &&
                                       hds[i].Y() < hds[j].Y())
}

// Swap does just that. Nothing to see here...
func (hds Heads) Swap(i int, j int) {
    hds[i], hds[j] = hds[j], hds[i]
}

// Heads returns the list of heads in a physical ordering.
// Namely, left to right then top to bottom. (Defined by (X, Y).)
func PhysicalHeads(xu *xgbutil.XUtil) (Heads, error) {
    xinfo, err := xu.Conn().XineramaQueryScreens()
    if err != nil {
        return nil, err
    }

    hds := make(Heads, 0)
    for _, info := range xinfo.ScreenInfo {
        head := xrect.Make(int(info.XOrg), int(info.YOrg),
                           int(info.Width), int(info.Height))

        // Maybe Xinerama is enabled, but we have cloned displays...
        unique := true
        for _, h := range hds {
            if h.X() == head.X() && h.Y() == head.Y() {
                unique = false
                break
            }
        }

        if unique {
            hds = append(hds, head)
        }
    }

    sort.Sort(hds)
    return hds, nil
}

