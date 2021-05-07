package spot

import "github.com/NadiaSama/ccexgo/exchange/binance"

type (
	RestClient struct {
		*binance.RestClient
	}
)

func NewSpotClient(key, secret string) *RestClient {
	return &RestClient{
		binance.NewRestClient(key, secret, "api.binance.com"),
	}
}
