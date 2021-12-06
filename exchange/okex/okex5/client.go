package okex5

import (
	"context"
	"io"
	"net/url"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/okex"
	"github.com/pkg/errors"
)

type (
	RestClient struct {
		client *okex.RestClient
	}

	GetRequest struct {
		fields map[string]string
	}
	RestResponse struct {
		Code string      `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}
)

func NewGetRequest() *GetRequest {
	return &GetRequest{
		fields: make(map[string]string),
	}
}

func (gr *GetRequest) Add(key, val string) {
	gr.fields[key] = val
}

func (gr *GetRequest) Values() url.Values {
	ret := url.Values{}
	for k, v := range gr.fields {
		ret.Add(k, v)
	}
	return ret
}

func NewRestClient(key, secret, pass string) *RestClient {
	return &RestClient{
		client: okex.NewRestClient(key, secret, pass),
	}
}

func NewTestRestClient(key, secret, pass string) *RestClient {
	return &RestClient{
		client: okex.NewTESTRestClient(key, secret, pass),
	}
}

//Request do okexv5 rest request. response data field will be store into dst
func (rc *RestClient) Request(ctx context.Context, method string, endPoint string, params url.Values, body io.Reader, sign bool, dst interface{}) error {
	resp := RestResponse{
		Data: dst,
	}

	if err := rc.client.Request(ctx, method, endPoint, params, body, sign, &resp); err != nil {
		return errors.WithMessagef(err, "request %s fail", endPoint)
	}

	if resp.Code != "0" {
		return errors.Errorf("request: %s fail code: %s msg: %s", endPoint, resp.Code, resp.Msg)
	}

	return nil
}

func (rc *RestClient) Property() exchange.Property {
	return exchange.Property{
		Trades: &exchange.TradesProp{
			MaxDuration: time.Hour * 168,
			SuportID:    true,
			SupportTime: false,
		},

		Finance: &exchange.FinanceProp{
			MaxDuration: time.Hour * 168,
			SuportID:    true,
			SupportTime: false,
		},
	}
}
