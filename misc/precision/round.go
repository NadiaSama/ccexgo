package precision

import (
	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/shopspring/decimal"
)

func RoundAmountFloat(symbol exchange.Symbol, amount float64) float64 {
	return roundFloat(symbol.AmountPrecision(), amount)
}

func RoundAmount(symbol exchange.Symbol, amount decimal.Decimal) decimal.Decimal {
	return roundDecimal(symbol.AmountPrecision(), amount)
}

func RoundPrice(symbol exchange.Symbol, price decimal.Decimal) decimal.Decimal {
	return roundDecimal(symbol.PricePrecision(), price)
}

func RoundPriceFloat(symbol exchange.Symbol, price float64) float64 {
	return roundFloat(symbol.PricePrecision(), price)
}

func roundFloat(exp decimal.Decimal, val float64) float64 {
	ret := roundDecimal(exp, decimal.NewFromFloat(val))
	r, _ := ret.Float64()
	return r
}

func roundDecimal(exp decimal.Decimal, val decimal.Decimal) decimal.Decimal {
	e := exp.Exponent()
	return val.Round(e * -1)
}
