package jex

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/NadiaSama/ccexgo/misc/tconv"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	OrderID int

	OrderResult struct {
		Symbol              string `json:"symbol"`
		OrderID             int    `json:"orderId"`
		TranscatTime        int64  `json:"transcatTime"`
		Price               string `json:"price"`
		OrigQty             string `json:"origQty"`
		ExecutedQty         string `json:"executedQty"`
		CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
		Status              string `json:"status"`
		TimeInForce         string `json:"timeInForce"`
		Type                string `json:"type"`
		Side                string `json:"side"`
		Time                int64  `json:"time"`
		UpdateTime          int64  `json:"updateTime"`
	}
)

const (
	newOrderRespTypeResult = "RESULT"
	orderURI               = "/api/v1/option/order"
)

var (
	//set in init
	side2Str = map[exchange.OrderSide]string{}
	str2Side = map[string]exchange.OrderSide{}
	type2Str = map[exchange.OrderType]string{}
	str2Type = map[string]exchange.OrderType{}

	str2Status = map[string]exchange.OrderStatus{
		"NEW":              exchange.OrderStatusOpen,
		"PARTIALLY_FILLED": exchange.OrderStatusOpen,
		"FILLED":           exchange.OrderStatusDone,
		"CANCELED":         exchange.OrderStatusCancel,
		"CANCLEFILLED":     exchange.OrderStatusCancel,
		"FAIL":             exchange.OrderStatusFailed,
	}
)

func init() {
	sides := []struct {
		Side exchange.OrderSide
		Str  string
	}{
		{exchange.OrderSideBuy, "BUY"},
		{exchange.OrderSideSell, "SELL"},
	}
	for _, e := range sides {
		side2Str[e.Side] = e.Str
		str2Side[e.Str] = e.Side
	}

	types := []struct {
		Type exchange.OrderType
		Str  string
	}{
		{exchange.OrderTypeLimit, "LIMIT"},
		{exchange.OrderTypeMarket, "MARKET"},
		{exchange.OrderTypeStopLimit, "STOPLIMIT"},
		{exchange.OrderTypeStopMarket, "STOPMARKET"},
	}
	for _, e := range types {
		type2Str[e.Type] = e.Str
		str2Type[e.Str] = e.Type
	}
}

func (c *Client) OptionCreateOrder(ctx context.Context, req *exchange.OrderRequest) (*exchange.Order, error) {
	params := map[string]string{
		"price":            fmt.Sprintf("%f", req.Price),
		"quantity":         fmt.Sprintf("%f", req.Amount),
		"symbol":           req.Symbol.String(),
		"side":             side2Str[req.Side],
		"type":             type2Str[req.Type],
		"newOrderRespType": newOrderRespTypeResult,
	}
	var or OrderResult
	if err := c.request(ctx, "POST", orderURI, params, &or, true); err != nil {
		return nil, err
	}
	return or.Transform()
}

func (c *Client) OptionCancelOrder(ctx context.Context, order *exchange.Order) (*exchange.Order, error) {
	params := map[string]string{
		"symbol":  order.Symbol.String(),
		"orderId": order.ID.String(),
	}

	var or OrderResult
	if err := c.request(ctx, "DELETE", orderURI, params, &or, true); err != nil {
		return nil, err
	}

	result, err := or.Transform()
	if err != nil {
		return nil, err
	}
	if result.Status != exchange.OrderStatusCancel {
		return nil, errors.Errorf("cancel order fail")
	}
	return result, nil
}

func (or *OrderResult) Transform() (*exchange.Order, error) {
	side, ok := str2Side[or.Side]
	if !ok {
		return nil, errors.Errorf(`bad order side: "%s"`, or.Side)
	}
	typ, ok := str2Type[or.Type]
	if !ok {
		return nil, errors.Errorf(`bad order type "%s"`, or.Type)
	}
	sym, err := ParseSymbol(or.Symbol)
	if err != nil {
		return nil, errors.WithMessagef(err, "parse order fail")
	}
	status, ok := str2Status[or.Status]
	if !ok {
		return nil, errors.Errorf(`bad order status "%s"`, or.Status)
	}
	price, err := strconv.ParseFloat(or.Price, 64)
	if err != nil {
		return nil, errors.Errorf("bad order price %s", or.Price)
	}
	amount, err := strconv.ParseFloat(or.OrigQty, 64)
	if err != nil {
		return nil, errors.Errorf("bad order amount %s", or.OrigQty)
	}
	exec, err := strconv.ParseFloat(or.ExecutedQty, 64)
	if err != nil {
		return nil, errors.Errorf("bad executed amount %s", or.ExecutedQty)
	}

	var updated time.Time
	var created time.Time
	if or.TranscatTime != 0 {
		created = tconv.Milli2Time(or.TranscatTime)
	} else {
		created = tconv.Milli2Time(or.Time)
	}
	if or.UpdateTime != 0 {
		updated = tconv.Milli2Time(or.UpdateTime)
	} else {
		updated = time.Now()
	}

	return &exchange.Order{
		ID:      NewOrderID(or.OrderID),
		Side:    side,
		Status:  status,
		Symbol:  sym,
		Price:   decimal.NewFromFloat(price),
		Amount:  decimal.NewFromFloat(amount),
		Filled:  decimal.NewFromFloat(exec),
		Created: created,
		Updated: updated,
		Type:    typ,
		Raw:     or,
	}, nil
}

func NewOrderID(id int) OrderID {
	return OrderID(id)
}

func (id OrderID) String() string {
	return fmt.Sprintf("%d", id)
}
