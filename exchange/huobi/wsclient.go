package huobi

import (
	"context"
	"strconv"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
)

type (
	//WSClient with auto response ping support
	WSClient struct {
		*exchange.WSClient
		data chan interface{}
	}

	//CallParam carry params which used by huobi websocket sub and pong
	CallParam struct {
		Pong  int    `json:"pong,omitempty"`
		Sub   string `json:"sub,omitempty"`
		UnSub string `json:"unsub,omitempty"`
		ID    string `json:"id,omitempty"`
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

func (ws *WSClient) Subscribe(ctx context.Context, channels ...exchange.Channel) error {
	if len(channels) != 1 {
		return errors.Errorf("only one channel subscribe support")
	}
	for i, ch := range channels {
		param := CallParam{
			ID:  strconv.Itoa(i),
			Sub: ch.String(),
		}

		var dest Response

		if err := ws.Call(ctx, param.ID, MethodSubscibe, &param, &dest); err != nil {
			return err
		}

		if dest.Status != "ok" {
			return errors.Errorf("subscirbe error %+v", dest)
		}
	}
	return nil
}

func (ws *WSClient) UnSubscribe(ctx context.Context, channels ...exchange.Channel) error {
	if len(channels) != 1 {
		return errors.Errorf("only one channel subscribe support")
	}
	for i, ch := range channels {
		param := CallParam{
			ID:    strconv.Itoa(i),
			UnSub: ch.String(),
		}

		var dest Response

		if err := ws.Call(ctx, param.ID, MethodUnSubscribe, &param, &dest); err != nil {
			return err
		}

		if dest.Status != "ok" {
			return errors.Errorf("subscirbe error %+v", dest)
		}
	}
	return nil
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
