/*
    A collection of functions that make working with property replies
    much easier.

    Technically, not all possible property replies are supported (yet).
    But everything needed to implement EWMH, ICCCM and Motif is here.
*/
package xgbutil

import (
    // "fmt" 

    "code.google.com/p/x-go-binding/xgb"
)

// GetProperty abstracts the messiness of calling xgb.GetProperty.
func (xu *XUtil) GetProperty(win xgb.Id, atom string) (*xgb.GetPropertyReply) {
    reply, err := xu.conn.GetProperty(false, win, xu.Atm(atom),
                                      xgb.GetPropertyTypeAny, 0, (1 << 32) - 1)

    if err != nil {
        panic(xerr(err, "GetProperty",
                   "Error retrieving property '%s' on window %x", atom, win))
    }

    if reply.Format == 0 {
        panic(xuerr("GetProperty", "No such property '%s' on window %x.",
                    atom, win))
    }

    return reply
}

// SafeGetProperty is a wrapper for GetProperty that will transform panics
// into Go idiomatic errors.
func (xu *XUtil) SafeGetProperty(win xgb.Id, atom string) (
     reply *xgb.GetPropertyReply, err error) {
    maybeReply, err := Safe(func() interface{} {
        return xu.GetProperty(win, atom)
    })

    if err != nil {
        reply = nil
    } else {
        reply = maybeReply.(*xgb.GetPropertyReply)
    }

    return
}

// ChangeProperty abstracts the semi-nastiness of xgb.ChangeProperty.
func (xu *XUtil) ChangeProperty(win xgb.Id, format byte, prop string,
                                typ string, data []byte) {
    xu.conn.ChangeProperty(xgb.PropModeReplace, win, xu.Atm(prop),
                           xu.Atm(typ), format, data)
}

// ChangeProperty32 makes changing 32 bit formatted properties easier
// by constructing the raw X data for you.
func (xu *XUtil) ChangeProperty32(win xgb.Id, prop string, typ string,
                                  data ...uint32) {
    buf := make([]byte, len(data) * 4)
    for i, datum := range data {
        put32(buf[(i * 4):], datum)
    }

    xu.ChangeProperty(win, 32, prop, typ, buf)
}

// IdTo32 is a covenience function for converting []xgb.Id to []uint32.
func IdTo32(ids []xgb.Id) (ids32 []uint32) {
    ids32 = make([]uint32, len(ids))
    for i, v := range ids {
        ids32[i] = uint32(v)
    }
    return
}

// StrToAtoms is a convenience function for converting
// []string to []uint32 atoms.
// NOTE: If an atom name in the list doesn't exist, it will be created.
func (xu *XUtil) StrToAtoms(atomNames []string) (atoms []uint32) {
    atoms = make([]uint32, len(atomNames))
    for i, atomName := range atomNames {
        atoms[i] = uint32(xu.Atom(atomName, false))
    }
    return
}

// PropValAtom transforms a GetPropertyReply struct into an ATOM name.
// The property reply must be in 32 bit format.
// This is a method of an XUtil struct, unlike the other 'PropVal...' functions.
func (xu *XUtil) PropValAtom(reply *xgb.GetPropertyReply) string {
    if reply.Format != 32 {
        panic(xuerr("PropValAtom", "Expected format 32 but got %d",
                    reply.Format))
    }

    return xu.AtomName(xgb.Id(get32(reply.Value)))
}

// PropValAtoms is the same as PropValAtom, except that it returns a slice
// of atom names. Also must be 32 bit format.
// This is a method of an XUtil struct, unlike the other 'PropVal...' functions.
func (xu *XUtil) PropValAtoms(reply *xgb.GetPropertyReply) []string {
    if reply.Format != 32 {
        panic(xuerr("PropValAtoms", "Expected format 32 but got %d",
                    reply.Format))
    }

    ids := make([]string, reply.ValueLen)
    vals := reply.Value
    for i := 0; len(vals) >= 4; i++ {
        ids[i] = xu.AtomName(xgb.Id(get32(vals)))
        vals = vals[4:]
    }

    return ids
}

// PropValId transforms a GetPropertyReply struct into an X resource
// identifier (typically a window id). 
// The property reply must be in 32 bit format.
func PropValId(reply *xgb.GetPropertyReply) (xgb.Id) {
    if reply.Format != 32 {
        panic(xuerr("PropValId", "Expected format 32 but got %d", reply.Format))
    }

    return xgb.Id(get32(reply.Value))
}

// PropValIds is the same as PropValId, except that it returns a slice
// of identifiers. Also must be 32 bit format.
func PropValIds(reply *xgb.GetPropertyReply) []xgb.Id {
    if reply.Format != 32 {
        panic(xuerr("PropValIds", "Expected format 32 but got %d",
                    reply.Format))
    }

    ids := make([]xgb.Id, reply.ValueLen)
    vals := reply.Value
    for i := 0; len(vals) >= 4; i++ {
        ids[i] = xgb.Id(get32(vals))
        vals = vals[4:]
    }

    return ids
}

// PropValNum transforms a GetPropertyReply struct into an unsigned
// integer. Useful when the property value is a single integer.
func PropValNum(reply *xgb.GetPropertyReply) (uint32) {
    if reply.Format != 32 {
        panic(xuerr("PropValNum", "Expected format 32 but got %d",
                    reply.Format))
    }

    return get32(reply.Value)
}

// PropValNums is the same as PropValNum, except that it returns a slice
// of integers. Also must be 32 bit format.
func PropValNums(reply *xgb.GetPropertyReply) []uint32 {
    if reply.Format != 32 {
        panic(xuerr("PropValIds", "Expected format 32 but got %d",
                    reply.Format))
    }

    nums := make([]uint32, reply.ValueLen)
    vals := reply.Value
    for i := 0; len(vals) >= 4; i++ {
        nums[i] = get32(vals)
        vals = vals[4:]
    }

    return nums
}

// PropValStr transforms a GetPropertyReply struct into a string.
// Useful when the property value is a null terminated string represented
// by integers. Also must be 8 bit format.
func PropValStr(reply *xgb.GetPropertyReply) string {
    if reply.Format != 8 {
        panic(xuerr("PropValStr", "Expected format 8 but got %d", reply.Format))
    }

    return string(reply.Value)
}

// PropValStrs is the same as PropValStr, except that it returns a slice
// of strings. The raw byte string is a sequence of null terminated strings,
// which is translated into a slice of strings.
func PropValStrs(reply *xgb.GetPropertyReply) []string {
    if reply.Format != 8 {
        panic(xuerr("PropValStrs", "Expected format 8 but got %d",
                    reply.Format))
    }

    var strs []string
    sstart := 0
    for i, c := range reply.Value {
        if c == 0 {
            strs = append(strs, string(reply.Value[sstart:i]))
            sstart = i + 1
        }
    }

    return strs
}

