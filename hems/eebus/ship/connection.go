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

type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type Connection struct {
	conn *websocket.Conn
	Log  Logger
}

func New(conn *websocket.Conn) (c *Connection) {
	return &Connection{
		conn: conn,
	}
}

func (c *Connection) log() Logger {
	return c.Log
}

// Connect performs the client connection handshake
func (c *Connection) Connect() error {
	err := c.init()
	if err == nil {
		err = c.hello()
	}
	if err == nil {
		err = c.clientProtocolHandshake()
	}

	// close connection if handshake or hello fails
	if err != nil {
		_ = c.conn.Close()
	}

	return err
}

// Serve performs the server connection handshake
func (c *Connection) Serve() error {
	err := c.init()
	if err == nil {
		c.log().Println("serve: hello")
		err = c.hello()
	}
	if err == nil {
		c.log().Println("serve: handshake")
		err = c.serverProtocolHandshake()
	}

	// close connection if handshake or hello fails
	if err != nil {
		_ = c.conn.Close()
	}

	timer := time.NewTimer(time.Minute)
	for {
		select {
		case <-timer.C:
			return nil
		default:
			var req CmiMessage
			typ, err := c.readJSON(&req)
			if err != nil {
				return err
			}

			typed := CmiDecode(req)
			c.log().Printf("serve: %d %+v", typ, typed)

			if _, ok := typed.(CmiCloseMsg); ok {
				return c.Close(false)
			}
		}
	}
}

func (c *Connection) writeBinary(msg []byte) error {
	if len(msg) < 3 {
		c.log().Printf("send: %0 x", msg)
	} else {
		c.log().Println("send:", string(msg))
	}

	err := c.conn.SetWriteDeadline(time.Now().Add(cmiReadWriteTimeout))
	if err == nil {
		c.conn.WriteMessage(websocket.BinaryMessage, msg)
	}
	return err
}

func (c *Connection) writeJSON(typ byte, jsonMsg interface{}) error {
	msg, err := json.Marshal(jsonMsg)
	if err != nil {
		return err
	}

	// add header
	b := bytes.NewBuffer([]byte{typ})
	b.WriteString(strconv.Quote(string(msg)))

	return c.writeBinary(b.Bytes())
}

func (c *Connection) readBinary() ([]byte, error) {
	err := c.conn.SetReadDeadline(time.Now().Add(cmiReadWriteTimeout))
	if err != nil {
		return nil, err
	}

	typ, msg, err := c.conn.ReadMessage()

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

func (c *Connection) readJSON(jsonMsg interface{}) (byte, error) {
	b, err := c.readBinary()
	if err != nil {
		return 0, err
	}

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

// Close closes the service connection
func (c *Connection) Close(active bool) error {
	var err error
	if active {
		err = c.closeActive()
	} else {
		err = c.closePassive()
	}

	// close physical connection
	_ = c.conn.Close()

	return err
}
