package ship

import (
	"bytes"
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

func (c *Server) log() Logger {
	if c.Log == nil {
		return &NopLogger{}
	}
	return c.Log
}

func (c *Server) init() error {
	init := []byte{CmiTypeInit, 0x00}

	// CMI_STATE_CLIENT_EVALUATE
	msg, err := c.readBinary()
	if err != nil {
		return err
	}

	if bytes.Compare(init, msg) != 0 {
		return fmt.Errorf("init: invalid response: %0 x", msg)
	}

	// CMI_STATE_CLIENT_SEND
	return c.writeBinary(init)
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

	switch hs := msg.(type) {
	case MessageProtocolHandshake:
		if hs.HandshakeType != ProtocolHandshakeTypeAnnounceMax || len(hs.Formats) != 1 || hs.Formats[0] != ProtocolHandshakeFormatJSON {
			msg := CmiProtocolHandshakeError{
				Error: CmiProtocolHandshakeErrorUnexpectedMessage,
			}

			_ = c.writeJSON(CmiTypeControl, msg)
			err = errors.New("handshake: invalid response")
			break
		}

		// send selection to client
		hs.HandshakeType = ProtocolHandshakeTypeSelect
		err = c.writeJSON(CmiTypeControl, CmiHandshakeMsg{
			MessageProtocolHandshake: []MessageProtocolHandshake{hs},
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
	c.Transport = &Transport{
		Conn:   conn,
		Log:    c.log(),
		inC:    make(chan []byte, 1),
		errC:   make(chan error, 1),
		closeC: make(chan struct{}, 1),
	}

	if err := c.init(); err != nil {
		return err
	}

	// start consuming messages
	go c.readPump()

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
