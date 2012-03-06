/*
    evtypes provides an interface to convert nice event structs into
    byte streams for use in xgb.SendEvent.

    As with property.go, evtypes is not feature complete yet. Namely,
    it doesn't support all events. (It probably only has enough events to
    make xgbutil's core functionality work.)

    The pattern here is (or should be):
        Define an event type by embedding the corresponding event type from XGB.
        Create a new constant variable holding that event's unique code.
        Define a 'New...' function that creates a value of that type.
        Define a 'Bytes' method on that type to satisfy the 'XEvent' interface.

*/
package xevent

import "fmt"

import (
    "code.google.com/p/jamslam-x-go-binding/xgb"
    "github.com/BurntSushi/xgbutil"
)

// XEvent is an interface whereby an event struct ought to be convertible into
// a raw slice of bytes, X protocol style.
type XEvent interface {
    Bytes() []byte
}

// ClientMessageEvent embeds the struct by the same name from the xgb library.
type ClientMessageEvent struct {
    *xgb.ClientMessageEvent
}

// The unique code for a ClientMessage event.
const ClientMessage = 33

// NewClientMessage takes all arguments required to build a ClientMessageEvent 
// struct and hides the messy details.
// The varidic parameters coincide with the "data" part of a client message.
// Right now, this function only supports a list of up to 5 uint32s.
// XXX: Use type assertions to support bytes and uint16s.
func NewClientMessage(Format byte, Window xgb.Id, Type xgb.Id,
                      data ...interface{}) (*ClientMessageEvent, error) {
    // Create the client data list first
    clientData := new(xgb.ClientMessageData)

    // Don't support formats 8 or 16 yet. They aren't used in EWMH anyway.
    switch Format {
    case 8:
        // copy(clientData.Data8[:], data.([]byte)) 
        // Using a loop here instead of a straight copy because
        // it appears I can't use type assertions like 'data.([]byte)'.
        // I'm still on my second day with Go, so I'm not sure why that is yet.
        for i := 0; i < 20; i++ {
            if i >= len(data) {
                clientData.Data8[i] = 0
            } else {
                clientData.Data8[i] = data[i].(byte)
            }
        }
    case 16:
        for i := 0; i < 10; i++ {
            if i >= len(data) {
                clientData.Data16[i] = 0
            } else {
                clientData.Data16[i] = data[i].(uint16)
            }
        }
    case 32:
        for i := 0; i < 5; i++ {
            if i >= len(data) {
                clientData.Data32[i] = 0
            } else {
                clientData.Data32[i] = data[i].(uint32)
            }
        }
    default:
        return nil, xgbutil.Xuerr("NewClientMessage",
                                  "Unsupported format '%d'.", Format)
    }

    return &ClientMessageEvent{&xgb.ClientMessageEvent{
        Format: 32,
        Window: Window,
        Type: Type,
        Data: *clientData,
    }}, nil
}

// Bytes transforms a ClientMessageEvent struct into a 32 byte slice.
func (ev *ClientMessageEvent) Bytes() []byte {
    buf := make([]byte, 32)

    buf[0] = xgb.ClientMessage
    buf[1] = ev.Format
    xgbutil.Put32(buf[4:], uint32(ev.Window))
    xgbutil.Put32(buf[8:], uint32(ev.Type))

    // ClientMessage data is a 20 byte list and can be one of:
    // 20 8-bit values
    // 10 16-bit values
    // 5  32-bit values
    // Therefore, check 'Format' and grab the appropriate data and copy
    data := buf[12:]
    switch ev.Format {
    case 8:
        copy(data, ev.Data.Data8[:])
    case 16:
        for i, datum := range ev.Data.Data16 {
            xgbutil.Put16(data[(i * 2):], datum)
        }
    case 32:
        for i, datum := range ev.Data.Data32 {
            xgbutil.Put32(data[(i * 4):], datum)
        }
    default:
        panic(xgbutil.Xuerr("Bytes", "Unsupported format '%d'.", ev.Format))
    }

    return buf
}

// String just shows the embedded type from XGB.
func (ev ClientMessageEvent) String() string {
    return fmt.Sprintf("%v", ev.ClientMessageEvent)
}

// The rest of the types don't implement 'Bytes' yet, but they should.
// These are also exposed to the user when constructing event callbacks.
// Namely, one should use these types instead of XGB event types.
// There isn't any particular reason why *now*, but this will make it future
// proof if more needs to be done with events.
// TODO: Write 'Bytes' functions for the rest of the events.
// XXX: The following is generated from scripts/write-events ATM.
// N.B. These types are in 'types_auto.go' to simplify generation.

