package deribit

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
)

type (
	Client struct {
		*exchange.Client
		tokenMu     sync.Mutex
		accessToken string
		expire      time.Time
	}
)

func (c *Client) call(method string, params interface{}, dest interface{}, private bool) error {
	if private {
		token, ok := params.(Token)
		if !ok {
			return fmt.Errorf("method %s private no access_token specific", method)
		}
		token.SetToken(c.accessToken)
	}
	ctx, _ := context.WithTimeout(c.Ctx, time.Second*5)
	err := c.Conn.Call(ctx, method, params, dest)
	return err
}
