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
		Strike() int
		Index() string
		SettleTime() time.Time
		Type() OptionType
	}

	BaseOptionSymbol struct {
		strike     int
		index      string
		settleTime time.Time
		typ        OptionType
	}
)

func NewBaseOptionSymbol(strike int, index string, st time.Time, typ OptionType) *BaseOptionSymbol {
	return &BaseOptionSymbol{
		strike:     strike,
		index:      index,
		settleTime: st,
		typ:        typ,
	}
}

func (bos *BaseOptionSymbol) Strike() int {
	return bos.strike
}
func (bos *BaseOptionSymbol) Index() string {
	return bos.index
}
func (bos *BaseOptionSymbol) SettleTime() time.Time {
	return bos.settleTime
}
func (bos *BaseOptionSymbol) Type() OptionType {
	return bos.typ
}

const (
	OptionTypeCall = iota
	OptionTypePut
)
