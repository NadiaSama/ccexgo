package swap

import (
	"github.com/NadiaSama/ccexgo/exchange/binance"
)

type (
	WSClient struct {
		*binance.NotifyClient
	}
)

const (
	WSClientEndPoint = "wss://fstream.binance.com/ws"
)

func NewWSClient(data chan interface{}) *WSClient {
	ret := &WSClient{
		NotifyClient: binance.NewNotifyClient(WSClientEndPoint, NewCodeC(), data),
	}
	return ret
}
