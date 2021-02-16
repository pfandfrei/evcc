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
	HandshakeType string    `json:"handshakeType"`
	Version       []Version `json:"version"`
	Formats       []Format  `json:"formats"`
}

type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
}

type Format struct {
	Format string `json:"format"`
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
