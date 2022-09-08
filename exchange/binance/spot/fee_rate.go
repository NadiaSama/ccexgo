package spot

import (
	"context"
	"net/http"
	"net/url"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	TradeFee struct {
		Symbol string          `json:"symbol"`
		Maker  decimal.Decimal `json:"makerCommission"`
		Taker  decimal.Decimal `json:"takerCommission"`
	}
)

const (
	TradeFeeEndPoint = "/sapi/v1/asset/tradeFee"
)

func (rc *RestClient) TradeFee(ctx context.Context, symbol string) ([]TradeFee, error) {
	var tfr []TradeFee
	values := url.Values{}
	if symbol != "" {
		values.Add("symbol", symbol)
	}
	if err := rc.Request(ctx, http.MethodGet, TradeFeeEndPoint, values, nil, true, &tfr); err != nil {
		return nil, errors.WithMessage(err, "fetch trade fee fail")
	}

	return tfr, nil
}

func (rc *RestClient) FeeRate(ctx context.Context, symbols []exchange.Symbol) ([]*exchange.TradeFee, error) {
	tfs, err := rc.TradeFee(ctx, "")
	if err != nil {
		return nil, err
	}

	set := make(map[string]struct{})
	for _, s := range symbols {
		set[s.String()] = struct{}{}
	}

	var ret []*exchange.TradeFee

	for i := range tfs {
		tf := tfs[i]

		r, err := tf.Parse()
		if err != nil {
			return nil, errors.WithMessage(err, "parse trade fee fail")
		}

		if _, ok := set[tf.Symbol]; !ok && len(set) != 0 {
			continue
		}
		ret = append(ret, r)
	}
	return ret, nil
}

func (tf *TradeFee) Parse() (*exchange.TradeFee, error) {
	s, err := ParseSymbol(tf.Symbol)
	if err != nil {
		return nil, err
	}

	return &exchange.TradeFee{
		Symbol: s,
		Taker:  tf.Taker,
		Maker:  tf.Maker,
		Raw:    *tf,
	}, nil
}
