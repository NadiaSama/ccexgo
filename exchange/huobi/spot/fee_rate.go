package spot

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	TransactFeeRate struct {
		Symbol          string
		MakerFeeRate    decimal.Decimal `json:"makerFeeRate"`
		TakerFeeRate    decimal.Decimal `json:"takerFeeRate"`
		ActualMakerRate decimal.Decimal `json:"actualMakerRate"`
		ActualTakerRate decimal.Decimal `json:"actualTakerRate"`
	}
)

const (
	TransactFeeRateEndPoint = "/v2/reference/transact-fee-rate"
)

func (rc *RestClient) TransactFeeRate(ctx context.Context, symbols []string) ([]TransactFeeRate, error) {
	if len(symbols) <= 0 || len(symbols) > 10 {
		return nil, errors.Errorf("invalid symbols len=%d", len(symbols))
	}

	values := url.Values{}
	values.Add("symbols", strings.Join(symbols, ","))

	var ret []TransactFeeRate
	if err := rc.Request(ctx, http.MethodGet, TransactFeeRateEndPoint, values, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "fetch transact fee rate fail")
	}

	return ret, nil
}

func (rc *RestClient) FeeRate(ctx context.Context, symbols []exchange.Symbol) ([]*exchange.TradeFee, error) {
	var pairs []string
	for _, s := range symbols {
		pairs = append(pairs, s.String())
	}
	tfrs, err := rc.TransactFeeRate(ctx, pairs)
	if err != nil {
		return nil, err
	}

	var ret []*exchange.TradeFee
	for _, tfr := range tfrs {
		r, err := tfr.Parse()
		if err != nil {
			return nil, errors.WithMessage(err, "parse trade fee fail")
		}
		ret = append(ret, r)
	}
	return ret, nil
}

func (tfr *TransactFeeRate) Parse() (*exchange.TradeFee, error) {
	s, err := ParseSymbol(tfr.Symbol)
	if err != nil {
		return nil, err
	}

	return &exchange.TradeFee{
		Symbol: s,
		Maker:  tfr.ActualMakerRate,
		Taker:  tfr.ActualTakerRate,
		Raw:    *tfr,
	}, nil
}
