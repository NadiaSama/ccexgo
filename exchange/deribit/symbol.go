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

	SpotSymbol struct {
		*exchange.BaseSpotSymbol
	}
)

var (
	timeLayout = "02Jan06"
)

//instrument_name-settle_time-strike-type
func ParseOptionSymbol(val string) (exchange.OptionSymbol, error) {
	fields := strings.Split(val, "-")
	failed := true
	msg := "bad symbol"
	var (
		arg    interface{} = nil
		typ    exchange.OptionType
		strike float64
		st     time.Time
		err    error
	)

	for {
		if len(fields) != 4 {
			break
		}
		if fields[0] != "BTC" && fields[0] != "ETH" {
			break
		}
		if fields[3] == "C" {
			typ = exchange.OptionTypeCall
		} else if fields[3] == "P" {
			typ = exchange.OptionTypePut
		} else {
			break
		}
		strike, err = strconv.ParseFloat(fields[2], 64)
		if err != nil {
			msg = "parse float error"
			arg = err
		}
		st, err = time.Parse(timeLayout, fields[1])
		if err != nil {
			msg = "parse time error"
			arg = err
		}
		failed = false
		break
	}

	if failed {
		return nil, exchange.NewBadArg(msg, arg)
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

func NewSpotSymbol(base, quote string) exchange.SpotSymbol {
	return &SpotSymbol{
		exchange.NewBaseSpotSymbol(strings.ToLower(base), strings.ToLower(quote)),
	}
}

func ParseSpotSymbol(sym string) (exchange.SpotSymbol, error) {
	fields := strings.Split(strings.ToLower(sym), "_")
	if len(fields) != 2 {
		return nil, exchange.NewBadArg("bad spot symbol field len", len(fields))
	}
	return &SpotSymbol{
		exchange.NewBaseSpotSymbol(fields[0], fields[1]),
	}, nil
}

func (sym *SpotSymbol) String() string {
	return fmt.Sprintf("%s_%s", sym.Base(), sym.Quote())
}
