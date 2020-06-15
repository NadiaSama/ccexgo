package deribit

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
)

//deribit message json serizlize

const (
	JsonRPCVersion     = "2.0"
	subscriptionMethod = "subscription"
)

type (
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

	Codec struct {
	}

	notifyParseCB func(*Notify) (*rpc.Notify, error)
)

var (
	notifyParseMap map[string]notifyParseCB = make(map[string]notifyParseCB)
)

func reigisterCB(key string, cb notifyParseCB) {
	_, ok := notifyParseMap[key]
	if ok {
		panic(fmt.Sprintf("duplicate cb %s register", key))
	}
	notifyParseMap[key] = cb
}

func (cc *Codec) Decode(raw []byte) (rpc.Response, error) {
	var resp Response
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}

	if resp.Method == subscriptionMethod {
		return parseNotify(&resp)
	}

	return &rpc.Result{
		ID: rpc.ID{Num: resp.ID},
		Error: rpc.Error{
			Code:    resp.Error.Code,
			Message: resp.Error.Message,
		},
		Result: resp.Result,
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

func parseNotify(resp *Response) (rpc.Response, error) {
	fields := strings.Split(resp.Params.Channel, ".")
	if len(fields) == 0 {
		return nil, errors.Errorf("bad message %v", resp.Params)
	}

	cb, ok := notifyParseMap[fields[0]]
	if !ok {
		return nil, errors.Errorf("unsupport channel %s", resp.Params.Channel)
	}
	return cb(&resp.Params)
}
