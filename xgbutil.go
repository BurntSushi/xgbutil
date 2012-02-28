/*
    This is a utility library designed to work with the X Go Binding. This 
    project's main goal is to make various X related tasks easier. For example, 
    binding keys, using the EWMH or ICCCM specs with the window manager, 
    moving/resizing windows, assigning function callbacks to particular events, 
    etc.
*/
package xgbutil

import (
    "fmt"

    "code.google.com/p/x-go-binding/xgb"
)

// An XUtil represents the state of xgbutil. It keeps track of the current 
// X connection, the root window, event callbacks, key/mouse bindings, etc.
type XUtil struct {
    conn *xgb.Conn
    root xgb.Id
    atoms map[string]xgb.Id
    atomNames map[xgb.Id]string
}

type XError struct {
    funcName string // some identifier so we know where the error comes from
    err string // free form string explaining the error
    XGBError *xgb.Error // error struct from XGB - to get the raw X error
}

func (xe *XError) Error() string {
    return fmt.Sprintf("%s: %s", xe.funcName, xe.err)
}

// Constructs an error struct from an X error
func xerr (xgberr interface{}, funcName string, err string,
           params ...interface{}) *XError {
    switch e := xgberr.(type) {
    case *xgb.Error:
        return &XError{
            funcName: funcName,
            err: fmt.Sprintf("%s: %v", fmt.Sprintf(err, params...), e),
            XGBError: e,
        }
    }

    panic(xuerr("xerr", "Unsupported error type: %T", err))
}

// Constructs an error struct from an error inside xgbutil (i.e., user error)
func xuerr (funcName string, err string, params ...interface{}) *XError {
    return &XError{
        funcName: funcName,
        err: fmt.Sprintf(err, params...),
        XGBError: nil,
    }
}

// Dial connects to the X server and creates a new XUtil.
func Dial(display string) (*XUtil, error) {
    c, err := xgb.Dial(display)

    if err != nil {
        return nil, err
    }

    xu := &XUtil{
        conn: c,
        root: c.DefaultScreen().Root,
        atoms: make(map[string]xgb.Id, 50), // start with a nice size
        atomNames: make(map[xgb.Id]string, 50),
    }

    return xu, nil
}

// Conn returns the xgb connection object.
func (xu *XUtil) Conn() (*xgb.Conn) {
    return xu.conn
}

// RootWin returns the current root window.
func (xu *XUtil) RootWin() (xgb.Id) {
    return xu.root
}

// SetRootWin will change the current root window to the one provided.
// N.B. This probably shouldn't be used unless you're desperately trying
// to support multiple X screens. (This is *not* the same as Xinerama/RandR or
// TwinView. All of those have a single root window.)
func (xu *XUtil) SetRootWin(root xgb.Id) {
    xu.root = root
}

// Atm is a short alias for Atom in the common case of interning an atom.
// Namely, only_if_exists is set to true, so that if "name" is an atom that
// does not exist, X will return "0" as an atom identifier. In which case,
// we panic because that isn't what anyone wants.
func (xu *XUtil) Atm(name string) (xgb.Id) {
    if aid := xu.Atom(name, true); aid > 0 {
        return aid
    }

    panic(xuerr("Atm", "'%s' returned an identifier of 0.", name))
}

// Atom interns an atom and panics if there is any error.
func (xu *XUtil) Atom(name string, only_if_exists bool) (xgb.Id) {
    // Check the cache first
    if aid, ok := xu.atoms[name]; ok {
        return aid
    }

    reply, err := xu.conn.InternAtom(only_if_exists, name)

    if err != nil {
        panic(xerr(err, "Atom", "Error interning atom '%s'", name))
    }

    // If we're here, it means we didn't have this atom cached. So cache it!
    xu.atoms[name] = reply.Atom
    xu.atomNames[reply.Atom] = name

    return reply.Atom
}

// AtomName fetches a string representation of an ATOM given its integer id.
func (xu *XUtil) AtomName(aid xgb.Id) string {
    // Check the cache first
    if atomName, ok := xu.atomNames[aid]; ok {
        return string(atomName)
    }

    reply, err := xu.conn.GetAtomName(aid)

    if err != nil {
        panic(xerr(err, "AtomName", "Error fetching name for ATOM id '%d'",
                   aid))
    }

    // If we're here, it means we didn't have ths ATOM id cached. So cache it.
    atomName := string(reply.Name)
    xu.atoms[atomName] = aid
    xu.atomNames[aid] = atomName

    return atomName
}

// GetEwmhWM uses the EWMH spec to find if a conforming window manager
// is currently running or not. If it is, then its name will be returned.
// Otherwise, an error will be returned explaining why one couldn't be found.
// (This function is safe.)
func (xu *XUtil) GetEwmhWM() (wmName string, err error) {
    defer func() {
        if r:= recover(); r != nil {
            wmName = ""
            err = xuerr("GetEwmhWM", "Failed because: %v", r)
        }
    }()

    childCheck := xu.EwmhSupportingWmCheck(xu.root)
    childCheck2 := xu.EwmhSupportingWmCheck(childCheck)

    if childCheck != childCheck2 {
        return "", xuerr("GetEwmhWM",
                         "_NET_SUPPORTING_WM_CHECK value on the root window " +
                         "(%x) does not match _NET_SUPPORTING_WM_CHECK value " +
                         "on the child window (%x).", childCheck, childCheck2)
    }

    return xu.EwmhWmName(childCheck), nil
}

// Safe will recover from any panic produced by xgb or xgbutil and transform
// it into an idiomatic Go error as a second return value.
// NOTE: Generality comes at a cost. The return value will need to be
//       type asserted.
func Safe(fun func() interface{}) (val interface{}, err error) {
    defer func() {
        if r := recover(); r != nil {
            val = nil

            // If we get an error that isn't from xgbutil or xgb itself,
            // then let the panic happen.
            var ok bool
            err, ok = r.(*XError)
            if !ok {
                err, ok = r.(*xgb.Error)
                if !ok { // some other error, panic!
                    panic(r)
                }
            }
        }
    }()

    return fun(), nil
}

// put16 adds a 16 bit integer to a byte slice.
// Lifted from the xgb package.
func put16(buf []byte, v uint16) {
	buf[0] = byte(v)
	buf[1] = byte(v >> 8)
}

// put32 adds a 32 bit integer to a byte slice.
// Lifted from the xgb package.
func put32(buf []byte, v uint32) {
	buf[0] = byte(v)
	buf[1] = byte(v >> 8)
	buf[2] = byte(v >> 16)
	buf[3] = byte(v >> 24)
}

// get16 extracts a 16 bit integer from a byte slice.
// Lifted from the xgb package.
func get16(buf []byte) uint16 {
	v := uint16(buf[0])
	v |= uint16(buf[1]) << 8
	return v
}

// get32 extracts a 32 bit integer from a byte slice.
// Lifted from the xgb package.
func get32(buf []byte) uint32 {
	v := uint32(buf[0])
	v |= uint32(buf[1]) << 8
	v |= uint32(buf[2]) << 16
	v |= uint32(buf[3]) << 24
	return v
}

