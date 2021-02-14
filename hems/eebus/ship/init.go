package ship

import (
	"bytes"
	"errors"
	"fmt"
	"time"
)

const (
	CmiTypeInit byte = 0

	CmiHelloInitTimeout         = 60 * time.Second
	CmiHelloProlongationTimeout = 30 * time.Second

	CmiHelloPhasePending = "pending"
	CmiHelloPhaseReady   = "ready"
	CmiHelloPhaseAborted = "aborted"
)

type CmiHelloMsg struct {
	ConnectionHello []ConnectionHello `json:"connectionHello"`
}

type ConnectionHello struct {
	Phase               string `json:"phase"`
	Waiting             int    `json:"waiting,omitempty"`
	ProlongationRequest bool   `json:"prolongationRequest,omitempty"`
}

func (c *Connection) init() error {
	init := []byte{CmiTypeInit, CmiTypeInit}

	// CMI_STATE_CLIENT_SEND
	if err := c.writeBinary(init); err != nil {
		return err
	}

	// CMI_STATE_CLIENT_EVALUATE
	msg, err := c.readBinary()
	if err != nil {
		return err
	}

	if bytes.Compare(init, msg) != 0 {
		return fmt.Errorf("invalid init response: %0 x", msg)
	}

	return nil
}

func (c *Connection) hello() (err error) {
	// send ABORT if hello fails
	defer func() {
		if err != nil {
			// TODO
			_ = c.writeJSON(CmiTypeEnd, CmiHelloMsg{
				[]ConnectionHello{
					{Phase: CmiHelloPhaseAborted},
				},
			})
		}
	}()

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

			if err != nil {
				return err
			}

			hello := resp.ConnectionHello[0]

			switch hello.Phase {
			case "":
				return errors.New("hello: invalid response")

			case CmiHelloPhaseAborted:
				return errors.New("hello: aborted by peer")

			case CmiHelloPhaseReady:
				return nil

			case CmiHelloPhasePending:
				if hello.ProlongationRequest {
					timer = time.NewTimer(CmiHelloProlongationTimeout)
				}
			}
		}
	}
}
