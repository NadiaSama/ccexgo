package exchange

import (
	"time"
)

type (
	//Symbol is used to unit different exchange markets symbol serialize
	Symbol interface {
		String() string
	}

	SpotSymbol interface {
		Symbol
		Base() string
		Quote() string
	}

	OptionType int

	OptionSymbol interface {
		Symbol
		Strike() float64
		Index() string
		SettleTime() time.Time
		Type() OptionType
	}

	//BaseOptionSymbol define common property of option symbol
	BaseOptionSymbol struct {
		strike     float64
		index      string
		settleTime time.Time
		typ        OptionType
	}

	//BaseSpotSymbol define common property of spot symbol
	BaseSpotSymbol struct {
		base  string
		quote string
	}

	FuturesSymbol interface {
		Symbol
		Index() string
		SettleTime() time.Time
	}
	BaseFutureSymbol struct {
		index      string
		settleTime time.Time
	}
)

const (
	OptionTypeCall = iota
	OptionTypePut
)

func NewBaseOptionSymbol(index string, st time.Time, strike float64, typ OptionType) *BaseOptionSymbol {
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

func (ot OptionType) String() string {
	if ot == OptionTypeCall {
		return "CALL"
	} else if ot == OptionTypePut {
		return "PUT"
	} else {
		return "UNKOWN"
	}
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

func NewBaseFutureSymbol(index string, st time.Time) *BaseFutureSymbol {
	return &BaseFutureSymbol{
		index:      index,
		settleTime: st,
	}
}

func (bfs *BaseFutureSymbol) Index() string {
	return bfs.index
}

func (bfs *BaseFutureSymbol) SettleTime() time.Time {
	return bfs.settleTime
}
