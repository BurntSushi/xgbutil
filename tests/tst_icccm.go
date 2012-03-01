package main

import "fmt"
import "code.google.com/p/x-go-binding/xgb"

import (
    "github.com/BurntSushi/xgbutil"
    "github.com/BurntSushi/xgbutil/ewmh"
    "github.com/BurntSushi/xgbutil/icccm"
)

func showTest(vals...interface{}) {
    fmt.Printf("%s\n\t%v\n\t%v\n", vals...)
}
func main() {
    X, _ := xgbutil.Dial("")

    active, err := ewmh.ActiveWindowGet(X)

    wmName, err := icccm.WmNameGet(X, active)
    showTest("WM_NAME get", wmName, err)

    err = icccm.WmNameSet(X, active, "hooblah")
    wmName, _ = icccm.WmNameGet(X, active)
    showTest("WM_NAME set", wmName, err)

    wmNormHints, err := icccm.WmNormalHintsGet(X, active)
    showTest("WM_NORMAL_HINTS get", wmNormHints, err)

    wmNormHints.Width += 5
    err = icccm.WmNormalHintsSet(X, active, wmNormHints)
    showTest("WM_NORMAL_HINTS set", wmNormHints, err)

    wmHints, err := icccm.WmHintsGet(X, active)
    showTest("WM_HINTS get", wmHints, err)

    wmHints.InitialState = icccm.StateNormal
    err = icccm.WmHintsSet(X, active, wmHints)
    showTest("WM_NORMAL_HINTS set", wmHints, err)

    wmClass, err := icccm.WmClassGet(X, active)
    showTest("WM_CLASS get", wmClass, err)

    wmClass.Instance = "hoopdy hoop"
    err = icccm.WmClassSet(X, active, wmClass)
    showTest("WM_CLASS set", wmClass, err)

    wmTrans, err := icccm.WmTransientForGet(X, active)
    showTest("WM_TRANSIENT_FOR get", wmTrans, err)

    wmProts, err := icccm.WmProtocolsGet(X, active)
    showTest("WM_PROTOCOLS get", wmProts, err)

    wmClient, err := icccm.WmClientMachineGet(X, active)
    showTest("WM_CLIENT_MACHINE get", wmClient, err)

    err = icccm.WmClientMachineSet(X, active, "Leopard")
    wmClient, _ = icccm.WmClientMachineGet(X, active)
    showTest("WM_CLIENT_MACHINE set", wmClient, err)

    wmState, err := icccm.WmStateGet(X, active)
    showTest("WM_STATE get", wmState, err)

    wmState.Icon = xgb.Id(8365538)
    wmState.State = icccm.StateNormal
    err = icccm.WmStateSet(X, active, wmState)
    wmState, _ = icccm.WmStateGet(X, active)
    showTest("WM_STATE set", wmState, err)
}

