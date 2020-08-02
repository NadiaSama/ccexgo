package swap

import (
	"context"
	"fmt"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/huobi/future"
	"github.com/pkg/errors"
)

type (
	WSClient struct {
		*future.WSClientDeriv
	}
)

const (
	SwapWSAddr = "wss://api.hbdm.com/swap-ws"
)

func NewWSClient(data chan interface{}) *WSClient {
	return &WSClient{
		future.NewWSClientDeriv(SwapWSAddr, NewCodeC(), data),
	}
}

func (ws *WSClient) Subscribe(ctx context.Context, typ exchange.SubType, syms ...exchange.Symbol) error {
	if typ != exchange.SubTypeTrade {
		return errors.Errorf("unsupport subscribe type %d", typ)
	}

	channels := make([]string, len(syms))
	for i, s := range syms {
		sym, ok := s.(*Symbol)
		if !ok {
			return errors.Errorf("bad symbol type %v", s)
		}
		channels[i] = fmt.Sprintf("market.%s.trade.detail", sym.Index())
	}

	return ws.DoSubscribe(ctx, channels)
}
