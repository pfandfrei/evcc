package ship

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

// Client is the SHIP client
type Client struct {
	mux sync.Mutex
	Log Logger
	*Transport
	closed bool
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

	return nil
}

func (c *Client) protocolHandshake() error {
	req := CmiHandshakeMsg{
		MessageProtocolHandshake: []MessageProtocolHandshake{
			{
				HandshakeType: ProtocolHandshakeTypeAnnounceMax,
				Version:       Version{Major: 1, Minor: 0},
				Formats:       []string{ProtocolHandshakeFormatJSON},
			},
		},
	}
	err := c.writeJSON(CmiTypeControl, req)

	// receive server selection
	var resp CmiHandshakeMsg
	if err == nil {
		resp, err = c.handshakeReceiveSelect()
	}

	// send selection back to server
	if err == nil {
		err = c.writeJSON(CmiTypeControl, resp)
	}

	return err
}

// Close performs ordered close of client connection
func (c *Client) Close() error {
	return c.close()
}

// Connect performs the client connection handshake
func (c *Client) Connect(conn *websocket.Conn) error {
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
		_ = c.Close()
	}

	// close connection if handshake or hello fails
	if err != nil {
		_ = c.Close()
	}

	return err
}
