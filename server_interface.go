package transporter

import (
	"github.com/ednailson/transporter-go/connection"
)

type Server interface {
	Listen(fn connection.ListenFn) <-chan error
	Close() error
}
