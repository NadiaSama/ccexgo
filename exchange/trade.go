package exchange

import (
	"time"

	"github.com/shopspring/decimal"
)

type (
	Trade struct {
		ID          string
		OrderID     string
		Symbol      Symbol
		Price       decimal.Decimal
		Amount      decimal.Decimal
		Fee         decimal.Decimal
		FeeCurrency string
		Time        time.Time
		Side        OrderSide
		IsMaker     bool
		Raw         interface{}
	}

	TradeReqParam struct {
		Symbol    Symbol
		StartTime time.Time
		EndTime   time.Time
		StartID   string
		EndID     string
		Limit     int
	}
)

func NewTradeReqParam() *TradeReqParam {
	return &TradeReqParam{}
}

func (trp *TradeReqParam) SetSymbol(s Symbol) *TradeReqParam {
	trp.Symbol = s
	return trp
}

func (trp *TradeReqParam) SetStartTime(st time.Time) *TradeReqParam {
	trp.StartTime = st
	return trp
}

func (trp *TradeReqParam) SetEndTime(et time.Time) *TradeReqParam {
	trp.EndTime = et
	return trp
}

func (trp *TradeReqParam) SetStartID(sid string) *TradeReqParam {
	trp.StartID = sid
	return trp
}

func (trp *TradeReqParam) SetEndID(eid string) *TradeReqParam {
	trp.EndID = eid
	return trp
}

func (trp *TradeReqParam) SetLimit(l int) *TradeReqParam {
	trp.Limit = l
	return trp
}
