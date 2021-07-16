package swap

import (
	"context"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/huobi/future"
)

type (
	WSClient struct {
		*future.WSClientDeriv
	}
)

const (
	SwapWSAddr = "wss://api.hbdm.com/swap-ws"
)

func NewWSClient(data chan interface{}) *WSClient {
	return &WSClient{
		WSClientDeriv: future.NewWSClientDeriv(SwapWSAddr, NewCodeC(), data),
	}
}

func (ws *WSClient) Subscribe(ctx context.Context, cs ...exchange.Channel) error {
	var channels []string
	for _, c := range cs {
		channels = append(channels, c.String())
	}
	return ws.DoSubscribe(ctx, channels)
}
