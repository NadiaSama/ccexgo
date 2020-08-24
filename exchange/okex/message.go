package okex

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
)

type (
	CodeC struct {
		LastSUBID string //okex op fail do not return operate type. we have to record
	}

	callParam struct {
		OP   string   `json:"op"`
		Args []string `json:"args"`
	}

	response struct {
		Event     string          `json:"event"`
		Table     string          `json:"table"`
		Action    string          `json:"action"`
		Data      json.RawMessage `json:"data"`
		Message   string          `json:"message"`
		ErrorCode int             `json:"errorCode"`
	}

	pingReq struct {
	}

	ResponseParseCB func(string, string, json.RawMessage) (*rpc.Notify, error)
)

const (
	eventError = "error"
	idPingPong = "ping-pong"
	pingMsg    = "ping"
	pongMsg    = "pong"
	timeFMT    = "2006-01-02T15:04:05.000Z"
)

var (
	rcbMap map[string]ResponseParseCB = map[string]ResponseParseCB{}
)

func NewCodeC() *CodeC {
	return &CodeC{}
}

func (cc *CodeC) Encode(req rpc.Request) ([]byte, error) {
	param := req.Params()
	r, ok := param.(*callParam)
	if ok {
		cc.LastSUBID = r.OP
	}

	if _, ok := param.(*pingReq); ok {
		return []byte("ping"), nil
	}

	return json.Marshal(req)
}

func (cc *CodeC) Decode(raw []byte) (rpc.Response, error) {
	reader := flate.NewReader(bytes.NewReader(raw))
	resp, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, rpc.NewStreamError(errors.WithMessagef(err, "decompress error"))
	}

	if len(resp) == 4 && string(resp) == "pong" {
		return &rpc.Result{
			ID: idPingPong,
			//return a empty json for json.Unmarshal work
			Result: []byte(`{}`),
		}, nil
	}

	var r response
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, rpc.NewMsgError(raw, err)
	}

	if r.Event == eventError {
		return &rpc.Result{
			ID:    cc.LastSUBID,
			Error: errors.Errorf("error: %s code: %d", r.Message, r.ErrorCode),
		}, nil
	}

	if r.Event == opSubscribe || r.Event == opUnSubscribe {
		return &rpc.Result{
			ID:     r.Event,
			Result: resp,
		}, nil
	}

	if len(r.Table) != 0 {
		return &rpc.Notify{
			Method: r.Table,
			Params: &r,
		}, nil
	}
	return nil, errors.Errorf("unkown message '%s'", string(resp))
}

func (r *response) transfer() (*rpc.Notify, error) {
	cb, ok := rcbMap[r.Table]
	if !ok {
		return nil, errors.Errorf("unkown channel '%s'", r.Table)
	}

	return cb(r.Table, r.Action, r.Data)
}

func ParseTime(timestamp string) (time.Time, error) {
	t, err := time.Parse(timeFMT, timestamp)
	return t, err
}

func SubscribeCB(channel string, cb ResponseParseCB) {
	rcbMap[channel] = cb
}
