package exchange

type (
	OrderElem struct {
		Price  float64
		Amount float64
	}

	//OrderBookNotify change of current orderbook
	//OrderElem.Amount == 0 means delete
	OrderBookNotify struct {
		Symbol string
		Bids   []OrderElem
		Asks   []OrderElem
	}
)
