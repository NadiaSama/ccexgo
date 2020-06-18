package jex

import (
	"fmt"
	"sync"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
)

type (
	JexOptionSymbol struct {
		*exchange.BaseOptionSymbol
	}
)

var (
	symbolMap = map[string]*JexOptionSymbol{}
	mux       = sync.Mutex{}
)

func SetSymbol(sym *JexOptionSymbol) {
	mux.Lock()
	defer mux.Unlock()
	symbolMap[sym.String()] = sym
}

func NewOptionSymbol(index string, strike float64, settle time.Time, typ exchange.OptionType) *JexOptionSymbol {
	return &JexOptionSymbol{
		exchange.NewBaseOptionSymbol(strike, index, settle, typ),
	}
}

//ParseSymbol jex option symbol does not carry strike price
//the symbol info need get via crawler and set by SetSymbol
func ParseSymbol(sym string) (exchange.Symbol, error) {
	mux.Lock()
	defer mux.Unlock()
	ret, ok := symbolMap[sym]
	if !ok {
		return nil, errors.Errorf("unsupport symbol %s", sym)
	}
	return ret, nil
}

func (jos *JexOptionSymbol) String() string {
	var typ string
	if jos.Type() == exchange.OptionTypeCall {
		typ = "CALL"
	} else {
		typ = "PUT"
	}
	return fmt.Sprintf("%s%s%s", jos.Index(), jos.SettleTime().Format("0102"), typ)
}
