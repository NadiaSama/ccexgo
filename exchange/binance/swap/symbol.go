package swap

import (
	"context"
	"net/http"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	SwapSymbol struct {
		*exchange.BaseSwapSymbol
		Symbol string
	}

	//TODO add additinoal field parse
	ExchangeInfo struct {
		Timezone   string   `json:"timezone"`
		ServerTime int64    `json:"serverTime"`
		Symbols    []Symbol `json:"symbols"`
	}

	Symbol struct {
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
		Filters               []Filter        `json:"filters"`
		OrderTypes            []string        `json:"orderTypes"`
		TimeInForce           []string        `json:"timeInForce"`
	}

	Filter struct {
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

var (
	symbolMap = map[string]exchange.SwapSymbol{}
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

func ParseSymbol(symbol string) (exchange.SwapSymbol, error) {
	sym, ok := symbolMap[symbol]
	if !ok {
		return nil, errors.Errorf("unsupport symbol %s", symbol)
	}
	return sym, nil
}

func (rc *RestClient) ExchangeInfo(ctx context.Context) (*ExchangeInfo, error) {
	var info ExchangeInfo
	if err := rc.Request(ctx, http.MethodGet, "/fapi/v1/exchangeInfo", nil, nil, false, &info); err != nil {
		return nil, errors.WithMessage(err, "fetch exchangeInfo fail")
	}
	return &info, nil
}

func (rc *RestClient) Symbols(ctx context.Context) ([]exchange.SwapSymbol, error) {
	info, err := rc.ExchangeInfo(ctx)
	if err != nil {
		return nil, err
	}

	var ret []exchange.SwapSymbol
	for i := range info.Symbols {
		sym := info.Symbols[i]
		s, err := sym.Parse()
		if err != nil {
			return nil, err
		}
		ret = append(ret, s)
	}
	return ret, nil
}

func (s *Symbol) Parse() (exchange.SwapSymbol, error) {
	return &SwapSymbol{
		exchange.NewBaseSwapSymbol(s.BaseAsset),
		s.Symbol,
	}, nil
}

func (s *SwapSymbol) String() string {
	return s.Symbol
}
