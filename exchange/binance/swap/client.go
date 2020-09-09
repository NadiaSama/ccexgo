package swap

import (
	"context"

	"github.com/NadiaSama/ccexgo/exchange/binance"
)

type (
	//RestClient struct
	RestClient struct {
		*binance.RestClient
		symbols map[string]*SwapSymbol
	}
)

const (
	SwapAPIHost string = "fapi.binance.com"
)

func NewRestClient(key, secret string) *RestClient {
	return &RestClient{
		binance.NewRestClient(key, secret, SwapAPIHost),
		map[string]*SwapSymbol{},
	}
}

func (rc *RestClient) Init(ctx context.Context) error {
	return rc.loadSymbol(ctx)
}
