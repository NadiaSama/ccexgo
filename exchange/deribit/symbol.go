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
	KindOptionCombo       = "option_combo"
	KindFutureCombo       = "future_combo"
	SettlePeriodPerpetual = "perpetual"
	SettlePeriodDay       = "day"
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
	rc        *RestClient

	Currencies = []string{"BTC", "ETH"}
)

func Init(ctx context.Context, testNet bool) error {
	if testNet {
		rc = NewTestRestClient("", "")
	} else {
		rc = NewRestClient("", "")
	}
	if err := updateSymbolMap(ctx, rc); err != nil {
		return err
	}
	return nil
}

//SymbolLoop start loop update symbol map periodly
func SymbolLoop(ctx context.Context) {
	ticker := time.NewTicker(time.Minute * 5)
	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			UpdateSymbolMap(ctx)
		}
	}
}

//UpdateSymbolMap rebuild symbol map. the symbol map update is leave to user
func UpdateSymbolMap(ctx context.Context) error {
	return updateSymbolMap(ctx, rc)
}

//Symbols get all symbols from symbol map
func Symbols() []exchange.Symbol {
	ret := []exchange.Symbol{}
	symbolMu.Lock()
	m := symbolMap
	symbolMu.Unlock()
	for _, v := range m {
		ret = append(ret, v)
	}
	return ret
}

//OptionSymbolWithIndex get all option symbol with specific index from symbol map
func OptionSymbolsWithIndex(index string) []exchange.OptionSymbol {
	ret := []exchange.OptionSymbol{}
	symbolMu.Lock()
	m := symbolMap
	symbolMu.Unlock()
	for _, v := range m {
		sym, ok := v.(exchange.OptionSymbol)
		if !ok {
			continue
		}

		if sym.Index() != index {
			continue
		}
		ret = append(ret, sym)
	}
	return ret
}

//NewOptionSymbol create a option symbol string with curreny, st, strike, typ. and parse it with ParseOptionSymbol
func NewOptionSymbol(currency string, st time.Time, strike float64, typ exchange.OptionType) (exchange.OptionSymbol, error) {
	var suffix string
	if typ == exchange.OptionTypeCall {
		suffix = "C"
	} else if typ == exchange.OptionTypePut {
		suffix = "P"
	} else {
		return nil, errors.Errorf("unkown option typ='%s'", typ)
	}
	symbol := fmt.Sprintf("%s-%s-%d-%s", currency, strings.ToUpper(st.Format(timeLayout)), int(strike), suffix)

	sym, err := ParseSymbol(symbol)
	if err != nil {
		return nil, errors.WithMessagef(err, "parse symbol fail symbol='%s'", symbol)
	}
	return sym.(exchange.OptionSymbol), nil
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

func updateSymbolMap(ctx context.Context, restClient *RestClient) error {
	newMap := map[string]exchange.Symbol{}
	for _, c := range Currencies {
		symbols, err := restClient.Symbols(ctx, c)
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
