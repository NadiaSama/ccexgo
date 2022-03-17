package option

import (
	"context"
	"net/http"

	"github.com/NadiaSama/ccexgo/exchange/binance"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type (
	GetOrderReq struct {
		*binance.RestReq
	}

	PostOrdreReq struct {
		*binance.RestReq
	}

	OrderResp struct {
		ID            string `json:"id"`
		Symbol        string `json:"symbol"`
		Price         string `json:"price"`
		Quantity      string `json:"quantity"`
		ExecutedQty   string `json:"executedQty"`
		Fee           string `json:"fee"`
		Side          string `json:"side"`
		Type          string `json:"type"`
		TimeInForce   string `json:"timeInForce"`
		CreateDate    int64  `json:"createDate"`
		Status        string `json:"status"`
		AvgPrice      string `json:"avgPrice"`
		Source        string `json:"source"`
		ReduceOnly    bool   `json:"reduceOnly"`
		ClientOrderID string `json:"clientOrderId"`
	}
)

const (
	OrderEndPoint  = "/vapi/v1/order"
	OrderSideBuy   = "BUY"
	OrderSideSell  = "SELL"
	OrderTypeLimit = "LIMIT"
)

//NewPostOrderReq build create order request, the maount and price param will be formatted according to symbol precision
func NewPostOrderReq(symbol string, side string, typ string, amount float64, price float64) (*PostOrdreReq, error) {
	s, err := ParseSymbol(symbol)
	if err != nil {
		return nil, errors.WithMessage(err, "parse symbol fail")
	}

	at := decimal.NewFromFloat(amount)
	pr := decimal.NewFromFloat(price)

	req := binance.NewRestReq()
	req.AddFields("symbol", symbol)
	req.AddFields("side", side)
	req.AddFields("type", typ)
	req.AddFields("quantity", at.StringFixed(int32(s.AmountExponent())))
	req.AddFields("price", pr.StringFixed(int32(s.PriceExponent())))

	return &PostOrdreReq{
		RestReq: req,
	}, nil
}

func (or *PostOrdreReq) ReduceOnly(r bool) *PostOrdreReq {
	or.AddFields("reduceOnly", r)
	return or
}

//NewGetOrderReq build get and cancel order request
func NewGetOrderReq(symbol string) *GetOrderReq {
	req := binance.NewRestReq()
	req.AddFields("symbol", symbol)

	return &GetOrderReq{
		RestReq: req,
	}
}

func (gor *GetOrderReq) OrderID(orderId string) *GetOrderReq {
	gor.AddFields("orderId", orderId)
	return gor
}

func (gor *GetOrderReq) ClientOrderID(cid string) *GetOrderReq {
	gor.AddFields("clientOrderId", cid)
	return gor
}

func (rc *RestClient) PostOrder(ctx context.Context, req *PostOrdreReq) (*OrderResp, error) {
	var ret OrderResp

	resp := RestResp{
		Data: &ret,
	}

	values, err := req.Values()
	if err != nil {
		return nil, errors.WithMessage(err, "get values fail")
	}

	if err := rc.Request(ctx, http.MethodPost, OrderEndPoint, values, nil, true, &resp); err != nil {
		return nil, errors.WithMessage(err, "request fail")
	}

	if resp.Code != 0 {
		return nil, errors.Errorf("invalid resp code=%d msg=%s", resp.Code, resp.Msg)
	}

	return &ret, nil
}

func (rc *RestClient) GetOrder(ctx context.Context, req *GetOrderReq) (*OrderResp, error) {
	var ret OrderResp

	if err := rc.GetRequest(ctx, OrderEndPoint, req, true, &ret); err != nil {
		return nil, errors.WithMessage(err, "get request fail")
	}

	return &ret, nil
}

func (rc *RestClient) DeleteOrder(ctx context.Context, req *GetOrderReq) (*OrderResp, error) {
	var ret OrderResp
	resp := RestResp{
		Data: &ret,
	}

	values, err := req.Values()
	if err != nil {
		return nil, errors.WithMessage(err, "get values fail")
	}

	if err := rc.Request(ctx, http.MethodDelete, OrderEndPoint, values, nil, true, &resp); err != nil {
		return nil, errors.WithMessage(err, "request fail")
	}

	if resp.Code != 0 {
		return nil, errors.Errorf("invalid resp code=%d msg=%s", resp.Code, resp.Msg)
	}

	return &ret, nil
}
