package deribit

import (
	"context"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
)

type (
	Client struct {
		*exchange.WSClient
		tokenMu     sync.Mutex
		accessToken string
		expire      time.Time
		seq         int64
		key         string
		secret      string
		data        chan interface{}
	}
)

func NewWSClient(key, secret string, data chan interface{}) *Client {
	return newWSClient(WSAddr, key, secret, data)
}

func NewTestWSClient(key, secret string, data chan interface{}) *Client {
	return newWSClient(WSTestAddr, key, secret, data)
}

func newWSClient(addr, key, secret string, data chan interface{}) *Client {
	codec := &Codec{}
	ret := &Client{
		key:    key,
		secret: secret,
		data:   data,
	}
	ret.WSClient = exchange.NewWSClient(addr, codec, ret)
	return ret
}

func (c *Client) Exchange() string {
	return "deribit"
}

func (c *Client) Handle(ctx context.Context, notify *rpc.Notify) {
	data := &exchange.WSNotify{
		Exchange: c.Exchange(),
		Chan:     notify.Method,
		Data:     notify.Params,
	}
	select {
	case c.data <- data:
	default:
		return
	}
}

//Auth is done by client.call
func (c *Client) Auth(ctx context.Context) error {
	return nil
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
			return errors.Errorf("method %s private no access_token specific", method)
		}

	}
	id := atomic.AddInt64(&c.seq, 1)
	err := c.Conn.Call(ctx, strconv.FormatInt(id, 10), method, params, dest)
	return exchange.NewBadExResp(err)
}
