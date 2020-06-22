package deribit

import (
	"context"
	"strings"

	"github.com/pkg/errors"
)

func (c *Client) Subscribe(ctx context.Context, channels ...string) error {
	var result []string
	if err := c.call(ctx, "public/subscribe", map[string]interface{}{
		"channels": channels,
	}, &result, false); err != nil {
		return err
	}

	if len(result) != len(channels) {
		return errors.Errorf("subscribe [%s] error bad result [%s]",
			strings.Join(channels, ","), strings.Join(result, ","))
	}
	set := map[string]struct{}{}
	for _, r := range result {
		set[r] = struct{}{}
	}
	for _, r := range channels {
		if _, ok := set[r]; !ok {
			return errors.Errorf("failed subscribe channel %s", r)
		}
	}
	return nil
}
