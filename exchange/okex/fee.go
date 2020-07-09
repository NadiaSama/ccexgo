package okex

import (
	"context"
	"net/http"
	"strconv"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
)

type (
	TradeFee struct {
		Maker     string `json:"maker"`
		Taker     string `json:"taker"`
		Timestamp string `json:"timestamp"`
	}
)

func (rc *RestClient) FeeRate(ctx context.Context, syms ...exchange.Symbol) ([]exchange.TradeFee, error) {
	if len(syms) > 0 {
		return nil, errors.New("okex trade fee do not symbols")
	}
	var result TradeFee
	if err := rc.request(ctx, http.MethodGet, "/api/spot/v3/trade_fee", nil, nil, true, &result); err != nil {
		return nil, err
	}

	mf, err := strconv.ParseFloat(result.Maker, 64)
	if err != nil {
		return nil, errors.WithMessagef(err, "parse maker fee %s fail", result.Maker)
	}

	tf, err := strconv.ParseFloat(result.Taker, 64)
	if err != nil {
		return nil, errors.WithMessagef(err, "parse taker fee %s fail", result.Taker)
	}

	return []exchange.TradeFee{
		{Symbol: nil, Maker: mf, Taker: tf},
	}, nil
}
