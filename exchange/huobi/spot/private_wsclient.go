package spot

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/huobi"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
)

type (
	PrivateWSClient struct {
		key    string
		secret string
		*exchange.WSClient
		data chan interface{}
	}
)

const (
	PrivateWSClientAddr = "wss://api.huobi.pro/ws/v2"
)

func NewPrivateWSClient(key, secret string, data chan interface{}) *PrivateWSClient {
	ret := &PrivateWSClient{
		key:    key,
		secret: secret,
		data:   data,
	}

	ret.WSClient = exchange.NewWSClient(PrivateWSClientAddr, NewPrivateCodeC(), ret)
	return ret
}

func (pws *PrivateWSClient) Run(ctx context.Context) error {
	if err := pws.WSClient.Run(ctx); err != nil {
		return err
	}

	return pws.Auth(ctx)
}

func (pws *PrivateWSClient) Subscribe(ctx context.Context, channels ...exchange.Channel) error {
	if len(channels) != 1 {
		return errors.Errorf("invalid channeld num=%d", len(channels))
	}

	params := PrivateWSReq{
		Action: ActionSub,
		Ch:     channels[0].String(),
	}

	var resp PrivateWSResp
	if err := pws.Call(ctx, params.Ch, params.Action, &params, &resp); err != nil {
		return errors.WithMessage(err, "subscribe fail")
	}

	return nil
}

func (pws *PrivateWSClient) Handle(ctx context.Context, n *rpc.Notify) {
	if n.Method == ActionPing {
		go func() {
			pws.Call(ctx, "", "", map[string]interface{}{
				"action": ActionPong,
				"data":   n.Params,
			}, nil)
		}()
		return
	}

	en := exchange.WSNotify{
		Exchange: huobi.Huobi,
		Chan:     n.Method,
		Data:     n.Params,
	}
	select {
	case pws.data <- &en:
	default:
	}
}

func (pws *PrivateWSClient) genSignatureParmas() map[string]string {
	ts := time.Now().UTC()
	ret := map[string]string{
		"accessKey":        pws.key,
		"signatureMethod":  "HmacSHA256",
		"signatureVersion": "2.1",
		"timestamp":        ts.Format("2006-01-02T15:04:05"),
	}
	values := url.Values{}
	for k, v := range ret {
		values.Add(k, v)
	}
	sig := huobi.Signature(pws.secret, http.MethodGet, "api.huobi.pro", "/ws/v2", values.Encode())

	ret["signature"] = sig
	ret["authType"] = "api"
	return ret
}

func (pws *PrivateWSClient) Auth(ctx context.Context) error {
	param := pws.genSignatureParmas()
	req := PrivateWSReq{
		Action: ActionReq,
		Ch:     "auth",
		Params: param,
	}

	var resp PrivateWSResp
	if err := pws.Call(ctx, req.Ch, req.Action, &req, &resp); err != nil {
		return errors.WithMessage(err, "auth failed")
	}

	return nil
}
