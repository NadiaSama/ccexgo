package spot

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	Trade struct {
		Symbol          string          `json:"symbol"`
		ID              int64           `json:"id"`
		OrderID         int64           `json:"orderId"`
		OrderListID     int64           `json:"orderListId"`
		Price           decimal.Decimal `json:"price"`
		Qty             decimal.Decimal `json:"qty"`
		QuoteQty        decimal.Decimal `json:"quoteQty"`
		Commission      decimal.Decimal `json:"commission"`
		CommissionAsset string          `json:"commissionAsset"`
		Time            int64           `json:"time"`
		IsBuyer         bool            `json:"isBuyer"`
		IsMaker         bool            `json:"isMaker"`
		IsBestMatch     bool            `json:"isBestMatch"`
	}
)

const (
	MyTradesEndPoint = "/api/v3/myTrades"
)

func (rc *RestClient) MyTrades(ctx context.Context, req *exchange.TradeReqParam) ([]Trade, error) {
	var ret []Trade
	value := url.Values{}
	value.Add("symbol", req.Symbol.String())
	if !req.StartTime.IsZero() {
		value.Add("startTime", fmt.Sprintf("%d", req.StartTime.UnixNano()/1e6))
	}
	if !req.EndTime.IsZero() {
		value.Add("endTime", fmt.Sprintf("%d", req.EndTime.UnixNano()/1e6))
	}
	if req.StartID != "" {
		value.Add("fromId", req.StartID)
	}
	if req.Limit != 0 {
		value.Add("limit", fmt.Sprintf("%d", req.Limit))
	}

	if err := rc.Request(ctx, http.MethodGet, MyTradesEndPoint, value, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "fetch myTrades fail")
	}
	return ret, nil
}

func (rc *RestClient) Trades(ctx context.Context, req *exchange.TradeReqParam) ([]*exchange.Trade, error) {
	trades, err := rc.MyTrades(ctx, req)
	if err != nil {
		return nil, err
	}
	ret := []*exchange.Trade{}
	for i := range trades {
		trade := trades[i]
		t, err := trade.Parse()
		if err != nil {
			return nil, err
		}

		ret = append(ret, t)
	}
	return ret, nil
}

func (t *Trade) Parse() (*exchange.Trade, error) {
	s, err := ParseSpotSymbol(t.Symbol)
	if err != nil {
		return nil, err
	}

	ret := &exchange.Trade{
		ID:          exchange.NewIntID(t.ID),
		OrderID:     exchange.NewIntID(t.OrderID),
		Symbol:      s,
		Amount:      t.Qty,
		Price:       t.Price,
		Fee:         t.Commission,
		FeeCurrency: t.CommissionAsset,
		Time:        time.Unix(t.Time/1000, t.Time%1000*1e6),
		Raw:         t,
	}
	return ret, nil
}
