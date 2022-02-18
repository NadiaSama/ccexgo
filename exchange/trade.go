package exchange

import (
	"time"

	"github.com/shopspring/decimal"
)

type (
	//Trade private trade
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

	PublicTrade struct {
		Symbol Symbol
		Price  decimal.Decimal
		Amount decimal.Decimal
		Side   OrderSide
		ID     string
		Time   time.Time
		Raw    interface{}
	}

	TradeReqParam struct {
		Symbol    Symbol
		StartTime time.Time
		EndTime   time.Time
		StartID   string
		EndID     string
		Limit     int
	}

	TradeNotify struct {
		Symbol      Symbol
		Price       string `json:"price"`
		Size        string `json:"size"`
		Side        string `json:"side"`
		Liquidation bool   `json:"liquidation"`
		Time        string `json:"time"`
	}

	TradeDS struct {
		symbol      Symbol
		updated     time.Time
		Price       string
		Size        string
		Side        string
		Liquidation bool
		Time        string
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

func NewTradeDS(notify *TradeNotify) *TradeDS {
	return &TradeDS{
		symbol:      notify.Symbol,
		Price:       notify.Price,
		Size:        notify.Size,
		Side:        notify.Side,
		Liquidation: notify.Liquidation,
		Time:        notify.Time,
	}
}

func (ts *TradeDS) Update(notify *TradeNotify) {
	ts.Price = notify.Price
	ts.Size = notify.Size
	ts.Side = notify.Side
	ts.Liquidation = notify.Liquidation
	ts.updated = time.Now()
}

func (ts *TradeDS) Snapshot() *Trade {
	price, _ := decimal.NewFromString(ts.Price)
	amount, _ := decimal.NewFromString(ts.Size)
	ret := &Trade{
		Symbol: ts.symbol,
		Price:  price,
		Amount: amount,
		Side:   toOrderSide(ts.Side),
	}
	return ret
}

func toOrderSide(side string) OrderSide {
	var o OrderSide
	switch side {
	case "buy":
		o = OrderSideBuy
	case "sell":
		o = OrderSideSell
	case "closeLong":
		o = OrderSideCloseLong
	case "closeShort":
		o = OrderSideCloseShort
	}
	return o
}
