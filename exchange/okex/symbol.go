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

	SwapSymbol struct {
		*exchange.BaseSwapSymbol
	}
)

func (c *RestClient) NewSpotSymbol(base, quote string) exchange.SpotSymbol {
	return &SpotSymbol{
		exchange.NewBaseSpotSymbol(strings.ToUpper(base), strings.ToUpper(quote)),
	}
}

func (ss *SpotSymbol) String() string {
	return fmt.Sprintf("%s-%s", ss.Base(), ss.Quote())
}

func (ss *RestClient) NewSwapSymbol(index string) exchange.SwapSymbol {
	return &SwapSymbol{
		exchange.NewBaseSwapSymbol(index),
	}
}

func (ss *SwapSymbol) String() string {
	return fmt.Sprintf("%s-SWAP", ss.Index())
}
