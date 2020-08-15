package ftx

import (
	"encoding/json"
	"fmt"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
)

type (
	CodeC struct {
		*exchange.CodeC
	}

	callParam struct {
		Channel string `json:"channel,omitempty"`
		Market  string `json:"market,omitempty"`
		OP      string `json:"op,omitempty"`
	}

	callResponse struct {
		Channel string          `json:"channel"`
		Market  string          `json:"market"`
		Type    string          `json:"type"`
		Code    int             `json:"code"`
		Msg     string          `json:"msg"`
		Data    json.RawMessage `json:"data"`
	}

	authArgs struct {
		Key  string `json:"key"`
		Sign string `json:"sign"`
		Time int64  `json:"time"`
	}

	authParam struct {
		Args authArgs `json:"args"`
		OP   string   `json:"op"`
	}
)

const (
	typeError        = "error"
	typeSubscribed   = "subscribed"
	typeUnSubscribed = "unsubscribed"
	typePong         = "pong"
	typeInfo         = "info"

	codeReconnet = 20001
)

func NewCodeC() *CodeC {
	return &CodeC{
		exchange.NewCodeC(),
	}
}

func (cc *CodeC) Decode(raw []byte) (rpc.Response, error) {
	var cr callResponse
	if err := json.Unmarshal(raw, &cr); err != nil {
		return nil, err
	}

	id := fmt.Sprintf("%s%s", cr.Channel, cr.Market)
	if cr.Type == typeError {
		ret := &rpc.Result{
			ID:     id,
			Error:  errors.Errorf("error msg: %s code: %d", cr.Msg, cr.Code),
			Result: raw,
		}
		return ret, nil
	}

	switch cr.Type {
	case typeSubscribed:
		fallthrough
	case typeUnSubscribed:
		ret := &rpc.Result{
			ID:     id,
			Result: raw,
		}
		return ret, nil

	case typePong:
		ret := &rpc.Notify{
			Method: typePong,
		}
		return ret, nil

	case typeInfo:
		if cr.Code == codeReconnet {
			return nil, rpc.NewStreamError(errors.Errorf("ftx ws reset info %s", string(raw)))
		}
		ret := &rpc.Notify{
			Method: id,
			Params: cr.Data,
		}
		return ret, nil

	default:
		ret := &rpc.Notify{
			Method: id,
			Params: cr.Data,
		}
		return ret, nil
	}
}
