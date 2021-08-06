package huobi

import (
	"context"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/internal/rpc"
)

type (
	//WSClient with auto response ping support
	WSClient struct {
		*exchange.WSClient
		data chan interface{}
	}

	//CallParam carry params which used by huobi websocket sub and pong
	CallParam struct {
		Pong int    `json:"pong,omitempty"`
		Sub  string `json:"sub,omitempty"`
		ID   string `json:"id,omitempty"`
	}
)

func NewWSClient(addr string, codec rpc.Codec, data chan interface{}) *WSClient {
	ret := &WSClient{
		data: data,
	}
	wc := exchange.NewWSClient(addr, codec, ret)

	ret.WSClient = wc
	return ret
}

func (ws *WSClient) Handle(ctx context.Context, notify *rpc.Notify) {
	if notify.Method == MethodPing {
		p := notify.Params.(int)
		ws.Call(ctx, MethodPong, "", &CallParam{Pong: p}, nil)
		return
	}

	ws.data <- &exchange.WSNotify{
		Exchange: Huobi,
		Chan:     notify.Method,
		Data:     notify.Params,
	}
}
