package exchange

import (
	"errors"
	"time"
)

const (
	OptionTypeCall = iota
	OptionTypePut
)

var (
	ErrBadSymbol = errors.New("bad symbol")
)

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

	BaseOptionSymbol struct {
		strike     float64
		index      string
		settleTime time.Time
		typ        OptionType
	}

	BaseSpotSymbol struct {
		base  string
		quote string
	}
)

func NewBaseOptionSymbol(strike float64, index string, st time.Time, typ OptionType) *BaseOptionSymbol {
	return &BaseOptionSymbol{
		strike:     strike,
		index:      index,
		settleTime: st,
		typ:        typ,
	}
}

func (bos *BaseOptionSymbol) Strike() float64 {
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

func NewBaseSpotSymbol(base, quote string) *BaseSpotSymbol {
	return &BaseSpotSymbol{
		base:  base,
		quote: quote,
	}
}
func (bss *BaseSpotSymbol) Base() string {
	return bss.base
}

func (bss *BaseSpotSymbol) Quote() string {
	return bss.quote
}
