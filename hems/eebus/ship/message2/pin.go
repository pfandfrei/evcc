package message

import "github.com/andig/evcc/hems/eebus/util"

type CmiConnectionPinState struct {
	ConnectionPinState ConnectionPinState `json:"connectionPinState"`
}

type ConnectionPinState struct {
	PinState        string `json:"pinState"`
	InputPermission string `json:"inputPermission,omitempty"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m ConnectionPinState) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *ConnectionPinState) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

type CmiConnectionPinInput struct {
	ConnectionPinInput ConnectionPinInput `json:"connectionPinInput"`
}

type ConnectionPinInput struct {
	Pin string `json:"pin"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m ConnectionPinInput) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *ConnectionPinInput) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

type CmiConnectionPinError struct {
	ConnectionPinError ConnectionPinError `json:"connectionPinError"`
}

type ConnectionPinError struct {
	Error byte `json:"error"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m ConnectionPinError) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *ConnectionPinError) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}
