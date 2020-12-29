package spot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/huobi"
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

	SymbolResp struct {
		Status string   `json:"status"`
		Data   []Symbol `json:"data"`
	}
)

func (rc *RestClient) initSymbol(ctx context.Context) error {
	var resp SymbolResp
	if err := rc.Request(ctx, http.MethodGet, "/v1/common/symbols", nil, nil, false, &resp); err != nil {
		return err
	}

	if resp.Status != huobi.StatusOK {
		ret, _ := json.Marshal(&resp)
		return huobi.NewError(string(ret))
	}

	for _, sym := range resp.Data {
		rc.pair2Symbol[sym.Symbol] = rc.NewSpotSymbol(sym.BaseCurrency, sym.QuoteCurreny)
	}
	return nil
}

func (rc *RestClient) NewSpotSymbol(base, quote string) exchange.SpotSymbol {
	return &SpotSymbol{
		exchange.NewBaseSpotSymbol(strings.ToLower(base), strings.ToLower(quote), exchange.SymbolConfig{}, nil),
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
