package ship

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
)

// Server is the SHIP server
type Server struct {
	Log Logger
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
	var req CmiHandshakeMsg
	typ, err := c.readJSON(&req)

	if err == nil && typ != CmiTypeControl {
		err = fmt.Errorf("handshake: invalid type: %0x", typ)
	}

	if err == nil && len(req.MessageProtocolHandshake) != 1 {
		return errors.New("handshake: invalid length")
	}

	hs := req.MessageProtocolHandshake[0]

	if hs.HandshakeType != ProtocolHandshakeTypeAnnounceMax || len(hs.Formats) != 1 || hs.Formats[0] != ProtocolHandshakeFormatJSON {
		msg := CmiProtocolHandshakeError{
			Error: CmiProtocolHandshakeErrorUnexpectedMessage,
		}

		_ = c.writeJSON(CmiTypeControl, msg)
		return errors.New("handshake: invalid response")
	}

	// send selection to client
	req.MessageProtocolHandshake[0].HandshakeType = ProtocolHandshakeTypeSelect
	err = c.writeJSON(CmiTypeControl, req)

	// receive selection back from client
	if err == nil {
		_, err = c.handshakeReceiveSelect()
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
		Conn: conn,
		Log:  c.log(),
	}

	err := c.init()
	if err == nil {
		err = c.hello()
	}
	if err == nil {
		err = c.protocolHandshake()
	}

	if err == nil {
		for {
			var typ byte
			var req CmiMessage
			typ, err = c.waitJSON(&req)
			if err != nil {
				break
			}

			var typed interface{}
			typed, err = DecodeMessage(req)

			c.log().Printf("serv: %d %+v", typ, typed)

			if err != nil {
				break
			}

			if _, ok := typed.(ConnectionClose); ok {
				return c.acceptClose()
			}

			if c.Handler == nil {
				err = errors.New("no handler")
				break
			}

			if err = c.Handler(typed); err != nil {
				break
			}
		}
	}

	// close connection if handshake or hello fails
	if err != nil {
		_ = c.Close()
	}

	return err
}
