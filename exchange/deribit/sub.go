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

func (c *Client) Subscribe(ctx context.Context, chs ...exchange.Channel) error {
	return c.subInternal(ctx, methodSubscribe, chs...)
}

func (c *Client) UnSubscribe(ctx context.Context, chs ...exchange.Channel) error {
	return c.subInternal(ctx, methodUnSubscribe, chs...)
}

func (c *Client) subInternal(ctx context.Context, op string, chs ...exchange.Channel) error {
	channels := []string{}
	for _, c := range chs {
		channels = append(channels, c.String())
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
