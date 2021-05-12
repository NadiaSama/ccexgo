package swap

import (
	"context"
	"net/http"
	"net/url"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	CommisionRate struct {
		Symbol               string
		MakerCommissionRate  decimal.Decimal `json:"makerCommissionRate"`
		TakerCommissioinRate decimal.Decimal `json:"takerCommissionRate"`
	}
)

const (
	CommiionRateEndPoint = "/fapi/v1/commissionRate"
)

func (rc *RestClient) CommisionRate(ctx context.Context, symbol string) (*CommisionRate, error) {
	var resp CommisionRate
	values := url.Values{}
	values.Add("symbol", symbol)
	if err := rc.Request(ctx, http.MethodGet, CommiionRateEndPoint, values, nil, true, &resp); err != nil {
		return nil, errors.WithMessage(err, "fetch trade fee fail")
	}

	return &resp, nil
}

func (rc *RestClient) FeeRate(ctx context.Context, symbols []exchange.Symbol) ([]*exchange.TradeFee, error) {
	if len(symbols) != 1 {
		return nil, errors.Errorf("symbols must be 1")
	}
	cr, err := rc.CommisionRate(ctx, symbols[0].String())
	if err != nil {
		return nil, err
	}

	rate, err := cr.Parse()
	if err != nil {
		return nil, err
	}

	return []*exchange.TradeFee{rate}, nil
}

func (tf *CommisionRate) Parse() (*exchange.TradeFee, error) {
	s, err := ParseSymbol(tf.Symbol)
	if err != nil {
		return nil, err
	}

	return &exchange.TradeFee{
		Symbol: s,
		Taker:  tf.TakerCommissioinRate,
		Maker:  tf.MakerCommissionRate,
		Raw:    *tf,
	}, nil
}
