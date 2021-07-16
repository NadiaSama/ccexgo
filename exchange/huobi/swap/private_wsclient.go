package swap

import (
	"context"
	"encoding/json"
	"fmt"
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
		Op             string      `json:"op"`
		TS             interface{} `json:"ts"` //ping message and auth message ts field have different type
		CID            string      `json:"cid"`
		ErrCode        int         `json:"err-code"`
		Topic          string      `json:"topci"`
		Symbol         string      `json:"symbol"`
		ContractCode   string      `json:"contract_code"`
		Volume         float64     `json:"volumen"`
		Price          float64     `json:"price"`
		OrderPriceType string      `json:"order_price_type"`
		Direction      string      `json:"direction"`
		Offset         string      `json:"offset"`
		Status         int         `json:"status"`
		LeverRate      int         `json:"lever_rate"`
		OrderID        int64       `json:"order_id"`
		OrderIDStr     string      `json:"order_id_str"`
		OrderType      int         `json:"order_type"`
		CreatedAt      int64       `json:"created_at"`
		TradeVolume    int         `json:"trade_volume"`
		TradeTurnOver  int         `json:"trade_turnover"`
		Fee            float64     `json:"fee"`
		TradeAvgPrice  float64     `json:"trade_avg_price"`
		CanceledAt     int64       `json:"canceled_at"`
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
	fmt.Printf("%s\n", string(r))
	return r, e
}

func (pcc *PrivateCodeC) Decode(raw []byte) (rpc.Response, error) {
	msg, err := pcc.Decompress(raw)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s\n", string(msg))

	var resp Response
	if err := json.Unmarshal(msg, &resp); err != nil {
		return nil, err
	}

	if resp.Op == "ping" {
		return &rpc.Notify{
			Method: "pong",
			Params: resp.TS,
		}, nil
	}

	if resp.Op == "auth" {
		var err error
		if resp.ErrCode != 0 {
			err = errors.Errorf("error happend %s", string(msg))
		}

		return &rpc.Result{
			ID:     "auth",
			Error:  err,
			Result: msg,
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
