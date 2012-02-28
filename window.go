/*
    A few utility functions related to client windows. In 
    particular, getting an accurate geometry of a client window
    including the decorations (this can vary with the window
    manager). Also, a functon to move and/or resize a window
    accurately by the top-left corner. (Also can change based on
    the currently running window manager.) 

    This module also contains a function 'Listen' that must be used 
    in order to receive certain events from a window.

    {SHOW EXAMPLE}

    The idea here is to tell X that you want events that fall under
    the 'PropertyChange' category. Then you bind 'func' to the 
    particular event 'PropertyNotify'.
*/
package xgbutil

import "code.google.com/p/x-go-binding/xgb"

func (xu *XUtil) ParentWindow(win xgb.Id) xgb.Id {
    tree, err := xu.conn.QueryTree(win)

    if err != nil {
        panic(xerr(err, "ParentWindow",
                   "Error retrieving parent window for %x", win))
    }

    return tree.Parent
}

