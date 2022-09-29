package spot

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

type (
	// CodeC used to decode binance websocket notify message to coresponding struct
	// and encode sbuscribe request
	CodeC struct {
		lastID string
	}

	// TickerNotify binance spot ticker notify
	TickerNotify struct {
		UpdateID   int64  `json:"u"`
		Symbol     string `json:"s"`
		Bid1Price  string `json:"b"`
		Bid1Amount string `json:"B"`
		Ask1Price  string `json:"a"`
		Ask1Amount string `json:"A"`
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
	MethodSubscibe = "SUBSCRIBE"
)

func NewCodeC() *CodeC {
	return &CodeC{}
}

// Decode binance websocket notify message
func (cc *CodeC) Decode(raw []byte) (rpc.Response, error) {
	var tn TickerNotify

	g := gjson.ParseBytes(raw)

	// by now only handle subscribe response which result is nil
	if g.Get("id").Exists() && g.Get("result").Exists() {
		return &rpc.Result{
			ID: g.Get("id").String(),
		}, nil
	}
	if err := json.Unmarshal(raw, &tn); err != nil {
		return nil, errors.WithMessage(err, "unmarshal json fail")
	}
	return &rpc.Notify{Params: &tn, Method: "ticker"}, nil
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
		Method: MethodSubscibe,
	}

	cc.lastID = req.ID()

	return json.Marshal(&sub)
}
