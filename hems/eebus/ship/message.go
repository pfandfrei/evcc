package ship

import "errors"

// CmiMessage is used to identify the defined message types
type CmiMessage struct {
	*CmiHelloMsg
	*CmiHandshakeMsg
	*CmiConnectionPinState
	*CmiConnectionPinInput
	*CmiConnectionPinError
	*CmiAccessMethodsRequest
	*CmiAccessMethods
	*CmiDataMsg
	*CmiCloseMsg
}

// decodeCmi extract ship message core
func decodeCmi(msg CmiMessage) (res interface{}, err error) {
	switch {
	case msg.CmiHelloMsg != nil:
		res = ConnectionHello{}
		if len(msg.CmiHelloMsg.ConnectionHello) == 1 {
			res = msg.CmiHelloMsg.ConnectionHello[0]
		}

	case msg.CmiHandshakeMsg != nil:
		res = MessageProtocolHandshake{}
		if len(msg.CmiHandshakeMsg.MessageProtocolHandshake) == 1 {
			res = msg.CmiHandshakeMsg.MessageProtocolHandshake[0]
		}

	case msg.CmiConnectionPinState != nil:
		res = ConnectionPinState{}
		if len(msg.CmiConnectionPinState.ConnectionPinState) == 1 {
			res = msg.CmiConnectionPinState.ConnectionPinState[0]
		}

	case msg.CmiConnectionPinInput != nil:
		res = ConnectionPinInput{}
		if len(msg.CmiConnectionPinInput.ConnectionPinInput) == 1 {
			res = msg.CmiConnectionPinInput.ConnectionPinInput[0]
		}

	case msg.CmiConnectionPinError != nil:
		res = ConnectionPinError{}
		if len(msg.CmiConnectionPinError.ConnectionPinError) == 1 {
			res = msg.CmiConnectionPinError.ConnectionPinError[0]
		}

	case msg.CmiAccessMethodsRequest != nil:
		res = AccessMethodsRequest{}

	case msg.CmiAccessMethods != nil:
		res = AccessMethods{}
		if len(msg.CmiAccessMethods.AccessMethods) == 1 {
			res = msg.CmiAccessMethods.AccessMethods[0]
		}

	case msg.CmiDataMsg != nil:
		res = Data{}
		if len(msg.CmiDataMsg.Data) == 1 {
			res = msg.CmiDataMsg.Data[0]
		}

	case msg.CmiCloseMsg != nil:
		res = ConnectionClose{}
		if len(msg.CmiCloseMsg.ConnectionClose) == 1 {
			res = msg.CmiCloseMsg.ConnectionClose[0]
		}

	default:
		err = errors.New("invalid type")
	}

	return
}
