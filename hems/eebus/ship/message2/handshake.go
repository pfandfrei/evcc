package message

import (
	"github.com/andig/evcc/hems/eebus/util"
	"github.com/thoas/go-funk"
)

// handshake constants
const (
	CmiTypeControl byte = 1

	ProtocolHandshakeFormatJSON = "JSON-UTF8"

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
	MessageProtocolHandshake MessageProtocolHandshake `json:"messageProtocolHandshake"`
}

type MessageProtocolHandshake struct {
	HandshakeType string  `json:"handshakeType"`
	Version       Version `json:"version"`
	Formats       Format  `json:"formats"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m MessageProtocolHandshake) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *MessageProtocolHandshake) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// Version of the protocol
type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m Version) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *Version) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

// Format of the protocol
type Format struct {
	Format []string `json:"format"`
}

// IsSupported validates if format is supported
func (m Format) IsSupported(format string) bool {
	return funk.ContainsString(m.Format, format)
}

// MarshalJSON is the SHIP serialization marshaller
func (m Format) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *Format) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}
