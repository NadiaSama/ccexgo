package spot

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/misc/tconv"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	MatchResult struct {
		ID                int64           `json:"id"`
		MatchID           int64           `json:"match-id"`
		OrderID           int64           `json:"order-id"`
		TradeID           int64           `json:"trade-id"`
		CreatedAt         int64           `json:"created-at"`
		FilledAmount      decimal.Decimal `json:"filled-amount"`
		FilledFees        decimal.Decimal `json:"filled-fees"`
		FilledPoints      decimal.Decimal `json:"filled-points"`
		FeeCurrency       string          `json:"fee-currency"`
		Price             decimal.Decimal `json:"price"`
		Source            string          `json:"source"`
		Symbol            string          `json:"symbol"`
		Type              string          `json:"type"`
		Role              string          `json:"role"`
		FeeDeductCurrency string          `json:"fee-deduct-currency"`
		FeeDeductState    string          `json:"fee-deduct-state"`
	}

	Direct string
)

const (
	DirectNone Direct = ""
	DirectPrev Direct = "forward"
	DirectNext Direct = "next"

	MatchResutEndPoint = "/v1/order/matchresults"
)

func (rc *RestClient) MatchResults(ctx context.Context, symbol string, types []string, st int64, et int64, from string, d Direct, size int) ([]MatchResult, error) {
	if st != 0 && et != 0 && (et < st || et-st > 2*86400*1000) {
		return nil, errors.Errorf("invalid st=%d et=%d config", st, et)
	}

	values := url.Values{}
	values.Add("symbol", symbol)
	if len(types) != 0 {
		values.Add("types", strings.Join(types, ","))
	}
	if st != 0 {
		values.Add("start-time", fmt.Sprintf("%d", st))
	}
	if et != 0 {
		values.Add("end-time", fmt.Sprintf("%d", et))
	}
	if from != "" {
		values.Add("from", from)
	}
	if d != DirectNone {
		values.Add("direct", string(d))
	}
	if size != 0 {
		values.Add("size", strconv.Itoa(size))
	}

	var mr []MatchResult
	if err := rc.Request(ctx, http.MethodGet, MatchResutEndPoint, values, nil, true, &mr); err != nil {
		return nil, errors.WithMessage(err, "get matchresult failed")
	}
	return mr, nil
}

func (rc *RestClient) Trades(ctx context.Context, req *exchange.TradeReqParam) ([]*exchange.Trade, error) {
	if req.EndID != "" || req.StartID != "" {
		return nil, errors.Errorf("id is not support yet")
	}

	var (
		from string
		dire Direct
	)

	from = ""
	dire = DirectNone

	mr, err := rc.MatchResults(ctx, req.Symbol.String(), nil, tconv.Time2Milli(req.StartTime),
		tconv.Time2Milli(req.EndTime), from, dire, req.Limit)
	if err != nil {
		return nil, err
	}

	ret := []*exchange.Trade{}
	for _, m := range mr {
		t, err := m.Parse()
		if err != nil {
			return nil, errors.WithMessage(err, "parse match result fail")
		}
		ret = append(ret, t)
	}
	return ret, nil
}

func (mr MatchResult) Parse() (*exchange.Trade, error) {
	s, err := ParseSymbol(mr.Symbol)
	if err != nil {
		return nil, err
	}

	var (
		fee         decimal.Decimal
		feeCurrency string
		side        exchange.OrderSide
	)

	if strings.HasPrefix(mr.Type, "buy") {
		side = exchange.OrderSideBuy
	} else if strings.HasPrefix(mr.Type, "sell") {
		side = exchange.OrderSideSell
	} else {
		return nil, errors.Errorf("unkown match type %s", mr.Type)
	}
	if mr.FeeDeductCurrency != "" {
		fee = mr.FilledPoints
		feeCurrency = mr.FeeDeductCurrency
	} else {
		fee = mr.FilledFees
		if side == exchange.OrderSideBuy {
			feeCurrency = s.Base()
		} else {
			feeCurrency = s.Quote()
		}
	}

	return &exchange.Trade{
		Symbol:      s,
		OrderID:     strconv.FormatInt(mr.OrderID, 10),
		ID:          strconv.FormatInt(mr.ID, 10),
		Price:       mr.Price,
		Amount:      mr.FilledAmount,
		Fee:         fee,
		FeeCurrency: feeCurrency,
		Time:        tconv.Milli2Time(mr.CreatedAt),
		Side:        side,
		Raw:         mr,
	}, nil
}
