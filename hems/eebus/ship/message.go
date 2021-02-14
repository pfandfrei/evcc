package ship

type CmiMessage struct {
	CmiHandshakeMsg
	CmiConnectionPinState
	CmiAccessMethodsRequest
	CmiAccessMethods
	CmiCloseMsg
}

func CmiDecode(msg CmiMessage) interface{} {
	switch {
	case len(msg.CmiHandshakeMsg.MessageProtocolHandshake) > 0:
		return msg.CmiHandshakeMsg

	case len(msg.CmiConnectionPinState.ConnectionPinState) > 0:
		return msg.CmiConnectionPinState

	case len(msg.CmiAccessMethodsRequest.AccessMethodsRequest) > 0:
		return msg.CmiAccessMethodsRequest

	case len(msg.CmiAccessMethods.AccessMethods) > 0:
		return msg.CmiAccessMethods

	case len(msg.CmiCloseMsg.ConnectionClose) > 0:
		return msg.CmiCloseMsg
	}

	return nil
}
