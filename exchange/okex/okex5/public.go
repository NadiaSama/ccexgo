package okex5

import (
	"context"
	"net/http"
	"net/url"

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
	spotSymbols   = map[string]exchange.Symbol{}
	swapSymbols   = map[string]exchange.Symbol{}
	marginSymbols = map[string]exchange.Symbol{}
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

func InitSymbols(ctx context.Context) error {
	params := [][]interface{}{
		{InstTypeSpot, spotSymbols},
		{InstTypeSwap, swapSymbols},
	}

	for _, ps := range params {
		typ := ps[0].(InstType)
		sm := ps[1].(map[string]exchange.Symbol)
		if err := initSymbols(ctx, typ, sm); err != nil {
			return errors.WithMessagef(err, "init %s fail", typ)
		}
	}

	return nil
}

func ParseSpotSymbol(sym string) (exchange.SpotSymbol, error) {
	ret, ok := spotSymbols[sym]
	if !ok {
		return nil, errors.Errorf("symbol '%s' not support", sym)
	}
	return ret.(exchange.SpotSymbol), nil
}

func ParseMarginSymbol(sym string) (exchange.MarginSymbol, error) {
	ret, ok := marginSymbols[sym]
	if !ok {
		return nil, errors.Errorf("symbol '%s' not support", sym)
	}
	return ret.(exchange.MarginSymbol), nil
}

func ParseSwapSymbol(sym string) (exchange.SwapSymbol, error) {
	ret, ok := swapSymbols[sym]
	if !ok {
		return nil, errors.Errorf("symbol '%s' not support", sym)
	}
	return ret.(exchange.SwapSymbol), nil
}

func (it *Instrument) Config() (*exchange.SymbolConfig, error) {
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

	return &exchange.SymbolConfig{
		AmountPrecision: ap,
		PricePrecision:  pp,
		AmountMin:       amin,
	}, nil
}

func (ss *SpotSymbol) String() string {
	it := ss.Raw().(Instrument)
	return it.InstID
}

func (ms *MarginSymbol) String() string {
	it := ms.Raw().(Instrument)
	return it.InstID
}

func (ss *SwapSymbol) String() string {
	it := ss.Raw().(Instrument)
	return it.InstID
}

func initSymbols(ctx context.Context, typ InstType, sm map[string]exchange.Symbol) error {
	client := NewRestClient("", "", "")
	t := typ
	if t == InstTypeMargin {
		t = InstTypeSpot
	}
	instruments, err := client.Instruments(ctx, t)
	if err != nil {
		return err
	}

	for _, it := range instruments {
		cfg, err := it.Config()
		if err != nil {
			return errors.WithMessagef(err, "parse %+v fail", it)
		}

		var symbol exchange.Symbol
		if typ == InstTypeSpot {
			if it.Lever != "" {
				lv, err := decimal.NewFromString(it.Lever)
				if err != nil {
					return errors.WithMessagef(err, "invalid lever %s", it.Lever)
				}
				ms := &MarginSymbol{
					exchange.NewBaseMarginSymbol(it.BaseCcy, it.QuoteCcy, *cfg, lv, it),
				}
				marginSymbols[ms.String()] = ms
			}
			symbol = &SpotSymbol{
				exchange.NewBaseSpotSymbol(it.BaseCcy, it.QuoteCcy, *cfg, it),
			}
		} else if typ == InstTypeSwap {
			ctVal, err := decimal.NewFromString(it.CtVal)
			if err != nil {
				return errors.WithMessagef(err, "parse ctrVal fail %+v", it)
			}
			symbol = &SwapSymbol{
				exchange.NewBaseSwapSymbolWithCfg(it.Uly, ctVal, *cfg, it),
			}
		} else {
			return errors.Errorf("not support yet %+v", typ)
		}

		sm[symbol.String()] = symbol
	}
	return nil
}
