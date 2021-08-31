package ftx

import (
	"encoding/json"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
)

type (
	OrderBook struct {
		*exchange.OrderBookDS
		symbol exchange.Symbol
	}

	OrderBookData struct {
		Action    string       `json:"action"`
		Bids      [][2]float64 `json:"bids"`
		Asks      [][2]float64 `json:"asks"`
		Timestamp int64        `json:"timestamp"`
	}

	OrderBookChannel struct {
		symbol exchange.Symbol
	}
)

func NewOrderBookChannel(sym exchange.Symbol) exchange.Channel {
	return &OrderBookChannel{
		symbol: sym,
	}
}

func (obc *OrderBookChannel) String() string {
	return obc.symbol.String()
}

func NewOrderBook(sym exchange.Symbol) *OrderBook {
	return &OrderBook{symbol: sym}
}

func (ob *OrderBook) Init(cr *callResponse) (*exchange.OrderBook, error) {
	var obd OrderBookData
	if err := json.Unmarshal(cr.Data, &obd); err != nil {
		return nil, err
	}

	if obd.Action != typePartial {
		return nil, errors.Errorf("bad action '%s' expect '%s'", obd.Action, typePartial)
	}

	notify := obd.Transfer(ob.symbol)
	ob.OrderBookDS = exchange.NewOrderBookDS(notify)
	return ob.Snapshot(), nil
}

func (ob *OrderBook) Update(cr *callResponse) (*exchange.OrderBook, error) {
	var obd OrderBookData
	if err := json.Unmarshal(cr.Data, &obd); err != nil {
		return nil, err
	}
	if obd.Action != typeUpdate {
		return nil, errors.Errorf("bad action '%s' expect '%s'", obd.Action, typeUpdate)
	}

	notify := obd.Transfer(ob.symbol)

	ob.OrderBookDS.Update(notify)
	return ob.Snapshot(), nil
}

func (obd *OrderBookData) Transfer(sym exchange.Symbol) *exchange.OrderBookNotify {
	bids := make([]exchange.OrderElem, len(obd.Bids))
	for i, v := range obd.Bids {
		bids[i].Price = v[0]
		bids[i].Amount = v[1]
	}

	asks := make([]exchange.OrderElem, len(obd.Asks))
	for i, v := range obd.Asks {
		asks[i].Price = v[0]
		asks[i].Amount = v[1]
	}

	return &exchange.OrderBookNotify{
		Symbol: sym,
		Bids:   bids,
		Asks:   asks,
	}
}
