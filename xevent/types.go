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
    "code.google.com/p/x-go-binding/xgb"
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

type KeyPressEvent struct {
    *xgb.KeyPressEvent
}

const KeyPress = xgb.KeyPress

func (ev *KeyPressEvent) Bytes() []byte { return nil }

func (ev KeyPressEvent) String() string {
    return fmt.Sprintf("%v", ev.KeyPressEvent)
}

type KeyReleaseEvent struct {
    *xgb.KeyReleaseEvent
}

const KeyRelease = xgb.KeyRelease

func (ev *KeyReleaseEvent) Bytes() []byte { return nil }

func (ev KeyReleaseEvent) String() string {
    return fmt.Sprintf("%v", ev.KeyReleaseEvent)
}

type ButtonPressEvent struct {
    *xgb.ButtonPressEvent
}

const ButtonPress = xgb.ButtonPress

func (ev *ButtonPressEvent) Bytes() []byte { return nil }

func (ev ButtonPressEvent) String() string {
    return fmt.Sprintf("%v", ev.ButtonPressEvent)
}

type ButtonReleaseEvent struct {
    *xgb.ButtonReleaseEvent
}

const ButtonRelease = xgb.ButtonRelease

func (ev *ButtonReleaseEvent) Bytes() []byte { return nil }

func (ev ButtonReleaseEvent) String() string {
    return fmt.Sprintf("%v", ev.ButtonReleaseEvent)
}

type MotionNotifyEvent struct {
    *xgb.MotionNotifyEvent
}

const MotionNotify = xgb.MotionNotify

func (ev *MotionNotifyEvent) Bytes() []byte { return nil }

func (ev MotionNotifyEvent) String() string {
    return fmt.Sprintf("%v", ev.MotionNotifyEvent)
}

type EnterNotifyEvent struct {
    *xgb.EnterNotifyEvent
}

const EnterNotify = xgb.EnterNotify

func (ev *EnterNotifyEvent) Bytes() []byte { return nil }

func (ev EnterNotifyEvent) String() string {
    return fmt.Sprintf("%v", ev.EnterNotifyEvent)
}

type LeaveNotifyEvent struct {
    *xgb.LeaveNotifyEvent
}

const LeaveNotify = xgb.LeaveNotify

func (ev *LeaveNotifyEvent) Bytes() []byte { return nil }

func (ev LeaveNotifyEvent) String() string {
    return fmt.Sprintf("%v", ev.LeaveNotifyEvent)
}

type FocusInEvent struct {
    *xgb.FocusInEvent
}

const FocusIn = xgb.FocusIn

func (ev *FocusInEvent) Bytes() []byte { return nil }

func (ev FocusInEvent) String() string {
    return fmt.Sprintf("%v", ev.FocusInEvent)
}

type FocusOutEvent struct {
    *xgb.FocusOutEvent
}

const FocusOut = xgb.FocusOut

func (ev *FocusOutEvent) Bytes() []byte { return nil }

func (ev FocusOutEvent) String() string {
    return fmt.Sprintf("%v", ev.FocusOutEvent)
}

type KeymapNotifyEvent struct {
    *xgb.KeymapNotifyEvent
}

const KeymapNotify = xgb.KeymapNotify

func (ev *KeymapNotifyEvent) Bytes() []byte { return nil }

func (ev KeymapNotifyEvent) String() string {
    return fmt.Sprintf("%v", ev.KeymapNotifyEvent)
}

type ExposeEvent struct {
    *xgb.ExposeEvent
}

const Expose = xgb.Expose

func (ev *ExposeEvent) Bytes() []byte { return nil }

func (ev ExposeEvent) String() string {
    return fmt.Sprintf("%v", ev.ExposeEvent)
}

type GraphicsExposureEvent struct {
    *xgb.GraphicsExposureEvent
}

const GraphicsExposure = xgb.GraphicsExposure

func (ev *GraphicsExposureEvent) Bytes() []byte { return nil }

func (ev GraphicsExposureEvent) String() string {
    return fmt.Sprintf("%v", ev.GraphicsExposureEvent)
}

type NoExposureEvent struct {
    *xgb.NoExposureEvent
}

const NoExposure = xgb.NoExposure

func (ev *NoExposureEvent) Bytes() []byte { return nil }

func (ev NoExposureEvent) String() string {
    return fmt.Sprintf("%v", ev.NoExposureEvent)
}

type VisibilityNotifyEvent struct {
    *xgb.VisibilityNotifyEvent
}

const VisibilityNotify = xgb.VisibilityNotify

func (ev *VisibilityNotifyEvent) Bytes() []byte { return nil }

func (ev VisibilityNotifyEvent) String() string {
    return fmt.Sprintf("%v", ev.VisibilityNotifyEvent)
}

type CreateNotifyEvent struct {
    *xgb.CreateNotifyEvent
}

const CreateNotify = xgb.CreateNotify

func (ev *CreateNotifyEvent) Bytes() []byte { return nil }

func (ev CreateNotifyEvent) String() string {
    return fmt.Sprintf("%v", ev.CreateNotifyEvent)
}

type DestroyNotifyEvent struct {
    *xgb.DestroyNotifyEvent
}

const DestroyNotify = xgb.DestroyNotify

func (ev *DestroyNotifyEvent) Bytes() []byte { return nil }

func (ev DestroyNotifyEvent) String() string {
    return fmt.Sprintf("%v", ev.DestroyNotifyEvent)
}

type UnmapNotifyEvent struct {
    *xgb.UnmapNotifyEvent
}

