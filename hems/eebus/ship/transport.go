package ship

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

const cmiReadWriteTimeout = 10 * time.Second

// ErrTimeout is the timeout error
var ErrTimeout = errors.New("timeout")

// Transport is the physical transport layer
type Transport struct {
	conn   *websocket.Conn
	logger Logger

	recv    chan []byte
	recvErr chan error
	send    chan []byte
	sendErr chan error
	closeC  chan struct{}
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

func NewTransport(log Logger, conn *websocket.Conn) *Transport {
	t := &Transport{
		conn:    conn,
		logger:  log,
		send:    make(chan []byte, 1),
		recv:    make(chan []byte, 1),
		sendErr: make(chan error, 1),
		recvErr: make(chan error, 1),
		closeC:  make(chan struct{}),
	}

	go t.readPump()
	go t.writePump()

	return t
}

func (c *Transport) log() Logger {
	if c.logger == nil {
		return &NopLogger{}
	}
	return c.logger
}

func (c *Transport) readPump() {
	defer func() {
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		select {
		case <-c.closeC:
			return

		default:
			typ, b, err := c.conn.ReadMessage()
			if err == nil {
				if len(b) > 2 {
					c.log().Println("recv:", string(b))
				}

				if typ != websocket.BinaryMessage {
					err = fmt.Errorf("invalid message type: %d", typ)
				}
			}

			if err == nil {
				c.recv <- b
			} else {
				c.recvErr <- err
			}
		}
	}
}

func (c *Transport) readBinary(timerC <-chan time.Time) ([]byte, error) {
	select {
	case <-timerC:
		return nil, ErrTimeout

	case <-c.closeC:
		return nil, net.ErrClosed

	case b := <-c.recv:
		return b, nil

	case err := <-c.recvErr:
		return nil, err
	}
}

func (c *Transport) readMessage(timerC <-chan time.Time) (interface{}, error) {
	select {
	case <-timerC:
		return nil, ErrTimeout

	case <-c.closeC:
		return nil, net.ErrClosed

	case b := <-c.recv:
		if len(b) < 2 {
			return nil, errors.New("invalid length")
		}
		if b[0] < 1 {
			return nil, errors.New("invalid phase")
		}

		return decodeMessage(b[1:])

	case err := <-c.recvErr:
		return nil, err
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Transport) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteMessage(websocket.BinaryMessage, msg)
			c.sendErr <- err

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Transport) writeBinary(msg []byte) error {
	c.send <- msg

	timer := time.NewTimer(10 * time.Second)
	select {
	case <-timer.C:
		return ErrTimeout

	case <-c.closeC:
		return net.ErrClosed

	case err := <-c.sendErr:
		return err
	}
}

func (c *Transport) writeJSON(typ byte, jsonMsg interface{}) error {
	msg, err := json.Marshal(jsonMsg)
	if err != nil {
		return err
	}

	// add header
	b := bytes.NewBuffer([]byte{typ})
	if _, err = b.Write(msg); err == nil {
		err = c.writeBinary(b.Bytes())
	}

	return err
}

func decodeMessage(b []byte) (interface{}, error) {
	var sum map[string]json.RawMessage

	// fmt.Println(string(b))
	if err := json.Unmarshal(b, &sum); err != nil {
		return nil, err
	}

	var typ string
	var raw json.RawMessage
	for k, v := range sum {
		typ = k
		raw = v
	}

	// fmt.Println(typ, sum)

	switch typ {
	case "accessMethods":
		res := []AccessMethods{}
		err := json.Unmarshal(raw, &res)
		if len(res) > 0 {
			return res[0], err
		}
		return AccessMethods{}, nil

	case "accessMethodsRequest":
		res := []AccessMethodsRequest{}
		err := json.Unmarshal(raw, &res)
		if len(res) > 0 {
			return res[0], err
		}
		return AccessMethodsRequest{}, nil

	case "connectionPinState":
		res := []ConnectionPinState{}
		err := json.Unmarshal(raw, &res)
		if len(res) > 0 {
			return res[0], err
		}
		return ConnectionPinState{}, nil

	case "connectionPinInput":
		res := []ConnectionPinInput{}
		err := json.Unmarshal(raw, &res)
		if len(res) > 0 {
			return res[0], err
		}
		return ConnectionPinInput{}, nil

	case "connectionPinError":
		res := []ConnectionPinError{}
		err := json.Unmarshal(raw, &res)
		if len(res) > 0 {
			return res[0], err
		}
		return ConnectionPinError{}, nil

	case "connectionHello":
		res := []ConnectionHello{}
		err := json.Unmarshal(raw, &res)
		if len(res) > 0 {
			return res[0], err
		}
		return ConnectionHello{}, nil

	case "connectionClose":
		res := []ConnectionClose{}
		err := json.Unmarshal(raw, &res)
		if len(res) > 0 {
			return res[0], err
		}
		return ConnectionClose{}, nil

	case "messageProtocolHandshake":
		res := MessageProtocolHandshake{}
		err := json.Unmarshal(raw, &res)
		return res, err

	default:
		return nil, errors.New("invalid type")
	}
}
