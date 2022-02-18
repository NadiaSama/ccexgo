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

func (ws *WSClient) Subscribe(ctx context.Context, channels ...exchange.Channel) error {
	if len(channels) != 1 {
		return errors.Errorf("ftx multi subscribe not support")
	}

	ch := channels[0]

	var result subscribeResult
	var req callParam

	switch t := ch.(type) {
	case *OrderBookChannel: // orderbook[public]
		req = callParam{
			Channel: channelOrderBook,
			OP:      "subscribe",
			Market:  t.symbol.String(),
		}

	case *FillChannel: // fills[private]
		req = callParam{
			Channel: channelFills,
			OP:      "subscribe",
			Market:  t.symbol.String(),
		}

	case *OrderChannel: // orders[private]
		req = callParam{
			Channel: channelOrders,
			OP:      "subscribe",
			Market:  t.symbol.String(),
		}

	case *TradeChannel: // trades[public]
		req = callParam{
			Channel: channelTrades,
			OP:      "subscribe",
			Market:  t.symbol.String(),
		}

	case *TickerChannel: // ticker[public]
		req = callParam{
			Channel: channelTicker,
			OP:      "subscribe",
			Market:  t.symbol.String(),
		}

	case *MarketChannel: // markets[public]
		req = callParam{
			Channel: channelMarkets,
			OP:      "subscribe",
			Market:  t.symbol.String(),
		}

	case *OrderbookGroupedChannel: // orderbookGrouped[public]
		req = callParam{
			Channel: channelOrderbookGrouped,
			OP:      "subscribe",
			Market:  t.symbol.String(),
		}

	default:
		return errors.Errorf("unsupport typ %+v", t)
	}

	if err := ws.Conn.Call(ctx, subID(req.Channel, req.Market), req.OP, &req, &result); err != nil {
		return errors.WithMessagef(err, "subscribe orders fail")
	}

	if result.Type != typeSubscribed {
		return errors.Errorf("bad result %v", result)
	}
	return nil
}

func (ws *WSClient) UnSubscribe(ctx context.Context, channels ...exchange.Channel) error {
	if len(channels) != 1 {
		return errors.Errorf("ftx multi unsubscribe not support")
	}

	ch := channels[0]

	var result subscribeResult
	var req callParam

	switch t := ch.(type) {
	case *OrderBookChannel:
		req = callParam{
			Channel: channelOrderBook,
			OP:      "unsubscribe",
			Market:  t.symbol.String(),
		}

	case *FillChannel:
		req = callParam{
			Channel: channelFills,
			OP:      "unsubscribe",
			Market:  t.symbol.String(),
		}

	case *OrderChannel:
		req = callParam{
			Channel: channelOrders,
			OP:      "unsubscribe",
			Market:  t.symbol.String(),
		}

	case *TradeChannel:
		req = callParam{
			Channel: channelTrades,
			OP:      "unsubscribe",
			Market:  t.symbol.String(),
		}

	case *TickerChannel:
		req = callParam{
			Channel: channelTicker,
			OP:      "unsubscribe",
			Market:  t.symbol.String(),
		}

	case *MarketChannel:
		req = callParam{
			Channel: channelMarkets,
			OP:      "unsubscribe",
			Market:  t.symbol.String(),
		}

	case *OrderbookGroupedChannel:
		req = callParam{
			Channel: channelOrderbookGrouped,
			OP:      "subscribe",
			Market:  t.symbol.String(),
		}

	default:
		return errors.Errorf("unsupport typ %+v", t)
	}

	if err := ws.Conn.Call(ctx, subID(req.Channel, req.Market), req.OP, &req, &result); err != nil {
		return errors.WithMessagef(err, "Unsubscribe orders fail")
	}

	if result.Type != typeUnSubscribed {
		return errors.Errorf("bad result %v", result)
	}
	return nil
}

func (ws *WSClient) Handle(ctx context.Context, notify *rpc.Notify) {
	// if notify.Method == typePong || notify.Method == typeInfo {
	// 	return
	// }

	ws.data <- &exchange.WSNotify{
		Exchange: ftxExchange,
		Chan:     notify.Method,
		Data:     notify.Params,
	}
}
