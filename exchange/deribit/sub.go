package deribit

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const (
	methodSubscribe   = "subscribe"
	methodUnSubscribe = "unsubscribe"
)

func (c *Client) Subscribe(ctx context.Context, channels ...string) error {
	return c.subInternal(ctx, methodSubscribe, channels...)
}

func (c *Client) UnSubscribe(ctx context.Context, channels ...string) error {
	return c.subInternal(ctx, methodUnSubscribe, channels...)
}

func (c *Client) subInternal(ctx context.Context, op string, channels ...string) error {
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
