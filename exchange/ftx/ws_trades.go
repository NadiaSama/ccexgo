package ftx

import "github.com/NadiaSama/ccexgo/exchange"

type (
	Trade struct {
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

const (
	TradesChannel = "trades"
)

func NewTradesChannel(sym exchange.Symbol) exchange.Channel {
	return &TradeChannel{
		symbol: sym,
	}
}

func (t *TradeChannel) String() string {
	return t.symbol.String()
}
