package deribit

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	OptionSymbol struct {
		*exchange.BaseOptionSymbol
	}

	SpotSymbol struct {
		*exchange.BaseSpotSymbol
	}

	SwapSymbol struct {
		*exchange.BaseSwapSymbol
	}

	FuturesSymbol struct {
		*exchange.BaseFutureSymbol
	}
)

const (
	OptionTypeCall = "call"
	OptionTypePut  = "put"
	timeLayout     = "2Jan06"

	KindOption            = "option"
	KindFuture            = "future"
	SettlePeriodPerpetual = "perpetual"
	SettlePeriodWeek      = "week"
	SettlePeriodMonth     = "month"
)

var (
	opMap = map[string]exchange.OptionType{
		OptionTypeCall: exchange.OptionTypeCall,
		OptionTypePut:  exchange.OptionTypePut,
	}

	symbolMu  = sync.Mutex{}
	symbolMap = map[string]exchange.Symbol{}

	Currencies = []string{"BTC", "ETH"}
)

type (
	InstrumentResult struct {
		TickSize            decimal.Decimal `json:"tick_size"`
		TakerCommision      decimal.Decimal `json:"taker_commision"`
		MakerCommision      decimal.Decimal `json:"maker_commision"`
		Strike              decimal.Decimal `json:"strike"`
		SettlementPeriod    string          `json:"settlement_period"`
		QuoteCurrency       string          `json:"quote_currency"`
		BaseCurreny         string          `json:"base_currency"`
		MinTradeAmount      decimal.Decimal `json:"min_trade_amount"`
		Kind                string          `json:"kind"`
		IsActive            bool            `json:"is_active"`
		InstrumentName      string          `json:"instrument_name"`
		ExpirationTimestamp int64           `json:"expiration_timestamp"`
		CreationTimestamp   int64           `json:"creation_timestamp"`
		ContractSize        decimal.Decimal `json:"contract_size"`
		OptionType          string          `json:"option_type"`
	}
)

func Init(ctx context.Context, testNet bool) error {
	return initSymbol(ctx, testNet)
}

func (i *InstrumentResult) Symbol() (exchange.Symbol, error) {
	cfg := exchange.SymbolConfig{
		AmountMax:       decimal.Zero,
		AmountMin:       i.MinTradeAmount,
		PricePrecision:  i.TickSize,
		AmountPrecision: i.MinTradeAmount,
		ValueMin:        decimal.Zero,
	}
	var st time.Time
	if i.Kind != KindFuture || i.SettlementPeriod != SettlePeriodPerpetual {
		st = time.Unix(i.ExpirationTimestamp/1e3, 0)
	}
	if i.Kind == KindOption {
		var t exchange.OptionType
		if i.OptionType == OptionTypeCall {
			t = exchange.OptionTypeCall
		} else if i.OptionType == OptionTypePut {
			t = exchange.OptionTypePut
		} else {
			return nil, errors.Errorf("unkown option type %s", i.OptionType)
		}
		ret := &OptionSymbol{
			exchange.NewBaseOptionSymbol(i.BaseCurreny, st, i.Strike, t, cfg, i),
		}
		return ret, nil

	} else if i.Kind == KindFuture {
		if i.SettlementPeriod == SettlePeriodPerpetual {
			return &SwapSymbol{
				exchange.NewBaseSwapSymbolWithCfg(i.BaseCurreny, cfg, i),
			}, nil
		}

		var ft exchange.FutureType
		diff := time.Until(st)
		if i.SettlementPeriod == SettlePeriodWeek {
			if (diff / time.Second / 86400) < 7 {
				ft = exchange.FutureTypeCW
			} else {
				ft = exchange.FutureTypeNW
			}
		} else {
			if (diff / time.Second / 86400) < 31 {
				ft = exchange.FutureTypeCQ
			} else if (diff / time.Second / 86400) < 90 {
				ft = exchange.FutureTypeNQ
			} else {
				ft = exchange.FutureTypeNNQ
			}
		}
		return &FuturesSymbol{
			exchange.NewBaseFuturesSymbolWithCfg(i.BaseCurreny, st, ft, cfg, i),
		}, nil
	}
	return nil, errors.Errorf("unkown kind '%s'", i.Kind)
}

