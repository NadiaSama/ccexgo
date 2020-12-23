package deribit

import (
	"encoding/json"
	"fmt"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/NadiaSama/ccexgo/misc/tconv"
	"github.com/pkg/errors"
)

type (
	IndexResult struct {
		IndexName string  `json:"index_name"`
		Price     float64 `json:"price"`
		Timestamp int64   `json:"timestamp"`
	}

	ChIndex struct {
		sym exchange.Symbol
	}
)

func init() {
	reigisterCB("deribit_price_index", parseNotifyIndex)
}

func NewIndexChannel(sym exchange.Symbol) exchange.Channel {
	return &ChIndex{
		sym: sym,
	}
}

func (ci *ChIndex) String() string {
	return fmt.Sprintf("deribit_price_index.%s", ci.sym.String())
}

func parseNotifyIndex(resp *Notify) (*rpc.Notify, error) {
	var ir IndexResult
	if err := json.Unmarshal(resp.Data, &ir); err != nil {
		return nil, errors.WithMessagef(err, "unmarshal index result")
	}

	sym, err := parseSpotSymbol(ir.IndexName)
	if err != nil {
		return nil, errors.WithMessagef(err, "bad indexName %s", ir.IndexName)
	}

	param := &exchange.IndexNotify{
		Price:   ir.Price,
		Created: tconv.Milli2Time(ir.Timestamp),
		Symbol:  sym,
	}
	return &rpc.Notify{
		Method: subscriptionMethod,
		Params: param,
	}, nil
}
