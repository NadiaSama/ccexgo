package spot

import (
	"context"
	"net/http"
	"strconv"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/binance"
	"github.com/NadiaSama/ccexgo/misc/tconv"
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

func (rc *RestClient) MyTrades(ctx context.Context, symbol string, st int64, et int64, fid int64, limit int) ([]Trade, error) {
	var ret []Trade
	value := binance.TradeParam(symbol, st, et, fid, limit)
	if err := rc.Request(ctx, http.MethodGet, MyTradesEndPoint, value, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "fetch myTrades fail")
	}
	return ret, nil
}

func (rc *RestClient) Trades(ctx context.Context, req *exchange.TradeReqParam) ([]*exchange.Trade, error) {
	var fid int64
	if req.StartID != "" {
		var err error
		fid, err = strconv.ParseInt(req.StartID, 10, 64)
		if err != nil {
			return nil, errors.Errorf("invalid startID %s", req.StartID)
		}
	}

	trades, err := rc.MyTrades(ctx, req.Symbol.String(), tconv.Time2Milli(req.StartTime),
		tconv.Time2Milli(req.EndTime), fid, req.Limit)
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
	s, err := ParseSymbol(t.Symbol)
	if err != nil {
		return nil, err
	}

	var side exchange.OrderSide
	if t.IsBuyer {
		side = exchange.OrderSideBuy
	} else {
		side = exchange.OrderSideSell
	}

	ret := &exchange.Trade{
		ID:          strconv.FormatInt(t.ID, 10),
		OrderID:     strconv.FormatInt(t.OrderID, 10),
		Symbol:      s,
		Amount:      t.Qty,
		Price:       t.Price,
		Fee:         t.Commission.Neg(),
		FeeCurrency: t.CommissionAsset,
		Time:        tconv.Milli2Time(t.Time),
		Side:        side,
		IsMaker:     t.IsMaker,
		Raw:         *t,
	}
	return ret, nil
}
