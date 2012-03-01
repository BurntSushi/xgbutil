package main

import "fmt"
import "code.google.com/p/x-go-binding/xgb"
import "github.com/BurntSushi/xgbutil"

func showTest(vals...interface{}) {
    fmt.Printf("%s\n\t%v\n\t%v\n", vals...)
}
func main() {
    X, _ := xgbutil.Dial("")

    active, err := X.EwmhActiveWindow()

    wmName, err := X.IcccmWmName(active)
    showTest("WM_NAME get", wmName, err)

    err = X.IcccmWmNameSet(active, "hooblah")
    wmName, _ = X.IcccmWmName(active)
    showTest("WM_NAME set", wmName, err)

    wmNormHints, err := X.IcccmWmNormalHints(active)
    showTest("WM_NORMAL_HINTS get", wmNormHints, err)

    wmNormHints.Width += 5
    err = X.IcccmWmNormalHintsSet(active, wmNormHints)
    showTest("WM_NORMAL_HINTS set", wmNormHints, err)

    wmHints, err := X.IcccmWmHints(active)
    showTest("WM_HINTS get", wmHints, err)

    wmHints.InitialState = xgbutil.StateNormal
    err = X.IcccmWmHintsSet(active, wmHints)
    showTest("WM_NORMAL_HINTS set", wmHints, err)

    wmClass, err := X.IcccmWmClass(active)
    showTest("WM_CLASS get", wmClass, err)

    wmClass.Instance = "kOnSoLe!!?!"
    err = X.IcccmWmClassSet(active, wmClass)
    showTest("WM_CLASS set", wmClass, err)

    wmTrans, err := X.IcccmWmTransientFor(active)
    showTest("WM_TRANSIENT_FOR get", wmTrans, err)

    wmProts, err := X.IcccmWmProtocols(active)
    showTest("WM_PROTOCOLS get", wmProts, err)

    wmClient, err := X.IcccmWmClientMachine(active)
    showTest("WM_CLIENT_MACHINE get", wmClient, err)

    err = X.IcccmWmClientMachineSet(active, "Leopard")
    wmClient, _ = X.IcccmWmClientMachine(active)
    showTest("WM_CLIENT_MACHINE set", wmClient, err)

    wmState, err := X.IcccmWmState(active)
    showTest("WM_STATE get", wmState, err)

    wmState.Icon = xgb.Id(8365538)
    wmState.State = xgbutil.StateNormal
    err = X.IcccmWmStateSet(active, wmState)
    wmState, _ = X.IcccmWmState(active)
    showTest("WM_STATE set", wmState, err)
}

