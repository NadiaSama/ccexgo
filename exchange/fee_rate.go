package exchange

import "github.com/shopspring/decimal"

type (
	//TradeFee for the symbol
	TradeFee struct {
		Symbol Symbol
		Maker  decimal.Decimal
		Taker  decimal.Decimal
		Raw    interface{}
	}
)
