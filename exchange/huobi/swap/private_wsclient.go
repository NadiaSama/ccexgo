package swap

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/huobi"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
)

type (
	PrivateCodeC struct {
		*huobi.CodeC
	}

	PrivateWSClient struct {
		key    string
		secret string
		*exchange.WSClient
		data chan interface{}
	}

	Response struct {
		Op string `json:"op"`
	}

	subParam struct {
		Op    string `json:"op"`
		Cid   string `json:"cid"`
		Topic string `json:"topic"`
	}

	PingResponse struct {
		Op string `json:"op"`
		TS string `json:"ts"`
	}

	AuthResponse struct {
		Op      string `json:"op"`
		ErrCode int    `json:"err-code"`
		TS      int64  `json:"ts"`
		ErrMsg  string `json:"err-msg"`
	}
)

const (
	SwapPrivateAddr = "wss://api.hbdm.com/swap-notification"
)

func NewPrivateCodeC() *PrivateCodeC {
	return &PrivateCodeC{
		CodeC: huobi.NewCodeC(),
	}
}

func (pcc *PrivateCodeC) Encode(req rpc.Request) ([]byte, error) {
	param := req.Params()
	r, e := json.Marshal(param)
	return r, e
}

func (pcc *PrivateCodeC) Decode(raw []byte) (rpc.Response, error) {
	msg, err := pcc.Decompress(raw)
	if err != nil {
		return nil, err
	}

	var resp Response
	if err := json.Unmarshal(msg, &resp); err != nil {
		return nil, err
	}

	if resp.Op == "ping" {
		var pr PingResponse
		if err := json.Unmarshal(msg, &pr); err != nil {
			return nil, err
		}
		return &rpc.Notify{
			Method: "pong",
			Params: pr.TS,
		}, nil
	}

	if resp.Op == "auth" {
		var ar AuthResponse
		if err := json.Unmarshal(msg, &ar); err != nil {
			return nil, err
		}
		var err error
		if ar.ErrCode != 0 {
			err = errors.Errorf("error happend %s", string(msg))
		}

		return &rpc.Result{
			ID:     "auth",
			Error:  err,
			Result: msg,
		}, nil
	}

	if resp.Op == "notify" {
		r, err := ParseOrder(msg)
		if err != nil {
			return nil, err
		}

		op := r.Raw.(*OrderNotify)
		return &rpc.Notify{
			Method: op.Topic,
			Params: r,
		}, nil
	}

	return nil, errors.Errorf("unkown op %s", resp.Op)
}

func NewPrivateWSClient(key, secret string, data chan interface{}) *PrivateWSClient {
	ret := &PrivateWSClient{
		key:    key,
		secret: secret,
		data:   data,
	}

	ret.WSClient = exchange.NewWSClient(SwapPrivateAddr, NewPrivateCodeC(), ret)
	return ret
}

func (ws *PrivateWSClient) Run(ctx context.Context) error {
	if err := ws.WSClient.Run(ctx); err != nil {
		return err
	}
	return ws.Auth(ctx)
}

func (ws *PrivateWSClient) Auth(ctx context.Context) error {
	param := ws.genSignatureParmas()
	var resp Response
	if err := ws.Call(ctx, "auth", "", param, &resp); err != nil {
		return err
	}

	return nil
}

func (ws *PrivateWSClient) Subscribe(ctx context.Context, channels ...exchange.Channel) error {
	if len(channels) != 1 {
		return errors.Errorf("only 1 channel can be subcribed")
	}

	param := subParam{
		Op:    "sub",
		Cid:   "123",
		Topic: channels[0].String(),
	}

	err := ws.Call(ctx, "", "", param, nil)
	return errors.WithMessage(err, "subscribe fail")
}

func (ws *PrivateWSClient) Handle(ctx context.Context, notify *rpc.Notify) {
	if notify.Method == huobi.MethodPong {
		go func() {
			ws.Call(ctx, "", "", map[string]interface{}{
				"op": "pong",
				"ts": notify.Params,
			}, nil)
		}()
	}

	d := exchange.WSNotify{
		Exchange: huobi.Huobi,
		Chan:     notify.Method,
		Data:     notify.Params,
	}
	select {
	case ws.data <- &d:
	default:
	}
}

func (pws *PrivateWSClient) genSignatureParmas() map[string]string {
	ts := time.Now().UTC()
	ret := map[string]string{
		"AccessKeyId":      pws.key,
		"SignatureMethod":  "HmacSHA256",
		"SignatureVersion": "2",
		"Timestamp":        ts.Format("2006-01-02T15:04:05"),
	}
	values := url.Values{}
	for k, v := range ret {
		values.Add(k, v)
	}
	sig := huobi.Signature(pws.secret, http.MethodGet, "api.hbdm.com", "/swap-notification", values.Encode())

	ret["Signature"] = sig
	ret["type"] = "api"
	ret["op"] = "auth"
	return ret
}
