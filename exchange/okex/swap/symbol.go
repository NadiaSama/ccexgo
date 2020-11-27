package swap

import (
	"context"
	"net/http"

	"github.com/shopspring/decimal"
)

type (
	Symbol struct {
		InstrumentID        string          `json:"instrument_id"`
		Underlying          string          `json:"underlying"`
		BaseCurrency        string          `json:"base_currency"`
		QuoteCurrency       string          `json:"quote_currency"`
		SettlementCurrency  string          `json:"settlement_currency"`
		ContractVal         decimal.Decimal `json:"contract_val"`
		Listing             string          `json:"listing"`
		Delivery            string          `json:"delivery"`
		SizeIncrement       string          `json:"size_increment"`
		TickSize            string          `json:"tick_size"`
		IsInverse           string          `json:"is_inverse"`
		Category            string          `json:"category"`
		ContractValCurrency string          `json:"contract_val_currency"`
		CurrencyIndex       string          `json:"currency_index"`
	}
)

const (
	symbolEndPoint = "/api/swap/v3/instruments"
)

//Symbols return swap symbol
func (rc *RestClient) Symbols(ctx context.Context) ([]Symbol, error) {
	var ret []Symbol

	if err := rc.Request(ctx, http.MethodGet, symbolEndPoint, nil, nil, false, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}
