package ftx

import (
	"context"
	"net/http"
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
		Restricted     bool    `json:"restricted"`
	}
)

func (rc *RestClient) Markets(ctx context.Context) ([]Market, error) {
	var resp []Market
	if err := rc.request(ctx, http.MethodGet, "/markets", nil, nil, false, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
