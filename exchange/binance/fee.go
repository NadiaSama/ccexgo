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

func (rc *RestClient) FeeRate(ctx context.Context, sym exchange.SpotSymbol) ([]exchange.TradeFee, error) {
	var param map[string]string
	if sym != nil {
		param = map[string]string{
			"symbol": sym.String(),
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
