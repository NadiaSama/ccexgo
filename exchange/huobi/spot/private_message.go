package spot

import (
	"encoding/json"

	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
)

type (
	PrivateCodeC struct {
	}

	PrivateWSReq struct {
		Action string      `json:"action"`
		Ch     string      `json:"ch"`
		Params interface{} `json:"params,omitempty"`
	}
	PrivateWSResp struct {
		Action string          `json:"action"`
		Code   int             `json:"code"`
		Ch     string          `json:"ch"`
		Data   json.RawMessage `json:"data"`
	}
)

const (
	ActionPing = "ping"
	ActionPong = "pong"
	ActionReq  = "req"
	ActionSub  = "sub"
	ActionPush = "push"
)

func NewPrivateCodeC() *PrivateCodeC {
	return &PrivateCodeC{}
}

func (pcc *PrivateCodeC) Encode(req rpc.Request) ([]byte, error) {
	raw, err := json.Marshal(req.Params())
	return raw, err
}

func (pcc *PrivateCodeC) Decode(raw []byte) (rpc.Response, error) {
	var resp PrivateWSResp
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}

	if resp.Action == ActionPing {
		return &rpc.Notify{Method: ActionPing, Params: resp.Data}, nil
	}
	if resp.Action == ActionReq || resp.Action == ActionSub {
		var err error
		if resp.Code != 200 {
			err = errors.Errorf("invalid response %s", string(raw))
		}
		return &rpc.Result{
			ID:     resp.Ch,
			Error:  err,
			Result: raw,
		}, nil
	}

	if resp.Action == ActionPush {
		r, err := ParseOrder(resp.Data)
		if err != nil {
			return nil, err
		}

		return &rpc.Notify{
			Method: resp.Ch,
			Params: r,
		}, nil
	}
	return nil, nil
}
