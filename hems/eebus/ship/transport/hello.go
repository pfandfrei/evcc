package transport

import (
	"errors"
	"fmt"
	"time"

	"github.com/andig/evcc/hems/eebus/ship/message"
)

// hello is the common hello exchange
func (c *Transport) Hello() error {
	if err := c.WriteJSON(message.CmiTypeControl, message.CmiHelloMsg{
		message.ConnectionHello{
			Phase: message.CmiHelloPhaseReady,
		},
	}); err != nil {
		return fmt.Errorf("hello: %w", err)
	}

	timer := time.NewTimer(message.CmiHelloInitTimeout)
	for {
		msg, err := c.ReadMessage(timer.C)
		if err != nil {
			if errors.Is(err, ErrTimeout) {
				_ = c.WriteJSON(message.CmiTypeControl, message.CmiHelloMsg{
					message.ConnectionHello{
						Phase: message.CmiHelloPhaseAborted,
					},
				})
			}

			return err
		}

		switch hello := msg.(type) {
		case message.ConnectionHello:
			switch hello.Phase {
			case message.CmiHelloPhaseReady:
				return nil

			case message.CmiHelloPhaseAborted:
				return errors.New("hello: aborted")

			case message.CmiHelloPhasePending:
				if hello.ProlongationRequest {
					timer = time.NewTimer(message.CmiHelloProlongationTimeout)
				}
			}

		case message.ConnectionClose:
			err = errors.New("hello: remote closed")

		default:
			return errors.New("hello: invalid type")
		}
	}
}
