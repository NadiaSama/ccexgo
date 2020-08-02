package exchange

import (
	"context"
	"fmt"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
)

type (
	WSClient struct {
		handler rpc.Handler
		conn    rpc.Conn
		codec   rpc.Codec
		addr    string
	}

	NotifyTrade struct {
		Trades []Trade
		Chan   string
	}

	subscribeReq struct {
	}
)

func NewWSClient(addr string, codec rpc.Codec, handler rpc.Handler) *WSClient {
	return &WSClient{
		addr:    addr,
		codec:   codec,
		handler: handler,
	}
}

func (wc *WSClient) Run(ctx context.Context) error {
	stream, err := rpc.NewWebsocketStream(wc.addr, wc.codec)
	if err != nil {
		return err
	}

	conn := rpc.NewConn(stream)
	wc.conn = conn
	go wc.conn.Run(ctx, wc.handler)
	return nil
}

func (ws *WSClient) Subscribe(ctx context.Context, typ exchange.SubType, syms ...exchange.Symbol) error {
	if typ != exchange.SubTypeTrade {
		return errors.New("unsupport subtype")
	}
	fsym := make([]*FutureSymbol, len(syms))
	for i := range syms {
		sym := syms[i]
		s, ok := sym.(*FutureSymbol)
		if !ok {
			return errors.New("bad symbol type")
		}
		fsym[i] = s
	}

	for _, s := range fsym {
		//huobi ws subscribe do not send response, so just write one subid.
		cp := &callParam{
			ID:  "s1",
			Sub: fmt.Sprintf("market.%s.trade.detail", s.WSSub()),
		}
		if err := ws.conn.Call(ctx, methodSubscibe, cp, nil); err != nil {
			return errors.WithMessagef(err, "subscribe fail")
		}
	}
	return nil
}

func (ws *WSClient) Error() error {
	return ws.conn.Error()
}

func (ws *WSClient) Done() <-chan struct{} {
	return ws.conn.Done()
}

func (ws *WSClient) Close() error {
	return ws.conn.Close()
}

func (wc *WSClient) Handle(ctx context.Context, notify *rpc.Notify) {
	if notify.Method == methodPING {
		wc.conn.Call(ctx, methodPONG, &callParam{Pong: notify.Params.(int)}, nil)
		return
	}

	trades, ok := notify.Params.([]Trade)
	if !ok {
		return
	}

	wc.data <- &NotifyTrade{Chan: notify.Method, Trades: trades}
}
