package exchange

type (
	//TradeFee for the symbol
	TradeFee struct {
		Symbol Symbol
		Maker float64
		Taker float64
	}
)