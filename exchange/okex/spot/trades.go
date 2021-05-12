package spot

import (
	"context"
	"net/http"

	"github.com/NadiaSama/ccexgo/exchange/okex"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	Fill struct {
		LedgerID     string          `json:"ledger_id"`
		TradeID      string          `json:"trade_id"`
		InstrumentID string          `json:"instrument_id"`
		Price        decimal.Decimal `json:"price"`
		Size         decimal.Decimal `json:"size"`
		OrderID      string          `json:"order_id"`
		ExecType     string          `json:"exec_type"`
		Timestamp    string          `json:"timestamp"`
		Fee          decimal.Decimal `json:"fee"`
		Side         string          `json:"side"`
		Currency     string          `json:"currency"`
	}
)

const (
	FillsEndPoint = "/api/spot/v3/fills"
)

func (rc *RestClient) Fills(ctx context.Context, instrumentID, orderID string, before, after, limit string) ([]Fill, error) {
	values := okex.FillsParam(instrumentID, orderID, before, after, limit)

	var ret []Fill
	if err := rc.Request(ctx, http.MethodGet, FillsEndPoint, values, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "fetch fills failed")
	}
	return ret, nil
}
