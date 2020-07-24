package binance

import (
	"context"

	"github.com/NadiaSama/ccexgo/exchange"
)

type (
	TradeFee struct {
		Symbol string  `json:"symbol"`
		Maker  float64 `json:"maker"`
		Taker  float64 `json:"taker"`
	}

	TradeFeeResp struct {
		Success  bool       `json:"success"`
		TradeFee []TradeFee `json:"tradeFee"`
	}
)

func (rc *RestClient) FeeRate(ctx context.Context, syms ...exchange.Symbol) ([]exchange.TradeFee, error) {
	var param map[string]string
	if len(syms) != 0 && len(syms) != 1 {
		return nil, exchange.NewBadArg("unsupport multi symbols", syms)
	}
	if len(syms) == 1 {
		param = map[string]string{
			"symbol": syms[0].String(),
		}
	}
	var resp TradeFeeResp
	if err := rc.request(ctx, "/wapi/v3/tradeFee.html", param, true, &resp); err != nil {
		return nil, err
	}

	ret := make([]exchange.TradeFee, len(resp.TradeFee))
	for i, tf := range resp.TradeFee {
		sym, err := rc.ParseSpotSymbol(tf.Symbol)
		if err != nil {
			return nil, err
		}
		ret[i] = exchange.TradeFee{
			Symbol: sym,
			Maker:  tf.Maker,
			Taker:  tf.Taker,
		}
	}
	return ret, nil
}
