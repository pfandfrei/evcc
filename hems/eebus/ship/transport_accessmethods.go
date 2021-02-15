package ship

import (
	"errors"
)

// accessMethodsRequest
func (c *Transport) accessMethodsRequest() error {
	req := CmiAccessMethodsRequest{
		AccessMethodsRequest: []AccessMethodsRequest{},
	}
	if err := c.writeJSON(CmiTypeControl, req); err != nil {
		return err
	}

	var resp CmiAccessMethodsRequest
	typ, err := c.readJSON(&resp)

	if err == nil && typ != CmiTypeControl {
		err = errors.New("access methods request: invalid type")
	}

	return err
}

// accessMethods
func (c *Transport) accessMethods() error {
	req := CmiAccessMethods{
		AccessMethods: []AccessMethods{
			{
				ID: "foo-bar",
			},
		},
	}
	if err := c.writeJSON(CmiTypeControl, req); err != nil {
		return err
	}

	var resp CmiAccessMethods
	typ, err := c.readJSON(&resp)

	if err == nil && typ != CmiTypeControl {
		err = errors.New("access methods: invalid type")
	}

	return err
}
