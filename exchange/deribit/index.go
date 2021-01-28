package deribit

import (
	"encoding/json"
	"fmt"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/NadiaSama/ccexgo/misc/tconv"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	IndexResult struct {
		IndexName string          `json:"index_name"`
		Price     decimal.Decimal `json:"price"`
		Timestamp int64           `json:"timestamp"`
	}

	ChIndex struct {
		index string
	}
)

func init() {
	reigisterCB("deribit_price_index", parseNotifyIndex)
}

func NewIndexChannel(index string) exchange.Channel {
	return &ChIndex{
		index: index,
	}
}

func (ci *ChIndex) String() string {
	return fmt.Sprintf("deribit_price_index.%s", ci.index)
}

func parseNotifyIndex(resp *Notify) (*rpc.Notify, error) {
	var ir IndexResult
	if err := json.Unmarshal(resp.Data, &ir); err != nil {
		return nil, errors.WithMessagef(err, "unmarshal index result")
	}

	sym, err := ParseIndexSymbol(ir.IndexName)
	if err != nil {
		return nil, err
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
