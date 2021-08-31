package ftx

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	FuturesSymbol struct {
		*exchange.BaseFutureSymbol
	}

	SwapSymbol struct {
		*exchange.BaseSwapSymbol
	}

	SpotSymbol struct {
		*exchange.BaseSpotSymbol
	}
)

const (
	typeFuture = "future"
	typeMove   = "move"
	typeSwap   = "perpetual"
	typeSpot   = "spot"
)

var (
	mu        sync.Mutex
	symbolMap map[string]exchange.Symbol
)

func Init(ctx context.Context) error {
	if err := initSymbols(ctx); err != nil {
		return err
	}

	go func() {
		for {
			next := time.Now().Add(time.Hour)
			next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), 0, 5, 0, next.Location())
			time.Sleep(time.Until(next))
			initSymbols(ctx)
		}
	}()
	return nil
}

func ParseSymbol(symbol string) (exchange.Symbol, error) {
	mu.Lock()
	defer mu.Unlock()
	s, ok := symbolMap[symbol]
	if !ok {
		return nil, errors.Errorf("bad %s", symbol)
	}
	return s, nil
}

func initSymbols(ctx context.Context) error {
	d := make(map[string]exchange.Symbol, 0)
	if err := initSpotSymbols(ctx, d); err != nil {
		return err
	}

	if err := initFutureSymbol(ctx, d); err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()
	symbolMap = d
	return nil
}

func initSpotSymbols(ctx context.Context, dst map[string]exchange.Symbol) error {
	rc := NewRestClient("", "")
	markets, err := rc.Markets(ctx)
	if err != nil {
		return err
	}

	for i := range markets {
		m := markets[i]

		if m.Type != typeSpot {
			continue
		}

		s, err := m.ToSymbol()
		if err != nil {
			return errors.WithMessagef(err, "parse market %s fail", m.Name)
		}
		dst[s.String()] = s
	}
	return nil
}

func initFutureSymbol(ctx context.Context, dst map[string]exchange.Symbol) error {
	rc := NewRestClient("", "")
	futures, err := rc.Futures(ctx)
	if err != nil {
		return err
	}

	for i := range futures {
		info := futures[i]
		symbol, err := info.ToSymbol()
		if symbol == nil {
			continue
		}

		if err != nil {
			return errors.WithMessagef(err, "parse futures %s fail", info.Name)
		}
		dst[symbol.String()] = symbol
	}
	return nil
}

func (m *Market) ToSymbol() (exchange.Symbol, error) {
	cfg := exchange.SymbolConfig{
		PricePrecision:  decimal.NewFromFloat(m.PriceIncrement),
		AmountPrecision: decimal.NewFromFloat(m.SizeIncrement),
		AmountMin:       decimal.NewFromFloat(m.MinProvideSize),
	}
	switch m.Type {
	case typeSpot:
		return newSpotSymbol(m.BaseCurrency, m.QuoteCurrency, cfg, m), nil
	default:
		return nil, errors.Errorf("unsupport type %s", m.Type)
	}
}

func newSpotSymbol(base string, quote string, cfg exchange.SymbolConfig, m *Market) *SpotSymbol {
	return &SpotSymbol{
		exchange.NewBaseSpotSymbol(base, quote, cfg, m),
	}
}

func (ss *SpotSymbol) String() string {
	r := ss.Raw().(*Market)
	return r.Name
}

func (info *FutureInfo) ToSymbol() (exchange.Symbol, error) {
	cfg := exchange.SymbolConfig{
		AmountPrecision: decimal.NewFromFloat(info.SizeIncrement),
		PricePrecision:  decimal.NewFromFloat(info.PriceIncrement),
	}
	if info.Type == typeFuture || info.Type == typeMove {
		st, err := time.Parse("2006-01-02T15:04:05Z07:00", info.Expiry)
		if err != nil {
			return nil, errors.WithMessagef(err, "bad expire time '%s'", info.Expiry)
		}
		var typ exchange.FutureType
		now := time.Now()
		if st.Sub(now).Hours() > 3*30*24 {
			typ = exchange.FutureTypeNQ
		} else {
			typ = exchange.FutureTypeCQ
		}
		return &FuturesSymbol{
			BaseFutureSymbol: exchange.NewBaseFuturesSymbolWithCfg(info.Underlying, st, typ, cfg, info),
		}, nil
	}

	if info.Type == typeSwap {
		return &SwapSymbol{
			BaseSwapSymbol: exchange.NewBaseSwapSymbolWithCfg(info.Underlying, decimal.NewFromFloat(1.0), cfg, info),
		}, nil
	}
	//return nil, errors.Errorf("unkown type %s", info.Type)
	//unsupport type
	return nil, nil
}

func (fs *FuturesSymbol) String() string {
	r := fs.Raw().(*FutureInfo)
	return r.Name
}

func newSwapSymbol(index string) *SwapSymbol {
	return &SwapSymbol{
		exchange.NewBaseSwapSymbol(index),
	}
}

func (fs *SwapSymbol) String() string {
	return fmt.Sprintf("%s-PERP", fs.Index())
}
