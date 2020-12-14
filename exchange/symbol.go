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

	FutureType int
	//BaseSpotSymbol define common property of spot symbol
	BaseSpotSymbol struct {
		base  string
		quote string
	}

	FuturesSymbol interface {
		Symbol
		Index() string
		SettleTime() time.Time
		Type() FutureType
	}
	//BaseFutureSymbol define common property of future symbol
	BaseFutureSymbol struct {
		index      string
		settleTime time.Time
		typ        FutureType
	}

	SwapSymbol interface {
		Symbol
		Index() string
	}

	BaseSwapSymbol struct {
		index string
	}
)

const (
	OptionTypeCall = iota
	OptionTypePut

	//FutureTypeCW current week settle future
	FutureTypeCW
	//FutureTypeNW next week settle future
	FutureTypeNW
	//FutureTypeCQ current quarter settle future
	FutureTypeCQ
	//FutureTypeNQ next quarter settle future
	FutureTypeNQ
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

func NewBaseFutureSymbol(index string, st time.Time, typ FutureType) *BaseFutureSymbol {
	return &BaseFutureSymbol{
		index:      index,
		settleTime: st,
		typ:        typ,
	}
}

func (bfs *BaseFutureSymbol) Index() string {
	return bfs.index
}

func (bfs *BaseFutureSymbol) SettleTime() time.Time {
	return bfs.settleTime
}

func (bfs *BaseFutureSymbol) Type() FutureType {
	return bfs.typ
}

func NewBaseSwapSymbol(index string) *BaseSwapSymbol {
	return &BaseSwapSymbol{
		index: index,
	}
}
func (bsw *BaseSwapSymbol) Index() string {
	return bsw.index
}
