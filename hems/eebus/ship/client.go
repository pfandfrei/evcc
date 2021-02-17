package ship

import (
	"fmt"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

// Client is the ship client
type Client struct {
	mux                 sync.Mutex
	Log                 Logger
	LocalPin, RemotePin string
	AccessMethods       []string
	*Transport
	closed bool
}

func (c *Client) protocolHandshake() error {
	hs := CmiHandshakeMsg{
		MessageProtocolHandshake: []MessageProtocolHandshake{
			{
				HandshakeType: HandshakeType{ProtocolHandshakeTypeAnnounceMax},
				Version:       Version{Major: 1, Minor: 0},
				Formats:       []Format{{[]string{ProtocolHandshakeFormatJSON}}},
			},
		},
	}
	if err := c.writeJSON(CmiTypeControl, hs); err != nil {
		return fmt.Errorf("handshake: %w", err)
	}

	// receive server selection and send selection back to server
	err := c.handshakeReceiveSelect()
	if err == nil {
		hs.MessageProtocolHandshake[0].HandshakeType.HandshakeType = ProtocolHandshakeTypeSelect
		err = c.writeJSON(CmiTypeControl, hs)
	}

	return err
}

// Close performs ordered close of client connection
func (c *Client) Close() error {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.closed {
		return os.ErrClosed
	}

	c.closed = true

	// stop readPump
	defer close(c.closeC)

	return c.close()
}

// Connect performs the client connection handshake
func (c *Client) Connect(conn *websocket.Conn) error {
	c.Transport = NewTransport(c.Log, conn)

	if err := c.init(); err != nil {
		return err
	}

	// // start consuming messages
	// go c.readPump()

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

	// close connection if handshake or hello fails
	if err != nil {
		_ = c.Close()
	}

	return err
}

// func (c *Client) Write(req interface{}) error {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()

// 	if c.closed {
// 		return os.ErrClosed
// 	}

// 	return c.writeJSON(CmiTypeData, req)
// }

// func (c *Client) Read(res interface{}) error {
// 	c.mux.Lock()
// 	defer c.mux.Unlock()

// 	if c.closed {
// 		return os.ErrClosed
// 	}

// 	typ, err := c.readJSON(&res)
// 	if err == nil && typ != CmiTypeData {
// 		err = fmt.Errorf("read: invalid type: %0x", typ)
// 	}

// 	return err
// }
