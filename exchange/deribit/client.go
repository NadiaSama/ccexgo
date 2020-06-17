package deribit

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
)

type (
	Client struct {
		*exchange.Client
		tokenMu     sync.Mutex
		accessToken string
		expire      time.Time
	}
)

func NewClient(conn rpc.Conn, key, secret string) *Client {
	c := exchange.Client{
		Conn:   conn,
		Key:    key,
		Secret: secret,
	}
	return &Client{
		Client: &c,
	}
}

func (c *Client) call(ctx context.Context, method string, params interface{}, dest interface{}, private bool) error {
	if private {
		ac, err := c.getToken(ctx)
		if err != nil {
			return errors.WithMessage(err, "get access token fail")
		}

		switch token := params.(type) {
		case Token:
			token.SetToken(ac)

		case map[string]interface{}:
			token["access_token"] = ac

		default:
			return fmt.Errorf("method %s private no access_token specific", method)
		}

	}
	err := c.Conn.Call(ctx, method, params, dest)
	return err
}
