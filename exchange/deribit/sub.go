package deribit

import (
	"context"
	"fmt"
	"strings"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
)

const (
	methodSubscribe   = "subscribe"
	methodUnSubscribe = "unsubscribe"
)

type (
	subTypeCB func(syms ...exchange.Symbol) ([]string, error)
)

var (
	subType2CB map[exchange.SubType]subTypeCB = make(map[exchange.SubType]subTypeCB)
)

func (c *Client) Subscribe(ctx context.Context, typ exchange.SubType, syms ...exchange.Symbol) error {
	return c.subInternal(ctx, typ, methodSubscribe, syms...)
}

func (c *Client) UnSubscribe(ctx context.Context, typ exchange.SubType, syms ...exchange.Symbol) error {
	return c.subInternal(ctx, typ, methodUnSubscribe, syms...)
}

func (c *Client) subInternal(ctx context.Context, typ exchange.SubType, op string, syms ...exchange.Symbol) error {
	cb, ok := subType2CB[typ]
	if !ok {
		return exchange.NewBadArg("unsupport type", typ)
	}
	channels, err := cb(syms...)
	if err != nil {
		return err
	}

	var result []string
	method := fmt.Sprintf("public/%s", op)
	if err := c.call(ctx, method, map[string]interface{}{
		"channels": channels,
	}, &result, false); err != nil {
		return err
	}

	if len(result) != len(channels) {
		return errors.Errorf("%s [%s] error bad result [%s]",
			op, strings.Join(channels, ","), strings.Join(result, ","))
	}
	set := map[string]struct{}{}
	for _, r := range result {
		set[r] = struct{}{}
	}
	for _, r := range channels {
		if _, ok := set[r]; !ok {
			return errors.Errorf("failed %s channel %s", op, r)
		}
	}
	return nil
}

func registerSubTypeCB(typ exchange.SubType, cb subTypeCB) {
	subType2CB[typ] = cb
}
