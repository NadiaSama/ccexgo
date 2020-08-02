package future

import (
	"context"

	"github.com/NadiaSama/ccexgo/exchange/huobi"
)

type (
	RestClient struct {
		*huobi.RestClient
		futureSymbolMap map[string]*FutureSymbol
	}
)

const (
	FutureHost = "api.hbdm.com"
)

func NewRestClient(key, secret string) *RestClient {
	return &RestClient{
		RestClient: huobi.NewRestClient(key, secret, FutureHost),
		futureSymbolMap: make(map[string]*FutureSymbol),
	}
}

func (rc *RestClient) Init(ctx context.Context) error {
	return rc.initFutureSymbol(ctx)
}