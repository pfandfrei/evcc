package ship

import (
	"errors"
	"time"
)

// accessMethodsRequest
func (c *Transport) accessMethods(methods []string) error {
	err := c.writeJSON(CmiTypeControl, CmiAccessMethodsRequest{
		AccessMethodsRequest: []AccessMethodsRequest{},
	})

	for err == nil {
		timer := time.NewTimer(cmiReadWriteTimeout)
		msg, err := c.readMessage(timer.C)

		if err != nil {
			break
		}

		switch msg.(type) {
		case AccessMethods:
			// access methods received
			return nil

		case AccessMethodsRequest:
			am := make([]AccessMethods, 0, len(methods))
			for _, m := range methods {
				am = append(am, AccessMethods{ID: m})
			}
			err = c.writeJSON(CmiTypeControl, CmiAccessMethods{
				AccessMethods: am,
			})

		default:
			err = errors.New("access methods: invalid type")
		}
	}

	return err
}
