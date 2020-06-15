package deribit

import (
	"encoding/json"
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
}

func parseNotifyBook(resp *Notify) (*rpc.Notify, error) {
	fields := strings.Split(resp.Channel, ".")
	var bn BookData
	if err := json.Unmarshal(resp.Data, &bn); err != nil {
		return nil, err
	}
	notify := &rpc.Notify{
		Method: subscriptionMethod,
	}
	on := &exchange.OrderBookNotify{
		Symbol: fields[1],
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
