package ship

import (
	"errors"
	"fmt"
	"time"
)

// Close closes the service connection
func (c *Transport) acceptClose() error {
	msg := CmiCloseMsg{
		ConnectionClose: []ConnectionClose{
			{
				Phase: CmiClosePhaseConfirm,
			},
		},
	}
	return c.writeJSON(CmiTypeEnd, msg)
}

func (c *Transport) close() error {
	msg := CmiCloseMsg{
		ConnectionClose: []ConnectionClose{
			{
				Phase:   CmiClosePhaseAnnounce,
				MaxTime: int(CmiCloseTimeout / time.Millisecond),
			},
		},
	}
	if err := c.writeJSON(CmiTypeEnd, msg); err != nil {
		return err
	}

	timer := time.NewTimer(CmiCloseTimeout)
	for {
		select {
		case <-timer.C:
			return errors.New("close: timeout")

		default:
			var msg CmiCloseMsg
			typ, err := c.readJSONWithTimeout(CmiCloseTimeout, &msg)

			if err == nil && typ != CmiTypeEnd {
				err = fmt.Errorf("close: invalid type: %0x", typ)
			}

			if err == nil && len(msg.ConnectionClose) != 1 {
				err = errors.New("close: invalid length")
			}

			if err != nil {
				return err
			}

			close := msg.ConnectionClose[0]

			switch close.Phase {
			case CmiClosePhaseConfirm:
				return nil

			default:
				return errors.New("close: invalid response")
			}
		}
	}
}
