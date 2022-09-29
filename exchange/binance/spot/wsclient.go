package spot

import (
	"context"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/internal/rpc"
)

type (
	// WSClient binance spot websocket client for public notify
	WSClient struct {
		*exchange.WSClient
		data chan interface{}
	}
)

const (
	WSClientEndPoint = "wss://stream.binance.com:9443/ws"
)

func NewWSClient(data chan interface{}) *WSClient {
	ret := &WSClient{
		data: data,
	}
	ret.WSClient = exchange.NewWSClient(WSClientEndPoint, NewCodeC(), ret)
	return ret
}

func (wcl *WSClient) Handle(ctx context.Context, n *rpc.Notify) {
	select {
	case wcl.data <- exchange.WSNotify{Data: n.Params, Chan: n.Method}:
	default:
	}
}

func (wcl *WSClient) Subscribe(ctx context.Context, channels ...exchange.Channel) error {
	return nil
}
