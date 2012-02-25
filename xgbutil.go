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
}

// Alias for error printing
var perr = fmt.Sprintf

// Dial connects to the X server and creates a new XUtil.
func Dial(display string) (*XUtil, error) {
    c, err := xgb.Dial(display)

    if err != nil {
        return nil, err
    }

    xu := &XUtil{
        conn: c,
        root: c.Setup.Roots[0].Root,
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

    panic(perr("Atom '%s' returned an identifier of 0.", name))
}

// Atom interns an atom and panics if there is any error.
func (xu *XUtil) Atom(name string, only_if_exists bool) (xgb.Id) {
    reply, err := xu.conn.InternAtom(only_if_exists, name)

    if err != nil {
        panic(perr("Error interning atom '%s': %v", name, err))
    }

    return reply.Atom
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

