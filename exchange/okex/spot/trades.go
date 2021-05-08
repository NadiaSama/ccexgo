package spot

import (
	"context"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	Fills struct {
		LedgeID      string          `json:"ledger_id"`
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

func (rc *RestClient) Fills(ctx context.Context, orderID string, instrumentID string, before, after, limit string) ([]Fills, error) {
	values := url.Values{}
	if orderID != "" {
		values.Add("order_id", orderID)
	}
	if instrumentID != "" {
		values.Add("instrument_id", instrumentID)
	}

	if after != "" {
		values.Add("after", after)
	}

	if before != "" {
		values.Add("before", before)
	}

	if limit != "" {
		values.Add("limit", limit)
	}

	var ret []Fills
	if err := rc.Request(ctx, http.MethodGet, FillsEndPoint, values, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "fetch fills failed")
	}
	return ret, nil
}
