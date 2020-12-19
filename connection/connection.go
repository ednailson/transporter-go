package connection

type ListenFn func(conn Connection)

type Connection interface {
	Write(data []byte) error
	Message() []byte
	Close() error
}
