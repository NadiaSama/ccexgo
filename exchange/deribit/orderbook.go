package deribit

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
)

type (
	BookData struct {
		Timestamp      int              `json:"timestamp"`
		InstrumentName string           `json:"instrument_name"`
		ChangeID       int              `json:"charge_id"`
		Bids           [][3]interface{} `json:"bids"`
		Asks           [][3]interface{} `json:"asks"`
	}
)

func init() {
	reigisterCB("book", parseNotifyBook)
	registerSubTypeCB(exchange.SubTypeOrderBook, orderBookChannel)
}

func orderBookChannel(syms ...exchange.Symbol) ([]string, error) {
	ret := make([]string, len(syms))
	for i, sym := range syms {
		ret[i] = fmt.Sprintf("book.%s.raw", sym.String())
	}
	return ret, nil
}
func parseNotifyBook(resp *Notify) (*rpc.Notify, error) {
	fields := strings.Split(resp.Channel, ".")
	var bn BookData
	if err := json.Unmarshal(resp.Data, &bn); err != nil {
		return nil, errors.WithMessage(err, "unarshal orderbookNotify")
	}
	sym, err := parseOptionSymbol(fields[1])
	if err != nil {
		return nil, errors.WithMessagef(err, "parse orderbookNotify symbol %s", fields[1])
	}
	notify := &rpc.Notify{
		Method: subscriptionMethod,
	}
	on := &exchange.OrderBookNotify{
		Symbol: sym,
		Asks:   make([]exchange.OrderElem, len(bn.Asks)),
		Bids:   make([]exchange.OrderElem, len(bn.Bids)),
	}

	if err := processArr(on.Asks, bn.Asks); err != nil {
		return nil, err
	}
	if err := processArr(on.Bids, bn.Bids); err != nil {
		return nil, err
	}
	notify.Params = on
	return notify, nil
}

func processArr(d []exchange.OrderElem, s [][3]interface{}) (ret error) {
	defer func() {
		if err := recover(); err != nil {
			ret = err.(error)
		}
	}()

	for i, v := range s {
		op := v[0].(string)
		price := v[1].(float64)
		amount := v[2].(float64)

		if op == "new" || op == "change" {
			d[i].Amount = amount
			d[i].Price = price
		} else if op == "delete" {
			d[i].Amount = 0
			d[i].Price = price
		} else {
			ret = errors.Errorf("unkown op %s", op)
			return
		}
	}
	return
}
