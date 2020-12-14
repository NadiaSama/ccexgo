package future

import (
	"context"
	"fmt"
	"net/http"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/okex"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

const (
	apiRawSymbolURI = "/api/futures/v3/instruments"
)

type (
	OkexSymbol struct {
		InstrumentID        string          `json:"instrument_id"`
		UnderlyingIndex     string          `json:"underlying_index"`
		QuoteCurrency       string          `json:"quote_currency"`
		TickSize            decimal.Decimal `json:"tick_size"`
		ContractVal         decimal.Decimal `json:"contract_val"`
		Listing             string          `json:"listing"`
		Delivery            string          `json:"delivery"`
		TradeIncrement      decimal.Decimal `json:"trade_increment"`
		Alias               string          `json:"alias"`
		Underlying          string          `json:"underlying"`
		BaseCurrency        string          `json:"base_currency"`
		SettlementCurrency  string          `json:"settlement_currency"`
		IsInverse           string          `json:"is_inverse"`
		ContractValCurrency string          `json:"contract_val_currency"`
		Category            string          `json:"category"`
	}

	Symbol struct {
		*exchange.BaseFutureSymbol
	}
)

var (
	alias2Type map[string]exchange.FutureType = map[string]exchange.FutureType{
		"this_week":  exchange.FutureTypeCW,
		"next_week":  exchange.FutureTypeNW,
		"quarter":    exchange.FutureTypeCQ,
		"bi_quarter": exchange.FutureTypeNQ,
	}
)

func (rc *RestClient) Symbols(ctx context.Context) ([]exchange.Symbol, error) {
	symbols, err := rc.RawSymbols(ctx)
	if err != nil {
		return nil, err
	}

	ret := make([]exchange.Symbol, len(symbols))
	for i, raw := range symbols {
		s, err := raw.Parse()
		if err != nil {
			return nil, err
		}
		ret[i] = s
	}
	return ret, nil
}

//RawSymbols return okex future symbols
func (rc *RestClient) RawSymbols(ctx context.Context) ([]OkexSymbol, error) {
	var symbols []OkexSymbol
	if err := rc.Request(ctx, http.MethodGet, apiRawSymbolURI, nil, nil, false, &symbols); err != nil {
		return nil, err
	}
	return symbols, nil
}

func (os *OkexSymbol) Parse() (*Symbol, error) {
	typ, ok := alias2Type[os.Alias]
	if !ok {
		return nil, errors.Errorf("unkown alias='%s'", os.Alias)
	}

	dt := fmt.Sprintf("%sT08:00:00.000Z")
	st, err := okex.ParseTime(dt)
	if err != nil {
		return nil, err
	}

	return &Symbol{
		exchange.NewBaseFutureSymbol(os.UnderlyingIndex, st, typ),
	}, nil
}

func (s *Symbol) String() string {
	st := s.SettleTime()
	return fmt.Sprintf("%s-%s", s.Index(), st.Format("060102"))
}
