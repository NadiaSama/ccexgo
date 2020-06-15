package exchange

import "time"

type (
	//Symbol is used to unit different exchange markets symbol
	Symbol interface {
		String() string
	}

	OptionType int

	OptionSymbol interface {
		Symbol
		Strike() float64
		Index() string
		SettleTime() time.Time
		Type() OptionType
	}
)

const (
	OptionTypeCall = iota
	OptionTypePut
)
