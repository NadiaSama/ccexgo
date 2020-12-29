package spot

import "github.com/NadiaSama/ccexgo/exchange/okex"

type (
	RestClient struct {
		*okex.RestClient
	}
)

func NewRestClient(key, secret, pass string) *RestClient {
	return &RestClient{
		okex.NewRestClient(key, secret, pass),
	}
}

func NewTestRestClient(key, secret, pass string) *RestClient {
	return &RestClient{
		okex.NewTESTRestClient(key, secret, pass),
	}
}