const UnmapNotify = xgb.UnmapNotify

func (ev *UnmapNotifyEvent) Bytes() []byte { return nil }

func (ev UnmapNotifyEvent) String() string {
    return fmt.Sprintf("%v", ev.UnmapNotifyEvent)
}

type MapNotifyEvent struct {
    *xgb.MapNotifyEvent
}

const MapNotify = xgb.MapNotify

func (ev *MapNotifyEvent) Bytes() []byte { return nil }

func (ev MapNotifyEvent) String() string {
    return fmt.Sprintf("%v", ev.MapNotifyEvent)
}

type MapRequestEvent struct {
    *xgb.MapRequestEvent
}

const MapRequest = xgb.MapRequest

func (ev *MapRequestEvent) Bytes() []byte { return nil }

func (ev MapRequestEvent) String() string {
    return fmt.Sprintf("%v", ev.MapRequestEvent)
}

type ReparentNotifyEvent struct {
    *xgb.ReparentNotifyEvent
}

const ReparentNotify = xgb.ReparentNotify

func (ev *ReparentNotifyEvent) Bytes() []byte { return nil }

func (ev ReparentNotifyEvent) String() string {
    return fmt.Sprintf("%v", ev.ReparentNotifyEvent)
}

type ConfigureNotifyEvent struct {
    *xgb.ConfigureNotifyEvent
}

const ConfigureNotify = xgb.ConfigureNotify

func (ev *ConfigureNotifyEvent) Bytes() []byte { return nil }

func (ev ConfigureNotifyEvent) String() string {
    return fmt.Sprintf("%v", ev.ConfigureNotifyEvent)
}

type ConfigureRequestEvent struct {
    *xgb.ConfigureRequestEvent
}

const ConfigureRequest = xgb.ConfigureRequest

func (ev *ConfigureRequestEvent) Bytes() []byte { return nil }

func (ev ConfigureRequestEvent) String() string {
    return fmt.Sprintf("%v", ev.ConfigureRequestEvent)
}

type GravityNotifyEvent struct {
    *xgb.GravityNotifyEvent
}

const GravityNotify = xgb.GravityNotify

func (ev *GravityNotifyEvent) Bytes() []byte { return nil }

func (ev GravityNotifyEvent) String() string {
    return fmt.Sprintf("%v", ev.GravityNotifyEvent)
}

type ResizeRequestEvent struct {
    *xgb.ResizeRequestEvent
}

const ResizeRequest = xgb.ResizeRequest

func (ev *ResizeRequestEvent) Bytes() []byte { return nil }

func (ev ResizeRequestEvent) String() string {
    return fmt.Sprintf("%v", ev.ResizeRequestEvent)
}

type CirculateNotifyEvent struct {
    *xgb.CirculateNotifyEvent
}

const CirculateNotify = xgb.CirculateNotify

func (ev *CirculateNotifyEvent) Bytes() []byte { return nil }

func (ev CirculateNotifyEvent) String() string {
    return fmt.Sprintf("%v", ev.CirculateNotifyEvent)
}

type CirculateRequestEvent struct {
    *xgb.CirculateRequestEvent
}

const CirculateRequest = xgb.CirculateRequest

func (ev *CirculateRequestEvent) Bytes() []byte { return nil }

func (ev CirculateRequestEvent) String() string {
    return fmt.Sprintf("%v", ev.CirculateRequestEvent)
}

type PropertyNotifyEvent struct {
    *xgb.PropertyNotifyEvent
}

const PropertyNotify = xgb.PropertyNotify

func (ev *PropertyNotifyEvent) Bytes() []byte { return nil }

func (ev PropertyNotifyEvent) String() string {
    return fmt.Sprintf("%v", ev.PropertyNotifyEvent)
}

type SelectionClearEvent struct {
    *xgb.SelectionClearEvent
}

const SelectionClear = xgb.SelectionClear

func (ev *SelectionClearEvent) Bytes() []byte { return nil }

func (ev SelectionClearEvent) String() string {
    return fmt.Sprintf("%v", ev.SelectionClearEvent)
}

type SelectionRequestEvent struct {
    *xgb.SelectionRequestEvent
}

const SelectionRequest = xgb.SelectionRequest

func (ev *SelectionRequestEvent) Bytes() []byte { return nil }

func (ev SelectionRequestEvent) String() string {
    return fmt.Sprintf("%v", ev.SelectionRequestEvent)
}

type SelectionNotifyEvent struct {
    *xgb.SelectionNotifyEvent
}

const SelectionNotify = xgb.SelectionNotify

func (ev *SelectionNotifyEvent) Bytes() []byte { return nil }

func (ev SelectionNotifyEvent) String() string {
    return fmt.Sprintf("%v", ev.SelectionNotifyEvent)
}

type ColormapNotifyEvent struct {
    *xgb.ColormapNotifyEvent
}

const ColormapNotify = xgb.ColormapNotify

func (ev *ColormapNotifyEvent) Bytes() []byte { return nil }

func (ev ColormapNotifyEvent) String() string {
    return fmt.Sprintf("%v", ev.ColormapNotifyEvent)
}

type MappingNotifyEvent struct {
    *xgb.MappingNotifyEvent
}

const MappingNotify = xgb.MappingNotify

func (ev *MappingNotifyEvent) Bytes() []byte { return nil }

func (ev MappingNotifyEvent) String() string {
    return fmt.Sprintf("%v", ev.MappingNotifyEvent)
}

