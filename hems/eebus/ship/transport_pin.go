package ship

import (
	"errors"
	"time"
)

// read pin requirements
func (c *Transport) readPinState() (ConnectionPinState, error) {
	timer := time.NewTimer(cmiReadWriteTimeout)
	msg, err := c.readMessage(timer.C)

	switch typed := msg.(type) {
	case ConnectionPinState:
		return typed, err

	default:
		if err == nil {
			err = errors.New("pin: invalid type")
		}

		return ConnectionPinState{}, err
	}
}

const (
	pinReceived = 1 << iota
	pinSent

	pinCompleted = pinReceived | pinSent
)

func (c *Transport) pinState(local, remote string) error {
	pinState := ConnectionPinState{
		PinState: PinStateNone,
	}

	var status int
	if local != "" {
		pinState.PinState = PinStateRequired
		pinState.InputPermission = PinInputPermissionOk
	} else {
		// always received if not necessary
		status |= pinReceived
	}

	err := c.writeJSON(CmiTypeControl, CmiConnectionPinState{
		ConnectionPinState: []ConnectionPinState{pinState},
	})

	timer := time.NewTimer(10 * time.Second)
	for err == nil && status != pinCompleted {
		var msg interface{}
		msg, err = c.readMessage(timer.C)
		if err != nil {
			break
		}

		switch typed := msg.(type) {
		// local pin
		case ConnectionPinInput:
			// signal error to client
			if typed.Pin != local {
				err = c.writeJSON(CmiTypeControl, CmiConnectionPinError{
					ConnectionPinError: []ConnectionPinError{
						{Error: 1},
					},
				})
			}

			status |= pinReceived

		// remote pin
		case ConnectionPinState:
			if typed.PinState == PinStateOptional || typed.PinState == PinStateRequired {
				if remote != "" {
					err = c.writeJSON(CmiTypeControl, CmiConnectionPinInput{
						ConnectionPinInput: []ConnectionPinInput{
							{Pin: remote},
						},
					})
				} else {
					err = errors.New("pin: remote pin required")
				}
			}

			status |= pinSent

		case ConnectionPinError:
			err = errors.New("pin: remote pin mismatched")

		case ConnectionClose:
			err = errors.New("pin: remote closed")

		default:
			err = errors.New("pin: invalid type")
		}
	}

	return err
}
