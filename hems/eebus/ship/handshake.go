package ship

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
)

type CmiProtocolHandshakeError struct {
	Error int `json:"error"`
}

type CmiHandshakeMsg struct {
	MessageProtocolHandshake []MessageProtocolHandshake `json:"messageProtocolHandshake"`
}

type MessageProtocolHandshake struct {
	HandshakeType string   `json:"handshakeType"`
	Version       Version  `json:"version"`
	Formats       []string `json:"formats"`
}

type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
}

type CmiConnectionPinState struct {
	ConnectionPinState []ConnectionPinState `json:"connectionPinState"`
}

type ConnectionPinState struct {
	PinState string `json:"pinState"`
}

type CmiAccessMethodsRequest struct {
	AccessMethodsRequest []AccessMethodsRequest `json:"accessMethodsRequest"`
}

type AccessMethodsRequest struct {
	ID  string `json:"dnsSd_mDns,omitempty"`
	DNS struct {
		URI string `json:"uri"`
	} `json:"dns,omitempty"`
}

type CmiAccessMethods struct {
	AccessMethods []AccessMethods `json:"accessMethods"`
}

type AccessMethods struct {
	ID  string `json:"dnsSd_mDns,omitempty"`
	DNS struct {
		URI string `json:"uri"`
	} `json:"dns,omitempty"`
}
