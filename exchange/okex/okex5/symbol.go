package okex5

import (
	"context"
	"net/http"
	"net/url"
	"reflect"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	Instrument struct {
		InstType  InstType `json:"instType"`
		InstID    string   `json:"instId"`
		Uly       string   `json:"uly"`
		Category  string   `json:"category"`
		BaseCcy   string   `json:"baseCcy"`
		QuoteCcy  string   `json:"quoteCcy"`
		SettleCcy string   `json:"settleCcy"`
		CtVal     string   `json:"CtVal"`
		CtMul     string   `json:"CtMul"`
		CtValCcy  string   `json:"ctValCcy"`
		OptType   string   `json:"optType"`
		Stk       string   `json:"stk"`
		ListTime  string   `json:"listTime"`
		ExpTime   string   `json:"expTime"`
		Lever     string   `json:"lever"`
		TickSz    string   `json:"tickSz"`
		LotSz     string   `json:"lotSz"`
		MinSz     string   `json:"minSz"`
		CtType    string   `json:"ctType"`
		Alias     string   `json:"alias"`
		State     string   `json:"state"`
	}

	SpotSymbol struct {
		*exchange.BaseSpotSymbol
	}

	MarginSymbol struct {
		*exchange.BaseMarginSymbol
	}

	SwapSymbol struct {
		*exchange.BaseSwapSymbol
	}
)

const (
	InstrumentEndPoint = "/api/v5/public/instruments"
)

var (
	spotSymbolMap   = map[string]exchange.SpotSymbol{}
	swapSymbolMap   = map[string]exchange.SwapSymbol{}
	marginSymbolMap = map[string]exchange.MarginSymbol{}
)

