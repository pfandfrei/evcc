package ship

import (
	"errors"
	"fmt"
	"time"
)

const (
	CmiTypeEnd      byte = 3
	CmiCloseTimeout      = 100 * time.Millisecond
)

type CmiCloseMsg struct {
	ConnectionClose []ConnectionClose `json:"connectionClose"`
}

const (
	ConnectionCloseReasonUnspecific        = "unspecific"
	ConnectionCloseReasonRemovedConnection = "removedConnection"

	CmiClosePhaseAnnounce = "announce"
	CmiClosePhaseConfirm  = "confirm"
)

type ConnectionClose struct {
	Phase   string `json:"phase"`
	MaxTime int    `json:"maxTime,omitempty"`
	Reason  string `json:"reason,omitempty"`
}

func (c *Connection) closePassive() error {
	msg := CmiCloseMsg{
		ConnectionClose: []ConnectionClose{
			{
				Phase: CmiClosePhaseConfirm,
			},
		},
	}

	return c.writeJSON(CmiTypeEnd, msg)
}

func (c *Connection) closeActive() error {
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
			typ, err := c.readJSON(&msg)

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
