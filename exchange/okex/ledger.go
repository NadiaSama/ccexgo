package okex

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	LedgerType string

	//Leder response for margin, swap
	Ledger struct {
		LedgerID     string          `json:"ledger_id"`
		Amount       decimal.Decimal `json:"amount"`
		Type         LedgerType      `json:"type"`
		Fee          decimal.Decimal `json:"fee"`
		Timestamp    string          `json:"timestamp"`
		InstrumentID string          `json:"instrument_id"`
		Currency     string          `json:"currency"`
		Details      interface{}     `json:"details"`
		OrderID      string          `json:"order_id"`
		From         string          `json:"from"`
		To           string          `json:"to"`
		Balance      decimal.Decimal `json:"balance"`
	}

	Pcb func(string) (exchange.Symbol, error)

	Client interface {
		Request(ctx context.Context, method string, endPoint string, param url.Values,
			body io.Reader, sign bool, dst interface{}) error
	}
)

func FetchLedgers(ctx context.Context, cl Client, path string, before, after, limit, typ string) ([]Ledger, error) {
	var ret []Ledger
	values := FillsParam("", "", before, after, limit)
	if typ != "" {
		values.Add("type", typ)
	}
	if err := cl.Request(ctx, http.MethodGet, path, values, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "fetch ledgers fail")
	}
	return ret, nil
}

func (l *Ledger) Parse(pcb Pcb) (*exchange.Finance, error) {
	var (
		s   exchange.Symbol
		err error
		t   time.Time
	)
	if len(l.InstrumentID) != 0 {
		s, err = pcb(l.InstrumentID)
		if err != nil {
			return nil, err
		}
	}
	t, err = ParseTime(l.Timestamp)
	if err != nil {
		return nil, err
	}

	return &exchange.Finance{
		ID:       l.LedgerID,
		Time:     t,
		Amount:   l.Amount,
		Currency: l.Currency,
		Type:     l.Type.Parse(),
		Symbol:   s,
		Raw:      *l,
	}, nil
}

func (typ LedgerType) Parse() exchange.FinanceType {
	if string(typ) == "funding" {
		return exchange.FinanceTypeFunding
	}
	return exchange.FinanceTypeOther
}
