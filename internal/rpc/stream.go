package rpc

type (
	Stream interface {
		Read() (Response, error)
		Write(Request) error
		Close() error
	}
)
