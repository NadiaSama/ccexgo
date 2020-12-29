package spot

import (
	"context"
	"fmt"
	"net/http"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/shopspring/decimal"
)

type (
	OkexSymbol struct {
		InstrumentID  string          `json:"instrument_id"`
		BaseCurrency  string          `json:"base_currency"`
		QuoteCurrency string          `json:"quote_currency"`
		MinSize       decimal.Decimal `json:"min_size"`
		SizeIncrement decimal.Decimal `json:"size_increment"`
		TickSize      decimal.Decimal `json:"tick_size"`
		Category      string          `json:"category"`
	}

	Symbol struct {
		*exchange.BaseSpotSymbol
	}
)

func (rc *RestClient) Symbols(ctx context.Context) ([]exchange.SpotSymbol, error) {
	var oss []OkexSymbol
	if err := rc.Request(ctx, http.MethodGet, "/api/spot/v3/instruments", nil, nil, false, &oss); err != nil {
		return nil, err
	}

	ret := make([]exchange.SpotSymbol, len(oss))
	for i, os := range oss {
		s, err := os.Transform()
		if err != nil {
			return nil, err
		}
		ret[i] = s
	}
	return ret, nil

}

func (os *OkexSymbol) Transform() (exchange.SpotSymbol, error) {
	cfg := exchange.SymbolConfig{
		AmountMin:       os.MinSize,
		AmountMax:       decimal.Zero,
		ValueMin:        decimal.Zero,
		PricePrecision:  os.TickSize,
		AmountPrecision: os.SizeIncrement,
	}

	return &Symbol{
		exchange.NewBaseSpotSymbol(os.BaseCurrency, os.QuoteCurrency, cfg, os),
	}, nil
}

func (s *Symbol) String() string {
	return fmt.Sprintf("%s-%s", s.Base(), s.Quote())
}