func (c *Client) FutureFetchInstruments(ctx context.Context, currency string, expired bool) ([]InstrumentResult, error) {
	return c.fetchInstruments(ctx, currency, expired, "future")
}

func (c *Client) OptionFetchInstruments(ctx context.Context, currency string, expired bool) ([]InstrumentResult, error) {
	return c.fetchInstruments(ctx, currency, expired, "option")
}

func (c *Client) fetchInstruments(ctx context.Context, currency string, expired bool, kind string) ([]InstrumentResult, error) {
	var ir []InstrumentResult
	param := map[string]interface{}{
		"currency": strings.ToUpper(currency),
		"kind":     kind,
		"expired":  expired,
	}
	if err := c.call(ctx, "public/get_instruments", param, &ir, false); err != nil {
		return nil, err
	}

	return ir, nil
}

func (c *Client) OptionSymbols(ctx context.Context, currency string) ([]exchange.OptionSymbol, error) {
	ir, err := c.fetchSymbols(ctx, currency, false, KindOption)
	if err != nil {
		return nil, err
	}

	ret := make([]exchange.OptionSymbol, len(ir))
	for i, inst := range ir {
		ret[i] = inst.(exchange.OptionSymbol)
	}
	return ret, nil
}

func (c *Client) OptionExpireSymbols(ctx context.Context, currency string) ([]exchange.OptionSymbol, error) {
	ir, err := c.fetchSymbols(ctx, currency, true, KindOption)
	if err != nil {
		return nil, err
	}

	ret := make([]exchange.OptionSymbol, len(ir))
	for i, inst := range ir {
		ret[i] = inst.(exchange.OptionSymbol)
	}
	return ret, nil
}

func (c *Client) FuturesSymbols(ctx context.Context, currency string) ([]exchange.FuturesSymbol, error) {
	ir, err := c.fetchSymbols(ctx, currency, false, KindFuture)
	if err != nil {
		return nil, err
	}
	ret := make([]exchange.FuturesSymbol, 0)
	for _, inst := range ir {
		s, ok := inst.(exchange.FuturesSymbol)
		if ok {
			ret = append(ret, s)
		}
	}
	return ret, nil
}

func (c *Client) SwapSymbol(ctx context.Context, currency string) (exchange.SwapSymbol, error) {
	symbol := strings.ToUpper(fmt.Sprintf("%s-%s", currency, SettlePeriodPerpetual))
	ret, err := getSymbol(symbol, reflect.TypeOf((*exchange.SwapSymbol)(nil)).Elem())
	if err != nil {
		return nil, err
	}
	return ret.(exchange.SwapSymbol), nil
}

func (c *Client) fetchSymbols(ctx context.Context, currency string, expired bool, kind string) ([]exchange.Symbol, error) {
	ir, err := c.fetchInstruments(ctx, currency, expired, kind)
	if err != nil {
		return nil, err
	}

	ret := make([]exchange.Symbol, len(ir))
	for i, inst := range ir {
		sym, err := inst.Symbol()
		if err != nil {
			return nil, err
		}
		ret[i] = sym
	}
	return ret, nil
}

func ParseOptionSymbol(sym string) (exchange.OptionSymbol, error) {
	ret, err := getSymbol(sym, reflect.TypeOf((*exchange.OptionSymbol)(nil)).Elem())
	if err != nil {
		return nil, err
	}
	v := ret.(exchange.OptionSymbol)
	return v, nil
}

