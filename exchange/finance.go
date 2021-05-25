package exchange

import (
	"time"

	"github.com/shopspring/decimal"
)

type (
	FinanceType int
	Finance     struct {
		ID       string
		Time     time.Time
		Amount   decimal.Decimal
		Currency string
		Type     FinanceType
		Symbol   Symbol
		Raw      interface{}
	}

	FinanceReqParam struct {
		TradeReqParam
		Type FinanceType
	}
)

const (
	FinanceTypeOther FinanceType = iota
	FinanceTypeFunding
	FinanceTypeInterest
)
