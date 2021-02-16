package ship

import (
	"errors"
	"time"
)

// Close closes the service connection
func (c *Transport) acceptClose() error {
	return c.writeJSON(CmiTypeEnd, CmiCloseMsg{
		ConnectionClose: []ConnectionClose{
			{Phase: CmiClosePhaseConfirm},
		},
	})
}

func (c *Transport) close() error {
	err := c.writeJSON(CmiTypeEnd, CmiCloseMsg{
		ConnectionClose: []ConnectionClose{
			{
				Phase:   CmiClosePhaseAnnounce,
				MaxTime: int(CmiCloseTimeout / time.Millisecond),
			},
		},
	})

	timer := time.NewTimer(CmiCloseTimeout)
	for err == nil {
		msg, err := c.readMessage(timer.C)
		if err != nil {
			break
		}

		if typed, ok := msg.(ConnectionClose); ok && typed.Phase == CmiClosePhaseConfirm {
			return nil
		}

		err = errors.New("close: invalid response")
	}

	return err
}
