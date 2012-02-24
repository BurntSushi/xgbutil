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

// PropValId transforms a GetPropertyReply struct into an X resource
// identifier (typically a window id). 
// The property reply must be in 32 bit format.
func PropValId(reply *xgb.GetPropertyReply) (xgb.Id) {
    if reply.Format != 32 {
        panic(perr("PropValId: Expected format 32 but got %d",
                   reply.Format))
    }

    return xgb.Id(get32(reply.Value))
}

// PropValIds is the same as PropValId, except that it returns a slice
// of identifiers. Also must be 32 bit format.
func PropValIds(reply *xgb.GetPropertyReply) []xgb.Id {
    if reply.Format != 32 {
        panic(perr("PropValIds: Expected format 32 but got %d",
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
        panic(perr("PropValNum: Expected format 32 but got %d",
                   reply.Format))
    }

    return get32(reply.Value)
}

// PropValNums is the same as PropValNum, except that it returns a slice
// of integers. Also must be 32 bit format.
func PropValNums(reply *xgb.GetPropertyReply) []uint32 {
    if reply.Format != 32 {
        panic(perr("PropValIds: Expected format 32 but got %d",
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
        panic(perr("PropValStr: Expected format 8 but got %d", reply.Format))
    }

    return string(reply.Value)
}

// PropValStrs is the same as PropValStr, except that it returns a slice
// of strings. The raw byte string is a sequence of null terminated strings,
// which is translated into a slice of strings.
func PropValStrs(reply *xgb.GetPropertyReply) []string {
    if reply.Format != 8 {
        panic(perr("PropValStrs: Expected format 8 but got %d", reply.Format))
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

