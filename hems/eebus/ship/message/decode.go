package message

import (
	"encoding/json"
	"errors"
)

func Decode(b []byte) (interface{}, error) {
	var sum map[string]json.RawMessage

	if err := json.Unmarshal(b, &sum); err != nil {
		return nil, err
	}

	var typ string
	var raw json.RawMessage
	for k, v := range sum {
		typ = k
		raw = v
	}

	switch typ {
	case "accessMethods":
		res := []AccessMethods{}
		err := json.Unmarshal(raw, &res)
		if len(res) > 0 {
			return res[0], err
		}
		return AccessMethods{}, nil

	case "accessMethodsRequest":
		res := []AccessMethodsRequest{}
		err := json.Unmarshal(raw, &res)
		if len(res) > 0 {
			return res[0], err
		}
		return AccessMethodsRequest{}, nil

	case "connectionPinState":
		res := []ConnectionPinState{}
		err := json.Unmarshal(raw, &res)
		if len(res) > 0 {
			return res[0], err
		}
		return ConnectionPinState{}, nil

	case "connectionPinInput":
		res := []ConnectionPinInput{}
		err := json.Unmarshal(raw, &res)
		if len(res) > 0 {
			return res[0], err
		}
		return ConnectionPinInput{}, nil

	case "connectionPinError":
		res := []ConnectionPinError{}
		err := json.Unmarshal(raw, &res)
		if len(res) > 0 {
			return res[0], err
		}
		return ConnectionPinError{}, nil

	case "connectionHello":
		res := []ConnectionHello{}
		err := json.Unmarshal(raw, &res)
		if len(res) > 0 {
			return res[0], err
		}
		return ConnectionHello{}, nil

	case "connectionClose":
		res := []ConnectionClose{}
		err := json.Unmarshal(raw, &res)
		if len(res) > 0 {
			return res[0], err
		}
		return ConnectionClose{}, nil

	case "messageProtocolHandshake":
		res := MessageProtocolHandshake{}
		err := json.Unmarshal(raw, &res)
		return res, err

	default:
		return nil, errors.New("invalid type")
	}
}
