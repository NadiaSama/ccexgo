package spot

import (
	"context"
	"net/http"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	SpotSymbol struct {
		*exchange.BaseSpotSymbol
		Symbol string
	}

	//TODO add precision
	Symbol struct {
		BaseCurrency    string  `json:"base-currency"`
		QuoteCurreny    string  `json:"quote-currency"`
		Symbol          string  `json:"symbol"`
		MinOrderAmt     float64 `json:"min-order-amt"`
		MaxOrderAmt     float64 `json:"max-order-amt"`
		MinOrderValue   float64 `json:"min-order-value"`
		PricePrecision  int     `json:"price-precision"`
		AmountPrecision int     `json:"amount-precision"`
		ValuePrecision  int     `json:"value-precision"`
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
	cfg := exchange.SymbolConfig{
		AmountMin:       decimal.NewFromFloat(s.MinOrderAmt),
		AmountMax:       decimal.NewFromFloat(s.MaxOrderAmt),
		AmountPrecision: decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(s.AmountPrecision) * -1)),
		PricePrecision:  decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(s.PricePrecision) * -1)),
		ValuePrecision:  decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(s.ValuePrecision) * -1)),
	}
	return &SpotSymbol{
		exchange.NewBaseSpotSymbol(s.BaseCurrency, s.QuoteCurreny, cfg, s),
		s.Symbol,
	}, nil
}

func (ss *SpotSymbol) String() string {
	return ss.Symbol
}
