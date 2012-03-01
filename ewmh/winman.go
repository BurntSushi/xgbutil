/*
    A convenience function to find the name of an active EWMH compliant
    window manager.
*/
package ewmh

// GetEwmhWM uses the EWMH spec to find if a conforming window manager
// is currently running or not. If it is, then its name will be returned.
// Otherwise, an error will be returned explaining why one couldn't be found.
// (This function is safe.)
func (xu *XUtil) GetEwmhWM() (wmName string, err error) {
    childCheck, err := xu.EwmhSupportingWmCheck(xu.root)
    if err != nil {
        return "", xuerr("GetEwmhWM", "Failed because: %v", err)
    }

    childCheck2, err := xu.EwmhSupportingWmCheck(childCheck)
    if err != nil {
        return "", xuerr("GetEwmhWM", "Failed because: %v", err)
    }

    if childCheck != childCheck2 {
        return "", xuerr("GetEwmhWM",
                         "_NET_SUPPORTING_WM_CHECK value on the root window " +
                         "(%x) does not match _NET_SUPPORTING_WM_CHECK value " +
                         "on the child window (%x).", childCheck, childCheck2)
    }

    return xu.EwmhWmName(childCheck)
}

