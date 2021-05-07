package websocket

import (
	"context"

	"github.com/NadiaSama/ccexgo/exchange"
)

type (
	//Conn Unique exchange websocket client
	Conn interface {
		Run(ctx context.Context) error
		Close() error
		Error() error
		Done() <-chan struct{}
		Subscribe(ctx context.Context, channels ...exchange.Channel) error
		UnSubscribe(ctx context.Context, channels ...exchange.Channel) error
	}
)
