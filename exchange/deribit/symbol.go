package deribit

import (
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
	timeLayout = "02Jan06"
)

//instrument_name-settle_time-strike-type
func PraseOptionSymbol(val string) (exchange.OptionSymbol, error) {
	fields := strings.Split(val, "-")
	if len(fields) != 4 {
		return nil, exchange.ErrBadSymbol
	}
	var typ exchange.OptionType
	var strike float64
	var st time.Time

	if fields[0] != "BTC" && fields[0] != "ETH" {
		return nil, exchange.ErrBadSymbol
	}
	if fields[3] == "C" {
		typ = exchange.OptionTypeCall
	} else if fields[3] == "P" {
		typ = exchange.OptionTypePut
	} else {
		return nil, exchange.ErrBadSymbol
	}
	strike, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, exchange.ErrBadSymbol
	}
	st, err = time.Parse(timeLayout, fields[1])
	if err != nil {
		return nil, exchange.ErrBadSymbol
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
	strike := int(sym.Strike())
	return fmt.Sprintf("%s-%s-%d-%s", sym.Index(), st, strike, typ)
}
