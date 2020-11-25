package okex

import (
	"fmt"
	"strings"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
)

type (
	SpotSymbol struct {
		*exchange.BaseSpotSymbol
	}

	SwapSymbol struct {
		*exchange.BaseSwapSymbol
	}
)

func NewSpotSymbol(base, quote string) exchange.SpotSymbol {
	return &SpotSymbol{
		exchange.NewBaseSpotSymbol(strings.ToUpper(base), strings.ToUpper(quote)),
	}
}

func (ss *SpotSymbol) String() string {
	return fmt.Sprintf("%s-%s", ss.Base(), ss.Quote())
}

func ParseSwapSymbol(sym string) (exchange.SwapSymbol, error) {
	if !strings.HasSuffix(sym, "-SWAP") {
		return nil, errors.Errorf("bad okex swap symbol %s", sym)
	}

	l := len(sym)
	return NewSwapSymbol(sym[:l-5]), nil
}
func NewSwapSymbol(index string) exchange.SwapSymbol {
	return &SwapSymbol{
		exchange.NewBaseSwapSymbol(index),
	}
}

func (ss *SwapSymbol) String() string {
	return fmt.Sprintf("%s-SWAP", ss.Index())
}
