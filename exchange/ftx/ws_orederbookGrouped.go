package ftx

import "github.com/NadiaSama/ccexgo/exchange"

type (
	OrderbookGrouped        struct{}
	OrderbookGroupedChannel struct {
		symbol exchange.Symbol
	}
)

func NewOrederbookGroupedChannel(sym exchange.Symbol) exchange.Channel {
	return &OrderbookGroupedChannel{
		symbol: sym,
	}
}

func (o *OrderbookGroupedChannel) String() string {
	return o.symbol.String()
}
