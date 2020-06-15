package exchange

type (
	OrderElem struct {
		Price  float64
		Amount float64
	}

	OrderBookNotify struct {
		Symbol string
		Bids   []OrderElem
		Asks   []OrderElem
	}
)
