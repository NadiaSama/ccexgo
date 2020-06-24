package deribit

import (
	"context"
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
}

func (c *Client) SubscribeIndex(ctx context.Context, sym exchange.Symbol) error {
	return c.subInternal(ctx, methodSubscribe, indexChannel(sym))
}

func (c *Client) UnSubscribeIndex(ctx context.Context, sym exchange.Symbol) error {
	return c.subInternal(ctx, methodUnSubscribe, indexChannel(sym))
}

func indexChannel(sym exchange.Symbol) string {
	return fmt.Sprintf("deribit_price_index.%s", sym.String())
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
