package jex

import (
	"errors"

	"github.com/NadiaSama/ccexgo/exchange"
)

//ParseSymbol jex option symbol does not carry strike price
//$instrument_name$settle_time$type
func ParseSymbol(sym string) (exchange.Symbol, error) {
	return nil, errors.New("not impl")
}
