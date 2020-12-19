package tcp

import (
	"fmt"
	"github.com/ednailson/transporter-go/connection"
	"net"
)

type server struct {
	listener *net.TCPListener
}

func New(protocol, ip string, port int) (*server, error) {
	addr, err := net.ResolveTCPAddr(protocol, fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenTCP(protocol, addr)
	if err != nil {
		return nil, err
	}
	return &server{
		listener: listener,
	}, nil
}

func (s *server) Listen(fn connection.ListenFn) <-chan error {
	ch := make(chan error)
	go func() {
		defer close(ch)
		for {
			tcpConn, err := s.listener.Accept()
			if err != nil {
				ch <- err
				return
			}
			go func() {
				for {
					buf := make([]byte, 1024)
					readLen, err := tcpConn.Read(buf)
					if err != nil {
						ch <- err
						return
					}
					go func() {
						fn(newConn(buf[:readLen], tcpConn))
					}()
				}
			}()
		}
	}()
	return ch
}

func (s *server) Close() error {
	return s.listener.Close()
}
