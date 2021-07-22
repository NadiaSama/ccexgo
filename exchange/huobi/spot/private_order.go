package spot

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/exchange/huobi"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	OrdersChannel struct {
		symbol string
	}

	OrderData struct {
		EventType       string `json:"eventType"`
		Symbol          string `json:"symbol"`
		AccountID       int64  `json:"accountId"`
		OrderID         int64  `json:"orderId"`
		ClientOrderID   string `json:"clientOrderId"`
		OrderStatus     string `json:"orderStatus"`
		OrderPrice      string `json:"orderPrice"`
		OrderSize       string `json:"orderSize"`
		OrderValue      string `json:"orderValue"`
		Type            string `json:"type"`
		OrderCreateTime int64  `json:"orderCreateTime"`
		OrderSource     string `json:"orderSource"`
		TradePrice      string `json:"tradePrice"`
		TradeVolume     string `json:"tradeVolume"`
		TradeID         int64  `json:"tradeId"`
		TradeTime       int64  `json:"tradeTime"`
		Aggressor       bool   `json:"aggressor"`
		RemainAmt       string `json:"remainAmt"`
		ExecAmt         string `json:"execAmt"`
		LastActTime     int64  `json:"lastActTime"`
	}
)

var (
	orderStatusMap = map[string]exchange.OrderStatus{
		"created":          exchange.OrderStatusOpen,
		"submitted":        exchange.OrderStatusOpen,
		"partail-filled":   exchange.OrderStatusOpen,
		"partial-canceled": exchange.OrderStatusOpen,
		"canceling":        exchange.OrderStatusOpen,
		"filled":           exchange.OrderStatusDone,
		"canceled":         exchange.OrderStatusCancel,
	}
)

func NewOrdersChannel(sym string) *OrdersChannel {
	return &OrdersChannel{
		symbol: sym,
	}
}

func (oc *OrdersChannel) String() string {
	return fmt.Sprintf("orders#%s", oc.symbol)
}

func ParseOrder(data json.RawMessage) (interface{}, error) {
	var od OrderData
	if err := json.Unmarshal(data, &od); err != nil {
		return nil, err
	}

	var (
		order *exchange.Order
		err   error
	)
	switch od.EventType {
	case "creation":
		fallthrough
	case "trade":
		fallthrough
	case "cancellation":
		order, err = od.Parse()
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.Errorf("unkown eventType '%s'", od.EventType)
	}

	return order, nil
}

func (od *OrderData) Parse() (*exchange.Order, error) {
	symbol, err := ParseSymbol(od.Symbol)
	if err != nil {
		return nil, err
	}
	status, err := ParseOrderStatus(od.OrderStatus)
	if err != nil {
		return nil, err
	}
	side, typ, err := ParseOrderType(od.Type)
	if err != nil {
		return nil, err
	}

	ret := exchange.Order{
		ID:     exchange.NewIntID(od.OrderID),
		Symbol: symbol,
		Side:   side,
		Type:   typ,
		Status: status,
	}
	if od.TradePrice != "" {
		prc, err := decimal.NewFromString(od.TradePrice)
		if err != nil {
			return nil, errors.WithMessage(err, "invalid tradePrice")
		}
		ret.Price = prc
	}

	if od.TradeVolume != "" {
		vol, err := decimal.NewFromString(od.TradeVolume)
		if err != nil {
			return nil, errors.WithMessage(err, "invalid tradeVolume")
		}
		ret.Amount = vol
	}

	if od.OrderCreateTime != 0 {
		ret.Created = huobi.ParseTS(od.OrderCreateTime)
		ret.Updated = ret.Created
	}

	var ts int64
	if od.LastActTime != 0 {
		ts = od.LastActTime
	} else if od.TradeTime != 0 {
		ts = od.TradeTime
	}
	if ts != 0 {
		ret.Updated = huobi.ParseTS(ts)
	}
	ret.Raw = od
	return &ret, nil
}

func ParseOrderStatus(state string) (exchange.OrderStatus, error) {
	if status, ok := orderStatusMap[state]; !ok {
		return status, errors.Errorf("unkown order status '%s'", state)
	} else {
		return status, nil
	}
}

func ParseOrderType(oType string) (exchange.OrderSide, exchange.OrderType, error) {
	fields := strings.SplitN(oType, "-", 2)

	var side exchange.OrderSide
	var typ exchange.OrderType
	if fields[0] == "buy" {
		side = exchange.OrderSideBuy
	} else if fields[0] == "sell" {
		side = exchange.OrderSideSell
	} else {
		return side, typ, errors.Errorf("parse order side fail unkown order type '%s'", oType)
	}

	if strings.HasPrefix(fields[1], "limit") {
		typ = exchange.OrderTypeLimit
	} else if strings.HasPrefix(fields[1], "market") {
		typ = exchange.OrderTypeMarket
	} else {
		return side, typ, errors.Errorf("parse order type fail unkown order type '%s'", oType)
	}

	return side, typ, nil
}
