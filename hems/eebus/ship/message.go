package ship

type CmiMessage struct {
	CmiHandshakeMsg
	CmiConnectionPinState
	CmiAccessMethodsRequest
	CmiAccessMethods
	CmiDataMsg
	CmiCloseMsg
}

func CmiDecode(msg CmiMessage) interface{} {
	switch {
	case len(msg.CmiHandshakeMsg.MessageProtocolHandshake) == 1:
		return msg.CmiHandshakeMsg.MessageProtocolHandshake[0]

	case len(msg.CmiConnectionPinState.ConnectionPinState) == 1:
		return msg.CmiConnectionPinState.ConnectionPinState[0]

	case len(msg.CmiAccessMethodsRequest.AccessMethodsRequest) == 1:
		return msg.CmiAccessMethodsRequest.AccessMethodsRequest[0]

	case len(msg.CmiAccessMethods.AccessMethods) == 1:
		return msg.CmiAccessMethods.AccessMethods[0]

	case len(msg.CmiDataMsg.Data) == 1:
		return msg.CmiDataMsg.Data[0]

	case len(msg.CmiCloseMsg.ConnectionClose) == 1:
		return msg.CmiCloseMsg.ConnectionClose[0]
	}

	return nil
}
