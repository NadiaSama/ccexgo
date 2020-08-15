package ftx

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
)

type (
	WSClient struct {
		*exchange.WSClient
		data chan interface{}
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

func NewWSClient() *WSClient {
	ret := &WSClient{}
	ret.WSClient = exchange.NewWSClient(ftxWSAddr, NewCodeC(), ret)
	return ret
}

func (ws *WSClient) Auth(ctx context.Context, key string, secret string) error {
	ts := time.Now().UnixNano() / 1e6
	es := fmt.Sprintf("%dwebsocket_login", ts)
	param := authParam{
		OP: "login",
		Args: authArgs{
			Key:  key,
			Sign: signature(secret, es),
			Time: ts,
		},
	}
	if err := ws.Conn.Call(ctx, "", "login", &param, nil); err != nil {
		return err
	}
	return nil
}

func (ws *WSClient) Subscribe(ctx context.Context, typ exchange.SubType, syms ...exchange.Symbol) error {
	if typ != exchange.SubTypeOrder {
		return errors.Errorf("unsupport subtype '%d'", typ)
	}

	var result subscribeResult
	req := callParam{
		Channel: "orders",
		OP:      "subscribe",
	}
	if err := ws.Conn.Call(ctx, req.Channel, req.OP, &req, &result); err != nil {
		return errors.WithMessagef(err, "subscribe orders fail")
	}

	if result.Type != typeSubscribed {
		return errors.Errorf("bad result %v", result)
	}
	return nil
}

func (ws *WSClient) Handle(ctx context.Context, notify *rpc.Notify) {
	if notify.Method == typePong {
		return
	}

	fmt.Printf("%s %s\n", notify.Method, string(notify.Params.(json.RawMessage)))
}
