package ftx

import (
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	FillNotify struct {
		Fee       decimal.Decimal `json:"fee"`
		FeeRate   decimal.Decimal `json:"feeRate"`
		Future    string          `json:"future"`
		ID        int64           `json:"id"`
		Liquidity string          `json:"liquidity"`
		Market    string          `json:"market"`
		OrderID   int64           `json:"orderId"`
		TradeID   int64           `json:"tradeId"`
		Price     decimal.Decimal `json:"price"`
		Side      string          `json:"side"`
		Size      decimal.Decimal `json:"size"`
		Time      string          `json:"time"`
		Type      string          `json:"type"`
	}

	Fill struct {
		ID      exchange.OrderID
		Symbol  exchange.Symbol
		OrderID exchange.OrderID
		TradeID exchange.OrderID
		Side    exchange.OrderSide
		Price   decimal.Decimal
		Size    decimal.Decimal
		Time    time.Time
		Fee     decimal.Decimal
		FeeRate decimal.Decimal
	}
)

func parseFillInternal(notify *FillNotify) (*Fill, error) {
	side, ok := sideMap[notify.Side]
	if !ok {
		return nil, errors.Errorf("unknown order side '%s'", notify.Side)
	}
	symbol, err := ParseSymbol(notify.Future)
	if err != nil {
		return nil, errors.Errorf("unkown order symbol '%s'", notify.Future)
	}

	t, err := parseTime(notify.Time)
	if err != nil {
		return nil, err
	}

	return &Fill{
		ID:      exchange.NewIntID(notify.ID),
		Symbol:  symbol,
		OrderID: exchange.NewIntID(notify.OrderID),
		TradeID: exchange.NewIntID(notify.TradeID),
		Side:    side,
		Price:   notify.Price,
		Size:    notify.Size,
		Time:    t,
		Fee:     notify.Fee,
		FeeRate: notify.FeeRate,
	}, nil
}
