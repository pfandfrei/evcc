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
	Conn *websocket.Conn
	Log  Logger

	phase  byte
	inC    chan []byte
	errC   chan error
	closeC chan struct{}
}

func (c *Transport) log() Logger {
	return c.Log
}

func (c *Transport) writeBinary(msg []byte) error {
	if len(msg) > 2 {
		c.log().Println("send:", string(msg))
	}

	err := c.Conn.SetWriteDeadline(time.Now().Add(cmiReadWriteTimeout))
	if err == nil {
		err = c.Conn.WriteMessage(websocket.BinaryMessage, msg)
	}

	return err
}

func (c *Transport) writeJSON(typ byte, jsonMsg interface{}) error {
	// time.Sleep(time.Duration(rand.Int31n(int32(time.Second))))

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

func (c *Transport) readBinaryNoDeadline() ([]byte, error) {
	typ, msg, err := c.Conn.ReadMessage()

	if err == nil {
		if len(msg) > 2 {
			c.log().Println("recv:", string(msg))
		}

		if typ != websocket.BinaryMessage {
			err = fmt.Errorf("invalid message type: %d", typ)
		}
	}

	return msg, err
}

func (c *Transport) readBinary() ([]byte, error) {
	err := c.Conn.SetReadDeadline(time.Now().Add(cmiReadWriteTimeout))
	if err != nil {
		return nil, err
	}

	return c.readBinaryNoDeadline()
}

func (c *Transport) readPump() {
	for {
		select {
		case <-c.closeC:
			return

		default:
			if b, err := c.readBinaryNoDeadline(); err != nil {
				c.errC <- err
			} else {
				c.inC <- b
			}
		}
	}
}

func (c *Transport) readMessage(timerC <-chan time.Time) (interface{}, error) {
	select {
	case <-timerC:
		return nil, ErrTimeout

	case <-c.closeC:
		return nil, net.ErrClosed

	case b := <-c.inC:
		if len(b) < 2 {
			return nil, errors.New("invalid length")
		}
		if b[0] < 1 {
			return nil, errors.New("invalid phase")
		}

		var cmi CmiMessage
		if err := json.Unmarshal(b[1:], &cmi); err != nil {
			return nil, err
		}

		return decodeCmi(cmi)

	case err := <-c.errC:
		return nil, err
	}
}
