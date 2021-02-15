package ship

import (
	"errors"
	"fmt"
)

func (c *Transport) handshakeReceiveSelect() (CmiHandshakeMsg, error) {
	var resp CmiHandshakeMsg
	typ, err := c.readJSON(&resp)

	if err == nil && typ != CmiTypeControl {
		err = fmt.Errorf("handshake: invalid type: %0x", typ)
	}

	if err == nil && len(resp.MessageProtocolHandshake) != 1 {
		return resp, errors.New("handshake: invalid length")
	}

	hs := resp.MessageProtocolHandshake[0]

	if hs.HandshakeType != ProtocolHandshakeTypeSelect || len(hs.Formats) != 1 || hs.Formats[0] != ProtocolHandshakeFormatJSON {
		msg := CmiProtocolHandshakeError{
			Error: CmiProtocolHandshakeErrorUnexpectedMessage,
		}

		_ = c.writeJSON(CmiTypeControl, msg)
		err = errors.New("handshake: invalid response")
	}

	return resp, err
}
