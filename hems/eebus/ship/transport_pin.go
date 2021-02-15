package ship

import (
	"errors"
)

// read pin requirements
func (c *Transport) readPinState() (ConnectionPinState, error) {
	var resp CmiConnectionPinState
	typ, err := c.readJSON(&resp)

	if err == nil && typ != CmiTypeControl {
		err = errors.New("pin: invalid type")
	}

	if err == nil && len(resp.ConnectionPinState) != 1 {
		err = errors.New("pin: invalid length")
	}

	var ps ConnectionPinState
	if err == nil {
		ps = resp.ConnectionPinState[0]
	}

	return ps, err
}
