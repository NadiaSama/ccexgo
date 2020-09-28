package swap

import (
	"context"
	"net/http"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/shopspring/decimal"
)

type (
	SwapSymbol struct {
		*exchange.BaseSwapSymbol
		BaseAsset       string
		MinAmount       decimal.Decimal
		PricePrecision  decimal.Decimal
		AmountPrecision decimal.Decimal
	}

	//only parse required field
	exchangeInfo struct {
		Timezone   string       `json:"timezone"`
		ServerTime int64        `json:"serverTime"`
		Symbols    []swapSymbol `json:"symbols"`
	}

	swapSymbol struct {
		Symbol                string          `json:"symbol"`
		Status                string          `json:"status"`
		MaintMarginPercent    decimal.Decimal `json:"maintMarginPercent"`
		RequiredMarginPercent decimal.Decimal `json:"requiredMarginPercent"`
		BaseAsset             string          `json:"baseAsset"`
		QuoteAsset            string          `json:"quoteAsset"`
		PricePrecision        int             `json:"pricePrecision"`
		QuantityPrecision     int             `json:"quantityPrecision"`
		BaseAssetPrecision    int             `json:"baseAssetPrecision"`
		QuotePrecision        int             `json:"quotePrecision"`
		UnderlyingType        string          `json:"COIN"`
		Filters               []filter        `json:"filters"`
		OrderTypes            []string        `json:"orderTypes"`
		TimeInForce           []string        `json:"timeInForce"`
	}

	filter struct {
		MinPrice          decimal.Decimal `json:"minPrice"`
		MaxPrice          decimal.Decimal `json:"maxPrice"`
		FilterType        string          `json:"filterType"`
		TickSize          decimal.Decimal `json:"tickSize"`
		StepSize          decimal.Decimal `json:"stepSize"`
		MaxQty            decimal.Decimal `json:"maxQty"`
		MinQty            decimal.Decimal `json:"minQty"`
		Limit             int             `json:"limit"`
		MultiplierDown    decimal.Decimal `json:"multiplierDown"`
		MultiplierUp      decimal.Decimal `json:"multiplierUp"`
		MultiplierDecimal decimal.Decimal `json:"multiplierDecimal"`
	}
)

const (
	priceFilter = "PRICE_FILTER"
	lotSize     = "LOT_SIZE"
)

func (rc *RestClient) loadSymbol(ctx context.Context) error {
	var info exchangeInfo
	if err := rc.Request(ctx, http.MethodGet, "/fapi/v1/exchangeInfo", nil, nil, false, &info); err != nil {
		return err
	}

	for _, sym := range info.Symbols {
		symbol := &SwapSymbol{
			BaseSwapSymbol: exchange.NewBaseSwapSymbol(sym.Symbol),
		}

		for _, f := range sym.Filters {
			if f.FilterType == priceFilter {
				symbol.PricePrecision = f.TickSize
				continue
			}
			if f.FilterType == lotSize {
				symbol.AmountPrecision = f.StepSize
				symbol.MinAmount = f.MinQty
				continue
			}
		}
		symbol.BaseAsset = sym.BaseAsset
		rc.symbols[symbol.String()] = symbol
	}
	return nil
}

func (rc *RestClient) Symbols() map[string]*SwapSymbol {
	return rc.symbols
}

func (s *SwapSymbol) String() string {
	return s.Index()
}
