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

func NewClient(ctx context.Context, conn rpc.Conn, key, secret string) *Client {
	c := exchange.Client{
		Conn:    conn,
		Key:     key,
		Secret:  secret,
		Ctx:     ctx,
		Timeout: time.Second * 2,
	}
	return &Client{
		Client: &c,
	}
}

func (c *Client) call(method string, params interface{}, dest interface{}, private bool) error {
	if private {
		ac, err := c.getToken()
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
	ctx, cancel := context.WithTimeout(c.Ctx, c.Timeout)
	defer cancel()
	err := c.Conn.Call(ctx, method, params, dest)
	return err
}
