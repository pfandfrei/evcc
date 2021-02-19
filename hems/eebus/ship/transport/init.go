package transport

import (
	"bytes"
	"fmt"
	"time"

	"github.com/andig/evcc/hems/eebus/ship/message"
)

func (c *Transport) Init() error {
	init := []byte{message.CmiTypeInit, 0x00}

	// CMI_STATE_CLIENT_SEND
	if err := c.WriteBinary(init); err != nil {
		return err
	}

	timer := time.NewTimer(message.CmiHelloInitTimeout)

	// CMI_STATE_CLIENT_EVALUATE
	msg, err := c.ReadBinary(timer.C)
	if err != nil {
		return err
	}

	if bytes.Compare(init, msg) != 0 {
		return fmt.Errorf("init: invalid response")
	}

	return nil
}
