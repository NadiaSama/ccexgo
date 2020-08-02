package swap

import (
	"context"

	"github.com/NadiaSama/ccexgo/exchange/huobi"
)

type (
	RestClient struct {
		*huobi.RestClient
		swapCodeMap map[string]*Symbol
	}
)

const (
	SwapHost = "api.hbdm.com"
)

func NewRestClient(key string, secret string) *RestClient {
	return &RestClient{
		RestClient:  huobi.NewRestClient(key, secret, SwapHost),
		swapCodeMap: make(map[string]*Symbol),
	}
}

func (rc *RestClient) Init(ctx context.Context) error {
	return rc.initSymbol(ctx)
}
