/*
   A collection of functions that make working with property replies
   much easier.

   Technically, not all possible property replies are supported (yet).
   But everything needed to implement EWMH, ICCCM and Motif is here.
*/
package xprop

import (
	"burntsushi.net/go/x-go-binding/xgb"
	"burntsushi.net/go/xgbutil"
)

// GetProperty abstracts the messiness of calling xgb.GetProperty.
func GetProperty(xu *xgbutil.XUtil, win xgb.Id, atom string) (
	*xgb.GetPropertyReply, error) {
	atomId, err := Atm(xu, atom)
	if err != nil {
		return nil, err
	}

	reply, err := xu.Conn().GetProperty(false, win, atomId,
		xgb.GetPropertyTypeAny, 0,
		(1<<32)-1)

	if err != nil {
		return nil, xgbutil.Xerr(err, "GetProperty",
			"Error retrieving property '%s' on window %x",
			atom, win)
	}

	if reply.Format == 0 {
		return nil, xgbutil.Xuerr("GetProperty",
			"No such property '%s' on window %x.",
			atom, win)
	}

	return reply, nil
}

// ChangeProperty abstracts the semi-nastiness of xgb.ChangeProperty.
func ChangeProp(xu *xgbutil.XUtil, win xgb.Id, format byte, prop string,
	typ string, data []byte) error {

	propAtom, err := Atm(xu, prop)
	if err != nil {
		return err
	}

	typAtom, err := Atm(xu, typ)
	if err != nil {
		return err
	}

	xu.Conn().ChangeProperty(xgb.PropModeReplace, win, propAtom,
		typAtom, format, data)
	return nil
}

// ChangeProperty32 makes changing 32 bit formatted properties easier
// by constructing the raw X data for you.
func ChangeProp32(xu *xgbutil.XUtil, win xgb.Id, prop string, typ string,
	data ...int) error {

	buf := make([]byte, len(data)*4)
	for i, datum := range data {
		xgbutil.Put32(buf[(i*4):], uint32(datum))
	}

	return ChangeProp(xu, win, 32, prop, typ, buf)
}

// Atm is a short alias for Atom in the common case of interning an atom.
func Atm(xu *xgbutil.XUtil, name string) (xgb.Id, error) {
	aid, err := Atom(xu, name, false)
	if err != nil {
		return 0, err
	}
	if aid == 0 {
		return 0, xgbutil.Xuerr("Atm", "'%s' returned an identifier of 0.",
			name)
	}

	return aid, err
}

// Atom interns an atom and panics if there is any error.
func Atom(xu *xgbutil.XUtil, name string, only_if_exists bool) (xgb.Id, error) {
	// Check the cache first
	if aid, ok := xu.GetAtom(name); ok {
		return aid, nil
	}

	reply, err := xu.Conn().InternAtom(only_if_exists, name)
	if err != nil {
		return 0, xgbutil.Xerr(err, "Atom", "Error interning atom '%s'", name)
	}

	// If we're here, it means we didn't have this atom cached. So cache it!
	xu.CacheAtom(name, reply.Atom)

	return reply.Atom, nil
}

// AtomName fetches a string representation of an ATOM given its integer id.
func AtomName(xu *xgbutil.XUtil, aid xgb.Id) (string, error) {
	// Check the cache first
	if atomName, ok := xu.GetAtomName(aid); ok {
		return string(atomName), nil
	}

	reply, err := xu.Conn().GetAtomName(aid)
	if err != nil {
		return "", xgbutil.Xerr(err, "AtomName",
			"Error fetching name for ATOM id '%d'", aid)
	}

	// If we're here, it means we didn't have ths ATOM id cached. So cache it.
	atomName := string(reply.Name)
	xu.CacheAtom(atomName, aid)

	return atomName, nil
}

// IdTo32 is a covenience function for converting []xgb.Id to []uint32.
func IdToInt(ids []xgb.Id) (ids32 []int) {
	ids32 = make([]int, len(ids))
	for i, v := range ids {
		ids32[i] = int(v)
	}
	return
}

// StrToAtoms is a convenience function for converting
// []string to []uint32 atoms.
// NOTE: If an atom name in the list doesn't exist, it will be created.
func StrToAtoms(xu *xgbutil.XUtil,
	atomNames []string) (atoms []int, err error) {

	atoms = make([]int, len(atomNames))
	for i, atomName := range atomNames {
		a, err := Atom(xu, atomName, false)
		if err != nil {
			return nil, err
		}

		atoms[i] = int(a)
	}
	return
}

