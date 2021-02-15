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
	mux sync.Mutex
	Log Logger
	Pin string
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

func (c *Client) pinState() error {
	req := CmiConnectionPinState{
		ConnectionPinState: []ConnectionPinState{
			{
				PinState: PinStateNone,
			},
		},
	}
	if err := c.writeJSON(CmiTypeControl, req); err != nil {
		return err
	}

	ps, err := c.readPinState()

	if err == nil {
		// ps := resp.ConnectionPinState[0]

		if ps.PinState == PinStateRequired || (ps.PinState == PinStateOptional && c.Pin != "") {
			req := CmiConnectionPinInput{
				ConnectionPinInput: []ConnectionPinInput{
					{
						Pin: c.Pin,
					},
				},
			}
			err = c.writeJSON(CmiTypeControl, req)
		}
	}

	// TODO check if next message is pin error

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

func (c *Client) pump() {
	for {
		var err error
		select {
		case req := <-c.send:
			err = c.writeJSON(CmiTypeData, req)
		}

		if err != nil {
			c.log().Println(err)
			break
		}
	}
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
		err = c.pinState()
	}
	if err == nil {
		err = c.accessMethodsRequest()
	}
	if err == nil {
		err = c.accessMethods()
	}

	// close connection if handshake or hello fails
	if err != nil {
		_ = c.Close()
	}

	return err
}

func (c *Client) Write(req interface{}) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.closed {
		return os.ErrClosed
	}

	return c.writeJSON(CmiTypeData, req)
}

func (c *Client) Read(res interface{}) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.closed {
		return os.ErrClosed
	}

	typ, err := c.readJSON(&res)
	if err == nil && typ != CmiTypeData {
		err = fmt.Errorf("read: invalid type: %0x", typ)
	}

	return err
}
