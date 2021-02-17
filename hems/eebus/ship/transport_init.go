package ship

import (
	"bytes"
	"fmt"
	"time"
)

func (c *Transport) init() error {
	init := []byte{CmiTypeInit, 0x00}

	// CMI_STATE_CLIENT_SEND
	if err := c.writeBinary(init); err != nil {
		return err
	}

	timer := time.NewTimer(CmiHelloInitTimeout)

	// CMI_STATE_CLIENT_EVALUATE
	msg, err := c.readBinary(timer.C)
	if err != nil {
		return err
	}

	if bytes.Compare(init, msg) != 0 {
		return fmt.Errorf("init: invalid response")
	}

	return nil
}