func ParseFutureSymbol(sym string) (exchange.FuturesSymbol, error) {
	ret, err := getSymbol(sym, reflect.TypeOf((*exchange.FuturesSymbol)(nil)).Elem())
	if err != nil {
		return nil, err
	}
	return ret.(exchange.FuturesSymbol), nil
}

func ParseSymbol(sym string) (exchange.Symbol, error) {
	symbolMu.Lock()
	defer symbolMu.Unlock()
	ret, ok := symbolMap[sym]
	if !ok {
		return nil, errors.Errorf("bad symbol %s", sym)
	}
	return ret, nil
}

func getSymbol(sym string, exType reflect.Type) (exchange.Symbol, error) {
	symbolMu.Lock()
	defer symbolMu.Unlock()
	ret, ok := symbolMap[sym]
	if !ok {
		return nil, errors.Errorf("bad symbol %s", sym)
	}

	typ := reflect.TypeOf(ret)
	if !typ.Implements(exType) {
		return nil, errors.Errorf("type mismatch typ=%s exType=%s", typ, exType)
	}
	return ret, nil
}

func initSymbol(ctx context.Context, testNet bool) error {
	var client *Client
	var newClientCB func(string, string, chan interface{}) *Client
	if testNet {
		newClientCB = NewTestWSClient
	} else {
		newClientCB = NewWSClient
	}
	client = newClientCB("", "", nil)
	if err := client.Run(ctx); err != nil {
		return err
	}
	if err := updateSymbolMap(ctx, client); err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(time.Hour)
		for {
			select {
			case <-client.Done():
				for {
					client = newClientCB("", "", nil)
					if err := client.Run(ctx); err != nil {
						time.Sleep(time.Second)
						continue
					}
					break
				}

			case <-ctx.Done():
				return

			case <-ticker.C:
				updateSymbolMap(ctx, client)
			}
		}
	}()
	return nil
}

func updateSymbolMap(ctx context.Context, client *Client) error {
	newMap := map[string]exchange.Symbol{}
	for _, c := range Currencies {
		symbols, err := client.OptionSymbols(ctx, c)
		if err != nil {
			return err
		}
		for _, s := range symbols {
			newMap[s.String()] = s
		}

		expired, err := client.OptionExpireSymbols(ctx, c)
		if err != nil {
			return err
		}

		for _, ex := range expired {
			newMap[ex.String()] = ex
		}

		futures, err := client.fetchSymbols(ctx, c, false, KindFuture)
		if err != nil {
			return err
		}
		for _, ex := range futures {
			newMap[ex.String()] = ex
		}
	}

	symbolMu.Lock()
	defer symbolMu.Unlock()
	symbolMap = newMap
	return nil
}

func (sym *OptionSymbol) String() string {
	typ := "P"
	if sym.Type() == exchange.OptionTypeCall {
		typ = "C"
	}
	st := strings.ToUpper(sym.SettleTime().Format(timeLayout))
	s, _ := sym.Strike().Float64()
	strike := int(s)
	return fmt.Sprintf("%s-%s-%d-%s", sym.Index(), st, strike, typ)
}

func ParseIndexSymbol(symbol string) (*SpotSymbol, error) {
	fields := strings.Split(symbol, "_")
	if len(fields) != 2 {
		return nil, errors.Errorf("invalid symbol '%s'", symbol)
	}

	return &SpotSymbol{
		exchange.NewBaseSpotSymbol(fields[0], fields[1], exchange.SymbolConfig{}, nil),
	}, nil
}

func (ss *SpotSymbol) String() string {
	return fmt.Sprintf("%s_%s", ss.Base(), ss.Quote())
}

func (ss *SwapSymbol) String() string {
	return fmt.Sprintf("%s-PERPETUAL", ss.Index())
}

func (fs *FuturesSymbol) String() string {
	return fmt.Sprintf("%s-%s", fs.Index(), strings.ToUpper(fs.SettleTime().Format(timeLayout)))
}
