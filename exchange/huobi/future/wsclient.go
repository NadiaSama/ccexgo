package future

import (
	"context"
	"fmt"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/huobi"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
)

type (
	WSClient struct {
		*WSClientDeriv
	}

	//WSClientDeriv define common logic used by huobi future swap ws
	WSClientDeriv struct {
		*huobi.WSClient
	}
)

const (
	FutureWSAddr = "wss://api.hbdm.com/ws"
)

func NewWSClient(codeMap map[string]string, data chan interface{}) *WSClient {
	codec := NewCodeC(codeMap)
	return &WSClient{
		NewWSClientDeriv(FutureWSAddr, codec, data),
	}
}

func NewWSClientDeriv(addr string, codec rpc.Codec, data chan interface{}) *WSClientDeriv {
	ret := &WSClientDeriv{
		WSClient: huobi.NewWSClient(addr, codec, data),
	}
	return ret
}

func (ws *WSClient) Subscribe(ctx context.Context, typ exchange.SubType, syms ...exchange.Symbol) error {
	if typ != exchange.SubTypeTrade {
		return errors.New("unsupport subtype")
	}
	channels := make([]string, len(syms))
	for i := range syms {
		sym := syms[i]
		s, ok := sym.(*FutureSymbol)
		if !ok {
			return errors.New("bad symbol type")
		}
		channels[i] = fmt.Sprintf("market.%s.trade.detail", s.WSSub())
	}

	return ws.WSClientDeriv.DoSubscribe(ctx, channels)
}

func (ws *WSClientDeriv) DoSubscribe(ctx context.Context, channels []string) error {
	for _, ch := range channels {
		//huobi ws future/swap subscribe do not send response, so just write one subid.
		cp := &huobi.CallParam{
			ID:  "s1",
			Sub: ch,
		}

		if err := ws.WSClient.Call(ctx, cp.ID, huobi.MethodSubscibe, cp, nil); err != nil {
			return errors.WithMessagef(err, "subscribe fail")
		}
	}
	return nil
}
