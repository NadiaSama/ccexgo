package ftx

import (
	"github.com/NadiaSama/ccexgo/exchange"
)

type (
	Trade struct {
		ID          int     `json:"id"`
		Price       float64 `json:"price"`
		Size        float64 `json:"size"`
		Side        string  `json:"side"`
		Liquidation bool    `json:"liquidation"`
		Time        string  `json:"time"`
	}
	TradeChannel struct {
		symbol exchange.Symbol
	}
	TradeNotify struct {
		ID          int     `json:"id"`
		Price       float64 `json:"price"`
		Size        float64 `json:"size"`
		Side        string  `json:"side"`
		Liquidation bool    `json:"liquidation"`
		Time        string  `json:"time"`
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

func parseTradesInternal(notify []*TradeNotify) ([]*Trade, error) {

	trades := make([]*Trade, len(notify))
	for k, v := range notify {
		trades[k] = &Trade{
			Price:       v.Price,
			Size:        v.Size,
			Side:        v.Side,
			Liquidation: v.Liquidation,
			Time:        v.Time,
		}
	}
	return trades, nil
}
