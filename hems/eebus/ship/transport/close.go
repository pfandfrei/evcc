package transport

import (
	"errors"
	"time"

	"github.com/andig/evcc/hems/eebus/ship/message"
)

// AcceptClose accepts connection close
func (c *Transport) AcceptClose() error {
	err := c.WriteJSON(message.CmiTypeEnd, message.CmiConnectionClose{
		ConnectionClose: message.ConnectionClose{
			Phase: string(message.ConnectionClosePhaseTypeConfirm),
		},
	})

	// stop read/write pump
	close(c.closeC)

	return err
}

// Close closes the connection
func (c *Transport) Close() error {
	err := c.WriteJSON(message.CmiTypeEnd, message.CmiConnectionClose{
		ConnectionClose: message.ConnectionClose{
			Phase: string(message.ConnectionClosePhaseTypeAnnounce),
			// MaxTime: int(message.CmiCloseTimeout / time.Millisecond),
		},
	})

	timer := time.NewTimer(message.CmiCloseTimeout)
	for err == nil {
		msg, err := c.ReadMessage(timer.C)
		if err != nil {
			break
		}

		if typed, ok := msg.(message.ConnectionClose); ok && typed.Phase == string(message.ConnectionClosePhaseTypeConfirm) {
			break
		}

		err = errors.New("close: invalid response")
	}

	// stop read/write pump
	close(c.closeC)

	return err
}
