package spot

import (
	"github.com/NadiaSama/ccexgo/exchange/binance"
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
	ret.NotifyClient = binance.NewNotifyClient(WSClientEndPoint, NewCodeC(), data, nil)
	return ret
}
