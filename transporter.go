package transporter

import (
	"errors"
	"github.com/ednailson/transporter-go/tcp"
	"github.com/ednailson/transporter-go/udp"
)

func New(protocol, host string, port int) (Server, error) {
	switch protocol {
	case "tcp", "tcp4", "tcp6":
		return tcp.New(protocol, host, port)
	case "udp", "udp4", "udp6":
		return udp.New(protocol, host, port)
	default:
		return nil, errors.New("invalid protocol")
	}
}
