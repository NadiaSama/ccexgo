package spot

import (
	"context"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/huobi"
)

type (
	RestClient struct {
		*huobi.RestClient
		pair2Symbol map[string]exchange.SpotSymbol
	}
)

const (
	SpotHost = "api.huobi.pro"
)

func NewRestClient(key, secret string) *RestClient {
	return &RestClient{
		RestClient:  huobi.NewRestClient(key, secret, SpotHost),
		pair2Symbol: make(map[string]exchange.SpotSymbol),
	}
}

func (rc *RestClient) Init(ctx context.Context) error {
	return rc.initSymbol(ctx)
}
