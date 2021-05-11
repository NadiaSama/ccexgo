package swap

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/okex"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	LedgerType string
	Ledger     struct {
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
)

func (rc *RestClient) Ledgers(ctx context.Context, instrumentID string, before, after, limit, typ string) ([]Ledger, error) {
	values := okex.FillsParam("", "", before, after, limit)
	if typ != "" {
		values.Add("type", typ)
	}
	endPoint := fmt.Sprintf("/api/swap/v3/accounts/%s/ledger", instrumentID)

	var ret []Ledger
	if err := rc.Request(ctx, http.MethodGet, endPoint, values, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "get ledger fail")
	}
	return ret, nil
}

func (rc *RestClient) Finance(ctx context.Context, req *exchange.FinanceReqParam) ([]*exchange.Finance, error) {
	var symbol string
	if req.Symbol != nil {
		symbol = req.Symbol.String()
	}

	var typ string
	if req.Type == exchange.FinanceTypeFunding {
		typ = "14"
	}
	ledgers, err := rc.Ledgers(ctx, symbol, req.StartID, req.EndID, strconv.Itoa(req.Limit), typ)
	if err != nil {
		return nil, err
	}
	ret := []*exchange.Finance{}

	for i := range ledgers {
		ledger := ledgers[i]
		f, err := ledger.Parse()
		if err != nil {
			return nil, errors.WithMessage(err, "parse ledger fail")
		}
		ret = append(ret, f)
	}
	return ret, nil
}

func (l *Ledger) Parse() (*exchange.Finance, error) {
	var (
		s   exchange.Symbol
		err error
		t   time.Time
	)
	if len(l.InstrumentID) != 0 {
		s, err = ParseSymbol(l.InstrumentID)
		if err != nil {
			return nil, err
		}
	}
	t, err = okex.ParseTime(l.Timestamp)
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
