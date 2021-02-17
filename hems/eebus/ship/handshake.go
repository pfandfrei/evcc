package ship

import (
	"encoding/json"

	"github.com/thoas/go-funk"
)

// handshake constants
const (
	CmiTypeControl byte = 1

	ProtocolHandshakeFormatJSON = "JSON-UTF8"

	ProtocolHandshakeTypeAnnounceMax = "announceMax"
	ProtocolHandshakeTypeSelect      = "select"

	CmiProtocolHandshakeErrorUnexpectedMessage = 2

	// Pin states
	PinStateRequired = "required"
	PinStateOptional = "optional"
	PinStatePinOk    = "pinok"
	PinStateNone     = "none"

	PinInputPermissionOk   = "ok"
	PinInputPermissionBusy = "busy"
)

type CmiProtocolHandshakeError struct {
	Error int `json:"error"`
}

type CmiHandshakeMsg struct {
	MessageProtocolHandshake []MessageProtocolHandshake `json:"messageProtocolHandshake"`
}

type MessageProtocolHandshake struct {
	HandshakeType HandshakeType `json:"handshakeType"`
	Version       Version       `json:"version"`
	Formats       Formats       `json:"formats"`
}

type HandshakeType struct {
	HandshakeType string
}

type Formats []Format

func (f Formats) Supports(required string) bool {
	for _, format := range f {
		if funk.ContainsString(format.Format, required) {
			return true
		}

	}
	return false
}

func (h CmiHandshakeMsg) MarshalJSON() ([]byte, error) {
	wrapper := struct {
		MessageProtocolHandshake []interface{} `json:"messageProtocolHandshake"`
	}{
		MessageProtocolHandshake: []interface{}{
			struct {
				HandshakeType string `json:"handshakeType"`
			}{
				h.MessageProtocolHandshake[0].HandshakeType.HandshakeType,
			},
			struct {
				Version Version `json:"version"`
			}{
				h.MessageProtocolHandshake[0].Version,
			},
			struct {
				Formats []Format `json:"formats"`
			}{
				h.MessageProtocolHandshake[0].Formats,
			},
		},
	}

	return json.Marshal(wrapper)
}

func (hs *MessageProtocolHandshake) UnmarshalJSON(b []byte) error {
	var wrapper []json.RawMessage

	err := json.Unmarshal(b, &wrapper)
	if err == nil && len(wrapper) == 0 {
		return &json.UnmarshalTypeError{Value: string(b)}
	}

	if err == nil && len(wrapper) > 0 {
		err = json.Unmarshal(wrapper[0], &hs.HandshakeType)
	}
	if err == nil && len(wrapper) > 1 {
		var v struct{ Version Version }
		if err = json.Unmarshal(wrapper[1], &v); err == nil {
			hs.Version = v.Version
		}
	}
	if err == nil && len(wrapper) > 2 {
		var v struct{ Formats Formats }
		if err = json.Unmarshal(wrapper[2], &v); err == nil {
			hs.Formats = v.Formats
		}
	}

	return err
}

type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
}

func (h Version) MarshalJSON() ([]byte, error) {
	wrapper := []interface{}{
		map[string]int{
			"major": h.Major,
		},
		map[string]int{
			"minor": h.Minor,
		},
	}

	return json.Marshal(wrapper)
}

func (h *Version) UnmarshalJSON(b []byte) error {
	var wrapper []json.RawMessage
	err := json.Unmarshal(b, &wrapper)

	type version Version
	h2 := (*version)(h)

	if err == nil {
		for _, v := range wrapper {
			if err = json.Unmarshal(v, &h2); err != nil {
				break
			}
		}
	}

	return err
}

type Format struct {
	Format []string `json:"format"`
}

type CmiConnectionPinState struct {
	ConnectionPinState []ConnectionPinState `json:"connectionPinState"`
}

type ConnectionPinState struct {
	PinState        string `json:"pinState"`
	InputPermission string `json:"inputPermission,omitempty"`
}

type CmiConnectionPinInput struct {
	ConnectionPinInput []ConnectionPinInput `json:"connectionPinInput"`
}

type ConnectionPinInput struct {
	Pin string `json:"pin"`
}

type CmiConnectionPinError struct {
	ConnectionPinError []ConnectionPinError `json:"connectionPinError"`
}

type ConnectionPinError struct {
	Error byte `json:"error"`
}

type CmiAccessMethodsRequest struct {
	AccessMethodsRequest []AccessMethodsRequest `json:"accessMethodsRequest"`
}

type AccessMethodsRequest struct{}

type CmiAccessMethods struct {
	AccessMethods []AccessMethods `json:"accessMethods"`
}

type AccessMethods struct {
	ID        string `json:"id"`
	DnsSDmDns string `json:"dnsSd_mDns,omitempty"`
	Dns       *struct {
		URI string `json:"uri"`
	} `json:"dns,omitempty"`
}
