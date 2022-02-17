package ftx

import "github.com/NadiaSama/ccexgo/exchange"

type (
	Ticker struct {
	}
	TickerChannel struct {
		symbol exchange.Symbol
	}
)

func NewTickersChannel(sym exchange.Symbol) exchange.Channel {
	return &TickerChannel{
		symbol: sym,
	}
}

func (t *TickerChannel) String() string {
	return t.symbol.String()
}
