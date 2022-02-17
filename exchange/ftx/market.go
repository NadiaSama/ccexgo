package ftx

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/NadiaSama/ccexgo/exchange"
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
		Depth  int
	}

	Markets       struct{}
	MarketChannel struct {
		symbol exchange.Symbol
	}
)

func NewBookReq(market string, depth int) *BookReq {
	return &BookReq{
		Market: market,
		Depth:  depth,
	}
}

func (rc *RestClient) Markets(ctx context.Context) ([]Market, error) {
	var resp []Market
	if err := rc.request(ctx, http.MethodGet, "/markets", nil, nil, false, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (rc *RestClient) Books(ctx context.Context, req *BookReq) (*Depth, error) {
	var ret Depth
	values := url.Values{}
	if req.Depth != 0 {
		depth := strconv.Itoa(req.Depth)
		values.Add("depth", depth)
	}

	uri := fmt.Sprintf("/markets/%s/orderbook", req.Market)
	if err := rc.request(ctx, http.MethodGet, uri, values, nil, false, &ret); err != nil {
		return nil, err
	}

	return &ret, nil
}

func NewMarketsChannel(sym exchange.Symbol) exchange.Channel {
	return &MarketChannel{
		symbol: sym,
	}
}

func (m *MarketChannel) String() string {
	return m.symbol.String()
}
