package deribit

import (
	"encoding/json"

	"github.com/NadiaSama/ccexgo/internal/rpc"
)

//deribit message json serizlize

const (
	JsonRPCVersion = "2.0"
)

type (
	Codec struct {
	}

	Notify struct {
		Data    json.RawMessage `json:"data"`
		Channel string          `json:"channel"`
	}

	Error struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}
	Response struct {
		JsonRPC string          `json:"jsonrpc"`
		ID      int64           `json:"id"`
		Error   Error           `json:"error"`
		Result  json.RawMessage `json:"result"`
		Method  string          `json:"method"`
		Params  Notify          `json:"params"`
	}

	Request struct {
		ID      int64       `json:"id"`
		Method  string      `json:"method"`
		JsonRPC string      `json:"jsonrpc"`
		Params  interface{} `json:"params"`
	}
)

func (cc *Codec) Decode(raw []byte) (rpc.Response, error) {
	var resp Response
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}

	if resp.Method == "subscription" {
		return &rpc.Notify{
			Method: resp.Method,
			Params: resp.Params,
		}, nil
	}

	return &rpc.Result{
		ID: rpc.ID{Num: resp.ID},
		Error: rpc.Error{
			Code:    resp.Error.Code,
			Message: resp.Error.Message,
		},
	}, nil
}

func (cc *Codec) Encode(req rpc.Request) ([]byte, error) {
	r := Request{
		ID:      req.ID().Num,
		Method:  req.Method(),
		Params:  req.Params(),
		JsonRPC: JsonRPCVersion,
	}

	return json.Marshal(&r)
}
