package xevent

import (
	"fmt"

	"github.com/BurntSushi/xgb/xproto"
)

// ClientMessageEvent embeds the struct by the same name from the xgb library.
type ClientMessageEvent struct {
	*xproto.ClientMessageEvent
}

const ClientMessage = xproto.ClientMessage

func number2byte(v interface{}) byte {
	switch rv := v.(type) {
	case int8:
		return byte(rv)
	case int16:
		return byte(rv)
	case int32:
		return byte(rv)
	case int:
		return byte(rv)
	case int64:
		return byte(rv)
	case uint8:
		return byte(rv)
	case uint16:
		return byte(rv)
	case uint32:
		return byte(rv)
	case uint:
		return byte(rv)
	case uint64:
		return byte(rv)
	}
	panic("number2byte: unspoorted type")
}
func number2uint16(v interface{}) uint16 {
	switch rv := v.(type) {
	case int8:
		return uint16(rv)
	case int16:
		return uint16(rv)
	case int32:
		return uint16(rv)
	case int:
		return uint16(rv)
	case int64:
		return uint16(rv)
	case uint8:
		return uint16(rv)
	case uint16:
		return uint16(rv)
	case uint32:
		return uint16(rv)
	case uint:
		return uint16(rv)
	case uint64:
		return uint16(rv)
	}
	panic("number2byte: unspoorted type")
}
func number2uint32(v interface{}) uint32 {
	switch rv := v.(type) {
	case int8:
		return uint32(rv)
	case int16:
		return uint32(rv)
	case int32:
		return uint32(rv)
	case int:
		return uint32(rv)
	case int64:
		return uint32(rv)
	case uint8:
		return uint32(rv)
	case uint16:
		return uint32(rv)
	case uint32:
		return uint32(rv)
	case uint64:
		return uint32(rv)
	case uint:
		return uint32(rv)
	}
	panic("number2byte: unspoorted type")
}

// NewClientMessage takes all arguments required to build a ClientMessageEvent
// struct and hides the messy details.
// The variadic parameters coincide with the "data" part of a client message.
// The type of the variadic parameters depends upon the value of Format.
// If Format is 8, 'data' should have type byte.
// If Format is 16, 'data' should have type int16.
// If Format is 32, 'data' should have type int.
// Any other value of Format returns an error.
func NewClientMessage(Format byte, Window xproto.Window, Type xproto.Atom,
	data ...interface{}) (*ClientMessageEvent, error) {

	// Create the client data list first
	var clientData xproto.ClientMessageDataUnion

	// Don't support formats 8 or 16 yet. They aren't used in EWMH anyway.
	switch Format {
	case 8:
		buf := make([]byte, 20)
		for i := 0; i < 20; i++ {
			if i >= len(data) {
				break
			}
			buf[i] = number2byte(data[i])
		}
		clientData = xproto.ClientMessageDataUnionData8New(buf)
	case 16:
		buf := make([]uint16, 10)
		for i := 0; i < 10; i++ {
			if i >= len(data) {
				break
			}
			buf[i] = number2uint16(data[i])
		}
		clientData = xproto.ClientMessageDataUnionData16New(buf)
	case 32:
		buf := make([]uint32, 5)
		for i := 0; i < 5; i++ {
			if i >= len(data) {
				break
			}
			buf[i] = number2uint32(data[i])
		}
		clientData = xproto.ClientMessageDataUnionData32New(buf)
	default:
		return nil, fmt.Errorf("NewClientMessage: Unsupported format '%d'.",
			Format)
	}

	return &ClientMessageEvent{&xproto.ClientMessageEvent{
		Format: Format,
		Window: Window,
		Type:   Type,
		Data:   clientData,
	}}, nil
}

// ConfigureNotifyEvent embeds the struct by the same name in XGB.
type ConfigureNotifyEvent struct {
	*xproto.ConfigureNotifyEvent
}

const ConfigureNotify = xproto.ConfigureNotify

// NewConfigureNotify takes all arguments required to build a
// ConfigureNotifyEvent struct and returns a ConfigureNotifyEvent value.
func NewConfigureNotify(Event, Window, AboveSibling xproto.Window,
	X, Y, Width, Height int, BorderWidth uint16,
	OverrideRedirect bool) *ConfigureNotifyEvent {

	return &ConfigureNotifyEvent{&xproto.ConfigureNotifyEvent{
		Event: Event, Window: Window, AboveSibling: AboveSibling,
		X: int16(X), Y: int16(Y), Width: uint16(Width), Height: uint16(Height),
		BorderWidth: BorderWidth, OverrideRedirect: OverrideRedirect,
	}}
}
