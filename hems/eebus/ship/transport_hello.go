package ship

import (
	"errors"
	"fmt"
	"time"
)

// hello is the common hello exchange
func (c *Transport) hello() error {
	if err := c.writeJSON(CmiTypeControl, CmiHelloMsg{
		[]ConnectionHello{
			{Phase: CmiHelloPhaseReady},
		},
	}); err != nil {
		return fmt.Errorf("hello: %w", err)
	}

	timer := time.NewTimer(CmiHelloInitTimeout)
	for {
		msg, err := c.readMessage(timer.C)
		if err != nil {
			if errors.Is(err, ErrTimeout) {
				_ = c.writeJSON(CmiTypeControl, CmiHelloMsg{
					[]ConnectionHello{
						{Phase: CmiHelloPhaseAborted},
					},
				})
			}

			return err
		}

		switch hello := msg.(type) {
		case ConnectionHello:
			switch hello.Phase {
			case CmiHelloPhaseReady:
				return nil

			case CmiHelloPhaseAborted:
				return errors.New("hello: aborted")

			case CmiHelloPhasePending:
				if hello.ProlongationRequest {
					timer = time.NewTimer(CmiHelloProlongationTimeout)
				}
			}

		default:
			return errors.New("hello: invalid type")
		}
	}
}
