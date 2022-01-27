package ftx

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

type (
	Market struct {
		Name           string  `json:"name"`
		BaseCurrency   string  `json:"baseCurrency"`
		QuoteCurrency  string  `json:"quoteCurrency"`
		Type           string  `json:"type"`
		Underlying     string  `json:"underlying"`
		Enabled        bool    `json:"enabled"`
		Ask            float64 `json:"ask"`
		Bid            float64 `json:"bid"`
		Last           float64 `json:"last"`
		PostOnly       bool    `json:"postOnly"`
		PriceIncrement float64 `json:"priceIncrement"`
		SizeIncrement  float64 `json:"sizeIncrement"`
		MinProvideSize float64 `json:"minProvideSize"`
		Restricted     bool    `json:"restricted"`
	}

	Depth struct {
		Asks [][2]float64
		Bids [][2]float64
	}

	BookReq struct {
		Market string
		Sz     string
	}
)

func NewBookReq(market, sz string) *BookReq {
	return &BookReq{
		Market: market,
		Sz:     sz,
	}
}

func (rc *RestClient) Markets(ctx context.Context) ([]Market, error) {
	var resp []Market
	if err := rc.request(ctx, http.MethodGet, "/markets", nil, nil, false, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (rc *RestClient) Books(ctx context.Context, req BookReq) (*Depth, error) {
	var ret Depth
	values := url.Values{}
	values.Add("market_name", req.Market)
	if req.Sz != "" {
		if _, err := strconv.Atoi(req.Sz); err != nil {
			return nil, errors.WithMessagef(err, "invalid sz '%s'", req.Sz)
		}
		values.Add("sz", req.Sz)
	}

	uri := fmt.Sprintf("/markets/%s/orderbook", req.Market)
	if err := rc.request(ctx, http.MethodGet, uri, values, nil, false, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}
