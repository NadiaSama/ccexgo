package exchange

import (
	"context"
	"time"

	"github.com/NadiaSama/ccexgo/internal/rpc"
)

type (
	Client struct {
		Conn    rpc.Conn
		Key     string
		Secret  string
		Ctx     context.Context
		Timeout time.Duration
	}
)
