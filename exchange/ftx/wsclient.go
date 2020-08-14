package ftx

import (
	"github.com/NadiaSama/ccexgo/exchange"
)

type (
	WSClient struct {
		*exchange.WSClient
		data chan interface{}
	}
)

const (
	ftxWSAddr = "wss://ftx.com/ws/"
)

func NewWSClient() *WSClient {
	return &WSClient{
		exchange.NewWSClient(ftxWSAddr, nil, nil)
	}
}
