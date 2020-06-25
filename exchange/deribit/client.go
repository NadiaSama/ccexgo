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

func NewClient(key, secret string, test bool) *Client {
	var addr string
	if test {
		addr = WSTestAddr
	} else {
		addr = WSAddr
	}

	ret := &Client{
		Client: exchange.NewClient(newDeribitConn, addr, key, secret),
	}
	return ret
}

func (c *Client) Exchange() string {
	return "deribit"
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
			panic(fmt.Sprintf("method %s private no access_token specific", method))
		}

	}
	err := c.Conn.Call(ctx, method, params, dest)
	return exchange.NewBadExResp(err)
}

func newDeribitConn(addr string) (rpc.Conn, error) {
	stream, err := rpc.NewWebsocketStream(addr, &Codec{})
	if err != nil {
		return nil, err
	}
	conn := rpc.NewConn(stream)
	return conn, nil
}
