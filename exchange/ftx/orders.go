package ftx

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/NadiaSama/ccexgo/exchange"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	Order struct {
		CreatedAt     string          `json:"createdAt"`
		FilledSize    decimal.Decimal `json:"filledSize"`
		Future        string          `json:"future"`
		ID            int64           `json:"id"`
		Market        string          `json:"market"`
		Price         decimal.Decimal `json:"price"`
		AvgFillPrice  decimal.Decimal `json:"avgFillPrice"`
		RemainingSize decimal.Decimal `json:"remainingSize"`
		Side          string          `json:"side"`
		Size          decimal.Decimal `json:"size"`
		Status        string          `json:"status"`
		Type          string          `json:"type"`
		ReduceOnly    bool            `json:"reduceOnly"`
		IOC           bool            `json:"ioc"`
		PostOnly      bool            `json:"postOnly"`
		ClientID      string          `json:"clientId"`
	}

	OrderReq struct {
		Market string  `json:"market"`
		Side   string  `json:"side"`
		Price  float64 `json:"price"`
		Type   string  `json:"type"`
		Size   float64 `json:"size"`
	}
)

const (
	ftxOrderNew   = "new"
	ftxOrderOpen  = "open"
	ftxOrderClose = "closed"

	orderEndPoint = "/orders"
)

var (
	typeMap map[string]exchange.OrderType = map[string]exchange.OrderType{
		"limit":  exchange.OrderTypeLimit,
		"market": exchange.OrderTypeMarket,
	}

	typeRMap map[exchange.OrderType]string = map[exchange.OrderType]string{
		exchange.OrderTypeLimit:  "limit",
		exchange.OrderTypeMarket: "market",
	}

	sideMap map[string]exchange.OrderSide = map[string]exchange.OrderSide{
		"buy":  exchange.OrderSideBuy,
		"sell": exchange.OrderSideSell,
	}

	sideRMap map[exchange.OrderSide]string = map[exchange.OrderSide]string{
		exchange.OrderSideBuy:  "buy",
		exchange.OrderSideSell: "sell",
	}
)

func (rc *RestClient) OrderNew(ctx context.Context, req *exchange.OrderRequest, options ...exchange.OrderReqOption) (*exchange.Order, error) {
	side, ok := sideRMap[req.Side]
	if !ok {
		return nil, errors.Errorf("unkown orderside '%d'", req.Side)
	}
	typ, ok := typeRMap[req.Type]
	if !ok {
		return nil, errors.Errorf("unkown order type '%d'", req.Type)
	}

	p, _ := req.Price.Float64()
	a, _ := req.Amount.Float64()
	or := OrderReq{
		Market: req.Symbol.String(),
		Price:  p,
		Size:   a,
		Side:   side,
		Type:   typ,
	}
	b, _ := json.Marshal(or)
	buf := bytes.NewBuffer(b)

	var o Order
	if err := rc.request(ctx, http.MethodPost, orderEndPoint, nil, buf, true, &o); err != nil {
		return nil, err
	}
	return rc.parseOrder(&o)
}

func (rc *RestClient) parseOrder(o *Order) (*exchange.Order, error) {
	ct, err := time.Parse("2006-01-02T15:04:05.000000Z07:00", o.CreatedAt)
	if err != nil {
		return nil, errors.WithMessagef(err, "bad create time '%s'", o.CreatedAt)
	}
	var os exchange.OrderStatus
	if o.Status == ftxOrderNew || o.Status == ftxOrderOpen {
		os = exchange.OrderStatusOpen
	} else {
		if o.FilledSize == o.Size {
			os = exchange.OrderStatusDone
		} else {
			os = exchange.OrderStatusCancel
		}
	}

	symbol, ok := rc.symbols[o.Market]
	if !ok {
		return nil, errors.WithMessagef(err, "parse symbol '%s' fail", o.Market)
	}

	typ, ok := typeMap[o.Type]
	if !ok {
		return nil, errors.Errorf("unkown order type '%s'", o.Type)
	}

	side, ok := sideMap[o.Side]
	if !ok {
		return nil, errors.Errorf("unkown order side '%s'", o.Side)
	}

	order := &exchange.Order{
		ID:       exchange.NewIntID(o.ID),
		Symbol:   symbol,
		Amount:   o.Size,
		Filled:   o.FilledSize,
		Price:    o.Price,
		AvgPrice: o.AvgFillPrice,
		Created:  ct,
		Updated:  ct,
		Status:   os,
		Side:     side,
		Type:     typ,
	}
	return order, nil
}
