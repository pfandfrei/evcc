package ship

import (
	"errors"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

// Server is the SHIP server
type Server struct {
	Log                 Logger
	LocalPin, RemotePin string
	AccessMethods       []string
	*Transport
	Handler func(req interface{}) error
}

func (c *Server) protocolHandshake() error {
	timer := time.NewTimer(cmiReadWriteTimeout)
	msg, err := c.readMessage(timer.C)
	if err != nil {
		if errors.Is(err, ErrTimeout) {
			_ = c.writeJSON(CmiTypeControl, CmiProtocolHandshakeError{
				Error: CmiProtocolHandshakeErrorUnexpectedMessage,
			})
		}

		return err
	}

	switch typed := msg.(type) {
	case MessageProtocolHandshake:
		if typed.HandshakeType.HandshakeType != ProtocolHandshakeTypeAnnounceMax || !typed.Formats.Supports(ProtocolHandshakeFormatJSON) {
			msg := CmiProtocolHandshakeError{
				Error: CmiProtocolHandshakeErrorUnexpectedMessage,
			}

			_ = c.writeJSON(CmiTypeControl, msg)
			err = errors.New("handshake: invalid response")
			break
		}

		// send selection to client
		typed.HandshakeType.HandshakeType = ProtocolHandshakeTypeSelect
		err = c.writeJSON(CmiTypeControl, CmiHandshakeMsg{
			MessageProtocolHandshake: []MessageProtocolHandshake{typed},
		})

	default:
		return fmt.Errorf("handshake: invalid type")
	}

	// receive selection back from client
	if err == nil {
		err = c.handshakeReceiveSelect()
	}

	return err
}

// Close performs ordered close of server connection
func (c *Server) Close() error {
	return c.close()
}

// Serve performs the server connection handshake
func (c *Server) Serve(conn *websocket.Conn) error {
	c.Transport = NewTransport(c.Log, conn)

	if err := c.init(); err != nil {
		return err
	}

	err := c.hello()
	if err == nil {
		err = c.protocolHandshake()
	}
	if err == nil {
		err = c.pinState(c.LocalPin, c.RemotePin)
	}
	if err == nil {
		err = c.accessMethods(c.AccessMethods)
	}

	for err == nil {
		endless := make(chan time.Time)

		var msg interface{}
		msg, err = c.readMessage(endless)
		if err != nil {
			break
		}

		switch typed := msg.(type) {
		case ConnectionClose:
			return c.acceptClose()

		case Data:
			c.log().Printf("serv: %+v", msg)
			if c.Handler == nil {
				err = errors.New("no handler")
				break
			}

			if err = c.Handler(typed); err != nil {
				break
			}

		default:
			err = errors.New("invalid type")
		}
	}

	// close connection if handshake or hello fails
	if err != nil {
		_ = c.Close()
	}

	return err
}
