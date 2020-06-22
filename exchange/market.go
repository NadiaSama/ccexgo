package exchange

import "time"

var (
	TimeNoExpire = time.Time{}
)

type (
	BaseMarket struct {
		symbol     Symbol
		expired    time.Time
		priceSize  float64
		amountSize float64
		makerFee   float64
		takerFee   float64
	}
)

func (m *BaseMarket) Expire() bool {
	if m.expired.Equal(TimeNoExpire) {
		return false
	}
	now := time.Now()
	return m.expired.Before(now)
}

func (m *BaseMarket) Symbol() Symbol {
	return m.symbol
}
