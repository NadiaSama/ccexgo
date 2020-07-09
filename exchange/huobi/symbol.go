package huobi

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
)

type (
	SpotSymbol struct {
		*exchange.BaseSpotSymbol
	}

	Symbol struct {
		BaseCurrency string `json:"base-currency"`
		QuoteCurreny string `json:"quote-currency"`
		Symbol       string `json:"symbol"`
	}
)

func (rc *RestClient) initSymbol(ctx context.Context) error {
	var syms []Symbol
	if err := rc.request(ctx, http.MethodGet, "/v1/common/symbols", nil, nil, false, &syms); err != nil {
		return err
	}

	for _, sym := range syms {
		rc.pair2Symbol[sym.Symbol] = rc.NewSpotSymbol(sym.BaseCurrency, sym.QuoteCurreny)
	}
	return nil
}

func (rc *RestClient) NewSpotSymbol(base, quote string) exchange.SpotSymbol {
	return &SpotSymbol{
		exchange.NewBaseSpotSymbol(strings.ToLower(base), strings.ToLower(quote)),
	}
}

func (rc *RestClient) ParseSpotSymbol(pair string) (exchange.SpotSymbol, error) {
	ret, ok := rc.pair2Symbol[pair]
	if !ok {
		return nil, errors.Errorf("unsupport pair %s", pair)
	}
	return ret, nil
}

func (ss *SpotSymbol) String() string {
	return fmt.Sprintf("%s%s", ss.Base(), ss.Quote())
}
