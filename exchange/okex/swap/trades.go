package swap

import (
	"context"
	"net/http"

	"github.com/NadiaSama/ccexgo/exchange/okex"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	Fill struct {
		TradeID      string          `json:"trade_id"`
		FillID       string          `json:"fill_id"`
		InstrumentID string          `json:"instrument_id"`
		OrderID      string          `json:"order_id"`
		Price        decimal.Decimal `json:"price"`
		OrderQty     decimal.Decimal `json:"order_qty"`
		Fee          decimal.Decimal `json:"fee"`
		Timestamp    string          `json:"timestamp"`
		ExecType     string          `json:"exec_type"`
		Side         string          `json:"side"`
		OrderSide    string          `json:"order_side"`
		Type         string          `json:"type"`
	}
)

const (
	FillsEndPoint = "/api/swap/v3/fills"
)

func (rc *RestClient) Fills(ctx context.Context, instrumentID string, orderID string, before, after, limit string) ([]Fill, error) {
	values := okex.FillsParam(instrumentID, orderID, before, after, limit)
	var ret []Fill
	if err := rc.Request(ctx, http.MethodGet, FillsEndPoint, values, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "fetch fills fail")
	}
	return ret, nil
}
