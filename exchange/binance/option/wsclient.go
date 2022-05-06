package option

import (
	"context"
	"fmt"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/binance"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
)

type (
	WSClient struct {
		*binance.WSClient
		data chan interface{}
	}
)

//NewTestWSClient return a wsclient which connect to binance option testnet
func NewTestWSClient(data chan interface{}, key, secret string) *WSClient {
	ret := &WSClient{
		data: data,
	}
	ret.WSClient = binance.NewWSClient(NewCodeC(), ret, NewTestRestClient(key, secret))

	return ret
}

func (wl *WSClient) Handle(ctx context.Context, notify *rpc.Notify) {
	fmt.Printf("%s %+v\n", notify.Method, notify.Params)
}

func (wl *WSClient) Subscribe(ctx context.Context, channels ...exchange.Channel) error {
	if len(channels) == 0 {
		return nil
	}

	params := make([]string, len(channels))
	for i, c := range channels {
		params[i] = c.String()
	}

	if err := wl.Call(ctx, "1", "", params, nil); err != nil {
		return errors.WithMessage(err, "rpc call fail")
	}
	return nil
}
