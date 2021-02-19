package transport

import (
	"errors"
	"time"

	"github.com/andig/evcc/hems/eebus/ship/message"
)

// accessMethodsRequest
func (c *Transport) AccessMethodsRequest(methods []string) ([]string, error) {
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
			return []string{typed.ID}, nil

		case message.AccessMethodsRequest:
			am := make([]message.AccessMethods, 0, len(methods))
			for _, m := range methods {
				am = append(am, message.AccessMethods{ID: m})
			}
			err = c.WriteJSON(message.CmiTypeControl, message.CmiAccessMethods{
				AccessMethods: am,
			})

		default:
			err = errors.New("access methods: invalid type")
		}
	}

	return nil, err
}
