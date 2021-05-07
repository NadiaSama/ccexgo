package spot

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
		Symbol string
	}

	//Symbol info
	//TODO: refactor Symbol add more info pricePrecison ...
	Symbol struct {
		Symbol              string `json:"symbol"`
		BaseAsset           string `json:"baseAsset"`
		QuoteAsset          string `json:"quoteAsset"`
		BaseAssetPrecision  int    `json:"baseAssetPrecision"`
		QuoteAssetPrecision int    `json:"quoteAssetPrecision"`
	}

	ExchangeInfo struct {
		Symbols []Symbol `json:"symbols"`
	}
)

var (
	ErrPair   = errors.New("symbol pair not support")
	symbolMap = map[string]exchange.SpotSymbol{}
)

func Init(ctx context.Context) error {
	client := NewRestClient("", "")
	symbols, err := client.Symbols(ctx)
	if err != nil {
		return err
	}

	for i := range symbols {
		s := symbols[i]
		symbolMap[s.String()] = s
	}
	return nil
}

func (rc *RestClient) ExchangeInfo(ctx context.Context) (*ExchangeInfo, error) {
	var exInfo ExchangeInfo
	if err := rc.Request(ctx, http.MethodGet, "/api/v3/exchangeInfo", nil, nil, false, &exInfo); err != nil {
		return nil, errors.WithMessage(err, "get exchange info fail")

	}
	return &exInfo, nil
}

func (rc *RestClient) Symbols(ctx context.Context) ([]exchange.SpotSymbol, error) {
	exInfo, err := rc.ExchangeInfo(ctx)
	if err != nil {
		return nil, err
	}

	var ret []exchange.SpotSymbol
	for i := range exInfo.Symbols {
		sym := exInfo.Symbols[i]
		s, err := sym.Parse()
		if err != nil {
			return nil, err
		}
		ret = append(ret, s)
	}
	return ret, nil
}

func NewSymbol(base, quote string) exchange.SpotSymbol {
	base = strings.ToUpper(base)
	quote = strings.ToUpper(quote)
	return &SpotSymbol{
		exchange.NewBaseSpotSymbol(base, quote, exchange.SymbolConfig{}, nil),
		fmt.Sprintf("%s%s", base, quote),
	}
}

func ParseSymbol(sym string) (exchange.SpotSymbol, error) {
	ret, ok := symbolMap[sym]
	if !ok {
		return nil, ErrPair
	}
	return ret, nil
}

func (sym *Symbol) Parse() (exchange.SpotSymbol, error) {
	return &SpotSymbol{
		exchange.NewBaseSpotSymbol(sym.BaseAsset, sym.QuoteAsset, exchange.SymbolConfig{}, sym),
		sym.Symbol,
	}, nil
}

func (ss *SpotSymbol) String() string {
	return ss.Symbol
}
