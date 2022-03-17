package option

import "github.com/NadiaSama/ccexgo/exchange/binance"

type (
	RestClient struct {
		*binance.RestClient
	}
)

func NewRestClient(key, secret string) *RestClient {
	return &RestClient{
		RestClient: binance.NewRestClient(key, secret, "vapi.binance.com"),
	}
}