func (rc *RestClient) Instruments(ctx context.Context, typ InstType) ([]Instrument, error) {
	var ret []Instrument
	v := url.Values{}
	v.Add("instType", string(typ))
	if err := rc.Request(ctx, http.MethodGet, InstrumentEndPoint, v, nil, false, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}

//Symbols return all spot + margin + swap symbols
func (rc *RestClient) Symbols(ctx context.Context) ([]exchange.Symbol, error) {
	var ret []exchange.Symbol
	spots, err := rc.SpotSymbols(ctx)
	if err != nil {
		return nil, err
	}

	for _, s := range spots {
		ret = append(ret, s)
	}

	margins, err := rc.MarginSymbols(ctx)
	if err != nil {
		return nil, err
	}
	for _, m := range margins {
		ret = append(ret, m)
	}

	swaps, err := rc.SwapSymbols(ctx)
	if err != nil {
		return nil, err
	}
	for _, s := range swaps {
		ret = append(ret, s)
	}
	return ret, nil
}

func (rc *RestClient) SpotSymbols(ctx context.Context) ([]exchange.SpotSymbol, error) {
	symbols, err := rc.symbols(ctx, InstTypeSpot)
	if err != nil {
		return nil, err
	}

	return symbols.([]exchange.SpotSymbol), nil
}

func (rc *RestClient) MarginSymbols(ctx context.Context) ([]exchange.MarginSymbol, error) {
	symbols, err := rc.symbols(ctx, InstTypeMargin)
	if err != nil {
		return nil, err
	}
	return symbols.([]exchange.MarginSymbol), nil
}

func (rc *RestClient) SwapSymbols(ctx context.Context) ([]exchange.SwapSymbol, error) {
	symbols, err := rc.symbols(ctx, InstTypeSwap)
	if err != nil {
		return nil, err
	}
	return symbols.([]exchange.SwapSymbol), nil
}

func (rc *RestClient) symbols(ctx context.Context, it InstType) (interface{}, error) {
	its, err := rc.Instruments(ctx, it)
	if err != nil {
		return nil, err
	}
	var (
		arr interface{}
	)
	switch it {
	case InstTypeSpot:
		arr = make([]exchange.SpotSymbol, 0, len(its))

	case InstTypeMargin:
		arr = make([]exchange.MarginSymbol, 0, len(its))

	case InstTypeSwap:
		arr = make([]exchange.SwapSymbol, 0, len(its))
	}

	arrValue := reflect.ValueOf(arr)

	for i := range its {
		it := its[i]
		sym, err := it.Parse()
		if err != nil {
			return nil, err
		}

		arrValue = reflect.Append(arrValue, reflect.ValueOf(sym))
	}

	return arrValue.Interface(), nil
}

func InitSymbols(ctx context.Context) error {
	return initSymbols(ctx, false)
}

func InitTestSymbols(ctx context.Context) error {
	return initSymbols(ctx, true)
}

func initSymbols(ctx context.Context, isTest bool) error {
	var rc *RestClient
	if isTest {
		rc = NewTestRestClient("", "", "")
	} else {
		rc = NewRestClient("", "", "")
	}

	spots, err := rc.SpotSymbols(ctx)
	if err != nil {
		return err
	}

	for _, s := range spots {
		spotSymbolMap[s.String()] = s
	}

	swaps, err := rc.SwapSymbols(ctx)
	if err != nil {
		return err
	}
	for _, s := range swaps {
		swapSymbolMap[s.String()] = s
	}

	margins, err := rc.MarginSymbols(ctx)
	if err != nil {
		return err
	}

	for _, m := range margins {
		marginSymbolMap[m.String()] = m
	}

	return nil
}

func ParseSpotSymbol(sym string) (exchange.SpotSymbol, error) {
	s, ok := spotSymbolMap[sym]
	if !ok {
		return nil, errors.Errorf("unsupport symbol '%s'", sym)
	}

	return s, nil
}

func ParseMarginSymbol(sym string) (exchange.MarginSymbol, error) {
	s, ok := marginSymbolMap[sym]
	if !ok {
		return nil, errors.Errorf("unsupport symbol '%s'", sym)
	}
	return s, nil
}

func ParseSwapSymbol(sym string) (exchange.SwapSymbol, error) {
	s, ok := swapSymbolMap[sym]
	if !ok {
		return nil, errors.Errorf("unsupport symbol '%s'", sym)
	}
	return s, nil
}

func (it *Instrument) Parse() (exchange.Symbol, error) {
	pp, err := decimal.NewFromString(it.TickSz)
	if err != nil {
		return nil, errors.WithMessage(err, "invalid tickSz")
	}

	ap, err := decimal.NewFromString(it.LotSz)
	if err != nil {
		return nil, errors.WithMessage(err, "invalid lotSz")
	}

	amin, err := decimal.NewFromString(it.MinSz)
	if err != nil {
		return nil, errors.WithMessage(err, "invalid minSz")
	}

	cfg := &exchange.SymbolConfig{
		AmountPrecision: ap,
		PricePrecision:  pp,
		AmountMin:       amin,
	}

	switch it.InstType {
	case InstTypeSpot:
		return &SpotSymbol{
			BaseSpotSymbol: exchange.NewBaseSpotSymbol(it.BaseCcy, it.QuoteCcy, *cfg, it),
		}, nil

	case InstTypeMargin:
		lv, err := decimal.NewFromString(it.Lever)
		if err != nil {
			return nil, errors.WithMessagef(err, "invalid lever %s", it.Lever)
		}
		return &MarginSymbol{
			exchange.NewBaseMarginSymbol(it.BaseCcy, it.QuoteCcy, *cfg, lv, it),
		}, nil

	case InstTypeSwap:
		ctVal, err := decimal.NewFromString(it.CtVal)
		if err != nil {
			return nil, errors.WithMessagef(err, "parse ctrVal fail %+v", it)
		}
		return &SwapSymbol{
			exchange.NewBaseSwapSymbolWithCfg(it.Uly, ctVal, *cfg, it),
		}, nil

	default:
		return nil, errors.Errorf("unsupport instType '%s'", it.InstType)
	}
}

func (ss *SpotSymbol) String() string {
	it := ss.Raw().(*Instrument)
	return it.InstID
}

func (ms *MarginSymbol) String() string {
	it := ms.Raw().(*Instrument)
	return it.InstID
}

func (ss *SwapSymbol) String() string {
	it := ss.Raw().(*Instrument)
	return it.InstID
}
