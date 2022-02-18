package ftx

import (
	"encoding/json"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
)

type (
	Trade struct {
		*exchange.TradeDS
		Symbol exchange.Symbol
	}
	TradeData struct {
		Action      string `json:"action"`
		Price       string `json:"price"`
		Size        string `json:"size"`
		Side        string `json:"side"`
		Liquidation bool   `json:"liquidation"`
		Time        string `json:"time"`
	}
	TradeChannel struct {
		symbol exchange.Symbol
	}
)

func NewTradesChannel(sym exchange.Symbol) exchange.Channel {
	return &TradeChannel{
		symbol: sym,
	}
}

func (t *TradeChannel) String() string {
	return t.symbol.String()
}

func NewTrade(sym exchange.Symbol) *Trade {
	return &Trade{Symbol: sym}
}

func (t *Trade) Init(cr *callResponse) (*exchange.Trade, error) {
	var td TradeData
	if err := json.Unmarshal(cr.Data, &td); err != nil {
		return nil, err
	}

	if td.Action != typePartial {
		return nil, errors.Errorf("bad action '%s' expect '%s'", td.Action, typePartial)
	}

	notify := td.Transfer(t.Symbol)
	t.TradeDS = exchange.NewTradeDS(notify)
	return t.Snapshot(), nil
}

func (t *Trade) Update(cr *callResponse) (*exchange.Trade, error) {
	var td TradeData
	if err := json.Unmarshal(cr.Data, &td); err != nil {
		return nil, err
	}
	if td.Action != typeUpdate {
		return nil, errors.Errorf("bad action '%s' expect '%s'", td.Action, typeUpdate)
	}

	notify := td.Transfer(t.Symbol)

	t.TradeDS.Update(notify)
	return &exchange.Trade{}, nil
}

func (t *TradeData) Transfer(sym exchange.Symbol) *exchange.TradeNotify {
	return &exchange.TradeNotify{
		Symbol:      sym,
		Price:       t.Price,
		Size:        t.Size,
		Side:        t.Side,
		Liquidation: t.Liquidation,
		Time:        t.Time,
	}
}
