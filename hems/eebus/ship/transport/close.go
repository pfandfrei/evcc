package transport

import (
	"errors"
	"time"

	"github.com/andig/evcc/hems/eebus/ship/message"
)

// AcceptClose accepts connection close
func (c *Transport) AcceptClose() error {
	return c.WriteJSON(message.CmiTypeEnd, message.CmiCloseMsg{
		ConnectionClose: message.ConnectionClose{
			Phase: message.CmiClosePhaseConfirm,
		},
	})
}

// Close closes the connection
func (c *Transport) Close() error {
	err := c.WriteJSON(message.CmiTypeEnd, message.CmiCloseMsg{
		ConnectionClose: message.ConnectionClose{
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
