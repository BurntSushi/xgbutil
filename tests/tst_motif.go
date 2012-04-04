package main

import "fmt"

import "burntsushi.net/go/x-go-binding/xgb"
import "burntsushi.net/go/xgbutil"
import "burntsushi.net/go/xgbutil/ewmh"
import "burntsushi.net/go/xgbutil/motif"

func DoDecor(mh motif.Hints) bool {
    if mh.Flags & motif.HintDecorations > 0 &&
        (mh.Decoration == motif.DecorationNone ||
         (mh.Decoration & motif.DecorationAll == 0 &&
          mh.Decoration & motif.DecorationTitle == 0 &&
          mh.Decoration & motif.DecorationResizeH == 0)) {
        return false
    }

    return true
}

func main() {
    X, _ := xgbutil.Dial("")

    gChrome := xgb.Id(0x2e00047)
    active, _ := ewmh.ActiveWindowGet(X)

    mh, err := motif.WmHintsGet(X, gChrome)
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(mh)
        fmt.Println("Does Chrome want decorations?", DoDecor(mh))
    }

    mh, err = motif.WmHintsGet(X, active)
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println(mh)
        fmt.Println("Does Active window want decorations?", DoDecor(mh))
    }
}

