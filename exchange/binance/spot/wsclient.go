package spot

import (
	"context"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
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
	case wcl.data <- &exchange.WSNotify{Data: n.Params, Chan: n.Method}:
	default:
	}
}

func (wcl *WSClient) Subscribe(ctx context.Context, channels ...exchange.Channel) error {
	param := make([]string, 0, len(channels))
	for _, c := range channels {
		param = append(param, c.String())
	}

	if err := wcl.Call(ctx, "1", MethodSubscibe, param, nil); err != nil {
		return errors.WithMessage(err, "subscribe error")
	}
	return nil
}
