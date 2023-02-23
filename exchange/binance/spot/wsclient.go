package spot

import (
	"context"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/binance"
	"github.com/NadiaSama/ccexgo/internal/rpc"
)

type (
	// WSClient binance spot websocket client for public notify
	WSClient struct {
		*binance.NotifyClient
	}
)

const (
	WSClientEndPoint = "wss://stream.binance.com:9443/ws"
)

func NewWSClient(data chan interface{}) *WSClient {
	ret := &WSClient{}
	ret.NotifyClient = binance.NewNotifyClient(WSClientEndPoint, NewCodeC(), data, ret)
	return ret
}

func (wsc *WSClient) Handle(ctx context.Context, notify *rpc.Notify) {
	if notify.Method == TradeEvent || notify.Method == AggTradeEvent {
		trades, ok := notify.Params.([]*exchange.Trade)
		if !ok || len(trades) != 2 {
			return
		}
		wsc.Push(notify.Method, trades[0])
		wsc.Push(notify.Method, trades[1])
		return
	}

	wsc.Push(notify.Method, notify.Params)
}
