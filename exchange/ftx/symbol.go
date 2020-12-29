package ftx

import (
	"context"
	"fmt"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
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
	typeSwap   = "perpetual"
	typeSpot   = "spot"
)

func (rc *RestClient) initFutureSymbol(ctx context.Context) error {
	infos, err := rc.Futures(ctx)
	if err != nil {
		return err
	}
	for _, info := range infos {
		if !info.Enabled {
			continue
		}

		if info.Type == typeFuture && !info.Expired {
			name := info.Name
			st, err := time.Parse("2006-01-02T15:04:05Z07:00", info.Expiry)
			if err != nil {
				return errors.WithMessagef(err, "bad expire time '%s'", info.Expiry)
			}
			var typ exchange.FutureType
			now := time.Now()
			if st.Sub(now).Hours() > 3*30*24 {
				typ = exchange.FutureTypeNQ
			} else {
				typ = exchange.FutureTypeCQ
			}
			rc.symbols[name] = newFutureSymbol(info.Underlying, st, typ)
			continue
		}

		if info.Type == typeSwap {
			name := info.Name
			rc.symbols[name] = newSwapSymbol(info.Underlying)
			continue
		}
	}
	return nil
}

func (rc *RestClient) initSpotSymbol(ctx context.Context) error {
	markets, err := rc.Markets(ctx)
	if err != nil {
		return err
	}

	for _, m := range markets {
		if m.Type != typeSpot {
			continue
		}
		rc.symbols[m.Name] = newSpotSymbol(m.BaseCurrency, m.QuoteCurrency, &m)
	}
	return nil
}

func (rc *RestClient) ParseSymbol(symbol string) (exchange.Symbol, error) {
	sym, ok := rc.symbols[symbol]
	if !ok {
		return nil, errors.Errorf("unkown future symbol '%s'", symbol)
	}
	return sym, nil
}

func (rc *RestClient) ParseFutureSymbol(symbol string) (exchange.FuturesSymbol, error) {
	sym, err := rc.ParseSymbol(symbol)
	if err != nil {
		return nil, err
	}

	ret, ok := sym.(exchange.FuturesSymbol)
	if !ok {
		return nil, errors.Errorf("bad symbol for '%s'", symbol)
	}

	return ret, nil
}

func (rc *RestClient) ParseSwapSymbol(symbol string) (exchange.SwapSymbol, error) {
	sym, err := rc.ParseSymbol(symbol)
	if err != nil {
		return nil, err
	}

	ret, ok := sym.(exchange.SwapSymbol)
	if !ok {
		return nil, errors.Errorf("bad symbol for '%s'", symbol)
	}
	return ret, nil
}

func newSpotSymbol(base string, quote string, m *Market) *SpotSymbol {
	return &SpotSymbol{
		exchange.NewBaseSpotSymbol(base, quote, exchange.SymbolConfig{}, m),
	}
}

func (ss *SpotSymbol) String() string {
	return fmt.Sprintf("%s/%s", ss.Base(), ss.Quote())
}

func newFutureSymbol(index string, st time.Time, typ exchange.FutureType) *FuturesSymbol {
	return &FuturesSymbol{
		exchange.NewBaseFutureSymbol(index, st, typ),
	}
}

func (fs *FuturesSymbol) String() string {
	st := fs.SettleTime()
	return fmt.Sprintf("%s-%s", fs.Index(), st.Format("0102"))
}

func newSwapSymbol(index string) *SwapSymbol {
	return &SwapSymbol{
		exchange.NewBaseSwapSymbol(index),
	}
}

func (fs *SwapSymbol) String() string {
	return fmt.Sprintf("%s-PERP", fs.Index())
}
