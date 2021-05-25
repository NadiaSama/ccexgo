package margin

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
		InstrumentID string          `json:"instrument_id"`
		Price        decimal.Decimal `json:"price"`
		Size         decimal.Decimal `json:"size"`
		OrderID      string          `json:"order_id"`
		Timestamp    string          `json:"timestamp"`
		ExecType     string          `json:"exec_type"`
		Fee          string          `json:"fee"`
		Side         string          `json:"side"`
		Currency     string          `json:"currency"`
	}
)

const (
	FillsEndPoint = "/api/margin/v3/fills"
)

func (rc *RestClient) Fills(ctx context.Context, orderID string, instrumentID string, before string, after string, limit string) ([]Fill, error) {
	var ret []Fill
	values := okex.FillsParam(instrumentID, orderID, before, after, limit)
	if err := rc.Request(ctx, http.MethodGet, FillsEndPoint, values, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "fetch fills fail")
	}
	return ret, nil
}
