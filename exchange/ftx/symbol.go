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
)

const (
	typeFuture = "future"
	typeSwap   = "perpetual"
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
			rc.symbols[name] = newFutureSymbol(info.Underlying, st)
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

func (rc *RestClient) ParseFutureSymbol(symbol string) (exchange.FuturesSymbol, error) {
	sym, ok := rc.symbols[symbol]
	if !ok {
		return nil, errors.Errorf("unkown future symbol '%s'", symbol)
	}

	ret, ok := sym.(exchange.FuturesSymbol)
	if !ok {
		return nil, errors.Errorf("bad symbol for '%s'", symbol)
	}

	return ret, nil
}

func (rc *RestClient) ParseSwapSymbol(symbol string) (exchange.SwapSymbol, error) {
	sym, ok := rc.symbols[symbol]
	if !ok {
		return nil, errors.Errorf("unkown swap symbol '%s'", symbol)
	}

	ret, ok := sym.(exchange.SwapSymbol)
	if !ok {
		return nil, errors.Errorf("bad symbol for '%s'", symbol)
	}
	return ret, nil
}

func newFutureSymbol(index string, st time.Time) *FuturesSymbol {
	return &FuturesSymbol{
		exchange.NewBaseFutureSymbol(index, st),
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
