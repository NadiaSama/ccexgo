package swap

import (
	"github.com/NadiaSama/ccexgo/exchange/huobi"
)

type (
	RestClient struct {
		*huobi.RestClient
	}
)

const (
	SwapHost = "api.hbdm.com"
)

func NewRestClient(key string, secret string) *RestClient {
	return &RestClient{
		RestClient: huobi.NewRestClient(key, secret, SwapHost),
	}
}
