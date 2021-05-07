package swap

import (
	"github.com/NadiaSama/ccexgo/exchange/binance"
)

type (
	//RestClient struct
	RestClient struct {
		*binance.RestClient
	}
)

const (
	SwapAPIHost string = "fapi.binance.com"
)

func NewRestClient(key, secret string) *RestClient {
	return &RestClient{
		binance.NewRestClient(key, secret, SwapAPIHost),
	}
}
