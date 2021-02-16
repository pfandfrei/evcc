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

func (c *Transport) pinState(local, remote string) error {
	pinState := PinStateNone
	var inputPermission string
	if local != "" {
		pinState = PinStateRequired
		inputPermission = PinInputPermissionOk
	}

	err := c.writeJSON(CmiTypeControl, CmiConnectionPinState{
		ConnectionPinState: []ConnectionPinState{
			{
				PinState:        pinState,
				InputPermission: inputPermission,
			},
		},
	})

	var pinEntered string
	timer := time.NewTimer(10 * time.Second)

	for err == nil && local != pinEntered {
		msg, err := c.readMessage(timer.C)
		if err != nil {
			break
		}

		switch typed := msg.(type) {
		// local pin
		case ConnectionPinInput:
			pinEntered = typed.Pin

			// signal error to client
			if typed.Pin != local {
				err = c.writeJSON(CmiTypeControl, CmiConnectionPinError{
					ConnectionPinError: []ConnectionPinError{
						{Error: 1},
					},
				})
			}

		// remote pin
		// case ConnectionPinState:

		case ConnectionPinError:
			err = errors.New("pin: remote pin mismatched")

		default:
			return errors.New("pin: invalid type")
		}
	}

	return err
}
