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
	if len(syms) != 1 {
		return nil, errors.Errorf("only 1 symbol support")
	}

	var result TradeFee
	//if err := rc.request(ctx, http.MethodGet, "/api/spot/v3/trade_fee", map[string]string{"instrument_id": syms[0].String()}, nil, true, &result); err != nil {
	if err := rc.request(ctx, http.MethodGet, "/api/spot/v3/trade_fee", map[string]string{"category": "1"}, nil, true, &result); err != nil {
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

	if len(syms) == 0 {
		return []exchange.TradeFee{
			{Symbol: nil, Maker: mf, Taker: tf},
		}, nil
	}
	ret := make([]exchange.TradeFee, len(syms))
	for i := range syms {
		ret[i].Symbol = syms[i]
		ret[i].Maker = mf
		ret[i].Taker = tf
	}
	return ret, nil
}
