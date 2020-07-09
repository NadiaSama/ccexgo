package binance

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/NadiaSama/ccexgo/exchange"
)

type (
	SpotSymbol struct {
		*exchange.BaseSpotSymbol
	}

	//Symbol info
	//TODO: refactor Symbol add more info pricePrecison ...?
	Symbol struct {
		Symbol     string `json:"symbol"`
		BaseAsset  string `json:"baseAsset"`
		QuoteAsset string `json:"quoteAsset"`
	}

	ExchangeInfo struct {
		Symbols []Symbol `json:"symbols"`
	}
)

var (
	ErrPair = errors.New("symbol pair not support")
)

func (c *RestClient) initPair(ctx context.Context) error {
	var exInfo ExchangeInfo
	if err := c.request(ctx, "/api/v3/exchangeInfo", nil, false, &exInfo); err != nil {
		return err
	}

	for _, sym := range exInfo.Symbols {
		c.pair2Symbol[sym.Symbol] = c.NewSpotSymbol(sym.BaseAsset, sym.QuoteAsset)
	}
	return nil
}

func (c *RestClient) NewSpotSymbol(base, quote string) exchange.SpotSymbol {
	base = strings.ToUpper(base)
	quote = strings.ToUpper(quote)
	return &SpotSymbol{
		exchange.NewBaseSpotSymbol(base, quote),
	}
}

func (c *RestClient) ParseSpotSymbol(sym string) (exchange.SpotSymbol, error) {
	ret, ok := c.pair2Symbol[sym]
	if !ok {
		return nil, ErrPair
	}
	return ret, nil
}

func (ss *SpotSymbol) String() string {
	return fmt.Sprintf("%s%s", ss.Base(), ss.Quote())
}
