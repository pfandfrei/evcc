package message

import "github.com/andig/evcc/hems/eebus/util"

// data constants
const (
	CmiTypeData byte = 2
)

type CmiData struct {
	Datagram Datagram `json:"datagram"`
}

type Datagram struct {
	Header    Header    `json:"header"`
	Payload   Payload   `json:"payload"`
	Extension Extension `json:"extension"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m Datagram) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *Datagram) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

type Header struct {
	SpecificationVersion string  `json:"specificationVersion"`
	AddressSource        Address `json:"addressSource"`
	AddressDestination   Address `json:"addressDestination"`
	MsgCounter           int     `json:"msgCounter"`
	CmdClassifier        string  `json:"cmdClassifier"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m Header) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *Header) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

type Address struct {
	Device  string `json:"device"`
	Entity  []int  `json:"entity"`
	Feature int    `json:"feature"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m Address) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *Address) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

type Payload struct {
	Cmd []interface{} `json:"cmd"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m Payload) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *Payload) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}

type Extension struct {
	ExtensionID string `json:"extensionId"`
	Binary      []byte `json:"binary,omitempty"`
	String      string `json:"string,omitempty"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m Extension) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *Extension) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}
