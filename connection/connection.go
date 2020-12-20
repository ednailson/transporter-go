package connection

import "net"

type ListenFn func(conn Connection)

type Connection interface {
	Write(data []byte) error
	Message() []byte
	RemoteAddr() net.Addr
	Close() error
}
