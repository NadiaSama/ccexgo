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
)

func init() {
	reigisterCB("deribit_price_index", parseNotifyIndex)
	registerSubTypeCB(exchange.SubTypeIndex, indexChannel)
}

func indexChannel(syms ...exchange.Symbol) ([]string, error) {
	ret := make([]string, len(syms))
	for i, sym := range syms {
		ret[i] = fmt.Sprintf("deribit_price_index.%s", sym.String())
	}
	return ret, nil
}

func parseNotifyIndex(resp *Notify) (*rpc.Notify, error) {
	var ir IndexResult
	if err := json.Unmarshal(resp.Data, &ir); err != nil {
		return nil, errors.WithMessagef(err, "bad index %s", string(resp.Data))
	}

	sym, err := ParseSpotSymbol(ir.IndexName)
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
