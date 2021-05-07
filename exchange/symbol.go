package exchange

import (
	"time"

	"github.com/shopspring/decimal"
)

type (
	//Symbol is used to unit different exchange markets symbol serialize
	Symbol interface {
		Raw() interface{}
		AmountPrecision() decimal.Decimal
		PricePrecision() decimal.Decimal
		AmountMax() decimal.Decimal
		AmountMin() decimal.Decimal
		ValueMin() decimal.Decimal
		String() string
	}

	RawMixin struct {
		raw interface{}
	}
	SpotSymbol interface {
		Symbol
		Base() string
		Quote() string
	}

	OptionType int

	OptionSymbol interface {
		Symbol
		Strike() decimal.Decimal
		Index() string
		SettleTime() time.Time
		Type() OptionType
	}

	//BaseSymbolProperty define common property of all kind symbol
	BaseSymbolProperty struct {
		pricePrecision  decimal.Decimal
		amountPrecision decimal.Decimal
		amountMin       decimal.Decimal
		amountMax       decimal.Decimal
		valueMin        decimal.Decimal //minuim price * amount
	}

	//SymbolConfig used to specific symbol property
	SymbolConfig struct {
		PricePrecision  decimal.Decimal
		AmountPrecision decimal.Decimal
		AmountMin       decimal.Decimal
		AmountMax       decimal.Decimal
		ValueMin        decimal.Decimal
	}

	//BaseOptionSymbol define common property of option symbol
	BaseOptionSymbol struct {
		RawMixin
		BaseSymbolProperty
		strike     decimal.Decimal
		index      string
		settleTime time.Time
		typ        OptionType
	}

	FutureType int
	//BaseSpotSymbol define common property of spot symbol
	BaseSpotSymbol struct {
		RawMixin
		BaseSymbolProperty
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
		RawMixin
		BaseSymbolProperty
		index      string
		settleTime time.Time
		typ        FutureType
	}

	SwapSymbol interface {
		Symbol
		Index() string
	}

	BaseSwapSymbol struct {
		RawMixin
		BaseSymbolProperty
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
	//FutureTypeNNQ next next quart settle future (deribit only)
	FutureTypeNNQ
)

func (r *RawMixin) Raw() interface{} {
	return r.raw
}

func (p *SymbolConfig) Property() BaseSymbolProperty {
	return BaseSymbolProperty{
		amountMax:       p.AmountMax,
		amountMin:       p.AmountMin,
		pricePrecision:  p.PricePrecision,
		amountPrecision: p.AmountPrecision,
		valueMin:        p.ValueMin,
	}
}

//AmountMin minium order amount
func (p *BaseSymbolProperty) AmountMin() decimal.Decimal {
	return p.amountMin
}

//AmountMax minum order amount zero means no limit
func (p *BaseSymbolProperty) AmountMax() decimal.Decimal {
	return p.amountMax
}

//PricePrecision return price precision value
func (p *BaseSymbolProperty) PricePrecision() decimal.Decimal {
	return p.pricePrecision
}

//AmountPrecision return amount precision value
func (p *BaseSymbolProperty) AmountPrecision() decimal.Decimal {
	return p.amountPrecision
}

//ValueMin return minium amount * price value zero means no limit
func (p *BaseSymbolProperty) ValueMin() decimal.Decimal {
	return p.valueMin
}

func NewBaseOptionSymbol(index string, st time.Time, strike decimal.Decimal, typ OptionType, prop SymbolConfig, raw interface{}) *BaseOptionSymbol {
	return &BaseOptionSymbol{
		RawMixin:           RawMixin{raw},
		BaseSymbolProperty: prop.Property(),
		strike:             strike,
		index:              index,
		settleTime:         st,
		typ:                typ,
	}
}

func (bos *BaseOptionSymbol) Strike() decimal.Decimal {
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

func NewBaseSpotSymbol(base, quote string, cfg SymbolConfig, raw interface{}) *BaseSpotSymbol {
	return &BaseSpotSymbol{
		RawMixin:           RawMixin{raw},
		BaseSymbolProperty: cfg.Property(),
		base:               base,
		quote:              quote,
	}
}
func (bss *BaseSpotSymbol) Base() string {
	return bss.base
}

func (bss *BaseSpotSymbol) Quote() string {
	return bss.quote
}

func NewBaseFuturesSymbolWithCfg(index string, st time.Time, typ FutureType, cfg SymbolConfig, raw interface{}) *BaseFutureSymbol {
	return &BaseFutureSymbol{
		RawMixin:           RawMixin{raw},
		BaseSymbolProperty: cfg.Property(),
		index:              index,
		settleTime:         st,
		typ:                typ,
	}
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

func NewBaseSwapSymbolWithCfg(index string, cfg SymbolConfig, raw interface{}) *BaseSwapSymbol {
	return &BaseSwapSymbol{
		RawMixin:           RawMixin{raw},
		BaseSymbolProperty: cfg.Property(),
		index:              index,
	}
}

func NewBaseSwapSymbol(index string) *BaseSwapSymbol {
	return &BaseSwapSymbol{
		index: index,
	}
}
func (bsw *BaseSwapSymbol) Index() string {
	return bsw.index
}

func Round(val decimal.Decimal, p decimal.Decimal) decimal.Decimal {
	times, _ := val.Div(p).Float64()
	return decimal.NewFromInt(int64(times)).Mul(p)
}
