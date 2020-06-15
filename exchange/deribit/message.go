package deribit

import (
	"encoding/json"
	"strings"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/internal/rpc"
	"github.com/pkg/errors"
)

//deribit message json serizlize

const (
	JsonRPCVersion     = "2.0"
	subscriptionMethod = "subscription"
)

type (
	BookData struct {
		Timestamp      int              `json:"timestamp"`
		InstrumentName string           `json:"instrument_name"`
		ChangeID       int              `json:"charge_id"`
		Bids           [][3]interface{} `json:"bids"`
		Asks           [][3]interface{} `json:"asks"`
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

	Codec struct {
	}
)

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

func parseNotify(resp *Response) (*rpc.Notify, error) {
	fields := strings.Split(resp.Params.Channel, ".")
	if len(fields) == 0 {
		return nil, errors.Errorf("bad message %v", resp.Params)
	}
	//book.${instrument_name}.${gap}
	if fields[0] == "book" && len(fields) == 3 {
		return parseNotifyBook(resp.Params.Data, fields[1])
	}
	return nil, errors.Errorf("unsupport channel %s", resp.Params.Channel)

}

func parseNotifyBook(data json.RawMessage, instrument string) (*rpc.Notify, error) {
	processArr := func(d []exchange.OrderElem, s [][3]interface{}) (ret error) {
		defer func() {
			if err := recover(); err != nil {
				ret = err.(error)
			}
		}()

		for i, v := range s {
			op := v[0].(string)
			price := v[1].(float64)
			amount := v[2].(float64)

			if op == "new" || op == "change" {
				d[i].Amount = amount
				d[i].Price = price
			} else if op == "delete" {
				d[i].Amount = 0
				d[i].Price = price
			} else {
				ret = errors.Errorf("unkown op %s", op)
				return
			}
		}
		return
	}
	var bn BookData
	if err := json.Unmarshal(data, &bn); err != nil {
		return nil, err
	}
	notify := &rpc.Notify{
		Method: subscriptionMethod,
	}
	on := &exchange.OrderBookNotify{
		Symbol: instrument,
		Asks:   make([]exchange.OrderElem, len(bn.Asks)),
		Bids:   make([]exchange.OrderElem, len(bn.Bids)),
	}

	if err := processArr(on.Asks, bn.Asks); err != nil {
		return nil, err
	}
	if err := processArr(on.Bids, bn.Bids); err != nil {
		return nil, err
	}
	notify.Params = on
	return notify, nil
}
