package ship

import (
	"errors"
	"fmt"
	"time"
)

// hello is the common hello exchange
func (c *Transport) hello() error {
	req := CmiHelloMsg{
		[]ConnectionHello{
			{Phase: CmiHelloPhaseReady},
		},
	}
	if err := c.writeJSON(CmiTypeControl, req); err != nil {
		return err
	}

	timer := time.NewTimer(CmiHelloInitTimeout)
	for {
		select {
		case <-timer.C:
			req := CmiHelloMsg{
				[]ConnectionHello{
					{Phase: CmiHelloPhaseAborted},
				},
			}
			_ = c.writeJSON(CmiTypeControl, req)
			return errors.New("hello: timeout")

		default:
			var resp CmiHelloMsg
			typ, err := c.readJSON(&resp)

			if err == nil && typ != CmiTypeControl {
				err = fmt.Errorf("hello: invalid type: %0x", typ)
			}

			if err == nil && len(resp.ConnectionHello) != 1 {
				err = errors.New("hello: invalid length")
			}

			if err == nil {
				hello := resp.ConnectionHello[0]

				switch hello.Phase {
				case CmiHelloPhaseAborted:
					return errors.New("hello: aborted by peer")

				case CmiHelloPhaseReady:
					return nil

				case CmiHelloPhasePending:
					if hello.ProlongationRequest {
						timer = time.NewTimer(CmiHelloProlongationTimeout)
					}

				default:
					return errors.New("hello: invalid response")
				}
			}
		}
	}
}
