package transport

import (
	"errors"
	"time"

	"github.com/andig/evcc/hems/eebus/ship/message"
)

// AccessMethodsRequest
func (c *Transport) AccessMethodsRequest(methods string) (string, error) {
	err := c.WriteJSON(message.CmiTypeControl, message.CmiAccessMethodsRequest{
		AccessMethodsRequest: []message.AccessMethodsRequest{},
	})

	for err == nil {
		timer := time.NewTimer(CmiReadWriteTimeout)
		msg, err := c.ReadMessage(timer.C)
		if err != nil {
			break
		}

		switch typed := msg.(type) {
		case message.AccessMethods:
			// access methods received
			return typed.ID, nil

		case message.AccessMethodsRequest:
			err = c.WriteJSON(message.CmiTypeControl, message.CmiAccessMethods{
				// AccessMethods: message.AccessMethods{ID: methods},
			})

		default:
			err = errors.New("access methods: invalid type")
		}
	}

	return "", err
}
