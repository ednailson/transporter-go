package tcp

import (
	"net"
)

type conn struct {
	message    []byte
	connection net.Conn
}

func newConn(message []byte, connection net.Conn) *conn {
	return &conn{
		message:    message,
		connection: connection,
	}
}

func (c *conn) Message() []byte {
	return c.message
}

func (c *conn) Write(data []byte) error {
	_, err := c.connection.Write(data)
	return err
}

func (c *conn) RemoteAddr() net.Addr {
	return c.connection.RemoteAddr()
}

func (c *conn) Close() error {
	return c.connection.Close()
}
