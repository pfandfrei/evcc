package ship

import (
	"errors"
	"fmt"
)

const (
	CmiTypeControl byte = 1
)

const (
	ProtocolHandshakeFormatJSON = "JSON-UTF8"

	ProtocolHandshakeTypeAnnounceMax = "announceMax"
	ProtocolHandshakeTypeSelect      = "select"

	SubProtocol = "ship"
)

type MessageProtocolHandshake struct {
	HandshakeType string   `json:"handshakeType"`
	Version       Version  `json:"version"`
	Formats       []string `json:"formats"`
}

type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
}

const (
	CmiProtocolHandshakeErrorUnexpectedMessage = 2
)

type CmiProtocolHandshakeError struct {
	Error int `json:"error"`
}

type CmiHandshakeMsg struct {
	MessageProtocolHandshake []MessageProtocolHandshake `json:"messageProtocolHandshake"`
}

type CmiConnectionPinState struct {
	ConnectionPinState []ConnectionPinState `json:"connectionPinState"`
}

type ConnectionPinState struct {
	PinState string `json:"pinState"` // required, optional, pinOk, none
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

func (c *Connection) handshakeReceiveSelect() (CmiHandshakeMsg, error) {
	var resp CmiHandshakeMsg
	typ, err := c.readJSON(&resp)

	if err == nil && typ != CmiTypeControl {
		err = fmt.Errorf("handshake: invalid type: %0x", typ)
	}

	if err == nil {
		if len(resp.MessageProtocolHandshake) != 1 {
			return resp, errors.New("handshake: invalid length")
		}

		handshake := resp.MessageProtocolHandshake[0]

		if handshake.HandshakeType != ProtocolHandshakeTypeSelect ||
			len(handshake.Formats) != 1 || handshake.Formats[0] != ProtocolHandshakeFormatJSON {
			msg := CmiProtocolHandshakeError{
				Error: CmiProtocolHandshakeErrorUnexpectedMessage,
			}

			_ = c.writeJSON(CmiTypeControl, msg)
			err = errors.New("handshake: invalid response")

		}
	}

	return resp, err
}

func (c *Connection) clientProtocolHandshake() error {
	req := CmiHandshakeMsg{
		MessageProtocolHandshake: []MessageProtocolHandshake{
			{
				HandshakeType: ProtocolHandshakeTypeAnnounceMax,
				Version:       Version{Major: 1, Minor: 0},
				Formats:       []string{ProtocolHandshakeFormatJSON},
			},
		},
	}
	err := c.writeJSON(CmiTypeControl, req)

	// receive server selection
	var resp CmiHandshakeMsg
	if err == nil {
		resp, err = c.handshakeReceiveSelect()
	}

	// send selection back to server
	if err == nil {
		err = c.writeJSON(CmiTypeControl, resp)
	}

	return err
}

func (c *Connection) serverProtocolHandshake() error {
	var req CmiHandshakeMsg
	typ, err := c.readJSON(&req)

	if err == nil && typ != CmiTypeControl {
		err = fmt.Errorf("handshake: invalid type: %0x", typ)
	}

	if err == nil {
		if len(req.MessageProtocolHandshake) != 1 {
			return errors.New("handshake: invalid length")
		}

		handshake := req.MessageProtocolHandshake[0]

		if handshake.HandshakeType != ProtocolHandshakeTypeAnnounceMax ||
			len(handshake.Formats) != 1 || handshake.Formats[0] != ProtocolHandshakeFormatJSON {
			msg := CmiProtocolHandshakeError{
				Error: CmiProtocolHandshakeErrorUnexpectedMessage,
			}

			_ = c.writeJSON(CmiTypeControl, msg)
			err = errors.New("handshake: invalid response")
		} else {
			// send selection to client
			req.MessageProtocolHandshake[0].HandshakeType = ProtocolHandshakeTypeSelect
			err = c.writeJSON(CmiTypeControl, req)
		}
	}

	// receive selection back from client
	if err == nil {
		_, err = c.handshakeReceiveSelect()
	}

	return err
}
