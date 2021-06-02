package okex5

import (
	"encoding/json"

	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
)

type (
	CodeC struct {
		LastEvent string
	}

	wsReq struct {
		Op   string      `json:"op"`
		Args interface{} `json:"args"`
	}

	wsRespArg struct {
		Channel string `json:"channel"`
		InstId  string `json:"instId"`
		Uly     string `json:"uly"`
	}

	wsResp struct {
		Event  string          `json:"event"`
		Code   string          `json:"code"`
		Msg    string          `json:"msg"`
		Arg    wsRespArg       `json:"arg"`
		Action string          `json:"action"`
		Data   json.RawMessage `json:"data"`
	}

	parseCB func(*wsResp) (*rpc.Notify, error)
)

const (
	eventError = "error"
)

var (
	parseCBMap map[string]parseCB = make(map[string]parseCB)
)

func NewCodec() *CodeC {
	return &CodeC{}
}

func (cc *CodeC) Encode(req rpc.Request) ([]byte, error) {
	method := req.Method()
	data := req.Params()

	r := wsReq{
		Op:   method,
		Args: data,
	}
	cc.LastEvent = method

	return json.Marshal(&r)
}

func (cc *CodeC) Decode(raw []byte) (rpc.Response, error) {
	var resp wsResp
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}

	if resp.Event != "" {
		if resp.Event == eventError {
			return &rpc.Result{
				ID:    cc.LastEvent,
				Error: errors.Errorf("rpc error code: %s msg: %s", resp.Code, resp.Msg),
			}, nil
		}
		return &rpc.Result{
			ID:     cc.LastEvent,
			Result: raw,
		}, nil
	}

	//websocket notify
	return resp.transfer()
}

func (r *wsResp) transfer() (*rpc.Notify, error) {
	cb, ok := parseCBMap[r.Arg.Channel]

	if !ok {
		return nil, errors.Errorf("unknown channel %s", r.Arg.Channel)
	}
	return cb(r)
}
