package udp

import (
	"fmt"
	"github.com/ednailson/transporter-go/connection"
	"net"
)

type server struct {
	conn     *net.UDPConn
	protocol string
}

func New(protocol, ip string, port int) (*server, error) {
	serverAddr, err := net.ResolveUDPAddr(protocol, fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return nil, err
	}
	udpConn, err := net.ListenUDP(protocol, serverAddr)
	if err != nil {
		return nil, err
	}
	return &server{
		conn:     udpConn,
		protocol: protocol,
	}, nil
}

func (s *server) Listen(fn connection.ListenFn) <-chan error {
	ch := make(chan error)
	go func() {
		defer close(ch)
		for {
			buf := make([]byte, 1024)
			readLen, remoteAddr, err := s.conn.ReadFromUDP(buf)
			if err != nil {
				ch <- err
				break
			}
			go func() {
				fn(newConn(buf[:readLen], s.conn, remoteAddr, s.protocol))
			}()
		}
	}()
	return ch
}

func (s *server) Close() error {
	return s.conn.Close()
}
