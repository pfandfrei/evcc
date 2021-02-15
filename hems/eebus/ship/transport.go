package ship

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const cmiReadWriteTimeout = 10 * time.Second

// Transport is the physical transport layer
type Transport struct {
	Conn *websocket.Conn
	Log  Logger
}

func (c *Transport) log() Logger {
	return c.Log
}

func (c *Transport) writeBinary(msg []byte) error {
	if len(msg) < 3 {
		c.log().Printf("send: %0 x", msg)
	} else {
		c.log().Println("send:", string(msg))
	}

	err := c.Conn.SetWriteDeadline(time.Now().Add(cmiReadWriteTimeout))
	if err == nil {
		err = c.Conn.WriteMessage(websocket.BinaryMessage, msg)
	}

	return err
}

func (c *Transport) writeJSON(typ byte, jsonMsg interface{}) error {
	msg, err := json.Marshal(jsonMsg)
	if err != nil {
		return err
	}

	// add header
	b := bytes.NewBuffer([]byte{typ})

	_, err = b.WriteString(strconv.Quote(string(msg)))
	if err == nil {
		err = c.writeBinary(b.Bytes())
	}

	return err
}

func (c *Transport) readBinaryNoDeadline() ([]byte, error) {
	typ, msg, err := c.Conn.ReadMessage()

	if err == nil {
		if len(msg) < 3 {
			c.log().Printf("recv: %0 x", msg)
		} else {
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

func (c *Transport) handleJSON(b []byte, jsonMsg interface{}) (byte, error) {
	if len(b) < 2 {
		return 0, errors.New("invalid message")
	}

	typ := b[0]

	q, err := strconv.Unquote(string(b[1:]))
	if err == nil {
		msg := []byte(q)
		err = json.Unmarshal(msg, &jsonMsg)
	}

	return typ, err
}

func (c *Transport) readJSON(jsonMsg interface{}) (byte, error) {
	b, err := c.readBinary()
	if err != nil {
		return 0, err
	}

	return c.handleJSON(b, &jsonMsg)
}

func (c *Transport) waitJSON(jsonMsg interface{}) (byte, error) {
	b, err := c.readBinaryNoDeadline()
	if err != nil {
		return 0, err
	}

	return c.handleJSON(b, &jsonMsg)
}

func (c *Transport) readJSONWithTimeout(timeout time.Duration, jsonMsg interface{}) (byte, error) {
	err := c.Conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return 0, err
	}

	return c.waitJSON(&jsonMsg)
}
