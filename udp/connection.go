package udp

import (
	"net"
)

type conn struct {
	protocol   string
	message    []byte
	remoteAddr net.Addr
	udpConn    *net.UDPConn
}

func newConn(message []byte, udpConn *net.UDPConn, remoteAddr net.Addr, protocol string) *conn {
	return &conn{
		message:    message,
		remoteAddr: remoteAddr,
		udpConn:    udpConn,
		protocol:   protocol,
	}
}

func (c *conn) Message() []byte {
	return c.message
}

func (c *conn) Write(data []byte) error {
	udpAddr, _ := net.ResolveUDPAddr(c.protocol, c.remoteAddr.String())
	_, err := c.udpConn.WriteToUDP(data, udpAddr)
	return err
}

func (c *conn) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c *conn) Close() error {
	return c.udpConn.Close()
}
