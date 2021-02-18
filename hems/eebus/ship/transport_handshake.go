package ship

import (
	"errors"
	"time"
)

func (c *Transport) handshakeReceiveSelect() error {
	timer := time.NewTimer(cmiReadWriteTimeout)
	msg, err := c.readMessage(timer.C)
	if err != nil {
		return err
	}

	switch typed := msg.(type) {
	case MessageProtocolHandshake:
		if typed.HandshakeType.HandshakeType != ProtocolHandshakeTypeSelect || !typed.Formats.Supports(ProtocolHandshakeFormatJSON) {
			_ = c.writeJSON(CmiTypeControl, CmiProtocolHandshakeError{
				Error: CmiProtocolHandshakeErrorUnexpectedMessage,
			})

			err = errors.New("handshake: invalid format")
		}

		return nil

	case ConnectionClose:
		err = errors.New("handshake: remote closed")

	default:
		err = errors.New("handshake: invalid type")
	}

	return err
}
