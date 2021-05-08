package spot

import (
	"context"
	"net/http"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
)

type (
	SpotSymbol struct {
		*exchange.BaseSpotSymbol
		Symbol string
	}

	//TODO add precision
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

var (
	symbolMap map[string]exchange.SpotSymbol = map[string]exchange.SpotSymbol{}
)

func Init(ctx context.Context) error {
	client := NewRestClient("", "")
	symbols, err := client.Symbols(ctx)
	if err != nil {
		return err
	}
	for _, s := range symbols {
		symbolMap[s.String()] = s
	}
	return nil
}

func ParseSymbol(symbol string) (exchange.SpotSymbol, error) {
	ret, ok := symbolMap[symbol]
	if !ok {
		return nil, errors.Errorf("unsupport symbol %s", symbol)
	}
	return ret, nil
}

func (rc *RestClient) FetchSymbols(ctx context.Context) ([]Symbol, error) {
	var resp []Symbol
	if err := rc.Request(ctx, http.MethodGet, "/v1/common/symbols", nil, nil, false, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (rc *RestClient) Symbols(ctx context.Context) ([]exchange.SpotSymbol, error) {
	symbols, err := rc.FetchSymbols(ctx)
	if err != nil {
		return nil, err
	}

	ret := []exchange.SpotSymbol{}
	for _, symbol := range symbols {
		s, err := symbol.Parse()
		if err != nil {
			return nil, err
		}
		ret = append(ret, s)
	}
	return ret, nil
}

func (s *Symbol) Parse() (exchange.SpotSymbol, error) {
	return &SpotSymbol{
		exchange.NewBaseSpotSymbol(s.BaseCurrency, s.QuoteCurreny, exchange.SymbolConfig{}, s),
		s.Symbol,
	}, nil
}

func (ss *SpotSymbol) String() string {
	return ss.Symbol
}
