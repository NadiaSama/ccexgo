package swap

import (
	"github.com/NadiaSama/ccexgo/exchange/binance"
)

type (
	//RestClient struct
	RestClient struct {
		*binance.RestClient
		side *GetPositionSideResp
	}
)

const (
	SwapAPIHost     string = "fapi.binance.com"
	SwapTestAPIHost string = "testnet.binancefuture.com"
)

func NewRestClient(key, secret string) *RestClient {
	return &RestClient{
		RestClient: binance.NewRestClient(key, secret, SwapAPIHost),
	}
}

func NewTestRestClient(key, secret string) *RestClient {
	return &RestClient{
		RestClient: binance.NewRestClient(key, secret, SwapTestAPIHost),
	}
}
