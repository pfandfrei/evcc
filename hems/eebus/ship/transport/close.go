package transport

import (
	"errors"
	"time"

	"github.com/andig/evcc/hems/eebus/ship/message"
)

// Close closes the service connection
func (c *Transport) AcceptClose() error {
	return c.WriteJSON(message.CmiTypeEnd, message.CmiCloseMsg{
		message.ConnectionClose{
			Phase: message.CmiClosePhaseConfirm,
		},
	})
}

// Close closes the service connection
func (c *Transport) Close() error {
	err := c.WriteJSON(message.CmiTypeEnd, message.CmiCloseMsg{
		message.ConnectionClose{
			Phase:   message.CmiClosePhaseAnnounce,
			MaxTime: int(message.CmiCloseTimeout / time.Millisecond),
		},
	})

	timer := time.NewTimer(message.CmiCloseTimeout)
	for err == nil {
		msg, err := c.ReadMessage(timer.C)
		if err != nil {
			break
		}

		if typed, ok := msg.(message.ConnectionClose); ok && typed.Phase == message.CmiClosePhaseConfirm {
			return nil
		}

		err = errors.New("close: invalid response")
	}

	return err
}
