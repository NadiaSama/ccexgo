package swap

import (
	"context"
	"net/http"

	"github.com/NadiaSama/ccexgo/exchange/okex"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	OkexSymbol struct {
		InstrumentID        string          `json:"instrument_id"`
		Underlying          string          `json:"underlying"`
		BaseCurrency        string          `json:"base_currency"`
		QuoteCurrency       string          `json:"quote_currency"`
		SettlementCurrency  string          `json:"settlement_currency"`
		ContractVal         decimal.Decimal `json:"contract_val"`
		Listing             string          `json:"listing"`
		Delivery            string          `json:"delivery"`
		SizeIncrement       decimal.Decimal `json:"size_increment"`
		TickSize            decimal.Decimal `json:"tick_size"`
		IsInverse           string          `json:"is_inverse"`
		Category            string          `json:"category"`
		ContractValCurrency string          `json:"contract_val_currency"`
		CurrencyIndex       string          `json:"currency_index"`
	}

	Symbol struct {
		*exchange.BaseSwapSymbol
		instrumentID string
	}
)

const (
	symbolEndPoint = "/api/swap/v3/instruments"
)

var (
	symbolMap map[string]*Symbol = map[string]*Symbol{}
)

func Init(ctx context.Context) error {
	client := RestClient{
		okex.NewRestClient("", "", ""),
	}
	syms, err := client.Symbols(ctx)
	if err != nil {
		return err
	}

	for _, s := range syms {
		symbolMap[s.String()] = s.(*Symbol)
	}
	return nil
}

//Symbols return swap symbol
func (rc *RestClient) Symbols(ctx context.Context) ([]exchange.SwapSymbol, error) {
	var oss []OkexSymbol

	if err := rc.Request(ctx, http.MethodGet, symbolEndPoint, nil, nil, false, &oss); err != nil {
		return nil, err
	}

	ret := make([]exchange.SwapSymbol, len(oss))
	for i, os := range oss {
		s, err := os.Parse()
		if err != nil {
			return nil, err
		}

		ret[i] = s
	}
	return ret, nil
}

func (os *OkexSymbol) Parse() (*Symbol, error) {
	cfg := exchange.SymbolConfig{
		PricePrecision:  os.TickSize,
		AmountPrecision: os.SizeIncrement,
	}
	return &Symbol{
		exchange.NewBaseSwapSymbolWithCfg(os.Underlying, os.ContractVal, cfg, os),
		os.InstrumentID,
	}, nil
}
func (s *Symbol) String() string {
	return s.instrumentID
}

func ParseSymbol(symbol string) (exchange.SwapSymbol, error) {
	sym, ok := symbolMap[symbol]
	if !ok {
		return nil, errors.Errorf("unkown symbol=%s", symbol)
	}

	return sym, nil
}
