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
		if typed.HandshakeType.HandshakeType != ProtocolHandshakeTypeSelect ||
			len(typed.Formats) != 1 || typed.Formats[0].Format != ProtocolHandshakeFormatJSON {
			_ = c.writeJSON(CmiTypeControl, CmiProtocolHandshakeError{
				Error: CmiProtocolHandshakeErrorUnexpectedMessage,
			})

			return errors.New("handshake: invalid format")
		}

		return nil

	default:
		return errors.New("handshake: invalid type")
	}
}
