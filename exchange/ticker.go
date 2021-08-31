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
		MarkPrice   decimal.Decimal
		Time        time.Time
		LastPrice   decimal.Decimal
		Raw         interface{}
	}
)
