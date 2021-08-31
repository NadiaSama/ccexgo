package ftx

import (
	"context"
	"fmt"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
)

type (
	WSClient struct {
		*exchange.WSClient
		data   chan interface{}
		key    string
		secret string
	}

	subscribeResult struct {
		Type    string `json:"type"`
		Channel string `json:"channel"`
		Market  string `json:"market"`
	}
)

const (
	ftxWSAddr = "wss://ftx.com/ws/"
)

var (
	subMaps map[exchange.SubType]string = map[exchange.SubType]string{
		exchange.SubTypePrivateOrder: channelOrders,
		exchange.SubTypePrivateTrade: channelFills,
		exchange.SubTypeOrderBook:    channelOrderBook,
	}
)

func NewWSClient(key, secret string, data chan interface{}) *WSClient {
	ret := &WSClient{
		key:    key,
		secret: secret,
	}
	ret.WSClient = exchange.NewWSClient(ftxWSAddr, NewCodeC(), ret)
	ret.data = data
	return ret
}

func (ws *WSClient) Run(ctx context.Context) error {
	if err := ws.WSClient.Run(ctx); err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(time.Second * 15)
		for {
			select {
			case <-ctx.Done():
				return

			case <-ws.Done():
				return

			case <-ticker.C:
				param := &callParam{
					OP: "ping",
				}
				ws.Call(ctx, "", "ping", &param, nil)
			}
		}
	}()
	return nil
}
func (ws *WSClient) Auth(ctx context.Context) error {
	ts := time.Now().UnixNano() / 1e6
	es := fmt.Sprintf("%dwebsocket_login", ts)
	param := authParam{
		OP: "login",
		Args: authArgs{
			Key:  ws.key,
			Sign: signature(ws.secret, es),
			Time: ts,
		},
	}
	if err := ws.Conn.Call(ctx, "", "login", &param, nil); err != nil {
		return err
	}
	return nil
}

func (ws *WSClient) Subscribe(ctx context.Context, typ exchange.SubType, syms ...exchange.Symbol) error {
	channel, ok := subMaps[typ]
	if !ok {
		return errors.Errorf("unsupport subtype '%d'", typ)
	}

	var result subscribeResult
	var req callParam
	if channel == channelOrderBook {
		if len(syms) != 1 {
			return errors.Errorf("ftx multi subscribe not support")
		}
		req = callParam{
			Channel: channel,
			OP:      "subscribe",
			Market:  syms[0].String(),
		}
	} else {
		req = callParam{
			Channel: channel,
			OP:      "subscribe",
		}
	}

	if err := ws.Conn.Call(ctx, subID(req.Channel, req.Market), req.OP, &req, &result); err != nil {
		return errors.WithMessagef(err, "subscribe orders fail")
	}

	if result.Type != typeSubscribed {
		return errors.Errorf("bad result %v", result)
	}
	return nil
}

func (ws *WSClient) Handle(ctx context.Context, notify *rpc.Notify) {
	if notify.Method == typePong || notify.Method == typeInfo {
		return
	}

	ws.data <- &exchange.WSNotify{
		Exchange: ftxExchange,
		Chan:     notify.Method,
		Data:     notify.Params,
	}
}
