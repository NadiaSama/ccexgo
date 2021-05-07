package swap

import (
	"context"
	"net/http"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/binance"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	Trade struct {
		Symbol          string          `json:"symbol"`
		ID              int64           `json:"id"`
		OrderID         int64           `json:"orderId"`
		Price           decimal.Decimal `json:"price"`
		Qty             decimal.Decimal `json:"qty"`
		QuoteQty        decimal.Decimal `json:"quoteQty"`
		Commission      decimal.Decimal `json:"commission"`
		CommissionAsset string          `json:"commissionAsset"`
		RealizedPnl     decimal.Decimal `json:"realizedPnl"`
		Side            string          `json:"SIDE"`
		PositionSide    string          `json:"positionSide"`
		Time            int64           `json:"time"`
		Maker           bool            `json:"Maker"`
	}
)

const (
	UserTradesEndPoint = "/fapi/v1/userTrades"
)

func (rc *RestClient) UserTrades(ctx context.Context, req *exchange.TradeReqParam) ([]Trade, error) {
	value := binance.TradeParam(req)
	var ret []Trade
	if err := rc.Request(ctx, http.MethodGet, UserTradesEndPoint, value, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "fetch myTrades fail")
	}
	return ret, nil
}

func (t *Trade) Parse() (*exchange.Trade, error) {
	s, err := ParseSymbol(t.Symbol)
	if err != nil {
		return nil, err
	}

	var side exchange.OrderSide
	if t.Side == "BUY" {
		side = exchange.OrderSideBuy
	} else if t.Side == "SELL" {
		side = exchange.OrderSideSell
	} else {
		return nil, errors.Errorf("unkown side '%s'", t.Side)
	}

	return &exchange.Trade{
		ID:          exchange.NewIntID(t.ID),
		OrderID:     exchange.NewIntID(t.OrderID),
		Symbol:      s,
		Side:        side,
		Amount:      t.Qty,
		Price:       t.Price,
		Fee:         t.Commission,
		FeeCurrency: t.CommissionAsset,
		Time:        binance.ParseTimestamp(t.Time),
		Raw:         t,
	}, nil
}

func (rc *RestClient) Trades(ctx context.Context, req *exchange.TradeReqParam) ([]*exchange.Trade, error) {
	trades, err := rc.UserTrades(ctx, req)
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
