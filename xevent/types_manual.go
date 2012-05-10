package xevent

import (
	"fmt"

	"github.com/BurntSushi/xgb"
)

// ClientMessageEvent embeds the struct by the same name from the xgb library.
type ClientMessageEvent struct {
	*xgb.ClientMessageEvent
}

// The unique code for a ClientMessage event.
const ClientMessage = xgb.ClientMessage

// NewClientMessage takes all arguments required to build a ClientMessageEvent 
// struct and hides the messy details.
// The varidic parameters coincide with the "data" part of a client message.
// Right now, this function only supports a list of up to 5 uint32s.
// XXX: Use type assertions to support bytes and uint16s.
func NewClientMessage(Format byte, Window xgb.Id, Type xgb.Id,
	data ...interface{}) (*ClientMessageEvent, error) {

	// Create the client data list first
	var clientData xgb.ClientMessageDataUnion

	// Don't support formats 8 or 16 yet. They aren't used in EWMH anyway.
	switch Format {
	case 8:
		buf := make([]byte, 20)
		for i := 0; i < 20; i++ {
			if i >= len(data) {
				break
			}
			buf[i] = data[i].(byte)
		}
		clientData = xgb.NewClientMessageDataUnionData8(buf)
	case 16:
		buf := make([]uint16, 10)
		for i := 0; i < 10; i++ {
			if i >= len(data) {
				break
			}
			buf[i] = uint16(data[i].(int16))
		}
		clientData = xgb.NewClientMessageDataUnionData16(buf)
	case 32:
		buf := make([]uint32, 5)
		for i := 0; i < 5; i++ {
			if i >= len(data) {
				break
			}
			buf[i] = uint32(data[i].(int))
		}
		clientData = xgb.NewClientMessageDataUnionData32(buf)
	default:
		return nil, fmt.Errorf("NewClientMessage: Unsupported format '%d'.",
			Format)
	}

	return &ClientMessageEvent{&xgb.ClientMessageEvent{
		Format: Format,
		Window: Window,
		Type:   Type,
		Data:   clientData,
	}}, nil
}

// ConfigureNotifyEvent embeds the struct by the same name in XGB.
type ConfigureNotifyEvent struct {
	*xgb.ConfigureNotifyEvent
}

// The unique code for a ConfigureNotify event.
const ConfigureNotify = xgb.ConfigureNotify

// NewConfigureNotify takes all arguments required to build a 
// ConfigureNotifyEvent struct and hides the messy details.
func NewConfigureNotify(Event, Window, AboveSibling xgb.Id,
	X, Y, Width, Height int, BorderWidth uint16,
	OverrideRedirect bool) *ConfigureNotifyEvent {

	return &ConfigureNotifyEvent{&xgb.ConfigureNotifyEvent{
		Event: Event, Window: Window, AboveSibling: AboveSibling,
		X: int16(X), Y: int16(Y), Width: uint16(Width), Height: uint16(Height),
		BorderWidth: BorderWidth, OverrideRedirect: OverrideRedirect,
	}}
}
