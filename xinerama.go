/*
    A file with a couple of functions related to the Xinerama extension.
    Namely, just grab a list of rectangles representing all active heads.

    The cool thing about Xinerama is that both RandR and TwinView provide
    Xinerama data even when the Xinerama extension isn't the thing powering
    your displays. Considering this is the only extension I have working with
    XGB, this is good news.
*/
package xgbutil

import "sort"

// Head is a struct representing an X head rectangle
// (the top left corner is the origin).
type Head struct {
    X, Y, Width, Height uint32
}

// Alias so we use it as a receiver to satisfy sort.Interface
type Heads []Head

// Len satisfies 'Len' in sort.Interface.
func (hds Heads) Len() int {
    return len(hds)
}

// Less satisfies 'Less' in sort.Interface.
func (hds Heads) Less(i int, j int) bool {
    return hds[i].X < hds[j].X || (
            hds[i].X == hds[j].X && hds[i].Y < hds[j].Y)
}

// Swap does just that. Nothing to see here...
func (hds Heads) Swap(i int, j int) {
    hds[i], hds[j] = hds[j], hds[i]
}

// Heads returns the list of heads in a physical ordering.
// Namely, left to right then top to bottom. (Defined by (X, Y).)
func (xu *XUtil) Heads() (Heads, error) {
    xinfo, err := xu.conn.XineramaQueryScreens()
    if err != nil {
        return nil, err
    }

    hds := make(Heads, len(xinfo.ScreenInfo))
    for i, info := range xinfo.ScreenInfo {
        hds[i] = Head{
            X: uint32(info.XOrg),
            Y: uint32(info.YOrg),
            Width: uint32(info.Width),
            Height: uint32(info.Height),
        }
    }

    sort.Sort(hds)
    return hds, nil
}

