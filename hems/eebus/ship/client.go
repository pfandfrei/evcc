package ship

import (
	"bytes"
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
	send   <-chan interface{}
}

func (c *Client) log() Logger {
	if c.Log == nil {
		return &NopLogger{}
	}
	return c.Log
}

func (c *Client) init() error {
	init := []byte{CmiTypeInit, 0x00}

	// CMI_STATE_CLIENT_SEND
	if err := c.writeBinary(init); err != nil {
		return err
	}

	// CMI_STATE_CLIENT_EVALUATE
	msg, err := c.readBinary()
	if err != nil {
		return err
	}

	if bytes.Compare(init, msg) != 0 {
		return fmt.Errorf("init: invalid response: %0 x", msg)
	}

	// move to control phase
	c.phase = CmiTypeControl

	return nil
}

func (c *Client) protocolHandshake() error {
	hs := CmiHandshakeMsg{
		MessageProtocolHandshake: []MessageProtocolHandshake{
			{
				HandshakeType: ProtocolHandshakeTypeAnnounceMax,
				Version:       []Version{{Major: 1, Minor: 0}},
				Formats:       []Format{{Format: ProtocolHandshakeFormatJSON}},
			},
		},
	}
	if err := c.writeJSON(CmiTypeControl, hs); err != nil {
		return fmt.Errorf("handshake: %w", err)
	}

	// receive server selection and send selection back to server
	err := c.handshakeReceiveSelect()
	if err == nil {
		hs.MessageProtocolHandshake[0].HandshakeType = ProtocolHandshakeTypeSelect
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

	return c.close()
}

// Connect performs the client connection handshake
func (c *Client) Connect(conn *websocket.Conn) error {
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
