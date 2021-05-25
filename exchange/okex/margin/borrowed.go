package margin

import (
	"context"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	BorrowedStatus int

	Borrowed struct {
		BorrowID         string          `json:"borrow_id"`
		InstrumentID     string          `json:"instrument_id"`
		Currency         string          `json:"currency"`
		Timestamp        string          `json:"timestamp"`
		Amount           decimal.Decimal `json:"amount"`
		Interest         decimal.Decimal `json:"interest"`
		ReturnedAmount   decimal.Decimal `json:"returned_amount"`
		PaidInterest     decimal.Decimal `json:"paid_interest"`
		LastInterestTime string          `json:"last_interest_time"`
		ForceRepayTime   string          `json:"force_repay_time"`
		Rate             decimal.Decimal `json:"rate"`
	}
)

const (
	BorrowedStatusOpen BorrowedStatus = iota
	BorrowedStatusClose

	BorrowedEndPoint = "/api/margin/v3/accounts/borrowed"
)

func (rc *RestClient) Borrowed(ctx context.Context, status BorrowedStatus, before, after, limit string) ([]Borrowed, error) {
	var ret []Borrowed

	values := url.Values{}
	if status == BorrowedStatusOpen {
		values.Add("status", "0")
	} else if status == BorrowedStatusClose {
		values.Add("status", "1")
	}

	if before != "" {
		values.Add("before", before)
	}

	if after != "" {
		values.Add("after", after)
	}

	if limit != "" {
		values.Add("limit", limit)
	}

	if err := rc.Request(ctx, http.MethodGet, BorrowedEndPoint, values, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "fetch borrowed fail")
	}

	return ret, nil
}
