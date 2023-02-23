package delivery

import "github.com/NadiaSama/ccexgo/exchange/binance"

type (
	WSClient struct {
		*binance.NotifyClient
	}
)

const (
	WSClientEndPoint = "wss://dstream.binance.com/ws"
)

func NewWSClient(data chan interface{}) *WSClient {
	ret := &WSClient{}
	ret.NotifyClient = binance.NewNotifyClient(WSClientEndPoint, NewCodeC(), data, nil)
	return ret
}
