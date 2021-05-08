package exchange

import (
	"time"

	"github.com/shopspring/decimal"
)

type (
	Trade struct {
		ID          OrderID
		OrderID     OrderID
		Symbol      Symbol
		Price       decimal.Decimal
		Amount      decimal.Decimal
		Fee         decimal.Decimal
		FeeCurrency string
		Time        time.Time
		Side        OrderSide
		Raw         interface{}
	}

	TradeReqParam struct {
		Symbol    Symbol
		StartTime time.Time
		EndTime   time.Time
		StartID   interface{}
		EndID     interface{}
		Limit     int
	}
)
