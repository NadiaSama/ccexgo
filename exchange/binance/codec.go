package binance

import (
	"encoding/json"
	"strconv"

	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

type (
	// CodeC common codec for spot swap notify
	CodeC struct {
		lastID string
	}

	SubscribeRequest struct {
		Method string      `json:"method"`
		Params interface{} `json:"params"`
		ID     int64       `json:"id"`
	}

	CallResult struct {
		Result interface{} `json:"result"`
		ID     int64       `json:"id"`
	}
)

const (
	MethodSubscibe    = "SUBSCRIBE"
	MethodUnSubscribe = "UNSUBSCRIBE"
)

func NewCodeC() *CodeC {
	return &CodeC{}
}

func (cc *CodeC) LastID() string {
	return cc.lastID
}

// Encode req to binance subscribe message
func (cc *CodeC) Encode(req rpc.Request) ([]byte, error) {
	id, err := strconv.ParseInt(req.ID(), 10, 64)
	if err != nil {
		return nil, errors.WithMessage(err, "invalid id")
	}

	sub := SubscribeRequest{
		ID:     id,
		Params: req.Params(),
		Method: req.Method(),
	}

	cc.lastID = req.ID()

	return json.Marshal(&sub)
}

func (cc *CodeC) DecodeByCB(raw []byte, cb func(g *gjson.Result) (rpc.Response, error)) (rpc.Response, error) {
	g := gjson.ParseBytes(raw)

	// by now only handle subscribe response which result is nil
	if g.Get("id").Exists() && g.Get("result").Exists() {
		return &rpc.Result{
			ID: g.Get("id").String(),
		}, nil
	}

	if g.Get("code").Exists() {
		return &rpc.Result{
			ID:    cc.LastID(),
			Error: errors.Errorf("error code: %d msg: %s", g.Get("code").Int(), g.Get("msg").String()),
		}, nil
	}

	return cb(&g)
}
