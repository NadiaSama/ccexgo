package okex

import (
	"fmt"
	"strings"

	"github.com/NadiaSama/ccexgo/exchange"
)

type (
	SpotSymbol struct {
		*exchange.BaseSpotSymbol
	}
)

func NewSpotSymbol(base, quote string) exchange.SpotSymbol {
	return &SpotSymbol{
		exchange.NewBaseSpotSymbol(strings.ToUpper(base), strings.ToUpper(quote), exchange.SymbolConfig{}, nil),
	}
}

func (ss *SpotSymbol) String() string {
	return fmt.Sprintf("%s-%s", ss.Base(), ss.Quote())
}
