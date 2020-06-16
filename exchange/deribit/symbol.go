package deribit

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
)

type (
	OptionSymbol struct {
		*exchange.BaseOptionSymbol
	}
)

var (
	ErrBadSymbol = errors.New("bad symbol")
	timeLayout   = "02Jan06"
)

//instrument_name-settle_time-strike-type
func PraseOptionSymbol(val string) (exchange.OptionSymbol, error) {
	fields := strings.Split(val, "-")
	if len(fields) != 4 {
		return nil, ErrBadSymbol
	}
	var typ exchange.OptionType
	var strike int
	var st time.Time

	if fields[0] != "BTC" && fields[0] != "ETH" {
		return nil, ErrBadSymbol
	}
	if fields[3] == "C" {
		typ = exchange.OptionTypeCall
	} else if fields[3] == "P" {
		typ = exchange.OptionTypePut
	} else {
		return nil, ErrBadSymbol
	}
	strike, err := strconv.Atoi(fields[2])
	if err != nil {
		return nil, ErrBadSymbol
	}
	st, err = time.Parse(timeLayout, fields[1])
	if err != nil {
		return nil, ErrBadSymbol
	}
	//deribit settle at utc 8:00
	st = st.UTC()
	st = st.Add(time.Hour * 8)
	osym := exchange.NewBaseOptionSymbol(strike, fields[0], st, typ)
	return &OptionSymbol{
		osym,
	}, nil
}

func (sym *OptionSymbol) String() string {
	typ := "P"
	if sym.Type() == exchange.OptionTypeCall {
		typ = "C"
	}
	st := strings.ToUpper(sym.SettleTime().Format(timeLayout))
	return fmt.Sprintf("%s-%s-%d-%s", sym.Index(), st, sym.Strike(), typ)
}
