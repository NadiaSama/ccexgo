package exchange

import (
	"github.com/NadiaSama/ccexgo/internal/rpc"
)

type (
	Client struct {
		Conn   rpc.Conn
		Key    string
		Secret string
	}
)
