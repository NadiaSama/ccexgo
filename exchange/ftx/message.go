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
		codeMap   map[string]exchange.Symbol
		orderBook map[string]*OrderBook
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
	typePartial      = "partial"
	typeUpdate       = "update"

	codeReconnet = 20001

	channelOrderBook = "orderbook"
	channelOrders    = "orders"
	channelFills     = "fills"
)

func NewCodeC(codeMap map[string]exchange.Symbol) *CodeC {
	return &CodeC{
		exchange.NewCodeC(),
		codeMap,
		make(map[string]*OrderBook),
	}
}

func (cc *CodeC) Decode(raw []byte) (rpc.Response, error) {
	var cr callResponse
	if err := json.Unmarshal(raw, &cr); err != nil {
		return nil, err
	}

	id := subID(cr.Channel, cr.Market)

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

	case typePartial:
		if cr.Channel == channelOrderBook {
			sym, ok := cc.codeMap[cr.Market]
			if !ok {
				return nil, errors.Errorf("unknow market '%s'", cr.Market)
			}
			ob := NewOrderBook(sym)
			notify, err := ob.Init(&cr)
			if err != nil {
				return nil, err
			}
			cc.orderBook[cr.Market] = ob

			return &rpc.Notify{
				Method: id,
				Params: notify,
			}, nil
		}
		return nil, errors.Errorf("unsupport partial data %s %s", cr.Channel, cr.Market)

	case typeUpdate:
		var param interface{}
		switch cr.Channel {
		case channelOrders:
			o, err := cc.parseOrder(cr.Data)
			if err != nil {
				return nil, err
			}
			param = o

		case channelFills:
			f, err := cc.parseFills(cr.Data)
			if err != nil {
				return nil, err
			}
			param = f

		case channelOrderBook:
			ob, ok := cc.orderBook[cr.Market]
			if !ok {
				return nil, errors.Errorf("unkown market '%s'", cr.Market)
			}
			f, err := ob.Update(&cr)
			if err != nil {
				return nil, err
			}
			param = f

		default:
			return nil, errors.Errorf("unsupport channel '%s'", cr.Channel)
		}
		ret := &rpc.Notify{
			Method: id,
			Params: param,
		}
		return ret, nil

	default:
		return nil, errors.Errorf("unsupport type '%s'", cr.Type)
	}
}

func (cc *CodeC) parseOrder(raw []byte) (*exchange.Order, error) {
	var order Order
	if err := json.Unmarshal(raw, &order); err != nil {
		return nil, err
	}
	return parseOrderInternal(&order, cc.codeMap)
}

func (cc *CodeC) parseFills(raw []byte) (*Fill, error) {
	var fill FillNotify
	if err := json.Unmarshal(raw, &fill); err != nil {
		return nil, err
	}

	return parseFillInternal(&fill, cc.codeMap)
}

func subID(channel string, market string) string {
	if len(market) == 0 {
		return channel
	}

	return fmt.Sprintf("%s.%s", channel, market)
}
