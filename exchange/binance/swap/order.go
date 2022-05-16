package swap

import (
	"context"
	"net/http"

	"github.com/NadiaSama/ccexgo/exchange/binance"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	AddOrderReq struct {
		*binance.RestReq
	}

	OrderResp struct {
		binance.APIError                 //in case of error
		ClientOrderID    string          `json:"newOrderId"`
		CumQty           decimal.Decimal `json:"cumQty"`
		CumQuote         decimal.Decimal `json:"cumQuote"`
		ExecutedQty      decimal.Decimal `json:"executedQty"`
		OrderID          int64           `json:"orderId"`
		AvgPrice         decimal.Decimal `json:"avgPrice"`
		OrigQty          decimal.Decimal `json:"origQty"`
		Price            decimal.Decimal `json:"price"`
		ReduceOnly       bool            `json:"reduceOnly"`
		Side             string          `json:"side"`
		PositionSide     string          `json:"positionSide"`
		Status           string          `json:"status"`
		StopPrice        decimal.Decimal `json:"stopPrice"`
		ClosePosition    bool            `json:"closePosition"`
		Symbol           string          `json:"symbol"`
		TimeInForce      string          `json:"timeInForce"`
		Type             string          `json:"type"`
		OrigType         string          `json:"origType"`
		ActivatePrice    decimal.Decimal `json:"activatePrice"`
		PriceRate        decimal.Decimal `json:"priceRate"`
		UpdateTime       int64           `json:"updateTime"`
		WorkingType      string          `json:"workingType"`
		PriceProtect     bool            `json:"priceProtect"`
	}

	OrderReq struct {
		*binance.RestReq
	}
)

const (
	OrderEndPoint = "/fapi/v1/order"
)

//NewAddOrderReq according symbol, side, type
func NewAddOrderReq(symbol string, side string, typ string) *AddOrderReq {
	req := binance.NewRestReq()
	req.AddFields("symbol", symbol)
	req.AddFields("side", side)
	req.AddFields("type", typ)

	return &AddOrderReq{
		RestReq: req,
	}
}

func (req *AddOrderReq) TimeInForce(tif string) *AddOrderReq {
	req.AddFields("timeInForce", tif)
	return req
}

func (req *AddOrderReq) PositionSide(side string) *AddOrderReq {
	req.AddFields("positionSide", side)
	return req
}

func (req *AddOrderReq) Price(prc decimal.Decimal) *AddOrderReq {
	req.AddFields("price", prc.String())
	return req
}

func (req *AddOrderReq) Quantity(q decimal.Decimal) *AddOrderReq {
	req.AddFields("quantity", q.String())
	return req
}

func (cl *RestClient) AddOrder(ctx context.Context, req *AddOrderReq) (*OrderResp, error) {
	values, err := req.Values()
	if err != nil {
		return nil, errors.WithMessage(err, "get param fail")
	}

	var ret OrderResp
	if err := cl.Request(ctx, http.MethodPost, OrderEndPoint, values, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "add order fail")
	}

	return &ret, nil
}

func NewOrderReq(symbol string) *OrderReq {
	req := binance.NewRestReq()
	req.AddFields("symbol", symbol)
	return &OrderReq{
		RestReq: req,
	}
}

func (r *OrderReq) OrderID(id int64) *OrderReq {
	r.AddFields("orderId", id)
	return r
}

func (cl *RestClient) GetOrder(ctx context.Context, req *OrderReq) (*OrderResp, error) {
	var ret OrderResp
	if err := cl.GetRequest(ctx, OrderEndPoint, req, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "get order req fail")
	}

	return &ret, nil
}

func (cl *RestClient) CancelOrder(ctx context.Context, req *OrderReq) (*OrderResp, error) {
	var ret OrderResp
	values, err := req.Values()
	if err != nil {
		return nil, errors.WithMessage(err, "get req fail")
	}

	if err := cl.Request(ctx, http.MethodDelete, OrderEndPoint, values, nil, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "cancel order fail")
	}

	return &ret, nil
}
