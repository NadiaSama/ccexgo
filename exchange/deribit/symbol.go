package deribit

import (
	"context"
	"fmt"
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
)

const (
	OptionTypeCall = "call"
	OptionTypePut  = "put"
	timeLayout     = "2Jan06"
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

func Init(ctx context.Context) error {
	return initSymbol(ctx)
}

func (i *InstrumentResult) Symbol() (exchange.OptionSymbol, error) {
	var t exchange.OptionType
	if i.OptionType == OptionTypeCall {
		t = exchange.OptionTypeCall
	} else if i.OptionType == OptionTypePut {
		t = exchange.OptionTypePut
	} else {
		return nil, errors.Errorf("unkown option type %s", i.OptionType)
	}
	st := time.Unix(i.ExpirationTimestamp/1e3, 0)
	cfg := exchange.SymbolConfig{
		AmountMax:       decimal.Zero,
		AmountMin:       i.MinTradeAmount,
		PricePrecision:  i.TickSize,
		AmountPrecision: i.MinTradeAmount,
		ValueMin:        decimal.Zero,
	}
	ret := &OptionSymbol{
		exchange.NewBaseOptionSymbol(i.BaseCurreny, st, i.Strike, t, cfg),
	}

	return ret, nil
}

func (c *Client) OptionFetchInstruments(ctx context.Context, currency string) ([]InstrumentResult, error) {
	var ir []InstrumentResult
	param := map[string]interface{}{
		"currency": strings.ToUpper(currency),
		"kind":     "option",
		"expired":  false,
	}
	if err := c.call(ctx, "public/get_instruments", param, &ir, false); err != nil {
		return nil, err
	}

	return ir, nil
}

func (c *Client) OptionSymbols(ctx context.Context, currency string) ([]exchange.OptionSymbol, error) {
	ir, err := c.OptionFetchInstruments(ctx, currency)
	if err != nil {
		return nil, err
	}

	ret := make([]exchange.OptionSymbol, len(ir))
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
	symbolMu.Lock()
	defer symbolMu.Unlock()
	ret, ok := symbolMap[sym]
	if !ok {
		return nil, errors.Errorf("bad symbol %s", sym)
	}
	return ret.(exchange.OptionSymbol), nil
}

func initSymbol(ctx context.Context) error {
	client := NewWSClient("", "", nil)
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
				client = NewWSClient("", "", nil)
				client.Run(ctx)

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

func (c *Client) NewSpotSymbol(base, quote string) exchange.SpotSymbol {
	return newSpotSymbol(base, quote)
}

func (c *Client) ParseSpotSymbol(sym string) (exchange.SpotSymbol, error) {
	return parseSpotSymbol(sym)
}

func newSpotSymbol(base, quote string) exchange.SpotSymbol {
	return &SpotSymbol{
		exchange.NewBaseSpotSymbol(strings.ToLower(base), strings.ToLower(quote)),
	}
}

func (sym *SpotSymbol) String() string {
	return fmt.Sprintf("%s_%s", sym.Base(), sym.Quote())
}

func parseSpotSymbol(sym string) (exchange.SpotSymbol, error) {
	fields := strings.Split(strings.ToLower(sym), "_")
	if len(fields) != 2 {
		return nil, exchange.NewBadArg("bad spot symbol field len", len(fields))
	}
	return &SpotSymbol{
		exchange.NewBaseSpotSymbol(fields[0], fields[1]),
	}, nil
}
