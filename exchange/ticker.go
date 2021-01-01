package exchange

import (
	"time"

	"github.com/shopspring/decimal"
)

type (
	Ticker struct {
		Symbol      Symbol
		BestBid     decimal.Decimal
		BestBidSize decimal.Decimal
		BestAsk     decimal.Decimal
		BestAskSize decimal.Decimal
		Time        time.Time
		Raw         interface{}
	}
)