// PropValAtom transforms a GetPropertyReply struct into an ATOM name.
// The property reply must be in 32 bit format.
// This is a method of an XUtil struct, unlike the other 'PropVal...' functions.
func PropValAtom(xu *xgbutil.XUtil, reply *xgb.GetPropertyReply,
	err error) (string, error) {

	if err != nil {
		return "", err
	}
	if reply.Format != 32 {
		return "", xgbutil.Xuerr("PropValAtom", "Expected format 32 but got %d",
			reply.Format)
	}

	return AtomName(xu, xgb.Id(xgbutil.Get32(reply.Value)))
}

// PropValAtoms is the same as PropValAtom, except that it returns a slice
// of atom names. Also must be 32 bit format.
// This is a method of an XUtil struct, unlike the other 'PropVal...' functions.
func PropValAtoms(xu *xgbutil.XUtil, reply *xgb.GetPropertyReply,
	err error) ([]string, error) {

	if err != nil {
		return nil, err
	}
	if reply.Format != 32 {
		return nil, xgbutil.Xuerr("PropValAtoms",
			"Expected format 32 but got %d", reply.Format)
	}

	ids := make([]string, reply.ValueLen)
	vals := reply.Value
	for i := 0; len(vals) >= 4; i++ {
		ids[i], err = AtomName(xu, xgb.Id(xgbutil.Get32(vals)))
		if err != nil {
			return nil, err
		}

		vals = vals[4:]
	}

	return ids, nil
}

// PropValId transforms a GetPropertyReply struct into an X resource
// identifier (typically a window id). 
// The property reply must be in 32 bit format.
func PropValId(reply *xgb.GetPropertyReply, err error) (xgb.Id, error) {
	if err != nil {
		return 0, err
	}
	if reply.Format != 32 {
		return 0, xgbutil.Xuerr("PropValId", "Expected format 32 but got %d",
			reply.Format)
	}

	return xgb.Id(xgbutil.Get32(reply.Value)), nil
}

// PropValIds is the same as PropValId, except that it returns a slice
// of identifiers. Also must be 32 bit format.
func PropValIds(reply *xgb.GetPropertyReply, err error) ([]xgb.Id, error) {
	if err != nil {
		return nil, err
	}
	if reply.Format != 32 {
		return nil, xgbutil.Xuerr("PropValIds", "Expected format 32 but got %d",
			reply.Format)
	}

	ids := make([]xgb.Id, reply.ValueLen)
	vals := reply.Value
	for i := 0; len(vals) >= 4; i++ {
		ids[i] = xgb.Id(xgbutil.Get32(vals))
		vals = vals[4:]
	}

	return ids, nil
}

// PropValNum transforms a GetPropertyReply struct into an unsigned
// integer. Useful when the property value is a single integer.
func PropValNum(reply *xgb.GetPropertyReply, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	if reply.Format != 32 {
		return 0, xgbutil.Xuerr("PropValNum", "Expected format 32 but got %d",
			reply.Format)
	}

	return int(xgbutil.Get32(reply.Value)), nil
}

// PropValNums is the same as PropValNum, except that it returns a slice
// of integers. Also must be 32 bit format.
func PropValNums(reply *xgb.GetPropertyReply, err error) ([]int, error) {
	if err != nil {
		return nil, err
	}
	if reply.Format != 32 {
		return nil, xgbutil.Xuerr("PropValIds", "Expected format 32 but got %d",
			reply.Format)
	}

	nums := make([]int, reply.ValueLen)
	vals := reply.Value
	for i := 0; len(vals) >= 4; i++ {
		nums[i] = int(xgbutil.Get32(vals))
		vals = vals[4:]
	}

	return nums, nil
}

// PropValStr transforms a GetPropertyReply struct into a string.
// Useful when the property value is a null terminated string represented
// by integers. Also must be 8 bit format.
func PropValStr(reply *xgb.GetPropertyReply, err error) (string, error) {
	if err != nil {
		return "", err
	}
	if reply.Format != 8 {
		return "", xgbutil.Xuerr("PropValStr", "Expected format 8 but got %d",
			reply.Format)
	}

	return string(reply.Value), nil
}

// PropValStrs is the same as PropValStr, except that it returns a slice
// of strings. The raw byte string is a sequence of null terminated strings,
// which is translated into a slice of strings.
func PropValStrs(reply *xgb.GetPropertyReply, err error) ([]string, error) {
	if err != nil {
		return nil, err
	}
	if reply.Format != 8 {
		return nil, xgbutil.Xuerr("PropValStrs", "Expected format 8 but got %d",
			reply.Format)
	}

	var strs []string
	sstart := 0
	for i, c := range reply.Value {
		if c == 0 {
			strs = append(strs, string(reply.Value[sstart:i]))
			sstart = i + 1
		}
	}

	return strs, nil
}
